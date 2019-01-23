package main

import (
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// NoopDoorController for testing
type AutoCloseDoorController struct {
	state DoorState
}

// toggleDoor is noop for testing
func (autoCloseDoorController *AutoCloseDoorController) toggleDoor() {
	if autoCloseDoorController.state == closed {
		log.Infof("Opening door before state=%b", autoCloseDoorController.state)
		autoCloseDoorController.state = open
		log.Infof("Opening door after state=%b", autoCloseDoorController.state)
	} else {
		log.Infof("Closing door before state=%b", autoCloseDoorController.state)
		autoCloseDoorController.state = closed
		log.Infof("Opening door after state=%b", autoCloseDoorController.state)
	}
}

// getDoorState is noop for testing always returning closed
func (autoCloseDoorController AutoCloseDoorController) getDoorState() DoorState {
	log.Infof("Getting state=%b", autoCloseDoorController.state)
	return autoCloseDoorController.state
}

// close is noop for testing
func (autoCloseDoorController AutoCloseDoorController) close() {
	log.Info("Noop close")
}

func TestShouldAutoCloseAfter10(t *testing.T) {
	var now = time.Now()
	var shouldClose = time.Date(now.Year(), now.Month(), now.Day(), 22, 0, 1, 0, time.UTC)
	var controller DoorController = &AutoCloseDoorController{open}
	var autoclose = Autoclose{now: shouldClose, doorController: controller}

	autoclose.autoClose()

	assert.Equal(t, closed, controller.getDoorState(), "Should be closed")
}

func TestShouldNotAutoCloseBefore10(t *testing.T) {
	var now = time.Now()
	var shouldClose = time.Date(now.Year(), now.Month(), now.Day(), 21, 0, 1, 0, time.UTC)
	var controller DoorController = &AutoCloseDoorController{open}
	var autoclose = Autoclose{now: shouldClose, doorController: controller}

	autoclose.autoClose()

	assert.Equal(t, open, controller.getDoorState(), "Should not be closed")
}

func TestShouldAutoCloseBefore8(t *testing.T) {
	var now = time.Now()
	var shouldClose = time.Date(now.Year(), now.Month(), now.Day(), 7, 0, 1, 0, time.UTC)
	var controller DoorController = &AutoCloseDoorController{open}
	var autoclose = Autoclose{now: shouldClose, doorController: controller}

	autoclose.autoClose()

	assert.Equal(t, closed, controller.getDoorState(), "Should be closed")
}

func TestShouldNotOpenAfter10(t *testing.T) {
	var now = time.Now()
	var shouldClose = time.Date(now.Year(), now.Month(), now.Day(), 22, 0, 1, 0, time.UTC)
	var controller DoorController = &AutoCloseDoorController{open}
	var autoclose = Autoclose{now: shouldClose, doorController: controller}

	autoclose.autoClose()

	assert.Equal(t, closed, controller.getDoorState(), "Should be closed")
}
