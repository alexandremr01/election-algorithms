go run main.go --id $NODE_ID --port 8000 --node_timeout $NODE_TIMEOUT \
                --config docker_config.json \
                --election_timeout $ELECTION_TIMEOUT \
                --heartbeat_interval $HEARTBEAT_INTERVAL \
                --autofailure 15 --autofailure_duration 9000 \
                --algorithm $ALGORITHM 