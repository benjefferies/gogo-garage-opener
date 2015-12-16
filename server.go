package main
import (
	"github.com/emicklei/go-restful"
	"net/http"
	 log "github.com/Sirupsen/logrus"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"flag"
	"strconv"
)

const database = "gogo-garage-opener.db"
const port = 8080

func main() {
	log.SetLevel(log.DebugLevel)
	databaseFlag, relayPinFlag, portFlag := flags();
	db, err := sql.Open("sqlite3", databaseFlag)
	if err != nil {
		log.Fatalf("Failed to open db [%s] errors [%s]", databaseFlag, err)
	}

	defer db.Close()
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS user (email TEXT NOT NULL PRIMARY KEY, password TEXT, longitude REAL, latitude REAL, time DATETIME, duration INTEGER, distance INTEGER);")
	if err != nil {
		log.Fatalf("Could not create table user [%s]", err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS event (timestamp DATETIME  NOT NULL PRIMARY KEY, email TEXT);")
	if err != nil {
		log.Fatal("Could not create table event [%s]", err)
	}

	userDao := UserDao{*db};
	u := UserResource{userDao}
	e := EventResource{eventDao:EventDao{*db}, userDao:userDao, doorController:DoorController{relayPinFlag}}

	wsContainer := restful.NewContainer()
	u.Register(wsContainer)
	e.Register(wsContainer)

	server := &http.Server{Addr: ":"+strconv.Itoa(portFlag), Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}

// TODO consider wrapping in struct
func flags() (databaseFlag string, relayPinFlag int, portFlag int) {
	flag.IntVar(&relayPinFlag, "r", -1, "The relay pin number on the raspberry pi")
	flag.IntVar(&portFlag, "p", port, "The port the server is listening on")
	flag.StringVar(&databaseFlag, "db", database, "The database file")
	flag.Parse()
	if (relayPinFlag == -1) {
		log.Fatal("Relay pin not set")
	} else {
		log.Debugf("Relay pin set to %d", relayPinFlag)
	}
	log.Debugf("Port set to [%d]", portFlag)
	log.Debugf("Database set to [%s]", databaseFlag)
	return databaseFlag, relayPinFlag, portFlag
}