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
	Ports     map[string]string `json:"ports"`
}

type ParsedJSONConfig struct {
	Addresses map[int]string
	Ports     map[int]string
}

func parseJSONConfig(fileName string) (*ParsedJSONConfig, error) {
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

	ports := make(map[int]string)
	for key, value := range result.Ports {
		intKey, err := strconv.Atoi(key)
		if err != nil {
			return nil, errors.New("port key must be integer")
		}
		ports[intKey] = value
	}

	return &ParsedJSONConfig{
		Addresses: addresses,
		Ports:     ports,
	}, nil
}

func GetConfig() (*types.Config, error) {
	cliArguments := parseCLI()

	jsonConfigs, err := parseJSONConfig(cliArguments.configFile)
	if err != nil {
		return nil, fmt.Errorf("error parsing json config: %w", err)
	}
	ids := []int{}
	for id := range jsonConfigs.Addresses {
		ids = append(ids, id)
	}

	port := cliArguments.port
	jsonPort, hasJsonPort := jsonConfigs.Ports[cliArguments.id]
	if port == "" && hasJsonPort {
		port = jsonPort
	} else if port == "" {
		return nil, fmt.Errorf("no port defined")
	}

	return &types.Config{
		HeartbeatDuration:   cliArguments.heartbeatIntervalDuration,
		TimeoutDuration:     cliArguments.nodeTimeoutDuration,
		ElectionDuration:    cliArguments.electionTimeoutDuration,
		Addresses:           jsonConfigs.Addresses,
		IDs:                 ids,
		Port:                port,
		NodeID:              cliArguments.id,
		AlgorithmName:       cliArguments.algorithmName,
		AutoFailure:         cliArguments.autoFailure,
		AutoFailureDuration: cliArguments.autoFailureDuration,
	}, nil
}
