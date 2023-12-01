#!/usr/bin/env bash
set -e
set -o nounset
set -o pipefail

# e.g.,
# ./scripts/build.sh
# ./scripts/tests.e2e.sh ./build/metalgo
if ! [[ "$0" =~ scripts/tests.e2e.sh ]]; then
  echo "must be run from repository root"
  exit 255
fi

METALGO_PATH="${1-}"
if [[ -z "${METALGO_PATH}" ]]; then
  echo "Missing METALGO_PATH argument!"
  echo "Usage: ${0} [METALGO_PATH]" >>/dev/stderr
  exit 255
fi

#################################
echo "installing metal-network-runner"
ANR_WORKDIR="/tmp"
./scripts/install_anr.sh

# Sourcing constants.sh ensures that the necessary CGO flags are set to
# build the portable version of BLST. Without this, ginkgo may fail to
# build the test binary if run on a host (e.g. github worker) that lacks
# the instructions to build non-portable BLST.
source ./scripts/constants.sh

#################################
echo "building e2e.test"
# to install the ginkgo binary (required for test build and run)
go install -v github.com/onsi/ginkgo/v2/ginkgo@v2.1.4
ACK_GINKGO_RC=true ginkgo build ./tests/e2e
./tests/e2e/e2e.test --help

#################################
# run "metal-network-runner" server
echo "launch metal-network-runner in the background"
$ANR_WORKDIR/metal-network-runner \
  server \
  --log-level debug \
  --port=":12342" \
  --disable-grpc-gateway &
PID=${!}

#################################
echo "running e2e tests against the local cluster with ${METALGO_PATH}"
./tests/e2e/e2e.test \
  --ginkgo.v \
  --log-level debug \
  --network-runner-grpc-endpoint="0.0.0.0:12342" \
  --network-runner-metalgo-path=${METALGO_PATH} \
  --network-runner-metalgo-log-level="WARN" \
  --test-keys-file=tests/test.insecure.secp256k1.keys &&
  EXIT_CODE=$? || EXIT_CODE=$?

kill ${PID}

if [[ ${EXIT_CODE} -gt 0 ]]; then
  echo "FAILURE with exit code ${EXIT_CODE}"
  exit ${EXIT_CODE}
else
  echo "ALL SUCCESS!"
fi
