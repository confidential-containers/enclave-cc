FROM rust:1.63-bullseye as builder

RUN git clone --depth 1 -b enclave-cc https://github.com/occlum/occlum.git

WORKDIR occlum/demos/runtime_boot/init
ADD v1-0001-init-code-change.patch .
RUN GIT_COMMITTER_EMAIL=foo@example.com git am v1-0001-init-code-change.patch && \
    rustup component add rustfmt && \
    cargo build --release

# Start preparing boot-instance "bundle"
FROM ubuntu:20.04
RUN apt-get update && \
    env DEBIAN_FRONTEND=noninteractive apt-get install -y \
    wget \
    gnupg

RUN echo "deb [arch=amd64 signed-by=/usr/share/keyrings/intel-sgx.gpg] https://download.01.org/intel-sgx/sgx_repo/ubuntu focal main" | tee -a /etc/apt/sources.list.d/intel-sgx.list \
 && wget -qO - https://download.01.org/intel-sgx/sgx_repo/ubuntu/intel-sgx-deb.key | gpg --dearmor --output /usr/share/keyrings/intel-sgx.gpg \
 && apt-get update \
 && env DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    libsgx-enclave-common \
    libsgx-urts \
    libsgx-uae-service \
    libsgx-quote-ex \
    libsgx-dcap-ql \
    libsgx-dcap-quote-verify \
    libsgx-dcap-default-qpl

RUN echo "deb [arch=amd64 signed-by=/usr/share/keyrings/occlum.gpg] http://mirrors.openanolis.cn/inclavare-containers/ubuntu20.04 focal main" | tee -a /etc/apt/sources.list.d/occlum.list \
 && wget -qO - http://mirrors.openanolis.cn/inclavare-containers/ubuntu20.04/DEB-GPG-KEY.key | gpg --dearmor --output /usr/share/keyrings/occlum.gpg \
 && apt-get update
RUN env DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    jq \
    tzdata \
    make \
    binutils \
    libfuse2 \
    libfuse3-3 \
    rsync \
    occlum-toolchains-glibc \
    occlum-runtime \
    occlum

WORKDIR /run/occlum
COPY --from=builder /occlum/demos/runtime_boot/init/target/release/init .
COPY --from=builder /occlum/demos/runtime_boot/boot_template.yaml .
ADD init.yaml .
ADD jq.filter .

RUN export PATH="$PATH:/opt/occlum/build/bin" && \
    occlum new boot_instance && cd boot_instance && \
    copy_bom -f ../boot_template.yaml --root image \
    --include-dir /opt/occlum/etc/template && \
    copy_bom -f ../init.yaml --root initfs \
    --include-dir /opt/occlum/etc/template && \
    mv Occlum.json Occlum.json.orig && \
    jq -a -f ../jq.filter Occlum.json.orig > Occlum.json && \
    occlum build && \
    occlum package

# TODO: we run occlum package to get a "clean" instance copied to a target
# directory. The optimization here could be that we just use the occlum_instance
# "working directory".
WORKDIR /run/rune
RUN tar xzf /run/occlum/boot_instance/boot_instance.tar.gz && \
    rm -rf /run/occlum

RUN apt-get remove -y wget gnupg make tzdata jq occlum occlum-toolchains-glibc binutils libfuse2 libfuse3-3 rsync && apt-get autoremove -y
RUN echo "/run/rune/boot_instance/build/lib/" | tee /etc/ld.so.conf.d/occlum-pal.conf && \
    ldconfig

ENTRYPOINT ["/bin/init"]
