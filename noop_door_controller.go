package main

import (
	log "github.com/Sirupsen/logrus"
)

// NoopDoorController for testing
type NoopDoorController struct {
}

// toggleDoor is noop for testing
func (noopDoorController NoopDoorController) toggleDoor() {
	log.Info("Noop toggleDoor")
}

// getDoorState is noop for testing always returning closed
func (noopDoorController NoopDoorController) getDoorState() DoorState {
	log.Info("Noop getDoorState")
	return closed
}

// close is noop for testing
func (noopDoorController NoopDoorController) close() {
	log.Info("Noop close")
}
