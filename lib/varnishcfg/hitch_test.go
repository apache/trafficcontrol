package varnishcfg

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

	"github.com/apache/trafficcontrol/v8/lib/go-atscfg"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
)

func TestGetHitchConfig(t *testing.T) {
	ds1 := &atscfg.DeliveryService{}
	ds1.XMLID = "ds1"
	ds1.Protocol = util.Ptr(1)
	ds1Type := "HTTP"
	ds1.Type = &ds1Type
	ds1.ExampleURLs = []string{"https://ds1.example.org"}
	deliveryServices := []atscfg.DeliveryService{*ds1}
	txt, warnings := GetHitchConfig(deliveryServices, "/ssl")
	expectedTxt := strings.Join([]string{
		`frontend = {`,
		`	host = "*"`,
		`	port = "443"`,
		`}`,
		`backend = "[127.0.0.1]:6081"`,
		`write-proxy-v2 = on`,
		`user = "root"`,
		`pem-file = {`,
		`	cert = "/ssl/ds1_example_org_cert.cer"`,
		`	private-key = "/ssl/ds1.example.org.key"`,
		`}`,
	}, "\n")
	expectedTxt += "\n"
	if len(warnings) != 0 {
		t.Errorf("expected no warnings got %v", warnings)
	}
	if txt != expectedTxt {
		t.Errorf("expected: %s got: %s", expectedTxt, txt)
	}
}
