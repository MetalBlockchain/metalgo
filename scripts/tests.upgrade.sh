#!/usr/bin/env bash

set -euo pipefail

# e.g.,
# ./scripts/tests.upgrade.sh                                               # Use default version
# ./scripts/tests.upgrade.sh 1.11.0                                        # Specify a version
# METALGO_PATH=./path/to/metalgo ./scripts/tests.upgrade.sh 1.11.0 # Customization of metalgo path
if ! [[ "$0" =~ scripts/tests.upgrade.sh ]]; then
  echo "must be run from repository root"
  exit 255
fi

# The MetalGo local network does not support long-lived
# backwards-compatible networks. When a breaking change is made to the
# local network, this flag must be updated to the last compatible
# version with the latest code.
#
# v1.11.3 fixes a regression in Coreth genesis for custom networks.
DEFAULT_VERSION="1.11.3"

VERSION="${1:-${DEFAULT_VERSION}}"
if [[ -z "${VERSION}" ]]; then
  echo "Missing version argument!"
  echo "Usage: ${0} [VERSION]" >>/dev/stderr
  exit 255
fi

METALGO_PATH="$(realpath "${METALGO_PATH:-./build/metalgo}")"

#################################
# download metalgo
# https://github.com/MetalBlockchain/metalgo/releases
GOARCH=$(go env GOARCH)
GOOS=$(go env GOOS)
DOWNLOAD_URL=https://github.com/MetalBlockchain/metalgo/releases/download/v${VERSION}/metalgo-linux-${GOARCH}-v${VERSION}.tar.gz
DOWNLOAD_PATH=/tmp/metalgo.tar.gz
if [[ ${GOOS} == "darwin" ]]; then
  DOWNLOAD_URL=https://github.com/MetalBlockchain/metalgo/releases/download/v${VERSION}/metalgo-macos-v${VERSION}.zip
  DOWNLOAD_PATH=/tmp/metalgo.zip
fi

rm -f ${DOWNLOAD_PATH}
rm -rf "/tmp/metalgo-v${VERSION}"
rm -rf /tmp/metalgo-build

echo "downloading metalgo ${VERSION} at ${DOWNLOAD_URL}"
curl -L "${DOWNLOAD_URL}" -o "${DOWNLOAD_PATH}"

echo "extracting downloaded metalgo"
if [[ ${GOOS} == "linux" ]]; then
  tar xzvf ${DOWNLOAD_PATH} -C /tmp
elif [[ ${GOOS} == "darwin" ]]; then
  unzip ${DOWNLOAD_PATH} -d /tmp/metalgo-build
  mv /tmp/metalgo-build/build "/tmp/metalgo-v${VERSION}"
fi
find "/tmp/metalgo-v${VERSION}"

# Sourcing constants.sh ensures that the necessary CGO flags are set to
# build the portable version of BLST. Without this, ginkgo may fail to
# build the test binary if run on a host (e.g. github worker) that lacks
# the instructions to build non-portable BLST.
source ./scripts/constants.sh

#################################
echo "building upgrade.test"
# to install the ginkgo binary (required for test build and run)
go install -v github.com/onsi/ginkgo/v2/ginkgo@v2.13.1
ACK_GINKGO_RC=true ginkgo build ./tests/upgrade
./tests/upgrade/upgrade.test --help

#################################
# By default, it runs all upgrade test cases!
echo "running upgrade tests against the local cluster with ${METALGO_PATH}"
./tests/upgrade/upgrade.test \
  --ginkgo.v \
  --metalgo-path="/tmp/metalgo-v${VERSION}/metalgo" \
  --metalgo-path-to-upgrade-to="${METALGO_PATH}"
