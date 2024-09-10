#!/bin/bash

# Sliver Implant Framework
# Copyright (C) 2019  Bishop Fox

set -e

# Creates the static go asset archives

GO_VER="1.22.5"
ZIG_VER="0.13.0"
SGN_VER="0.0.3"
BLOAT_FILES="AUTHORS CONTRIBUTORS PATENTS VERSION favicon.ico robots.txt SECURITY.md CONTRIBUTING.md LICENSE README.md ./doc ./test ./api ./misc"

REPO_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
OUTPUT_DIR=$REPO_DIR/server/assets/fs
mkdir -p $OUTPUT_DIR
WORK_DIR=`mktemp -d`

echo "-----------------------------------------------------------------"
echo "$WORK_DIR (Output: $OUTPUT_DIR)"
echo "-----------------------------------------------------------------"
cd $WORK_DIR

# --- Android (arm64) ---
curl --output go$GO_VER.android-arm64.tar.gz https://dl.google.com/go/go$GO_VER.linux-arm64.tar.gz
tar xvf go$GO_VER.android-arm64.tar.gz
cd go
rm -rf $BLOAT_FILES
rm -rf ./src
rm -f ./pkg/tool/linux_arm64/doc
rm -f ./pkg/tool/linux_arm64/tour
rm -f ./pkg/tool/linux_arm64/test2json
cd ..
zip -r android-go.zip ./go
mkdir -p $OUTPUT_DIR/android/arm64
cp -vv android-go.zip $OUTPUT_DIR/android/arm64/go.zip
rm -rf ./go
rm -f android-go.zip go$GO_VER.android-arm64.tar.gz

# --- Cleanup ---
echo -e "clean up: $WORK_DIR"
rm -rf $WORK_DIR
echo -e "\n[*] All done\n"
