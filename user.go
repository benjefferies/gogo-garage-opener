package main

import (
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type User struct {
	Email, Password, Token string
}

func (u User) hashPassword() (User, error) {
	passwordBytes := []byte(u.Password)
	password, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	return User{Email: u.Email, Password: string(password), Token: u.Password}, err
}

func (u User) getEmail() string {
	return strings.ToLower(u.Email)
}
