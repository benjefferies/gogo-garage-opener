package main

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

// UserDao data access for users
type UserDao struct {
	db *sql.DB
}

func (userDao UserDao) createUser(email string) {
	log.WithField("email", email).Debug("inserting user")

	tx, _ := userDao.db.Begin()
	prepStmt, err := userDao.db.Prepare("insert into user(email) SELECT ? WHERE NOT EXISTS(SELECT 1 FROM user WHERE email = ?);")
	if err != nil {
		log.WithError(err).Error("Could prepare statement to create user")
	}
	defer prepStmt.Close()
	_, err = prepStmt.Exec(email, email)
	if err != nil {
		log.WithError(err).Error("Could not create user")
		tx.Rollback()
	} else {
		tx.Commit()
	}
}

func (userDao UserDao) getUserByEmail(email string) User {
	log.WithField("email", email).Debug("getting user")

	tx, _ := userDao.db.Begin()
	rows, _ := userDao.db.Query("select lower(email) from user where lower(email) = lower(?)", email)
	defer rows.Close()

	user := getUserFromRows(rows)
	tx.Commit()
	return user
}

func (userDao UserDao) getUserByToken(token string) User {
	log.WithField("token", token).Debug("getting user")

	tx, _ := userDao.db.Begin()
	rows, _ := userDao.db.Query("select lower(email) from user where token = ?", token)
	defer rows.Close()

	user := getUserFromRows(rows)
	tx.Commit()
	return user
}

func (userDao UserDao) getSubscribedUserEmails() []string {
	tx, _ := userDao.db.Begin()
	rows, _ := userDao.db.Query("select lower(email) from user where subscribed = 1")
	defer rows.Close()
	var emails []string
	for rows.Next() {
		var userEmail string
		rows.Scan(&userEmail)
		log.WithField("email", userEmail).Debug("Found user")
		emails = append(emails, userEmail)
	}
	tx.Commit()
	return emails
}

func getUserFromRows(rows *sql.Rows) User {
	for rows.Next() {
		var userEmail string

		rows.Scan(&userEmail)
		user := User{Email: userEmail}
		log.WithField("email", user.Email).Debug("Found user")
		return user
	}
	return User{}
}

func (userDao UserDao) setToken(user User) {
	log.WithField("token", user.Token).WithField("email", user.Email).Debug("Updating token")

	tx, _ := userDao.db.Begin()
	prepStmt, err := userDao.db.Prepare("update user set token = ? where email = ?;")
	defer prepStmt.Close()
	_, err = prepStmt.Exec(user.Token, user.Email)
	if err != nil {
		log.Error(err)
		tx.Rollback()
	} else {
		tx.Commit()
	}
}

func (userDao UserDao) tokenExists(token string) bool {
	log.WithField("token", token).Debug("Checking valid token")

	tx, _ := userDao.db.Begin()
	rows, _ := userDao.db.Query("select 1 from user where token = ?", token)
	defer rows.Close()
	tx.Commit()
	return rows.Next()
}
