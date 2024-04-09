package client

import (
	"fmt"
	"net/rpc"
)

type Client struct {
	clients map[int]*rpc.Client
	nodeID int
}

func NewClient(nodeID int) *Client{
	return &Client{
		nodeID: nodeID,
		clients: make(map[int]*rpc.Client),
	}
}

func (c *Client) Init(ids []int) {
	for _, id := range ids {
		if id == c.nodeID {
			continue
		}
		hostname := fmt.Sprintf("p%d:8000", id)
		client, _ := rpc.DialHTTP("tcp", hostname)
		c.clients[id] = client // can be nil
	}
}

func (c *Client) Broadcast(ids []int, serviceMethod string, args any){
	for _, id := range ids {
		if id == c.nodeID {
			continue
		}
		c.Send(id, serviceMethod, args)
	}
}

func (c *Client) Send(id int, serviceMethod string, args any) {
	// tries to connect - not guaranteed
	if c.clients[id] == nil {
		hostname := fmt.Sprintf("p%d:8000", id)
		client, _ := rpc.DialHTTP("tcp", hostname)
		c.clients[id] = client
	}
	if c.clients[id] != nil {
		_ =  c.clients[id].Call(serviceMethod, args, nil)
	}
}