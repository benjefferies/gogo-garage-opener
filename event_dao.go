package main
import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
	"time"
)

type EventDao struct {
	db sql.DB
}

func (e EventDao) createEvent(event Event) {
	log.Debugf("inserting event, timestamp:[%s], open:[%s]", event.EventTime, event.Open)
	tx,_ := e.db.Begin()
	prepStmt,err := e.db.Prepare("insert into event values (?, ?)")
	defer prepStmt.Close()
	_,err = prepStmt.Exec(event.EventTime, event.Open)
	if (err != nil) {
		log.Error(err)
		tx.Rollback()
	} else {
		tx.Commit()
	}
}

func (e EventDao) getEvents() []Event {
	tx,_ := e.db.Begin()
	rows,err := e.db.Query("select timestamp, open from event order by timestamp desc")
	defer rows.Close()
	if (err != nil) {
		log.Error(err)
	}
	var events []Event;
	if (!rows.Next()) {
		return events;
	}
	for rows.Next() {
		var timestamp time.Time
		var open bool
		rows.Scan(&timestamp, &open)
		events = append(events, Event{timestamp, open})
	}
	tx.Commit()
	return events
}