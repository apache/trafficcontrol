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

/* ToCheck.go
   This is a simple app that allows you to submit arbitrary
   check data to the Traffic Ops API. You need to supply check
   name, server ID for the check, and the check value as an
   integer. Useful for testing.
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tc "github.com/apache/trafficcontrol/lib/go-tc"
	toclient "github.com/apache/trafficcontrol/traffic_ops/v3-client"
	"github.com/romana/rlog"
)

// Traffic Ops connection params
const AllowInsecureConnections = false
const UserAgent = "go/tc-dscp-monitor"
const UseClientCache = false
const TrafficOpsRequestTimeout = time.Second * time.Duration(10)

var cpath_new string
var statusData tc.ServercheckRequestNullable

type Config struct {
	URL    string `json:"to_url"`
	User   string `json:"user"`
	Passwd string `json:"passwd"`
}

func LoadConfig(file string) (Config, error) {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		return config, err
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config, err
}

func main() {
	// define default config file path
	cpath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		rlog.Error("Config error:", err)
		os.Exit(1)
	}
	cpath_new = strings.Replace(cpath, "/bin/checks", "/conf/check-config.json", 1)

	// command-line flags
	confPtr := flag.String("conf", cpath_new, "Config file path")
	confName := flag.String("name", "undef", "Check name to pass to TO, e.g. 'DSCP'")
	confHost := flag.String("host", "undef", "Name of server to update)")
	confValue := flag.Int("value", -1, "value to send")
	flag.Parse()

	if *confHost == "undef" {
		rlog.Error("Must specify host name to update")
		os.Exit(1)
	}
	if *confName == "undef" {
		rlog.Error("Must specify check name for update to send to TO")
		os.Exit(1)
	}

	// load config json
	config, err := LoadConfig(*confPtr)
	if err != nil {
		rlog.Error("Error loading config:", err)
		os.Exit(1)
	}

	// connect to TO
	session, _, err := toclient.LoginWithAgent(
		config.URL,
		config.User,
		config.Passwd,
		AllowInsecureConnections,
		UserAgent,
		UseClientCache,
		TrafficOpsRequestTimeout)
	if err != nil {
		rlog.Criticalf("An error occurred while logging in: %v\n", err)
		os.Exit(1)
	}

	// Make TO API call for server details
	server, _, err := session.GetServerByHostName(*confHost)
	if err != nil {
		rlog.Criticalf("An error occurred while getting servers: %v\n", err)
		os.Exit(1)
	}

	statusData.ID = &server[0].ID
	statusData.Name = confName
	statusData.Value = confValue
	_, _, err = session.InsertServerCheckStatus(statusData)
	if err != nil {
		rlog.Error("Error updating server check status with TO:", err)
		os.Exit(1)
	}
	fmt.Printf("ID:%d name:%s value:%v\n", server[0].ID, *confName, *confValue)
	os.Exit(0)
}
