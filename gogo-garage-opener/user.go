package main

import (
	"strings"
)

// User holds email and token
type User struct {
	Email, Token string
}

func (user User) getEmail() string {
	return strings.ToLower(user.Email)
}
