#!/bin/bash
set -e

CI=${CI:-no}
PUSH=${PUSH:-no}
SGX_MODE=${SGX_MODE:-HW}
if [ "${CI}" == "yes" ]; then
	DEFAULT_IMAGE=quay.io/confidential-containers/runtime-payload-ci:enclave-cc-${SGX_MODE}-$(git rev-parse HEAD)
	DEFAULT_LATEST_IMAGE=quay.io/confidential-containers/runtime-payload-ci:enclave-cc-${SGX_MODE}-latest
else
	DEFAULT_IMAGE=quay.io/confidential-containers/runtime-payload:enclave-cc-${SGX_MODE}-$(git describe --tags --abbrev=0)
	DEFAULT_LATEST_IMAGE=quay.io/confidential-containers/runtime-payload:enclave-cc-${SGX_MODE}-latest
fi
IMAGE=${IMAGE:-${DEFAULT_IMAGE}}

docker rmi ${IMAGE} -f

export SCRIPT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
export ENCLAVE_CC_ROOT="${SCRIPT_ROOT}/../../../"

export PAYLOAD_ARTIFACTS="${SCRIPT_ROOT}/payload_artifacts"
mkdir -p ${PAYLOAD_ARTIFACTS}

# build pre-installed OCI bundle for agent enclave container
pushd ${SCRIPT_ROOT}/agent-enclave-bundle
./gen_bundle.sh agent_enclave_container
popd

# build pre-installed OCI bundle for boot instance
pushd ${SCRIPT_ROOT}/boot-instance-bundle
./gen_bundle.sh app_enclave_container
popd

# build shim-rune binary: "containerd-shim-rune-v2"
pushd ${ENCLAVE_CC_ROOT}/src/shim
make binaries
cp ./bin/containerd-shim-rune-v2 ${PAYLOAD_ARTIFACTS}
# prepare shim-rune configuration.
cp ./config/config.toml ${PAYLOAD_ARTIFACTS}/shim-rune-config.toml
popd

cp ${SCRIPT_ROOT}/../deploy/enclave-cc-deploy.sh ${PAYLOAD_ARTIFACTS}/scripts

# prepare payload artifacts static tarball
pushd $PAYLOAD_ARTIFACTS
tar cfJ enclave-cc-static.tar.xz *
cp ${SCRIPT_ROOT}/Dockerfile .
docker build . -t ${IMAGE} -t ${DEFAULT_LATEST_IMAGE}
if [ "${PUSH}" == "yes" ]; then
	docker push ${IMAGE}
	docker push ${DEFAULT_LATEST_IMAGE}
fi
popd

#cleanup
rm -rf payload_artifacts
