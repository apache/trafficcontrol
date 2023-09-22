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
	"flag"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	client "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
)

type Benchmark struct {
	RequestsPerSecondThreshold int
	BenchmarkSeconds           int
	ThreadCount                int
	DSType                     tc.Type
	TrafficRouters             []*TRDetails
	CoverageZoneLocation       string
}

var IPv4Only *bool
var IPv6Only *bool
var CDNName *string
var DeliveryServiceName *string
var TrafficRouterName *string
var UseCoverageZone *bool
var CoverageZoneLocation *string
var HTTPRequestsPerSecondThreshold *int
var DNSRequestsPerSecondThreshold *int
var BenchTime *int
var ThreadCount *int

func init() {
	rand.Seed(time.Now().UnixNano())
	IPv4Only = flag.Bool("ipv4only", false, "test IPv4 addresses only")
	IPv6Only = flag.Bool("ipv6only", false, "test IPv4 addresses only")
	CDNName = flag.String("cdn", "", "the name of a CDN to search for Delivery Services")
	DeliveryServiceName = flag.String("ds", "", "the name (XMLID) of a Delivery Service to use for tests")
	TrafficRouterName = flag.String("hostname", "", "the hostname of a Traffic Router to use")
	UseCoverageZone = flag.Bool("coverage_zone", false, "whether to use an IP address from the Traffic Router's Coverage Zone File")
	CoverageZoneLocation = flag.String("coverage_zone_location", "", "the Coverage Zone location to use (implies coverage_zone=true)")
	HTTPRequestsPerSecondThreshold = flag.Int("http_requests_threshold", 8000, "the minimum number of HTTP requests per second a Traffic Router must successfully respond to")
	DNSRequestsPerSecondThreshold = flag.Int("dns_requests_threshold", 25000, "the minimum number of DNS requests per second a Traffic Router must successfully respond to")
	BenchTime = flag.Int("benchmark_time", 15, "the duration of each load test in seconds")
	ThreadCount = flag.Int("thread_count", 12, "the number of threads to use for each test")

	log.Init(os.Stderr, os.Stderr, os.Stderr, os.Stderr, os.Stderr)
}

func GetTrafficRouterDetails(t *testing.T) []*TRDetails {
	var err error
	var trafficRouterDetails []*TRDetails

	ipAddressMaps := map[tc.CDNName]IPAddressMap{}

	trafficRouters, err := getTrafficRouters(*TrafficRouterName, tc.CDNName(*CDNName))
	if err != nil {
		t.Fatalf("could not get Traffic Routers: %s", err.Error())
	}

	for _, trafficRouter := range trafficRouters {
		var ipAddresses []string
		for _, serverInterface := range trafficRouter.Interfaces {
			if !serverInterface.Monitor {
				log.Warnf("skipping server interface %s of Traffic Router %s because it is unmonitored\n", serverInterface.Name, *trafficRouter.HostName)
				continue
			}
			ipv4, ipv6 := serverInterface.GetDefaultAddress()
			if ipv4 != "" && !*IPv6Only {
				ipAddresses = append(ipAddresses, ipv4)
			}
			if ipv6 != "" && !*IPv4Only {
				ipAddresses = append(ipAddresses, "["+ipv6+"]")
			}
		}
		if len(ipAddresses) < 1 {
			log.Warnf("need at least 1 monitored service address on an interface of Traffic Router '%s' to use it for benchmarks, but %d such addresses were found\n", *trafficRouter.HostName, len(ipAddresses))
			continue
		}
		dsTypeName := tc.DSTypeHTTP
		httpDSes := getDSes(t, *trafficRouter.CDNID, dsTypeName, tc.DeliveryServiceName(*DeliveryServiceName))
		if len(httpDSes) < 1 {
			t.Errorf("at least 1 Delivery Service with type '%s' is required to run HTTP load tests on Traffic Router '%s', but %d were found", dsTypeName, *trafficRouter.HostName, len(httpDSes))
			continue
		}
		if len(httpDSes[0].ExampleURLs) < 1 {
			log.Warnf("No Example URLs for Delivery Service '%s'. Skipping...\n", *DeliveryServiceName)
			continue
		}
		dsURL, err := url.Parse(httpDSes[0].ExampleURLs[0])
		if err != nil {
			t.Fatalf("parsing Delivery Service URL %s: %s", dsURL, err.Error())
		}
		cdnName := tc.CDNName(*trafficRouter.CDNName)

		singleTrafficRouterDetails := TRDetails{
			Hostname:    *trafficRouter.HostName,
			IPAddresses: ipAddresses,
			ClientIP:    *ClientIPAddress,
			Port:        *trafficRouter.TCPPort,
			DSHost:      dsURL.Host,
			CDNName:     cdnName,
		}
		if *UseCoverageZone {
			_, ok := ipAddressMaps[cdnName]
			if !ok {
				ipAddressMaps[cdnName], err = buildIPAddressMap(cdnName)
				if err != nil {
					t.Fatalf("building IP Address map for CDN '%s': %s", cdnName, err.Error())
				}
			}
			singleTrafficRouterDetails.ClientIPAddressMap = ipAddressMaps[cdnName]
		}
		trafficRouterDetails = append(trafficRouterDetails, &singleTrafficRouterDetails)
	}
	if len(trafficRouterDetails) < 1 {
		t.Fatalf("no Traffic Router with at least 1 HTTP Delivery Service and at least 1 monitored service address was found")
	}
	return trafficRouterDetails
}

func TestMain(m *testing.M) {
	var err error
	if err = flag.Set("test.v", "true"); err != nil {
		fmt.Printf("settings flags 'test.v': %s\n", err.Error())
		os.Exit(1)
	}
	flag.Parse()
	if *CoverageZoneLocation != "" {
		*UseCoverageZone = true
	}

	TOSession, _, err = client.LoginWithAgent(TOConfig.TOURL, TOConfig.TOUser, TOConfig.TOPassword, TOConfig.TOInsecure, UserAgent, true, time.Second*time.Duration(TOConfig.TOTimeout))
	if err != nil {
		fmt.Printf("logging into Traffic Ops server %s: %s\n", TOConfig.TOURL, err.Error())
		os.Exit(1)
	}

	// Run tests
	os.Exit(m.Run())
}
