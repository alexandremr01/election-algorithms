package bully

import (
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

type HearbeatArgs struct{ Sender int }
type ElectionArgs struct{ Sender int }
type RespondElectionArgs struct{ Sender int }
type NotifyNewCoordinatorArgs struct{ Sender int }

func NewServer(nodeID int, client *client.Client, elections *Elections, state *types.State) *Server {
	return &Server{NodeID: nodeID, Client: client, Elections: elections, state: state}
}

func (s *Server) SendHeartbeat(args *HearbeatArgs, _ *int64) error {
	log.Printf("Node %d: Received heartbeat from node %d\n", s.NodeID, args.Sender)
	now := time.Now()
	s.state.LastHearbeat = &now
	return nil
}

func (s *Server) RespondElection(args *RespondElectionArgs, _ *int64) error {
	log.Printf("Node %d: Received OK from node %d\n", s.NodeID, args.Sender)
	s.Elections.Answered = true
	return nil
}

func (s *Server) NotifyNewCoordinator(args *NotifyNewCoordinatorArgs, _ *int64) error {
	log.Printf("Node %d: Received NewCoordinator from node %d\n", s.NodeID, args.Sender)
	s.state.CoordinatorID = args.Sender
	return nil
}

func (s *Server) CallForElection(args *ElectionArgs, _ *int64) error {
	log.Printf("Node %d: Received call for elections from node %d\n", s.NodeID, args.Sender)
	s.Client.Send(
		args.Sender,
		"Server.RespondElection",
		RespondElectionArgs{Sender: s.NodeID},
		nil,
	)
	if !s.Elections.Happening {
		s.Elections.StartElections()
	}
	return nil
}
