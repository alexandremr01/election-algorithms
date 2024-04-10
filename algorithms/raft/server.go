package raft

import (
	"time"
	"fmt"
	"log"
	"github.com/alexandremr01/user-elections/client"
	"github.com/alexandremr01/user-elections/state"
)

type Server struct {
	NodeID int 
	LastHearbeat *time.Time
	Client *client.Client
	Elections *Elections
	state *state.State
}

type HearbeatArgs struct {
	Term int
	Sender int
}

func NewServer(nodeID int, client *client.Client, elections *Elections, state *state.State) *Server {
	return &Server{NodeID: nodeID, Client: client, Elections: elections, state: state}
}

func (s *Server) SendHeartbeat(args *HearbeatArgs, reply *int64) error {
	message := ""
	if args.Term >= s.Elections.CurrentTerm {
		now := time.Now()
		s.LastHearbeat = &now
		message += fmt.Sprintf("Node %d: Received heartbeat from node %d", s.NodeID, args.Sender)
	}
	if args.Term > s.Elections.CurrentTerm {
		s.Elections.CurrentTerm = args.Term
		s.state.CoordinatorID = args.Sender
		message += fmt.Sprintf("Updated term to %d", args.Term)
	}
	message += "\n"
	log.Printf(message)
    return nil
}

// func (s *Server) RequestVote(args *messages.ElectionArgs, reply *int64) error {
//     log.Printf("Node %d: Received call for elections from node %d\n", s.NodeID, args.Sender)
// 	s.Client.Send(
// 		args.Sender,
// 		"Server.RespondElection", 
// 		messages.RespondElectionArgs{Sender: s.NodeID},
// 	)
// 	if !s.Elections.Happening{
// 		s.Elections.StartElections()
// 	}
//     return nil
// }
