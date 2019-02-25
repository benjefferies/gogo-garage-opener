package main

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Autoclose is to auto close the garage door
type Autoclose struct {
	lastOpen       time.Time
	doorController DoorController
}

// NewAutoclose AutoClose with default lastOpen set to max date
func NewAutoclose(doorcontroller DoorController) Autoclose {
	return Autoclose{lastOpen: time.Unix(1<<63-1, 0), doorController: doorcontroller}
}

func (autoclose Autoclose) shouldClose(now time.Time) bool {
	shouldClose := time.Date(now.Year(), now.Month(), now.Day(), 22, 0, 0, 0, time.Local)
	canStayOpen := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, time.Local)
	openTooLong := now.Sub(autoclose.lastOpen) > 2*time.Minute
	if openTooLong && (now.After(shouldClose) || now.Before(canStayOpen)) {
		return true
	}
	return false
}

func (autoclose Autoclose) autoClose(now time.Time) bool {
	shouldClose := autoclose.shouldClose(now)
	state := autoclose.doorController.getDoorState()
	log.Infof("Autoclose shouldClose=%t state=%b", shouldClose, state)
	if state == open {
		autoclose.lastOpen = time.Now()
		if shouldClose {
			autoclose.doorController.toggleDoor()
		}
		return true
	}
	return false
}
