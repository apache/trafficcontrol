/*

     Licensed under the Apache License, Version 2.0 (the "License");
     you may not use this file except in compliance with the License.
     You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

     Unless required by applicable law or agreed to in writing, software
     distributed under the License is distributed on an "AS IS" BASIS,
     WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
     See the License for the specific language governing permissions and
     limitations under the License.
 */

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Configuration struct {
	Description string `json:"description"`
	DbName      string `json:"dbname"`
	DbHostname  string `json:"hostname"`
	DbUser      string `json:"user"`
	DbPassword  string `json:"password"`
	DbPort      string `json:"port"`
	DbType      string `json:"type"`
}

func main() {
	//use command line args
	var adminUser string
	var adminPassword string
	if len(os.Args) != 3 {
		fmt.Println("Usage ./create_db AdminUserName AdminPassword")
		os.Exit(1)
	} else {
		adminUser = os.Args[1]
		adminPassword = os.Args[2]
	}

	//read prop file for database credentials
	//file, _ := os.Open("conf/database.conf")
	file, _ := os.Open("/opt/traffic_ops/app/conf/production/database.conf")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	dbHostname := configuration.DbHostname
	dbUsername := configuration.DbUser
	dbPassword := configuration.DbPassword
	dbPort := configuration.DbPort
	dbName := configuration.DbName
	debug := false

	//connect to DB
	connectString := adminUser + ":" + adminPassword + "@" + "tcp(" + dbHostname + ":" + dbPort + ")" + "/mysql"
	if debug {
		fmt.Println("connect string:" + connectString)
	}

	db, err := sql.Open("mysql", connectString)
	if err != nil {
		fmt.Println("An error occurred")
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Can't ping the database")
		panic(err)
	}
	//create database
	connectString = adminUser + ":" + adminPassword + "@" + "tcp(" + dbHostname + ":" + dbPort + ")" + "/mysql"
	if debug {
		fmt.Println("connect string:" + connectString)
	}

	db, err = sql.Open("mysql", connectString)
	if err != nil {
		fmt.Println("An error occurred")
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Can't ping the database")
		panic(err)
	}

	fmt.Println("Creating database...")
	createResult, err := db.Exec("create database if not exists " + dbName)
	if err != nil {
		fmt.Println("An error occured creating the database")
		panic(err)
	}
	if debug {
		fmt.Println("the database create result was", createResult)
	}

	fmt.Println("Creating user...")
	createResult, err = db.Exec("GRANT ALL PRIVILEGES ON " + dbName + ".* TO '" + dbUsername + "'@'%' IDENTIFIED BY '" + dbPassword + "' WITH GRANT OPTION")
	if err != nil {
		fmt.Println("An error occurred creating the user")
		panic(err)
	}

	createResult, err = db.Exec("GRANT ALL PRIVILEGES ON " + dbName + ".* TO '" + dbUsername + "'@'localhost' IDENTIFIED BY '" + dbPassword + "' WITH GRANT OPTION")
	if err != nil {
		fmt.Println("An error occurred granting privs")
		panic(err)
	}

	fmt.Println("Flushing privileges...")
	createResult, err = db.Exec("flush privileges")
	if err != nil {
		fmt.Println("An error occurred flushing privileges")
		panic(err)
	}
	//create schema
	// fmt.Println("Creating schema...")
	// createResult, err = db.Exec("use twelve_monkeys")
	// if err != nil {
	// 	fmt.Println("An error occured setting the database to use")
	// 	panic(err)
	// }
	// //	//TODO Couldn't get this to work via the commented-out attempts;
	// //	createResult, err = db.Exec("source /opt/traffic_ops/install/schema.sql")
	// //	if err != nil {
	// //		fmt.Println("An error occured creating the schema", err)
	// //	}
	// output, err := exec.Command("/bin/bash", "-c", "mysql --user="+adminUser+" --password="+adminPassword+" "+dbName+" < /opt/traffic_ops/install/data/sql/schema.sql").Output()
	// if err != nil {
	// 	println("exec error:  " + err.Error())
	// 	panic(err)
	// } else {
	// 	fmt.Println(string(output))
	// }

	// if debug {
	// 	fmt.Println("createResult: ", createResult)
	// }

	err = db.Close()
	if err != nil {
		fmt.Println("couldn't close the database connection")
		panic(err)
	}
	//connect to newly created database
	connectString = dbUsername + ":" + dbPassword + "@" + "tcp(" + dbHostname + ":" + dbPort + ")" + "/" + dbName
	if debug {
		fmt.Println("new connect string:" + connectString)
	}

	db, err = sql.Open("mysql", connectString)
	if err != nil {
		fmt.Println("Couldnt create db handle")
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Can't ping the new database")
		panic(err)
	}
}
