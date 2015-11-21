package main

import (
	"github.com/emicklei/go-restful"
	"time"
	"github.com/golang/geo/s2"
	"math"
)

type Event struct {
	EventTime time.Time
	Open bool
}

type EventResource struct {
	eventDao EventDao
	userDao UserDao
}

func (e EventResource) Register (container *restful.Container) {
	ws := new(restful.WebService)

	ws.Path("/event").
	Consumes(restful.MIME_JSON, restful.MIME_JSON).
	Produces(restful.MIME_JSON, restful.MIME_JSON)

	ws.Route(ws.POST("{email}/{latitude}/{longitude}").To(e.garageEvent))
	ws.Route(ws.GET("").To(e.findEvents))

	container.Add(ws)
}

func (e EventResource) garageEvent(request *restful.Request, response *restful.Response) {
	user := e.userDao.getUser(request.PathParameter("email"))
	latitude := request.PathParameter("latitude")
	longitude := request.PathParameter("longitude")
	currentLatLong := s2.LatLng{latitude, longitude}
	garageLatLong := s2.LatLng{user.Latitude, user.Longitude}

	distanceInMetres := distance(currentLatLong.Lat, currentLatLong.Lng, garageLatLong.Lat, garageLatLong.Lng)

	if (distanceInMetres > 100) {
		return
	}
	
	events := e.eventDao.getEvents();

	var open bool
	if (len(events) == 0) {
		open = true
	} else {
		open = !events[0].Open;
	}

	event := Event{time.Now(), open}
	e.eventDao.createEvent(event)
}

// haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// Distance function returns the distance (in meters) between two points of
//     a given longitude and latitude relatively accurately (using a spherical
//     approximation of the Earth) through the Haversin Distance Formula for
//     great arc distance on a sphere with accuracy for small distances
//
// point coordinates are supplied in degrees and converted into rad. in the func
//
// distance returned is METERS!!!!!!
// http://en.wikipedia.org/wiki/Haversine_formula
func distance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}

func (e EventResource) findEvents(request *restful.Request, response *restful.Response) {
	response.WriteEntity(e.eventDao.getEvents());
}