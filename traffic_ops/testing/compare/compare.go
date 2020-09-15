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
	"io"
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

const __version__ = "4.0.0"
const SHORT_HEADER = "# DO NOT EDIT"
const LONG_HEADER = "# TRAFFIC OPS NOTE:"
const LUA_HEADER = "-- DO NOT EDIT"
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

// keeps result along with instance -- no guarantee on order collected
type result struct {
	TO    *Connect
	Res   *http.Response
	Error error
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
	url := to.URL + `/api/2.0/user/login`
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
			resp, err := to.get(route)
			ch <- result{to, resp, err}
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

	ctypeA, ctypeB := res[0].Res.Header.Get("Content-Type"), res[1].Res.Header.Get("Content-Type")
	if ctypeA != ctypeB {
		log.Printf("ERROR: Differing content types for route %s - %s reports %s but %s reports %s !\n",
			route, res[0].TO.URL, ctypeA, res[1].TO.URL, ctypeB)
		return
	}

	// Handle JSON data - note that this WILL NOT be used for endpoints that report the wrong content-type
	// (ignores charset encoding)
	if strings.Contains(ctypeA, "application/json") {
		handleJSONResponse(&res, route)

		// WARNING: treats ALL non-JSON responses as plaintext - should usually operate as expected, but
		// optimizations could be made for other structures
	} else {
		handlePlainTextResponse(&res, route)
	}
}

// Reads in the bodies of responses, closing them as soon as possible
func readRespBodies(a *io.ReadCloser, b *io.ReadCloser) ([]byte, []byte, error) {
	defer (*a).Close()
	defer (*b).Close()

	aBody, err := ioutil.ReadAll(*a)
	if err != nil {
		return nil, nil, err
	}

	bBody, err := ioutil.ReadAll(*b)
	if err != nil {
		return nil, nil, err
	}

	return aBody, bBody, nil
}

// Scrubs out the traffic ops headers from the passed lines
// Note that this assumes UNIX line endings
func scrubPlainText(lines []string) string {
	r := ""
	for _, l := range lines {
		if len(l) >= len(SHORT_HEADER) && l[:len(SHORT_HEADER)] == SHORT_HEADER {
			continue
		}

		if len(l) >= len(LUA_HEADER) && l[:len(LUA_HEADER)] == LUA_HEADER {
			continue
		}

		if len(l) >= len(LONG_HEADER) && l[:len(LONG_HEADER)] == LONG_HEADER {
			continue
		}

		r += l + "\n"
	}

	return r
}

func hashlines(lines []string) map[string]struct{} {
	m := make(map[string]struct{})

	for _, l := range lines {
		m[l] = struct{}{}
	}

	return m
}

// Given a slice of (exactly two) result objects, compares the plain text content of their responses
// and write them to files if they differ. Ignores Traffic Ops headers in the response (to the
// degree possible)
func handlePlainTextResponse(responses *[]result, route string) {

	// I avoid using `defer` to close the bodies because I want to do it as quickly as possible
	result0, result1, err := readRespBodies(&(*responses)[0].Res.Body, &(*responses)[1].Res.Body)
	if err != nil {
		log.Fatalf("Failed to read response body from %s: %s\n", route, err.Error())
	}

	// Check for Traffic Ops headers and remove them before comparison
	result0Str, result1Str := string(result0), string(result1)
	scrubbedResult0, scrubbedResult1 := "", ""
	if strings.Contains(route, "configfiles") {
		lines0 := strings.Split(result0Str, "\n")
		lines1 := strings.Split(result1Str, "\n")

		// If the two files have different numbers of lines, they definitely differ
		if len(lines0) != len(lines1) {
			writeAllResults(route, result0Str, (*responses)[0].TO, result1Str, (*responses)[1].TO, false)
			return
		}

		scrubbedResult0 = scrubPlainText(lines0)
		scrubbedResult1 = scrubPlainText(lines1)

	} else {
		scrubbedResult0 = result0Str
		scrubbedResult1 = result1Str
	}

	if scrubbedResult0 == scrubbedResult1 {
		log.Printf("Identical results (%d bytes) from %s\n", len(result0), route)
	} else {
		writeAllResults(route,
			result0Str,
			(*responses)[0].TO,
			result1Str,
			(*responses)[1].TO,
			checkOrderDiffs(scrubbedResult0, scrubbedResult1))
	}
}

// This function checks for order-only differences in the passed plaintext API responses - the message
// output will indicate whether the difference was caused purely by an ordering of lines or not.
func checkOrderDiffs(s0 string, s1 string) bool {
	m0 := hashlines(strings.Split(s0, "\n"))
	m1 := hashlines(strings.Split(s1, "\n"))

	for k, _ := range m0 {
		if _, ok := m1[k]; !ok {
			return false
		}
		delete(m1, k)
	}

	for k, _ := range m1 {
		if _, ok := m0[k]; !ok {
			return false
		}
	}
	return true
}

// Removes keys that generate false positives in comparisons from the passed JSON object
func sanitizeJSON(m map[string]interface{}) map[string]interface{} {
	// Need to make a full copy so we don't modify while iterating
	object := m

	// ... Now we have to iterate over every key in each map to determine if it should be removed...
	for key, value := range m {

		// handles timestamp/hostname/version/user differences in snapshot and snapshot/new
		if key == "response" {
			switch value.(type) {
			case map[string]interface{}:
				response := value.(map[string]interface{})

				if k, in := response["stats"]; in {
					switch k.(type) {
					case map[string]interface{}:
						stats := k.(map[string]interface{})

						if v, ok := stats["date"]; ok {
							switch v.(type) {
							case float64:
								delete(object["response"].(map[string]interface{})["stats"].(map[string]interface{}), "date")
							}
						}

						if v, ok := stats["tm_host"]; ok {
							switch v.(type) {
							case string:
								delete(object["response"].(map[string]interface{})["stats"].(map[string]interface{}), "tm_host")
							}
						}

						if v, ok := stats["tm_version"]; ok {
							switch v.(type) {
							case string:
								delete(object["response"].(map[string]interface{})["stats"].(map[string]interface{}), "tm_version")
							}
						}

						if v, ok := stats["tm_user"]; ok {
							switch v.(type) {
							case string:
								delete(object["response"].(map[string]interface{})["stats"].(map[string]interface{}), "tm_user")
							}
						}
					}
				}
			}

			// Handles hostname differences in api/1.x/servers/{{server}}/configfiles/ats endpoints
		} else if key == "info" {
			switch value.(type) {
			case map[string]interface{}:
				info := value.(map[string]interface{})

				if v, ok := info["toUrl"]; ok {
					switch v.(type) {
					case string:
						delete(object["info"].(map[string]interface{}), "toUrl")
					}
				}

				if v, ok := info["toRevProxyUrl"]; ok {
					switch v.(type) {
					case string:
						delete(object["info"].(map[string]interface{}), "toRevProxyUrl")
					}
				}
			}
		}
	}

	return object
}

// Given a slice of (exactly two) result objects, compares the JSON content of their responses
// and write them to files if they differ. Ignores timestamps and Traffic Ops hostnames (to the
// degree possible)
func handleJSONResponse(responses *[]result, route string) {

	// I avoid using `defer` to close the bodies because I want to do it as quickly as possible

	result0Orig, result1Orig, err := readRespBodies(&(*responses)[0].Res.Body, &(*responses)[1].Res.Body)
	if err != nil {
		log.Fatalf("Failed to read response body from %s: %s\n", route, err.Error())
	}

	var result0, result1 map[string]interface{}
	if err = json.Unmarshal(result0Orig, &result0); err != nil {
		log.Fatalf("Failed to parse response body from %s/%s as JSON: %s\n", (*responses)[0].TO.URL, route, err.Error())
	}

	if err = json.Unmarshal(result1Orig, &result1); err != nil {
		log.Fatalf("Failed to parse response body from %s/%s as JSON: %s\n", (*responses)[1].TO.URL, route, err.Error())
	}

	result0Bytes, err := json.Marshal(sanitizeJSON(result0))
	if err != nil {
		log.Fatalf("Error re-encoding JSON response from %s/%s: %s\n", (*responses)[0].TO.URL, route, err.Error())
	}

	result1Bytes, err := json.Marshal(sanitizeJSON(result1))
	if err != nil {
		log.Fatalf("Error re-encoding JSON response from %s/%s: %s\n", (*responses)[1].TO.URL, route, err.Error())
	}

	if string(result0Bytes) == string(result1Bytes) {
		log.Printf("Identical results (%d bytes) from %s\n", len(result0Bytes), route)
	} else {
		writeAllResults(route, string(result0Orig), (*responses)[0].TO, string(result1Orig), (*responses)[1].TO, false)
	}
}

// Writes out a set of results for a given route, and logs to stderr information about what was
// written
func writeAllResults(route string, result0 string, connect0 *Connect, result1 string, connect1 *Connect, orderOnly bool) {
	p0, err := connect0.writeResults(route, result0)
	if err != nil {
		log.Fatalf("Error writing results for %s: %s", route, err.Error())
	}

	p1, err := connect1.writeResults(route, result1)
	if err != nil {
		log.Fatalf("Error writing results for %s: %s", route, err.Error())
	}

	if orderOnly {
		log.Println("Order-only diffs from ", route, " written to ", p0, " and ", p1)
	} else {
		log.Println("Diffs from ", route, " written to ", p0, " and ", p1)
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

func (to *Connect) get(route string) (*http.Response, error) {
	url := to.URL + "/" + route

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	// Should wait for any retries to complete before sending a request
	to.mutex.Lock()
	defer to.mutex.Unlock()

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
			return nil, err
		}
	}

	// check for protocol-level errors
	if err == nil && (resp.StatusCode < 200 || resp.StatusCode >= 300) {
		log.Fatalf("Got status %s from %s\n", resp.Status, url)
	}

	return resp, err
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

	scanner := bufio.NewScanner(routesFile)
	for scanner.Scan() {
		wg.Add(1)
		go func(r string) {
			testRoute(tos, r)
			wg.Done()
		}(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}
