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
	log.Debugf("inserting user, email:[%s]", email)

	tx, _ := userDao.db.Begin()
	prepStmt, err := userDao.db.Prepare("insert into user(email) SELECT ? WHERE NOT EXISTS(SELECT 1 FROM user WHERE email = ?);")
	if err != nil {
		log.Error(err)
	}
	defer prepStmt.Close()
	_, err = prepStmt.Exec(email, email)
	if err != nil {
		log.Error(err)
		tx.Rollback()
	} else {
		tx.Commit()
	}
}

func (userDao UserDao) getUserByEmail(email string) User {
	log.Debugf("getting user for email [%s]", email)

	tx, _ := userDao.db.Begin()
	rows, _ := userDao.db.Query("select lower(email) from user where lower(email) = lower(?)", email)
	defer rows.Close()

	user := getUserFromRows(rows)
	tx.Commit()
	return user
}

func (userDao UserDao) getUserByToken(token string) User {
	log.Debugf("getting user for token [%s]", token)

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
		log.Debugf("Found user %v", emails)
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
		log.Debugf("Found user %v", user)
		return user
	}
	return User{}
}

func (userDao UserDao) setToken(user User) {
	log.Debugf("Updating token [%s]  for email [%s]", user.Token, user.Email)

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
	log.Debugf("Checking valid token [%s]", token)

	tx, _ := userDao.db.Begin()
	rows, _ := userDao.db.Query("select 1 from user where token = ?", token)
	defer rows.Close()
	tx.Commit()
	return rows.Next()
}
