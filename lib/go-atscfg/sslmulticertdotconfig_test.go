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

	"github.com/apache/trafficcontrol/lib/go-util"
)

func TestMakeSSLMultiCertDotConfig(t *testing.T) {
	hdr := "myHeaderComment"

	server := makeGenericServer()
	server.CDNName = util.Ptr("mycdn")

	makeDS := func(name string, dsType *string, protocol int, exampleURL string) DeliveryService {
		ds := makeGenericDS()
		ds.XMLID = name
		ds.Type = dsType
		ds.Protocol = util.Ptr(protocol)
		ds.ExampleURLs = []string{exampleURL}
		return *ds
	}

	dses := []DeliveryService{
		makeDS("my-https-ds", util.Ptr("HTTP"), 1 /* https */, "https://my-https-ds.example.net"),
		makeDS("my-https-and-http-ds", util.Ptr("HTTP"), 2 /* https and http */, "https://my-https-and-http-ds.example.net"),
		makeDS("my-https-to-http-ds", util.Ptr("HTTP"), 3 /* https to http */, "https://my-https-to-http-ds.example.net"),
	}

	cfg, err := MakeSSLMultiCertDotConfig(server, dses, &SSLMultiCertDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	if !strings.Contains(txt, hdr) {
		t.Errorf("expected: header comment text '" + hdr + "', actual: missing")
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
	hdr := "myHeaderComment"

	server := makeGenericServer()
	server.CDNName = util.Ptr("mycdn")

	makeDS := func(name string, dsType *string, protocol int, exampleURL string) DeliveryService {
		ds := makeGenericDS()
		ds.XMLID = "name"
		ds.Type = dsType
		ds.Protocol = util.IntPtr(protocol)
		ds.ExampleURLs = []string{exampleURL}
		return *ds
	}

	// http-only DSes have no SSL, should be excluded
	dses := []DeliveryService{
		makeDS("myds", util.Ptr("HTTP"), 0 /* http */, "https://myds.example.net"),
	}

	cfg, err := MakeSSLMultiCertDotConfig(server, dses, &SSLMultiCertDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	if !strings.Contains(txt, hdr) {
		t.Errorf("expected: header comment text '" + hdr + "', actual: missing")
	}
	if !strings.HasPrefix(strings.TrimSpace(txt), "#") {
		t.Errorf("expected: header comment, actual: missing")
	}

	if strings.Contains(txt, "myds") {
		t.Errorf("expected HTTP-only DS to be excluded, actual '%v'", txt)
	}
}
