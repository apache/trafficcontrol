package integration

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

import (
	// "database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/traffic_monitor/tests/integration/config"
	"github.com/apache/trafficcontrol/traffic_monitor/tmclient"
)

var Config config.Config
var TMClient *tmclient.TMClient

func TestMain(m *testing.M) {
	var err error
	configFileName := flag.String("cfg", "traffic-monitor-test.conf", "The config file path")
	flag.Parse()

	if Config, err = config.LoadConfig(*configFileName); err != nil {
		fmt.Printf("Error Loading Config %v %v\n", Config, err)
		os.Exit(1)
	}

	if err = log.InitCfg(Config); err != nil {
		fmt.Printf("Error initializing loggers: %v\n", err)
		os.Exit(1)
	}

	log.Infof(`Using Config values:
			   TM Config File:       %s
			   TM URL:               %s
			   TM Session Timeout:   %ds
`, *configFileName, Config.TrafficMonitor.URL, Config.Default.Session.TimeoutInSecs)

	// //Load the test data
	// LoadFixtures(*tcFixturesFileName)

	// var db *sql.DB
	// db, err = OpenConnection()
	// if err != nil {
	// 	fmt.Printf("\nError opening connection to %s - %s, %v\n", Config.TrafficOps.URL, Config.TrafficOpsDB.User, err)
	// 	os.Exit(1)
	// }
	// defer db.Close()

	// err = Teardown(db)
	// if err != nil {
	// 	fmt.Printf("\nError tearingdown data %s - %s, %v\n", Config.TrafficOps.URL, Config.TrafficOpsDB.User, err)
	// 	os.Exit(1)
	// }

	// err = SetupTestData(db)
	// if err != nil {
	// 	fmt.Printf("\nError setting up data %s - %s, %v\n", Config.TrafficOps.URL, Config.TrafficOpsDB.User, err)
	// 	os.Exit(1)
	// }

	tmReqTimeout := time.Second * time.Duration(Config.Default.Session.TimeoutInSecs)

	// err = SetupSession(toReqTimeout, Config.TrafficOps.URL, Config.TrafficOps.Users.Admin, Config.TrafficOps.UserPassword)
	// if err != nil {
	// 	fmt.Printf("\nError creating session to %s - %s, %v\n", Config.TrafficOps.URL, Config.TrafficOpsDB.User, err)
	// 	os.Exit(1)
	// }

	monitorWaitSpan := 30 * time.Second // TODO make configurable?

	if !WaitForMonitor(Config.TrafficMonitor.URL, monitorWaitSpan) {
		fmt.Printf("\nError communicating with Monitor '%v' - didn't return a 200 OK in %v\n", Config.TrafficMonitor.URL, monitorWaitSpan)
		os.Exit(1)
	}

	TMClient = tmclient.New(Config.TrafficMonitor.URL, tmReqTimeout)

	// Now run the test case
	rc := m.Run()
	os.Exit(rc)
}

// WaitForMonitor waits for the monitor to fully start, and stop serving 5xx codes.
// If the monitor does not return a 200 from an API endpoint by timeout, returns false.
func WaitForMonitor(url string, timeout time.Duration) bool {
	httpClient := http.Client{Timeout: timeout}

	tryInterval := time.Second // TODO make configurable?

	start := time.Now()
	for {
		if time.Now().After(start.Add(timeout)) {
			return false
		}
		time.Sleep(tryInterval)
		resp, err := httpClient.Get(strings.TrimSuffix(url, "/") + "/api/version")
		if err != nil {
			continue
		}
		resp.Body.Close()
		if resp.StatusCode != 200 {
			continue
		}
		return true
	}
}
