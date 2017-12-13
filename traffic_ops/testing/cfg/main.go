package main

//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"
	"unicode"

	"github.com/kelseyhightower/envconfig"
	"github.comcast.com/cdn/trafficcontrol/lib/go-tc"
	"golang.org/x/net/publicsuffix"
)

var testRoutes = []string{
	`api/1.2/asns?orderby=id`,
	`api/1.2/cdns?orderby=id`,
	`api/1.2/divisions?orderby=id`,
	`api/1.2/parameters?orderby=id`,
	`api/1.2/phys_locations?orderby=id`,
	`api/1.2/regions?orderby=id`,
	`api/1.2/servers?orderby=id`,
	`api/1.2/statuses?orderby=id`,
	`api/1.2/profiles?orderby=id`,
}

// Environment variables used:
//   TO_URL      -- URL for reference Traffic Ops
//   TEST_URL    -- URL for test Traffic Ops
//   TO_USER     -- Username for both instances
//   TO_PASSWORD -- Password for both instances
type Creds struct {
	// common user/password
	User     string `json:"u" required:"true"`
	Password string `json:"p" required:"true"`
}

// Credentials to login to both servers
var creds Creds

type Connect struct {
	// URL of reference traffic_ops
	URL         string `required:"true"`
	Client      *http.Client
	ResultsPath string
}

// refTO, newTO are connections to the two Traffic Ops instances
var refTO = &Connect{ResultsPath: `/tmp/gofiles/ref`}
var newTO = &Connect{ResultsPath: `/tmp/gofiles/new`}

// ResultsPath ...
//var ResultsPath = `/tmp/gofiles/`

func (to *Connect) login(creds Creds) error {
	body, err := json.Marshal(creds)
	if err != nil {
		return err
	}
	to.Client = &http.Client{}

	url := to.URL + `/api/1.2/user/login`
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Create cookiejar for created cookie to be placed into
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return err
	}
	to.Client.Jar = jar

	resp, err := to.Client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Printf("Logged in to %s: %s", to.URL, string(data))
	return nil
}

func doGetRoute(to *Connect, route string, res *[]byte) {
	var err error
	*res, err = to.get(route)
	if err != nil {
		*res = []byte(fmt.Sprintf("Error from %s : %s", to.URL+route, err))
	}
}

func testRoute(route string) {
	var wg sync.WaitGroup
	wg.Add(2)

	var res1, res2 []byte
	go func() {
		doGetRoute(refTO, route, &res1)
		wg.Done()
	}()

	go func() {
		doGetRoute(newTO, route, &res2)
		wg.Done()
	}()

	wg.Wait()

	if bytes.Equal(res1, res2) {
		log.Printf("Identical results (%d bytes) from %s", len(res1), route)
	} else {
		refPath, err := refTO.writeResults(route, res1)
		if err != nil {
			log.Fatal("Error writing results for ", route)
		}

		newPath, err := newTO.writeResults(route, res2)
		if err != nil {
			log.Fatal("Error writing results for ", route)
		}
		log.Print("Diffs from ", route, " written to")
		log.Print("  ", refPath)
		log.Print("  ", newPath)
	}
}

func (to *Connect) writeResults(route string, res []byte) (string, error) {
	var dst bytes.Buffer
	json.Indent(&dst, res, "", "  ")

	m := func(r rune) rune {
		if unicode.IsPunct(r) && r != '.' || unicode.IsSymbol(r) {
			return '-'
		}
		return r
	}

	p := to.ResultsPath + "/" + strings.Map(m, route)
	err := ioutil.WriteFile(p, dst.Bytes(), 0644)
	return p, err
}

func (to *Connect) get(route string) ([]byte, error) {
	url := to.URL + `/` + route
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := to.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	return data, err
}

func getCDNNames(c *Connect) ([]string, error) {
	var res []byte
	doGetRoute(c, `api/1.2/cdns`, &res)
	fmt.Println(string(res))

	var cdnResp tc.CDNsResponse

	err := json.Unmarshal(res, &cdnResp)
	if err != nil {
		return nil, err
	}
	var cdnNames []string
	for _, c := range cdnResp.Response {
		cdnNames = append(cdnNames, c.Name)
	}
	return cdnNames, nil
}

func main() {
	err := envconfig.Process("TO", &creds)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = envconfig.Process("TO", refTO)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = envconfig.Process("TEST", newTO)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Login to the 2 Traffic Ops instances concurrently
	var wg sync.WaitGroup

	tos := []*Connect{refTO, newTO}
	wg.Add(len(tos))
	for _, t := range tos {
		go func(to *Connect) {
			log.Print("Login to ", to.URL)
			err := to.login(creds)
			if err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}(t)
	}

	wg.Wait()

	wg.Add(len(testRoutes))
	for _, route := range testRoutes {
		go func(r string) {
			testRoute(r)
			wg.Done()
		}(route)
	}
	wg.Wait()

	cdnNames, err := getCDNNames(refTO)
	if err != nil {
		panic(err)
	}
	log.Printf("CDNNames are %+v", cdnNames)
	wg.Add(len(cdnNames))
	for _, cdnName := range cdnNames {
		log.Print("CDN ", cdnName)
		go func(c string) {
			testRoute(`api/1.2/cdns/` + c + `/snapshot/new`)
			wg.Done()
		}(cdnName)
	}
	wg.Wait()

}
