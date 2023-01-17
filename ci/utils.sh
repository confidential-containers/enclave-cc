#!/bin/bash

ECC_RC_NAME="enclave-cc"
ECC_RC_VER="v0.2.0"
TIMEOUT_SECS=300
CI_DEBUG_MODE=true

is_pod_exist() {
    local POD_NAME=$1
    kubectl get pods -A 2>&1 | grep $POD_NAME
    return $?
}

is_runtimeclass_exist() {
    local RUNTIMECLASS_NAME=$1
    kubectl get runtimeclass 2>&1 | grep $RUNTIMECLASS_NAME
    returnh $?
}

wait_pod_ready() {
    local POD_NAME=$1
    local TIMEOUT_SECS=$2

    timeout $TIMEOUT_SECS bash -c '
        pod_name=$0
        while [ true ]
        do
            pod_info=$(kubectl get pods -A | grep $pod_name | awk '"'"'{print $3" "$4}'"'"')
            IFS=' '
            read -a $temp_arr <<< $pod_info
            pod_status=${temp_arr[1]}
            pod_status=$(echo "$pod_status" | awk '"'"'{print tolower($0)}'"'"')
            temp=${temp_arr[0]}
            IFS='/'
            read -a $temp_arr <<< $temp
            num_pod_ready=${temp_arr[0]}
            num_pod_total=${temp_arr[1]}
            
            echo "[Debug] pod_ready_status: $num_pod_ready/$num_pod_total $pod_status"

            if [[ num_pod_ready = num_pod_total && pod_status=running ]]; then
                break
            fi

            sleep 2
        done
    ' $POD_NAME

    return $?
}

wait_pod_terminating() {

}

wait_pod_log() {

}

wait_cc_operator_controller_manager_pod_ready() {
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
    return $?
}

wait_cc_operator_pre_install_daemon(){
    timeout $TIMEOUT_SECS bash -c '
        CI_DEBUG_MODE=$0
        while [ true ]
        do
            pod_ready_status=$(kubectl get pods -n confidential-containers-system 2>&1 | grep cc-operator-pre-install-daemon | awk '"'"'{print $2}'"'"')
            if [ $CI_DEBUG_MODE = true ]; then
                echo "[Debug] pod_ready_status: $pod_ready_status"
            fi

            if [ $pod_ready_status = "1/1" ]; then
                break
            fi
            sleep 1
        done
    ' $CI_DEBUG_MODE
    return $?
}

is_enclave_cc_runtimeclass_exist() {
    ecc_rc_info=$(kubectl get runtimeclass 2>&1 | grep $ECC_RC_NAME)
    return $?
}

wait_enclave_cc_runtimeclass_ready() {
    timeout $TIMEOUT_SECS bash -c '
        ECC_RC_NAME=$0
        while [ true ]
        do
            kubectl get runtimeclass 2>&1 | grep $ECC_RC_NAME > /dev/null
            if [ $? = 0 ]; then
                break
            fi
            sleep 1
        done 
    ' $ECC_RC_NAME
    return $?
}

wait_enclave_cc_runtimeclass_terminating() {
    # add detect cc-operator-post-uninstall-daemon, cc-operator-pre-install
    timeout $TIMEOUT_SECS bash -c '
        ECC_RC_NAME=$0
        while [ true ]
        do
            kubectl get pods -n confidential-containers-system 2>&1 | grep "cc-operator-daemon-install"
            is_cc_pod_exist=$?  
            kubectl get runtimeclass 2>&1 | grep "enclave-cc"
            is_cc_runtimeclass_exist=$? 

            if [[ $is_cc_pod_exist != 0 && $is_cc_runtimeclass_exist != 0 ]]; then
                break
            fi
            sleep 1
        done 
    ' $ECC_RC_NAME
    return $?
}

wait_workload_output() {
    timeout $TIMEOUT_SECS bash -c '
        while [ true ]
        do
            kubectl logs enclave-cc-pod | head -n 5 | grep "Hello world!"
            if [ $? = 0 ]; then
                break
            fi
            sleep 1
        done
    '
    return $?
}
