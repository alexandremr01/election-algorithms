export HEARTBEAT_INTERVAL=2000 
export ELECTION_TIMEOUT=5000 
export ALGORITHM=bully

go run main.go --id 1 --port 1231 --node_timeout 2500 \
                --config process_config.json \
                --election_timeout $ELECTION_TIMEOUT \
                --heartbeat_interval $HEARTBEAT_INTERVAL \
                --algorithm $ALGORITHM &\
go run main.go --id 2 --port 1232 --node_timeout 3000 \
                --config process_config.json \
                --election_timeout $ELECTION_TIMEOUT \
                --heartbeat_interval $HEARTBEAT_INTERVAL \
                --algorithm $ALGORITHM &\
go run main.go --id 3 --port 1233 --node_timeout 3500 \
                --config process_config.json \
                --election_timeout $ELECTION_TIMEOUT \
                --heartbeat_interval $HEARTBEAT_INTERVAL \
                --algorithm $ALGORITHM &\
go run main.go --id 4 --port 1234 --node_timeout 4000 \
                --config process_config.json \
                --election_timeout $ELECTION_TIMEOUT \
                --heartbeat_interval $HEARTBEAT_INTERVAL \
                --algorithm $ALGORITHM
                