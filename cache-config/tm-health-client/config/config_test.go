package config

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
	"testing"

	"github.com/apache/trafficcontrol/cache-config/tm-health-client/util"
)

const (
	test_config_file = "test_files/tm-health-client.json"
)

func TestLoadConfig(t *testing.T) {
	cf := util.ConfigFile{
		Filename:       test_config_file,
		LastModifyTime: 0,
	}

	cfg := Cfg{
		TrafficMonitors:        make(map[string]bool, 0),
		HealthClientConfigFile: cf,
	}

	_, err := LoadConfig(&cfg)
	if err != nil {
		t.Fatalf("failed to load %s: %s\n", test_config_file, err.Error())
	}

	expect := "over-the-top"
	cdn := cfg.CDNName
	if cdn != expect {
		t.Fatalf("expected '%s', got %s\n", expect, cdn)
	}

	expectb := false
	markDowns := cfg.EnableActiveMarkdowns
	if markDowns != expectb {
		t.Fatalf("expected '%v', got %v\n", expectb, markDowns)
	}

	expect = "active"
	reasonCode := cfg.ReasonCode
	if reasonCode != expect {
		t.Fatalf("expected '%s', got %s\n", expect, reasonCode)
	}

	expect = "test_files/credentials"
	creds := cfg.TOCredentialFile
	if creds != expect {
		t.Fatalf("expected '%s', got %s\n", expect, creds)
	}

	ReadCredentials(&cfg)

	expect = "https://tp.cdn.com:443"
	tourl := cfg.TOUrl
	if tourl != expect {
		t.Fatalf("expected '%s', got %s\n", expect, tourl)
	}

	expect = "to_user"
	touser := cfg.TOUser
	if touser != expect {
		t.Fatalf("expected '%s', got %s\n", expect, touser)
	}

	expect = "to_pass"
	topass := cfg.TOPass
	if topass != expect {
		t.Fatalf("expected '%s', got %s\n", expect, topass)
	}

	expect = "15s"
	poll := cfg.TmPollIntervalSeconds
	if poll != expect {
		t.Fatalf("expected '%s', got %s\n", expect, poll)
	}

	expect = "5s"
	rto := cfg.TORequestTimeOutSeconds
	if rto != expect {
		t.Fatalf("expected '%s', got %s\n", expect, rto)
	}

	expect = "./test_files/etc"
	cfgdir := cfg.TrafficServerConfigDir
	if cfgdir != expect {
		t.Fatalf("expected '%s', got %s\n", expect, cfgdir)
	}

	expect = "./test_files/bin"
	bindir := cfg.TrafficServerBinDir
	if bindir != expect {
		t.Fatalf("expected '%s', got %s\n", expect, bindir)
	}
}
