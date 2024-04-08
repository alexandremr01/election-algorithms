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

type Server struct {NodeID string}

type HearbeatArgs struct {Sender string}
func (s *Server) SendHeartbeat(args *HearbeatArgs, reply *int64) error {
    log.Printf("Node %s: Received heartbeat from node %s\n", s.NodeID, args.Sender)
    return nil
}

func main() {
	ids := [3]string{"1", "2", "3"}

	port := os.Getenv("PORT")
	nodeID := os.Getenv("NODE_ID")
	heartbeatTime, err := strconv.Atoi(os.Getenv("HEARTBEAT_TIME"))
	if err != nil {
		log.Fatal("error parsing heartbeat time:", err)
	}

    log.Printf("My ID: %s\n", nodeID)

	go func() {
		clients := make(map[string]*rpc.Client)
		time.Sleep(time.Duration(heartbeatTime) * time.Second)
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
		}
    }()
	
	server := &Server{NodeID: nodeID}
	rpc.Register(server)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", ":"+port)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}
