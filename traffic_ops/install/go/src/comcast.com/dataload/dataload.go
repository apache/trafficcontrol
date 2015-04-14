/*
   Copyright 2015 Comcast Cable Communications Management, LLC

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

type Configuration struct {
	Description string `json:"description"`
	DbName      string `json:"dbname"`
	DbHostname  string `json:"hostname"`
	DbUser      string `json:"user"`
	DbPassword  string `json:"password"`
	DbPort      string `json:"port"`
	DbType      string `json:"type"`
}

type Profile struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Parameter struct {
	Name        string `json:"name"`
	Config_File string `json:"config_file"`
	Value       string `json:"value"`
}

type ProfileParameter struct {
	Profile    string `json:"profile"`
	Parameter  string `json:"parameter"`
	ConfigFile string `json:"config_file"`
	Value      string `json:"value"`
}

type CustomParams struct {
	CdnName                string `json:"cdnname"`
	TmInfoUrl              string `json:"tminfo.url"`
	CoverageZonePollingUrl string `json:"coveragezone.polling.url"`
	GeoLocationPollingUrl  string `json:"geolocation.polling.url"`
	DomainName             string `json:"domainname"`
	TmUrl                  string `json:"tmurl.url"`
	GeoLocation6PollingUrl string `json:"geolocation6.polling.url"`
}

type DefaultUsers struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type DataType struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	UseInTable  string `json:"use_in_table"`
}

func main() {
	//read prop file for database credentials
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

	//connect to database
	connectString := dbUsername + ":" + dbPassword + "@" + "tcp(" + dbHostname + ":" + dbPort + ")" + "/" + dbName
	if debug {
		fmt.Println("new connect string:" + connectString)
	}

	db, err := sql.Open("mysql", connectString)
	if err != nil {
		fmt.Println("Couldnt create db handle")
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Can't ping the new database")
		panic(err)
	}

	//read profile json file
	fmt.Println("seeding profile data...")
	file, _ = os.Open("/opt/traffic_ops/install/data/json/profile.json")
	lineCount := 0
	profile := Profile{}
	decoder = json.NewDecoder(file)
	profileInsert, err := db.Prepare("insert ignore into " + dbName + ".profile (name, description) values (?,?)")
	if err != nil {
		fmt.Println("Couldn't prepare profile insert statment")
		panic(err)
	}
	for {
		err := decoder.Decode(&profile)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("name", profile.Name, "description", profile.Description)
		//load profile table
		_, err = profileInsert.Exec(profile.Name, profile.Description)
		if err != nil {
			fmt.Println("The profile Insert failed")
			panic(err)
		}
		lineCount += 1
	}
	//read parameter json file
	fmt.Println("seeding parameter data...")
	file, err = os.Open("/opt/traffic_ops/install/data/json/parameter.json")
	if err != nil {
		fmt.Println("trouble reading parameter.json")
		panic(err)
	}
	lineCount = 0
	parameter := Parameter{}
	decoder = json.NewDecoder(file)
	parameterInsert, err := db.Prepare("insert ignore into " + dbName + ".parameter (name, config_file, value) values (?,?,?)")
	if err != nil {
		fmt.Println("Couldn't prepare parameter insert statment")
		panic(err)
	}
	for {
		err := decoder.Decode(&parameter)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		//load parameter table
		rows, err := db.Query("select count(name) from "+dbName+".parameter where name = ? and config_file = ? and value = ?", parameter.Name, parameter.Config_File, parameter.Value)
		if err != nil {
			fmt.Println("Couldn't prepare parameter select statment")
			panic(err)
		}
		defer rows.Close()
		var count int
		for rows.Next() {
			if err := rows.Scan(&count); err != nil {
				panic(err)
			}
			// fmt.Println("row count ", count)
		}
		if count == 0 {
			fmt.Println("inserting parameter: name = ", parameter.Name, "config_file = ", parameter.Config_File, "value = ", parameter.Value)
			_, err = parameterInsert.Exec(parameter.Name, parameter.Config_File, parameter.Value)
		} else {
			fmt.Println("parameter already exists!  name = ", parameter.Name, "config_file = ", parameter.Config_File, "value = ", parameter.Value)
		}
		if err != nil {
			fmt.Println("The parameter insert failed")
			panic(err)
		}
		count = 0
		lineCount += 1
	}
	//seed profile_parameter data
	fmt.Println("seeding profile_parameter data...")
	file, _ = os.Open("/opt/traffic_ops/install/data/json/profile_parameter.json")
	lineCount = 0
	profileParameter := ProfileParameter{}
	decoder = json.NewDecoder(file)
	profileParameterInsert, err := db.Prepare("insert ignore into " + dbName + ".profile_parameter (profile, parameter) values ((select id from profile where name = ?), (select id from parameter where name = ? and config_file = ? and value = ?))")
	if err != nil {
		fmt.Println("Couldn't prepare profile_parameter insert statment")
		panic(err)
	}
	for {
		err := decoder.Decode(&profileParameter)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		rows, err := db.Query("select count(profile) from "+dbName+".profile_parameter where profile = (select id from profile where name = ?) and parameter = (select id from parameter where name = ? and config_file = ? and value = ?)", profileParameter.Profile, profileParameter.Parameter, profileParameter.ConfigFile, profileParameter.Value)
		if err != nil {
			fmt.Println("Couldn't prepare profile_parameter select statment")
			panic(err)
		}
		defer rows.Close()
		var count int
		for rows.Next() {
			if err := rows.Scan(&count); err != nil {
				panic(err)
			}
			// fmt.Println("row count ", count)
		}
		if count == 0 {
			fmt.Println("inserting profile parameter value: profile_name =", profileParameter.Profile, "parameter_name =", profileParameter.Parameter, "config_file =", profileParameter.ConfigFile, "value =", profileParameter.Value)
			//load parameter table
			_, err = profileParameterInsert.Exec(profileParameter.Profile, profileParameter.Parameter, profileParameter.ConfigFile, profileParameter.Value)
			if err != nil {
				fmt.Println("The insert failed")
				panic(err)
			}
		} else {
			fmt.Printf("the profile_parameter combination already exists.  Profile Name = %s, Parameter Name = %s, Parameter Config_File = %s, Parameter Value = %s\n", profileParameter.Profile, profileParameter.Parameter, profileParameter.ConfigFile, profileParameter.Value)
			count = 0
		}
		lineCount += 1
	}
	//read type json file
	fmt.Println("seeding type data...")
	file, _ = os.Open("/opt/traffic_ops/install/data/json/type.json")
	lineCount = 0
	dataType := DataType{}
	decoder = json.NewDecoder(file)
	dataTypeInsert, err := db.Prepare("insert ignore into " + dbName + ".type (name, description, use_in_table) values (?,?,?)")
	if err != nil {
		fmt.Println("Couldn't prepare data type insert statment")
		panic(err)
	}
	for {
		err := decoder.Decode(&dataType)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("name", dataType.Name, "description", dataType.Description, "use_in_table", dataType.UseInTable)
		//load profile table
		_, err = dataTypeInsert.Exec(dataType.Name, dataType.Description, dataType.UseInTable)
		if err != nil {
			fmt.Println("The data type Insert failed")
			panic(err)
		}
		lineCount += 1
	}
	//load custom data
	fmt.Println("creating custom parameters...")
	//read param file into struct
	file, _ = os.Open("/opt/traffic_ops/install/data/json/parameters.json") //should probably rename this file to be more meaningful
	decoder = json.NewDecoder(file)
	customParams := CustomParams{}
	err = decoder.Decode(&customParams)
	if err != nil {
		fmt.Println("error:", err)
	}
	//setup constants
	var (
		cdnName                = "CDN_name"
		tmInfoUrl              = "tm.infourl"
		coverageZonePollingUrl = "coveragezone.polling.url"
		geoLocationPollingUrl  = "geolocation.polling.url"
		domainName             = "domain_name"
		tmUrl                  = "tm.url"
		geoLocation6PollingUrl = "geolocation6.polling.url"
		rascalConfig           = "rascal-config.txt"
		crConfig               = "CRConfig.json"
	)
	// insert cdnname data
	fmt.Println("inserting data for ", cdnName)
	//update
	updateParam, err := db.Prepare("UPDATE parameter SET value=? WHERE name = ? and config_file = ?")
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	_, err = tx.Stmt(updateParam).Exec(customParams.CdnName, cdnName, rascalConfig)
	if err != nil {
		fmt.Println("There was an issue creating paramter for", cdnName, " with a value of ", customParams.CdnName)
		panic(err)
	}

	//insert tmInfoUrl data
	//tminfo.url = tm.infourl, global
	fmt.Println("inserting data for ", tmInfoUrl)
	// _, err = parameterInsert.Exec(tmInfoUrl, "global", customParams.TmInfoUrl)
	_, err = tx.Stmt(updateParam).Exec(customParams.TmInfoUrl, "global", tmInfoUrl)
	if err != nil {
		fmt.Println("There was an issue creating paramter for", tmInfoUrl, " with a value of ", customParams.TmInfoUrl)
		panic(err)
	}
	//insert coverageZonePollingUrl data
	//coveragezone.polling.url, CRConfig.json; CCR1
	fmt.Println("inserting data for ", coverageZonePollingUrl)
	// _, err = parameterInsert.Exec(coverageZonePollingUrl, crConfig, customParams.CoverageZonePollingUrl)
	_, err = tx.Stmt(updateParam).Exec(customParams.CoverageZonePollingUrl, crConfig, coverageZonePollingUrl)
	if err != nil {
		fmt.Println("There was an issue creating paramter for", coverageZonePollingUrl, " with a value of ", customParams.CoverageZonePollingUrl)
		panic(err)
	}

	//insert geolocation polling url data
	//geolocation.polling.url, CRConfig.json; CCR1
	fmt.Println("inserting data for ", geoLocationPollingUrl)
	// _, err = parameterInsert.Exec(geoLocationPollingUrl, crConfig, customParams.GeoLocationPollingUrl)
	_, err = tx.Stmt(updateParam).Exec(customParams.GeoLocationPollingUrl, crConfig, geoLocationPollingUrl)
	if err != nil {
		fmt.Println("There was an issue creating paramter for", geoLocationPollingUrl, " with a value of ", customParams.GeoLocationPollingUrl)
		panic(err)
	}

	//insert domain name data
	//domainname = domain_name, CRConfig.json; EDGE1, CCR1, RASCAL1, MID1
	fmt.Println("inserting data for ", domainName)
	// _, err = parameterInsert.Exec(domainName, crConfig, customParams.DomainName)
	_, err = tx.Stmt(updateParam).Exec(customParams.DomainName, crConfig, domainName)
	if err != nil {
		fmt.Println("There was an issue creating paramter for", domainName, " with a value of ", customParams.DomainName)
		panic(err)
	}

	//insert tm url data
	//tmurl =  tm.url, global
	fmt.Println("inserting data for ", tmUrl)
	// _, err = parameterInsert.Exec(tmUrl, "global", customParams.TmUrl)
	_, err = tx.Stmt(updateParam).Exec(customParams.TmUrl, "global", tmUrl)
	if err != nil {
		fmt.Println("There was an issue creating paramter for", tmUrl, " with a value of ", customParams.TmUrl)
		panic(err)
	}
	//insert geoLocation6 data
	//geolocation6.polling.url = geolocation6.polling.url, CRConfig.json; CCR1
	fmt.Println("inserting data for ", geoLocation6PollingUrl)
	// _, err = parameterInsert.Exec(geoLocation6PollingUrl, crConfig, customParams.GeoLocation6PollingUrl)
	_, err = tx.Stmt(updateParam).Exec(customParams.GeoLocation6PollingUrl, crConfig, geoLocation6PollingUrl)
	if err != nil {
		fmt.Println("There was an issue creating paramter for", geoLocation6PollingUrl, " with a value of ", customParams.GeoLocation6PollingUrl)
		panic(err)
	}

	//add default user data
	fmt.Println("Adding default user data")
	file, _ = os.Open("/opt/traffic_ops/install/data/json/users.json") //should probably rename this file to be more meaningful
	decoder = json.NewDecoder(file)
	defaultUsers := DefaultUsers{}
	err = decoder.Decode(&defaultUsers)
	if err != nil {
		fmt.Println("Error reading users data")
		panic(err)
	}
	userInsert, err := db.Prepare("insert ignore into " + dbName + ".tm_user (username, role, local_passwd, new_user, local_user) values (?,(select id from role where name = ?), ?, 1, 0)")
	if err != nil {
		fmt.Println("An error occurred preparing the statmement")
		panic(err)
	}
	_, err = userInsert.Exec(defaultUsers.Username, "admin", defaultUsers.Password)
	if err != nil {
		fmt.Println("An error occured inserting user with name ", defaultUsers.Username)
		panic(err)
	}
}
