/*
   Copyright 2015 Comcast Cable Communications Management, LLC

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

package test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/client/fixtures"
	"github.com/jheitz200/test_helper"
)

func TestServer(t *testing.T) {
	resp := fixtures.Servers()
	server := testHelper.ValidHTTPServer(resp)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	testHelper.Context(t, "Given the need to test a successful Traffic Ops request for Servers")

	servers, err := to.Servers()
	if err != nil {
		testHelper.Error(t, "Should be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should be able to make a request to Traffic Ops")
	}

	if len(servers) != 3 {
		testHelper.Error(t, "Should get back \"3\" Server, got: %d", len(servers))
	} else {
		testHelper.Success(t, "Should get back \"3\" Server")
	}

	if servers[0].HostName != "edge-alb-01" {
		testHelper.Error(t, "Should get \"edge-alb-01\" for \"HostName\", got: %s", servers[0].HostName)
	} else {
		testHelper.Success(t, "Should get \"edge-alb-01\" for \"HostName\"")
	}

	if servers[0].DomainName != "albuquerque.nm.albuq.kabletown.com" {
		testHelper.Error(t, "Should get \"albuquerque.nm.albuq.kabletown.com\" for \"DomainName\", got: %s", servers[0].DomainName)
	} else {
		testHelper.Success(t, "Should get \"albuquerque.nm.albuq.kabletown.com\" for \"DomainName\"")
	}

	if servers[0].Type != "EDGE" {
		testHelper.Error(t, "Should get \"EDGE\" for \"Type\", got: %s", servers[0].Type)
	} else {
		testHelper.Success(t, "Should get \"EDGE\" for \"Type\"")
	}
}

func TestServersUnauthorized(t *testing.T) {
	server := testHelper.InvalidHTTPServer(http.StatusUnauthorized)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	testHelper.Context(t, "Given the need to test a failed Traffic Ops request for Servers")

	_, err := to.Servers()
	if err == nil {
		testHelper.Error(t, "Should not be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should not be able to make a request to Traffic Ops")
	}
}

func TestServerFQDN(t *testing.T) {
	resp := fixtures.Servers()
	server := testHelper.ValidHTTPServer(resp)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	shortName := "edge-alb-01"
	testHelper.Context(t, "Given the need to test a successful Traffic Ops request for the FQDN of Server: \"%s\"", shortName)

	s, err := to.ServersFqdn("edge-alb-01")
	if err != nil {
		testHelper.Error(t, "Should be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should be able to make a request to Traffic Ops")
	}

	if s != "edge-alb-01.albuquerque.nm.albuq.kabletown.com" {
		testHelper.Error(t, "Should get back \"edge-alb-01.albuquerque.nm.albuq.kabletown.com\", got: %s", s)
	} else {
		testHelper.Success(t, "Should get back \"edge-alb-01.albuquerque.nm.albuq.kabletown.com\"")
	}
}

func TestServerFQDNError(t *testing.T) {
	var resp client.ServerResponse
	server := testHelper.ValidHTTPServer(resp)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	shortName := "edge-alb-01"
	testHelper.Context(t, "Given the need to test a failed Traffic Ops request for the FQDN of Server: \"%s\"", shortName)

	_, err := to.ServersFqdn(shortName)
	if err == nil {
		testHelper.Error(t, "Should not be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should not be able to make a request to Traffic Ops")
	}
}

func TestServerFQDNUnauthorized(t *testing.T) {
	server := testHelper.InvalidHTTPServer(http.StatusUnauthorized)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	shortName := "edge-alb-01"
	testHelper.Context(t, "Given the need to test a failed Traffic Ops request for the FQDN of Server: \"%s\"", shortName)

	_, err := to.ServersFqdn(shortName)
	if err == nil {
		testHelper.Error(t, "Should not be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should not be able to make a request to Traffic Ops")
	}
}

func TestServerShortName(t *testing.T) {
	resp := fixtures.Servers()
	server := testHelper.ValidHTTPServer(resp)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	pattern := "edge"
	testHelper.Context(t, "Given the need to test a successful Traffic Ops request for servers that match Short Name: \"%s\"", pattern)

	servers, err := to.ServersShortNameSearch(pattern)
	if err != nil {
		testHelper.Error(t, "Should be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should be able to make a request to Traffic Ops")
	}

	if servers[0] != "edge-alb-01" {
		testHelper.Error(t, "Should get back \"edge-alb-01\", got: %s", servers[0])
	} else {
		testHelper.Success(t, "Should get back \"edge-alb-01\"")
	}

	if servers[1] != "edge-alb-02" {
		testHelper.Error(t, "Should get back \"edge-alb-02\", got: %s", servers[1])
	} else {
		testHelper.Success(t, "Should get back \"edge-alb-02\"")
	}
}

func TestServerShortNameError(t *testing.T) {
	var resp client.ServerResponse
	server := testHelper.ValidHTTPServer(resp)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	pattern := "edge"
	testHelper.Context(t, "Given the need to test a failed Traffic Ops request for servers that match Short Name: \"%s\"", pattern)

	_, err := to.ServersShortNameSearch(pattern)
	if err == nil {
		testHelper.Error(t, "Should not be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should not be able to make a request to Traffic Ops")
	}
}

func TestServerShortNameUnauthorized(t *testing.T) {
	server := testHelper.InvalidHTTPServer(http.StatusUnauthorized)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	pattern := "edge"
	testHelper.Context(t, "Given the need to test a failed Traffic Ops request for servers that match Short Name: \"%s\"", pattern)

	_, err := to.ServersShortNameSearch(pattern)
	if err == nil {
		testHelper.Error(t, "Should not be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should not be able to make a request to Traffic Ops")
	}
}

func TestServerByType(t *testing.T) {
	resp := fixtures.LogstashServers()
	server := testHelper.ValidHTTPServer(resp)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	testHelper.Context(t, "Given the need to test a successful Traffic Ops request for \"Logstash\" Servers")

	params := make(url.Values)
	params.Add("type", "Logstash")

	servers, err := to.ServersByType(params)
	if err != nil {
		testHelper.Error(t, "Should be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should be able to make a request to Traffic Ops")
	}

	if len(servers) != 2 {
		testHelper.Error(t, "Should get back \"2\" Server, got: %d", len(servers))
	} else {
		testHelper.Success(t, "Should get back \"2\" Server")
	}

	if servers[0].HostName != "logstash-01" {
		testHelper.Error(t, "Should get \"logstash-01\" for \"HostName\", got: %s", servers[0].HostName)
	} else {
		testHelper.Success(t, "Should get \"logstash-01\" for \"HostName\"")
	}

	if servers[0].DomainName != "albuquerque.nm.albuq.kabletown.com" {
		testHelper.Error(t, "Should get \"albuquerque.nm.albuq.kabletown.com\" for \"DomainName\", got: %s", servers[0].DomainName)
	} else {
		testHelper.Success(t, "Should get \"albuquerque.nm.albuq.kabletown.com\" for \"DomainName\"")
	}

	if servers[0].Type != "LOGSTASH" {
		testHelper.Error(t, "Should get \"LOGSTASH\" for \"Type\", got: %s", servers[0].Type)
	} else {
		testHelper.Success(t, "Should get \"LOGSTASH\" for \"Type\"")
	}
}

func TestServerByTypeUnauthorized(t *testing.T) {
	server := testHelper.InvalidHTTPServer(http.StatusUnauthorized)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	testHelper.Context(t, "Given the need to test a failed Traffic Ops request for \"Logstash\" servers")

	params := make(url.Values)
	params.Add("type", "Logstash")

	_, err := to.ServersByType(params)
	if err == nil {
		testHelper.Error(t, "Should not be able to make a request to Traffic Ops")
	} else {
		testHelper.Success(t, "Should not be able to make a request to Traffic Ops")
	}
}
