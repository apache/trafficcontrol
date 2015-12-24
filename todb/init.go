// Copyright 2015 Comcast Cable Communications Management, LLC

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
