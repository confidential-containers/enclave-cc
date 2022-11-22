#!/bin/bash
set -e

if [ ! -n "$1" ] ;then
    echo "error: missing input parameter, please input image tag, such as zhiwei/occlum-enclave-agent-app:v1.0."
    exit 1
fi

./build_occlum_instance.sh

./gen_docker_image.sh $1

# we typically name the [path-to-agent-bundle] by 
# /**/agent-container, and we will use "agent-container" 
# to refer the bundle.
agentContainerPath="[path-to-agent-bundle]"

if [ ! -d "$agentContainerPath" ]; then
  mkdir -p $agentContainerPath
fi


echo $agentContainerPath
pushd $agentContainerPath
rm -rf rootfs && mkdir rootfs

docker export $(docker create $1) | sudo tar -C rootfs -xvf -

cp /etc/resolv.conf rootfs/etc/
cp /etc/hostname    rootfs/etc/

# cp [path-to-sgx-qcnl]/sgx_default_qcnl.conf rootfs/etc/

# NOTE: if we directly launch the agent by rune/runc
# we need to create image directory in manual. The 
# procedure is carried by `shim` in normal.
# mkdir -p rootfs/images/[image-name]/sefs/lower
# mkdir -p rootfs/images/[image-name]/sefs/upper

popd