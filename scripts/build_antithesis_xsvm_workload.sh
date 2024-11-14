#!/usr/bin/env bash

set -euo pipefail

# Directory above this script
METAL_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd .. && pwd )
# Load the constants
source "$METAL_PATH"/scripts/constants.sh

echo "Building Workload..."
go build -o "$METAL_PATH/build/antithesis-xsvm-workload" "$METAL_PATH/tests/antithesis/xsvm/"*.go
