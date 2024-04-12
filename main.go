package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"
	"os"

	"github.com/alexandremr01/user-elections/algorithms"
	"github.com/alexandremr01/user-elections/types"
	"github.com/alexandremr01/user-elections/client"
	"github.com/alexandremr01/user-elections/config"
	"github.com/alexandremr01/user-elections/state"
)

func main() {
	// get config from CLI, json and env vars
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
	log.Printf("My ID: %d, My PID: %d\n", config.NodeID, os.Getpid())
	algorithm.InitializeNode()
	for {
		if config.NodeID == state.CoordinatorID {
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

type cliArguments struct {
	configFile    string
	algorithmName string
	id int
	port string
}

func parseCLI() *cliArguments {
	// format algorithms list
	algorithmList := algorithms.GetAlgorithmsList()
	algorithmListStr := ""
	for _, alg := range algorithmList {
		algorithmListStr += alg + ","
	}
	algorithmListStr = algorithmListStr[:len(algorithmListStr)-1]
	algorithmHelp := fmt.Sprintf("Algorithm name (%s)", algorithmListStr)
	// get command line arguments
	algorithmName := flag.String("algorithm", "raft", algorithmHelp)
	configFile := flag.String("config", "config.json", "Configuration file")
	port := flag.String("port", "8000", "Port to connect")
	id := flag.Int("id", 1, "Node ID")
	flag.Parse()
	return &cliArguments{
		configFile: *configFile, 
		algorithmName: *algorithmName,
		port: *port,
		id: *id,
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
