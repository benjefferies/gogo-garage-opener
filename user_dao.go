package main
import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
)

type UserDao struct {
	db sql.DB
}

func (u UserDao) createUser(user User) {
	log.Debug("inserting user, email:[%s], password:[%s], longitude:[%s], latitude:[%s]", user.Email, user.Password, user.Latitude, user.Longitude)
	tx,_ := u.db.Begin()
	prepStmt,_ := u.db.Prepare("insert into user values (?, ?, ?, ?)")
	defer prepStmt.Close()
	prepStmt.Exec(user.Email, user.Password)
	err := tx.Commit()
	if (err != nil) {
		log.Error(err)
		tx.Rollback()
		return
	}
}

func (u UserDao) getUser(email string) User {
	log.Debug("getting user for email [%s]", email)
	rows, err := u.db.Query("select email, password, latitude, longitude from user")
	if (err != nil) {
		log.Error(err)
	}
	cols,err := rows.Columns()
	if (err != nil) {
		log.Error(err)
	}
	return User{cols[0], cols[1], cols[2], cols[3]}
}