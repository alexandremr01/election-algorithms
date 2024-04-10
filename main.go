package main

import (
	"log"
	"net/rpc"
	"net/http"
	"net"
	"time"

	"github.com/alexandremr01/user-elections/config"
	"github.com/alexandremr01/user-elections/client"
	"github.com/alexandremr01/user-elections/algorithms/bully"
	"github.com/alexandremr01/user-elections/algorithms/raft"
	"github.com/alexandremr01/user-elections/state"
)

type Algorithm interface{
	InitializeNode()
	StartElections()
	SendHeartbeat()
}

type Server interface {
	GetLastHearbeat() *time.Time
}

func main() {
	ids := []int{1, 2, 3, 4}

	config, err := config.GetConfig()
	if err != nil {
		log.Fatal("error parsing config: ", err)
	}
    log.Printf("My ID: %d\n", config.NodeID)

	algorithmName := "raft"
	connection := client.NewClient(config.NodeID)
	state := state.NewState()

	var algorithm Algorithm
	var server any
	if algorithmName == "bully"{
		bullyAlg := bully.NewBullyElections(ids, config.NodeID, state, connection, config.ElectionDuration)
		server = bully.NewServer(config.NodeID, connection, bullyAlg, state)
		algorithm = bullyAlg
	} else if algorithmName == "raft" {
		raftAlg := raft.NewElections(ids, config.NodeID, state, connection, config.ElectionDuration)
		server = raft.NewServer(config.NodeID, connection, raftAlg, state)
		algorithm = raftAlg
	} else {
		log.Fatalf("Unrecognized algorithm %s", algorithmName)
	}

	go func() {
		connection.Init(ids[:])
		algorithm.InitializeNode()
		for {
			if (config.NodeID == state.CoordinatorID) {
				algorithm.SendHeartbeat()
				time.Sleep(config.HeartbeatDuration)
			} else {
				time.Sleep(config.TimeoutDuration)
				if (state.LastHearbeat == nil) || (time.Now().Sub(*state.LastHearbeat) > config.TimeoutDuration) {
					log.Printf("Leader timed out")
					algorithm.StartElections()					
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
