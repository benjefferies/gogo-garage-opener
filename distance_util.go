package main
import (
	"strconv"
	log "github.com/Sirupsen/logrus"
	"math"
)


func parseFloat64(floatValue string) float64 {
	float, err := strconv.ParseFloat(floatValue, 64)
	if (err != nil) {
		log.Error(err)
	}
	return float
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

func withinRange(user User, lat, lon float64) bool {
	distanceFrom := distance(user.Latitude, user.Longitude, lat, lon)
	log.Debugf("%s is %sm from garage", user.Email, distanceFrom)

	return distanceFrom <= float64(user.Distance)
}