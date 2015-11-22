package main

import (
	"github.com/emicklei/go-restful"
	"time"
	log "github.com/Sirupsen/logrus"
)

type Event struct {
	EventTime time.Time
	Open bool
}

type EventResource struct {
	eventDao EventDao
	userDao UserDao
}

func (e EventResource) Register (container *restful.Container) {
	ws := new(restful.WebService)

	ws.Path("/event").
	Consumes(restful.MIME_JSON, restful.MIME_JSON).
	Produces(restful.MIME_JSON, restful.MIME_JSON)

	ws.Route(ws.POST("{email}/{latitude}/{longitude}").To(e.garageEvent))
	ws.Route(ws.GET("").To(e.findEvents))

	container.Add(ws)
}

func (e EventResource) garageEvent(request *restful.Request, response *restful.Response) {
	email := request.PathParameter("email")
	user := e.userDao.getUser(email)
	log.Debugf("Found user, email:[%s]", user.Email)
	latitude := request.PathParameter("latitude")
	longitude := request.PathParameter("longitude")
	log.Debugf("Request geolocation [%s, %s]", latitude, longitude)


	distanceInMetres := distance(parseFloat64(latitude), parseFloat64(longitude), user.Latitude, user.Longitude)

	log.Infof("%s is %s metres away", user.Email, distanceInMetres)

	if (distanceInMetres > 100) {
		return
	}

	events := e.eventDao.getEvents();

	log.Debugf("Found [%s] events", len(events))
	var open bool
	if (len(events) == 0) {
		open = true
	} else {
		open = !events[0].Open;
	}

	event := Event{time.Now(), open}
	e.eventDao.createEvent(event)
}

func (e EventResource) findEvents(request *restful.Request, response *restful.Response) {
	response.WriteEntity(e.eventDao.getEvents());
}