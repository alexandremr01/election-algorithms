package algorithms

import (
	"errors"

	"github.com/alexandremr01/user-elections/algorithms/bully"
	"github.com/alexandremr01/user-elections/algorithms/raft"
	"github.com/alexandremr01/user-elections/client"
	"github.com/alexandremr01/user-elections/state"
	"github.com/alexandremr01/user-elections/types"
)

type AlgorithmBuilder func(*types.Config, *state.State, *client.Client) types.Algorithm

var algorithmBuilders = map[string]AlgorithmBuilder{
	"raft":  raft.NewElections,
	"bully": bully.NewElections,
}

func GetAlgorithmsList() []string {
	algorithms := []string{}
	for k := range algorithmBuilders {
		algorithms = append(algorithms, k)
	}
	return algorithms
}

func GetAlgorithm(algorithmName string,
	conf *types.Config,
	state *state.State,
	client *client.Client) (types.Algorithm, error) {
	builder, ok := algorithmBuilders[algorithmName]
	if !ok {
		return nil, errors.New("algorithm not registered")
	}
	return builder(conf, state, client), nil
}
