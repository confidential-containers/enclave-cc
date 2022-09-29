FROM rust:1.63-bullseye as builder

# Build enclave-agent
RUN apt-get update && \
    env DEBIAN_FRONTEND=noninteractive apt-get install -y \
    protobuf-compiler
RUN git clone --depth 1 https://github.com/confidential-containers/enclave-cc.git && \
    cd enclave-cc/src/enclave-agent && \
    rustup component add rustfmt && \
    make

# Start preparing enclave-agent "bundle"
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

RUN echo "deb [arch=amd64 signed-by=/usr/share/keyrings/occlum.gpg] https://occlum.io/occlum-package-repos/debian focal main" | tee -a /etc/apt/sources.list.d/occlum.list \
 && wget -qO - https://occlum.io/occlum-package-repos/debian/public.key | gpg --dearmor --output /usr/share/keyrings/occlum.gpg \
 && apt-get update
RUN env DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    jq \
    tzdata \
    make \
    binutils \
    libfuse2 \
    libfuse3-3 \
    ca-certificates \
    rsync \
    occlum-runtime \
    occlum-toolchains-glibc \
    occlum

WORKDIR /run/enclave-agent
COPY --from=builder /enclave-cc/src/enclave-agent/target/release/enclave-agent .
ADD enclave-agent.yaml .
ADD jq.filter .

RUN export PATH="$PATH:/opt/occlum/build/bin" && \
    occlum new occlum_instance && cd occlum_instance && \
    copy_bom -f ../enclave-agent.yaml --root image \
    --include-dir /opt/occlum/etc/template && \
    mv Occlum.json Occlum.json.orig && \
    jq -a -f ../jq.filter Occlum.json.orig > Occlum.json && \
    occlum build && \
    occlum package

# TODO: we run occlum package to get a "clean" instance copied to a target
# directory. The optimization here could be that we just use the occlum_instance
# "working directory".
WORKDIR /run/rune
RUN tar xzf /run/enclave-agent/occlum_instance/occlum_instance.tar.gz && \
    rm -rf /run/enclave-agent

RUN apt-get remove -y wget gnupg tzdata jq occlum make binutils libfuse2 libfuse3-3 ca-certificates rsync && apt-get autoremove -y
RUN echo "/run/rune/occlum_instance/build/lib/" | tee /etc/ld.so.conf.d/occlum-pal.conf && \
    ldconfig

ENTRYPOINT ["/bin/enclave-agent"]
