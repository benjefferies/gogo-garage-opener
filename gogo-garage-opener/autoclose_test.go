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

type noOpGarageDoorDao struct {
	shouldClose  time.Time
	canStayOpen  time.Time
	openDuration int64
	enabled      bool
}

func (noOpGarageDoorDao noOpGarageDoorDao) updateConfiguration(updateConfig []GarageConfiguration) error {
	return nil
}
func (noOpGarageDoorDao noOpGarageDoorDao) getConfiguration() ([]GarageConfiguration, error) {
	var config []GarageConfiguration
	for _, day := range [7]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"} {
		var theDay = day
		dayConfig := GarageConfiguration{Day: &theDay, OpenDuration: &noOpGarageDoorDao.openDuration, ShouldCloseTime: &noOpGarageDoorDao.shouldClose,
			CanStayOpenTime: &noOpGarageDoorDao.canStayOpen, Enabled: &noOpGarageDoorDao.enabled}
		config = append(config, dayConfig)
	}
	return config, nil
}

func TestShouldAutoCloseAfterShouldCloseTime(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	shouldClose := now.Add(-time.Minute)
	canStayOpen := now.Add(-time.Hour)
	openDuration := int64(120)
	garageDoorDao := noOpGarageDoorDao{openDuration: openDuration, shouldClose: shouldClose, canStayOpen: canStayOpen}
	autoclose := NewAutoclose(controller, garageDoorDao)
	autoclose.openDuration = 120 * time.Second

	autoclose.autoClose()

	assert.True(t, controller.getDoorState().isClosed(), "Should be closed")
}

func TestShouldAutoCloseWhenShouldCloseTimeAfterCanStayOpenTime(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	shouldClose := now.Add(-time.Minute)
	canStayOpen := now.Add(-time.Minute * 2)
	openDuration := int64(120)
	garageDoorDao := noOpGarageDoorDao{openDuration: openDuration, shouldClose: shouldClose, canStayOpen: canStayOpen}
	autoclose := NewAutoclose(controller, garageDoorDao)
	autoclose.openDuration = 120 * time.Second

	autoclose.autoClose()

	assert.True(t, controller.getDoorState().isClosed(), "Should be closed")
}

func TestShouldNotAutoCloseBeforeShouldCloseTime(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	shouldClose := now.Add(time.Minute)
	canStayOpen := now.Add(-time.Hour)
	openDuration := int64(120)
	garageDoorDao := noOpGarageDoorDao{openDuration: openDuration, shouldClose: shouldClose, canStayOpen: canStayOpen}
	autoclose := NewAutoclose(controller, garageDoorDao)
	autoclose.openDuration = 120 * time.Second

	autoclose.autoClose()

	assert.True(t, controller.getDoorState().isOpen(), "Should not be closed")
}

func TestShouldAutoCloseBeforeCanStayOpenTime(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	shouldClose := now.Add(-time.Minute)
	canStayOpen := now.Add(-time.Hour)
	openDuration := int64(120)
	garageDoorDao := noOpGarageDoorDao{openDuration: openDuration, shouldClose: shouldClose, canStayOpen: canStayOpen}
	autoclose := NewAutoclose(controller, garageDoorDao)
	autoclose.openDuration = 120 * time.Second

	autoclose.autoClose()

	assert.True(t, controller.getDoorState().isClosed(), "Should be closed")
}

func TestShouldNotOpenAfterShouldCloseTime(t *testing.T) {
	controller := &AutoCloseDoorController{closed}
	now := time.Now()
	shouldClose := now.Add(-time.Minute)
	canStayOpen := now.Add(-time.Hour)
	openDuration := int64(120)
	garageDoorDao := noOpGarageDoorDao{openDuration: openDuration, shouldClose: shouldClose, canStayOpen: canStayOpen}
	autoclose := NewAutoclose(controller, garageDoorDao)
	autoclose.openDuration = 120 * time.Second

	autoclose.autoClose()

	assert.True(t, controller.getDoorState().isClosed(), "Should be closed")
}

func TestShouldNotCloseShouldCloseTimeWhenLeftOpenFor1Minutes(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	shouldClose := now.Add(-time.Minute)
	canStayOpen := now.Add(-time.Hour)
	openDuration := int64(120)
	garageDoorDao := noOpGarageDoorDao{openDuration: openDuration, shouldClose: shouldClose, canStayOpen: canStayOpen}
	autoclose := NewAutoclose(controller, garageDoorDao)
	autoclose.openDuration = 60 * time.Second

	autoclose.autoClose()

	assert.True(t, controller.getDoorState().isOpen(), "Should not be closed")
}

func TestShouldNotCloseShouldCloseTimeWhenLeftOpenFor3MinutesAndNotShouldCloseYet(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	shouldClose := now.Add(time.Hour)
	canStayOpen := now.Add(-time.Hour)
	openDuration := int64(120)
	garageDoorDao := noOpGarageDoorDao{openDuration: openDuration, shouldClose: shouldClose, canStayOpen: canStayOpen}
	autoclose := NewAutoclose(controller, garageDoorDao)
	autoclose.openDuration = 120 * time.Second

	autoclose.autoClose()

	assert.True(t, controller.getDoorState().isOpen(), "Should not be closed")
}

func TestShouldCloseAfterShouldCloseWhenLeftOpenFor3Minutes(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	shouldClose := now.Add(-time.Minute)
	canStayOpen := now.Add(-time.Hour)
	openDuration := int64(120)
	garageDoorDao := noOpGarageDoorDao{openDuration: openDuration, shouldClose: shouldClose, canStayOpen: canStayOpen}
	autoclose := NewAutoclose(controller, garageDoorDao)
	autoclose.openDuration = 120 * time.Second

	autoclose.autoClose()

	assert.True(t, controller.getDoorState().isClosed(), "Should be closed")
}

func TestShouldReturnTrueWhenClose(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	shouldClose := now.Add(-time.Minute)
	canStayOpen := now.Add(-time.Hour)
	openDuration := int64(120)
	garageDoorDao := noOpGarageDoorDao{openDuration: openDuration, shouldClose: shouldClose, canStayOpen: canStayOpen}
	autoclose := NewAutoclose(controller, garageDoorDao)
	autoclose.openDuration = 120 * time.Second

	closing := autoclose.autoClose()

	assert.Equal(t, true, closing, "Should be closed")
}

func TestShouldReturnFalseWhenNotClose(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	shouldClose := now.Add(-time.Minute)
	canStayOpen := now.Add(-time.Hour)
	openDuration := int64(120)
	garageDoorDao := noOpGarageDoorDao{openDuration: openDuration, shouldClose: shouldClose, canStayOpen: canStayOpen}
	autoclose := NewAutoclose(controller, garageDoorDao)
	autoclose.openDuration = 60 * time.Second

	closing := autoclose.autoClose()

	assert.Equal(t, false, closing, "Should be closed")
}

func TestIncreaseOpenDurationWhenNotClosed(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	shouldClose := now.Add(-time.Minute)
	canStayOpen := now.Add(-time.Hour)
	openDuration := int64(120)
	garageDoorDao := noOpGarageDoorDao{openDuration: openDuration, shouldClose: shouldClose, canStayOpen: canStayOpen}
	autoclose := NewAutoclose(controller, garageDoorDao)
	autoclose.openDuration = 60 * time.Second

	autoclose.autoClose()

	assert.Equal(t, time.Minute*2, autoclose.openDuration, "Should increase open duration from 1 to 2 minutes")
}

func TestResetOpenDurationWhenClosed(t *testing.T) {
	controller := &AutoCloseDoorController{closed}
	now := time.Now()
	shouldClose := now.Add(-time.Minute)
	canStayOpen := now.Add(-time.Hour)
	openDuration := int64(120)
	garageDoorDao := noOpGarageDoorDao{openDuration: openDuration, shouldClose: shouldClose, canStayOpen: canStayOpen}
	autoclose := NewAutoclose(controller, garageDoorDao)
	autoclose.openDuration = 120 * time.Second

	autoclose.autoClose()

	assert.Equal(t, time.Minute*0, autoclose.openDuration, "Should reset open duration from 1 to 0 minutes")
}

func TestResetOpenDurationWhenClosing(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	shouldClose := now.Add(-time.Minute)
	canStayOpen := now.Add(-time.Hour)
	openDuration := int64(120)
	garageDoorDao := noOpGarageDoorDao{openDuration: openDuration, shouldClose: shouldClose, canStayOpen: canStayOpen}
	autoclose := NewAutoclose(controller, garageDoorDao)
	autoclose.openDuration = 120 * time.Second

	autoclose.autoClose()

	assert.Equal(t, time.Minute*0, autoclose.openDuration, "Should reset open duration from 1 to 0 minutes")
}

func TestShouldResetTimesToRolloverDay(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	shouldClose := now.Add(-time.Hour * 24).Add(time.Minute)
	canStayOpen := now.Add(-time.Hour * 25)
	openDuration := int64(120)
	garageDoorDao := noOpGarageDoorDao{openDuration: openDuration, shouldClose: shouldClose, canStayOpen: canStayOpen}
	autoclose := NewAutoclose(controller, garageDoorDao)
	autoclose.openDuration = 120 * time.Second

	should := autoclose.shouldClose()

	assert.Equal(t, true, should, "Should not close after resetting times")
}
