package main
import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
	"strconv"
	"time"
)

type EventDao struct {
	db sql.DB
}

func (e EventDao) createEvent(event Event) {
	log.Debug("inserting event timestamp:[%s],open:[%s]",event.EventTime, event.Open)
	tx,_ := e.db.Begin()
	prepStmt,_ := e.db.Prepare("insert into event values (?, ?)")
	defer prepStmt.Close()
	_,err := prepStmt.Exec(event.EventTime, event.Open)
	if (err != nil) {
		log.Error(err)
		tx.Rollback()
	} else {
		tx.Commit()
	}
}

func (e EventDao) getEvents() []Event {
	log.Debug("getting all events")
	rows,err := e.db.Query("select timestamp, event from event")
	if (err != nil) {
		log.Error(err)
	}
	var events []Event;
	for rows.Next() {
		cols,err := rows.Columns()
		if (err != nil) {
			log.Error(err)
		}
		time,err := time.Parse(time.RFC3339Nano, cols[0])
		if (err != nil) {
			log.Error(err)
		}
		open,err := strconv.ParseBool(cols[1])
		if (err != nil) {
			log.Error(err)
		}
		events = append(events, Event{time, open})
	}
	return events
}