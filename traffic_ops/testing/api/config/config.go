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

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/kelseyhightower/envconfig"
)

// Config reflects the structure of the test-to-api.conf file
type Config struct {
	TrafficOps   TrafficOps   `json:"trafficOps"`
	TrafficOpsDB TrafficOpsDB `json:"trafficOpsDB"`
	Default      Default      `json:"default"`
}

// TrafficOps - config section
type TrafficOps struct {
	// URL - The point to the Traffic Ops instance being tested
	URL string `json:"URL" envconfig:"TO_URL" default:"https://localhost:8443"`

	// User - The Traffic Ops test user hitting the API
	User string `json:"user" envconfig:"TO_USER"`

	// UserPassword - The Traffic Ops test user password hitting the API
	UserPassword string `json:"password" envconfig:"TO_USER_PASSWORD"`

	// Insecure - ignores insecure ssls certs that were self-generated
	Insecure bool `json:"sslInsecure" envconfig:"SSL_INSECURE"`
}

// TrafficOpsDB - config section
type TrafficOpsDB struct {
	// Name - Traffic Ops Database name where the test data will be setup
	Name string `json:"dbname" envconfig:"TODB_NAME"`

	// Hostname - Traffic Ops Database hostname where Postgres is running
	Hostname string `json:"hostname" envconfig:"TODB_HOSTNAME"`

	// User - database user that Traffic Ops is using to point to the 'to_test' database
	User string `json:"user" envconfig:"TODB_USER"`

	// Password - database password for the User above
	Password string `json:"password" envconfig:"TODB_PASSWORD"`

	// Port - Postgres port running that has the to_test schema
	Port string `json:"port" envconfig:"TODB_PORT"`

	// DBType - will be 'Pg' by default but tells the Golang database driver
	//          the database type
	DBType string `json:"type" envconfig:"TODB_TYPE"`

	// SSL - Flag that tells the database driver that the Postgres instances has TLS enabled
	SSL bool `json:"ssl" envconfig:"TODB_SSL"`

	// Description - database description
	Description string `json:"description" envconfig:"TODB_DESCRIPTION"`
}

// Default - config section
type Default struct {
	Session Session   `json:"session"`
	Log     Locations `json:"logLocations"`
}

// Session - config section
type Session struct {
	TimeoutInSecs int `json:"timeoutInSecs" envconfig:"SESSION_TIMEOUT_IN_SECS"`
}

// Locations - reflects the structure of the database.conf file
type Locations struct {
	Debug   string `json:"debug"`
	Event   string `json:"event"`
	Error   string `json:"error"`
	Info    string `json:"info"`
	Warning string `json:"warning"`
}

// LoadConfig - reads the config file into the Config struct
func LoadConfig(confPath string) (Config, error) {
	var cfg Config

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
	err := envconfig.Process("traffic-ops-client-tests", &cfg)
	if err != nil {
		fmt.Printf("Cannot parse config: %v\n", err)
		os.Exit(0)
	}

	return cfg, err
}

// ErrorLog - critical messages
func (c Config) ErrorLog() log.LogLocation {
	return log.LogLocation(c.Default.Log.Error)
}

// WarningLog - warning messages
func (c Config) WarningLog() log.LogLocation {
	return log.LogLocation(c.Default.Log.Warning)
}

// InfoLog - information messages
func (c Config) InfoLog() log.LogLocation {
	return log.LogLocation(c.Default.Log.Info)
}

// DebugLog - troubleshooting messages
func (c Config) DebugLog() log.LogLocation {
	return log.LogLocation(c.Default.Log.Debug)
}

// EventLog - access.log high level transactions
func (c Config) EventLog() log.LogLocation {
	return log.LogLocation(c.Default.Log.Event)
}
