package main
import (
	"github.com/emicklei/go-restful"
	"net/http"
	 log "github.com/Sirupsen/logrus"
	"time"
)

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