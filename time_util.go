package main
import (
	"time"
	log "github.com/Sirupsen/logrus"
)

func timeToOpen(user User) bool {
	hour, minute, second := time.Now().Clock()
	now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hour, minute, second, 0, time.Now().Location())
	openStartTime := startTime(user)
	openEndTime := endTime(user)

	log.Debugf("Open start time [%s], open end time [%s], now [%s]", openStartTime, openEndTime, now)
	return now.After(openStartTime) && now.Before(openEndTime)
}

func startTime(user User) time.Time {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), user.Time.Hour(), user.Time.Minute(), user.Time.Second(), 0, time.Now().Location())
}

func endTime(user User) time.Time {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), user.Time.Hour(), user.Time.Minute(), user.Time.Second(), 0, time.Now().Location()).Add(time.Duration(user.Duration) * time.Minute)
}