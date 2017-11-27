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

package client_tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	log "github.com/apache/incubator-trafficcontrol/lib/go-log"
	tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	_ "github.com/lib/pq"
)

//TODO: drichardson - put these in the config
var (
	to *client.Session
)

func TestMain(m *testing.M) {

	configFileName := flag.String("cfg", "", "The config file path")
	flag.Parse()

	var cfg Config
	var err error
	fmt.Printf("configFileName ---> %v\n", configFileName)
	if cfg, err = LoadConfig(*configFileName); err != nil {
		fmt.Printf("Error Loading Config %v %v\n", cfg, err)
	}

	if err = log.InitCfg(cfg); err != nil {
		fmt.Printf("Error initializing loggers: %v\n", err)
		return
	}
	log.Debugln("cfg ---> %v\n", cfg)

	log.Infof(`Using Config values:
			   TO URL:               %s
			   Db Server:            %s
			   Db User:              %s
			   Db Name:              %s
			   Db Ssl:               %t`, cfg.TOURL, cfg.DB.Hostname, cfg.DB.User, cfg.DB.DBName, cfg.DB.SSL)

	//log.Debugln("Setting up Data")
	prepareDatabase(&cfg)

	TOSession, netAddr, err := setupSession(cfg, cfg.TOURL, cfg.TOUser, cfg.TOUserPassword)
	fmt.Printf("TOSession ---> %v\n", TOSession)
	fmt.Printf("netAddr ---> %v\n", netAddr)
	if err != nil {
		fmt.Printf("\nError logging in to %v: %v\nMake sure toURL, toUser, and toPass flags are included and correct.\nExample:  go test -toUser=%s -toURL=http://localhost:3000\n\n", cfg.TOURL, cfg.TOUserPassword, err)
		os.Exit(1)
	}

}

func getConfigOptionsFromEnv() {
}

func setupSession(cfg Config, toURL string, toUser string, toPass string) (*client.Session, net.Addr, error) {
	var err error
	var TOSession *client.Session
	var netAddr net.Addr
	toReqTimeout := time.Second * time.Duration(30)
	TOSession, netAddr, err = client.LoginWithAgent(toURL, toUser, toPass, true, "traffic-ops-client-integration-tests", true, toReqTimeout)
	if err != nil {
		return nil, nil, err
	}
	log.Debugln("%v-->", toURL)

	return TOSession, netAddr, err
}

func loadFixtureData() TrafficControl {

	fixtureData, err := ioutil.ReadFile("./sample_cdn.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var tc TrafficControl
	err = json.Unmarshal(fixtureData, &tc)
	if err != nil {
		log.Errorf("Cannot unmarshal the json ", err)
	}
	return tc
}

//GetCDN returns a Cdn struct
func GetCDN() (tc.CDN, error) {
	cdns, err := to.CDNs()
	if err != nil {
		return *new(tc.CDN), err
	}
	cdn := cdns[0]
	if cdn.Name == "ALL" {
		cdn = cdns[1]
	}
	return cdn, nil
}

//GetProfile returns a Profile Struct
func GetProfile() (tc.Profile, error) {
	profiles, err := to.Profiles()
	if err != nil {
		return *new(tc.Profile), err
	}
	return profiles[0], nil
}

//GetType returns a Type Struct
func GetType(useInTable string) (tc.Type, error) {
	types, err := to.Types()
	if err != nil {
		return *new(tc.Type), err
	}
	for _, myType := range types {
		if myType.UseInTable == useInTable {
			return myType, nil
		}
	}
	nfErr := fmt.Sprintf("No Types found for useInTable %s\n", useInTable)
	return *new(tc.Type), errors.New(nfErr)
}

//GetDeliveryService returns a DeliveryService Struct
func GetDeliveryService(cdn string) (tc.DeliveryService, error) {
	dss, err := to.DeliveryServices()
	if err != nil {
		return *new(tc.DeliveryService), err
	}
	if cdn != "" {
		for _, ds := range dss {
			if ds.CDNName == cdn {
				return ds, nil
			}
		}
	}
	return dss[0], nil
}

//Request sends a request to TO and returns a response.
//This is basically a copy of the private "request" method in the tc.go \
//but I didn't want to make that one public.
func Request(to client.Session, method, path string, body []byte) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", to.URL, path)

	var req *http.Request
	var err error

	if body != nil && method != "GET" {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest("GET", url, nil)
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
