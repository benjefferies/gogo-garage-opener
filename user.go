package main

import (
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type User struct {
	Email, Password, Token string
}

func (this User) hashPassword() (User, error) {
	passwordBytes := []byte(this.Password)
	password, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	return User{Email: this.Email, Password: string(password), Token: this.Password}, err
}

func (this User) getEmail() string {
	return strings.ToLower(this.Email)
}
