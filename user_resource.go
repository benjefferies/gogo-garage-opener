package main
import (
	"github.com/emicklei/go-restful"
	"github.com/satori/go.uuid"
	log "github.com/Sirupsen/logrus"
	"time"
	"crypto/sha256"
	"fmt"
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

func (u UserResource) register(container *restful.Container) {
	ws := new(restful.WebService)

	ws.Path("/user").
	Consumes(restful.MIME_JSON).
	Produces(restful.MIME_JSON)

	ws.Route(ws.POST("login").To(u.login))
	container.Add(ws)
}

func (u UserResource) createUser(request *restful.Request, response *restful.Response) {
	user := new(User)
	request.ReadEntity(&user)
	user.Password = hashedPassword(*user)
	user.Approved = false
	u.userDao.createUser(*user)
}

func (u UserResource) login(request *restful.Request, response *restful.Response) {
	loginUser := new(User)
	request.ReadEntity(&loginUser)
	user := u.userDao.getUserByEmail(loginUser.Email)
	hashedPassword := hashedPassword(*loginUser)
	if (user.Password == hashedPassword) {
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

func hashedPassword(user User) string {
	hashedBytes := sha256.Sum256([]byte(user.Password))
	return fmt.Sprintf("%s", hashedBytes)
}