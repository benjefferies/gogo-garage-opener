package main

import (
	"database/sql"
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/ventu-io/go-shortid"
)

// PinDao is the data access for one time pins
type PinDao struct {
	db *sql.DB
}

func (pinDao PinDao) newOneTimePin(email string) (string, error) {
	pin := shortid.MustGenerate()
	now := time.Now()
	log.Debugf("inserting one time pin:[%s], created_by:[%s], created:[%s]", pin, email, now.Local())

	tx, _ := pinDao.db.Begin()
	prepStmt, err := pinDao.db.Prepare("insert into one_time_pin(pin, created_by, created) values (?, ?, ?)")
	defer prepStmt.Close()
	_, err = prepStmt.Exec(pin, email, now.Unix())
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return pin, err
}

func (pinDao PinDao) use(pin string) error {
	now := time.Now()
	log.Debugf("using one time pin: [%s] at: [%s]", pin, now.Local())

	tx, _ := pinDao.db.Begin()
	prepStmt, err := pinDao.db.Prepare("update one_time_pin set used = ? where pin = ? and used is null")
	defer prepStmt.Close()
	resp, err := prepStmt.Exec(now.Unix(), pin)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	rowsAffected, err := resp.RowsAffected()
	if err == nil && rowsAffected == 0 {
		err = errors.New("Pin has already been used")
	}
	return err
}

func (pinDao PinDao) getPinUsedDate(pin string) (int64, error) {
	log.Debugf("Getting used date for pin: [%s]", pin)

	tx, _ := pinDao.db.Begin()
	prepStmt, err := pinDao.db.Prepare("select used from one_time_pin where pin = ?")
	defer prepStmt.Close()
	row := prepStmt.QueryRow(pin)
	var usedDate int64
	row.Scan(&usedDate)
	if err != nil {
		log.WithError(err).Errorf("Could not get pin used date for [%s]", pin)
		tx.Rollback()
	}
	tx.Commit()
	return usedDate, err
}