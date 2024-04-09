package algorithms

import (
	"log"
	"github.com/alexandremr01/user-elections/client"
	"github.com/alexandremr01/user-elections/messages"
	"time"
)
 
type Elections struct {
	Happening bool
	Answered bool
	CoordinatorID int

	connection *client.Client

	nodeID int
	electionDuration time.Duration
	higherIds []int
	ids []int
}

func NewElections(ids []int, nodeID int, coordinatorID int, connection *client.Client, electionDuration time.Duration) *Elections {
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
		CoordinatorID: coordinatorID,
		connection: connection,
		electionDuration: electionDuration,
	}
}

func (e *Elections) Start() {
	e.Answered = false
	e.Happening = true

	e.connection.Broadcast(e.higherIds[:], "Server.CallForElection", messages.ElectionArgs{Sender: e.nodeID})
	time.Sleep(e.electionDuration)
	if e.Answered {
		log.Printf("Node %d: Election finished with responses, going back to normal.", e.nodeID)
	} else {
		e.CoordinatorID = e.nodeID
		e.connection.Broadcast(e.ids[:], "Server.NotifyNewCoordinator", messages.NotifyNewCoordinatorArgs{Sender: e.nodeID})
		log.Printf("Node %d: Election finished without responses, becoming leader.", e.nodeID)
	}

	e.Happening = false
}

