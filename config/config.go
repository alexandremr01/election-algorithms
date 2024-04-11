package config

import (
	"fmt"
	"os"
	"errors"
	"io/ioutil"
	"encoding/json"
	"time"
	"strconv"
)

type Config struct {
	TimeoutDuration time.Duration
	ElectionDuration time.Duration
	HeartbeatDuration time.Duration
	Port string
	NodeID int
	Addresses map[int]string
	IDs []int
}

type JsonConfig struct {
	Addresses map[string]string `json:"addresses"`
}

func parseJsonConfig(fileName string) (map[int]string, error) {
	jsonFile, err := os.Open(fileName)
    if err != nil {
        return nil, err
    }
    defer jsonFile.Close()

    byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
        return nil, err
    }

    var result JsonConfig
    json.Unmarshal(byteValue, &result)
	
	addresses := make(map[int]string)
	for key, value := range result.Addresses {
		intKey, err := strconv.Atoi(key)
		if err != nil {
			return nil, errors.New("address key must be integer")
		}
		addresses[intKey] = value
	}

	return addresses, nil
}

func GetConfig(jsonConfigName string) (*Config, error) {
	addresses, err := parseJsonConfig(jsonConfigName)
	if err != nil {
		return nil, fmt.Errorf("error parsing json config: %s", err)
	}
	ids := []int{}
	for id := range addresses {
		ids = append(ids,id)
	}

	nodeID, err := strconv.Atoi(os.Getenv("NODE_ID")) 
	if err != nil {
		return nil, fmt.Errorf("error parsing node id:", err)
	}
	timeout, err := strconv.Atoi(os.Getenv("NODE_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("error parsing timeout:", err)
	}
	timeoutDuration := time.Duration(timeout)*time.Millisecond

	electionTimeout, err := strconv.Atoi(os.Getenv("ELECTION_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("error parsing election timeout:", err)
	}
	electionDuration := time.Duration(electionTimeout)*time.Millisecond

	heartbeatTime, err := strconv.Atoi(os.Getenv("HEARTBEAT_TIME"))
	if err != nil {
		return nil, fmt.Errorf("error parsing heartbeat time:", err)
	}
	heartbeatTimeDuration := time.Duration(heartbeatTime)*time.Second
	port := os.Getenv("PORT")

	return &Config{
		Port: port,
		NodeID: nodeID, 
		HeartbeatDuration: heartbeatTimeDuration,
		TimeoutDuration: timeoutDuration,
		ElectionDuration: electionDuration,
		Addresses: addresses,
		IDs: ids,
	}, nil
}