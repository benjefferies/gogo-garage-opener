package main

import (
	"database/sql"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// GarageDoorDao for garage configuration
type GarageDoorDao interface {
	updateConfiguration(updateConfig []GarageConfiguration) error
	getConfiguration() ([]GarageConfiguration, error)
}

// SqliteGarageDoorDao for garage configuration
type SqliteGarageDoorDao struct {
	db *sql.DB
}

// GarageConfiguration configuration for the garage
type GarageConfiguration struct {
	Day             *string    `json:"day"`
	OpenDuration    *int64     `json:"openDuration"`
	ShouldCloseTime *time.Time `json:"shouldCloseTime,omitempty"`
	CanStayOpenTime *time.Time `json:"canStayOpenTime,omitempty"`
	Enabled         *bool      `json:"enabled"`
}

func (garageConfiguration GarageConfiguration) String() string {
	return fmt.Sprintf("day=%s openDuration=%d shouldCloseTime=%v canStayOpenTime=%v enabled=%t",
		*garageConfiguration.Day, *garageConfiguration.OpenDuration, garageConfiguration.ShouldCloseTime,
		garageConfiguration.CanStayOpenTime, *garageConfiguration.Enabled)
}

func (garageDoorDao SqliteGarageDoorDao) updateConfiguration(updateConfig []GarageConfiguration) error {
	log.Debug("Updating garage configuration")
	tx, _ := garageDoorDao.db.Begin()
	for _, config := range updateConfig {
		if config.Day == nil {
			continue
		}
		if config.OpenDuration != nil {
			err := garageDoorDao.updateValue("open_duration", config.OpenDuration, *config.Day, tx)
			if err != nil {
				return err
			}
		}
		if config.ShouldCloseTime != nil {
			err := garageDoorDao.updateValue("should_close_time", config.ShouldCloseTime.Format(time.RFC3339), *config.Day, tx)
			if err != nil {
				return err
			}
		}
		if config.CanStayOpenTime != nil {
			err := garageDoorDao.updateValue("can_stay_open_time", config.CanStayOpenTime.Format(time.RFC3339), *config.Day, tx)
			if err != nil {
				return err
			}
		}
		if config.Enabled != nil {
			err := garageDoorDao.updateValue("enabled", config.Enabled, *config.Day, tx)
			if err != nil {
				return err
			}
		}
	}
	tx.Commit()
	return nil
}

func (garageDoorDao SqliteGarageDoorDao) updateValue(field string, value interface{}, day string, tx *sql.Tx) error {
	log.WithField(field, value).Info("Updating value")
	prepStmt, err := garageDoorDao.db.Prepare("update garage set " + field + " = ? where day = ?")
	defer prepStmt.Close()
	if err != nil {
		log.WithError(err).Errorf("Failed to update %s", day)
		tx.Rollback()
		return err
	}
	_, err = prepStmt.Exec(value, day)
	if err != nil {
		log.WithError(err).Errorf("Failed to update %s", day)
		tx.Rollback()
		return err
	}
	return nil
}

func (garageDoorDao SqliteGarageDoorDao) getConfiguration() ([]GarageConfiguration, error) {
	log.Debug("Getting garage configuration")
	tx, _ := garageDoorDao.db.Begin()
	prepStmt, err := garageDoorDao.db.Prepare("select day, open_duration, should_close_time, can_stay_open_time, enabled from garage")
	defer prepStmt.Close()
	rows, err := prepStmt.Query()
	var config []GarageConfiguration
	for rows.Next() {
		var dayConfig GarageConfiguration
		var shouldCloseTime string
		var canStayOpenTime string
		rows.Scan(&dayConfig.Day, &dayConfig.OpenDuration, &shouldCloseTime, &canStayOpenTime, &dayConfig.Enabled)
		parsedShouldCloseTime, _ := time.Parse(time.RFC3339, shouldCloseTime)
		parsedCanStayOpenTime, _ := time.Parse(time.RFC3339, canStayOpenTime)
		dayConfig.ShouldCloseTime = &parsedShouldCloseTime
		dayConfig.CanStayOpenTime = &parsedCanStayOpenTime
		config = append(config, dayConfig)
	}
	if err != nil {
		log.Error("Could not get garage config")
		tx.Rollback()
	}
	tx.Commit()
	log.WithField("config", config).Debug("Found garage config")
	return config, err
}

func (garageDoorDao SqliteGarageDoorDao) init() {
	log.Debug("Initialising garage configuration")

	tx, _ := garageDoorDao.db.Begin()
	prepStmt, err := garageDoorDao.db.Prepare("insert into garage(day, should_close_time, can_stay_open_time) SELECT ?, ?, ? WHERE NOT EXISTS(SELECT 1 FROM garage WHERE day = ?);")
	if err != nil {
		log.WithError(err).Error("Could prepare statement to create garage configuration")
	}
	defer prepStmt.Close()
	createGarageConfiguration(tx, prepStmt, time.Monday.String())
	createGarageConfiguration(tx, prepStmt, time.Tuesday.String())
	createGarageConfiguration(tx, prepStmt, time.Wednesday.String())
	createGarageConfiguration(tx, prepStmt, time.Thursday.String())
	createGarageConfiguration(tx, prepStmt, time.Friday.String())
	createGarageConfiguration(tx, prepStmt, time.Saturday.String())
	createGarageConfiguration(tx, prepStmt, time.Sunday.String())
	tx.Commit()
}

func createGarageConfiguration(tx *sql.Tx, prepStmt *sql.Stmt, day string) {
	now := time.Now()
	shouldCloseTime := time.Date(now.Year(), now.Month(), now.Day(), 22, 0, 0, 0, time.Local)
	canStayOpenTime := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, time.Local)
	_, err := prepStmt.Exec(day, shouldCloseTime.Format(time.RFC3339), canStayOpenTime.Format(time.RFC3339), day)
	if err != nil {
		tx.Rollback()
		log.WithError(err).WithField("day", day).Fatal("Could not create garage configuration")
	}
}
