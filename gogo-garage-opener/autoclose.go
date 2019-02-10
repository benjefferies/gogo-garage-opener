package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

// Autoclose is to auto close the garage door
type Autoclose struct {
	now            time.Time
	doorController DoorController
}

func (autoclose Autoclose) shouldClose() bool {
	now := autoclose.now
	shouldClose := time.Date(now.Year(), now.Month(), now.Day(), 22, 0, 0, 0, time.Local)
	canStayOpen := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, time.Local)
	if now.After(shouldClose) || now.Before(canStayOpen) {
		return true
	}
	return false
}

func (autoclose Autoclose) autoClose() bool {
	shouldClose := autoclose.shouldClose()
	state := autoclose.doorController.getDoorState()
	log.Infof("Autoclose shouldClose=%t state=%b", shouldClose, state)
	if shouldClose && state == open {
		autoclose.doorController.toggleDoor()
		return true
	}
	return false
}
