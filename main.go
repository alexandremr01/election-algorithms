package main

import (
	"log"
	"net/rpc"
	"net/http"
	"net"
	"fmt"
	"time"
	"flag"

	"github.com/alexandremr01/user-elections/config"
	"github.com/alexandremr01/user-elections/client"
	"github.com/alexandremr01/user-elections/algorithms"
	"github.com/alexandremr01/user-elections/algorithms/types"
	"github.com/alexandremr01/user-elections/state"
)

func main() {
	cliArguments := parseCLI()

	// get config from json and env vars
	config, err := config.GetConfig(cliArguments.configFile)
	if err != nil {
		log.Fatal("error parsing config: ", err)
	}

	// build necessary dependencies
	connection := client.NewClient(config.NodeID, config.Addresses)
	state := state.NewState()
	algorithm, err := algorithms.GetAlgorithm(cliArguments.algorithmName, config, state, connection)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// run in a second thread
	go mainLoop(algorithm, state, config)

	// listen to RPC requests
	server := algorithm.GetServer()
	registerAndServe(server, config.Port)
}

func mainLoop(algorithm types.Algorithm, state *state.State, config *config.Config) {
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


type cliArguments struct {
	configFile string
	algorithmName string
}

func parseCLI() *cliArguments{
	// format algorithms list
	algorithmList := algorithms.GetAlgorithmsList()
	algorithmListStr := ""
	for _, alg := range algorithmList {
		algorithmListStr += alg + ","
	}
	algorithmListStr = algorithmListStr[:len(algorithmListStr)-1]
	algorithmHelp := fmt.Sprintf("Algorithm name (%s)", algorithmListStr)
	// get command line arguments
	algorithmName := flag.String("algorithm", "raft" , algorithmHelp)
	configFile := flag.String("config", "config.json" , "Configuration file")
	flag.Parse()
	return &cliArguments{configFile: *configFile, algorithmName: *algorithmName}
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