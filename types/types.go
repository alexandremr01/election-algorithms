package types

import "time"

type Algorithm interface {
	InitializeNode()
	StartElections()
	SendHeartbeat()
	GetServer() any
}

type Config struct {
	TimeoutDuration   time.Duration
	ElectionDuration  time.Duration
	HeartbeatDuration time.Duration
	Port              string
	Addresses         map[int]string
	IDs               []int
	NodeID int
	AlgorithmName string
}