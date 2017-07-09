package main

import (
	"database/sql"
	"log" // TODO change to traffic_monitor_golang/common/log
)

func preparePrivLevelStmt(db *sql.DB) (*sql.Stmt, error) {
	return db.Prepare("select r.priv_level from tm_user as u join role as r on u.role = r.id where u.username = $1")
}

func hasPrivLevel(privLevelStmt *sql.Stmt, user string, level int) bool {
	var privLevel int
	err := privLevelStmt.QueryRow(user).Scan(&privLevel)
	switch {
	case err == sql.ErrNoRows:
		log.Println("Error checking user " + user + " priv level: user not in database")
		return false
	case err != nil:
		log.Println("Error checking user " + user + " priv level: " + err.Error())
		return false
	default:
		return privLevel >= level
	}
}
