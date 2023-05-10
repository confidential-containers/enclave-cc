#!/bin/bash
set -e

CI=${CI:-no}
PUSH=${PUSH:-no}
SGX_MODE=${SGX_MODE:-HW}
KBC=${KBC:-cc-kbc}
GO_VERSION=${GO_VERSION:-1.19}
if [ "${CI}" == "yes" ]; then
	DEFAULT_IMAGE=quay.io/confidential-containers/runtime-payload-ci:enclave-cc-${SGX_MODE}-${KBC}-$(git rev-parse HEAD)
	DEFAULT_LATEST_IMAGE=quay.io/confidential-containers/runtime-payload-ci:enclave-cc-${SGX_MODE}-${KBC}-latest
else
	DEFAULT_IMAGE=quay.io/confidential-containers/runtime-payload:enclave-cc-${SGX_MODE}-${KBC}-$(git describe --tags --abbrev=0)
	DEFAULT_LATEST_IMAGE=quay.io/confidential-containers/runtime-payload:enclave-cc-${SGX_MODE}-${KBC}-latest
fi
IMAGE=${IMAGE:-${DEFAULT_IMAGE}}

export SCRIPT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
export ENCLAVE_CC_ROOT="${SCRIPT_ROOT}/../../../"

export PAYLOAD_ARTIFACTS="${SCRIPT_ROOT}/payload_artifacts"
mkdir -p ${PAYLOAD_ARTIFACTS}

# build pre-installed OCI bundle for agent enclave container
pushd ${SCRIPT_ROOT}/agent-enclave-bundle
docker build ${ENCLAVE_CC_ROOT} -f ${SCRIPT_ROOT}/agent-enclave-bundle/Dockerfile --build-arg SGX_MODE=${SGX_MODE} --build-arg KBC=${KBC} -t agent-instance
jq -a -f sgx-mode-config.filter config.json.template | tee ${PAYLOAD_ARTIFACTS}/config.json
docker export $(docker create agent-instance) | tee > ${PAYLOAD_ARTIFACTS}/agent-instance.tar
popd

# build pre-installed OCI bundle for boot instance
docker build ${ENCLAVE_CC_ROOT} -f ${SCRIPT_ROOT}/boot-instance-bundle/Dockerfile --build-arg SGX_MODE=${SGX_MODE} -t boot-instance
docker export $(docker create boot-instance) | tee > ${PAYLOAD_ARTIFACTS}/boot-instance.tar

# build shim-rune binary: "containerd-shim-rune-v2"
pushd ${ENCLAVE_CC_ROOT}/src/shim
docker run --pull always -t -v ${PWD}:/build --workdir /build golang:${GO_VERSION}-bullseye make binaries
cp ./bin/containerd-shim-rune-v2 ${PAYLOAD_ARTIFACTS}
# prepare shim-rune configuration.
cp ./config/config.toml ${PAYLOAD_ARTIFACTS}/shim-rune-config.toml
popd

install -D ${SCRIPT_ROOT}/../deploy/enclave-cc-deploy.sh ${PAYLOAD_ARTIFACTS}/scripts/enclave-cc-deploy.sh

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
docker rmi ${IMAGE} boot-instance agent-instance -f
rm -rf payload_artifacts
