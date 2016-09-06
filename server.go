package main

import (
	"database/sql"
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"strconv"
)

const database = "gogo-garage-opener.db"
const port = 8080

var (
	relayPinFlag         = flag.Int("r", 14, "The relay pin number on the raspberry pi")
	contactSwitchPinFlag = flag.Int("s", 7, "The contact switch pin number on the raspberry pi")
	portFlag             = flag.Int("p", port, "The port the server is listening on")
	databaseFlag         = flag.String("db", database, "The database file")
	email                = flag.String("email", "", "Specify email to create account")
	password             = flag.String("password", "", "Specify email to create account")
	noop                 = flag.Bool("noop", false, "Noop can be ran without the raspberry pi")
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

	setupTables(*db)

	userDao := UserDao{*db}
	created := createUser(userDao)
	if (!created) {
		pinDao := PinDao{*db}
		userResource := UserResource{userDao: userDao, pinDao: pinDao}
		doorController := getDoorController(*noop)

		defer doorController.close()
		garageDoorResource := GarageDoorResource{userDao: userDao, pinDao: pinDao, doorController: doorController}
		authFilter := AuthFilter{userDao: userDao}
		wsContainer := restful.NewContainer()
		userResource.register(wsContainer)
		garageDoorResource.register(wsContainer)

		cors := restful.CrossOriginResourceSharing{
			ExposeHeaders:  []string{"X-Auth-Token"},
			AllowedHeaders: []string{"Content-Type", "Accept", "X-Auth-Token"},
			CookiesAllowed: false,
			Container:      wsContainer}
		wsContainer.Filter(cors.Filter)
		wsContainer.Filter(authFilter.tokenFilter)

		server := &http.Server{Addr: ":" + strconv.Itoa(*portFlag), Handler: wsContainer}
		log.Fatal(server.ListenAndServe())
	}
}

func logConfiguration() {
	log.Debugf("Relay pin %d", *relayPinFlag)
	log.Debugf("Sensor pin %d", *contactSwitchPinFlag)
	log.Debugf("Database file %s", *databaseFlag)
	log.Debugf("Webserver port %d", *portFlag)
}

func setupTables(db sql.DB) {
	// Create user table
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS user (email TEXT NOT NULL PRIMARY KEY, password TEXT, token TEXT);")
	if err != nil {
		log.WithError(err).Fatal("Could not create user table")
	}

	// Create one time pin table
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS one_time_pin (pin TEXT NOT NULL PRIMARY KEY, created_by TEXT, created INTEGER, used INTEGER);")
	if err != nil {
		log.WithError(err).Fatal("Could not create one_time_pin table")
	}
}

func createUser(userDao UserDao) bool {
	if (*email != "" && password != nil) && (email != nil && *password != "") {
		user, err := User{Email: *email, Password: *password}.hashPassword()
		if err != nil {
			log.Fatalf("Failed to create user: %s" , email)
		}
		userDao.createUser(user)
		log.Infof("Created account email:%s. Exiting...", *email)
		return true
	}
	return false
}

func getDoorController(noop bool) DoorController {
	var doorController DoorController
	if (noop) {
		log.Info("Running in noop mode")
		doorController = NoopDoorController{}
	} else {
		doorController = NewRaspberryPiDoorController(*relayPinFlag, *contactSwitchPinFlag)
	}
	return doorController
}