package main
import (
	"github.com/emicklei/go-restful"
	"net/http"
	 log "github.com/Sirupsen/logrus"
	"time"
	"strconv"
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

type Event struct {
	EventTime time.Time
	Open bool
}
// TODO use dao
type EventResource struct {
	events map[time.Time]Event
}

func (e EventResource) Register (container *restful.Container) {
	ws := new(restful.WebService)

	ws.Path("/event").
	Consumes(restful.MIME_JSON, restful.MIME_JSON).
	Produces(restful.MIME_JSON, restful.MIME_JSON)

	ws.Route(ws.POST("").To(e.createEvent))
	ws.Route(ws.GET("").To(e.findEvents))

	container.Add(ws)
}

func (e EventResource) createEvent(request *restful.Request, response *restful.Response) {
	event := Event{time.Now(), true}
	log.Debug("Storing event: " + event.EventTime.String() + ", open: " + strconv.FormatBool(event.Open))
	e.events[event.EventTime] = event
}

func (e EventResource) findEvents(request *restful.Request, response *restful.Response) {
	var events []Event
	for _,v := range e.events {
		events = append(events, v)
	}
	response.WriteEntity(events);
}

func main() {
	log.SetLevel(log.DebugLevel)
	wsContainer := restful.NewContainer()

	u := UserResource{map[string]User{}}
	e := EventResource{map[time.Time]Event{}}

	u.Register(wsContainer)
	e.Register(wsContainer)

	server := &http.Server{Addr: ":8080", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}