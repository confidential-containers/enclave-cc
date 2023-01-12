#!/bin/bash

source  ./ci/utils.sh

echo "ci debug mode is $CI_DEBUG_MODE"

if [ $CI_DEBUG_MODE = true ]; then
    echo "[Debug] ECC_RC_VER is $ECC_RC_VER "
fi

install_coco_operator() {
    # if is_cc_operator_controller_manager_pod_exist; then 
    #     echo "[Error] Found CcCo operator pod."
    #     echo "Please uninstall operator pod to ensure a clean CI environment."
    #     return 1
    # fi

    if [ $CI_DEBUG_MODE = true ]; then
        echo "[Debug] Start to install the operator..."
    fi

    if [ $CI_DEBUG_MODE = true ]; then
        echo "[Debug] ECC_RC_VER is $ECC_RC_VER "
    fi

    logs=$(timeout $TIMEOUT_SECS kubectl apply -k github.com/confidential-containers/operator/config/release?ref=$ECC_RC_VER 2>&1)
    case $? in 
        124)
            echo "[Error] Timeout when installing operator."
            echo "$logs"
            return 1
        ;;
        1)
            echo "[Error] Something is wrong when installing operator."
            echo "$logs"
            return 1
        ;;
    esac
    
    wait_cc_operator_controller_manager_pod_ready
    case $? in 
        124)
            echo "[Error] Timeout when installing operator."
            return 1
        ;;
        1)
            echo "[Error] Something is wrong when installing operator."
            return 1
        ;;
    esac

    # workaround
    # todo: check cc-operator-pre-install-daemon status

    if [ $CI_DEBUG_MODE = true  ]; then
        echo "Successfully install CoCo operator."
    fi

    return 0
}

install_enclave_cc_runtimeclass() {
    # if is_enclave_cc_runtimeclass_exist; then
    #     echo "[Error] Found enclave-cc runtimeclass."
    #     echo "Please uninstall enclave-cc runtime to ensure a clean CI environment."
    #     return 1
    # fi

    if [ $CI_DEBUG_MODE = true ]; then
        echo "[Debug] Start to install enclave-cc runtimeclass."
    fi

    # logs=$(timeout $TIMEOUT_SECS kubectl apply -f https://raw.githubusercontent.com/confidential-containers/operator/main/config/samples/enclave-cc/base/ccruntime-enclave-cc.yaml 2>&1)
    logs=$(timeout $TIMEOUT_SECS kubectl apply -f /home/ecc/operator/hlc/ccruntime-enclave-cc.yaml 2>&1)
    case $? in 
        124)
            echo "[Error] Timeout when installing enclave-cc runtimeclass."
            echo "$logs"
            return 1
        ;;
        1)
            echo "[Error] Something is wrong when installing enclave-cc runtimeclass."
            echo "$logs"
            return 1
        ;;
    esac

    wait_enclave_cc_runtimeclass_ready
    case $? in
        124)
            echo "[Error] Timeout when installing enclave-cc runtimeclass."
            return 1
        ;;
        1)
            echo "[Error] Something is wrong when installing enclave-cc runtimeclass."
            return 1
        ;;
    esac

    if [ $CI_DEBUG_MODE = true ]; then
        echo "[OK] Successfully install enclave-cc runtimeclass."
    fi

    return 0
}

uninstall_coco_operator() {
    if ! is_cc_operator_controller_manager_pod_exist; then 
        echo "[Error] Not found CoCo operator pod."
        return 1
    fi

    if [ $CI_DEBUG_MODE = true ]; then 
        echo "[Debug] Start to delete operator..."
    fi

    logs=$(timeout $TIMEOUT_SECS kubectl delete -k github.com/confidential-containers/operator/config/release?ref=$ECC_RC_VER 2>&1)
    case $? in 
        124)
            echo "[Error] Timeout when uninstalling operator."
            echo "$logs"
            return 1
        ;;
        1)   
            echo "[Error] Something is wrong when uninstalling operator."
            echo "$logs"
            return 1
        ;;
    esac

    if [ $CI_DEBUG_MODE = true ]; then
        echo "[Debug] Successfully delete operator..."
    fi

    return 0
}

uninstall_enclave_cc_runtimeclass() {
    # if ! is_enclave_cc_runtimeclass_exist; then
    #     echo "[Error] Not found enclave-cc runtimeclass."
    #     return 1
    # fi

    if [ $CI_DEBUG_MODE = true ]; then
        echo "[Debug] Start to delete enclave-cc runtimeclass..."
    fi

    # workaround
    # if [ ! -d /opt/confidential-containers/bin ]; then
    #     echo "[Debug] The directory doesn't exist. Start to create /opt/confidential-containers/bin dir."
    #     mkdir /opt/confidential-containers/bin
    # fi

    # logs=$(timeout $TIMEOUT_SECS kubectl delete -f https://raw.githubusercontent.com/confidential-containers/operator/main/config/samples/enclave-cc/base/ccruntime-enclave-cc.yaml 2>&1)
    logs=$(timeout $TIMEOUT_SECS kubectl delete -f /home/ecc/operator/hlc/ccruntime-enclave-cc.yaml 2>&1)
    case $? in 
        124)
            echo "[Error] Timeout when uninstalling enclave-cc runtimeclass."
            echo "$logs"

            # workaround
            echo "[Debug] Start to delete kubernetes stuck CRD..."
            kubectl get ccruntimes.confidentialcontainers.org -o yaml | sed '/finalizers/{s/$/ []/;n;d;}' > /tmp/my-resource.yaml
            kubectl apply -f /tmp/my-resource.yaml

            # return 1
        ;;
        1)
            echo "[Error] Something is wrong when uninstalling enclave-cc runtimeclass."
            echo "$logs"
            return 1
        ;;
    esac

    wait_enclave_cc_runtimeclass_terminating
    case $? in
        124) 
            echo "[Error] Timeout when uninstalling enclave-cc runtimeclass."
            return 1
        ;;
        1)
            echo "[Error] Something is wrong when uninstalling enclave-cc runtimeclass."
            return 1
        ;;
    esac

    if [ $CI_DEBUG_MODE = true ]; then
        echo "[Debug] Successfully delete runtimeclass."
    fi 

    return 0
}

apply_eaa_cosign_encryped_hello_world_workload() {
    # workaround
    

    kubectl apply -f ./ci/case_configs/eaa_cosign_encrypted_hello_world.yaml
    wait_workload_output
    if [ $? != 0 ]; then
        echo "[Error] Fail to run the eaa cosign encrypted hello-world workload"
        return 1
    fi
    echo "[OK] Successfully run the eaa cosign encrypted hello-world workload."
    return 0
}

delete_eaa_cosign_encrypted_hello_world_workload() {
    kubectl delete -f ./ci/case_configs/eaa_cosign_encrypted_hello_world.yaml
    if [ $? != 0 ]; then
        echo "[Error] Fail to delete the eaa cosign encrypted hello-world workload"
        return 1
    fi
    echo "[OK] Successfully apply the eaa cosign encrypted hello-world workload."
    return 0
}

if [ $# != 1 ]; then
    echo "[CI script error] Please specify one task."
    exit 1
fi

$1
