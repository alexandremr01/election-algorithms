package config

import (
	"fmt"
	"os"
	"time"
	"strconv"
)

type Config struct {
	TimeoutDuration time.Duration
	ElectionDuration time.Duration
	HeartbeatDuration time.Duration
	Port string
	NodeID int
}

func GetConfig() (*Config, error) {
	nodeID, err := strconv.Atoi(os.Getenv("NODE_ID")) 
	if err != nil {
		fmt.Errorf("error parsing node id:", err)
	}
	timeout, err := strconv.Atoi(os.Getenv("NODE_TIMEOUT"))
	if err != nil {
		fmt.Errorf("error parsing timeout:", err)
	}
	timeoutDuration := time.Duration(timeout)*time.Millisecond

	electionTimeout, err := strconv.Atoi(os.Getenv("ELECTION_TIMEOUT"))
	if err != nil {
		fmt.Errorf("error parsing election timeout:", err)
	}
	electionDuration := time.Duration(electionTimeout)*time.Millisecond

	heartbeatTime, err := strconv.Atoi(os.Getenv("HEARTBEAT_TIME"))
	if err != nil {
		fmt.Errorf("error parsing heartbeat time:", err)
	}
	heartbeatTimeDuration := time.Duration(heartbeatTime)*time.Second
	port := os.Getenv("PORT")

	return &Config{
		Port: port,
		NodeID: nodeID, 
		HeartbeatDuration: heartbeatTimeDuration,
		TimeoutDuration: timeoutDuration,
		ElectionDuration: electionDuration,
	}, nil
}