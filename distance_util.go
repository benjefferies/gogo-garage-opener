package main
import (
	"strconv"
	log "github.com/Sirupsen/logrus"
	"googlemaps.github.io/maps"
	"golang.org/x/net/context"
	"time"
	"fmt"
)

type DistanceUtil struct {
	apiKey string
}

func parseFloat64(floatValue string) float64 {
	float, err := strconv.ParseFloat(floatValue, 64)
	if (err != nil) {
		log.Error(err)
	}
	return float
}

func (du DistanceUtil) getArrivalTime(user User, destLat, destLon float64) time.Duration {
	garageLat, garageLon := user.Latitude, user.Longitude
	c, err := maps.NewClient(maps.WithAPIKey(du.apiKey))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	r := &maps.DistanceMatrixRequest{
		Origins:      []string{fmt.Sprintf("%f,%f", garageLat, garageLon)},
		Destinations: []string{fmt.Sprintf("%f,%f", destLat, destLon)},
		Mode: maps.TravelModeDriving,

	}
	log.Debugf("Distance matrix with origins [%v] and destination [%v]", r.Origins, r.Destinations)
	resp, err := c.DistanceMatrix(context.Background(), r)
	if err != nil {
		log.Warnf("Could not find distance: %s", err)
	}
	for _,row := range resp.Rows {
		for _, element := range row.Elements {
			log.Debugf("User location [%s],[%s] to [%s],[%v] is %s away", user.Latitude, user.Longitude, destLat, destLon, element.Duration)
			return element.Duration
		}
	}

	return *new(time.Duration)
}
