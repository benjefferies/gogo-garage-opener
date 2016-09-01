package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/stianeikeland/go-rpio"
	"time"
)

type DoorController struct {
	relayPin         int
	contactSwitchPin int
}

func (c DoorController) toggleDoor() {

	log.Debugf("Using relay pin %d to toggle relay", c.relayPin)
	pin := rpio.Pin(c.relayPin)
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		log.Error(err)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	// Set pin to output mode
	pin.Output()

	// Toggle pin on/off
	pin.Toggle()
	log.Infof("Toggle relay switch on")
	time.Sleep(time.Millisecond * 500)
	pin.Toggle()
	log.Infof("Toggle relay switch off")
}

func (c DoorController) getDoorState() rpio.State {

	log.Debugf("Using pin %d to read contact switch pin", c.contactSwitchPin)
	pin := rpio.Pin(c.contactSwitchPin)
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		log.Error(err)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	// Set pin to input mode
	pin.Input()

	// Read state
	state := pin.Read()
	log.Infof("Sensor reading state: %s", state)
	return state
}
