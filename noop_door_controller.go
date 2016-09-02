package main

import (
	log "github.com/Sirupsen/logrus"
)

type NoopDoorController struct {

}


func (c NoopDoorController) toggleDoor() {
	log.Info("Noop toggleDoor")
}

func (c NoopDoorController) getDoorState() DoorState {
	log.Info("Noop getDoorState")
	return closed
}
