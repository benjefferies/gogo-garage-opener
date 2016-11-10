package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/stianeikeland/go-rpio"
	"time"
)


type RaspberryPiDoorController struct {
	relayPin         rpio.Pin
	contactSwitchPin rpio.Pin
}

func NewRaspberryPiDoorController(relayPinId int, contactSwitchPinId int) RaspberryPiDoorController {
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		log.Error(err)
	}
	relayPin := rpio.Pin(relayPinId)
	contactSwitchPin := rpio.Pin(contactSwitchPinId)
	relayPin.Output()
	contactSwitchPin.Input()
	return RaspberryPiDoorController{relayPin: relayPin, contactSwitchPin: contactSwitchPin}
}

func (this RaspberryPiDoorController) toggleDoor() {
	// Toggle pin on/off
	this.relayPin.Toggle()
	log.Infof("Toggle relay switch on")
	time.Sleep(time.Millisecond * 500)
	this.relayPin.Toggle()
	log.Infof("Toggle relay switch off")
}

func (this RaspberryPiDoorController) getDoorState() DoorState {

	log.Debugf("Using pin %d to read contact switch pin", this.contactSwitchPin)
	pin := rpio.Pin(this.contactSwitchPin)
	// Open and map  memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		log.Error(err)
	}

	// Read state
	state := pin.Read()
	log.Infof("Sensor reading state: %v", state)
	doorState := int8(state)
	if (doorState == int8(closed)) {
		return closed
	} else {
		return open
	}
}

func (this RaspberryPiDoorController) close() {
	err := rpio.Close()
	if err != nil {
		log.WithError(err).Error("Could not close pins")
	}
}