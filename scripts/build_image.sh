#!/usr/bin/env bash

set -euo pipefail

# Metalgo root folder
METAL_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd .. && pwd )
# Load the constants
source "$METAL_PATH"/scripts/constants.sh

# build_image_from_remote.sh is deprecated
source "$METAL_PATH"/scripts/build_local_image.sh
