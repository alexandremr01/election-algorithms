#!/bin/bash

export HEARTBEAT_INTERVAL=2000 
export ELECTION_TIMEOUT=5000 
export ALGORITHM=${1:-bully}
export NODE_ID=${2}

go build

if [ -z "$NODE_ID" ] # if second parameter is not set, then run all
    then     
    mkdir -p tmp
    for id in {1..4} 
    do
        node_timeout=$((2000 + id * 500))

        ./user-elections --id $id --node_timeout $node_timeout \
            --config process_config.json \
            --election_timeout $ELECTION_TIMEOUT \
            --heartbeat_interval $HEARTBEAT_INTERVAL \
            --algorithm $ALGORITHM &

        # writes PID of the last background process to file
        pid=$!
        echo $pid > "tmp/node_${id}_pid.txt"
    done


    else # if it is set, run only that node id
    ./user-elections -id $NODE_ID --node_timeout 2500 \
        --config process_config.json \
        --election_timeout $ELECTION_TIMEOUT \
        --heartbeat_interval $HEARTBEAT_INTERVAL \
        --algorithm $ALGORITHM 
fi

