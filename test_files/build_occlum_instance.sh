#!/bin/bash
set -e

# compile enclave-agent supporting rats-tls
pushd enclave-agent
make rats-tls
popd

# build and package an occlum instance
rm -rf occlum_instance
occlum new occlum_instance

pushd occlum_instance
# tune the parameters for the occlum runtime 
# "metadata.debuggable = true/false" opens/closes the debug log.
# In product env, set it to false
new_json="$(jq '.resource_limits.user_space_size = "1024MB" |
            .resource_limits.kernel_space_heap_size= "512MB" |
            .resource_limits.kernel_space_stack_size= "128MB" |
            .resource_limits.max_num_of_threads = 16 |
            .metadata.debuggable = true ' Occlum.json)" && \
echo "${new_json}" > Occlum.json

rm -rf image
# /opt/occlum/etc is the default directory where the occlum is 
# installed. Replace it if you customize the directory.
copy_bom -f ../enclave-agent.yaml --root image --include-dir /opt/occlum/etc/template

occlum build
# option "--debug" helps build debuggable apps.
# In product env, remove it
occlum package --debug  occlum_instance.tar.gz
popd