package raft

import (
	"log"
	"github.com/alexandremr01/user-elections/client"
	"github.com/alexandremr01/user-elections/state"
	"time"
)
 
type Elections struct {
	CurrentTerm int
	VotedFor int
	VotesCount int 
	Happening bool

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
	e.Happening = true
	e.CurrentTerm += 1
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
			e.VotesCount += 1
		}
	}
	
	time.Sleep(e.electionDuration)
	if ! e.Happening {
		log.Printf("Finished aborted election")
		return		
	} 
	
	if e.VotesCount > len(e.ids) / 2 {
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