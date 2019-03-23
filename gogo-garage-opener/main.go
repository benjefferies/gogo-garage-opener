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
	noop := *relayPin == 0 && *contactSwitchPin == 0
	log.WithField("NOOP", noop).Info("Running in mode")
	doorController := getDoorController(noop)
	router := registerResources(userDao, pinDao, doorController)

	shouldNotify := *notification != time.Second*0
	if shouldNotify {
		log.Info("Monitoring garage door to send alerts")
		go leftOpenAlertMonitoring(doorController, userDao)
	}

	if *autoclose {
		log.Info("Monitoring garage door to autoclose")
		go autoCloseMonitoring(doorController, userDao)
	}

	server := &http.Server{Addr: ":" + strconv.Itoa(*port), Handler: router}
	defer doorController.close()
	log.Fatal(server.ListenAndServe())
}

func autoCloseMonitoring(doorController DoorController, userDao UserDao) {
	autoclose := NewAutoclose(doorController)
	for true {
		if autoclose.autoClose() {
			sendMail(userDao, "Autoclose: Garage door left open", fmt.Sprintf("Garage door appears to be left open at %s", time.Now().Format("3:04 PM")))
		}
		time.Sleep(time.Minute)
	}
}

func registerResources(userDao UserDao, pinDao PinDao, doorController DoorController) *mux.Router {
	userResource := UserResource{userDao: userDao, pinDao: pinDao}
	garageDoorResource := GarageDoorResource{userDao: userDao, pinDao: pinDao, doorController: doorController}
	router := mux.NewRouter()
	userResource.register(router)
	garageDoorResource.register(router)
	return router
}

func leftOpenAlertMonitoring(doorController DoorController, userDao UserDao) {
	nilTime := time.Time{}
	lastOpened := nilTime
	for true {
		if doorController.getDoorState().isOpen() {
			if lastOpened == nilTime {
				log.Info("Setting lastOpened")
				lastOpened = time.Now()
			}
			now := time.Now()
			openTooLong := lastOpened.Add(*notification)
			if now.After(openTooLong) {
				log.Info("Sending emails for open notification")
				sendMail(userDao, "Notification: Garage door left open", fmt.Sprintf("Garage door has been left open since %s", lastOpened.Format("3:04 PM")))
			}
		} else {
			lastOpened = nilTime
		}
		time.Sleep(*notification)
	}
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
