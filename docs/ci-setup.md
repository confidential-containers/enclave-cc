# Notes to Setup an SGX VM for CI Testing

# SGX PSW

- install `aesmd`, configure with `default quoting type = ecdsa_256` and the preferred DCAP quote provider lib (with its config)

# containerd
We pre-install the CoCo containerd and tools with a simple command:

```bash
curl -fsSL https://github.com/confidential-containers/containerd/releases/download/v1.6.8.1/cri-containerd-cni-1.6.8.1-linux-amd64.tar.gz | sudo tar zx -C /
```

And it's configured with the following runtimehandler:
```
$ tail -3 /etc/containerd/config.toml
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.enclavecc]
  cri_handler = "cc"
  runtime_type = "io.containerd.rune.v2"
```

# KBS (for CC-KBC testing)

We can set up a KBS cluster with `docker-compose` quickly.

```
git clone https://github.com/confidential-containers/kbs.git && cd kbs
```

1. change the [attestation-service's](https://github.com/confidential-containers/attestation-service/blob/main/Dockerfile.as) image in `docker-compose.yml` to `docker.io/xynnn007/attestation-service:sgx-v0.6.0`
2. change the PCCS configuration volume of as in `docker-compose.yml` to `/etc/sgx_default_qcnl.conf:/etc/sgx_default_qcnl.conf:rw`
3. Run

```bash
openssl genpkey -algorithm ed25519 > config/private.key
openssl pkey -in config/private.key -pubout -out config/public.pub
```

4. Run with `docker compose up -d`, and KBS will listen to port `8080` for requests

We then add default configs under `data/kbs-storage`.

1. Image policy
```bash
mkdir -p data/kbs-storage/default/security-policy
cat << EOF > data/kbs-storage/default/security-policy/test
{
    "default": [{"type": "insecureAcceptAnything"}],
    "transports": {}
}
EOF
```
2. Image decryption key (for ghcr.io/confidential-containers/test-container-enclave-cc:encrypted)

The key and the key id are defined in the [test image's Dockerfile](../tools/packaging/build/test-image/Dockerfile)
```bash
mkdir -p data/kbs-storage/default/image-kek
echo LieOhvkqFcGMzZrVzt6vPWlj/F/bgYMNe45vhQpdxAA= | base64 -d > data/kbs-storage/default/image-kek/11032d96-dccd-46a3-9244-93644d76745f
```
# Github Runner Service
For `enclave-cc` e2e tests, we run a "job-started" pre-cleanup job configured
for the runner:

```
Environment=ACTIONS_RUNNER_HOOK_JOB_STARTED=<path/to>/job-started-cleanup.sh
```

The script content:
```
#!/bin/bash

echo "delete previous workspace $GITHUB_WORKSPACE"
pushd $GITHUB_WORKSPACE
sudo rm -rf coco src
popd

echo "delete lingering pods"
for i in $(sudo crictl pods -q); do
    sudo crictl -t 10s stopp $i;
    sudo crictl -t 10s rmp $i;
done

echo "docker system prune"
docker system prune -a -f
```
