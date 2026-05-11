#!/usr/bin/env bash

set -euo pipefail

# MetalGo root folder
METAL_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd .. && pwd )
# Load the constants
source "$METAL_PATH"/scripts/constants.sh
source "$METAL_PATH"/scripts/git_commit.sh

echo "Building tmpnetctl..."
go build -ldflags\
   "-X github.com/MetalBlockchain/metalgo/version.GitCommit=$git_commit $static_ld_flags"\
   -o "$METAL_PATH/build/tmpnetctl"\
   "$METAL_PATH/tests/fixture/tmpnet/tmpnetctl/"*.go
