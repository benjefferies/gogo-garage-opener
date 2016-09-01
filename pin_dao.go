package main

import (
	log "github.com/Sirupsen/logrus"
	"database/sql"
	"github.com/satori/go.uuid"
	"time"
	"errors"
)

type PinDao struct {
	db sql.DB
}

func (p PinDao) newOneTimePin(user User) (uuid.UUID, error) {
	pin := uuid.NewV4()
	now := time.Now()
	log.Debugf("inserting one time pin:[%s], created_by:[%s], created:[%s]", pin.String(), user.Email, now.Local())

	tx, _ := p.db.Begin()
	prepStmt, err := p.db.Prepare("insert into one_time_pin(pin, created_by, created) values (?, ?, ?)")
	defer prepStmt.Close()
	_, err = prepStmt.Exec(pin, user.getEmail(), now.Unix())
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return pin, err
}

func (p PinDao) use(pin uuid.UUID) error {
	now := time.Now()
	log.Debugf("using one time pin: [%s] at: [%s]", pin.String(), now.Local())

	tx, _ := p.db.Begin()
	prepStmt, err := p.db.Prepare("update one_time_pin set used = ? where pin = ? and used is null")
	defer prepStmt.Close()
	resp, err := prepStmt.Exec(now.Unix(), pin.String())
	if err != nil {
		tx.Rollback()
		return err
	} else {
		tx.Commit()
	}
	rowsAffected, err := resp.RowsAffected()
	if err == nil && rowsAffected == 0 {
		err = errors.New("Pin has already been used")
	}
	return err
}

func (p PinDao) getPinUsedDate(pin uuid.UUID) (int64, error) {
	log.Debugf("Getting used date for pin: [%s]", pin.String())

	tx, _ := p.db.Begin()
	prepStmt, err := p.db.Prepare("select used from one_time_pin where pin = ?")
	defer prepStmt.Close()
	row := prepStmt.QueryRow(pin.String())
	var usedDate int64
	row.Scan(&usedDate)
	if err != nil {
		log.WithError(err).Error("Could not get pin used date for [%s]", pin.String())
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return usedDate, err
}