package main
import (
	"time"
	log "github.com/Sirupsen/logrus"
)

func timeToOpen(times []TimeWindow) bool {
	now := time.Now()
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
func hasFiredOpenEvent(user User, times []TimeWindow) bool {
	lastOpen := user.LastOpen
	if (&lastOpen == nil) { return true }
	for _,timeWindow := range times {
		openStartTime := startTime(timeWindow)
		openEndTime := endTime(timeWindow)
		log.Debugf("Open start time [%s], open end time [%s], last_open time [%s]", openStartTime, openEndTime, user.LastOpen)
		if (lastOpen.After(openStartTime) && lastOpen.Before(openEndTime)) {
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