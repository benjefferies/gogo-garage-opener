package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"

	log "github.com/sirupsen/logrus"
)

var emailCache = cache.New(60*time.Minute, 120*time.Minute)

func getEmail(accessToken string) string {
	if email, found := emailCache.Get(accessToken); found {
		return email.(string)
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/userinfo", *as), nil)
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
	log.WithField("body", string(body)).Debug("Response from auth0 userinfo")
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
	email := string(json)[1 : len(json)-1]
	emailCache.Set(accessToken, email, cache.DefaultExpiration)
	return email
}
