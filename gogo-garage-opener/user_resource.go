package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// oneTimePinUse static html of one time pin page
const oneTimePinUse = `
<html>
	</head>
	<body>
		<h1>One time pin garage door opener</h1>
		<br>
		<p>By clicking the button below it will open the garage door. The garage door will automatically close %v seconds after clicking the button</p>
		<br>
		<br>
		<form name="myform" action="/garage/one-time-pin/%s" method="post">
			<button>Open</button>
		</form>
	</body>
</html>
`

// UserResource API for users
type UserResource struct {
	userDao UserDao
	pinDao  PinDao
}

func (userResource UserResource) register(router *mux.Router) {
	subRouter := router.PathPrefix("/user").Subrouter()
	subRouter.Path("/login").Methods("POST").Handler(jwtCheckHandleFunc(userResource.login))
	subRouter.Path("/one-time-pin").Methods("POST").Handler(jwtCheckHandleFunc(userResource.oneTimePin))
	subRouter.Path("/one-time-pin/{oneTimePin}").Methods("GET").HandlerFunc(userResource.useOneTimePin)
}

func (userResource UserResource) oneTimePin(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Creating new one time pin")
	accessToken := context.Get(r, "access_token")
	email := getEmail(fmt.Sprintf("%s", accessToken))
	log.Debugf("%s is creating new one time pin", email)
	pin, err := userResource.pinDao.newOneTimePin(email)
	if err != nil {
		log.WithError(err).Error("Could not create one time pin")
		w.WriteHeader(500)
	}
	pinMap := map[string]interface{}{
		"pin": pin,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pinMap)
	w.WriteHeader(http.StatusOK)
}

func (userResource UserResource) useOneTimePin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	oneTimePin := vars["oneTimePin"]
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, oneTimePinUse, timeToClose.Seconds(), oneTimePin)
}

func (userResource UserResource) login(w http.ResponseWriter, r *http.Request) {
	accessToken := context.Get(r, "access_token")
	email := getEmail(fmt.Sprintf("%s", accessToken))
	userResource.userDao.createUser(email)
}
