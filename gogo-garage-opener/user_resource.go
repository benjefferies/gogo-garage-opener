package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// PinPage model for enter pin page
type PinPage struct {
	CloseTime float64
	Pin       string
}

// UserResource API for users
type UserResource struct {
	userDao UserDao
	pinDao  PinDao
}

func (userResource UserResource) register(router *mux.Router) {
	subRouter := router.PathPrefix("/user").Subrouter()
	subRouter.Path("/login").Methods("POST").Handler(jwtCheckHandleFunc(userResource.login))
	subRouter.Path("/one-time-pin").Methods("POST").Handler(jwtCheckHandleFunc(userResource.oneTimePin))
	subRouter.Path("/one-time-pin").Methods("GET").Handler(jwtCheckHandleFunc(userResource.getOneTimePins))
	subRouter.Path("/one-time-pin/{oneTimePin}").Methods("DELETE").Handler(jwtCheckHandleFunc(userResource.deleteOneTimePin))
	subRouter.Path("/one-time-pin/{oneTimePin}").Methods("GET").HandlerFunc(userResource.oneTimePinPage)
}

func (userResource UserResource) oneTimePin(w http.ResponseWriter, r *http.Request) {
	log.Debug("Creating new one time pin")
	accessToken := context.Get(r, "access_token")
	email := getEmail(fmt.Sprintf("%s", accessToken))
	log.WithField("email", email).Debug("creating new one time pin")
	pin, err := userResource.pinDao.newOneTimePin(email)
	if err != nil {
		log.WithError(err).Error("Could not create one time pin")
		w.WriteHeader(500)
		return
	}
	pinMap := map[string]interface{}{
		"pin": pin,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pinMap)
	w.WriteHeader(http.StatusOK)
}

func (userResource UserResource) oneTimePinPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	oneTimePin := vars["oneTimePin"]
	log.WithField("one_time_pin", oneTimePin).Debug("Using one time pin")
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	page := PinPage{
		CloseTime: timeToClose.Seconds(),
		Pin:       oneTimePin,
	}
	tmpl, _ := template.ParseFiles("index.html")
	tmpl.Execute(w, page)
}

func (userResource UserResource) deleteOneTimePin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	oneTimePin := vars["oneTimePin"]
	log.WithField("one_time_pin", oneTimePin).Debug("Deleting one time pin")
	err := userResource.pinDao.delete(oneTimePin)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(500)
	}
	w.WriteHeader(200)
}

func (userResource UserResource) getOneTimePins(w http.ResponseWriter, r *http.Request) {
	pins, err := userResource.pinDao.getPins()
	if err != nil {
		log.WithError(err).Error("Could not get one time pins")
		w.WriteHeader(500)
		return
	}
	log.WithField("pins", pins).Info("Got pins")
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pins)
}

func (userResource UserResource) login(w http.ResponseWriter, r *http.Request) {
	accessToken := context.Get(r, "access_token")
	email := getEmail(fmt.Sprintf("%s", accessToken))
	log.WithField("email", email).Debug("User logged in")
	userResource.userDao.createUser(email)
}
