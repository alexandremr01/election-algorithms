#!/bin/bash

NODE_ID=$1

kill_process() {
    local NODE_ID=$1
    local PID_FILE="tmp/node_${NODE_ID}_pid.txt"

    if [ ! -f "$PID_FILE" ]; then
        echo "PID file $PID_FILE does not exist."
        return 1
    fi

    PID=$(cat "$PID_FILE")
    if kill "$PID" > /dev/null 2>&1; then
        echo "Process with PID $PID has been terminated."
    else
        echo "Failed to kill process with PID $PID."
    fi

    rm "$PID_FILE"
}

if [[ "$NODE_ID" == "all" ]]; then
    for id in {1..4} 
    do
        kill_process $id
    done
else
    kill_process $NODE_ID
fi
 