package server

import (
	"time"
	"log"
	"github.com/alexandremr01/user-elections/client"
	"github.com/alexandremr01/user-elections/algorithms"
	"github.com/alexandremr01/user-elections/messages"
)

type Server struct {
	NodeID int 
	LastHearbeat *time.Time
	Client *client.Client
	Elections *algorithms.Elections
}

func NewServer(nodeID int, client *client.Client, elections *algorithms.Elections) *Server {
	return &Server{NodeID: nodeID, Client: client, Elections: elections}
}

func (s *Server) SendHeartbeat(args *messages.HearbeatArgs, reply *int64) error {
    log.Printf("Node %d: Received heartbeat from node %d\n", s.NodeID, args.Sender)
	now := time.Now()
	s.LastHearbeat = &now
    return nil
}

func (s *Server) RespondElection(args *messages.RespondElectionArgs, reply *int64) error {
    log.Printf("Node %d: Received OK from node %d\n", s.NodeID, args.Sender)
	s.Elections.Answered = true
    return nil
}

func (s *Server) NotifyNewCoordinator(args *messages.NotifyNewCoordinatorArgs, reply *int64) error {
    log.Printf("Node %d: Received NewCoordinator from node %d\n", s.NodeID, args.Sender)
	s.Elections.CoordinatorID = args.Sender
    return nil
}

func (s *Server) CallForElection(args *messages.ElectionArgs, reply *int64) error {
    log.Printf("Node %d: Received call for elections from node %d\n", s.NodeID, args.Sender)
	s.Client.Send(
		args.Sender,
		"Server.RespondElection", 
		messages.RespondElectionArgs{Sender: s.NodeID},
	)
	if !s.Elections.Happening{
		s.Elections.Start()
	}
    return nil
}
