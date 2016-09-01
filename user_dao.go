package main

import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
)

type UserDao struct {
	db sql.DB
}

func (u UserDao) createUser(user User) {
	log.Debugf("inserting user, email:[%s], password:[%s]", user.Email, user.Password)

	tx, _ := u.db.Begin()
	prepStmt, err := u.db.Prepare("insert into user(email, password) values (?, ?)")
	defer prepStmt.Close()
	_, err = prepStmt.Exec(user.getEmail(), user.Password)
	if err != nil {
		log.Error(err)
		tx.Rollback()
	} else {
		tx.Commit()
	}
}

func (u UserDao) getUserByEmail(email string) User {
	log.Debugf("getting user for email [%s]", email)

	tx, _ := u.db.Begin()
	rows, _ := u.db.Query("select lower(email), password from user where lower(email) = lower(?)", email)
	defer rows.Close()

	user := getUserFromRows(rows)
	tx.Commit()
	return user
}

func (u UserDao) getUserByToken(token string) User {
	log.Debugf("getting user for token [%s]", token)

	tx, _ := u.db.Begin()
	rows, _ := u.db.Query("select lower(email), password from user where token = ?", token)
	defer rows.Close()

	user := getUserFromRows(rows)
	tx.Commit()
	return user
}

func getUserFromRows(rows *sql.Rows) User {
	for rows.Next() {
		var userEmail, password string

		rows.Scan(&userEmail, &password)
		user := User{Email: userEmail, Password: password}
		log.Debugf("Found user %v", user)
		return user
	}
	return User{}
}

func (u UserDao) setToken(user User) {
	log.Debugf("Updating token [%s]  for email [%s]", user.Token, user.Email)

	tx, _ := u.db.Begin()
	prepStmt, err := u.db.Prepare("update user set token = ? where email = ?;")
	defer prepStmt.Close()
	_, err = prepStmt.Exec(user.Token, user.Email)
	if err != nil {
		log.Error(err)
		tx.Rollback()
	} else {
		tx.Commit()
	}
}

func (u UserDao) tokenExists(token string) bool {
	log.Debugf("Checking valid token [%s]", token)

	tx, _ := u.db.Begin()
	rows, _ := u.db.Query("select 1 from user where token = ?", token)
	defer rows.Close()
	tx.Commit()
	return rows.Next()
}
