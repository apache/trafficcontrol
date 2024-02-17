package ultimate_test_harness

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"text/tabwriter"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

const (
	UserAgent = "Traffic Router Load Tests"
)

type TRDetails struct {
	Hostname           string
	IPAddresses        []string
	ClientIP           string
	ClientIPAddressMap IPAddressMap
	Port               int
	DSHost             string
	CDNName            tc.CDNName
}

type IPAddressMap struct {
	Zones []string
	Map   map[string]tc.CoverageZoneLocation
}

type HTTPBenchmark struct {
	Benchmark
	ClientIP      *string
	PathCount     int
	MaxPathLength int
}

var ClientIPAddress *string
var PathCount *int
var MaxPathLength *int

func init() {
	ClientIPAddress = flag.String("ip_address", "", "spoof your client IP address to Traffic Router's geolocation")
	PathCount = flag.Int("path_count", 10000, "the number of paths to generate for use in requests to Delivery Services")
	MaxPathLength = flag.Int("max_path_length", 100, "the maximum length for each generated path")
}

func getCoverageZoneURL(cdnName tc.CDNName) (string, error) {
	snapshot, _, err := TOSession.GetCRConfig(string(cdnName), client.RequestOptions{})
	if err != nil {
		return "", fmt.Errorf("getting the Snapshot of CDN '%s': %s", cdnName, err.Error())
	}
	czPollingURLInterface, ok := snapshot.Response.Config[tc.CoverageZonePollingURL]
	if !ok {
		return "", fmt.Errorf("parameter %s was not found in the Snapshot of CDN '%s'", tc.CoverageZonePollingURL, cdnName)
	}
	czPollingURL := czPollingURLInterface.(string)
	return czPollingURL, nil
}

func getCoverageZoneFile(czPollingURL string) (tc.CoverageZoneFile, error) {
	czMap := tc.CoverageZoneFile{}
	czMapRequest, err := http.NewRequest("GET", czPollingURL, nil)
	if err != nil {
		return czMap, fmt.Errorf("creating HTTP request for URL %s: %s", czPollingURL, err.Error())
	}
	czMapRequest.Header.Set("User-Agent", UserAgent)
	httpClient := http.Client{Timeout: time.Duration(TOConfig.TOTimeout) * time.Second, Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: TOConfig.TOInsecure}}}
	czMapResponse, err := httpClient.Do(czMapRequest)
	if err != nil {
		return czMap, fmt.Errorf("getting Coverage Zone File from URL %s: %s", czPollingURL, err.Error())
	}
	defer log.Close(czMapResponse.Body, "closing the Coverage Zone File response")
	czMapBytes, err := ioutil.ReadAll(czMapResponse.Body)
	if err != nil {
		return czMap, fmt.Errorf("reading Coverage Zone File bytes: %s", err.Error())
	}
	if err = json.Unmarshal(czMapBytes, &czMap); err != nil {
		return czMap, fmt.Errorf("unmarshalling Coverage Zone Map bytes: %s", err.Error())
	}
	return czMap, nil
}

func (i *IPAddressMap) buildFromCoverageZoneMap(czMap tc.CoverageZoneFile) error {
	i.Zones = make([]string, len(czMap.CoverageZones))
	i.Map = map[string]tc.CoverageZoneLocation{}
	zoneIndex := 0
	for location, networks := range czMap.CoverageZones {
		coverageZoneLocation := tc.CoverageZoneLocation{
			Network:  make([]string, 2*len(networks.Network)),
			Network6: make([]string, 2*len(networks.Network6)),
		}
		for index, ipAddress := range networks.Network {
			_, ipNet, err := net.ParseCIDR(ipAddress)
			if err != nil {
				return fmt.Errorf("parsing IP address %s in CIDR notation: %s", ipAddress, err.Error())
			}
			coverageZoneLocation.Network[index*2] = util.FirstIP(ipNet).To4().String()
			coverageZoneLocation.Network[index*2+1] = util.LastIP(ipNet).To4().String()
		}
		for index, ipAddress6 := range networks.Network6 {
			_, ipNet, err := net.ParseCIDR(ipAddress6)
			if err != nil {
				return fmt.Errorf("parsing IP address %s in CIDR notation: %s", ipAddress6, err.Error())
			}
			coverageZoneLocation.Network6[index*2] = util.FirstIP(ipNet).To16().String()
			coverageZoneLocation.Network6[index*2+1] = util.LastIP(ipNet).To16().String()
		}
		i.Map[location] = coverageZoneLocation
		i.Zones[zoneIndex] = location
		zoneIndex++
	}
	return nil
}

func buildIPAddressMap(cdnName tc.CDNName) (IPAddressMap, error) {
	ipAddressMap := IPAddressMap{}
	czPollingURL, err := getCoverageZoneURL(cdnName)
	if err != nil {
		return ipAddressMap, fmt.Errorf("getting Coverage Zone Polling URL from the Snapshot of CDN '%s': %s", cdnName, err.Error())
	}
	czMap, err := getCoverageZoneFile(czPollingURL)
	if err != nil {
		return ipAddressMap, fmt.Errorf("getting Coverage Zone File: %s", err.Error())
	}
	if err = ipAddressMap.buildFromCoverageZoneMap(czMap); err != nil {
		return ipAddressMap, fmt.Errorf("building IP Address Map from Coverage Zone File: %s", err.Error())
	}

	return ipAddressMap, nil
}

func RunHTTPBenchmarksAgainstTrafficRouters(t *testing.T, benchmark HTTPBenchmark) {
	passedTests := 0
	failedTests := 0

	fmt.Printf("Passing criteria: Routing at least %d requests per second\n", benchmark.RequestsPerSecondThreshold)
	writer := tabwriter.NewWriter(os.Stdout, 20, 8, 1, '\t', tabwriter.AlignRight)
	fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", "Traffic Router", "Protocol", "Delivery Service", "Passed?", "Requests Per Second", "Redirects", "Failures")
	for trafficRouterIndex, trafficRouter := range benchmark.TrafficRouters {
		for ipAddressIndex, ipAddress := range trafficRouter.IPAddresses {
			trafficRouterURL := fmt.Sprintf("http://%s:%d/", ipAddress, trafficRouter.Port)

			isIPv4 := strings.Contains(ipAddress, ".")
			if trafficRouter.ClientIP == "" && trafficRouter.ClientIPAddressMap.Zones != nil {
				if benchmark.CoverageZoneLocation != "" {
					location := trafficRouter.ClientIPAddressMap.Map[benchmark.CoverageZoneLocation]
					trafficRouter.ClientIP = location.GetFirstIPAddressOfType(isIPv4)
				}
				if trafficRouter.ClientIP == "" {
					for _, location := range trafficRouter.ClientIPAddressMap.Map {
						trafficRouter.ClientIP = location.GetFirstIPAddressOfType(isIPv4)
						if trafficRouter.ClientIP != "" {
							break
						}
					}
				}
			}

			redirects, failures := 0, 0
			redirectsChannels := make([]chan int, benchmark.ThreadCount)
			failuresChannels := make([]chan int, benchmark.ThreadCount)
			for threadIndex := 0; threadIndex < benchmark.ThreadCount; threadIndex++ {
				redirectsChannels[threadIndex] = make(chan int)
				failuresChannels[threadIndex] = make(chan int)
				go benchmark.Run(t, redirectsChannels[threadIndex], failuresChannels[threadIndex], trafficRouterIndex, trafficRouterURL, ipAddressIndex)
			}

			for threadIndex := 0; threadIndex < benchmark.ThreadCount; threadIndex++ {
				redirects += <-redirectsChannels[threadIndex]
				failures += <-failuresChannels[threadIndex]
			}
			protocol := "IPv6"
			if isIPv4 {
				protocol = "IPv4"
			}
			var passed string
			requestsPerSecond := redirects / benchmark.BenchmarkSeconds
			if requestsPerSecond > benchmark.RequestsPerSecondThreshold {
				passedTests++
				passed = "Yes"
			} else {
				failedTests++
				passed = "No"
			}
			fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%d\t%d\t%d\n", trafficRouter.Hostname, protocol, trafficRouter.DSHost, passed, requestsPerSecond, redirects, failures)
			writer.Flush()
		}
	}
	summary := fmt.Sprintf("%d out of %d load tests passed", passedTests, passedTests+failedTests)
	if failedTests < 1 {
		t.Logf(summary)
	} else {
		t.Fatal(summary)
	}
}

func (b HTTPBenchmark) Run(t *testing.T, redirectsChannel chan int, failuresChannel chan int, trafficRouterIndex int, trafficRouterURL string, ipAddressIndex int) {
	paths := generatePaths(b.PathCount, b.MaxPathLength)
	stopTime := time.Now().Add(time.Duration(b.BenchmarkSeconds) * time.Second)
	redirects, failures := 0, 0
	var req *http.Request
	var resp *http.Response
	var err error
	httpClient := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 10 * time.Second,
	}
	trafficRouter := b.TrafficRouters[trafficRouterIndex]
	for time.Now().Before(stopTime) {
		requestURL := trafficRouterURL + paths[rand.Intn(len(paths))]
		if req, err = http.NewRequest("GET", requestURL, nil); err != nil {
			t.Errorf("creating GET request to Traffic Router '%s' (IP address %s): %s",
				trafficRouter.Hostname, trafficRouter.IPAddresses[ipAddressIndex], err.Error())
		}
		req.Header.Set("User-Agent", UserAgent)
		if trafficRouter.ClientIP != "" {
			req.Header.Set(tc.X_MM_CLIENT_IP, trafficRouter.ClientIP)
		}
		req.Host = trafficRouter.DSHost
		resp, err = httpClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusFound {
			redirects++
		} else {
			failures++
		}
	}
	redirectsChannel <- redirects
	failuresChannel <- failures
}

func TestHTTPLoad(t *testing.T) {
	var trafficRouterDetails = GetTrafficRouterDetails(t)

	benchmark := HTTPBenchmark{
		Benchmark: Benchmark{
			RequestsPerSecondThreshold: *HTTPRequestsPerSecondThreshold,
			BenchmarkSeconds:           *BenchTime,
			ThreadCount:                *ThreadCount,
			TrafficRouters:             trafficRouterDetails,
			CoverageZoneLocation:       *CoverageZoneLocation,
		},
		PathCount:     *PathCount,
		MaxPathLength: *MaxPathLength,
	}

	RunHTTPBenchmarksAgainstTrafficRouters(t, benchmark)
}

func generatePaths(pathCount, maxPathLength int) []string {
	const alphabetSize = 26 + 26 + 10
	alphabet := make([]rune, alphabetSize)
	index := 0
	for char := 'A'; char <= 'Z'; char++ {
		alphabet[index] = char
		index++
	}
	for char := 'a'; char <= 'z'; char++ {
		alphabet[index] = char
		index++
	}
	for char := '0'; char <= '9'; char++ {
		alphabet[index] = char
		index++
	}
	paths := make([]string, pathCount)
	for index = 0; index < pathCount; index++ {
		pathLength := rand.Intn(maxPathLength)
		generatedURL := make([]rune, pathLength)
		for runeIndex := 0; runeIndex < pathLength; runeIndex++ {
			generatedURL[runeIndex] = alphabet[rand.Intn(alphabetSize)]
		}
		paths[index] = string(generatedURL)
	}
	return paths
}

func getTrafficRouters(trafficRouterName string, cdnName tc.CDNName) ([]tc.ServerV40, error) {
	requestOptions := client.RequestOptions{QueryParameters: url.Values{
		"type":   {tc.RouterTypeName},
		"status": {tc.CacheStatusOnline.String()},
	}}
	if trafficRouterName != "" {
		requestOptions.QueryParameters.Set("hostName", trafficRouterName)
	}
	if cdnName != "" {
		cdnRequestOptions := client.RequestOptions{QueryParameters: url.Values{
			"name": {string(cdnName)},
		}}
		cdnResponse, _, err := TOSession.GetCDNs(cdnRequestOptions)
		if err != nil {
			return nil, fmt.Errorf("requesting a CDN named '%s': %s", cdnName, err.Error())
		}
		cdns := cdnResponse.Response
		if len(cdns) != 1 {
			return nil, fmt.Errorf("did not find exactly 1 CDN with name '%s'", cdnName)
		}
		requestOptions.QueryParameters.Set("cdn", string(cdnName))
	}
	response, _, err := TOSession.GetServers(requestOptions)
	if err != nil {
		return nil, fmt.Errorf("requesting %s-status Traffic Routers: %s", requestOptions.QueryParameters["status"], err.Error())
	}
	trafficRouters := response.Response
	trafficRoutersV40 := make([]tc.ServerV40, 0)
	for _, tr := range trafficRouters {
		trafficRoutersV40 = append(trafficRoutersV40, tr)
	}
	if len(trafficRouters) < 1 {
		return trafficRoutersV40, fmt.Errorf("no Traffic Routers were found with these criteria: %v", requestOptions.QueryParameters)
	}
	return trafficRoutersV40, nil
}

func getDSes(t *testing.T, cdnId int, dsTypeName tc.DSType, dsName tc.DeliveryServiceName) []tc.DeliveryServiceV4 {
	requestOptions := client.RequestOptions{QueryParameters: url.Values{"name": {dsTypeName.String()}}}
	var dsType tc.Type
	{
		response, _, err := TOSession.GetTypes(requestOptions)
		if err != nil {
			t.Fatalf("getting type %s: %s", requestOptions.QueryParameters["name"], err.Error())
		}
		types := response.Response
		if len(types) != 1 {
			t.Fatalf("did not find exactly 1 type with name '%s'", requestOptions.QueryParameters["name"])
		}
		dsType = types[0]
	}

	requestOptions = client.RequestOptions{QueryParameters: url.Values{
		"cdn":    {strconv.Itoa(cdnId)},
		"type":   {strconv.Itoa(dsType.ID)},
		"status": {tc.CacheStatusOnline.String()},
	}}
	if dsName != "" {
		requestOptions.QueryParameters.Set("xmlId", dsName.String())
	}
	response, _, err := TOSession.GetDeliveryServices(requestOptions)
	if err != nil {
		t.Fatalf("getting Delivery Services with type '%s' (type ID %d): %s", dsType.Name, dsType.ID, err.Error())
	}
	httpDSes := response.Response
	return httpDSes
}
