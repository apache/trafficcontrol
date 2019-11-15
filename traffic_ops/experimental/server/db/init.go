
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package db

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func createConnectionStringMysql(server, database, user, pass string, port uint) (string, error) {
	defaultMysqlPort := uint(3306)
	if server == "" {
		server = "localhost"
	}
	if port == 0 {
		port = defaultMysqlPort
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True", user, pass, server, port, database), nil
}

func createConnectionStringPostgres(server, database, user, pass string, port uint) (string, error) {
	connString := fmt.Sprintf("dbname=%s user=%s password=%s sslmode=disable", database, user, pass)
	if server != "" {
		connString += fmt.Sprintf(" host=%s", server)
	}
	if server != "" {
		connString += fmt.Sprintf(" port=%d", port)
	}
	return connString, nil
}

func createConnectionString(dbtype, username, password, environment, server string, port uint) (string, error) {
	if dbtype == "mysql" {
		return createConnectionStringMysql(server, environment, username, password, port)
	} else if dbtype == "postgres" {
		return createConnectionStringPostgres(server, environment, username, password, port)
	}
	return "", errors.New("invalid database type")
}

// InitializeDatabase initializes the database and returns the db variable.
// The server is optional, and defaults to localhost if empty
// The port is optional, and defaults to the default database port if 0
func InitializeDatabase(dbtype, username, password, environment, server string, port uint) (*sqlx.DB, error) {
	connString, err := createConnectionString(dbtype, username, password, environment, server, port)
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Connect(dbtype, connString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
