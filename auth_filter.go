package main
import (
	"github.com/emicklei/go-restful"
	log "github.com/Sirupsen/logrus"
	"strings"
)

type AuthFilter struct {
	userDao UserDao
}

func (a AuthFilter) tokenFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	authToken := req.Request.Header.Get("X-Auth-Token")

	log.Debugf("Checking token [%s] for path request [%s]", authToken, req.Request.URL.RequestURI())
	isValidToken := a.userDao.validToken(authToken);
	if (strings.Contains(req.Request.URL.String(), "login")) {
		chain.ProcessFilter(req, resp)
	} else if len(authToken) == 0 || !isValidToken {
		log.Infof("Not authorized request from [%s]", authToken)
		resp.WriteErrorString(401, "401: Not Authorized")
		return
	} else {
		log.Debugf("Authorized request fromf [%s]", authToken)
	}
	chain.ProcessFilter(req, resp)
}