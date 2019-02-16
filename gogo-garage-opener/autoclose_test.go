package main

import (
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
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
	var controller DoorController = &AutoCloseDoorController{open}
	var lastOpen = time.Date(now.Year(), now.Month(), now.Day(), 22, 0, 0, 0, time.UTC)
	var autoclose = Autoclose{lastOpen: lastOpen, doorController: controller}
	var shouldClose = time.Date(now.Year(), now.Month(), now.Day(), 22, 2, 59, 0, time.UTC)

	autoclose.autoClose(shouldClose)

	assert.Equal(t, closed, controller.getDoorState(), "Should be closed")
}

func TestShouldNotAutoCloseBefore10(t *testing.T) {
	var now = time.Now()
	var controller DoorController = &AutoCloseDoorController{open}
	var lastOpen = time.Date(now.Year(), now.Month(), now.Day(), 21, 0, 0, 0, time.UTC)
	var autoclose = Autoclose{lastOpen: lastOpen, doorController: controller}
	var shouldNotClose = time.Date(now.Year(), now.Month(), now.Day(), 21, 2, 59, 0, time.UTC)

	autoclose.autoClose(shouldNotClose)

	assert.Equal(t, open, controller.getDoorState(), "Should not be closed")
}

func TestShouldAutoCloseBefore8(t *testing.T) {
	var now = time.Now()
	var controller DoorController = &AutoCloseDoorController{open}
	var lastOpen = time.Date(now.Year(), now.Month(), now.Day(), 7, 0, 0, 0, time.UTC)
	var autoclose = Autoclose{lastOpen: lastOpen, doorController: controller}
	var shouldClose = time.Date(now.Year(), now.Month(), now.Day(), 7, 2, 59, 0, time.UTC)

	autoclose.autoClose(shouldClose)

	assert.Equal(t, closed, controller.getDoorState(), "Should be closed")
}

func TestShouldNotOpenAfter10(t *testing.T) {
	var now = time.Now()
	var controller DoorController = &AutoCloseDoorController{closed}
	var lastOpen = time.Date(now.Year(), now.Month(), now.Day(), 22, 0, 0, 0, time.UTC)
	var autoclose = Autoclose{lastOpen: lastOpen, doorController: controller}
	var shouldNotOpen = time.Date(now.Year(), now.Month(), now.Day(), 22, 2, 59, 0, time.UTC)

	autoclose.autoClose(shouldNotOpen)

	assert.Equal(t, closed, controller.getDoorState(), "Should be closed")
}

func TestShouldNotCloseAfter10AndLeftOpenFor1Minutes(t *testing.T) {
	var now = time.Now()
	var controller DoorController = &AutoCloseDoorController{open}
	var lastOpen = time.Date(now.Year(), now.Month(), now.Day(), 22, 0, 0, 0, time.UTC)
	var autoclose = Autoclose{lastOpen: lastOpen, doorController: controller}
	var shouldNotClose = time.Date(now.Year(), now.Month(), now.Day(), 22, 0, 59, 0, time.UTC)

	autoclose.autoClose(shouldNotClose)

	assert.Equal(t, open, controller.getDoorState(), "Should not be closed")
}

func TestShouldCloseAfter10AndLeftOpenFor3Minutes(t *testing.T) {
	var now = time.Now()
	var controller DoorController = &AutoCloseDoorController{open}
	var lastOpen = time.Date(now.Year(), now.Month(), now.Day(), 22, 0, 0, 0, time.UTC)
	var autoclose = Autoclose{lastOpen: lastOpen, doorController: controller}
	var shouldNotClose = time.Date(now.Year(), now.Month(), now.Day(), 22, 3, 0, 0, time.UTC)

	autoclose.autoClose(shouldNotClose)

	assert.Equal(t, closed, controller.getDoorState(), "Should be closed")
}
