package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func getEmail(accessToken string) string {
	req, err := http.NewRequest("GET", "https://gogo-garage-opener.eu.auth0.com/userinfo", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(errors.New("Could not get email from access token"))
	}
	if resp.StatusCode == 401 {
		log.Warn("Token is not authorised to get email")
		panic(errors.New("Could not get email using token"))
	}
	var userInfo map[string]*json.RawMessage
	body, err := ioutil.ReadAll(resp.Body)
	log.Infof("Got %s", body)
	if err != nil {
		panic(errors.New("Could not parse user info"))
	}
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		panic(errors.New("Could not marsher user info"))
	}
	json, err := userInfo["email"].MarshalJSON()
	if err != nil {
		panic(errors.New("Could not get email marshelled user info"))
	}
	return string(json)
}
