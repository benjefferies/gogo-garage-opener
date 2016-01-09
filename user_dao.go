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
	log.Debugf("inserting user, email:[%s], password:[%s], longitude:[%s], latitude:[%s]", user.Email, user.Password, user.Latitude, user.Longitude,)

	tx,_ := u.db.Begin()
	prepStmt,err := u.db.Prepare("insert into user(email, password, longitude, latitude) values (?, ?, ?, ?)")
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
	rows, _ := u.db.Query("select * from user where email = ?", email)
	defer rows.Close()
	var userEmail string
	var password string
	var latitude float64
	var longitude float64

	for (rows.Next()) {
		rows.Scan(&userEmail, &password, &latitude, &longitude)
	}
	var user User = User{Email: userEmail, Password: password, Latitude: latitude, Longitude: longitude}
	tx.Commit()
	return user
}

func (u UserDao) updateToken(user User) {
	log.Debugf("Updating UUID [%s] token for email [%s]", user.Token, user.Email)

	tx,_ := u.db.Begin()
	prepStmt,err := u.db.Prepare("update user set token = ? where email = ?;")
	defer prepStmt.Close()
	_,err = prepStmt.Exec(user.Token, user.Email)
	if (err != nil) {
		log.Error(err)
		tx.Rollback()
	} else {
		tx.Commit()
	}
}

type TimeWindow struct {

	Time time.Time
	DurationWindow time.Duration

}

func (u UserDao) getTimes(user User) []TimeWindow {
	log.Debugf("Getting times for [%s]", user.Email)

	tx,_ := u.db.Begin()
	rows, _ := u.db.Query("select * from user_time where email = ?", user.Email)
	defer rows.Close()

	var timeWindows []TimeWindow = make([]TimeWindow, 1)
	for (rows.Next()) {
		time := TimeWindow{}
		rows.Scan(&time.Time, &time.DurationWindow)
		timeWindows = append(timeWindows, time)
	}
	tx.Commit()
	return timeWindows
}