package main

func newState(state DoorState) *State {
	var description = "Open"
	if state == closed {
		description = "Closed"
	}
	return &State{State: state, Description: description}
}

// State holds reference to the door state and description
type State struct {
	State       DoorState
	Description string
}
