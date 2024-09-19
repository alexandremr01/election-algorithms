package types

import "time"

type Algorithm interface {
	OnInitialization()
	OnLeaderTimeout()
	SendHeartbeat()
	GetServer() any
}

type Config struct {
	TimeoutDuration     time.Duration
	ElectionDuration    time.Duration
	HeartbeatDuration   time.Duration
	Port                string
	Addresses           map[int]string
	IDs                 []int
	NodeID              int
	AlgorithmName       string
	AutoFailure         int
	AutoFailureDuration time.Duration
}

type State struct {
	CoordinatorID int
	LastHearbeat  *time.Time
}

func NewState() *State {
	return &State{CoordinatorID: -1, LastHearbeat: nil}
}
