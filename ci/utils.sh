#!/bin/bash

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
    local TIMEOUT_SECS=$1
    local POD_NAME=$2

    timeout $TIMEOUT_SECS bash -c '
        pod_name=$0
        while [ true ]
        do
            pod_info=$(kubectl get pods -A 2>&1 | grep $pod_name | awk '"'"'{print $3" "tolower($4)}'"'"')
            IFS='"'"' '"'"'
            echo $pod_info
            read -a temp_arr <<< $pod_info
            pod_status=${temp_arr[1]}
            temp=${temp_arr[0]}
            IFS="/"
            read -a temp_arr <<< $temp
            num_pod_ready=${temp_arr[0]}
            num_pod_total=${temp_arr[1]}
            
            echo "[Debug] pod_ready_status: $num_pod_ready/$num_pod_total $pod_status"

            if [[ $num_pod_ready = $num_pod_total && $pod_status = running ]]; then
                break
            fi

            sleep 2
        done
    ' $POD_NAME

    return $?
}

wait_runtimeclass_ready() {
    local TIMEOUT_SECS=$1
    local RUNTIMECLASS_NAME=$2

    timeout $TIMEOUT_SECS bash -c '
        runtimeclass_name=$0
        while [ true ]
        do
            kubectl get runtimeclass 2>&1 | grep $runtimeclass_name
            if [ $? = 0 ]; then
                break
            fi
            sleep 2
        done
    ' $RUNTIMECLASS_NAME

    return $?
}

wait_runtimeclass_deleted() {
    local TIMEOUT_SECS=$1
    local RUNTIMECLASS_NAME=$2

    timeout $TIMEOUT_SECS bash -c '
        runtimeclass_name=$0
        while [ true ]
        do
            kubectl get runtimeclass 2>&1 | grep $runtimeclass_name
            if [ $? != 0 ]; then
                break
            fi
            sleep 2
        done
    ' $RUNTIMECLASS_NAME

    return $?
}

wait_pod_terminating() {
    local TIMEOUT_SECS=$1
    local POD_NAME=$2

    timeout $TIMEOUT_SECS bash -c '
        pod_name=$0
        while [ true ]
        do
            kubectl get pods -A 2>&1 | grep $pod_name
            if [ $? != 0 ]; then
                break
            fi
            sleep 2 
        done
    ' $POD_NAME

    return $?
}

wait_pod_log() {
    local TIMEOUT_SECS=$1
    local POD_NAME=$2
    local LOG_CONTENT=$3
    
    timeout $TIMEOUT_SECS bash -c '
        pod_name=$0
        log_content=$1
        while [ true ]
        do
            kubectl logs $pod_name 2>&1 | grep "$log_content"
            if [ $? == 0 ]; then
                break
            fi
            sleep 2
        done
    ' $POD_NAME $LOG_CONTENT

    return $?
}
