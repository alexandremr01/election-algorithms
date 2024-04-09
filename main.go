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
}

type HearbeatArgs struct {Sender string}
func (s *Server) SendHeartbeat(args *HearbeatArgs, reply *int64) error {
    log.Printf("Node %s: Received heartbeat from node %s\n", s.NodeID, args.Sender)
	now := time.Now()
	s.LastHearbeat = &now
    return nil
}

func main() {
	ids := [3]string{"1", "2", "3"}

	port := os.Getenv("PORT")
	nodeID := os.Getenv("NODE_ID")
	timeout, err := strconv.Atoi(os.Getenv("NODE_TIMEOUT"))
	if err != nil {
		log.Fatal("error parsing timeout:", err)
	}
	timeoutDuration := time.Duration(timeout)*time.Millisecond
	heartbeatTime, err := strconv.Atoi(os.Getenv("HEARTBEAT_TIME"))
	if err != nil {
		log.Fatal("error parsing heartbeat time:", err)
	}
	heartbeatTimeDuration := time.Duration(heartbeatTime)*time.Second


    log.Printf("My ID: %s\n", nodeID)
	server := &Server{NodeID: nodeID}

	go func() {
		clients := make(map[string]*rpc.Client)
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
				if (server.LastHearbeat == nil) || (time.Now().Sub(*server.LastHearbeat) > timeoutDuration) {
					log.Printf("Leader timed out")
				}
				time.Sleep(timeoutDuration)
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
