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

func TestShouldAutoCloseAfterShouldCloseTime_LateEveningExample(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	canStayOpen := now.Add(-time.Hour)
	shouldClose := now.Add(-time.Minute)
	openDuration := int64(120)
	garageDoorDao := noOpGarageDoorDao{openDuration: openDuration, shouldClose: shouldClose, canStayOpen: canStayOpen}
	autoclose := NewAutoclose(controller, garageDoorDao)
	autoclose.openDuration = 120 * time.Second

	autoclose.autoClose()

	assert.True(t, controller.getDoorState().isClosed(), "Should be closed")
}

func TestShouldAutoCloseWhenShouldCloseTimeAfterCanStayOpenTime_LateEveningExample(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	canStayOpen := now.Add(-time.Hour)
	shouldClose := now.Add(-time.Minute)
	openDuration := int64(120)
	garageDoorDao := noOpGarageDoorDao{openDuration: openDuration, shouldClose: shouldClose, canStayOpen: canStayOpen}
	autoclose := NewAutoclose(controller, garageDoorDao)
	autoclose.openDuration = 120 * time.Second

	autoclose.autoClose()

	assert.True(t, controller.getDoorState().isClosed(), "Should be closed")
}

func TestShouldNotAutoCloseBeforeShouldCloseTime_MiddayExample(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	canStayOpen := now.Add(-time.Hour)
	shouldClose := now.Add(time.Minute)
	openDuration := int64(120)
	garageDoorDao := noOpGarageDoorDao{openDuration: openDuration, shouldClose: shouldClose, canStayOpen: canStayOpen}
	autoclose := NewAutoclose(controller, garageDoorDao)
	autoclose.openDuration = 120 * time.Second

	autoclose.autoClose()

	assert.True(t, controller.getDoorState().isOpen(), "Should not be closed")
}

func TestShouldAutoCloseBeforeCanStayOpenTime_EarlyMorningExample(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	canStayOpen := now.Add(time.Minute)
	shouldClose := now.Add(time.Hour)
	openDuration := int64(120)
	garageDoorDao := noOpGarageDoorDao{openDuration: openDuration, shouldClose: shouldClose, canStayOpen: canStayOpen}
	autoclose := NewAutoclose(controller, garageDoorDao)
	autoclose.openDuration = 120 * time.Second

	autoclose.autoClose()

	assert.True(t, controller.getDoorState().isClosed(), "Should be closed")
}

func TestShouldNotOpenAfterShouldCloseTime_NoToggleLateEveningExample(t *testing.T) {
	controller := &AutoCloseDoorController{closed}
	now := time.Now()
	canStayOpen := now.Add(-time.Hour)
	shouldClose := now.Add(-time.Minute)
	openDuration := int64(120)
	garageDoorDao := noOpGarageDoorDao{openDuration: openDuration, shouldClose: shouldClose, canStayOpen: canStayOpen}
	autoclose := NewAutoclose(controller, garageDoorDao)
	autoclose.openDuration = 120 * time.Second

	autoclose.autoClose()

	assert.True(t, controller.getDoorState().isClosed(), "Should be closed")
}

func TestShouldNotCloseWhenLeftOpenFor1MinutesAndNotShouldCloseYet_MiddayExample(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	canStayOpen := now.Add(-time.Minute)
	shouldClose := now.Add(time.Hour)
	openDuration := int64(120)
	garageDoorDao := noOpGarageDoorDao{openDuration: openDuration, shouldClose: shouldClose, canStayOpen: canStayOpen}
	autoclose := NewAutoclose(controller, garageDoorDao)
	autoclose.openDuration = 60 * time.Second

	autoclose.autoClose()

	assert.True(t, controller.getDoorState().isOpen(), "Should not be closed")
}

func TestShouldNotCloseWhenLeftOpenFor3MinutesAndNotShouldCloseYet_MiddayExample(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	canStayOpen := now.Add(-time.Minute)
	shouldClose := now.Add(time.Hour)
	openDuration := int64(120)
	garageDoorDao := noOpGarageDoorDao{openDuration: openDuration, shouldClose: shouldClose, canStayOpen: canStayOpen}
	autoclose := NewAutoclose(controller, garageDoorDao)
	autoclose.openDuration = 120 * time.Second

	autoclose.autoClose()

	assert.True(t, controller.getDoorState().isOpen(), "Should not be closed")
}

func TestShouldNotCloseAfterShouldCloseWhenLeftOpenFor1Minutes_LateEveningExample(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	canStayOpen := now.Add(-time.Hour)
	shouldClose := now.Add(-time.Minute)
	openDuration := int64(120)
	garageDoorDao := noOpGarageDoorDao{openDuration: openDuration, shouldClose: shouldClose, canStayOpen: canStayOpen}
	autoclose := NewAutoclose(controller, garageDoorDao)
	autoclose.openDuration = 60 * time.Second

	autoclose.autoClose()

	assert.False(t, controller.getDoorState().isClosed(), "Should be closed")
}

func TestShouldCloseAfterShouldCloseWhenLeftOpenFor3Minutes_LateEveningExample(t *testing.T) {
	controller := &AutoCloseDoorController{open}
	now := time.Now()
	canStayOpen := now.Add(-time.Hour)
	shouldClose := now.Add(-time.Minute)
	openDuration := int64(120)
	garageDoorDao := noOpGarageDoorDao{openDuration: openDuration, shouldClose: shouldClose, canStayOpen: canStayOpen}
	autoclose := NewAutoclose(controller, garageDoorDao)
	autoclose.openDuration = 180 * time.Second

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

	assert.False(t, closing, "Should be closed")
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
	shouldClose := now.Add(-time.Hour * 23)
	canStayOpen := now.Add(-time.Hour * 24).Add(-time.Minute)
	openDuration := int64(120)
	garageDoorDao := noOpGarageDoorDao{openDuration: openDuration, shouldClose: shouldClose, canStayOpen: canStayOpen}
	autoclose := NewAutoclose(controller, garageDoorDao)
	autoclose.openDuration = 120 * time.Second

	should := autoclose.shouldClose()

	assert.Equal(t, false, should, "Should not close after resetting times")
}
