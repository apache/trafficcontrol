package atscfg

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
	"strings"
	"testing"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

func TestMakeSSLMultiCertDotConfig(t *testing.T) {
	cdnName := tc.CDNName("mycdn")
	toToolName := "my-to"
	toURL := "my-to.example.net"

	dses := map[tc.DeliveryServiceName]SSLMultiCertDS{
		"my-https-ds": SSLMultiCertDS{
			Type:        tc.DSTypeHTTP,
			Protocol:    1, // https and http
			ExampleURLs: []string{"https://my-https-ds.example.net"},
		},
		"my-https-and-http-ds": SSLMultiCertDS{
			Type:        tc.DSTypeHTTP,
			Protocol:    2, // https and http
			ExampleURLs: []string{"https://my-https-and-http-ds.example.net"},
		},
		"my-https-to-http-ds": SSLMultiCertDS{
			Type:        tc.DSTypeHTTP,
			Protocol:    3, // https to http
			ExampleURLs: []string{"https://my-https-to-http-ds.example.net"},
		},
	}

	txt := MakeSSLMultiCertDotConfig(cdnName, toToolName, toURL, dses)

	if !strings.Contains(txt, string(cdnName)) {
		t.Errorf("expected: cdnName '" + string(cdnName) + "', actual: missing")
	}
	if !strings.Contains(txt, toToolName) {
		t.Errorf("expected: toToolName '" + toToolName + "', actual: missing")
	}
	if !strings.Contains(txt, toURL) {
		t.Errorf("expected: toURL '" + toURL + "', actual: missing")
	}
	if !strings.HasPrefix(strings.TrimSpace(txt), "#") {
		t.Errorf("expected: header comment, actual: missing")
	}

	if !strings.Contains(txt, "my-https-ds") {
		t.Errorf("expected HTTPS DS to be included, actual '%v'", txt)
	}
	if !strings.Contains(txt, "my-https-and-http-ds") {
		t.Errorf("expected HTTPS DS to be included, actual '%v'", txt)
	}
	if !strings.Contains(txt, "my-https-to-http-ds") {
		t.Errorf("expected HTTPS DS to be included, actual '%v'", txt)
	}
}

func TestMakeSSLMultiCertDotConfigHTTPDeliveryService(t *testing.T) {
	cdnName := tc.CDNName("mycdn")
	toToolName := "my-to"
	toURL := "my-to.example.net"

	dses := map[tc.DeliveryServiceName]SSLMultiCertDS{
		"myds": SSLMultiCertDS{
			Type:        tc.DSTypeHTTP,
			Protocol:    0, // http-only DSes have no SSL, should be excluded
			ExampleURLs: []string{"https://myds.example.net"},
		},
	}

	txt := MakeSSLMultiCertDotConfig(cdnName, toToolName, toURL, dses)

	if !strings.Contains(txt, string(cdnName)) {
		t.Errorf("expected: cdnName '" + string(cdnName) + "', actual: missing")
	}
	if !strings.Contains(txt, toToolName) {
		t.Errorf("expected: toToolName '" + toToolName + "', actual: missing")
	}
	if !strings.Contains(txt, toURL) {
		t.Errorf("expected: toURL '" + toURL + "', actual: missing")
	}
	if !strings.HasPrefix(strings.TrimSpace(txt), "#") {
		t.Errorf("expected: header comment, actual: missing")
	}

	if strings.Contains(txt, "myds") {
		t.Errorf("expected HTTP-only DS to be excluded, actual '%v'", txt)
	}
}
