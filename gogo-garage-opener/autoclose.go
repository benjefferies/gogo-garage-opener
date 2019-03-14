package main

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Autoclose is to auto close the garage door
type Autoclose struct {
	openDuration    time.Duration
	shouldCloseTime time.Time
	canStayOpenTime time.Time
	doorController  DoorController
}

// NewAutoclose AutoClose with default openDuration set to zero and closing time between 10pm-8am
func NewAutoclose(doorcontroller DoorController) Autoclose {
	now := time.Now()
	shouldClose := time.Date(now.Year(), now.Month(), now.Day(), 22, 0, 0, 0, time.Local)
	canStayOpen := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, time.Local)
	return Autoclose{openDuration: time.Minute * 0, doorController: doorcontroller, shouldCloseTime: shouldClose, canStayOpenTime: canStayOpen}
}

func (autoclose Autoclose) shouldClose() bool {
	now := time.Now()
	openTooLong := autoclose.openDuration > time.Minute*2
	shouldClose := now.After(autoclose.shouldCloseTime)
	canStayOpen := now.Before(autoclose.canStayOpenTime)
	if shouldClose && openTooLong && !canStayOpen {
		return true
	}
	return false
}

func (autoclose *Autoclose) autoClose() bool {
	state := autoclose.doorController.getDoorState()
	if state == closed {
		autoclose.openDuration = time.Minute * 0
		return false
	}

	shouldClose := autoclose.shouldClose()
	log.Infof("Autoclose shouldClose=%t state=%b", shouldClose, state)
	if shouldClose {
		autoclose.doorController.toggleDoor()
		autoclose.openDuration = time.Minute * 0
		return true
	}
	autoclose.openDuration = autoclose.openDuration + time.Minute // Time increase needs to be in sync with sleep time
	log.Infof("Increased openDuration=%d", autoclose.openDuration)
	return false
}
