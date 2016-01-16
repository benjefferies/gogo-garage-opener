package main

import (
	"github.com/emicklei/go-restful"
	log "github.com/Sirupsen/logrus"
)

type EventResource struct {
	userDao UserDao
	doorController DoorController
	distanceUtil DistanceUtil
}

func NewState(state int8) *State {
	var description string
	if (state == 0) {
		description = "Door 1 - Closed"
	} else {
		description = "Door 1 - Open"
	}
	return &State{State: state, Description: description}
}

type State struct {
	State int8
	Description string
}

func (e EventResource) Register (container *restful.Container) {
	ws := new(restful.WebService)

	ws.Path("/garage").
	Consumes(restful.MIME_JSON, restful.MIME_JSON).
	Produces(restful.MIME_JSON, restful.MIME_JSON)

	ws.Route(ws.POST("geo/{email}/{latitude}/{longitude}").To(e.openGarageByLocation))
	ws.Route(ws.POST("toggle").To(e.toggleGarage))
	ws.Route(ws.GET("state").To(e.getState))

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

	if (timeToOpen(timeWindows) && !hasFiredOpenEvent(user, timeWindows)) {
		arrivalDuration := e.distanceUtil.getArrivalTime(user, parseFloat64(latitude), parseFloat64(longitude))
		if &arrivalDuration != nil && arrivalDuration.Seconds() < 60 {
			log.Debug("Within 60 seconds of destination, opening door")
			e.doorController.toggleDoor()
			e.userDao.updateLastOpen(user)
		}
	}
}

func (e EventResource) toggleGarage(request *restful.Request, response *restful.Response) {
	email := request.PathParameter("email")
	user := e.userDao.getUser(email)
	log.Debugf("%s is opening garage", user.Email)
	e.doorController.toggleDoor()
	e.userDao.updateLastOpen(user)
}

func (e EventResource) getState(request *restful.Request, response *restful.Response) {
	log.Debug("Getting garage state")
	state := e.doorController.getDoorState()
	response.WriteAsJson(*NewState(int8(state)))
}