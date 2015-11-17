package main
import (
	"github.com/emicklei/go-restful"
	"net/http"
	"log"
)

func main() {
	wsContainer := restful.NewContainer()
	ws := new(restful.WebService)

	ws.
	Path("/helloword").
	Consumes(restful.MIME_JSON, restful.MIME_XML).
	Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("").To(printHelloWorld))  // u is a UserResource

	wsContainer.Add(ws)

	server := &http.Server{Addr: ":8080", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}

func  printHelloWorld(request *restful.Request, response *restful.Response) {
	response.WriteEntity("Hello world")
}