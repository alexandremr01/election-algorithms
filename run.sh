#!/bin/bash

export HEARTBEAT_INTERVAL=2000 
export ELECTION_TIMEOUT=5000 
export ALGORITHM=${1:-bully}
export NODE_ID=${2}

if [ -z "$NODE_ID" ]; # if second parameter is not set, then run all
    then     
    go run main.go --id 1 --node_timeout 2500 \
        --config process_config.json \
        --election_timeout $ELECTION_TIMEOUT \
        --heartbeat_interval $HEARTBEAT_INTERVAL \
        --algorithm $ALGORITHM &\
    go run main.go --id 2 --node_timeout 3000 \
        --config process_config.json \
        --election_timeout $ELECTION_TIMEOUT \
        --heartbeat_interval $HEARTBEAT_INTERVAL \
        --algorithm $ALGORITHM &\
    go run main.go --id 3 --node_timeout 3500 \
        --config process_config.json \
        --election_timeout $ELECTION_TIMEOUT \
        --heartbeat_interval $HEARTBEAT_INTERVAL \
        --algorithm $ALGORITHM &\
    go run main.go --id 4 --node_timeout 4000 \
        --config process_config.json \
        --election_timeout $ELECTION_TIMEOUT \
        --heartbeat_interval $HEARTBEAT_INTERVAL \
        --algorithm $ALGORITHM;

    else # if it is set, run only that node id
    go run main.go --id $NODE_ID --node_timeout 2500 \
        --config process_config.json \
        --election_timeout $ELECTION_TIMEOUT \
        --heartbeat_interval $HEARTBEAT_INTERVAL \
        --algorithm $ALGORITHM 
fi

