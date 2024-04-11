package main

import (
	"log"
	"net/rpc"
	"net/http"
	"net"
	"errors"
	"time"
	"flag"

	"github.com/alexandremr01/user-elections/config"
	"github.com/alexandremr01/user-elections/client"
	"github.com/alexandremr01/user-elections/algorithms/bully"
	"github.com/alexandremr01/user-elections/algorithms/raft"
	"github.com/alexandremr01/user-elections/state"
)

type Algorithm interface{
	InitializeNode()
	StartElections()
	SendHeartbeat()
	GetServer() any
}

func main() {
	// get command line arguments
	algorithmName := flag.String("algorithm", "raft" , "Algorithm name (raft or bully)")
	configFile := flag.String("config", "config.json" , "Configuration file")
	flag.Parse()

	// get config from json and env vars
	config, err := config.GetConfig(*configFile)
	if err != nil {
		log.Fatal("error parsing config: ", err)
	}

	// build necessary dependencies
	connection := client.NewClient(config.NodeID, config.Addresses)
	state := state.NewState()
	algorithm, err := getAlgorithm(config.IDs, *algorithmName, state, connection, config)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// run in a second thread
	go mainLoop(algorithm, state, config)

	// listen to RPC requests
	server := algorithm.GetServer()
	registerAndServe(server, config.Port)
}

func getAlgorithm(ids []int, algorithmName string, state *state.State, connection *client.Client, config *config.Config) (Algorithm, error) {
	if algorithmName == "bully"{
		return bully.NewElections(ids, config.NodeID, state, connection, config.ElectionDuration), nil
	} else if algorithmName == "raft" {
		return raft.NewElections(ids, config.NodeID, state, connection, config.ElectionDuration), nil
	} 
	return nil, errors.New("unrecognized algorithm")
}

func mainLoop(algorithm Algorithm, state *state.State, config *config.Config) {
    log.Printf("My ID: %d\n", config.NodeID)
	algorithm.InitializeNode()
	for {
		if (config.NodeID == state.CoordinatorID) {
			algorithm.SendHeartbeat()
			time.Sleep(config.HeartbeatDuration)
		} else {
			time.Sleep(config.TimeoutDuration)
			if (state.LastHearbeat == nil) || (time.Now().Sub(*state.LastHearbeat) > config.TimeoutDuration) {
				log.Printf("Leader timed out")
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