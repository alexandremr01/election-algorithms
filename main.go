package main

import (
	"log"
	"net/rpc"
	"net/http"
	"net"
	"time"

	"github.com/alexandremr01/user-elections/config"
	"github.com/alexandremr01/user-elections/server"
	"github.com/alexandremr01/user-elections/client"
	"github.com/alexandremr01/user-elections/algorithms"
	"github.com/alexandremr01/user-elections/messages"
)

func main() {
	ids := []int{1, 2, 3, 4}

	config, err := config.GetConfig()
	if err != nil {
		log.Fatal("error parsing config: ", err)
	}
    log.Printf("My ID: %d\n", config.NodeID)
	
	leader := -1
	for _, id := range ids {
		if id > leader {
			leader = id
		}
	}

	connection := client.NewClient(config.NodeID)
	elections := algorithms.NewElections(ids, config.NodeID, leader, connection, config.ElectionDuration)
	server := server.NewServer(config.NodeID, connection, elections)

	go func() {
		time.Sleep(config.HeartbeatDuration)
		connection.Init(ids[:])
		elections.Start()
		for {
			if (config.NodeID == elections.CoordinatorID) {
				time.Sleep(config.HeartbeatDuration)
				connection.Broadcast(ids[:], "Server.SendHeartbeat", messages.HearbeatArgs{Sender: config.NodeID})
			} else {
				time.Sleep(config.TimeoutDuration)
				if (server.LastHearbeat == nil) || (time.Now().Sub(*server.LastHearbeat) > config.TimeoutDuration) {
					log.Printf("Leader timed out")
					elections.Start()					
				}			
			}
		}
    }()
	
	rpc.Register(server)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", ":"+config.Port)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}
