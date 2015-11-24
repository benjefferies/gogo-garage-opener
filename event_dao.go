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
	log.Debugf("inserting event, timestamp:[%s], email:[%s]", event.EventTime, event.Email)
	tx,_ := e.db.Begin()
	prepStmt,err := e.db.Prepare("insert into event values (?, ?)")
	defer prepStmt.Close()
	_,err = prepStmt.Exec(event.EventTime, event.Email)
	if (err != nil) {
		log.Error(err)
		tx.Rollback()
	} else {
		tx.Commit()
	}
}

func (e EventDao) getEvents() []Event {
	tx,_ := e.db.Begin()
	rows,err := e.db.Query("select timestamp, email from event order by timestamp desc")
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
		var email string
		rows.Scan(&timestamp, &email)
		events = append(events, Event{timestamp, email})
	}
	tx.Commit()
	return events
}