package main

import (
	"database/sql"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ventu-io/go-shortid"
)

// PinDao is the data access for one time pins
type PinDao struct {
	db *sql.DB
}

// Pin datamodel
type Pin struct {
	Pin       string    `json:"pin"  binding:"required"`
	CreatedBy string    `json:"createdBy"  binding:"required"`
	Created   time.Time `json:"created"  binding:"required"`
	Used      time.Time `json:"used"  binding:"required"`
}

func (pinDao PinDao) newOneTimePin(email string) (string, error) {
	pin := shortid.MustGenerate()
	now := time.Now()
	log.WithField("one_time_pin", pin).WithField("email", email).WithField("created", now.Local()).Debug("inserting one time pin")

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
	log.WithField("one_time_pin", pin).WithField("time", now.Local()).Debug("using one time pin")

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

func (pinDao PinDao) delete(pin string) error {
	log.WithField("one_time_pin", pin).Debug("Deleting one time pin")

	tx, _ := pinDao.db.Begin()
	prepStmt, err := pinDao.db.Prepare("delete from one_time_pin where pin = ?")
	defer prepStmt.Close()
	_, err = prepStmt.Exec(pin)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (pinDao PinDao) getPinUsedDate(pin string) (int64, error) {
	log.WithField("one_time_pin", pin).Debug("Getting used date for pin")

	tx, _ := pinDao.db.Begin()
	prepStmt, err := pinDao.db.Prepare("select used from one_time_pin where pin = ?")
	defer prepStmt.Close()
	row := prepStmt.QueryRow(pin)
	var usedDate int64
	row.Scan(&usedDate)
	if err != nil {
		log.WithError(err).WithField("one_time_pin", pin).Error("Could not get pin used date")
		tx.Rollback()
	}
	tx.Commit()
	return usedDate, err
}

func (pinDao PinDao) getPins() ([]Pin, error) {
	log.Debug("Getting pins")
	tx, _ := pinDao.db.Begin()
	prepStmt, err := pinDao.db.Prepare("select pin, created_by, created, used from one_time_pin")
	defer prepStmt.Close()
	rows, err := prepStmt.Query()
	var pins []Pin
	for rows.Next() {
		var pin Pin
		var created int64
		var used int64
		rows.Scan(&pin.Pin, &pin.CreatedBy, &created, &used)
		pin.Created = time.Unix(created, 0)
		pin.Used = time.Unix(used, 0)
		log.WithField("pin", pin).Debug("Found pin")
		pins = append(pins, pin)
	}
	if err != nil {
		log.Error("Could not get pins")
		tx.Rollback()
	}
	tx.Commit()
	return pins, err
}
