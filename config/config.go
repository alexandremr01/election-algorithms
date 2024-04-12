package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/alexandremr01/user-elections/types"
)

type JSONConfig struct {
	Addresses map[string]string `json:"addresses"`
}

func parseJSONConfig(fileName string) (map[int]string, error) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var result JSONConfig
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		return nil, err
	}

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

func GetConfig() (*types.Config, error) {
	cliArguments := parseCLI()

	addresses, err := parseJSONConfig(cliArguments.configFile)
	if err != nil {
		return nil, fmt.Errorf("error parsing json config: %w", err)
	}
	ids := []int{}
	for id := range addresses {
		ids = append(ids, id)
	}

	timeout, err := strconv.Atoi(os.Getenv("NODE_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("error parsing timeout: %w", err)
	}
	timeoutDuration := time.Duration(timeout) * time.Millisecond

	electionTimeout, err := strconv.Atoi(os.Getenv("ELECTION_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("error parsing election timeout: %w", err)
	}
	electionDuration := time.Duration(electionTimeout) * time.Millisecond

	heartbeatTime, err := strconv.Atoi(os.Getenv("HEARTBEAT_TIME"))
	if err != nil {
		return nil, fmt.Errorf("error parsing heartbeat time: %w", err)
	}
	heartbeatTimeDuration := time.Duration(heartbeatTime) * time.Second

	return &types.Config{
		HeartbeatDuration: heartbeatTimeDuration,
		TimeoutDuration:   timeoutDuration,
		ElectionDuration:  electionDuration,
		Addresses:         addresses,
		IDs:               ids,
		Port:              cliArguments.port,
		NodeID:            cliArguments.id,
		AlgorithmName:     cliArguments.algorithmName,
	}, nil
}
