#!/bin/bash
set -e

if [ ! -n "$1" ] ;then
    echo "error: missing input parameter, please input image tag, such as enclave-cc-enclave-agent-app:v1.0."
    exit 1
fi

#if image $1 exist, remove it.
docker rmi $1 -f

docker build . -t $1

pushd ${PAYLOAD_ARTIFACTS}

sudo cp ${SCRIPT_ROOT}/agent-enclave-bundle/config.json .

docker export $(docker create $1) -o /tmp/agent-instance.tar

sudo mv /tmp/agent-instance.tar .

popd
