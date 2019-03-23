package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/strategy"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gopkg.in/resty.v0"
)

var accessToken string

// NOTE Requires password grant flow to be enabled
func getAccessToken() string {
	url := "https://gogo-garage-opener.eu.auth0.com/oauth/token"
	email := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	log.WithField("clientSecret", len(clientSecret)).WithField("email", len(email)).WithField("password", len(password)).WithField("clientID", len(clientID)).Info("Making request to login")

	payloadString := "{\"grant_type\":\"http://auth0.com/oauth/grant-type/password-realm\",\"username\": \"" + email + "\",\"password\": \"" + password + "\",\"audience\": \"https://open.mygaragedoor.space/api\", \"scope\": \"openid email\", \"client_id\": \"" + clientID + "\", \"client_secret\": \"" + clientSecret + "\", \"realm\": \"Username-Password-Authentication\"}"
	payload := strings.NewReader(payloadString)
	log.WithField("payload", payload).Info("Making request to login")
	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)
	var accessToken map[string]interface{}
	body, err := ioutil.ReadAll(res.Body)
	log.WithField("body", string(body)).Info("Response for userinfo")
	if err != nil {
		panic(errors.New("Could not parse user info"))
	}
	err = json.Unmarshal(body, &accessToken)
	if err != nil {
		panic(errors.New("Could not marsher user info"))
	}

	defer res.Body.Close()
	json := accessToken["access_token"]
	return fmt.Sprint(json)
}

func TestMain(m *testing.M) {
	accessToken = getAccessToken()
	log.Info("Starting server")
	go main()
	authHeader := fmt.Sprintf("Bearer %s", accessToken)
	err := retry.Retry(func(attempt uint) error {
		_, err := resty.R().SetHeader("Authorization", authHeader).Get("http://localhost:8080/user/one-time-pin/my-pin")
		return err
	}, strategy.Limit(5), strategy.Delay(time.Second))
	if err != nil {
		log.WithError(err).Fatal("Application is not initialised")
	}
	exitCode := m.Run()
	log.Info("Started server in integration test")
	err = os.Remove("gogo-garage-opener.db")
	if err != nil {
		log.WithError(err).Fatal("Could not delete database file")
	}
	os.Exit(exitCode)
}

func TestOneTimePinAccess(t *testing.T) {
	// Given
	authHeader := fmt.Sprintf("Bearer %s", accessToken)

	// When
	response, err := resty.R().
		SetHeader("Authorization", authHeader).
		Get("http://localhost:8080/user/one-time-pin/my-pin")

	// Then
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode(), "Expecting OK http status")
	assert.Contains(t, string(response.Body()), "action=\"/garage/one-time-pin/my-pin\"", "Should contain link to use pin")
}

func TestNewOneTimePin(t *testing.T) {
	// Given
	authHeader := fmt.Sprintf("Bearer %s", accessToken)

	// When
	response, err := resty.R().
		SetHeader("Authorization", authHeader).
		SetHeader("Content-Type", "application/json").
		Post("http://localhost:8080/user/one-time-pin")

	// Then
	assert.Nil(t, err, "Not expecting an error")
	assert.Equal(t, 200, response.StatusCode(), "Expecting OK http status")
	assert.Contains(t, string(response.Body()), getPin(t), "Response should contain pin")
}

func TestGetOneTimePins(t *testing.T) {
	// Given
	authHeader := fmt.Sprintf("Bearer %s", accessToken)

	// When
	response, err := resty.R().
		SetHeader("Authorization", authHeader).
		SetHeader("Content-Type", "application/json").
		Get("http://localhost:8080/user/one-time-pin")

	// Then
	var pins []Pin
	json.Unmarshal(response.Body(), &pins)
	assert.Nil(t, err, "Not expecting an error")
	assert.Equal(t, 200, response.StatusCode(), "Expecting OK http status")
	assert.Equal(t, pins[0].Pin, getPin(t), "Response should contain pin")
	assert.Equal(t, pins[0].CreatedBy, "test@echosoft.uk", "Response should who created the pin")
}

func TestUseOneTimePin(t *testing.T) {
	// Given
	pin := getNewPin(t)

	// When
	response, err := resty.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Accept", "application/json").
		Post("http://localhost:8080/garage/one-time-pin/" + pin)

	// Then
	assert.Nil(t, err, "Not expecting an error")
	assert.Equal(t, 202, response.StatusCode(), "Expecting accepted http status")
}

func TestCannotUseOneTimePinTwice(t *testing.T) {
	// Given
	pin := getNewPin(t)
	response, err := resty.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Accept", "application/json").
		Post("http://localhost:8080/garage/one-time-pin/" + pin)

	// When
	response, err = resty.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Accept", "application/json").
		Post("http://localhost:8080/garage/one-time-pin/" + pin)

	// Then
	assert.Nil(t, err, "Not expecting an error")
	assert.Equal(t, 401, response.StatusCode(), "Should not be authorised")
}

func TestToggleGarage(t *testing.T) {
	// Given
	authHeader := fmt.Sprintf("Bearer %s", accessToken)

	// When
	response, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", authHeader).
		Post("http://localhost:8080/garage/toggle")

	// Then
	assert.Nil(t, err, "Not expecting an error")
	assert.Equal(t, 202, response.StatusCode(), "Expecting accepted http status")
}

func TestGarageStatus(t *testing.T) {
	// Given
	authHeader := fmt.Sprintf("Bearer %s", accessToken)

	// When
	response, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", authHeader).
		SetResult(map[string]interface{}{}).
		Get("http://localhost:8080/garage/state")

	// Then
	assert.Nil(t, err, "Not expecting an error")
	assert.Equal(t, 200, response.StatusCode(), "Expecting OK http status")
	result := (*response.Result().(*map[string]interface{}))
	assert.Equal(t, fmt.Sprintf("%v", closed), fmt.Sprintf("%.f", result["State"]), "Expecting closed status")
	assert.Equal(t, "Closed", result["Description"], "Expecting closed description")
}

func getPin(t *testing.T) string {
	db, err := sql.Open("sqlite3", *database)
	assert.Nil(t, err, "Not expecting an error")
	db.Begin()
	row := db.QueryRow("select pin from one_time_pin where created_by = ?", "test@echosoft.uk")
	assert.NotEqual(t, row, sql.ErrNoRows, "Shouldn't be error")
	db.Close()
	var pin string
	row.Scan(&pin)
	log.WithField("pin", pin).Info("Found pin")
	return pin
}

func getNewPin(t *testing.T) string {
	authHeader := fmt.Sprintf("Bearer %s", accessToken)
	response, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", authHeader).
		Post("http://localhost:8080/user/one-time-pin")

	t.Log(accessToken)
	assert.Nil(t, err, "Not expecting an error")
	var pin map[string]string
	err = json.Unmarshal(response.Body(), &pin)
	assert.Nil(t, err, "Not expecting an error")
	return pin["pin"]
}
