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
	"fmt"
	"os"
	"strings"
	"testing"
	"text/tabwriter"
	"time"

	"github.com/miekg/dns"
)

type DNSBenchmark struct {
	Benchmark
}

func RunDNSBenchmarksAgainstTrafficRouters(t *testing.T, benchmark DNSBenchmark) {
	passedTests := 0
	failedTests := 0

	fmt.Printf("Passing criteria: Routing at least %d requests per second\n", benchmark.RequestsPerSecondThreshold)
	writer := tabwriter.NewWriter(os.Stdout, 20, 8, 1, '\t', tabwriter.AlignRight)
	fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", "Traffic Router", "Protocol", "Delivery Service", "Passed?", "Requests Per Second", "Answers", "Failures")
	for trafficRouterIndex, trafficRouter := range benchmark.TrafficRouters {
		if (*trafficRouter).DSHost[len((*trafficRouter).DSHost)-1] != '.' {
			(*trafficRouter).DSHost += "."
		}
		for ipAddressIndex, ipAddress := range (*trafficRouter).IPAddresses {
			isIPv4 := strings.Contains(ipAddress, ".")
			if (*trafficRouter).ClientIP == "" && (*trafficRouter).ClientIPAddressMap.Zones != nil {
				if benchmark.CoverageZoneLocation != "" {
					location := (*trafficRouter).ClientIPAddressMap.Map[benchmark.CoverageZoneLocation]
					(*trafficRouter).ClientIP = location.GetFirstIPAddressOfType(isIPv4)
				}
				if (*trafficRouter).ClientIP == "" {
					for _, location := range (*trafficRouter).ClientIPAddressMap.Map {
						(*trafficRouter).ClientIP = location.GetFirstIPAddressOfType(isIPv4)
						if (*trafficRouter).ClientIP != "" {
							break
						}
					}
				}
			}

			answers, failures := 0, 0
			redirectsChannels := make([]chan int, benchmark.ThreadCount)
			failuresChannels := make([]chan int, benchmark.ThreadCount)
			for threadIndex := 0; threadIndex < benchmark.ThreadCount; threadIndex++ {
				redirectsChannels[threadIndex] = make(chan int)
				failuresChannels[threadIndex] = make(chan int)
				go benchmark.Run(t, redirectsChannels[threadIndex], failuresChannels[threadIndex], trafficRouterIndex, ipAddressIndex)
			}

			for threadIndex := 0; threadIndex < benchmark.ThreadCount; threadIndex++ {
				answers += <-redirectsChannels[threadIndex]
				failures += <-failuresChannels[threadIndex]
			}
			protocol := "IPv6"
			if isIPv4 {
				protocol = "IPv4"
			}
			var passed string
			requestsPerSecond := answers / benchmark.BenchmarkSeconds
			if requestsPerSecond > benchmark.RequestsPerSecondThreshold {
				passedTests++
				passed = "Yes"
			} else {
				failedTests++
				passed = "No"
			}
			fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%d\t%d\t%d\n", (*trafficRouter).Hostname, protocol, (*trafficRouter).DSHost, passed, requestsPerSecond, answers, failures)
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

func (b DNSBenchmark) Run(t *testing.T, answersChannel chan int, failuresChannel chan int, trafficRouterIndex int, ipAddressIndex int) {
	stopTime := time.Now().Add(time.Duration(b.BenchmarkSeconds) * time.Second)
	answers, failures := 0, 0
	client := dns.Client{Net: "udp", Timeout: 10 * time.Second}
	message := new(dns.Msg)
	trafficRouter := b.TrafficRouters[trafficRouterIndex]
	address := trafficRouter.IPAddresses[ipAddressIndex] + ":53"

	message.SetQuestion(trafficRouter.DSHost, dns.TypeA)
	for time.Now().Before(stopTime) {
		r, _, err := client.Exchange(message, address)
		if err == nil && len(r.Answer) > 0 {
			answers++
		} else {
			failures++
		}
	}
	answersChannel <- answers
	failuresChannel <- failures
}

func TestDNSLoad(t *testing.T) {
	var trafficRouterDetails = GetTrafficRouterDetails(t)

	benchmark := DNSBenchmark{Benchmark{
		RequestsPerSecondThreshold: *DNSRequestsPerSecondThreshold,
		BenchmarkSeconds:           *BenchTime,
		ThreadCount:                *ThreadCount,
		TrafficRouters:             trafficRouterDetails,
		CoverageZoneLocation:       *CoverageZoneLocation,
	}}
	RunDNSBenchmarksAgainstTrafficRouters(t, benchmark)
}
