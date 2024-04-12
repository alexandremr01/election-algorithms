HEARTBEAT_TIME=2 ELECTION_TIMEOUT=5000 NODE_TIMEOUT=2500 go run main.go --algorithm raft --config process_config.json --id 1 --port 1231 &\
HEARTBEAT_TIME=2 ELECTION_TIMEOUT=5000 NODE_TIMEOUT=3000 go run main.go --algorithm raft --config process_config.json --id 2 --port 1232 &\
HEARTBEAT_TIME=2 ELECTION_TIMEOUT=5000 NODE_TIMEOUT=3500 go run main.go --algorithm raft --config process_config.json --id 3 --port 1233 &\
HEARTBEAT_TIME=2 ELECTION_TIMEOUT=5000 NODE_TIMEOUT=4000 go run main.go --algorithm raft --config process_config.json --id 4 --port 1234