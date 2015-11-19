package main

import (
	"github.com/emicklei/go-restful"
	"time"
)

type Event struct {
	EventTime time.Time
	Open bool
}

type EventResource struct {
	eventDao EventDao
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
	e.eventDao.createEvent(event)
}

func (e EventResource) findEvents(request *restful.Request, response *restful.Response) {
	response.WriteEntity(e.eventDao.getEvents());
}