package main

func newState(state DoorState) *State {
	var description string = "Open"
	if state == closed {
		description = "Closed"
	}
	return &State{State: state, Description: description}
}

type State struct {
	State       DoorState
	Description string
}
