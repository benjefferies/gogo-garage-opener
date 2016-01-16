package main
import (
	"github.com/emicklei/go-restful"
	"github.com/satori/go.uuid"
	log "github.com/Sirupsen/logrus"
	"time"
)

type User struct {
	Email, Password, Token string
	Latitude, Longitude float64
	LastOpen time.Time
	Approved bool
}

type UserResource struct {
	userDao UserDao
}

func (u UserResource) Register (container *restful.Container) {
	ws := new(restful.WebService)

	ws.Path("/user").
	Consumes(restful.MIME_JSON, restful.MIME_JSON).
	Produces(restful.MIME_JSON, restful.MIME_JSON)

	ws.Route(ws.POST("register").To(u.CreateUser))
	ws.Route(ws.POST("login").To(u.login))
	container.Add(ws)
}

func (u UserResource) CreateUser(request *restful.Request, response *restful.Response) {
	user := new(User)
	request.ReadEntity(&user)
	user.Approved = false
	u.createUser(user)
}

func (u UserResource) createUser(user *User) {
//	u.userDao.createUser(*user)
}

func (u UserResource) login(request *restful.Request, response *restful.Response) {
	loginUser := new(User)
	request.ReadEntity(&loginUser)
	user := u.userDao.getUser(loginUser.Email)
	hashedPassword := hashedPassword(*loginUser)
	log.Debugf("Comparing passwords from db user [%s] to request user [%s]", user.Password, hashedPassword)
	if (user.Password == hashedPassword) {
		log.Debugf("Login successful for [%s]", user.Email)
		user.Token = uuid.NewV4().String()
		u.userDao.updateToken(user)
		response.Header().Set("X-Auth-Token", user.Token)
		log.Debugf("Setting X-Auth-Token to [%s]", user.Token)
	} else {
		response.WriteErrorString(401, "401: Not Authorized")
	}
}