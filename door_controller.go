package main

type DoorState uint8

const (
	closed DoorState = 0
	open DoorState = 1
)

type DoorController interface {
	toggleDoor()
	getDoorState() DoorState
	close()
}