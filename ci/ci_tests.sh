#!/bin/bash
source  ./ci/utils.sh

install_coco_operator() {
    if is_cc_operator_controller_manager_pod_exist; then 
        echo "[Error] Found CcCo operator pod."
        echo "Please uninstall operator pod to ensure a clean CI environment."
        return 1
    fi

    if [ $CI_DEBUG_MODE = true ]; then
        echo "[Debug] Start to install the operator..."
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

    if [ $CI_DEBUG_MODE = true  ]; then
        echo "Successfully install CoCo operator."
    fi

    return 0
}

install_enclave_cc_runtimeclass() {
    if is_enclave_cc_runtimeclass_exist; then
        echo "[Error] Found enclave-cc runtimeclass."
        echo "Please uninstall enclave-cc runtime to ensure a clean CI environment."
        return 1
    fi

    if [ $CI_DEBUG_MODE = true ]; then
        echo "[Debug] Start to install enclave-cc runtimeclass."
    fi

    logs=$(timeout $TIMEOUT_SECS kubectl apply -f https://raw.githubusercontent.com/confidential-containers/operator/main/config/samples/enclave-cc/base/ccruntime-enclave-cc.yaml 2>&1)
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

function uninstall_enclave_cc_runtimeclass(){
    if ! is_enclave_cc_runtimeclass_exist; then
        echo "[Error] Not found enclave-cc runtimeclass."
        return 1
    fi

    if [ $CI_DEBUG_MODE = true ]; then
        echo "[Debug] Start to delete enclave-cc runtimeclass..."
    fi

    # workaround
    if [ ! -f /opt/confidential-containers/bin ]; then
        echo "[Debug] The directory doesn't exist. Start to create /opt/confidential-containers/bin dir."
        mkdir /opt/confidential-containers/bin
    fi

    logs=$(timeout $TIMEOUT_SECS kubectl delete -f https://raw.githubusercontent.com/confidential-containers/operator/main/config/samples/enclave-cc/base/ccruntime-enclave-cc.yaml 2>&1)
    case $? in 
       124)
            echo "[Error] Timeout when uninstalling enclave-cc runtimeclass."
            echo "$logs"

            # workaround
            echo "[Debug] Start to delete kubernetes stuck CRD..."
            kubectl get ccruntimes.confidentialcontainers.org -o yaml | sed '/finalizers/{s/$/[]/;n;d;}' > /home/ecc/operator/hlc/my-resource.yaml
            kubectl apply -f /home/ecc/operator/hlc/my-resource.yaml

            # return 1
        ;;
        1)
            echo "[Error] Something is wrong when uninstalling enclave-cc runtimeclass."
            echo "$logs"
            return 1
        ;;
    esac

    # todo: wait enclave-cc runtimeclass deleted

    if [ $CI_DEBUG_MODE = true ]; then
        echo "[Debug] Successfully delete runtimeclass..."
    fi 

    return 0
}

if [ $# != 1 ]; then
    echo "[CI script error] Please specify a task."
    exit 1
fi

$1
