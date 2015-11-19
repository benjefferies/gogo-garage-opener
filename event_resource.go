package main

import (
	"github.com/emicklei/go-restful"
	log "github.com/Sirupsen/logrus"
	"time"
	"strconv"
)

type Event struct {
	EventTime time.Time
	Open bool
}

type EventResource struct {
	events map[time.Time]Event
}

func (e EventResource) Register (container *restful.Container) {
	ws := new(restful.WebService)

	ws.Path("/event").
	Consumes(restful.MIME_JSON, restful.MIME_JSON).
	Produces(restful.MIME_JSON, restful.MIME_JSON)

	ws.Route(ws.POST("").To(e.createEvent))
	ws.Route(ws.GET("").To(e.findEvents))

	container.Add(ws)
}

func (e EventResource) createEvent(request *restful.Request, response *restful.Response) {
	event := Event{time.Now(), true}
	log.Debug("Storing event: " + event.EventTime.String() + ", open: " + strconv.FormatBool(event.Open))
	e.events[event.EventTime] = event
}

func (e EventResource) findEvents(request *restful.Request, response *restful.Response) {
	var events []Event
	for _,v := range e.events {
		events = append(events, v)
	}
	response.WriteEntity(events);
}