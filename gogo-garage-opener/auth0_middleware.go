package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/codegangsta/negroni"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	log "github.com/sirupsen/logrus"
)

// Jwks contains an array of JWKs
type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

// JSONWebKeys is a struct for a JWK
type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

// CustomClaims JWT claims
type CustomClaims struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}

func jwtCheckHandleFunc(httpFunc http.HandlerFunc) *negroni.Negroni {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// Verify 'aud' claim
			aud := fmt.Sprintf("https://%s/api", *rs)
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAud {
				return token, errors.New("invalid audience")
			}
			// Verify 'iss' claim
			iss := fmt.Sprintf("https://%s/", *as)
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("invalid issuer")
			}

			cert, err := getPemCert(token)
			if err != nil {
				panic(err.Error())
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})
	return negroni.New(negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
			token := authHeaderParts[1]
			log.Info("Checking scope")
			hasScope := checkScope("email", token)

			if !hasScope {
				log.Info("Invalid scope")
				http.Error(w, "Insufficient scope.", http.StatusForbidden)
			} else {
				next(w, r)
			}
		}), negroni.WrapFunc(exportAccessToken), negroni.Wrap(httpFunc))
}

func exportAccessToken(w http.ResponseWriter, r *http.Request) {
	reqToken := r.Header.Get("Authorization")
	log.WithField("access_token", reqToken).Info("Got access token")
	splitToken := strings.Split(reqToken, "Bearer ")
	accessToken := splitToken[1]
	context.Set(r, "access_token", accessToken)
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get(fmt.Sprintf("https://%s/.well-known/jwks.json", *as))

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}

func checkScope(scope string, tokenString string) bool {
	token, _ := jwt.ParseWithClaims(tokenString, &CustomClaims{}, nil)

	claims, _ := token.Claims.(*CustomClaims)

	hasScope := false
	result := strings.Split(claims.Scope, " ")
	for i := range result {
		if result[i] == scope {
			hasScope = true
		}
	}

	return hasScope
}
