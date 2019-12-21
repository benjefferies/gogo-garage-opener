package main

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Autoclose is to auto close the garage door
type Autoclose struct {
	openDuration   time.Duration
	config         GarageConfiguration
	doorController DoorController
	garageDoorDao  GarageDoorDao
}

// NewAutoclose AutoClose with default openDuration set to zero and closing time between 10pm-8am
func NewAutoclose(doorcontroller DoorController, garageDoorDao GarageDoorDao) Autoclose {
	return Autoclose{openDuration: time.Second * 0, doorController: doorcontroller, garageDoorDao: garageDoorDao}.resetShouldCloseAndStayOpenTimes()
}

func (autoclose Autoclose) resetShouldCloseAndStayOpenTimes() Autoclose {
	now := time.Now()
	config, err := autoclose.garageDoorDao.getConfiguration()
	if err != nil {
		shouldCloseTime := time.Date(now.Year(), now.Month(), now.Day(), 22, 0, 0, 0, time.Local)
		canStayOpenTime := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, time.Local)
		enabled := true
		canStayOpenDuration := int64(120)
		day := now.Weekday().String()
		autoclose.config = GarageConfiguration{Day: &day, ShouldCloseTime: &shouldCloseTime, CanStayOpenTime: &canStayOpenTime, Enabled: &enabled, OpenDuration: &canStayOpenDuration}
		log.WithField("config", autoclose.config).Error("Could not get garage configuration using default")
	} else {
		autoclose.config = getConfigByDay(now.Weekday().String(), config)
	}
	return autoclose
}

func getConfigByDay(day string, config []GarageConfiguration) GarageConfiguration {
	for _, dc := range config {
		if *dc.Day == day {
			return dc
		}
	}
	return GarageConfiguration{}
}

func (autoclose Autoclose) shouldClose() bool {
	autoclose.resetShouldCloseAndStayOpenTimes()
	now := time.Now()
	openTooLong := autoclose.openDuration >= time.Second*time.Duration(autoclose.openDuration)
	canStayOpen := now.After(*autoclose.config.CanStayOpenTime) && now.Before(*autoclose.config.ShouldCloseTime)
	log.WithField("openTooLong", openTooLong).
		WithField("canStayOpen", canStayOpen).
		WithField("openDuration", autoclose.openDuration).
		WithField("shouldCloseTime", autoclose.config.ShouldCloseTime.Format("3:04:05 PM")).
		WithField("canStayOpenTime", autoclose.config.CanStayOpenTime.Format("3:04:05 PM")).
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
		autoclose.openDuration = time.Second * 0
		return false
	}

	shouldClose := autoclose.shouldClose()
	if shouldClose {
		log.Debug("Autoclosing garage")
		autoclose.doorController.toggleDoor()
		autoclose.openDuration = time.Second * 0
		return true
	}
	autoclose.openDuration = autoclose.openDuration + time.Minute // Time increase needs to be in sync with sleep time
	log.WithField("openDuration", autoclose.openDuration).Debug("Garage left open increasing openDuration")
	return false
}
