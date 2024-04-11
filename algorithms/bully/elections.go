package bully

import (
	"log"
	"github.com/alexandremr01/user-elections/client"
	"github.com/alexandremr01/user-elections/state"
	"time"
)
 
type BullyElections struct {
	Happening bool
	Answered bool
	state *state.State

	connection *client.Client

	nodeID int
	electionDuration time.Duration
	higherIds []int
	ids []int
	server *Server
}

func NewElections(ids []int, nodeID int, state *state.State, connection *client.Client, electionDuration time.Duration) *BullyElections {
	var higherIds []int
	for _, id := range ids {
		if id > nodeID {
			higherIds = append(higherIds, id)
		}
	}
	alg := &BullyElections{
		Happening: false,
		Answered: false,
		higherIds: higherIds,
		ids: ids,
		nodeID: nodeID,
		state: state,
		connection: connection,
		electionDuration: electionDuration,
		server: nil,
	}
	server := NewServer(nodeID, connection, alg, state)
	alg.server = server
	return alg
}

func (e *BullyElections) InitializeNode() {
	e.StartElections()
}

func (e *BullyElections) StartElections() {
	e.Answered = false
	e.Happening = true

	e.connection.Broadcast(e.higherIds, "Server.CallForElection", ElectionArgs{Sender: e.nodeID})
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

func (e *BullyElections) SendHeartbeat() {
	e.connection.Broadcast(e.ids, "Server.SendHeartbeat", HearbeatArgs{Sender: e.nodeID})
}

func (e *BullyElections) GetServer() any {
	return e.server
}