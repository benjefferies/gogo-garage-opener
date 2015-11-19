package main
import (
	"github.com/emicklei/go-restful"
	log "github.com/Sirupsen/logrus"
)


type User struct {
	Email, Password string
}

// TODO use dao
type UserResource struct {
	users map[string]User
}

func (u UserResource) Register (container *restful.Container) {
	ws := new(restful.WebService)

	ws.Path("/user").
	Consumes(restful.MIME_JSON, restful.MIME_JSON).
	Produces(restful.MIME_JSON, restful.MIME_JSON)

	ws.Route(ws.POST("").To(u.createUser))
	ws.Route(ws.GET("{email}").To(u.findUser))

	container.Add(ws)
}

func (u UserResource) createUser(request *restful.Request, response *restful.Response) {
	user := new(User)
	request.ReadEntity(&user)
	log.Debug("Storing email: " + user.Email)
	u.users[user.Email] = *user
}

func (u UserResource) findUser(request *restful.Request, response *restful.Response) {
	email := request.PathParameter("email")
	log.Debug("Finding user for email: " + email)
	user := u.users[email]
	response.WriteEntity(user);
}