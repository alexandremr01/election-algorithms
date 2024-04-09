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

type Server struct {
	NodeID string 
	LastHearbeat *time.Time
	Clients map[string]*rpc.Client
	ElectionHadResponse bool
}

type HearbeatArgs struct {Sender string}
type ElectionArgs struct {Sender string}
type RespondElectionArgs struct {Sender string}
func (s *Server) SendHeartbeat(args *HearbeatArgs, reply *int64) error {
    log.Printf("Node %s: Received heartbeat from node %s\n", s.NodeID, args.Sender)
	now := time.Now()
	s.LastHearbeat = &now
    return nil
}

func (s *Server) RespondElection(args *RespondElectionArgs, reply *int64) error {
    log.Printf("Node %s: Received OK from node %s\n", s.NodeID, args.Sender)
	s.ElectionHadResponse = true
    return nil
}

func (s *Server) CallForElection(args *ElectionArgs, reply *int64) error {
    log.Printf("Node %s: Received call for elections from node %s\n", s.NodeID, args.Sender)
	client, ok := s.Clients[args.Sender]
	if !ok {
		log.Printf("Node %s: Node %s not connected\n", s.NodeID, args.Sender)
		return nil
	}
	err := client.Call(
		"Server.RespondElection", 
		RespondElectionArgs{Sender: s.NodeID},
		nil,
	)
	if err != nil {
		log.Printf("Node %s: Error communicating with node %s: %s", s.NodeID, args.Sender, err)
	}
	// TODO: initiates an election
    return nil
}

func main() {
	ids := [3]string{"1", "2", "3"}
	numIds := [3]int{1, 2, 3}

	port := os.Getenv("PORT")
	nodeID := os.Getenv("NODE_ID")
	myID, err := strconv.Atoi(nodeID) 
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


    log.Printf("My ID: %s\n", nodeID)
	clients := make(map[string]*rpc.Client)
	server := &Server{NodeID: nodeID, Clients: clients}

	go func() {
		time.Sleep(heartbeatTimeDuration)
		for _, id := range ids {
			if id == nodeID {
				continue
			}
			hostname := fmt.Sprintf("p%s:%s", id, port)
			client, err := rpc.DialHTTP("tcp", hostname)
			if err != nil {
				log.Fatal("dialing:", err)
			}
			clients[id] = client
		}

		if (nodeID == "3") {
			for {
				time.Sleep(2* time.Second)
				for _, id := range ids {
					if id == nodeID {
						continue
					}
					err := clients[id].Call(
						"Server.SendHeartbeat", 
						HearbeatArgs{Sender: nodeID},
						nil,
					)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		} else {
			for {
				time.Sleep(timeoutDuration)
				if (server.LastHearbeat == nil) || (time.Now().Sub(*server.LastHearbeat) > timeoutDuration) {
					log.Printf("Leader timed out")
					server.ElectionHadResponse = false
					for _, numID := range numIds {
						id := strconv.Itoa(numID)
					
						if id == nodeID || numID < myID {
							continue
						}
						err := clients[id].Call(
							"Server.CallForElection", 
							ElectionArgs{Sender: nodeID},
							nil,
						)
						if err != nil {
							log.Printf("Node %s: Error communicating with node %s: %s", nodeID, id, err)
						}
					}
					time.Sleep(electionDuration)
					if server.ElectionHadResponse {
						log.Printf("Node %s: Election finished with responses, going back to normal.", nodeID)
					} else {
						// TODO: becomes leader
						log.Printf("Node %s: Election finished without responses, becoming leader.", nodeID)
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
