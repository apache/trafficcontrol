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

package client_tests

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kelseyhightower/envconfig"
)

var (
	db *sql.DB
)

// Config reflects the structure of the test-to-api.conf file
type Config struct {
	TOURL    string             `json:"toURL" envconfig:"TO_URL" default:"https://localhost:443"`
	TOUser   string             `json:"toUser" envconfig:"TO_USER"`
	Insecure bool               `json:"sslInsecure" envconfig:"SSL_INSECURE"`
	DB       TrafficOpsDatabase `json:"db"`
	Log      Locations          `json:"logLocations"`
}

type TrafficOpsDatabase struct {
	DBName      string `json:"dbname" envconfig:"TODB_NAME"`
	Hostname    string `json:"hostname" envconfig:"TODB_HOSTNAME"`
	User        string `json:"user" envconfig:"TODB_USER"`
	Password    string `json:"password" envconfig:"TODB_PASSWORD"`
	Port        string `json:"port" envconfig:"TODB_PORT"`
	DBType      string `json:"type" envconfig:"TODB_TYPE"`
	SSL         bool   `json:"ssl" envconfig:"TODB_SSL"`
	Description string `json:"description" envconfig:"TODB_DESCRIPTION"`
}

// ConfigDatabase reflects the structure of the database.conf file
type Locations struct {
	Debug   string `json:"debug"`
	Event   string `json:"event"`
	Error   string `json:"error"`
	Info    string `json:"info"`
	Warning string `json:"warning"`
}

// LoadConfig - reads the config file into the Config struct
func LoadConfig(confPath string) (Config, error) {
	var err error
	var cfg Config
	fmt.Printf("LoadConfig...\n")
	fmt.Printf("confPath ---> %v\n", confPath)

	// load json from cdn.conf
	if _, err := os.Stat(confPath); !os.IsNotExist(err) {
		confBytes, err := ioutil.ReadFile(confPath)
		if err != nil {
			return Config{}, fmt.Errorf("Reading CDN conf '%s': %v", confPath, err)
		}

		err = json.Unmarshal(confBytes, &cfg)
		if err != nil {
			return Config{}, fmt.Errorf("unmarshalling '%s': %v", confPath, err)
		}
	}
	//cfg = Config{}
	err = envconfig.Process("traffic-ops-client-tests", &cfg)
	if err != nil {
		fmt.Printf("Cannot parse config: %v\n", err)
		os.Exit(0)
	}

	//cfg, err = ParseConfig(cfg)

	return cfg, err
}
