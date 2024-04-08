package main

import (
	"fmt"
	"os"
	"log"
	"net/rpc"
	"net/http"
	"net"
	"time"
)

type Server struct { }

type HearbeatArgs struct {Sender string}
func (t *Server) SendHeartbeat(args *HearbeatArgs, reply *int64) error {
    fmt.Printf("Received heartbeat from %s\n", args.Sender)
    return nil
}

func main() {
	ids := [3]string{"1", "2", "3"}

	port := os.Getenv("PORT")
	nodeID := os.Getenv("NODE_ID")

    fmt.Printf("My IDs: %s\n", nodeID)

	go func() {
		clients := make(map[string]*rpc.Client)
		time.Sleep(2* time.Second)
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
	
	server := new(Server)
	rpc.Register(server)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", ":"+port)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}
