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

var (
	relayPinFlag = flag.Int("r", -1, "The relay pin number on the raspberry pi")
	portFlag = flag.Int("p", port, "The port the server is listening on")
	databaseFlag = flag.String("db", database, "The database file")
	apiKeyFlag = flag.String("key", "", "The google maps API key (with distance matrix api enabled)")
)

func main() {
	flag.Parse()
	log.SetLevel(log.DebugLevel)
	logConfiguration()
	db, err := sql.Open("sqlite3", *databaseFlag)
	if err != nil {
		log.Fatalf("Failed to open db [%s] errors [%s]", *databaseFlag, err)
	}

	defer db.Close()
	var errs []error = make([]error, 0)
	// Create user table
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS user (email TEXT NOT NULL PRIMARY KEY, password TEXT, token TEXT, longitude REAL, latitude REAL);")
	errs = append(errs, err)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS user_time (email TEXT NOT NULL, time DATETIME, duration INTEGER);")
	errs = append(errs, err)
	if (len(errs) > 0) {
		for _, err := range errs {
			log.Errorf("%v", err)
		}
	}

	userDao := UserDao{*db};
	u := UserResource{userDao}
	e := EventResource{eventDao:EventDao{*db}, userDao:userDao, doorController:DoorController{*relayPinFlag}, distanceUtil:DistanceUtil{apiKey:*apiKeyFlag}}

	wsContainer := restful.NewContainer()
	u.Register(wsContainer)
	e.Register(wsContainer)

	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-My-Header"},
		AllowedHeaders: []string{"Content-Type", "Accept"},
		CookiesAllowed: false,
		Container:      wsContainer}
	wsContainer.Filter(cors.Filter)

	server := &http.Server{Addr: ":"+strconv.Itoa(*portFlag), Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}

func logConfiguration() {
	log.Debugf("Relay pin %d", *relayPinFlag)
	log.Debugf("Api key %s", *apiKeyFlag)
	log.Debugf("Database file %s", *databaseFlag)
	log.Debugf("Webserver port %d", *portFlag)
}