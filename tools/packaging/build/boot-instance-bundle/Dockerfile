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
    make \
    wget \
    gnupg

RUN echo "deb [arch=amd64 signed-by=/usr/share/keyrings/intel-sgx.gpg] https://download.01.org/intel-sgx/sgx_repo/ubuntu focal main" | tee -a /etc/apt/sources.list.d/intel-sgx.list \
 && wget -qO - https://download.01.org/intel-sgx/sgx_repo/ubuntu/intel-sgx-deb.key | gpg --dearmor --output /usr/share/keyrings/intel-sgx.gpg \
 && apt-get update \
 && env DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    libsgx-dcap-ql \
    libsgx-dcap-quote-verify \
    libsgx-dcap-default-qpl

WORKDIR /opt
ARG SGX_SDK_URL=https://download.01.org/intel-sgx/sgx-linux/2.17.1/distro/ubuntu20.04-server/sgx_linux_x64_sdk_2.17.101.1.bin
RUN wget ${SGX_SDK_URL} \
 && export SGX_SDK_INSTALLER=$(basename $SGX_SDK_URL) \
 && chmod +x  $SGX_SDK_INSTALLER \
 && echo "yes" | ./$SGX_SDK_INSTALLER \
 && rm $SGX_SDK_INSTALLER

RUN echo "deb [arch=amd64 signed-by=/usr/share/keyrings/occlum.gpg] http://mirrors.openanolis.cn/inclavare-containers/ubuntu20.04 focal main" | tee -a /etc/apt/sources.list.d/occlum.list \
 && wget -qO - http://mirrors.openanolis.cn/inclavare-containers/ubuntu20.04/DEB-GPG-KEY.key | gpg --dearmor --output /usr/share/keyrings/occlum.gpg \
 && apt-get update
RUN env DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    jq \
    tzdata \
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

ARG SGX_MODE=HW
RUN export PATH="$PATH:/opt/occlum/build/bin" && \
    occlum new boot_instance && cd boot_instance && \
    copy_bom -f ../boot_template.yaml --root image \
    --include-dir /opt/occlum/etc/template && \
    copy_bom -f ../init.yaml --root initfs \
    --include-dir /opt/occlum/etc/template && \
    mv Occlum.json Occlum.json.orig && \
    jq -a -f ../jq.filter Occlum.json.orig > Occlum.json && \
    SGX_MODE=${SGX_MODE} occlum build && \
    if [ $SGX_MODE="SIM" ]; then \
        export DEBUG="--debug"; \
    fi; occlum package ${DEBUG}

# TODO: we run occlum package to get a "clean" instance copied to a target
# directory. The optimization here could be that we just use the occlum_instance
# "working directory".
WORKDIR /run/rune
RUN tar xzf /run/occlum/boot_instance/boot_instance.tar.gz && \
    rm -rf /run/occlum

RUN apt-get purge -y wget gnupg make tzdata jq occlum occlum-toolchains-glibc binutils libfuse2 libfuse3-3 rsync && apt-get autoremove -y
RUN echo "/run/rune/boot_instance/build/lib/" | tee /etc/ld.so.conf.d/occlum-pal.conf && \
    echo "/opt/sgxsdk/lib64" | tee /etc/ld.so.conf.d/sgxsdk.conf && \
    ldconfig

ENTRYPOINT ["/bin/init"]
