package client

import (
	"errors"
	"log"
	"net/rpc"
)

type Client struct {
	clients   map[int]*rpc.Client
	addresses map[int]string
	nodeID    int
}

func NewClient(nodeID int, addresses map[int]string) *Client {
	return &Client{
		nodeID:    nodeID,
		clients:   make(map[int]*rpc.Client),
		addresses: addresses,
	}
}

func (c *Client) Broadcast(ids []int, serviceMethod string, args any) {
	for _, id := range ids {
		if id == c.nodeID {
			continue
		}
		c.Send(id, serviceMethod, args, nil)
	}
}

func (c *Client) Send(id int, serviceMethod string, args any, resp any) {
	// tries to connect - not guaranteed
	if c.clients[id] == nil {
		client, _ := rpc.DialHTTP("tcp", c.addresses[id])
		// if err == nil {
		// 	log.Printf("Connection established with node %d at address %s", id, c.addresses[id])
		// }
		c.clients[id] = client
	}
	if c.clients[id] != nil {
		err := c.clients[id].Call(serviceMethod, args, resp)
		if errors.Is(err, rpc.ErrShutdown) {
			c.clients[id] = nil
		} else if err != nil {
			log.Print("error trying to talk to", id, err)
		}
	}
}
