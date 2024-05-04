package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

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

	return &types.Config{
		HeartbeatDuration:   cliArguments.heartbeatIntervalDuration,
		TimeoutDuration:     cliArguments.nodeTimeoutDuration,
		ElectionDuration:    cliArguments.electionTimeoutDuration,
		Addresses:           addresses,
		IDs:                 ids,
		Port:                cliArguments.port,
		NodeID:              cliArguments.id,
		AlgorithmName:       cliArguments.algorithmName,
		AutoFailure:         cliArguments.autoFailure,
		AutoFailureDuration: cliArguments.autoFailureDuration,
	}, nil
}
