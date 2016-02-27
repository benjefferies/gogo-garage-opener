package main
import (
	"github.com/emicklei/go-restful"
	"net/http"
	 log "github.com/Sirupsen/logrus"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"flag"
	"strconv"
	"os"
)

const database = "gogo-garage-opener.db"
const port = 8080

var (
	relayPinFlag = flag.Int("r", 14, "The relay pin number on the raspberry pi")
	contactSwitchPinFlag = flag.Int("s", 7, "The contact switch pin number on the raspberry pi")
	portFlag = flag.Int("p", port, "The port the server is listening on")
	databaseFlag = flag.String("db", database, "The database file")
	apiKeyFlag = flag.String("key", "", "The google maps API key (with distance matrix api enabled)")
	email = flag.String("email", "", "Specify email to create account")
	password = flag.String("password", "", "Specify email to create account")
)

func main() {
	flag.Parse()
	log.SetLevel(log.InfoLevel)
	logConfiguration()
	db, err := sql.Open("sqlite3", *databaseFlag)
	if err != nil {
		log.Fatalf("Failed to open db [%s] errors [%s]", *databaseFlag, err)
	}

	defer db.Close()

	setupTables(*db)

	userDao := UserDao{*db}
	createUser(userDao)
	u := UserResource{userDao}

	e := EventResource{userDao:userDao, doorController:DoorController{relayPin: *relayPinFlag, contactSwitchPin: *contactSwitchPinFlag}, distanceUtil:DistanceCalculator{apiKey:*apiKeyFlag}}
	authFilter := AuthFilter{userDao: userDao}
	wsContainer := restful.NewContainer()
	u.register(wsContainer)
	e.register(wsContainer)

	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-Auth-Token"},
		AllowedHeaders: []string{"Content-Type", "Accept", "X-Auth-Token"},
		CookiesAllowed: false,
		Container:      wsContainer}
	wsContainer.Filter(cors.Filter)
	wsContainer.Filter(authFilter.tokenFilter)

	server := &http.Server{Addr: ":"+strconv.Itoa(*portFlag), Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}

func logConfiguration() {
	log.Debugf("Relay pin %d", *relayPinFlag)
	log.Debugf("Api key %s", *apiKeyFlag)
	log.Debugf("Database file %s", *databaseFlag)
	log.Debugf("Webserver port %d", *portFlag)
}

func setupTables(db sql.DB) {
	// Create user table
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS user (email TEXT NOT NULL PRIMARY KEY, password TEXT, token TEXT, longitude REAL, latitude REAL, last_open DATETIME, approved BOOLEAN);")
	if (err != nil) {log.Fatalf("%v", err)}
}

func createUser(userDao UserDao) {
	if ((*email != "" && password != nil) && (email != nil && *password != "")) {
		userDao.createUser(User{Email: *email, Password: *password, Approved: true})
		log.Infof("Created account email:%s. Exiting...", *email)
		os.Exit(0)
	}
}