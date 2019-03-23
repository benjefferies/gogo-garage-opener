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
		log.WithField("state", autoCloseDoorController.state).Info("Opening door before")
		autoCloseDoorController.state = open
		log.WithField("state", autoCloseDoorController.state).Info("Closing door after")
	} else {
		log.WithField("state", autoCloseDoorController.state).Info("Closing door before")
		autoCloseDoorController.state = closed
		log.WithField("state", autoCloseDoorController.state).Info("Opening door after")
	}
}

// getDoorState is noop for testing always returning closed
func (autoCloseDoorController AutoCloseDoorController) getDoorState() DoorState {
	log.WithField("state", autoCloseDoorController.state).Info("Getting state")
	return autoCloseDoorController.state
}

// close is noop for testing
func (autoCloseDoorController AutoCloseDoorController) close() {
	log.Info("Noop close")
}

func TestNewInstanceSetsClosingTimeAt10pm(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	shouldClose := time.Date(now.Year(), now.Month(), now.Day(), 22, 0, 0, 0, time.Local)

	autoclose := NewAutoclose(controller)

	assert.Equal(t, shouldClose, autoclose.shouldCloseTime, "Should close at 10pm")
}

func TestNewInstanceSetsCanStayOpenTimeAt8am(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	canStayOpen := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, time.Local)

	autoclose := NewAutoclose(controller)

	assert.Equal(t, canStayOpen, autoclose.canStayOpenTime, "Should stay open at 8am")
}

func TestShouldAutoCloseAfterShouldCloseTime(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	shouldClose := now.Add(-time.Minute)
	canStayOpen := now.Add(-time.Hour)
	autoclose := Autoclose{openDuration: time.Minute * 3, doorController: controller, shouldCloseTime: shouldClose, canStayOpenTime: canStayOpen}

	autoclose.autoClose()

	assert.True(t, controller.getDoorState().isClosed(), "Should be closed")
}

func TestShouldAutoCloseWhenShouldCloseTimeAfterCanStayOpenTime(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	shouldClose := now.Add(-time.Minute)
	canStayOpen := now.Add(-time.Minute * 2)
	autoclose := Autoclose{openDuration: time.Minute * 3, doorController: controller, shouldCloseTime: shouldClose, canStayOpenTime: canStayOpen}

	autoclose.autoClose()

	assert.True(t, controller.getDoorState().isClosed(), "Should be closed")
}

func TestShouldNotAutoCloseBeforeShouldCloseTime(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	shouldClose := now.Add(time.Minute)
	canStayOpen := now.Add(-time.Hour)
	autoclose := Autoclose{openDuration: time.Minute * 3, doorController: controller, shouldCloseTime: shouldClose, canStayOpenTime: canStayOpen}

	autoclose.autoClose()

	assert.True(t, controller.getDoorState().isOpen(), "Should not be closed")
}

func TestShouldAutoCloseBeforeCanStayOpenTime(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	shouldClose := now.Add(-time.Minute)
	canStayOpen := now.Add(-time.Hour)
	autoclose := Autoclose{openDuration: time.Minute * 3, doorController: controller, shouldCloseTime: shouldClose, canStayOpenTime: canStayOpen}

	autoclose.autoClose()

	assert.True(t, controller.getDoorState().isClosed(), "Should be closed")
}

func TestShouldNotOpenAfterShouldCloseTime(t *testing.T) {
	controller := &AutoCloseDoorController{closed}
	now := time.Now()
	shouldClose := now.Add(-time.Minute)
	canStayOpen := now.Add(-time.Hour)
	autoclose := Autoclose{openDuration: time.Minute * 3, doorController: controller, shouldCloseTime: shouldClose, canStayOpenTime: canStayOpen}

	autoclose.autoClose()

	assert.True(t, controller.getDoorState().isClosed(), "Should be closed")
}

func TestShouldNotCloseShouldCloseTimeWhenLeftOpenFor1Minutes(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	shouldClose := now.Add(-time.Minute)
	canStayOpen := now.Add(-time.Hour)
	autoclose := Autoclose{openDuration: time.Minute, doorController: controller, shouldCloseTime: shouldClose, canStayOpenTime: canStayOpen}

	autoclose.autoClose()

	assert.True(t, controller.getDoorState().isOpen(), "Should not be closed")
}

func TestShouldNotCloseShouldCloseTimeWhenLeftOpenFor3MinutesAndNotShouldCloseYet(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	shouldClose := now.Add(time.Hour)
	canStayOpen := now.Add(-time.Hour)
	autoclose := Autoclose{openDuration: time.Minute * 3, doorController: controller, shouldCloseTime: shouldClose, canStayOpenTime: canStayOpen}

	autoclose.autoClose()

	assert.True(t, controller.getDoorState().isOpen(), "Should not be closed")
}

func TestShouldCloseAfterShouldCloseWhenLeftOpenFor3Minutes(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	shouldClose := now.Add(-time.Minute)
	canStayOpen := now.Add(-time.Hour)
	autoclose := Autoclose{openDuration: time.Minute * 3, doorController: controller, shouldCloseTime: shouldClose, canStayOpenTime: canStayOpen}

	autoclose.autoClose()

	assert.True(t, controller.getDoorState().isClosed(), "Should be closed")
}

func TestShouldReturnTrueWhenClose(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	shouldClose := now.Add(-time.Minute)
	canStayOpen := now.Add(-time.Hour)
	autoclose := Autoclose{openDuration: time.Minute * 3, doorController: controller, shouldCloseTime: shouldClose, canStayOpenTime: canStayOpen}

	closing := autoclose.autoClose()

	assert.Equal(t, true, closing, "Should be closed")
}

func TestShouldReturnFalseWhenNotClose(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	shouldClose := now.Add(-time.Minute)
	canStayOpen := now.Add(-time.Hour)
	autoclose := Autoclose{openDuration: time.Minute * 1, doorController: controller, shouldCloseTime: shouldClose, canStayOpenTime: canStayOpen}

	closing := autoclose.autoClose()

	assert.Equal(t, false, closing, "Should be closed")
}

func TestIncreaseOpenDurationWhenNotClosed(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	shouldClose := now.Add(-time.Minute)
	canStayOpen := now.Add(-time.Hour)
	autoclose := Autoclose{openDuration: time.Minute * 1, doorController: controller, shouldCloseTime: shouldClose, canStayOpenTime: canStayOpen}

	autoclose.autoClose()

	assert.Equal(t, time.Minute*2, autoclose.openDuration, "Should increase open duration from 1 to 2 minutes")
}

func TestResetOpenDurationWhenClosed(t *testing.T) {
	controller := &AutoCloseDoorController{closed}
	now := time.Now()
	shouldClose := now.Add(-time.Minute)
	canStayOpen := now.Add(-time.Hour)
	autoclose := Autoclose{openDuration: time.Minute * 3, doorController: controller, shouldCloseTime: shouldClose, canStayOpenTime: canStayOpen}

	autoclose.autoClose()

	assert.Equal(t, time.Minute*0, autoclose.openDuration, "Should reset open duration from 1 to 0 minutes")
}

func TestResetOpenDurationWhenClosing(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	shouldClose := now.Add(-time.Minute)
	canStayOpen := now.Add(-time.Hour)
	autoclose := Autoclose{openDuration: time.Minute * 3, doorController: controller, shouldCloseTime: shouldClose, canStayOpenTime: canStayOpen}

	autoclose.autoClose()

	assert.Equal(t, time.Minute*0, autoclose.openDuration, "Should reset open duration from 1 to 0 minutes")
}
