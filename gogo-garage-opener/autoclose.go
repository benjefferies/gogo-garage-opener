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
	canStayOpen := now.After(autoclose.canStayOpenTime) && now.Before(autoclose.shouldCloseTime)
	log.WithField("openTooLong", openTooLong).
		WithField("canStayOpen", canStayOpen).
		WithField("openDuration", autoclose.openDuration).
		WithField("shouldCloseTime", autoclose.shouldCloseTime.Format("3:04:05 PM")).
		WithField("canStayOpenTime", autoclose.canStayOpenTime.Format("3:04:05 PM")).
		WithField("now", now.Format("3:04:05 PM")).
		Debug("Evaluating if garage should change state")
	if !canStayOpen && openTooLong {
		log.Debug("Should close garage door")
		return true
	}
	log.Debug("Should not close garage door")
	return false
}

func (autoclose *Autoclose) autoClose() bool {
	state := autoclose.doorController.getDoorState()
	if state.isClosed() {
		autoclose.openDuration = time.Minute * 0
		return false
	}

	shouldClose := autoclose.shouldClose()
	if shouldClose {
		log.Debug("Autoclosing garage")
		autoclose.doorController.toggleDoor()
		autoclose.openDuration = time.Minute * 0
		return true
	}
	autoclose.openDuration = autoclose.openDuration + time.Minute // Time increase needs to be in sync with sleep time
	log.WithField("openDuration", autoclose.openDuration).Debug("Garage left open increasing openDuration")
	return false
}
