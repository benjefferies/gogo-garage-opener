package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// timeToClose time to automatically close the door after opening with pin
const timeToClose = 1 * time.Minute

// GarageDoorResource API for interacting with garage door
type GarageDoorResource struct {
	userDao        UserDao
	pinDao         PinDao
	doorController DoorController
}

func (garageDoorResource GarageDoorResource) register(router *mux.Router) {
	subRouter := router.PathPrefix("/garage").Subrouter()

	subRouter.Path("/toggle").Methods("POST").Handler(jwtCheckHandleFunc(garageDoorResource.toggleGarage))
	subRouter.Path("/state").Methods("GET").Handler(jwtCheckHandleFunc(garageDoorResource.getState))
	subRouter.Path("/one-time-pin/{oneTimePin}").Methods("POST").HandlerFunc(garageDoorResource.useOneTimePin)
}

func (garageDoorResource GarageDoorResource) useOneTimePin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	oneTimePin := vars["oneTimePin"]
	log.Infof("Using one time pin: [%s] to toggle garage", oneTimePin)
	usedDate, err := garageDoorResource.pinDao.getPinUsedDate(oneTimePin)
	if usedDate > 0 {
		log.Infof("Pin has already been used")
		w.WriteHeader(401)
		fmt.Fprintf(w, "Pin has already been used")
	} else if err != nil {
		log.WithError(err).Errorf("Could not get pin used date for [%s]", oneTimePin)
		w.WriteHeader(500)
		fmt.Fprintf(w, "Failed to open garage")
	} else {
		err = garageDoorResource.pinDao.use(oneTimePin)
		if err != nil {
			log.WithError(err).Errorf("Could not use pin: [%s]", oneTimePin)
			w.WriteHeader(401)
			fmt.Fprintf(w, "Pin has already been used")
			return
		}
		garageDoorResource.doorController.toggleDoor()
		w.WriteHeader(202)
		fmt.Fprintf(w, "Opening garage, it will close in %v seconds", timeToClose.Seconds())
		go garageDoorResource.closeGarage(oneTimePin)
	}
}

func (garageDoorResource GarageDoorResource) closeGarage(pin string) {
	time.Sleep(timeToClose)
	log.Infof("Closing garage for pin: [%s]", pin)
	garageDoorResource.doorController.toggleDoor()
}

func (garageDoorResource GarageDoorResource) toggleGarage(w http.ResponseWriter, r *http.Request) {
	accessToken := context.Get(r, "access_token")
	email := getEmail(fmt.Sprintf("%s", accessToken))
	log.Infof("%s is opening or closing garage", email)
	garageDoorResource.doorController.toggleDoor()
	w.WriteHeader(202)
}

func (garageDoorResource GarageDoorResource) getState(w http.ResponseWriter, r *http.Request) {
	log.Debug("Getting garage state")
	state := garageDoorResource.doorController.getDoorState()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(*newState(state))
}
