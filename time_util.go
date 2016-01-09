package main
import (
	"time"
	log "github.com/Sirupsen/logrus"
)

func timeToOpen(times []TimeWindow) bool {
	hour, minute, second := time.Now().Clock()
	now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hour, minute, second, 0, time.Now().Location())
	for _,timeWindow := range times {
		openStartTime := startTime(timeWindow)
		openEndTime := endTime(timeWindow)
		log.Debugf("Open start time [%s], open end time [%s], now [%s]", openStartTime, openEndTime, now)
		if (now.After(openStartTime) && now.Before(openEndTime)) {
			return true;
		}
	}
	return false;
}

func toTime(timeWindow TimeWindow) time.Time {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), timeWindow.Time.Hour(), timeWindow.Time.Minute(), timeWindow.Time.Second(), 0, time.Now().Location())
}

func startTime(timeWindow TimeWindow) time.Time {
	return toTime(timeWindow);
}

func endTime(timeWindow TimeWindow) time.Time {
	return toTime(timeWindow).Add(time.Duration(timeWindow.DurationWindow) * time.Minute)
}