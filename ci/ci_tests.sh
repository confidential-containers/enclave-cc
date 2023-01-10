#!/bin/bash
source  ./ci/utils.sh

if [ ! -n "$1" ]; then
    echo "error: missing input parameter, such as install_coco_operator."
    exit 1
fi

function install_coco_operator(){
    if ! is_cc_operator_controller_manager_pod_exist; then 
        exit 1
    fi
    if [ $CI_DEBUG_MODE = true ]; then
        echo "[Debug] Start to install the operator..."
    fi

    logs=$(timeout $TIMEOUT_SECS kubectl apply -k github.com/confidential-containers/operator/config/release?ref=$ECC_RC_VER 2>&1 )
    case $? in 
        124)
            echo "[Error] Timeout when installing operator."
            echo "$logs"
            exit 1
        ;;
        1)
            echo "[Error] Something is wrong when installing operator."
            echo "$logs"
            exit 1
        ;;
    esac
    
    if ! check_cc_operator_controller_manager_pod_ready; then 
        exit 1
    fi
    exit 0
}

function install_enclave_cc_runtimeclass(){
    if ! is_enclave_cc_runtimeclass_exist; then 
        exit 1
    fi
    if [ $CI_DEBUG_MODE = true ]; then
        echo "[Debug] Start to install enclave-cc runtimeclass."
    fi

    logs=$(timeout $TIMEOUT_SECS kubectl apply -f https://raw.githubusercontent.com/confidential-containers/operator/main/config/samples/enclave-cc/base/ccruntime-enclave-cc.yaml 2>&1)
    case $? in 
        124)
            echo "[Error] Timeout when installing enclave-cc runtimeclass."
            echo "$logs"
            exit  1
        ;;
        1)
            echo "[Error] Something is wrong when installing  enclave-cc runtimeclass."
            echo "$logs"
            exit  1
        ;;
    esac

    if ! check_enclave_cc_runtimeclass_exist; then 
        exit 1
    fi
    exit  0
}

function uninstall_coco_operator(){
    if  is_cc_operator_controller_manager_pod_exist; then 
        exit 1
    fi
    if [ $CI_DEBUG_MODE = true ]; then 
        echo "[Debug] Start to delete operator..."
    fi
    logs=$(timeout $TIMEOUT_SECS kubectl delete -k github.com/confidential-containers/operator/config/release?ref=$ECC_RC_VER 2>&1)
    case $? in 
        124)
            echo "[Error] Timeout when uninstalling operator."
            echo "$logs"
            exit 1
        ;;
        1)   
            echo "[Error] Something is wrong when uninstalling operator."
            echo "$logs"
            exit 1
        ;;
    esac
    if [ $CI_DEBUG_MODE = true ]; then
        echo "[Debug] Successfully delete operator..."
    fi 
    exit 0
}

function uninstall_enclave_cc_runtimeclass(){
    if is_enclave_cc_runtimeclass_exist; then 
        exit 1
    fi
    if [ $CI_DEBUG_MODE = true ]; then
        echo "[Debug] Start to delete enclave-cc runtimeclass..."
    fi

    logs=$(timeout $TIMEOUT_SECS kubectl delete -f https://raw.githubusercontent.com/confidential-containers/operator/main/config/samples/enclave-cc/base/ccruntime-enclave-cc.yaml 2>&1)
    rtn_code=$?
    case $? in 
       124)
            echo "[Error] Timeout when uninstalling enclave-cc runtimeclass."
            if [ $CI_DEBUG_MODE = true ]; then
                echo "[Debug] Start to delete kubernetes stuck CRD deletion..."
                kubectl get ccruntimes.confidentialcontainers.org -o yaml 2>&1 | sed '/finalizer/{n;d;}'
                kubectl apply -f /home/ecc/operator/hlc/my-resource.yaml
            fi
            echo "$logs"
            exit  1
        ;;
        1)
            echo "[Error] Something is wrong when uninstalling enclave-cc runtimeclass."
            echo "$logs"
            exit 1
        ;;
    esac
    if [ $CI_DEBUG_MODE = true ]; then
        echo "[Debug] Successfully delete runtimeclass..."
    fi 
    exit 0
}

$1
