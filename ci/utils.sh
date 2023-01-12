#!/bin/bash

ECC_RC_NAME="enclave-cc"
ECC_RC_VER="v0.2.0"
TIMEOUT_SECS=120
CI_DEBUG_MODE=true

is_cc_operator_controller_manager_pod_exist() {
    operator_pod_info=$(kubectl get pods -n confidential-containers-system 2>&1 | grep cc-operator-controller-manager)
    return $?
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
    timeout $TIMEOUT_SECS bash -c '
        ECC_RC_NAME=$0
        while [ true ]
        do
            kubectl get pods -n confidential-containers-system 2>&1 | grep cc-operator-daemon-install
            is_cc_pod_destroy=$?  
            kubectl get runtimeclass 2>&1 | grep $ECC_RC_NAME
            is_cc_runtimeclass_destroy=$? 

            if (( $is_cc_pod_destroy && $is_cc_runtimeclass_destroy )); then
                break
            fi
            sleep 1
        done 
    ' $ECC_RC_NAME
    return $?
}
