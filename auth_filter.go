package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
	"strings"
)

type AuthFilter struct {
	userDao UserDao
}

func (a AuthFilter) tokenFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	authToken := req.Request.Header.Get("X-Auth-Token")

	log.Debugf("Checking token [%s] for path request [%s]", authToken, req.Request.URL.RequestURI())
	tokenExists := a.userDao.tokenExists(authToken)
	urlPath := req.Request.URL.Path
	allowedPaths := strings.HasPrefix(urlPath, "/user/login") || strings.HasPrefix(urlPath, "/garage/one-time-pin") ||
		(strings.HasPrefix(urlPath, "/user/one-time-pin") && req.Request.Method == "GET")
	if !tokenExists && !allowedPaths {
		log.Infof("Not authorized request from [%s]", req.Request.RemoteAddr)
		resp.WriteErrorString(401, "401: Not Authorized")
		return
	} else {
		log.Debugf("Authorized request from [%s]", req.Request.RemoteAddr)
	}
	chain.ProcessFilter(req, resp)
}
