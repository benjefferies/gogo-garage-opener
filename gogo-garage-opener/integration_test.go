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
	resty "github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var accessToken string

var client *resty.Client = resty.New()

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
		_, err := client.R().SetHeader("Authorization", authHeader).Get("http://localhost:8080/user/one-time-pin/my-pin")
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
	response, err := client.R().
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
	var pin map[string]string

	// When
	response, err := client.R().
		SetHeader("Authorization", authHeader).
		SetHeader("Content-Type", "application/json").
		SetResult(&pin).
		Post("http://localhost:8080/user/one-time-pin")

	// Then
	assert.Nil(t, err, "Not expecting an error")
	log.WithField("pin", pin).Info("Got new pin")
	assert.Equal(t, 200, response.StatusCode(), "Expecting OK http status")
	assert.NotEmpty(t, pin["pin"], "Response should contain pin")
	assert.True(t, pinExists(t, pin["pin"]), "Pin should exist")
}

func TestGetOneTimePins(t *testing.T) {
	// Given
	authHeader := fmt.Sprintf("Bearer %s", accessToken)

	// When
	response, err := client.R().
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

func TestDeleteOneTimePins(t *testing.T) {
	// Given
	pin := getNewPin(t)
	authHeader := fmt.Sprintf("Bearer %s", accessToken)

	// When
	response, err := client.R().
		SetHeader("Authorization", authHeader).
		SetHeader("Content-Type", "application/json").
		Delete("http://localhost:8080/user/one-time-pin/" + pin)

	// Then
	assert.Nil(t, err, "Not expecting an error")
	assert.Equal(t, 200, response.StatusCode(), "Expecting OK http status")
	assert.False(t, pinExists(t, pin), "Response should who created the pin")
}

func TestUseOneTimePin(t *testing.T) {
	// Given
	pin := getNewPin(t)

	// When
	response, err := client.R().
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
	response, err := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Accept", "application/json").
		Post("http://localhost:8080/garage/one-time-pin/" + pin)

	// When
	response, err = client.R().
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
	response, err := client.R().
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
	response, err := client.R().
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

func TestDefaultGarageConfiguration(t *testing.T) {
	// Given
	authHeader := fmt.Sprintf("Bearer %s", accessToken)

	// When
	response, err := client.R().
		SetHeader("Authorization", authHeader).
		Get("http://localhost:8080/garage/config")

	// Then
	assert.Equal(t, 200, response.StatusCode(), "Expecting OK http status")
	assert.Nil(t, err, "Not expecting an error")
	var config []GarageConfiguration
	err = json.Unmarshal(response.Body(), &config)
	assert.Nil(t, err, "Not expecting an error")
	assert.Equal(t, 7, len(config), "Expecting configuration for all days of week")
	dayConfig := getConfigByDay("Sunday", config)
	assert.Equal(t, "Sunday", *dayConfig.Day, "Expecting day to be Sunday")
	assert.Equal(t, int64(180), *dayConfig.OpenDuration, "Expecting open duration to be 3")
	now := time.Now()
	shouldCloseTime := time.Date(now.Year(), now.Month(), now.Day(), 22, 0, 0, 0, time.UTC)
	assert.Equal(t, shouldCloseTime, *dayConfig.ShouldCloseTime, "Expecting ShouldCloseTime to be 10pm")
	canStayOpenTime := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, time.UTC)
	assert.Equal(t, canStayOpenTime, *dayConfig.CanStayOpenTime, "Expecting CanStayOpenTime to be 8am")
	assert.Equal(t, true, *dayConfig.Enabled, "Expecting enabled to be true")
}

func TestUpdateGarageConfiguration(t *testing.T) {
	// Given
	authHeader := fmt.Sprintf("Bearer %s", accessToken)
	now := time.Now()
	var day = "Sunday"
	var duration = int64(4)
	var enabled = false
	newConfig := [1]GarageConfiguration{GarageConfiguration{Day: &day, OpenDuration: &duration, ShouldCloseTime: &now, CanStayOpenTime: &now, Enabled: &enabled}}
	b, err := json.Marshal(newConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
	response, err := client.R().
		SetHeader("Authorization", authHeader).
		SetHeader("Content-Type", "application/json").
		SetBody(string(b)).
		Put("http://localhost:8080/garage/config")
	assert.Nil(t, err, "Not expecting an error")

	// When
	response, err = client.R().
		SetHeader("Authorization", authHeader).
		Get("http://localhost:8080/garage/config")

	// Then
	assert.Equal(t, 200, response.StatusCode(), "Expecting OK http status")
	assert.Nil(t, err, "Not expecting an error")
	var config []GarageConfiguration
	err = json.Unmarshal(response.Body(), &config)
	assert.Nil(t, err, "Not expecting an error")
	assert.Equal(t, 7, len(config), "Expecting configuration for all days of week")
	dayConfig := getConfigByDay("Sunday", config)
	assert.Equal(t, day, *dayConfig.Day, "Expecting day to be Sunday")
	assert.Equal(t, duration, *dayConfig.OpenDuration, "Expecting open duration to be 4")
	assert.Equal(t, now.Format(time.RFC3339), dayConfig.ShouldCloseTime.Format(time.RFC3339), "Expecting ShouldCloseTime to be now")
	assert.Equal(t, now.Format(time.RFC3339), dayConfig.CanStayOpenTime.Format(time.RFC3339), "Expecting CanStayOpenTime to be now")
	assert.Equal(t, enabled, *dayConfig.Enabled, "Expecting enabled to be false")
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

func pinExists(t *testing.T, pin string) bool {
	db, err := sql.Open("sqlite3", *database)
	assert.Nil(t, err, "Not expecting an error")
	db.Begin()
	row := db.QueryRow("select count(pin) from one_time_pin where pin = ?", pin)
	assert.NotEqual(t, row, sql.ErrNoRows, "Shouldn't be error")
	db.Close()
	var count int64
	row.Scan(&count)
	log.WithField("pin", pin).WithField("count", count).Info("Count pin")
	return count > 0
}

func getNewPin(t *testing.T) string {
	authHeader := fmt.Sprintf("Bearer %s", accessToken)
	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", authHeader).
		Post("http://localhost:8080/user/one-time-pin")

	assert.Nil(t, err, "Not expecting an error")
	var pin map[string]string
	err = json.Unmarshal(response.Body(), &pin)
	assert.Nil(t, err, "Not expecting an error")
	return pin["pin"]
}
