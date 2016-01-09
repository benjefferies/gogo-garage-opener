package main
import (
	"github.com/emicklei/go-restful"
	"github.com/satori/go.uuid"
	"strings"
	"github.com/Sirupsen/logrus"
)

type User struct {
	Email, Password, Token string
	Latitude, Longitude float64
}

type UserResource struct {
	userDao UserDao
}

func (u UserResource) Register (container *restful.Container) {
	ws := new(restful.WebService)

	ws.Path("/user").
	Consumes(restful.MIME_JSON, restful.MIME_JSON).
	Produces(restful.MIME_JSON, restful.MIME_JSON)

	ws.Route(ws.POST("register").To(u.createUser))
	ws.Route(ws.POST("login").To(u.login))
	ws.Route(ws.GET("{email}").To(u.getUser).Consumes("application/json").Produces("application/json"))
	container.Add(ws)
}

func (u UserResource) createUser(request *restful.Request, response *restful.Response) {
	user := new(User)
	request.ReadEntity(&user)
	u.userDao.createUser(*user)
}

func (u UserResource) login(request *restful.Request, response *restful.Response) {
	loginUser := new(User)
	request.ReadEntity(&loginUser)
	user := u.userDao.getUser(loginUser.Email)
	if (strings.EqualFold(user.Password, loginUser.Password)) {
		logrus.Debugf("Login successful for [%s]", user.Email)
		user.Token = uuid.NewV4().String()
		u.userDao.updateToken(user)
	}
}

func (u UserResource) getUser(request *restful.Request, response *restful.Response) {
	email := request.PathParameter("email")
	user := u.userDao.getUser(email)
	response.WriteEntity(user);
}