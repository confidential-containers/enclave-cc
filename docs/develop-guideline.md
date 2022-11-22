# 1. Overview

The `enclave-agent` runs on a SGX-based LibOS. In the runtime, the memory and the mounted filesystem used by the enclave agent is encrypted by the SGX. 

The supporting LibOS includes
- [occlum](https://github.com/occlum)
    - the encrypted filesystem is called by `unionfs`
- [gramine](https://github.com/gramineproject/gramine)

The dependencies installation process can be found in the [develop-env-install.md](develop-env-install.md).


## Image Management

The enclave agent plays the role of the container image management. It is responsible for four basic missions, including
- Pulling: it pulls the index, manifest, config and layers of a image, which is dependent on the crate [image-rs](https://github.com/confidential-containers/image-rs)
- Unpacking: it decrypts and decompresses the layers of a image if necessary. The decryption is dependent on
the crate [ocicrypt-rs](https://github.com/confidential-containers/ocicrypt-rs).
- Storing: it saves the unpacking layers on an encrypted filesystem
- Mounting: it merges layers to a rootfs for a application container and mounts the rootfs to a mounting point


## Security Management

The enclave agent need to satisfy security requirements, including
- attestation
    - remote attestation: it provides security information for the remove 
    - local attestation
    - the crate [attestation-agent](https://github.com/confidential-containers/attestation-agent) supports a native attesation client
- signature verification: 
    - some images are signed by the pulisher or the store or both, it need to verify these signatures
    - the crate [attestation-agent](https://github.com/confidential-containers/attestation-agent) supports a native signature client 

# 2. Development Guide

## Occlum-based Enclave Agent

### *Dependencies*

#### **rats-tls**

The [rats-tls](https://github.com/inclavare-containers/inclavare-containers/tree/master/rats-tls) is a library for communicating on a secured channel, which supports the native attestation agent with feature `eaa kbc`.

The rats-tls depends on the SGX SDK (following the [installing guide](develop-env-install.md)), libraries `libsgx-dcap-ql` and `libsgx-dcap-ql-dev` (install them if necessary).

The rats-tls provides a different version for the application on the occlum or on the host. We want to build and run a KBC (Key Broken Service) supporting `eaa kbc` on the host, we should build and install rats-tls for the host. The library will be installed in `/usr/local/lib/rats-tls`.
```bash
git clone https://github.com/inclavare-containers/rats-tls.git
cd rats-tls
cmake -DRATS_TLS_BUILD_MODE="host" -DBUILD_SAMPLES=on -H. -Bbuild
make -C build install
```

If we want to build the library for application running on the occlum, we need to build an occlum version. 
```bash
git clone https://github.com/inclavare-containers/rats-tls.git
cd rats-tls
cmake -DRATS_TLS_BUILD_MODE="occlum" -DBUILD_SAMPLES=on -H. -Bbuild
make -C build install
```
The library is also installed in the `/usr/local/lib/rats-tls` in default. If we want to keep two versions on the same machine, we need to install the occlum version instance into a customized place `[rats-tls-path]` and refer it in a config of the occlum instance.


### *Agent Bundle Building*

There are three steps to build a occlum-based enclave agent. Copy following three scripts under directory `enclave-cc/src`.

#### Step 1 Build A Occlum Instance

The script template `test_files/build_occlum_instance.sh` will use occlum dev tools to wrap the enclave agent as an occlum instance. The [page1](https://github.com/occlum/occlum) and [page2](https://github.com/occlum/occlum/blob/master/docs/resource_config_guide.md) can provide some suggestions to tune parameters.

The script template uses one config `test_files/enclave-agent.yaml`. It describes some copy opertions that copy dependencies from the host to the occlum instance. Complete the template and copy the config into enclave-cc/src.

`test_files/enclave-agent.yaml` need to refer two configs, including
- `test_files/ocicrypt.conf`: it is config for crate `ocicrypt-rs`.
- `test_files/agent.conf`: it is for the enclave agent.

Copy them to where you desire and refer them correctly 

`test_files/enclave-agent.yaml` need to refer the library `rats-tls` too.

#### Step 2 Build A Container

The script template is in `test_files/gen_docker_image.sh`.

#### Step 3 Build A Container Bundle

The script template is in `test_files/gen_bundle.sh`, it will use above two scripts. Build the container bundle of the enclave agent by
```bash
gen_bundle.sh [image_name]
```

After the script completing, the agent bundle is placed at `[agent-container]`, copy the `test_files/config.json` into it.


### *Image Pulling Testing*

#### Pulling Plain-text Images

Given that we want to launch the enclave agent by `rune`, and we have created the `[agent-container]`. Enter into `[agent-container]`, we will find a `config.json` and a `rootfs`. If we want to pull the image `docker.io/redis`, confirm that `rootfs/images/redis/sefs/lower` and `rootfs/images/redis/sefs/upper` is created. 

Launch the enclave agent by 
```bash
rune run [name]
```

If we use `test_files/config.json`, the agent will listen on `tcp://0.0.0.0:6661`. 

The image pulling client is located in `tools/debug`, build client and pulling `redis`.
```bash
cd tools/debug
make 
target/debug/examples/async-client -c tcp://0.0.0.0:6661 -i docker.io/redis
```

The image pulling success, when the server logs `pulling xxx successful`. 

#### Pulling Encrypted Images 

TODO: 