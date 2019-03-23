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
