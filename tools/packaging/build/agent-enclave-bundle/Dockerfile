FROM ubuntu:20.04 as builder

ARG KBC=cc-kbc

RUN apt-get update && \
    env DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    make \
    ca-certificates \
    curl \
    wget \
    gnupg

RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
RUN echo "deb [arch=amd64 signed-by=/usr/share/keyrings/intel-sgx.gpg] https://download.01.org/intel-sgx/sgx_repo/ubuntu focal main" | tee -a /etc/apt/sources.list.d/intel-sgx.list \
 && wget -qO - https://download.01.org/intel-sgx/sgx_repo/ubuntu/intel-sgx-deb.key | gpg --dearmor --output /usr/share/keyrings/intel-sgx.gpg \
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
    libdevmapper-dev \
    libfuse2 \
    libfuse3-3 \
    ca-certificates \
    protobuf-compiler \
    rsync \
    occlum-toolchains-glibc \
    occlum-pal=0.29.7-1 \
    occlum-sgx-tools=0.29.7-1 \
    occlum-runtime=0.29.7-1 \
    occlum=0.29.7-1

COPY src/ /enclave-cc/src/
RUN cd /enclave-cc/src/enclave-agent && \
    . $HOME/.cargo/env && \
    make ${KBC}

WORKDIR /run/enclave-agent
COPY tools/packaging/build/agent-enclave-bundle/enclave-agent-${KBC}.yaml .
COPY tools/packaging/build/agent-enclave-bundle/jq.filter .

ARG SGX_MODE=HW
RUN export PATH="$PATH:/opt/occlum/build/bin" && \
    occlum new occlum_instance && cd occlum_instance && \
    copy_bom -f ../enclave-agent-${KBC}.yaml --root image \
    --include-dir /opt/occlum/etc/template && \
    mv Occlum.json Occlum.json.orig && \
    jq -a -f ../jq.filter Occlum.json.orig > Occlum.json && \
    SGX_MODE=${SGX_MODE} occlum build && \
    if [ $SGX_MODE="SIM" ]; then \
        export DEBUG="--debug"; \
    fi; occlum package ${DEBUG}

# TODO: add new build stage and copy occlum_instance.tar.gz to it
WORKDIR /run/rune
RUN tar xzf /run/enclave-agent/occlum_instance/occlum_instance.tar.gz && \
    rm -rf /run/enclave-agent

RUN rm -rf $HOME/.cargo $HOME/.rustup /enclave-cc && sed -e '/cargo/d' -i /root/.profile && sed -e '/cargo/d' -i /root/.bashrc
RUN apt-get purge -y wget gnupg tzdata jq occlum occlum-pal occlum-toolchains-glibc make binutils libfuse2 libfuse3-3 ca-certificates rsync build-essential cmake git && apt-get autoremove -y
RUN echo "/run/rune/occlum_instance/build/lib/" | tee /etc/ld.so.conf.d/occlum-pal.conf && \
    echo "/opt/occlum/sgxsdk-tools/sdk_libs/" | tee -a /etc/ld.so.conf.d/occlum-pal.conf && \
    ldconfig
ENTRYPOINT ["/bin/enclave-agent"]
