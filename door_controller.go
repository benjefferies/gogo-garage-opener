package main

// DoorState is either 0 or 1 depending on if it is open or closed
type DoorState uint8

const (
	closed DoorState = 0
	open DoorState = 1
)

// DoorController abstraction on garage door functionality for testing
type DoorController interface {
	toggleDoor()
	getDoorState() DoorState
	close()
}