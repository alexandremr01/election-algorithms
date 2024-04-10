package raft

import (
	// "log"
	"github.com/alexandremr01/user-elections/client"
	"github.com/alexandremr01/user-elections/state"
	"time"
)
 
type Elections struct {
	CurrentTerm int
	VotedFor int

	connection *client.Client

	nodeID int
	electionDuration time.Duration
	ids []int
	state *state.State
}

func NewElections(ids []int, nodeID int, state *state.State, connection *client.Client, electionDuration time.Duration) *Elections {

	return &Elections{
		CurrentTerm: 0,
		VotedFor: -1,
		ids: ids,
		nodeID: nodeID,
		state: state,
		connection: connection,
		electionDuration: electionDuration,
	}
}

// RAFT has no action on startup: it will follow the current leader
func (e *Elections) InitializeNode() {
	return
}


func (e *Elections) StartElections() {
	// e.CurrentTerm += 1
	// e.VotedFor = e.nodeID

	// e.connection.Broadcast(e.higherIds[:], "Server.RequestVote", RequestVoteArgs{Sender: e.nodeID})
	// time.Sleep(e.electionDuration)
	// if e.Won {
	// 	log.Printf("Node %d: Election finished with responses, going back to normal.", e.nodeID)
	// } else {
	// 	e.CoordinatorID = e.nodeID
	// 	e.connection.Broadcast(e.ids[:], "Server.NotifyNewCoordinator", messages.NotifyNewCoordinatorArgs{Sender: e.nodeID})
	// 	log.Printf("Node %d: Election finished without responses, becoming leader.", e.nodeID)
	// }
}

func (e *Elections) SendHeartbeat() {
	e.connection.Broadcast(e.ids, "Server.SendHeartbeat", HearbeatArgs{Sender: e.nodeID, Term: e.CurrentTerm})
}