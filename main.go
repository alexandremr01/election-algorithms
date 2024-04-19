package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"

	"github.com/alexandremr01/user-elections/algorithms"
	"github.com/alexandremr01/user-elections/client"
	"github.com/alexandremr01/user-elections/config"
	"github.com/alexandremr01/user-elections/state"
	"github.com/alexandremr01/user-elections/types"
)

func main() {
	// get config from CLI and JSON
	config, err := config.GetConfig()
	if err != nil {
		log.Fatal("error parsing config: ", err)
	}

	// build necessary dependencies
	connection := client.NewClient(config.NodeID, config.Addresses)
	state := state.NewState()
	algorithm, err := algorithms.GetAlgorithm(config.AlgorithmName, config, state, connection)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// run in a second thread
	go mainLoop(algorithm, state, config)

	// listen to RPC requests
	server := algorithm.GetServer()
	registerAndServe(server, config.Port)
}

func mainLoop(algorithm types.Algorithm, state *state.State, config *types.Config) {
	log.Printf("Node %d: My PID: %d\n", config.NodeID, os.Getpid())
	algorithm.InitializeNode()
	for {
		if config.NodeID == state.CoordinatorID {
			algorithm.SendHeartbeat()
			time.Sleep(config.HeartbeatDuration)
		} else {
			time.Sleep(config.TimeoutDuration)

			heartbeatTimedout := false
			hasHeartbeat := state.LastHearbeat != nil
			if hasHeartbeat {
				heartbeatTimedout = time.Since(*state.LastHearbeat) > config.TimeoutDuration
			}

			if !hasHeartbeat || heartbeatTimedout {
				log.Printf("Node %d: Leader timed out", config.NodeID)
				algorithm.StartElections()
			}
		}
	}
}

func registerAndServe(server any, port string) {
	rpc.Register(server)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", ":"+port)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}
