package state

import (
	"time"
)

type State struct {
	CoordinatorID int
	LastHearbeat *time.Time
}

func NewState(leader int) *State {
	return &State{CoordinatorID: leader, LastHearbeat: nil}
}