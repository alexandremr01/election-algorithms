package state

import (
	"time"
)

type State struct {
	CoordinatorID int
	LastHearbeat *time.Time
}

func NewState() *State {
	return &State{CoordinatorID: -1, LastHearbeat: nil}
}