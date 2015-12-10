package main

import (
	"github.com/stianeikeland/go-rpio"
	log "github.com/Sirupsen/logrus"
	"time"
	"github.com/spf13/viper"
)


func toggleDoor() {

	relayPin := viper.GetInt("relayPin")
	log.Debugf("Using relay pin %s to open garage", relayPin)
	pin := rpio.Pin(relayPin)
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
	time.Sleep(time.Second * 2)
	pin.Toggle()
}