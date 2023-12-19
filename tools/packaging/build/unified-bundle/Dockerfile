FROM ubuntu:20.04 as builder

ARG KBC=cc-kbc

RUN apt-get update && \
    env DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    make \
    ca-certificates \
    curl \
    wget \
    gnupg

ARG OCCLUM_VERSION=0.30.0-1
ARG SGXSDK_VERSION=2_22_100
ARG RUST_VERSION=1.72.0

RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y --default-toolchain=${RUST_VERSION}
RUN echo "deb [arch=amd64 signed-by=/usr/share/keyrings/intel-sgx.gpg] https://download.01.org/intel-sgx/sgx_repo/ubuntu focal main" | tee -a /etc/apt/sources.list.d/intel-sgx.list \
 && wget -qO - https://download.01.org/intel-sgx/sgx_repo/ubuntu/intel-sgx-deb.key | gpg --dearmor --output /usr/share/keyrings/intel-sgx.gpg \
 && wget -qO - https://download.01.org/intel-sgx/sgx_repo/ubuntu/apt_preference_files/99sgx_${SGXSDK_VERSION}_focal_custom_version.cfg | tee -a /etc/apt/preferences.d/99sgx_sdk.cfg \
 && apt-get update \
 && env DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    libsgx-uae-service \
    libsgx-dcap-ql \
    libsgx-dcap-ql-dev \
    libsgx-dcap-quote-verify \
    libsgx-dcap-quote-verify-dev \
    libsgx-quote-ex \
    libsgx-dcap-default-qpl

RUN echo "deb [arch=amd64 signed-by=/usr/share/keyrings/occlum.gpg] https://occlum.io/occlum-package-repos/debian focal main" | tee -a /etc/apt/sources.list.d/occlum.list \
 && wget -qO - https://occlum.io/occlum-package-repos/debian/public.key | gpg --dearmor --output /usr/share/keyrings/occlum.gpg \
 && apt-get update
RUN env DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    binutils \
    build-essential \
    cmake \
    git \
    libssl-dev \
    jq \
    tzdata \
    binutils \
    libclang-dev \
    libfuse2 \
    libfuse3-3 \
    ca-certificates \
    protobuf-compiler \
    rsync \
    occlum-toolchains-glibc \
    occlum-pal=${OCCLUM_VERSION} \
    occlum-sgx-tools=${OCCLUM_VERSION} \
    occlum-runtime=${OCCLUM_VERSION} \
    occlum=${OCCLUM_VERSION}

COPY src/ /enclave-cc/src/
RUN cd /enclave-cc/src/enclave-agent && \
    . $HOME/.cargo/env && \
    make ${KBC}
RUN cd /enclave-cc/src/runtime-boot/init && \
    . $HOME/.cargo/env && \
    cargo build --release

WORKDIR /run/unified
COPY tools/packaging/build/unified-bundle/enclave-agent-${KBC}.yaml .
COPY tools/packaging/build/unified-bundle/boot_template.yaml .
COPY tools/packaging/build/unified-bundle/init.yaml .
COPY tools/packaging/build/unified-bundle/agent-instance.filter .

ARG SGX_MODE=HW
RUN export PATH="$PATH:/opt/occlum/build/bin" && \
    occlum new occlum_instance && cd occlum_instance && \
    copy_bom -f ../enclave-agent-${KBC}.yaml --root image \
    --include-dir /opt/occlum/etc/template && \
    copy_bom -f ../init.yaml --root initfs \
    --include-dir /opt/occlum/etc/template && \
    mv Occlum.json Occlum.json.orig && \
    jq -a -f ../agent-instance.filter Occlum.json.orig > Occlum.json && \
    SGX_MODE=${SGX_MODE} occlum build && \
    if [ $SGX_MODE="SIM" ]; then \
        export DEBUG="--debug"; \
    fi; occlum package ${DEBUG}

FROM ubuntu:20.04

RUN apt-get update && \
    env DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    ca-certificates \
    wget \
    gnupg

ARG OCCLUM_VERSION=0.30.0-1
ARG SGXSDK_VERSION=2_22_100

RUN echo "deb [arch=amd64 signed-by=/usr/share/keyrings/intel-sgx.gpg] https://download.01.org/intel-sgx/sgx_repo/ubuntu focal main" | tee -a /etc/apt/sources.list.d/intel-sgx.list \
 && echo "deb [arch=amd64 signed-by=/usr/share/keyrings/occlum.gpg] https://occlum.io/occlum-package-repos/debian focal main" | tee -a /etc/apt/sources.list.d/occlum.list \
 && wget -qO - https://occlum.io/occlum-package-repos/debian/public.key | gpg --dearmor --output /usr/share/keyrings/occlum.gpg \
 && wget -qO - https://download.01.org/intel-sgx/sgx_repo/ubuntu/intel-sgx-deb.key | gpg --dearmor --output /usr/share/keyrings/intel-sgx.gpg \
 && wget -qO - https://download.01.org/intel-sgx/sgx_repo/ubuntu/apt_preference_files/99sgx_${SGXSDK_VERSION}_focal_custom_version.cfg | tee -a /etc/apt/preferences.d/99sgx_sdk.cfg \
 && apt-get update \
 && env DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    libsgx-uae-service \
    libsgx-dcap-ql \
    libsgx-dcap-quote-verify \
    libsgx-quote-ex \
    libsgx-dcap-default-qpl \
    tzdata \
    libfuse2 \
    libfuse3-3 \
    occlum-runtime=${OCCLUM_VERSION} \
    occlum-sgx-tools=${OCCLUM_VERSION}

WORKDIR /run/rune
COPY --from=builder /run/unified/occlum_instance/occlum_instance.tar.gz .

RUN tar xf occlum_instance.tar.gz && \
    echo "/run/rune/occlum_instance/build/lib/" | tee /etc/ld.so.conf.d/occlum-pal.conf && \
    echo "/opt/occlum/sgxsdk-tools/sdk_libs/" | tee -a /etc/ld.so.conf.d/occlum-pal.conf && \
    ldconfig && \
    /opt/occlum/sgxsdk-tools/bin/x64/sgx_sign dump -dumpfile /run/agent.txt -cssfile /run/sigstruct1.bin -enclave occlum_instance/build/lib/libocclum-libos.signed.so && \
    mkdir /configs && \
    rm *.tar.gz
