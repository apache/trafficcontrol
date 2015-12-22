package todb

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	globalDB     sqlx.DB
	databaseName string
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func InitializeDatabase(username string, password string, environment string) {
	db, err := sqlx.Connect("mysql", username+":"+password+"@tcp(localhost:3306)/"+environment+"?parseTime=True")
	check(err)

	globalDB = *db
	databaseName = environment
}
