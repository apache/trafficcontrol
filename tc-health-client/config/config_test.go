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
	"os"
	"strconv"
	"testing"

	"github.com/apache/trafficcontrol/v8/tc-health-client/util"
)

const (
	test_config_file = "test_files/tc-health-client.json"
)

func TestLoadConfig(t *testing.T) {
	cf := util.ConfigFile{
		Filename:       test_config_file,
		LastModifyTime: 0,
	}

	cfg := Cfg{
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

	ReadCredentials(&cfg, false)

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

	expect = "true"
	monitorStrategisPeers := cfg.MonitorStrategiesPeers
	if expect != strconv.FormatBool(monitorStrategisPeers) {
		t.Fatalf("expected '%s', got %v\n", expect, monitorStrategisPeers)
	}
}

func TestGetCredentialsFromFile(t *testing.T) {

	fi, err := os.CreateTemp("", "creds")
	if err != nil {
		t.Fatalf("creating temp credentials file: %v", err)
	}
	defer os.Remove(fi.Name())

	credsFileContents := `
# credentials
export TO_URL="https://trafficops.example.net"
export TO_USER="myuser"
export TO_PASS="mypass"
`
	expectedURL := `https://trafficops.example.net`
	expectedUser := `myuser`
	expectedPass := `mypass`

	if _, err := fi.Write([]byte(credsFileContents)); err != nil {
		t.Fatalf("writing temp credentials file: %v", err)
	}

	if err := fi.Close(); err != nil {
		t.Fatalf("closing temp credentials file: %v", err)
	}

	credsFilePath := fi.Name()
	toURL, toUser, toPass, err := getCredentialsFromFile(credsFilePath)
	if err != nil {
		t.Fatalf("getting temp credentials file: %v", err)
	}
	if toURL != expectedURL {
		t.Errorf("credentials file TO URL expected '%v' actual '%v'", expectedURL, toURL)
	}
	if toUser != expectedUser {
		t.Errorf("credentials file TO User expected '%v' actual '%v'", expectedUser, toUser)
	}
	if toPass != expectedPass {
		t.Errorf("credentials file TO Pass expected '%v' actual '%v'", expectedPass, toPass)
	}

}
