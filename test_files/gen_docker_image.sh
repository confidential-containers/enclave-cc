#!/bin/bash
set -e

pushd occlum_instance

cat >Dockerfile <<EOF
FROM ubuntu

RUN mkdir -p /run/rune
WORKDIR /run/rune

#RUN mkdir -p rootfs/lower && mkdir -p rootfs/upper

ADD occlum_instance.tar.gz /run/rune

ENTRYPOINT ["/bin/enclave-agent"]
EOF

docker build . -t $1

popd