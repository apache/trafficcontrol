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
	"io"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Configuration contains parameters needed to connect to the Traffic Ops Database.
type Configuration struct {
	Description string `json:"description"`
	DbName      string `json:"dbname"`
	DbHostname  string `json:"hostname"`
	DbUser      string `json:"user"`
	DbPassword  string `json:"password"`
	DbPort      string `json:"port"`
	DbType      string `json:"type"`
}

func loadCdn(db *sql.DB, dbName string) (sql.Result, error) {
	fmt.Println("Seeding cdn data...")

	file, err := os.Open("/opt/traffic_ops/install/data/json/cdn.json")
	if err != nil {
		return nil, err
	}

	type cdnData struct {
		Name string `json:"name"`
	}
	var c cdnData
	if err := json.NewDecoder(file).Decode(&c); err != nil && err != io.EOF {
		return nil, err
	}

	fmt.Printf("\t Inserting cdn: %+v \n", c)
	cdn, err := db.Exec("insert ignore into "+dbName+".cdn (name) values (?)", c.Name)
	if err != nil {
		fmt.Println("\t An error occured inserting cdn with name ", c.Name)
		return nil, err
	}
	return cdn, nil
}

func loadProfile(db *sql.DB, dbName string) error {
	fmt.Println("Seeding profile data...")

	stmt, err := db.Prepare("insert ignore into " + dbName + ".profile (name, description) values (?,?)")
	if err != nil {
		fmt.Println("Couldn't prepare profile insert statment")
		return err
	}

	file, err := os.Open("/opt/traffic_ops/install/data/json/profile.json")
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)

	type profileData struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	for {
		var p profileData
		if err = decoder.Decode(&p); err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		fmt.Printf("\t Inserting profile: %+v \n", p)
		_, err = stmt.Exec(p.Name, p.Description)
		if err != nil {
			fmt.Println("\t An error occured inserting profile with name ", p.Name)
			return err
		}
	}
	return nil
}

func loadParameter(db *sql.DB, dbName string) error {
	fmt.Println("Seeding parameter data...")

	stmt, err := db.Prepare("insert ignore into " + dbName + ".parameter (name, config_file, value) values (?,?,?)")
	if err != nil {
		fmt.Println("Couldn't prepare parameter insert statment")
		return err
	}

	file, err := os.Open("/opt/traffic_ops/install/data/json/parameter.json")
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)

	type parameterData struct {
		Name       string `json:"name"`
		ConfigFile string `json:"config_file"`
		Value      string `json:"value"`
	}

	for {
		var p parameterData
		if err = decoder.Decode(&p); err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		query := "select count(name) from " + dbName + ".parameter where name = ? and config_file = ? and value = ?"
		rows, err := db.Query(query, p.Name, p.ConfigFile, p.Value)
		if err != nil {
			return err
		}
		defer rows.Close()

		var count int
		for rows.Next() {
			if err := rows.Scan(&count); err != nil {
				return err
			}
		}

		if count != 0 {
			fmt.Printf("\t Parameter already exists: %+v \n", p)
			count = 0
		} else {
			fmt.Printf("\t Inserting parameter: %+v \n", p)
			_, err = stmt.Exec(p.Name, p.ConfigFile, p.Value)
			if err != nil {
				fmt.Println("\t An error occured inserting parameter with name ", p.Name)
				return err
			}
		}
	}
	return nil
}

func loadProfileParameter(db *sql.DB, dbName string) error {
	fmt.Println("Seeding profile_parameter data...")

	stmt, err := db.Prepare("insert ignore into " + dbName + ".profile_parameter (profile, parameter) values ((select id from profile where name = ?), (select id from parameter where name = ? and config_file = ? and value = ?))")
	if err != nil {
		fmt.Println("Couldn't prepare profile_parameter insert statment")
		return err
	}

	file, err := os.Open("/opt/traffic_ops/install/data/json/profile_parameter.json")
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)

	type profileParameterData struct {
		Profile    string `json:"profile"`
		Parameter  string `json:"parameter"`
		ConfigFile string `json:"config_file"`
		Value      string `json:"value"`
	}

	for {
		var p profileParameterData
		if err = decoder.Decode(&p); err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		query := "select count(profile) from " + dbName + ".profile_parameter where profile = (select id from profile where name = ?) and parameter = (select id from parameter where name = ? and config_file = ? and value = ?)"
		rows, err := db.Query(query, p.Profile, p.Parameter, p.ConfigFile, p.Value)
		if err != nil {
			return err
		}
		defer rows.Close()

		var count int
		for rows.Next() {
			if err := rows.Scan(&count); err != nil {
				return err
			}
		}

		if count != 0 {
			fmt.Printf("\t Profile Parameter combination already exists:  %+v \n", p)
			count = 0
		} else {
			fmt.Printf("\t Inserting profile parameter: %+v \n", p)
			_, err = stmt.Exec(p.Profile, p.Parameter, p.ConfigFile, p.Value)
			if err != nil {
				fmt.Println("\t An error occured inserting profile parameter with profile", p.Profile, "parameter", p.Parameter)
				return err
			}
		}
	}
	return nil
}

func loadType(db *sql.DB, dbName string) error {
	fmt.Println("Seeding type data...")

	stmt, err := db.Prepare("insert ignore into " + dbName + ".type (name, description, use_in_table) values (?,?,?)")
	if err != nil {
		fmt.Println("Couldn't prepare type insert statment")
		return err
	}

	file, err := os.Open("/opt/traffic_ops/install/data/json/type.json")
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)

	type typeData struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		UseInTable  string `json:"use_in_table"`
	}

	for {
		var t typeData
		if err = decoder.Decode(&t); err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		fmt.Printf("\t Inserting type: %+v \n", t)
		_, err = stmt.Exec(t.Name, t.Description, t.UseInTable)
		if err != nil {
			fmt.Println("\t An error occured inserting type with name ", t.Name)
			return err
		}
	}
	return nil
}

func loadStatus(db *sql.DB, dbName string) error {
	fmt.Println("Seeding status data...")

	stmt, err := db.Prepare("insert ignore into " + dbName + ".status (name, description) values (?,?)")
	if err != nil {
		fmt.Println("Couldn't prepare status insert statment")
		panic(err)
	}

	file, err := os.Open("/opt/traffic_ops/install/data/json/status.json")
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)

	type statusData struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	for {
		var s statusData
		if err = decoder.Decode(&s); err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		fmt.Printf("\t Inserting status: %+v \n", s)
		_, err = stmt.Exec(s.Name, s.Description)
		if err != nil {
			fmt.Println("\t An error occured inserting status with name ", s.Name)
			return err
		}
	}
	return nil
}

func loadCustomParams(db *sql.DB, dbName string) error {
	fmt.Println("creating custom parameters...")

	//setup constants
	var (
		tmInfoURL              = "tm.infourl"
		coverageZonePollingURL = "coveragezone.polling.url"
		geoLocationPollingURL  = "geolocation.polling.url"
		domainName             = "domain_name"
		tmURL                  = "tm.url"
		geoLocation6PollingURL = "geolocation6.polling.url"
		crConfig               = "CRConfig.json"
	)

	updateParam, err := db.Prepare("UPDATE parameter SET value=? WHERE name = ? and config_file = ?")
	if err != nil {
		return err
	}

	file, err := os.Open("/opt/traffic_ops/install/data/json/post_install.json")
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)

	type customParams struct {
		TmInfoURL              string `json:"tminfo.url"`
		CoverageZonePollingURL string `json:"coveragezone.polling.url"`
		GeoLocationPollingURL  string `json:"geolocation.polling.url"`
		DomainName             string `json:"domainname"`
		TmURL                  string `json:"tmurl.url"`
		GeoLocation6PollingURL string `json:"geolocation6.polling.url"`
	}

	for {
		var c customParams
		if err = decoder.Decode(&c); err != nil {
			return err
		}

		tx, err := db.Begin()
		if err != nil {
			return err
		}

		fmt.Printf("\t Inserting custom param: value=%s, name=global, config_file=%s \n", c.TmInfoURL, tmInfoURL)
		_, err = tx.Stmt(updateParam).Exec(c.TmInfoURL, "global", tmInfoURL)
		if err != nil {
			fmt.Println("\t An error occured inserting parameter for", tmInfoURL, " with a value of ", c.TmInfoURL)
			return err
		}

		fmt.Printf("\t Inserting custom param: value=%s, name=%s, config_file=%s \n", c.CoverageZonePollingURL, crConfig, coverageZonePollingURL)
		_, err = tx.Stmt(updateParam).Exec(c.CoverageZonePollingURL, crConfig, coverageZonePollingURL)
		if err != nil {
			fmt.Println("\t An error occured inserting parameter for", coverageZonePollingURL, " with a value of ", c.CoverageZonePollingURL)
			return err
		}

		fmt.Printf("\t Inserting custom param: value=%s, name=%s, config_file=%s \n", c.GeoLocationPollingURL, crConfig, geoLocationPollingURL)
		_, err = tx.Stmt(updateParam).Exec(c.GeoLocationPollingURL, crConfig, geoLocationPollingURL)
		if err != nil {
			fmt.Println("\t An error occured inserting paramter for", geoLocationPollingURL, " with a value of ", c.GeoLocationPollingURL)
			return err
		}

		fmt.Printf("\t Inserting custom param: value=%s, name=%s, config_file=%s \n", c.DomainName, crConfig, domainName)
		_, err = tx.Stmt(updateParam).Exec(c.DomainName, crConfig, domainName)
		if err != nil {
			fmt.Println("\t An error occured inserting paramter for", domainName, " with a value of ", c.DomainName)
			return err
		}

		fmt.Printf("\t Inserting custom param: value=%s, name=global, config_file=%s \n", c.TmURL, tmURL)
		_, err = tx.Stmt(updateParam).Exec(c.TmURL, "global", tmURL)
		if err != nil {
			fmt.Println("\t An error occured inserting paramter for", tmURL, " with a value of ", c.TmURL)
			return err
		}

		fmt.Printf("\t Inserting custom param: value=%s, name=%s, config_file=%s \n", c.GeoLocation6PollingURL, crConfig, geoLocation6PollingURL)
		_, err = tx.Stmt(updateParam).Exec(c.GeoLocation6PollingURL, crConfig, geoLocation6PollingURL)
		if err != nil {
			fmt.Println("\t An error occured inserting paramter for", geoLocation6PollingURL, " with a value of ", c.GeoLocation6PollingURL)
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}
	}
}

func loadUsers(db *sql.DB, dbName string) error {
	fmt.Println("Adding default user data")
	file, err := os.Open("/opt/traffic_ops/install/data/json/users.json") //should probably rename this file to be more meaningful
	if err != nil {
		return err
	}

	type users struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(file)

	for {
		var u users
		if err = decoder.Decode(&u); err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		query := "insert ignore into " + dbName + ".tm_user (username, role, local_passwd, new_user) values (?,(select id from role where name = ?), ?, 1)"
		_, err = db.Exec(query, u.Username, "admin", u.Password)
		if err != nil {
			fmt.Println("\t An error occured inserting user with name ", u.Username)
			return err
		}
	}
	return nil
}

func main() {
	//read prop file for database credentials
	file, err := os.Open("/opt/traffic_ops/app/conf/production/database.conf")
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)

	var c Configuration
	if err = decoder.Decode(&c); err != nil {
		panic(err)
	}

	debug := false

	//connect to database
	connectString := c.DbUser + ":" + c.DbPassword + "@" + "tcp(" + c.DbHostname + ":" + c.DbPort + ")" + "/" + c.DbName
	if debug {
		fmt.Println("new connect string:" + connectString)
	}

	db, err := sql.Open("mysql", connectString)
	if err != nil {
		fmt.Println("Couldnt create db handle")
		panic(err)
	}

	if err = db.Ping(); err != nil {
		fmt.Println("Can't ping the new database")
		panic(err)
	}

	// read cdn json file
	_, err = loadCdn(db, c.DbName)
	if err != nil {
		fmt.Println(err)
	}

	// read type json file
	if err = loadType(db, c.DbName); err != nil {
		fmt.Println(err)
	}

	// read status json file
	if err = loadStatus(db, c.DbName); err != nil {
		fmt.Println(err)
	}

	// read params json file
	if err = loadCustomParams(db, c.DbName); err != nil {
		fmt.Println(err)
	}

	// read params json file
	if err = loadUsers(db, c.DbName); err != nil {
		fmt.Println(err)
	}
}
