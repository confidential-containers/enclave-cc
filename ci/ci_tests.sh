#!/bin/bash

CURRENT_DIR=${0%$(basename $0)}
source  $CURRENT_DIR"utils.sh"

COCO_OPERATOR_VERSION="v0.2.0"
COCO_OPERATOR_POD_NAME="cc-operator-controller-manager"
ECC_PAYLOAD_POD_NAME_1="cc-operator-daemon-install"
ECC_PAYLOAD_POD_NAME_2="cc-operator-pre-install-daemon"
ECC_RUNTIMECLASS_NAME="enclave-cc"
ECC_RUNTIMECLASS_CONFIG=$CURRENT_DIR"case_configs"/"ccruntime_for_ci.yaml"
TIMEOUT_SECS=300
CI_DEBUG_MODE=true

install_coco_operator() {
    if [ $CI_DEBUG_MODE = true ]; then
        echo "[Debug] Start to install the operator..."
    fi

    timeout $TIMEOUT_SECS kubectl apply -k github.com/confidential-containers/operator/config/release?ref=$COCO_OPERATOR_VERSION 2>&1
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
    
    wait_pod_ready $TIMEOUT_SECS $COCO_OPERATOR_POD_NAME
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

    if [ $CI_DEBUG_MODE = true ]; then
        echo "Successfully install CoCo operator."
    fi
    return 0
}

install_enclave_cc_runtimeclass() {
    if [ $CI_DEBUG_MODE = true ]; then
        echo "[Debug] Start to install enclave-cc runtimeclass."
    fi

    timeout $TIMEOUT_SECS kubectl apply -f $ECC_RUNTIMECLASS_CONFIG 2>&1
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

    wait_pod_ready $TIMEOUT_SECS $ECC_PAYLOAD_POD_NAME_1
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

    wait_pod_ready $TIMEOUT_SECS $ECC_PAYLOAD_POD_NAME_2
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

    wait_runtimeclass_ready $TIMEOUT_SECS $ECC_RUNTIMECLASS_NAME
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
    if [ $CI_DEBUG_MODE = true ]; then 
        echo "[Debug] Start to delete operator..."
    fi

    timeout $TIMEOUT_SECS kubectl delete -k github.com/confidential-containers/operator/config/release?ref=$COCO_OPERATOR_VERSION 2>&1
    case $? in 
        124)
            echo "[Error] Timeout when uninstalling operator."
            return 1
        ;;
        1)   
            echo "[Error] Something is wrong when uninstalling operator."
            return 1
        ;;
    esac

    wait_pod_terminating $TIMEOUT_SECS $COCO_OPERATOR_POD_NAME
    case $? in
        124)
            echo "[Error] Timeout when uninstalling operator."
            return 1
        ;;
        1)
            echo "[Error] Something is wrong when uninstalling operator"
            return 1
        ;;
    esac

    if [ $CI_DEBUG_MODE = true ]; then
        echo "[Debug] Successfully delete operator..."
    fi
    return 0
}

uninstall_enclave_cc_runtimeclass() {
    if [ $CI_DEBUG_MODE = true ]; then
        echo "[Debug] Start to delete enclave-cc runtimeclass..."
    fi

    timeout $TIMEOUT_SECS kubectl delete -f $ECC_RUNTIMECLASS_CONFIG 2>&1
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

    wait_pod_terminating $TIMEOUT_SECS $ECC_PAYLOAD_POD_NAME_1
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

    wait_pod_terminating $TIMEOUT_SECS $ECC_PAYLOAD_POD_NAME_2
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

    wait_runtimeclass_deleted $TIMEOUT_SECS $ECC_RUNTIMECLASS_NAME
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
        echo "[Debug] Successfully delete enclave-cc runtimeclass."
    fi 

    return 0
}

apply_eaa_cosign_encryped_hello_world_workload() {
    WORKLOAD_POD_NAME="enclave-cc-pod"
    kubectl apply -f $CURRENT_DIR"case_configs"/"eaa_cosign_encrypted_hello_world.yaml"

    wait_pod_ready $TIMEOUT_SECS $WORKLOAD_POD_NAME
    case $? in
        124) 
            echo "[Error] Timeout when running workload."
            return 1
        ;;
        1)
            echo "[Error] Something is wrong when running workload."
            return 1
        ;;
    esac

    wait_pod_log $TIMEOUT_SECS $WORKLOAD_POD_NAME "Hello world!"
    case $? in
        124) 
            echo "[Error] Timeout when running workload."
            return 1
        ;;
        1)
            echo "[Error] Something is wrong when running workload."
            return 1
        ;;
    esac

    if [ $CI_DEBUG_MODE = true ]; then
        echo "[OK] Successfully run the eaa cosign encrypted hello-world workload."
    fi

    return 0
}

delete_eaa_cosign_encrypted_hello_world_workload() {
    WORKLOAD_POD_NAME="enclave-cc-pod"
    kubectl delete -f $CURRENT_DIR"case_configs"/"eaa_cosign_encrypted_hello_world.yaml"

    wait_pod_terminating $TIMEOUT_SECS $WORKLOAD_POD_NAME
    case $? in
        124) 
            echo "[Error] Timeout when deleting workload."
            return 1
        ;;
        1)
            echo "[Error] Something is wrong when deleting workload."
            return 1
        ;;
    esac

    if [ $CI_DEBUG_MODE = true ]; then
        echo "[OK] Successfully delete the eaa cosign encrypted hello-world workload."
    fi
    
    return 0
}

if [ $# != 1 ]; then
    echo "[CI script error] Please specify one task."
    exit 1
fi

$1
