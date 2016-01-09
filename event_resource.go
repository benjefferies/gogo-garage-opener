package main

import (
	"github.com/emicklei/go-restful"
	"time"
	log "github.com/Sirupsen/logrus"
)

type Event struct {
	EventTime time.Time
	Email string
}

type EventResource struct {
	eventDao EventDao
	userDao UserDao
	doorController DoorController
	distanceUtil DistanceUtil
}

func (e EventResource) Register (container *restful.Container) {
	ws := new(restful.WebService)

	ws.Path("/garage").
	Consumes(restful.MIME_JSON, restful.MIME_JSON).
	Produces(restful.MIME_JSON, restful.MIME_JSON)

	ws.Route(ws.POST("geo/{email}/{latitude}/{longitude}").To(e.openGarageByLocation))
	ws.Route(ws.POST("toggle").To(e.toggleGarage))
	ws.Route(ws.GET("").To(e.findEvents))

	container.Add(ws)
}

func (e EventResource) openGarageByLocation(request *restful.Request, response *restful.Response) {
	email := request.PathParameter("email")
	user := e.userDao.getUser(email)
	timeWindows := e.userDao.getTimes(user)
	log.Debugf("Found user, email:[%s]", user.Email)
	latitude := request.PathParameter("latitude")
	longitude := request.PathParameter("longitude")
	log.Debugf("Request geolocation [%s, %s]", latitude, longitude)

	if (timeToOpen(timeWindows)) {
		arrivalDuration := e.distanceUtil.getArrivalTime(user, parseFloat64(latitude), parseFloat64(longitude))
		if &arrivalDuration != nil && arrivalDuration.Seconds() < 60 {
			log.Debug("Within 60 seconds of destination, opening door")
			e.doorController.toggleDoor()
		}
	}

	event := Event{time.Now(), email}
	e.eventDao.createEvent(event)
}

func (e EventResource) toggleGarage(request *restful.Request, response *restful.Response) {
	email := request.PathParameter("email")
	user := e.userDao.getUser(email)
	log.Debugf("%s is opening garage", user.Email)
	e.doorController.toggleDoor()
	event := Event{time.Now(), email}
	e.eventDao.createEvent(event)
}

func (e EventResource) findEvents(request *restful.Request, response *restful.Response) {
	response.WriteEntity(e.eventDao.getEvents());
}