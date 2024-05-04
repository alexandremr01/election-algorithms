package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/alexandremr01/user-elections/algorithms"
)

type cliArguments struct {
	configFile                string
	algorithmName             string
	id                        int
	port                      string
	nodeTimeoutDuration       time.Duration
	heartbeatIntervalDuration time.Duration
	electionTimeoutDuration   time.Duration
	autoFailure               int
	autoFailureDuration       time.Duration
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
	nodeTimeout := flag.Int("node_timeout", 2000, "Node Timeout")
	heartbeatInterval := flag.Int("heartbeat_interval", 2000, "Hearbeat Interval")
	electionTimeout := flag.Int("election_timeout", 5000, "Election Timeout")
	autoFailure := flag.Int("autofailure", 0, "Chance of automatic failure in %")
	autoFailureInterval := flag.Int("autofailure_duration", 0, "Duration of autofailure")

	flag.Parse()
	autoFailureIntervalDuration := time.Duration(*autoFailureInterval) * time.Millisecond
	heartbeatIntervalDuration := time.Duration(*heartbeatInterval) * time.Millisecond
	nodeTimeoutDuration := time.Duration(*nodeTimeout) * time.Millisecond
	electionTimeoutDuration := time.Duration(*electionTimeout) * time.Millisecond

	return &cliArguments{
		configFile:                *configFile,
		algorithmName:             *algorithmName,
		port:                      *port,
		id:                        *id,
		nodeTimeoutDuration:       nodeTimeoutDuration,
		heartbeatIntervalDuration: heartbeatIntervalDuration,
		electionTimeoutDuration:   electionTimeoutDuration,
		autoFailure:               *autoFailure,
		autoFailureDuration:       autoFailureIntervalDuration,
	}
}
