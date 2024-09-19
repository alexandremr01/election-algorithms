package bully

import (
	"log"
	"time"

	"github.com/alexandremr01/user-elections/client"
	"github.com/alexandremr01/user-elections/types"
)

type Elections struct {
	Happening bool
	Answered  bool
	state     *types.State

	connection *client.Client

	nodeID           int
	electionDuration time.Duration
	higherIDs        []int
	ids              []int
	server           *Server
}

func NewElections(conf *types.Config, state *types.State, connection *client.Client) types.Algorithm {
	var higherIDs []int
	ids := conf.IDs
	for _, id := range ids {
		if id > conf.NodeID {
			higherIDs = append(higherIDs, id)
		}
	}
	alg := &Elections{
		Happening:        false,
		Answered:         false,
		higherIDs:        higherIDs,
		ids:              ids,
		nodeID:           conf.NodeID,
		state:            state,
		connection:       connection,
		electionDuration: conf.ElectionDuration,
		server:           nil,
	}
	server := NewServer(conf.NodeID, connection, alg, state)
	alg.server = server
	return alg
}

func (e *Elections) InitializeNode() {
	e.StartElections()
}

func (e *Elections) StartElections() {
	e.Answered = false
	e.Happening = true

	e.connection.Broadcast(e.higherIDs, "Server.CallForElection", ElectionArgs{Sender: e.nodeID})
	time.Sleep(e.electionDuration)
	if e.Answered {
		log.Printf("Node %d: Election finished with responses, going back to normal.", e.nodeID)
	} else {
		e.state.CoordinatorID = e.nodeID
		e.connection.Broadcast(e.ids, "Server.NotifyNewCoordinator", NotifyNewCoordinatorArgs{Sender: e.nodeID})
		log.Printf("Node %d: Election finished without responses, becoming leader.", e.nodeID)
	}

	e.Happening = false
}

func (e *Elections) SendHeartbeat() {
	e.connection.Broadcast(e.ids, "Server.SendHeartbeat", HearbeatArgs{Sender: e.nodeID})
}

func (e *Elections) GetServer() any {
	return e.server
}
