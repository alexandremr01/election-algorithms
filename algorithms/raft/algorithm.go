package raft

import (
	"log"
	"time"

	"github.com/alexandremr01/user-elections/client"
	"github.com/alexandremr01/user-elections/types"
)

type NodeStatus string

const (
	CandidateStatus NodeStatus = "CandidateStatus"
	FollowerStatus  NodeStatus = "FollowerStatus"
	LeaderStatus    NodeStatus = "LeaderStatus"
)

type Elections struct {
	CurrentTerm int
	VotedFor    int
	VotesCount  int
	Status      NodeStatus

	connection *client.Client

	nodeID           int
	electionDuration time.Duration
	ids              []int
	state            *types.State
	server           *Server
}

func NewElections(conf *types.Config, state *types.State, connection *client.Client) types.Algorithm {
	alg := &Elections{
		CurrentTerm:      0,
		VotedFor:         -1,
		ids:              conf.IDs,
		nodeID:           conf.NodeID,
		state:            state,
		connection:       connection,
		electionDuration: conf.ElectionDuration,
		server:           nil,
		Status:           FollowerStatus,
	}
	server := NewServer(conf.NodeID, connection, alg, state)
	alg.server = server
	return alg
}

// RAFT has no action on startup: it will follow the current leader
func (e *Elections) OnInitialization() {

}

func (e *Elections) OnLeaderTimeout() {
	e.Status = CandidateStatus
	e.CurrentTerm++
	e.VotesCount = 1
	e.VotedFor = e.nodeID

	for _, id := range e.ids {
		if id == e.nodeID {
			continue
		}
		var resp RequestVoteResponse
		e.connection.Send(id,
			"Server.RequestVote",
			RequestVoteArgs{Sender: e.nodeID, Term: e.CurrentTerm},
			&resp,
		)
		log.Printf("Response from %d: Vote granted: %t", id, resp.VoteGranted)
		if resp.VoteGranted {
			e.VotesCount++
		}
	}

	// case 1: received an append entries during elections
	if e.Status != CandidateStatus {
		log.Printf("Node %d: Finished aborted election", e.nodeID)
		return
	}

	// case 2: win
	if e.VotesCount > len(e.ids)/2 {
		e.state.CoordinatorID = e.nodeID
		log.Printf("Node %d: Election finished with victory, becoming leader.", e.nodeID)
		e.SendHeartbeat()
		e.Status = LeaderStatus
	} else {
		// case 3: no victory
		log.Printf("Node %d: Election finished without a victory.", e.nodeID)
		e.Status = FollowerStatus
	}
	e.VotedFor = -1
}

func (e *Elections) SendHeartbeat() {
	e.connection.Broadcast(
		e.ids,
		"Server.AppendEntries",
		AppendEntriesArgs{Sender: e.nodeID, Term: e.CurrentTerm},
	)
}

func (e *Elections) GetServer() any {
	return e.server
}

func (e *Elections) updateTerm(newTerm int) {
	e.CurrentTerm = newTerm
	e.VotedFor = -1
}
