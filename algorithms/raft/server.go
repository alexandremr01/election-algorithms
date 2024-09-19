package raft

import (
	"fmt"
	"log"
	"time"

	"github.com/alexandremr01/user-elections/client"
	"github.com/alexandremr01/user-elections/types"
)

type Server struct {
	NodeID    int
	Client    *client.Client
	Elections *Elections
	state     *types.State
}

type AppendEntriesArgs struct {
	Term   int
	Sender int
}

func NewServer(nodeID int, client *client.Client, elections *Elections, state *types.State) *Server {
	return &Server{NodeID: nodeID, Client: client, Elections: elections, state: state}
}

func (s *Server) AppendEntries(args *AppendEntriesArgs, _ *int64) error {
	message := ""
	if args.Term >= s.Elections.CurrentTerm {
		s.Elections.interruptElection() // this will step down if candidate
		now := time.Now()
		s.state.LastHearbeat = &now
		s.state.CoordinatorID = args.Sender
		message += fmt.Sprintf("Node %d: Received heartbeat from node %d", s.NodeID, args.Sender)
	}
	if args.Term > s.Elections.CurrentTerm {
		s.Elections.updateTerm(args.Term)
		message += fmt.Sprintf("Updated term to %d", args.Term)
	}
	if args.Term < s.Elections.CurrentTerm {
		message = fmt.Sprintf(
			"Rejecting AppendEntries from outdated leader %d of term %d",
			args.Sender, args.Term,
		)
	}
	message += "\n"
	log.Print(message)
	return nil
}

type RequestVoteArgs struct {
	Sender int
	Term   int
}
type RequestVoteResponse struct{ VoteGranted bool }

func (s *Server) RequestVote(args *RequestVoteArgs, reply *RequestVoteResponse) error {
	if args.Term > s.Elections.CurrentTerm {
		s.Elections.updateTerm(args.Term)
	}
	willVote := (args.Term == s.Elections.CurrentTerm) && (s.Elections.VotedFor == -1)
	log.Printf(
		"Node %d: Received vote request from node %d for term %d, voting %t\n",
		s.NodeID, args.Sender, args.Term, willVote,
	)
	if willVote {
		s.Elections.VotedFor = args.Sender
	}
	// received signal from candidate: can reset timer
	if args.Term >= s.Elections.CurrentTerm {
		now := time.Now()
		s.state.LastHearbeat = &now
	}
	*reply = RequestVoteResponse{
		VoteGranted: willVote,
	}
	return nil
}
