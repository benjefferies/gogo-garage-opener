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

func (this GarageDoorResource) register(container *restful.Container) {
	ws := new(restful.WebService)

	ws.Path("/garage").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("toggle").To(this.toggleGarage))
	ws.Route(ws.GET("state").To(this.getState))
	ws.Route(ws.POST("one-time-pin/{oneTimePin}").Consumes("application/x-www-form-urlencoded").To(this.useOneTimePin))

	container.Add(ws)
}

func (this GarageDoorResource) useOneTimePin(request *restful.Request, response *restful.Response) {
	oneTimePin := request.PathParameter("oneTimePin")
	log.Infof("Using one time pin: [%s] to toggle garage", oneTimePin)
	usedDate, err := this.pinDao.getPinUsedDate(oneTimePin)
	if usedDate > 0 {
		log.Infof("Pin has already been used")
		response.WriteHeaderAndEntity(401, "Pin has already been used")
	} else if err != nil {
		log.WithError(err).Errorf("Could not get pin used date for [%s]", oneTimePin)
		response.WriteHeaderAndEntity(500, "Failed to open garage")
	} else {
		err = this.pinDao.use(oneTimePin)
		if err != nil {
			log.WithError(err).Errorf("Could not use pin: [%s]", oneTimePin)
			response.WriteHeaderAndEntity(401, "Pin has already been used")
			return
		}
		this.doorController.toggleDoor()
		response.WriteHeaderAndEntity(202, fmt.Sprintf("Opening garage, it will close in %v seconds",
			TIME_TO_CLOSE.Seconds()))
		go this.closeGarage(oneTimePin)
	}
}

func (this GarageDoorResource) closeGarage(pin string)  {
	time.Sleep(TIME_TO_CLOSE)
	log.Infof("Closing garage for pin: [%s]", pin)
	this.doorController.toggleDoor()
}

func (this GarageDoorResource) toggleGarage(request *restful.Request, response *restful.Response) {
	token := request.HeaderParameter("X-Auth-Token")
	user := this.userDao.getUserByToken(token)
	log.Infof("%s is opening or closing garage", user.Email)
	this.doorController.toggleDoor()
	response.WriteHeader(202)
}

func (this GarageDoorResource) getState(request *restful.Request, response *restful.Response) {
	log.Debug("Getting garage state")
	state := this.doorController.getDoorState()
	response.WriteAsJson(*newState(state))
}
