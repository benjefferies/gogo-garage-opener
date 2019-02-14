package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"os"
	"testing"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/strategy"
	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gopkg.in/resty.v0"
)

func TestMain(m *testing.M) {
	flag.Set("email", "test@example.com")
	flag.Set("password", "password")
	main()
	flag.Set("email", "")
	flag.Set("password", "")
	log.Info("created user")
	flag.Set("noop", "true")
	log.Info("Starting server")
	go main()
	err := retry.Retry(func(attempt uint) error {
		_, err := resty.R().Get("http://localhost:8080/user/one-time-pin/my-pin")
		return err
	}, strategy.Limit(5), strategy.Delay(time.Second))
	if err != nil {
		log.WithError(err).Fatal("Application is not initialised")
	}
	m.Run()
	log.Info("Started server")
	err = os.Remove("gogo-garage-opener.db")
	if err != nil {
		log.WithError(err).Fatal("Could not delete database file")
	}
	os.Exit(0)
}

func TestOneTimePinAccess(t *testing.T) {
	response, err := resty.R().Get("http://localhost:8080/user/one-time-pin/my-pin")

	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode(), "Expecting OK http status")
	assert.Contains(t, string(response.Body()), "action=\"/garage/one-time-pin/my-pin\"", "Should contain link to use pin")
}

func TestNewOneTimePin(t *testing.T) {
	response, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Auth-Token", getToken(t)).
		Post("http://localhost:8080/user/one-time-pin")

	assert.Nil(t, err, "Not expecting an error")
	assert.Equal(t, 200, response.StatusCode(), "Expecting OK http status")
	assert.Contains(t, string(response.Body()), getPin(t), "Response should contain pin")
}

func TestUseOneTimePin(t *testing.T) {
	pin := getNewPin(t)

	response, err := resty.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Accept", "application/json").
		Post("http://localhost:8080/garage/one-time-pin/" + pin)

	assert.Nil(t, err, "Not expecting an error")
	assert.Equal(t, 202, response.StatusCode(), "Expecting accepted http status")
}

func TestCannotUseOneTimePinTwice(t *testing.T) {
	pin := getNewPin(t)
	response, err := resty.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Accept", "application/json").
		Post("http://localhost:8080/garage/one-time-pin/" + pin)

	response, err = resty.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Accept", "application/json").
		Post("http://localhost:8080/garage/one-time-pin/" + pin)

	assert.Nil(t, err, "Not expecting an error")
	assert.Equal(t, 401, response.StatusCode(), "Should not be authorised")
}

func TestToggleGarage(t *testing.T) {
	response, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Auth-Token", getToken(t)).
		Post("http://localhost:8080/garage/toggle")

	assert.Nil(t, err, "Not expecting an error")
	assert.Equal(t, 202, response.StatusCode(), "Expecting accepted http status")
}

func TestGarageStatus(t *testing.T) {
	response, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Auth-Token", getToken(t)).
		Get("http://localhost:8080/garage/state")

	assert.Nil(t, err, "Not expecting an error")
	assert.Equal(t, 200, response.StatusCode(), "Expecting OK http status")
}

func getToken(t *testing.T) string {
	db, err := sql.Open("sqlite3", *databaseFlag)
	assert.Nil(t, err, "Not expecting an error")
	db.Begin()
	row := db.QueryRow("select token from user where email = ?", "test@example.com")
	db.Close()
	var token string
	row.Scan(&token)
	return token
}

func getPin(t *testing.T) string {
	db, err := sql.Open("sqlite3", *databaseFlag)
	assert.Nil(t, err, "Not expecting an error")
	db.Begin()
	row := db.QueryRow("select pin from one_time_pin where email = ?", "test@example.com")
	db.Close()
	var pin string
	row.Scan(&pin)
	return pin
}

func getNewPin(t *testing.T) string {
	response, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Auth-Token", getToken(t)).
		Post("http://localhost:8080/user/one-time-pin")
	assert.Nil(t, err, "Not expecting an error")
	var pin map[string]string
	err = json.Unmarshal(response.Body(), &pin)
	assert.Nil(t, err, "Not expecting an error")
	return pin["pin"]
}
