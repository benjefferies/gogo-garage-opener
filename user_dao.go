package main
import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
	"time"
)

type UserDao struct {
	db sql.DB
}

func (u UserDao) createUser(user User) {
	log.Debugf("inserting user, email:[%s], password:[%s], longitude:[%s], latitude:[%s], time:[%s], duration:[%s], distance:[%s]",
		user.Email, user.Password, user.Latitude, user.Longitude, user.Time, user.Duration, user.Distance)
	tx,_ := u.db.Begin()
	prepStmt,err := u.db.Prepare("insert into user values (?, ?, ?, ?, ?, ?, ?)")
	defer prepStmt.Close()
	_,err = prepStmt.Exec(user.Email, user.Password, user.Latitude, user.Longitude, user.Time, user.Duration, user.Distance)
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
	rows, _ := u.db.Query("select * from user where email = ?", email)
	defer rows.Close()
	var userEmail string
	var password string
	var latitude float64
	var longitude float64
	var time time.Time
	var duration int32
	var distance int32

	for (rows.Next()) {
		rows.Scan(&userEmail, &password, &latitude, &longitude, &time, &duration, &distance)
	}
	var user User = User{userEmail, password, latitude, longitude, time, duration, distance}
	tx.Commit()
	return user
}