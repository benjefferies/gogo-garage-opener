package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

// AssistantResource API for google assistant webhook
type AssistantResource struct {
	doorController DoorController
}

func (assistantResource AssistantResource) register(router *mux.Router) {
	router.Path("/webhook").Methods("POST").Handler(basicAuth(assistantResource.handleWebhook))
}

func basicAuth(h http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, _ := r.BasicAuth()

		if *webhookUsername != user || *webhookPassword != pass {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func (assistantResource AssistantResource) handleWebhook(w http.ResponseWriter, r *http.Request) {
	log.WithField("headers", r.Header).Info("Hit webhook")
	var err error
	var unmar jsonpb.Unmarshaler
	unmar.AllowUnknownFields = true

	wr := dialogflow.WebhookRequest{}
	if err = unmar.Unmarshal(r.Body, &wr); err != nil {
		log.WithError(err).Error("Couldn't Unmarshal request to jsonpb")
		w.WriteHeader(400)
		return
	}
	log.WithField("webhook", wr).Info("Webhook parsed request")
	if wr.QueryResult.Action == "input.toggle" {
		assistantResource.toggle(w, r)
		return
	}
	if wr.QueryResult.Action == "input.state" {
		assistantResource.getState(w, r)
		return
	}
	w.WriteHeader(200)
	state := assistantResource.doorController.getDoorState()
	json.NewEncoder(w).Encode(dialogflow.WebhookResponse{
		FulfillmentText: fmt.Sprintf("The garage door is %s, what's next?", state.description()),
	})
}

func (assistantResource AssistantResource) toggle(w http.ResponseWriter, r *http.Request) {
	state := assistantResource.doorController.getDoorState()
	action := "Closing"
	if state == closed {
		action = "Opening"
	}
	json.NewEncoder(w).Encode(dialogflow.WebhookResponse{
		FulfillmentText: fmt.Sprintf("%s the garage door", action),
	})
	assistantResource.doorController.toggleDoor()
	w.WriteHeader(200)
}

func (assistantResource AssistantResource) getState(w http.ResponseWriter, r *http.Request) {
	state := assistantResource.doorController.getDoorState()
	json.NewEncoder(w).Encode(dialogflow.WebhookResponse{
		FulfillmentText: fmt.Sprintf("The garage door is %s", state.description()),
	})
	w.WriteHeader(200)
}
