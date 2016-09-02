package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
	"time"
	"fmt"
)

const TIME_TO_CLOSE  = 1 * time.Minute

type GarageDoorResource struct {
	userDao		UserDao
	pinDao		PinDao
	doorController	DoorController
}

func (e GarageDoorResource) register(container *restful.Container) {
	ws := new(restful.WebService)

	ws.Path("/garage").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("toggle").To(e.toggleGarage))
	ws.Route(ws.GET("state").To(e.getState))
	ws.Route(ws.POST("one-time-pin/{oneTimePin}").Consumes("application/x-www-form-urlencoded").To(e.useOneTimePin))

	container.Add(ws)
}

func (e GarageDoorResource) useOneTimePin(request *restful.Request, response *restful.Response) {
	oneTimePin := request.PathParameter("oneTimePin")
	log.Infof("Using one time pin: [%s] to toggle garage", oneTimePin)
	usedDate, err := e.pinDao.getPinUsedDate(oneTimePin)
	if usedDate > 0 {
		log.Infof("Pin has already been used")
		response.WriteHeaderAndEntity(401, "Pin has already been used")
	} else if err != nil {
		log.WithError(err).Error("Could not get pin used date for [%s]", oneTimePin)
		response.WriteHeaderAndEntity(500, "Failed to open garage")
	} else {
		err = e.pinDao.use(oneTimePin)
		if err != nil {
			log.WithError(err).Errorf("Could not use pin: [%s]", oneTimePin)
			response.WriteHeaderAndEntity(401, "Pin has already been used")
			return
		}
		e.doorController.toggleDoor()
		response.WriteHeaderAndEntity(202, fmt.Sprintf("Opening garage, it will close in %v seconds",
			TIME_TO_CLOSE.Seconds()))
		go e.closeGarage(oneTimePin)
	}
}

func (e GarageDoorResource) closeGarage(pin string)  {
	time.Sleep(TIME_TO_CLOSE)
	log.Infof("Closing garage for pin: [%s]", pin)
	e.doorController.toggleDoor()
}

func (e GarageDoorResource) toggleGarage(request *restful.Request, response *restful.Response) {
	token := request.HeaderParameter("X-Auth-Token")
	user := e.userDao.getUserByToken(token)
	log.Infof("%s is opening or closing garage", user.Email)
	e.doorController.toggleDoor()
	response.WriteHeader(202)
}

func (e GarageDoorResource) getState(request *restful.Request, response *restful.Response) {
	log.Debug("Getting garage state")
	state := e.doorController.getDoorState()
	response.WriteAsJson(*newState(state))
}
