#!/bin/bash
set -e

SGX_MODE=${SGX_MODE:-HW}

if [ ! -n "$1" ] ;then
echo "error: missing input parameter, please input image tag, such as
enclave-cc-boot-instance:v1.0."
exit 1
fi

#if image $1 exist, remove it.
docker rmi $1 -f

docker build --build-arg SGX_MODE=${SGX_MODE} . -t $1

docker export $(docker create $1) | sudo tee > ${PAYLOAD_ARTIFACTS}/boot-instance.tar
