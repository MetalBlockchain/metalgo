#!/usr/bin/env bash

set -euo pipefail

# Metalgo root folder
METAL_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd .. && pwd )
# Load the constants
source "$METAL_PATH"/scripts/constants.sh

if [[ $current_branch == *"-race" ]]; then
  echo "Branch name must not end in '-race'"
  exit 1
fi

# WARNING: this will use the most recent commit even if there are un-committed changes present
full_commit_hash="$(git --git-dir="$METAL_PATH/.git" rev-parse HEAD)"
commit_hash="${full_commit_hash::8}"

echo "Building Docker Image with tags: $metalgo_dockerhub_repo:$commit_hash , $metalgo_dockerhub_repo:$current_branch"
docker build -t "$metalgo_dockerhub_repo:$commit_hash" \
        -t "$metalgo_dockerhub_repo:$current_branch" "$METAL_PATH" -f "$METAL_PATH/Dockerfile"

echo "Building Docker Image with tags: $metalgo_dockerhub_repo:$commit_hash-race , $metalgo_dockerhub_repo:$current_branch-race"
docker build --build-arg="RACE_FLAG=-r" -t "$metalgo_dockerhub_repo:$commit_hash-race" \
        -t "$metalgo_dockerhub_repo:$current_branch-race" "$METAL_PATH" -f "$METAL_PATH/Dockerfile"
