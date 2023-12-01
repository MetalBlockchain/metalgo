#!/usr/bin/env bash

set -euo pipefail

# Avalanche root directory
AVALANCHE_PATH=$(
  cd "$(dirname "${BASH_SOURCE[0]}")"
  cd .. && pwd
)

#################################
# download avalanche-network-runner
# https://github.com/ava-labs/avalanche-network-runner
GOARCH=$(go env GOARCH)
GOOS=$(go env GOOS)
NETWORK_RUNNER_VERSION=1.6.2
anr_workdir=${ANR_WORKDIR:-"/tmp"}
DOWNLOAD_PATH=${anr_workdir}/metal-network-runner-v${NETWORK_RUNNER_VERSION}.tar.gz
DOWNLOAD_URL="https://github.com/MetalBlockchain/metal-network-runner/releases/download/v${NETWORK_RUNNER_VERSION}/metal-network-runner_${NETWORK_RUNNER_VERSION}_${GOOS}_${GOARCH}.tar.gz"
echo "Installing metal-network-runner ${NETWORK_RUNNER_VERSION} to ${anr_workdir}/metal-network-runner"

# download only if not already downloaded
if [ ! -f "$DOWNLOAD_PATH" ]; then
  echo "downloading metal-network-runner ${NETWORK_RUNNER_VERSION} at ${DOWNLOAD_URL} to ${DOWNLOAD_PATH}"
  curl --fail -L ${DOWNLOAD_URL} -o ${DOWNLOAD_PATH}
else
  echo "metal-network-runner ${NETWORK_RUNNER_VERSION} already downloaded at ${DOWNLOAD_PATH}"
fi

rm -f ${anr_workdir}/metal-network-runner

echo "extracting downloaded metal-network-runner"
tar xzvf ${DOWNLOAD_PATH} -C ${anr_workdir}
