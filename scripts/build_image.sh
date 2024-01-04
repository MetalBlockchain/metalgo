#!/usr/bin/env bash

set -euo pipefail

# Metalgo root folder
METAL_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd .. && pwd )
# Load the constants
source "$METAL_PATH"/scripts/constants.sh

# WARNING: this will use the most recent commit even if there are un-committed changes present
full_commit_hash="$(git --git-dir="$METAL_PATH/.git" rev-parse HEAD)"
commit_hash="${full_commit_hash::8}"

echo "Building Docker Image with tags: $metalgo_dockerhub_repo:$commit_hash , $metalgo_dockerhub_repo:$current_branch"
docker build -t "$metalgo_dockerhub_repo:$commit_hash" \
        -t "$metalgo_dockerhub_repo:$current_branch" "$METAL_PATH" -f "$METAL_PATH/Dockerfile"
