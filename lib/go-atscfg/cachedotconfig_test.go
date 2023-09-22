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

	"github.com/apache/trafficcontrol/v8/lib/go-util"
)

func TestMakeCacheDotConfig(t *testing.T) {
	server := makeGenericServer()
	serverProfile := "myProfile"
	server.Profiles = []string{serverProfile}
	servers := []Server{*server}

	ds0 := makeGenericDS()
	ds0.ID = util.Ptr(420)
	ds0.XMLID = "ds0"
	ds0.OrgServerFQDN = util.Ptr("http://my.fqdn.example.net")
	ds0Type := "HTTP"
	ds0.Type = &ds0Type

	ds1 := makeGenericDS()
	ds1.ID = util.Ptr(421)
	ds1.XMLID = "ds1"
	ds1.OrgServerFQDN = util.Ptr("http://my.other.fqdn.example.net")
	ds1Type := "DNS"
	ds1.Type = &ds1Type

	ds2 := makeGenericDS()
	ds2.ID = util.Ptr(422)
	ds2.XMLID = "ds2"
	ds2.OrgServerFQDN = util.Ptr("http://nocache-fqn.example.net")
	ds2Type := "HTTP_NO_CACHE"
	ds2.Type = &ds2Type

	dses := []DeliveryService{*ds0, *ds1, *ds2}

	dss := makeDSS(servers, dses)

	hdr := "myHeaderComment"

	cfg, err := MakeCacheDotConfig(server, servers, dses, dss, &CacheDotConfigOpts{HdrComment: hdr})
	if err != nil {
		t.Fatal(err)
	}
	txt := cfg.Text

	testComment(t, txt, hdr)

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
