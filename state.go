package main

func newState(state int8) *State {
	var description string = "Open"
	if state == 0 {
		description = "Closed"
	}
	return &State{State: state, Description: description}
}

type State struct {
	State       int8
	Description string
}
