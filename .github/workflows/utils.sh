#!/bin/bash
ECC_RC_NAME="enclave-cc"
ECC_RC_VER="v0.2.0"
TIMEOUT_SECS=60
CI_DEBUG_MODE=true

function is_enclave_cc_runtimeclass_exist() {
    ecc_rc_info=$( kubectl get runtimeclass 2>&1 | grep $ECC_RC_NAME  )
    if [ $? = 0 ]
    then
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
    #operator_pod_status=$(kubectl get pods -n confidential-containers-system 2>&1 | grep cc-operator-controller-manager| awk '{print $2}')
    #operator_pod_status=$operator_pod_info | awk '{print $2}'
    if [ $? = 0  ]
    #if [ $operator_pod_status = "2/2" ]
    then
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
            if [ $CI_DEBUG_MODE = true ]
            then
                echo "[Debug] pod_ready_status: $pod_ready_status"
            fi

            if [ $pod_ready_status = "2/2" ]
            then
                break
            fi
            sleep 1
        done
    ' $CI_DEBUG_MODE
    rtn_code=$?
    if [ $rtn_code = 124 ]
    then
        echo "[Error] Timeout when installing operator pod."
        return 1
    elif [ $rtn_code != 0 ]
    then
        echo "[Error] Something is wrong when installing operator pod."
        return 1
    fi
    if [ CI_DEBUG_MODE=true ]
    then
        echo "[Debug] Succefully installed the operator pod."
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
            if [ $? = 0 ]
            then
                exit 0
            fi
            sleep 1
            (( counter++ ))
        done ' $TIMEOUT_SECS $ECC_RC_NAME
    rtn_code=$?
    if [ $rtn_code = 124 ]
    then
        echo "[Error] Timeout when installing Enclave-CC runtime."
        return 1
    elif [ $rtn_code != 0 ]
    then
        echo "[Error] Something is wrong when installing Enclave-CC runtime."
        return 1
    elif [ $rtn_code = 0 ]
    then
        echo "[OK] Successfully install Enclave-CC runtime."
    fi
    return 0
}