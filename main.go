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
	Clients map[int]*rpc.Client
	ElectionHadResponse bool
	Coordinator *Coordinator
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
	s.ElectionHadResponse = true
    return nil
}

func (s *Server) NotifyNewCoordinator(args *NotifyNewCoordinatorArgs, reply *int64) error {
    log.Printf("Node %d: Received NewCoordinator from node %d\n", s.NodeID, args.Sender)
	s.Coordinator.ID = args.Sender
    return nil
}

func (s *Server) CallForElection(args *ElectionArgs, reply *int64) error {
    log.Printf("Node %d: Received call for elections from node %d\n", s.NodeID, args.Sender)
	client, ok := s.Clients[args.Sender]
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
    return nil
}

func main() {
	ids := [3]int{1, 2, 3}
	leader := -1
	for _, id := range ids {
		if id > leader {
			leader = id
		}
	}
	coordinator := &Coordinator{ID: leader}

	port := os.Getenv("PORT")
	nodeID, err := strconv.Atoi(os.Getenv("NODE_ID")) 
	if err != nil {
		log.Fatal("error parsing node id:", err)
	}

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
	clients := make(map[int]*rpc.Client)
	server := &Server{NodeID: nodeID, Clients: clients, Coordinator: coordinator}

	go func() {
		time.Sleep(heartbeatTimeDuration)
		for _, id := range ids {
			if id == nodeID {
				continue
			}
			hostname := fmt.Sprintf("p%d:%s", id, port)
			client, err := rpc.DialHTTP("tcp", hostname)
			if err != nil {
				log.Fatal("dialing:", err)
			}
			clients[id] = client
		}

		for {
			if (nodeID == coordinator.ID) {
				time.Sleep(2* time.Second)
				for _, id := range ids {
					if id == nodeID {
						continue
					}
					_ =  clients[id].Call(
						"Server.SendHeartbeat", 
						HearbeatArgs{Sender: nodeID},
						nil,
					)
				}
			} else {
				time.Sleep(timeoutDuration)
				if (server.LastHearbeat == nil) || (time.Now().Sub(*server.LastHearbeat) > timeoutDuration) {
					log.Printf("Leader timed out")
					server.ElectionHadResponse = false
					for _, id := range ids {					
						if id <= nodeID {
							continue
						}
						err := clients[id].Call(
							"Server.CallForElection", 
							ElectionArgs{Sender: nodeID},
							nil,
						)
						if err != nil {
							log.Printf("Node %d: Error communicating with node %d: %s", nodeID, id, err)
						}
					}
					time.Sleep(electionDuration)
					if server.ElectionHadResponse {
						log.Printf("Node %d: Election finished with responses, going back to normal.", nodeID)
					} else {
						coordinator.ID = nodeID
						for _, id := range ids {					
							if id == nodeID {
								continue
							}
							err := clients[id].Call(
								"Server.NotifyNewCoordinator", 
								NotifyNewCoordinatorArgs{Sender: nodeID},
								nil,
							)
							if err != nil {
								log.Printf("Node %d: Error communicating with node %d: %s", nodeID, id, err)
							}
						}
						log.Printf("Node %d: Election finished without responses, becoming leader.", nodeID)
					}
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
