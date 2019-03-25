package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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
	log.WithField("one_time_pin", oneTimePin).Info("Using one time pin to toggle garage")
	usedDate, err := garageDoorResource.pinDao.getPinUsedDate(oneTimePin)
	if usedDate > 0 {
		log.Info("Pin has already been used")
		w.WriteHeader(401)
		fmt.Fprintf(w, "Pin has already been used")
		return
	}
	if err != nil {
		log.WithError(err).WithField("one_time_pin", oneTimePin).Error("Could not get pin used date")
		w.WriteHeader(500)
		fmt.Fprintf(w, "Failed to open garage")
		return
	}
	err = garageDoorResource.pinDao.use(oneTimePin)
	if err != nil {
		log.WithError(err).WithField("one_time_pin", oneTimePin).Error("Could not use pin")
		w.WriteHeader(401)
		fmt.Fprintf(w, "Pin has already been used")
		return
	}
	garageDoorResource.doorController.toggleDoor()
	w.WriteHeader(202)
	w.Header().Set("Content-Type", "text/html")
	page := PinPage{
		CloseTime: timeToClose.Seconds(),
		Pin:       oneTimePin,
	}
	tmpl, _ := template.ParseFiles("used.html")
	tmpl.Execute(w, page)
	go garageDoorResource.closeGarage(oneTimePin)
}

func (garageDoorResource GarageDoorResource) closeGarage(pin string) {
	time.Sleep(timeToClose)
	log.WithField("one_time_pin", pin).Info("Closing garage")
	garageDoorResource.doorController.toggleDoor()
}

func (garageDoorResource GarageDoorResource) toggleGarage(w http.ResponseWriter, r *http.Request) {
	accessToken := context.Get(r, "access_token")
	email := getEmail(fmt.Sprintf("%s", accessToken))
	log.WithField("email", email).Info("opening or closing garage")
	vars := r.URL.Query()
	autoclose := vars.Get("autoclose")
	go func() {
		garageDoorResource.doorController.toggleDoor()
		if autocloseBool, _ := strconv.ParseBool(autoclose); autocloseBool {
			log.WithField("email", email).Info("autoclosing garage in 60s")
			time.Sleep(60 * time.Second)
			garageDoorResource.doorController.toggleDoor()
		}
	}()
	w.WriteHeader(202)
}

func (garageDoorResource GarageDoorResource) getState(w http.ResponseWriter, r *http.Request) {
	log.Debug("Getting garage state")
	state := garageDoorResource.doorController.getDoorState()
	w.Header().Set("Content-Type", "application/json")
	stateResponse := map[string]interface{}{
		"State":       state,
		"Description": state.description(),
	}
	json.NewEncoder(w).Encode(stateResponse)
}
