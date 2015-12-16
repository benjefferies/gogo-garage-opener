package main

import (
	"github.com/stianeikeland/go-rpio"
	log "github.com/Sirupsen/logrus"
	"time"
)

type DoorController struct {
	relayPin int;
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
	time.Sleep(time.Millisecond * 500)
	pin.Toggle()
	log.Debugf("Toggled pin %d on/off", c.relayPin)
}