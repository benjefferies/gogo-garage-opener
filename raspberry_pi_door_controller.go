package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/stianeikeland/go-rpio"
)

// RaspberryPiDoorController the controller for accessing the garage door
type RaspberryPiDoorController struct {
	relayPin         rpio.Pin
	contactSwitchPin rpio.Pin
}

// NewRaspberryPiDoorController Constructor for RaspberryPiDoorController
func NewRaspberryPiDoorController(relayPinID int, contactSwitchPinID int) RaspberryPiDoorController {
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		log.Error(err)
	}
	relayPin := rpio.Pin(relayPinID)
	contactSwitchPin := rpio.Pin(contactSwitchPinID)
	relayPin.Output()
	contactSwitchPin.Input()
	return RaspberryPiDoorController{relayPin: relayPin, contactSwitchPin: contactSwitchPin}
}

// Open or close garage door
func (raspberryPiDoorController RaspberryPiDoorController) toggleDoor() {
	// Toggle pin on/off
	raspberryPiDoorController.relayPin.Toggle()
	log.Infof("Toggle relay switch on")
	time.Sleep(time.Millisecond * 500)
	raspberryPiDoorController.relayPin.Toggle()
	log.Infof("Toggle relay switch off")
}

// Get the state of the garage door
func (raspberryPiDoorController RaspberryPiDoorController) getDoorState() DoorState {

	log.Debugf("Using pin %d to read contact switch pin", raspberryPiDoorController.contactSwitchPin)
	pin := rpio.Pin(raspberryPiDoorController.contactSwitchPin)
	// Open and map  memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		log.Error(err)
	}

	// Read state
	state := pin.Read()
	log.Infof("Sensor reading state: %v", state)
	doorState := int8(state)
	if doorState == int8(closed) {
		return closed
	}
	return open
}

// Close rpio
func (raspberryPiDoorController RaspberryPiDoorController) close() {
	err := rpio.Close()
	if err != nil {
		log.WithError(err).Error("Could not close pins")
	}
}
