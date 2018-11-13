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
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"sync"
	"unicode"

	"github.com/kelseyhightower/envconfig"
	"golang.org/x/net/publicsuffix"
)

const __version__ = "2.1.0"
const SHORT_HEADER = "# DO NOT EDIT"
const LONG_HEADER = "# TRAFFIC OPS NOTE:"
const MAX_RETRIES = 5

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

type Connect struct {
	// URL of reference traffic_ops
	URL         string       `required:"true"`
	Client      *http.Client `ignore:"true"`
	ResultsPath string       `ignore:"true"`
	creds       Creds        `ignore:"true"`
	mutex       *sync.Mutex  `ignore:"true"`
}

func (to *Connect) login(creds Creds) error {
	body, err := json.Marshal(creds)
	if err != nil {
		return err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	to.Client = &http.Client{Transport: tr}
	url := to.URL + `/api/1.3/user/login`
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Create cookiejar so created cookie will be reused
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return err
	}
	to.Client.Jar = jar

	resp, err := to.Client.Do(req)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("Failed to login to Traffic Ops at " + to.URL + " : " + string(data))
	}

	log.Printf("Logged in to %s: %s\n", to.URL, string(data))
	return nil
}

func testRoute(tos []*Connect, route string) {
	// keeps result along with instance -- no guarantee on order collected
	type result struct {
		TO  *Connect
		Res string
		Error error
	}
	var res []result
	ch := make(chan result, len(tos))

	// sanitize routes
	if route[0] == '/' {
		route = route[1:]
	}

	var wg sync.WaitGroup
	var m sync.Mutex

	for _, to := range tos {
		wg.Add(1)
		go func(to *Connect) {
			s, err := to.get(route)
			ch <- result{to, s, err}
			wg.Done()
		}(to)

		wg.Add(1)
		go func() {
			m.Lock()
			defer m.Unlock()
			res = append(res, <-ch)
			wg.Done()
		}()
	}
	wg.Wait()
	close(ch)

	// preliminary error handling
	if len(res) != 2 {
		log.Fatalf("Something wicked happened - expected exactly 2 responses, but got %d!\n", len(res))
	}

	if res[0].Error != nil {
		log.Fatalf("Error occurred `GET`ting %s from %s: %s\n", route, res[0].TO.URL, res[0].Error.Error())
	}

	if res[1].Error != nil {
		log.Fatalf("Error occurred `GET`ting %s from %s: %s\n", route, res[1].TO.URL, res[1].Error.Error())
	}

	// Check for Traffic Ops headers and remove them before comparison
	refResult := res[0].Res
	testResult := res[1].Res
	if strings.Contains(route, "configfiles") {
		refLines := strings.Split(refResult, "\n")
		testLines := strings.Split(testResult, "\n")

		// If the two files have different numbers of lines, they definitely differ
		if len(refLines) != len(testLines) {
			log.Print("Diffs from ", route, " written to")
			p, err := res[0].TO.writeResults(route, refResult)
			if err != nil {
				log.Fatalf("Error writing results for %s: %s\n", route, err.Error())
			}
			log.Print(" ", p)
			p, err = res[1].TO.writeResults(route, testResult)
			if err != nil {
				log.Fatalf("Error writing results for %s: %s\n", route, err.Error())
			}
			log.Print(" and ", p)
			return
		}


		refResult = ""
		testResult = ""

		for _, refLine := range refLines {
			if len(refLine) < len(SHORT_HEADER) {
				refResult += refLine
			} else if refLine[:len(SHORT_HEADER)] != SHORT_HEADER {
				if len(refLine) >= len(LONG_HEADER) {
					if refLine[:len(LONG_HEADER)] != LONG_HEADER {
						refResult += refLine
					}
				} else {
					refResult += refLine
				}
			}
		}

		for _, testLine := range testLines {
			if len(testLine) < len(SHORT_HEADER) {
				testResult += testLine
			} else if testLine[:len(SHORT_HEADER)] != SHORT_HEADER {
				if len(testLine) >= len(LONG_HEADER) {
					if testLine[:len(LONG_HEADER)] != LONG_HEADER {
						testResult += testLine
					}
				} else {
					testResult += testLine
				}
			}
		}
	}

	if refResult == testResult {
		log.Printf("Identical results (%d bytes) from %s", len(res[0].Res), route)
	} else {
		log.Print("Diffs from ", route, " written to")
		for _, r := range res {
			p, err := r.TO.writeResults(route, r.Res)
			if err != nil {
				log.Fatalf("Error writing results for %s: %s", route, err.Error())
			}
			log.Print("  ", p)
		}
	}
}

func (to *Connect) writeResults(route string, res string) (string, error) {
	var dst bytes.Buffer
	if err := json.Indent(&dst, []byte(res), "", "  "); err != nil {
		dst.WriteString(res)
	}

	m := func(r rune) rune {
		if unicode.IsPunct(r) && r != '.' || unicode.IsSymbol(r) {
			return '-'
		}
		return r
	}

	err := os.MkdirAll(to.ResultsPath, 0755)
	if err != nil {
		return "", err
	}

	p := to.ResultsPath + "/" + strings.Map(m, route)
	err = ioutil.WriteFile(p, dst.Bytes(), 0644)
	return p, err
}

func (to *Connect) get(route string) (string, error) {
	url := to.URL + "/" + route

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
<<<<<<< HEAD

	// Should wait for any retries to complete before sending a request
	to.mutex.Lock()
	defer to.mutex.Unlock()
=======
>>>>>>> 4126d2f... Fixed 'compare' tool hiding errors as endpoint output - better error logging

	resp, err := to.Client.Do(req)
	if err != nil {
		log.Println("Connection to " + to.URL + "has been dropped - attempting to reconnect")
		retries := 1
		for ; retries <= MAX_RETRIES; retries++ {
			log.Printf("Retrying connection (#%d)...\n", retries)
			if err := to.login(to.creds); err == nil {
				break
			}
		}

		if retries > MAX_RETRIES {
			to.mutex.Unlock() // prevent zombie threads
			log.Fatalln("Cannot establish connection to " + to.URL + "!")
		}

		// if it fails this time, then I guess we're just done.
		resp, err = to.Client.Do(req)
		if err != nil {
			return "", err
		}
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	return string(data), err
}


func main() {

	routesFileLong := flag.String("file", "", "File listing routes to test (will read from stdin if not given)")
	routesFileShort := flag.String("f", "", "File listing routes to test (will read from stdin if not given)")
	resultsPathLong := flag.String("results_path", "", "Directory where results will be written")
	resultsPathShort := flag.String("r", "", "Directory where results will be written")
	refURL := flag.String("ref_url", "", "The URL for the reference Traffic Ops instance (overrides TO_URL environment variable)")
	testURL := flag.String("test_url", "", "The URL for the testing Traffic Ops instance (overrides TEST_URL environment variable)")
	refUser := flag.String("ref_user", "", "The username for logging into the reference Traffic Ops instance (overrides TO_USER environment variable)")
	refPasswd := flag.String("ref_passwd", "", "The password for logging into the reference Traffic Ops instance (overrides TO_PASSWORD environment variable)")
	testUser := flag.String("test_user", "", "The username for logging into the testing Traffic Ops instance (overrides TEST_USER environment variable)")
	testPasswd := flag.String("test_passwd", "", "The password for logging into the testing Traffic Ops instance (overrides TEST_PASSWORD environment variable)")
	versionLong := flag.Bool("version", false, "Print version information and exit")
	versionShort := flag.Bool("v", false, "Print version information and exit")
	helpLong := flag.Bool("help", false, "Print usage information and exit")
	helpShort := flag.Bool("h", false, "Print usage information and exit")
	flag.Parse()

	// Coalesce long/short form options
	version := *versionLong || *versionShort
	if version {
		fmt.Printf("Traffic Control 'compare' tool v%s\n", __version__)
		os.Exit(0)
	}

	help := *helpLong || *helpShort
	if help {
		flag.Usage()
		os.Exit(0)
	}

	var resultsPath string
	if *resultsPathLong == "" {
		if *resultsPathShort == "" {
			resultsPath = "results"
		} else {
			resultsPath = *resultsPathShort
		}
	} else if *resultsPathShort == "" || *resultsPathShort == *resultsPathLong {
		resultsPath = *resultsPathLong
	} else {
		log.Fatal("Duplicate specification of results path! (Hint: try '-h'/'--help')")
	}

	var routesFile *os.File
	var err error
	if *routesFileLong == "" {
		if *routesFileShort == "" {
			routesFile = os.Stdin
		} else {
			if routesFile, err = os.Open(*routesFileShort); err != nil {
				log.Fatal(err)
			}
			defer routesFile.Close()
		}
	} else if *routesFileShort == "" || *routesFileLong == *routesFileShort {
		if routesFile, err = os.Open(*routesFileLong); err != nil {
			log.Fatal(err)
		}
		defer routesFile.Close()
	} else {
		log.Fatal("Duplicate specification of input file! (Hint: try '-h'/'--help')")
	}

	// refTO, testTO are connections to the two Traffic Ops instances
	var refTO = &Connect{ResultsPath: resultsPath + `/ref`}
	var testTO = &Connect{ResultsPath: resultsPath + `/test`}

	if *refUser != "" && *refPasswd != "" {
		refTO.creds = Creds{*refUser, *refPasswd}
	} else {
		err := envconfig.Process("TO", &refTO.creds)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	if *testUser != "" && *testPasswd != "" {
		testTO.creds = Creds{*testUser, *testPasswd}
	} else {
		err := envconfig.Process("TEST", &testTO.creds)
		if err != nil {
			// if not provided, re-use the same credentials
			testTO.creds = refTO.creds
		}
	}

	if *refURL != "" {
		refTO.URL = *refURL
	} else {
		err := envconfig.Process("TO", refTO)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	if *testURL != "" {
		testTO.URL = *testURL
	} else {
		err := envconfig.Process("TEST", testTO)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	tos := []*Connect{refTO, testTO}

	// Login to the 2 Traffic Ops instances concurrently
	var wg sync.WaitGroup
	wg.Add(len(tos))
	for _, t := range tos {
		go func(to *Connect) {
			log.Print("Login to ", to.URL)
			err := to.login(to.creds)
			if err != nil {
				log.Fatal(err)
			}
			to.mutex = &sync.Mutex{}
			wg.Done()
		}(t)
	}
	wg.Wait()

	var testRoutes []string

	scanner := bufio.NewScanner(routesFile)
	for scanner.Scan() {
		testRoutes = append(testRoutes, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	wg.Add(len(testRoutes))
	for _, route := range testRoutes {
		go func(r string) {
			testRoute(tos, r)
			wg.Done()
		}(route)
	}
	wg.Wait()
}
