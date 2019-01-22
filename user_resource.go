package main

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// oneTimePinUse static html of one time pin page
const oneTimePinUse = `
<html>
	</head>
	<body>
		<h1>One time pin garage door opener</h1>
		<br>
		<p>By clicking the button below it will open the garage door. The garage door will automatically close %v seconds after clicking the button</p>
		<br>
		<br>
		<form name="myform" action="/garage/one-time-pin/%s" method="post">
			<button>Open</button>
		</form>
	</body>
</html>
`

// UserResource API for users
type UserResource struct {
	userDao UserDao
	pinDao  PinDao
}

func (userResource UserResource) register(container *restful.Container) {
	ws := new(restful.WebService)

	ws.Path("/user").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("login").To(userResource.login))
	ws.Route(ws.POST("one-time-pin").To(userResource.oneTimePin))
	ws.Route(ws.GET("one-time-pin/{oneTimePin}").To(userResource.useOneTimePin))
	container.Add(ws)
}

func (userResource UserResource) oneTimePin(request *restful.Request, response *restful.Response) {
	token := request.HeaderParameter("X-Auth-Token")
	user := userResource.userDao.getUserByToken(token)
	log.Debugf("%s is creating new one time pin", user.Email)
	pin, err := userResource.pinDao.newOneTimePin(user)
	if err != nil {
		log.WithError(err).Error("Could not create one time pin")
		response.WriteHeader(500)
	}
	pinMap := map[string]interface{}{
		"pin": pin,
	}
	payload, err := json.Marshal(pinMap)
	if err != nil {
		log.WithError(err).Error("Could not marshell one time pin")
		response.WriteHeader(500)
	}
	response.Write(payload)
}

func (userResource UserResource) useOneTimePin(request *restful.Request, response *restful.Response) {
	oneTimePin := request.PathParameter("oneTimePin")
	response.ResponseWriter.WriteHeader(200)
	response.ResponseWriter.Write([]byte(fmt.Sprintf(oneTimePinUse, timeToClose.Seconds(), oneTimePin)))
}

func (userResource UserResource) login(request *restful.Request, response *restful.Response) {
	loginUser := new(User)
	request.ReadEntity(&loginUser)
	user := userResource.userDao.getUserByEmail(loginUser.Email)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		log.Infof("Login failed for [%s]", user.Email)
		response.WriteErrorString(400, "400: Incorrect username or passwords")
	} else if err != nil {
		log.Infof("Login failed for [%s]", user.Email)
		log.Errorf("%v", err)
		response.WriteErrorString(400, "400: Incorrect username or passwords")
	} else {
		log.Infof("Login successful for [%s]", user.Email)
		user.Token = uuid.Must(uuid.NewV4()).String()
		userResource.userDao.setToken(user)
		response.Header().Set("X-Auth-Token", user.Token)
		log.Debugf("Setting X-Auth-Token to [%s]", user.Token)
	}
}
