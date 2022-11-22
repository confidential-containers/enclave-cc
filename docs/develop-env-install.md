# 1. SGX install

## 1.1 driver
The kernel at 5.11 or above has a built-in in-tree driver. If your kernel version is below 5.11, there are two options, including
- updating the kernel to 5.11 or above (recommanded)
- building and installing the out-of-tree driver by following the [guide](https://github.com/intel/linux-sgx-driver) 

## 1.2 SDK

Select and download the SDK installing script from [here](https://download.01.org/intel-sgx/latest/linux-latest/distro), such as `sgx_linux_x64_sdk_2.17.100.3.bin`. 

Run the script and specify the install directory `[path-to-install]`.  
- The directory `[path-to-install]/sgxsdk/lib64` contains dynamic libraries, run the `ldconfig -v -n [path-to-install]/sgxsdk/lib64` to update the library reference.
- The directory `[path-to-install]/sgxsdk/include` contains the header files, if compiler cannot find these files, use `export C_INCLUDE_PATH=$C_INCLUDE_PATH:[path-to-install]/include` to expose it.

## 1.3 PSW

### 1.3.1 Quote Generation and Verification

Following the [guide](https://github.com/intel/SGXDataCenterAttestationPrimitives) to install libraries for quote generation and verification, such as `libsgx_dcap_quoteverify.so` 

# 2. LibOS install

## 2.1 Occlum & Rune

The occlum containes runtime libOS and some dev tools. Follow the [guide](https://github.com/occlum/occlum) to install the occlum libOS dev tools on the branch `enclave-cc`.

The rune is used to launch a occlum contaienr. Follow [guide](https://github.com/inclavare-containers/inclavare-containers/tree/master/rune) to install the `rune`.

## 2.2 Gramine