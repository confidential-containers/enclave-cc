#!/bin/bash
ECC_RC_NAME="enclave-cc"
ECC_RC_VER="v0.2.0"
TIMEOUT_SECS=120
CI_DEBUG_MODE=true

function is_enclave_cc_runtimeclass_exist() {
    ecc_rc_info=$( kubectl get runtimeclass 2>&1 | grep $ECC_RC_NAME  )
    if [ $? = 0 ]; then
        echo "Found k8s runtimeclass: "
        echo "$ecc_rc_info"
        echo "Please uninstall enclave-cc runtime to ensure a clean CI environment."
        return 1
    else
        echo "Can not found k8s runtimeclass, please install."
        return 0
    fi       
}

function is_cc_operator_controller_manager_pod_exist() {
    operator_pod_info=$(kubectl get pods -n confidential-containers-system 2>&1 | grep cc-operator-controller-manager)
    if [ $? = 0  ]; then
        echo "Found operator pod :"
        echo "$operator_pod_info"
        echo "Please uninstall operator pod to ensure a clean CI environment."
        return 1
    else
        echo "Can not found operator pod, please install."
        return 0
    fi  
}

function check_cc_operator_controller_manager_pod_ready(){
   timeout $TIMEOUT_SECS bash -c '
        CI_DEBUG_MODE=$0
        while [ true ]
        do
            pod_ready_status=$(kubectl get pods -n confidential-containers-system 2>&1 | grep cc-operator-controller-manager | awk '"'"'{print $2}'"'"')
            if [ $CI_DEBUG_MODE = true ]; then
                echo "[Debug] pod_ready_status: $pod_ready_status"
            fi

            if [ $pod_ready_status = "2/2" ]; then
                break
            fi
            sleep 1
        done
    ' $CI_DEBUG_MODE
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
    if [ CI_DEBUG_MODE=true ]; then
        echo "[Debug] Succefully installed the operator."
    fi
    return 0
}

function check_enclave_cc_runtimeclass_exist() {
    timeout $TIMEOUT_SECS bash -c '
        TIMEOUT_SECS=$0
        ECC_RC_NAME=$1
        counter=0
        while [ $counter -lt $TIMEOUT_SECS ]
        do
            kubectl get runtimeclass 2>&1 | grep $ECC_RC_NAME > /dev/null
            if [ $? = 0 ]; then
                exit 0
            fi
            sleep 1
            (( counter++ ))
        done ' $TIMEOUT_SECS $ECC_RC_NAME
    rtn_code=$?
    case $rtn_code in
        124)
            echo "[Error] Timeout when installing enclave-cc runtimeclass."
            return 1
        ;;
        1)
            echo "[Error] Something is wrong when installing enclave-cc runtimeclass."
            return 1
        ;;
        0)
            echo "[OK] Successfully install enclave-cc runtimeclass."
        ;;
    esac
    return 0
}
