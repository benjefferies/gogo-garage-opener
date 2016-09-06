package main

import (
	log "github.com/Sirupsen/logrus"
)

type NoopDoorController struct {

}


func (this NoopDoorController) toggleDoor() {
	log.Info("Noop toggleDoor")
}

func (this NoopDoorController) getDoorState() DoorState {
	log.Info("Noop getDoorState")
	return closed
}


func (this NoopDoorController) close() {
	log.Info("Noop close")
}
