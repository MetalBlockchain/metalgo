#!/usr/bin/env bash

set -euo pipefail

print_usage() {
  printf "Usage: build_metal [OPTIONS]

  Build metalgo

  Options:

    -r  Build with race detector
"
}

race=''
while getopts 'r' flag; do
  case "${flag}" in
    r) race='-race' ;;
    *) print_usage
      exit 1 ;;
  esac
done

# MetalGo root folder
METAL_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd .. && pwd )
# Load the constants
source "$METAL_PATH"/scripts/constants.sh

build_args="$race"
echo "Building MetalGo..."
go build $build_args -ldflags "-X github.com/MetalBlockchain/metalgo/version.GitCommit=$git_commit $static_ld_flags" -o "$metalgo_path" "$METAL_PATH/main/"*.go
