package main

// DoorState is either 0 or 1 depending on if it is open or closed
type DoorState uint8

const (
	closed DoorState = 0
	open   DoorState = 1
)

func (doorState DoorState) isOpen() bool {
	return doorState == open
}

func (doorState DoorState) isClosed() bool {
	return doorState == closed
}

func (doorState DoorState) description() string {
	if doorState.isClosed() {
		return "Closed"
	}
	return "Open"
}

func (doorState DoorState) value() int8 {
	return int8(doorState)
}

// DoorController abstraction on garage door functionality for testing
type DoorController interface {
	toggleDoor()
	getDoorState() DoorState
	close()
}
