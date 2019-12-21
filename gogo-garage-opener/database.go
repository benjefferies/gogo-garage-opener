package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

func initialise(databasePath string) *sql.DB {
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.WithError(err).WithField("database", databasePath).Fatalf("Failed to open db")
	}

	setupTables(db)
	return db
}

func setupTables(db *sql.DB) {
	// Create user table
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS user (email TEXT NOT NULL PRIMARY KEY, token TEXT, subscribed BOOLEAN DEFAULT 1);")
	if err != nil {
		log.WithError(err).Fatal("Could not create user table")
	}

	// Create one time pin table
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS one_time_pin (pin TEXT NOT NULL PRIMARY KEY, created_by TEXT, created INTEGER, used INTEGER);")
	if err != nil {
		log.WithError(err).Fatal("Could not create one_time_pin table")
	}

	// Create garage configuration table
	// day - ISO 8601 day of week i.e. Sunday, Monday
	// open_duration - how long the garage door can be left open in seconds during close time
	// should_close_time - ISO 8601 time of day that the garage door should start auto closing
	// can_stay_open_time - ISO 8601 time of day that the garage door can stay open and should not autoclose
	// enabled - whether autoclose is enabled or not
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS garage (day TEXT NOT NULL PRIMARY KEY, open_duration INTEGER DEFAULT 180, should_close_time TEXT, can_stay_open_time TEXT, enabled BOOLEAN DEFAULT 1);")
	if err != nil {
		log.WithError(err).Fatal("Could not create garage table")
	}

}
