package main
import (
	"github.com/emicklei/go-restful"
	"time"
)


type User struct {
	Email, Password string
	Latitude, Longitude float64
	Time time.Time
	Duration, Distance int32
}

type UserResource struct {
	userDao UserDao
}

func (u UserResource) Register (container *restful.Container) {
	ws := new(restful.WebService)

	ws.Path("/user").
	Consumes(restful.MIME_JSON, restful.MIME_JSON).
	Produces(restful.MIME_JSON, restful.MIME_JSON)

	ws.Route(ws.POST("").To(u.createUser))
	ws.Route(ws.GET("{email}").To(u.getUser))
	container.Add(ws)
}

func (u UserResource) createUser(request *restful.Request, response *restful.Response) {
	user := new(User)
	request.ReadEntity(&user)
	u.userDao.createUser(*user)
}

func (u UserResource) getUser(request *restful.Request, response *restful.Response) {
	email := request.PathParameter("email")
	user := u.userDao.getUser(email)
	response.WriteEntity(user);
}