FROM ubuntu:20.04 as builder

ARG RATS_TLS_REV=5de6fc308f2a18f4105429451c85e6173b42ab9d

RUN apt-get update && \
    env DEBIAN_FRONTEND=noninteractive apt-get install -y \
    make \
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
    libfuse2 \
    libfuse3-3 \
    ca-certificates \
    protobuf-compiler \
    rsync \
    occlum-runtime \
    occlum-toolchains-glibc \
    occlum

WORKDIR /opt
ARG SGX_SDK_URL=https://download.01.org/intel-sgx/sgx-linux/2.18.1/distro/ubuntu20.04-server/sgx_linux_x64_sdk_2.18.101.1.bin
RUN wget ${SGX_SDK_URL} \
 && export SGX_SDK_INSTALLER=$(basename $SGX_SDK_URL) \
 && chmod +x $SGX_SDK_INSTALLER \
 && echo "yes" | ./$SGX_SDK_INSTALLER \
 && rm $SGX_SDK_INSTALLER \
 && cp /opt/sgxsdk/lib64/libsgx_uae_service_sim.so /lib/x86_64-linux-gnu/ \
 && /opt/sgxsdk/uninstall.sh

RUN git clone https://github.com/inclavare-containers/rats-tls.git && \
    cd rats-tls && \
    git checkout $RATS_TLS_REV && \
    cmake -DRATS_TLS_BUILD_MODE="occlum" -DBUILD_SAMPLES=on -H. -Bbuild && \
    make -C build install

COPY src/ /enclave-cc/src/
RUN cd /enclave-cc/src/enclave-agent && \
    . $HOME/.cargo/env && \
    make

WORKDIR /run/enclave-agent
COPY tools/packaging/build/agent-enclave-bundle/enclave-agent.yaml .
COPY tools/packaging/build/agent-enclave-bundle/jq.filter .

ARG SGX_MODE=HW
RUN export PATH="$PATH:/opt/occlum/build/bin" && \
    occlum new occlum_instance && cd occlum_instance && \
    copy_bom -f ../enclave-agent.yaml --root image \
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

RUN rm -rf $HOME/.cargo $HOME/.rustup /enclave-cc
RUN apt-get purge -y wget gnupg tzdata jq occlum occlum-toolchains-glibc make binutils libfuse2 libfuse3-3 ca-certificates rsync build-essential cmake git && apt-get autoremove -y
RUN echo "/run/rune/occlum_instance/build/lib/" | tee /etc/ld.so.conf.d/occlum-pal.conf && \
    ldconfig
ENTRYPOINT ["/bin/enclave-agent"]
