package main

import (
	"database/sql"
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

const database = "gogo-garage-opener.db"
const port = 8080

var (
	relayPinFlag         = flag.Int("relay", 14, "The relay pin number on the raspberry pi")
	contactSwitchPinFlag = flag.Int("switch", 7, "The contact switch pin number on the raspberry pi")
	portFlag             = flag.Int("port", port, "The port the server is listening on")
	databaseFlag         = flag.String("db", database, "The database file")
	noop                 = flag.Bool("noop", false, "Noop can be ran without the raspberry pi")
	notification         = flag.Duration("notification", time.Second*0, "The time to wait in minutes before sending a warning email")
	autoclose            = flag.Bool("autoclose", true, "Should auto close between 10pm-8am")
)

func main() {
	flag.Parse()
	log.SetLevel(log.InfoLevel)
	logConfiguration()
	db, err := sql.Open("sqlite3", *databaseFlag)
	if err != nil {
		log.WithError(err).Fatalf("Failed to open db [%s]", *databaseFlag)
	}

	defer db.Close()

	setupTables(db)

	userDao := UserDao{db}
	pinDao := PinDao{db}
	userResource := UserResource{userDao: userDao, pinDao: pinDao}
	doorController := getDoorController(*noop)

	defer doorController.close()
	garageDoorResource := GarageDoorResource{userDao: userDao, pinDao: pinDao, doorController: doorController}
	router := mux.NewRouter()
	userResource.register(router)
	garageDoorResource.register(router)

	server := &http.Server{Addr: ":" + strconv.Itoa(*portFlag), Handler: router}
	if *notification > time.Second*0 {
		log.Infof("Monitoring garage door to send alerts")
		go leftOpenAlertMonitoring(doorController, userDao)
	}

	if *autoclose {
		log.Infof("Monitoring garage door to autoclose")
		go autoCloseMonitoring(doorController, userDao)
	}
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

func leftOpenAlertMonitoring(doorController DoorController, userDao UserDao) {
	nilTime := time.Time{}
	lastOpened := nilTime
	for true {
		if doorController.getDoorState() == open {
			if lastOpened == nilTime {
				log.Infof("Setting lastOpened")
				lastOpened = time.Now()
			}
			now := time.Now()
			openTooLong := lastOpened.Add(*notification)
			if now.After(openTooLong) {
				log.Infof("Sending emails for open notification")
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
			log.Errorf("Error sending email: %s\n", err)
		}
	}
}

func logConfiguration() {
	log.Debugf("Relay pin %d", *relayPinFlag)
	log.Debugf("Sensor pin %d", *contactSwitchPinFlag)
	log.Debugf("Database file %s", *databaseFlag)
	log.Debugf("Webserver port %d", *portFlag)
}

func setupTables(db *sql.DB) {
	// Create user table
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS user (email TEXT NOT NULL PRIMARY KEY, token TEXT, subscribed BOOLEAN DEFAULT 1);")
	if err != nil {
		log.WithError(err).Fatal("Could not create user table")
	}

	// Create one time pin table
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS one_time_pin (pin TEXT NOT NULL PRIMARY KEY, created_by TEXT, created INTEGER, used INTEGER);")
	if err != nil {
		log.WithError(err).Fatal("Could not create one_time_pin table")
	}
}

func getDoorController(noop bool) DoorController {
	var doorController DoorController
	if noop {
		log.Info("Running in noop mode")
		doorController = NoopDoorController{}
	} else {
		doorController = NewRaspberryPiDoorController(*relayPinFlag, *contactSwitchPinFlag)
	}
	return doorController
}
