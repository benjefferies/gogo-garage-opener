package main
import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
)

type UserDao struct {
	db sql.DB
}

func (u UserDao) createUser(user User) {
	log.Debugf("inserting user, email:[%s], password:[%s], longitude:[%s], latitude:[%s]", user.Email, user.Password, user.Latitude, user.Longitude)
	tx,_ := u.db.Begin()
	prepStmt,err := u.db.Prepare("insert into user values (?, ?, ?, ?)")
	defer prepStmt.Close()
	_,err = prepStmt.Exec(user.Email, user.Password, user.Latitude, user.Longitude)
	if (err != nil) {
		log.Error(err)
		tx.Rollback()
	} else {
		tx.Commit()
	}
}

func (u UserDao) getUser(email string) User {
	log.Debugf("getting user for email [%s]", email)

	tx,_ := u.db.Begin()
	rows, _ := u.db.Query("select email, password, latitude, longitude from user where email = ?", email)
	defer rows.Close()
	var userEmail string
	var password string
	var latitiude float64
	var longitude float64

	for (rows.Next()) {
		rows.Scan(&userEmail, &password, &latitiude, &longitude)
	}
	var user User = User{userEmail, password, latitiude, longitude}
	tx.Commit()
	return user
}