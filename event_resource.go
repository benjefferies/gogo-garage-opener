package main

import (
	"github.com/emicklei/go-restful"
	log "github.com/Sirupsen/logrus"
	"strconv"
)

type EventResource struct {
	userDao        UserDao
	doorController DoorController
	distanceUtil   DistanceCalculator
}

func newState(state int8) *State {
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

func (e EventResource) register(container *restful.Container) {
	ws := new(restful.WebService)

	ws.Path("/garage").
	Consumes(restful.MIME_JSON, restful.MIME_JSON).
	Produces(restful.MIME_JSON, restful.MIME_JSON)

	ws.Route(ws.POST("geo/{latitude}/{longitude}/arrival/{arrival}").To(e.openGarageByLocation))
	ws.Route(ws.POST("toggle").To(e.toggleGarage))
	ws.Route(ws.GET("state").To(e.getState))

	container.Add(ws)
}

func parseFloat64(floatValue string) float64 {
	float, err := strconv.ParseFloat(floatValue, 64)
	if (err != nil) {
		log.Error(err)
	}
	return float
}

func (e EventResource) openGarageByLocation(request *restful.Request, response *restful.Response) {
	token := request.HeaderParameter("X-Auth-Token")
	user := e.userDao.getUserByToken(token)
	log.Debugf("Found user, email:[%s]", email)
	latitude := request.PathParameter("latitude")
	longitude := request.PathParameter("longitude")
	arrival := request.PathParameter("arrival")
	log.Debugf("Request geolocation [%s, %s]", latitude, longitude)

	arrivalDuration := e.distanceUtil.getTimeToArrive(user, parseFloat64(latitude), parseFloat64(longitude))
	if &arrivalDuration != nil && arrivalDuration.Seconds() < parseFloat64(arrival) {
		log.Infof("Within %s seconds of destination, opening door for user %s", arrival, email)
		e.doorController.toggleDoor()
		e.userDao.updateLastOpen(user)
	}
}

func (e EventResource) toggleGarage(request *restful.Request, response *restful.Response) {
	token := request.HeaderParameter("X-Auth-Token")
	user := e.userDao.getUserByToken(token)
	log.Infof("%s is opening or closing garage", user.Email)
	e.doorController.toggleDoor()
	e.userDao.updateLastOpen(user)
}

func (e EventResource) getState(request *restful.Request, response *restful.Response) {
	log.Debug("Getting garage state")
	state := e.doorController.getDoorState()
	response.WriteAsJson(*newState(int8(state)))
}