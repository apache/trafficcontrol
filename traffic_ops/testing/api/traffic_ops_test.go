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

package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	_ "github.com/lib/pq"
)

var (
	TOSession *client.Session
	cfg       Config
	testData  TrafficControl
)

func TestMain(m *testing.M) {
	var err error
	configFileName := flag.String("cfg", "traffic-ops-test.conf", "The config file path")
	tcFixturesFileName := flag.String("fixtures", "tc-fixtures.json", "The test fixtures for the API test tool")
	flag.Parse()

	if cfg, err = LoadConfig(*configFileName); err != nil {
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
	loadTestCDN(*tcFixturesFileName)

	var db *sql.DB
	db, err = openConnection(&cfg)
	if err != nil {
		fmt.Printf("\nError opening connection to %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}
	defer db.Close()

	err = teardownData(&cfg, db)
	if err != nil {
		fmt.Printf("\nError tearingdown data %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}

	err = setupUserData(&cfg, db)
	if err != nil {
		fmt.Printf("\nError setting up data %s - %s, %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, err)
		os.Exit(1)
	}

	TOSession, _, err = setupSession(cfg, cfg.TrafficOps.URL, cfg.TrafficOps.User, cfg.TrafficOps.UserPassword)
	if err != nil {
		fmt.Printf("\nError logging into TOURL: %s TOUser: %s/%s - %v\n", cfg.TrafficOps.URL, cfg.TrafficOps.User, cfg.TrafficOps.UserPassword, err)
		os.Exit(1)
	}

	// Now run the test case
	rc := m.Run()
	os.Exit(rc)

}

func setupSession(cfg Config, toURL string, toUser string, toPass string) (*client.Session, net.Addr, error) {
	var err error
	var session *client.Session
	var netAddr net.Addr
	toReqTimeout := time.Second * time.Duration(cfg.Default.Session.TimeoutInSecs)
	session, netAddr, err = client.LoginWithAgent(toURL, toUser, toPass, true, "to-api-client-tests", true, toReqTimeout)
	if err != nil {
		return nil, nil, err
	}

	return session, netAddr, err
}

func loadTestCDN(fixturesPath string) {

	f, err := ioutil.ReadFile(fixturesPath)
	if err != nil {
		log.Errorf("Cannot unmarshal fixtures json %s", err)
		os.Exit(1)
	}
	err = json.Unmarshal(f, &testData)
	if err != nil {
		log.Errorf("Cannot unmarshal fixtures json %v", err)
		os.Exit(1)
	}
}

//Request sends a request to TO and returns a response.
//This is basically a copy of the private "request" method in the tc.go \
//but I didn't want to make that one public.
func Request(to client.Session, method, path string, body []byte) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", TOSession.URL, path)

	var req *http.Request
	var err error

	if body != nil && method != "GET" {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, err
		}
	}

	resp, err := to.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		e := client.HTTPError{
			HTTPStatus:     resp.Status,
			HTTPStatusCode: resp.StatusCode,
			URL:            url,
		}
		return nil, &e
	}

	return resp, nil
}
