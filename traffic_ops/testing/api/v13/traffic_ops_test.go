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

package v13

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/testing/api/config"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/testing/api/todb"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/testing/api/towrap"
	_ "github.com/lib/pq"
)

var (
	TOSession *to.Session
	cfg       config.Config
	testData  TrafficControl
)

func TestMain(m *testing.M) {
	var err error
	configFileName := flag.String("cfg", "traffic-ops-test.conf", "The config file path")
	tcFixturesFileName := flag.String("fixtures", "tc-fixtures.json", "The test fixtures for the API test tool")
	flag.Parse()

	if cfg, err = config.LoadConfig(*configFileName); err != nil {
		fmt.Printf("Error Loading Config %v %v\n", cfg, err)
	}

	if err = log.InitCfg(cfg); err != nil {
		fmt.Printf("Error initializing loggers: %v\n", err)
		return
	}

	log.Infof(`Using Config values:
			   TO Config File:       %s
			   TO Fixtures:          %s
			   TO URL:               %s
			   TO Session Timeout In Secs:  %d
			   DB Server:            %s
			   DB User:              %s
			   DB Name:              %s
			   DB Ssl:               %t`, *configFileName, *tcFixturesFileName, cfg.TrafficOps.URL, cfg.Default.Session.TimeoutInSecs, cfg.TrafficOpsDB.Hostname, cfg.TrafficOpsDB.User, cfg.TrafficOpsDB.Name, cfg.TrafficOpsDB.SSL)

	//Load the test data
	LoadFixtures(*tcFixturesFileName)

	var db *sql.DB
	db, err = todb.OpenConnection(&cfg)
	if err != nil {
		fmt.Printf("\nError opening connection to %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}
	defer db.Close()

	err = todb.Teardown(&cfg, db)
	if err != nil {
		fmt.Printf("\nError tearingdown data %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}

	err = todb.SetupTestData(&cfg, db)
	if err != nil {
		fmt.Printf("\nError setting up data %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}

	TOSession, _, err = towrap.SetupSession(cfg, cfg.TrafficOps.URL, cfg.TrafficOps.User, cfg.TrafficOps.UserPassword)
	if err != nil {
		fmt.Printf("\nError logging into TOURL: %s TOUser: %s/%s - %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, cfg.TrafficOps.UserPassword, err)
		os.Exit(1)
	}

	// Now run the test case
	rc := m.Run()
	os.Exit(rc)

}
