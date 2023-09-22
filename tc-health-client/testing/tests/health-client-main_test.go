package hctest

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
	"database/sql"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/cache-config/testing/ort-tests/config"
	"github.com/apache/trafficcontrol/v8/cache-config/testing/ort-tests/tcdata"
	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc/totest"
)

const cfgFmt = `Using Config values:
  TO Config File:              %s
  TO Fixtures:                 %s
	TO URL:                      %s
	TO Session Timeout In Secs:  %d
	DB Server:                   %s
	DB User:                     %s
	DB Name:                     %s
	DB Ssl:                      %t
	UseIMS:                      %t`

var (
	Config             config.Config
	testData           totest.TrafficControl
	includeSystemTests bool
	tcd                *tcdata.TCData
	TCD                *tcdata.TCData
)

func TestMain(m *testing.M) {
	tcd = tcdata.NewTCData()
	configFileName := flag.String("cfg", "traffic-ops-test.conf", "The config file path")
	tcFixturesFileName := flag.String("fixtures", "tc-fixtures.json", "The test fixtures")
	cliIncludeSystemTests := *flag.Bool("includeSystemTests", false, "Whether to enable tests that have environment dependencies beyond a database")

	flag.Parse()

	// Skip loading configuration when run with `go test -list=<pat>`. The -list
	// flag does not actually run tests, so configuration data is not needed in
	// that mode. If the user is just trying to list the available tests we
	// don't want to abort with an error about a bad configuration the user
	// doesn't care about yet.
	if f := flag.Lookup("test.list"); f != nil {
		if f.Value.String() != "" {
			os.Exit(m.Run())
		}
	}

	var err error

	if *tcd.Config, err = config.LoadConfig(*configFileName); err != nil {
		fmt.Fprintf(os.Stderr, "Error Loading Config: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "Config: %v\n", tcd.Config)

	// CLI option overrides config
	includeSystemTests = tcd.Config.Default.IncludeSystemTests || cliIncludeSystemTests

	if err = log.InitCfg(tcd.Config); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing loggers: %v\n", err)
		os.Exit(1)
	}

	log.Infof(cfgFmt, *configFileName, *tcFixturesFileName, tcd.Config.TrafficOps.URL, tcd.Config.Default.Session.TimeoutInSecs, tcd.Config.TrafficOpsDB.Hostname, tcd.Config.TrafficOpsDB.User, tcd.Config.TrafficOpsDB.Name, tcd.Config.TrafficOpsDB.SSL, tcd.Config.UseIMS)

	//Load the test data
	tcd.LoadFixtures(*tcFixturesFileName)

	var db *sql.DB
	db, err = tcd.OpenConnection()
	if err != nil {
		log.Errorf("\nError opening connection to %s - %s, %v\n", tcd.Config.TrafficOps.URL, tcd.Config.TrafficOpsDB.User, err)
		os.Exit(1)
	}
	defer db.Close()

	err = tcd.Teardown(db)
	if err != nil {
		log.Errorf("\nError tearingdown data %s - %s, %v\n", tcd.Config.TrafficOps.URL, tcd.Config.TrafficOpsDB.User, err)
	}

	err = tcd.SetupTestData(db)
	if err != nil {
		log.Errorf("setting up data on TO instance %s as DB user '%s' failed: %v\n", tcd.Config.TrafficOps.URL, tcd.Config.TrafficOpsDB.User, err)
		os.Exit(1)
	}

	toReqTimeout := time.Second * time.Duration(tcd.Config.Default.Session.TimeoutInSecs)
	err = tcd.SetupSession(toReqTimeout, tcd.Config.TrafficOps.URL, tcd.Config.TrafficOps.Users.Admin, tcd.Config.TrafficOps.UserPassword)
	if err != nil {
		log.Errorf("\nError creating session to %s - %s, %v\n", tcd.Config.TrafficOps.URL, tcd.Config.TrafficOpsDB.User, err)
		os.Exit(1)
	}

	// Now run the test case
	rc := m.Run()
	os.Exit(rc)
}
