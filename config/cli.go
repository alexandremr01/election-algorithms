package config

import (
	"flag"
	"fmt"

	"github.com/alexandremr01/user-elections/algorithms"
)

type cliArguments struct {
	configFile    string
	algorithmName string
	id            int
	port          string
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
		configFile:    *configFile,
		algorithmName: *algorithmName,
		port:          *port,
		id:            *id,
	}
}
