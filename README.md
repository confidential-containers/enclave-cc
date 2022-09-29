# enclave-cc
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fconfidential-containers%2Fenclave-cc.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fconfidential-containers%2Fenclave-cc?ref=badge_shield)

## Introduction

Confidential containers are the product of the combination of confidential computing technology and cloud-native technology.
Confidential containers use hardware-based Trusted Execution Environments (HW-TEE) for resource isolation, data protection,
and remote attestation. They can protect data in execution from cloud service providers and other privileged third parties.

Confidential containers have the following characteristics: they provide confidentiality and integrity for data,
especially for runtime data; they facilitate end-to-end security for the deployment and operation of sensitive workloads;
they provide a way to prove that the environment in which confidential workloads are launched is authentic and trustworthy;
tenants maintain almost the same experience as using ordinary containers, but can deploy sensitive applications with greater
confidence in their security.

The Confidential Containers project provides a way to protect cloud native applications in a HW-TEE without additional modification
to the container image during development. HW-TEE-based confidential container isolation can take two forms: process and virtual
machine based isolation.

Enclave-cc provides a process-based confidential container solution leveraging Intel SGX. The term enclave is often associated with
process-based isolation.

## Motivation

Enclave-CC enables process-based isolation. Process-based isolation is beneficial by drawing the isolation boundary exactly around
each container process. This reduces the Trusted Computing Base, or the amount of the system your application depends on for security.

Consistent with the VM-based isolation design, the enclave-cc approach does not measure and attest to the workload (user application).
Instead it uses a generic enclave to isolate user's application from the rest of the system. Measurement and attestation of the empty
generic enclave are used in building trust to deliver the application securely.

The Confidential Containers project, including the Enclave-CC approach provides all the elements for a complete deployment flow,
from creating user container image with encryption/signature to uploading it to an image registry, to pulling the image, verifying its
signature and decrypting it, unpacking it into user rootfs in enclave, then booting up user containers and mounting encrypted rootfs,
all the container processes are isolated with enclaves.

## Documentation

- [Design](docs/design.md)


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fconfidential-containers%2Fenclave-cc.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fconfidential-containers%2Fenclave-cc?ref=badge_large)
