ARG IMAGE
FROM ${IMAGE:-ubuntu}:20.04
ARG ENCLAVE_CC_ARTIFACTS=./enclave-cc-static.tar.xz
ARG DESTINATION=/opt/enclave-cc-artifacts

COPY ${ENCLAVE_CC_ARTIFACTS} ${WORKDIR}
COPY runtimeclass ${DESTINATION}/runtimeclass

ENV DEBIAN_FRONTEND=noninteractive
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

RUN \
    apt-get update && \
    apt-get install -y --no-install-recommends apt-transport-https ca-certificates curl xz-utils systemd gpg && \
    mkdir -p /etc/apt/keyrings/ && \
    curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.28/deb/Release.key | gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg && \
    echo "deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.28/deb/ /" | tee /etc/apt/sources.list.d/kubernetes.list && \
    apt-get update && \
    apt-get install -y --no-install-recommends kubectl && \
    apt-get clean && rm -rf /var/lib/apt/lists/ && \
    mkdir -p ${DESTINATION} && \
    tar xvf ${ENCLAVE_CC_ARTIFACTS} -C ${DESTINATION} && \
    rm ${WORKDIR}/${ENCLAVE_CC_ARTIFACTS}

WORKDIR ${DESTINATION}
