package main

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
)

// timeToClose time to automatically close the door after opening with pin
const timeToClose = 1 * time.Minute

// GarageDoorResource API for interacting with garage door
type GarageDoorResource struct {
	userDao        UserDao
	pinDao         PinDao
	doorController DoorController
}

func (garageDoorResource GarageDoorResource) register(container *restful.Container) {
	ws := new(restful.WebService)

	ws.Path("/garage").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("toggle").To(garageDoorResource.toggleGarage))
	ws.Route(ws.GET("state").To(garageDoorResource.getState))
	ws.Route(ws.POST("one-time-pin/{oneTimePin}").Consumes("application/x-www-form-urlencoded").To(garageDoorResource.useOneTimePin))

	container.Add(ws)
}

func (garageDoorResource GarageDoorResource) useOneTimePin(request *restful.Request, response *restful.Response) {
	oneTimePin := request.PathParameter("oneTimePin")
	log.Infof("Using one time pin: [%s] to toggle garage", oneTimePin)
	usedDate, err := garageDoorResource.pinDao.getPinUsedDate(oneTimePin)
	if usedDate > 0 {
		log.Infof("Pin has already been used")
		response.WriteHeaderAndEntity(401, "Pin has already been used")
	} else if err != nil {
		log.WithError(err).Errorf("Could not get pin used date for [%s]", oneTimePin)
		response.WriteHeaderAndEntity(500, "Failed to open garage")
	} else {
		err = garageDoorResource.pinDao.use(oneTimePin)
		if err != nil {
			log.WithError(err).Errorf("Could not use pin: [%s]", oneTimePin)
			response.WriteHeaderAndEntity(401, "Pin has already been used")
			return
		}
		garageDoorResource.doorController.toggleDoor()
		response.WriteHeaderAndEntity(202, fmt.Sprintf("Opening garage, it will close in %v seconds",
			timeToClose.Seconds()))
		go garageDoorResource.closeGarage(oneTimePin)
	}
}

func (garageDoorResource GarageDoorResource) closeGarage(pin string) {
	time.Sleep(timeToClose)
	log.Infof("Closing garage for pin: [%s]", pin)
	garageDoorResource.doorController.toggleDoor()
}

func (garageDoorResource GarageDoorResource) toggleGarage(request *restful.Request, response *restful.Response) {
	token := request.HeaderParameter("X-Auth-Token")
	user := garageDoorResource.userDao.getUserByToken(token)
	log.Infof("%s is opening or closing garage", user.Email)
	garageDoorResource.doorController.toggleDoor()
	response.WriteHeader(202)
}

func (garageDoorResource GarageDoorResource) getState(request *restful.Request, response *restful.Response) {
	log.Debug("Getting garage state")
	state := garageDoorResource.doorController.getDoorState()
	response.WriteAsJson(*newState(state))
}
