package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/namsral/flag"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"github.com/sourcegraph/go-ses"
)

var (
	relayPin         = flag.Int("relay", 0, "The relay pin number on the raspberry pi")
	contactSwitchPin = flag.Int("switch", 0, "The contact switch pin number on the raspberry pi")
	port             = flag.Int("port", 8080, "The port the server is listening on")
	database         = flag.String("db", "gogo-garage-opener.db", "The database file")
	notification     = flag.Duration("notification", time.Second*0, "The time to wait in minutes before sending a warning email")
	autoclose        = flag.Bool("autoclose", true, "Should auto close between 10pm-8am")
	rs               = flag.String("rs", "open.mygaragedoor.space", "Domain of the resource sever (raspberry pi)")
	as               = flag.String("as", "gogo-garage-opener.eu.auth0.com", "Domain of the authorisation sever (auth0 api)")
)

func main() {
	log.SetLevel(log.DebugLevel)
	flag.Parse()
	logConfiguration()

	db := initialise(*database)
	userDao := UserDao{db}
	pinDao := PinDao{db}
	garageDoorDao := SqliteGarageDoorDao{db}
	garageDoorDao.init()
	noop := *relayPin == 0 && *contactSwitchPin == 0
	log.WithField("NOOP", noop).Info("Running in mode")
	doorController := getDoorController(noop)
	router := registerResources(userDao, pinDao, garageDoorDao, doorController)

	if *autoclose {
		log.Info("Monitoring garage door to autoclose")
		go autoCloseMonitoring(doorController, userDao, garageDoorDao)
	}

	server := &http.Server{Addr: ":" + strconv.Itoa(*port), Handler: router}
	defer doorController.close()
	log.Fatal(server.ListenAndServe())
}

func autoCloseMonitoring(doorController DoorController, userDao UserDao, garageDoorDao GarageDoorDao) {
	autoclose := NewAutoclose(doorController, garageDoorDao)
	shouldNotify := *notification != time.Second*0
	for true {
		autoclose = autoclose.resetShouldCloseAndStayOpenTimes()
		message := fmt.Sprintf("Garage door has been left open for %v", autoclose.openDuration)
		if autoclose.autoClose() {
			log.WithField("message", message).Info("Sending emails for close notification")
			sendMail(userDao, "Autoclose: Garage door left open", message)
		} else if shouldNotify && autoclose.openDuration > *notification {
			log.WithField("message", message).Info("Sending emails for open notification")
			sendMail(userDao, "Notification: Garage door left open", message)
		}
		time.Sleep(time.Minute)
	}
}

func registerResources(userDao UserDao, pinDao PinDao, garageDoorDao GarageDoorDao, doorController DoorController) *mux.Router {
	userResource := UserResource{userDao: userDao, pinDao: pinDao}
	garageDoorResource := GarageDoorResource{userDao: userDao, pinDao: pinDao, garageDoorDao: garageDoorDao, doorController: doorController}
	router := mux.NewRouter()
	userResource.register(router)
	garageDoorResource.register(router)
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}).Name("ping")
	return router
}

func sendMail(userDao UserDao, title string, body string) {
	for _, email := range userDao.getSubscribedUserEmails() {
		_, err := ses.EnvConfig.SendEmail("garagedoor@mygaragedoor.tech", email, title, body)
		if err != nil {
			log.WithError(err).Error("Error sending email")
		}
	}
}

func logConfiguration() {
	log.
		WithField("relayPin", *relayPin).
		WithField("contactSwitchPin", *contactSwitchPin).
		WithField("port", *port).
		WithField("database", *database).
		WithField("notification", *notification).
		WithField("autoclose", *autoclose).
		WithField("rs", *rs).
		WithField("as", *as).
		Debug("Configuration")
}

func getDoorController(noop bool) DoorController {
	var doorController DoorController
	if noop {
		doorController = NoopDoorController{}
	} else {
		doorController = NewRaspberryPiDoorController(*relayPin, *contactSwitchPin)
	}
	return doorController
}
