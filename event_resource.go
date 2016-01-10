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

func (e EventResource) Register (container *restful.Container) {
	ws := new(restful.WebService)

	ws.Path("/garage").
	Consumes(restful.MIME_JSON, restful.MIME_JSON).
	Produces(restful.MIME_JSON, restful.MIME_JSON)

	ws.Route(ws.POST("geo/{email}/{latitude}/{longitude}").To(e.openGarageByLocation))
	ws.Route(ws.POST("toggle").To(e.toggleGarage))

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