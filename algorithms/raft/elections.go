package raft

import (
	"log"
	"time"

	"github.com/alexandremr01/user-elections/algorithms/types"
	"github.com/alexandremr01/user-elections/client"
	"github.com/alexandremr01/user-elections/config"
	"github.com/alexandremr01/user-elections/state"
)

type Elections struct {
	CurrentTerm int
	VotedFor    int
	VotesCount  int
	Happening   bool

	connection *client.Client

	nodeID           int
	electionDuration time.Duration
	ids              []int
	state            *state.State
	server           *Server
}

func NewElections(conf *config.Config, state *state.State, connection *client.Client) types.Algorithm {
	alg := &Elections{
		CurrentTerm:      0,
		VotedFor:         -1,
		ids:              conf.IDs,
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

// RAFT has no action on startup: it will follow the current leader
func (e *Elections) InitializeNode() {

}

func (e *Elections) StartElections() {
	e.Happening = true
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

	time.Sleep(e.electionDuration)
	if !e.Happening {
		log.Printf("Finished aborted election")
		return
	}

	if e.VotesCount > len(e.ids)/2 {
		e.state.CoordinatorID = e.nodeID
		log.Printf("Node %d: Election finished with victory, becoming leader.", e.nodeID)
		e.SendHeartbeat()
	} else {
		log.Printf("Node %d: Election finished without a victory.", e.nodeID)
	}
	e.Interrupt()
}

func (e *Elections) SendHeartbeat() {
	e.connection.Broadcast(
		e.ids,
		"Server.AppendEntries",
		AppendEntriesArgs{Sender: e.nodeID, Term: e.CurrentTerm},
	)
}

func (e *Elections) Interrupt() {
	e.Happening = false
	e.VotedFor = -1
}

func (e *Elections) GetServer() any {
	return e.server
}
