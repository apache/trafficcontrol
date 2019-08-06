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

func TestMakeCacheDotConfig(t *testing.T) {
	profileName := "myProfile"
	toolName := "myToolName"
	toURL := "https://myto.example.net"
	originFQDN0 := "my.fqdn.example.net"
	originFQDN1 := "my.other.fqdn.example.net"
	originFQDNNoCache := "nocache-fqn.example.net"
	profileDSes := []ProfileDS{
		ProfileDS{Type: tc.DSTypeHTTP, OriginFQDN: &originFQDN0},
		ProfileDS{Type: tc.DSTypeDNS, OriginFQDN: &originFQDN1},
		ProfileDS{Type: tc.DSTypeHTTPNoCache, OriginFQDN: &originFQDNNoCache},
	}

	txt := MakeCacheDotConfig(profileName, profileDSes, toolName, toURL)

	testComment(t, txt, profileName, toolName, toURL)

	if strings.Contains(txt, "my.fqdn.example.net") {
		t.Errorf("expected cached DS type 'my.fqdn.example.net' omitted, actual: '%v'", txt)
	}
	if strings.Contains(txt, "my.other.fqdn.example.net") {
		t.Errorf("expected cached DS type 'my.fqdn.example.net' omitted, actual: '%v'", txt)
	}
	if strings.Contains(txt, "nocache-fqn-should-not-exist.example.net") {
		t.Errorf("expected config include NoCache DS origin 'nocache-fqn-should-not-exist.example.net', actual: '%v'", txt)
	}

}
