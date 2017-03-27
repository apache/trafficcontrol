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

package integration

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	traffic_ops "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
)

var (
	to *traffic_ops.Session
)

func init() {
	toURL := flag.String("toURL", "http://localhost:3000", "Traffic Ops URL")
	toUser := flag.String("toUser", "admin", "Traffic Ops user")
	toPass := flag.String("toPass", "password", "Traffic Ops password")
	flag.Parse()
	var loginErr error
	to, loginErr = traffic_ops.Login(*toURL, *toUser, *toPass, true)
	if loginErr != nil {
		fmt.Printf("\nError logging in to %v: %v\nMake sure toURL, toUser, and toPass flags are included and correct.\nExample:  go test -toUser=admin -toPass=pass -toURL=http://localhost:3000\n\n", *toURL, loginErr)
		os.Exit(1)
	}
}

//GetCdn returns a Cdn struct
func GetCdn() (traffic_ops.CDN, error) {
	cdns, err := to.CDNs()
	if err != nil {
		return *new(traffic_ops.CDN), err
	}
	cdn := cdns[0]
	if cdn.Name == "ALL" {
		cdn = cdns[1]
	}
	return cdn, nil
}

//GetProfile returns a Profile Struct
func GetProfile() (traffic_ops.Profile, error) {
	profiles, err := to.Profiles()
	if err != nil {
		return *new(traffic_ops.Profile), err
	}
	return profiles[0], nil
}

//GetType returns a Type Struct
func GetType(useInTable string) (traffic_ops.Type, error) {
	types, err := to.Types()
	if err != nil {
		return *new(traffic_ops.Type), err
	}
	for _, myType := range types {
		if myType.UseInTable == useInTable {
			return myType, nil
		}
	}
	nfErr := fmt.Sprintf("No Types found for useInTable %s\n", useInTable)
	return *new(traffic_ops.Type), errors.New(nfErr)
}

//GetDeliveryService returns a DeliveryService Struct
func GetDeliveryService(cdn string) (traffic_ops.DeliveryService, error) {
	dss, err := to.DeliveryServices()
	if err != nil {
		return *new(traffic_ops.DeliveryService), err
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
//This is basically a copy of the private "request" method in the traffic_ops.go \
//but I didn't want to make that one public.
func Request(to traffic_ops.Session, method, path string, body []byte) (*http.Response, error) {
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
		e := traffic_ops.HTTPError{
			HTTPStatus:     resp.Status,
			HTTPStatusCode: resp.StatusCode,
			URL:            url,
		}
		return nil, &e
	}

	return resp, nil
}
