package main

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// User holds email, password and token
type User struct {
	Email, Password, Token string
}

func (user User) hashPassword() (User, error) {
	passwordBytes := []byte(user.Password)
	password, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	return User{Email: user.Email, Password: string(password), Token: user.Password}, err
}

func (user User) getEmail() string {
	return strings.ToLower(user.Email)
}
