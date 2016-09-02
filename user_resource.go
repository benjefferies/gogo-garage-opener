package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"encoding/json"
	"fmt"
)

const ONE_TIME_PIN_USE  = `
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

type UserResource struct {
	userDao UserDao
	pinDao PinDao
}

func (u UserResource) register(container *restful.Container) {
	ws := new(restful.WebService)

	ws.Path("/user").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("login").To(u.login))
	ws.Route(ws.POST("one-time-pin").To(u.oneTimePin))
	ws.Route(ws.GET("one-time-pin/{oneTimePin}").To(u.useOneTimePin))
	container.Add(ws)
}

func (u UserResource) oneTimePin(request *restful.Request, response *restful.Response)  {
	token := request.HeaderParameter("X-Auth-Token")
	user := u.userDao.getUserByToken(token)
	log.Debugf("%s is creating new one time pin", user.Email)
	pin, err := u.pinDao.newOneTimePin(user)
	if err != nil {
		log.WithError(err).Error("Could not create one time pin")
		response.WriteHeader(500)
	}
	pinMap := map[string]interface{}{
		"pin": pin,
	}
	json, err := json.Marshal(pinMap)
	if err != nil {
		log.WithError(err).Error("Could not marshell one time pin")
		response.WriteHeader(500)
	}
	response.Write(json)
}

func (u UserResource) useOneTimePin(request *restful.Request, response *restful.Response)  {
	oneTimePin := request.PathParameter("oneTimePin")
	response.ResponseWriter.WriteHeader(200)
	response.ResponseWriter.Write([]byte(fmt.Sprintf(ONE_TIME_PIN_USE, TIME_TO_CLOSE.Seconds(), oneTimePin)))
}

func (u UserResource) login(request *restful.Request, response *restful.Response) {
	loginUser := new(User)
	request.ReadEntity(&loginUser)
	user := u.userDao.getUserByEmail(loginUser.Email)
	if bcrypt.CompareHashAndPassword([]byte(loginUser.Password), []byte(user.Password)) != nil {
		log.Infof("Login successful for [%s]", user.Email)
		user.Token = uuid.NewV4().String()
		u.userDao.setToken(user)
		response.Header().Set("X-Auth-Token", user.Token)
		log.Debugf("Setting X-Auth-Token to [%s]", user.Token)
	} else {
		log.Infof("Login failed for [%s]", user.Email)
		response.WriteErrorString(400, "400: Incorrect username or passwords")
	}
}
