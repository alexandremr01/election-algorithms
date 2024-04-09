package main

import (
	"fmt"
	"os"
	"log"
	"net/rpc"
	"net/http"
	"net"
	"time"
	"strconv"
)

type Coordinator struct {
	ID int
}

type Server struct {
	NodeID int 
	LastHearbeat *time.Time
	Connection *Connection
	Coordinator *Coordinator
	Elections *Elections
}

type HearbeatArgs struct {Sender int}
type ElectionArgs struct {Sender int}
type RespondElectionArgs struct {Sender int}
type NotifyNewCoordinatorArgs struct {Sender int}
func (s *Server) SendHeartbeat(args *HearbeatArgs, reply *int64) error {
    log.Printf("Node %d: Received heartbeat from node %d\n", s.NodeID, args.Sender)
	now := time.Now()
	s.LastHearbeat = &now
    return nil
}

func (s *Server) RespondElection(args *RespondElectionArgs, reply *int64) error {
    log.Printf("Node %d: Received OK from node %d\n", s.NodeID, args.Sender)
	s.Elections.Answered = true
    return nil
}

func (s *Server) NotifyNewCoordinator(args *NotifyNewCoordinatorArgs, reply *int64) error {
    log.Printf("Node %d: Received NewCoordinator from node %d\n", s.NodeID, args.Sender)
	s.Coordinator.ID = args.Sender
    return nil
}

func (s *Server) CallForElection(args *ElectionArgs, reply *int64) error {
    log.Printf("Node %d: Received call for elections from node %d\n", s.NodeID, args.Sender)
	client, ok := s.Connection.Clients[args.Sender]
	if !ok {
		log.Printf("Node %d: Node %d not connected\n", s.NodeID, args.Sender)
		return nil
	}
	err := client.Call(
		"Server.RespondElection", 
		RespondElectionArgs{Sender: s.NodeID},
		nil,
	)
	if err != nil {
		log.Printf("Node %d: Error communicating with node %d: %s", s.NodeID, args.Sender, err)
	}
	// TODO: initiates an election
	s.Elections.Start()
    return nil
}

type Elections struct {
	Happening bool
	Answered bool

	coordinator *Coordinator
	connection *Connection

	nodeID int
	electionDuration time.Duration
	higherIds []int
	ids []int
}

func NewElections(ids []int, nodeID int, coordinator *Coordinator, connection *Connection, electionDuration time.Duration) *Elections {
	var higherIds []int
	for _, id := range ids {
		if id > nodeID {
			higherIds = append(higherIds, id)
		}
	}
	return &Elections{
		Happening: false,
		Answered: false,
		higherIds: higherIds,
		ids: ids,
		nodeID: nodeID,
		coordinator: coordinator,
		connection: connection,
		electionDuration: electionDuration,
	}
}

func (e *Elections) Start() {
	e.Answered = false
	e.Happening = true

	e.connection.Broadcast(e.higherIds[:], "Server.CallForElection", ElectionArgs{Sender: e.nodeID})
	time.Sleep(e.electionDuration)
	if e.Answered {
		log.Printf("Node %d: Election finished with responses, going back to normal.", e.nodeID)
	} else {
		e.coordinator.ID = e.nodeID
		e.connection.Broadcast(e.ids[:], "Server.NotifyNewCoordinator", NotifyNewCoordinatorArgs{Sender: e.nodeID})
		log.Printf("Node %d: Election finished without responses, becoming leader.", e.nodeID)
	}

	e.Happening = false
}

type Connection struct {
	Clients map[int]*rpc.Client
	nodeID int
}
func NewConnection(nodeID int) *Connection{
	return &Connection{
		nodeID: nodeID,
		Clients: make(map[int]*rpc.Client),
	}
}

func (c *Connection) Init(ids []int) {
	for _, id := range ids {
		if id == c.nodeID {
			continue
		}
		hostname := fmt.Sprintf("p%d:8000", id)
		client, err := rpc.DialHTTP("tcp", hostname)
		if err != nil {
			log.Fatal("dialing:", err)
		}
		c.Clients[id] = client
	}
}

func (c *Connection) Broadcast(ids []int, serviceMethod string, args any){
	for _, id := range ids {
		if id == c.nodeID {
			continue
		}
		_ =  c.Clients[id].Call(serviceMethod, args, nil)
	}
}

func main() {
	ids := []int{1, 2, 3}
	nodeID, err := strconv.Atoi(os.Getenv("NODE_ID")) 
	if err != nil {
		log.Fatal("error parsing node id:", err)
	}

	leader := -1
	for _, id := range ids {
		if id > leader {
			leader = id
		}
	}
	coordinator := &Coordinator{ID: leader}

	port := os.Getenv("PORT")
	
	timeout, err := strconv.Atoi(os.Getenv("NODE_TIMEOUT"))
	if err != nil {
		log.Fatal("error parsing timeout:", err)
	}
	timeoutDuration := time.Duration(timeout)*time.Millisecond

	electionTimeout, err := strconv.Atoi(os.Getenv("ELECTION_TIMEOUT"))
	if err != nil {
		log.Fatal("error parsing election timeout:", err)
	}
	electionDuration := time.Duration(electionTimeout)*time.Millisecond

	heartbeatTime, err := strconv.Atoi(os.Getenv("HEARTBEAT_TIME"))
	if err != nil {
		log.Fatal("error parsing heartbeat time:", err)
	}
	heartbeatTimeDuration := time.Duration(heartbeatTime)*time.Second


    log.Printf("My ID: %d\n", nodeID)

	connection := NewConnection(nodeID)
	elections := NewElections(ids, nodeID, coordinator, connection, electionDuration)

	server := &Server{NodeID: nodeID, Connection: connection, Coordinator: coordinator, Elections: elections}

	go func() {
		time.Sleep(heartbeatTimeDuration)
		connection.Init(ids[:])

		for {
			if (nodeID == coordinator.ID) {
				time.Sleep(2* time.Second)
				connection.Broadcast(ids[:], "Server.SendHeartbeat", HearbeatArgs{Sender: nodeID})
			} else {
				time.Sleep(timeoutDuration)
				if (server.LastHearbeat == nil) || (time.Now().Sub(*server.LastHearbeat) > timeoutDuration) {
					log.Printf("Leader timed out")
					elections.Start()					
				}			
			}
		}
    }()
	
	rpc.Register(server)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", ":"+port)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}
