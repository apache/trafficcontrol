package tmagent

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
	"github.com/apache/trafficcontrol/v6/tc-health-client/config"
	"github.com/apache/trafficcontrol/v6/tc-health-client/util"
	"testing"
)

const (
	test_config_file = "test_files/tc-health-client.json"
)

func TestReadParentDotConfig(t *testing.T) {
	cf := util.ConfigFile{
		Filename:       test_config_file,
		LastModifyTime: 0,
	}

	cfg := config.Cfg{
		TrafficMonitors:        make(map[string]bool, 0),
		HealthClientConfigFile: cf,
	}

	parentDotConfig := util.ConfigFile{
		Filename:       "test_files/etc/parent.config",
		LastModifyTime: 1,
	}

	pi := ParentInfo{
		ParentDotConfig: parentDotConfig,
		Cfg:             cfg,
	}

	_, err := config.LoadConfig(&cfg)
	if err != nil {
		t.Fatalf("failed to load %s: %s\n", test_config_file, err.Error())
	}

	parentStatus := make(map[string]ParentStatus)
	if err := pi.readParentConfig(parentStatus); err != nil {
		t.Fatalf("failed readParentConfig(): %s\n", err.Error())
	}

	numParents := len(parentStatus)
	if numParents != 8 {
		t.Fatalf("failed readParentConfig(): expected 8 parents got %d\n", numParents)
	}
}

func TestReadStrategiesDotYaml(t *testing.T) {
	cf := util.ConfigFile{
		Filename:       test_config_file,
		LastModifyTime: 0,
	}

	cfg := config.Cfg{
		TrafficMonitors:        make(map[string]bool, 0),
		HealthClientConfigFile: cf,
	}

	strategiesDotYaml := util.ConfigFile{
		Filename:       "test_files/etc/strategies.yaml",
		LastModifyTime: 1,
	}

	pi := ParentInfo{
		StrategiesDotYaml: strategiesDotYaml,
		Cfg:               cfg,
	}

	_, err := config.LoadConfig(&cfg)
	if err != nil {
		t.Fatalf("failed to load %s: %s\n", test_config_file, err.Error())
	}

	parentStatus := make(map[string]ParentStatus)
	if err := pi.readStrategies(parentStatus); err != nil {
		t.Fatalf("failed readStrategies(): %s\n", err.Error())
	}

	numParents := len(parentStatus)
	if numParents != 6 {
		t.Fatalf("failed readStrategies(): expected 6 parents got %d\n", numParents)
	}
}

func TestReadHostStatus(t *testing.T) {
	cf := util.ConfigFile{
		Filename:       test_config_file,
		LastModifyTime: 0,
	}

	cfg := config.Cfg{
		TrafficMonitors:        make(map[string]bool, 0),
		HealthClientConfigFile: cf,
	}

	pi := ParentInfo{
		TrafficServerBinDir: "./test_files/bin",
		Cfg:                 cfg,
	}

	_, err := config.LoadConfig(&cfg)
	if err != nil {
		t.Fatalf("failed to load %s: %s\n", test_config_file, err.Error())
	}
	fmt.Println(cfg)

	parentStatus := make(map[string]ParentStatus)
	if err := pi.readHostStatus(parentStatus); err != nil {
		t.Fatalf("failed readHostStatus(): %s\n", err.Error())
	}

	numParents := len(parentStatus)
	if numParents != 14 {
		t.Fatalf("failed readHostStatus(): expected 14 parents got %d\n", numParents)
	}
}

func TestFindATrafficMonitor(t *testing.T) {
	cf := util.ConfigFile{
		Filename:       test_config_file,
		LastModifyTime: 0,
	}

	cfg := config.Cfg{
		TrafficMonitors:        make(map[string]bool, 0),
		HealthClientConfigFile: cf,
	}

	pi := ParentInfo{
		Cfg: cfg,
	}

	_, err := config.LoadConfig(&cfg)
	if err != nil {
		t.Fatalf("failed to load %s: %s\n", test_config_file, err.Error())
	}

	tm, err := pi.findATrafficMonitor()
	if err != nil {
		t.Fatalf("unexpected error calling findATrafficMonitor(): %s\n", err.Error())
	}
	if tm != "tm-01.foo.com" {
		t.Fatalf("expected result: 'tm-01.foo.com', got: %s\n", tm)
	}

	hostName := parseFqdn(tm)
	if hostName != "tm-01" {
		t.Fatalf("expected result: 'tm-01', got: %s\n", hostName)
	}
}

func TestParentStatus(t *testing.T) {
	pstat := ParentStatus{
		Fqdn:         "foo-01.bar.com",
		ActiveReason: true,
		LocalReason:  true,
		ManualReason: false,
	}

	available := pstat.available("manual")
	if available {
		t.Fatalf("expected parent status for %s to be false, got %v\n",
			pstat.Fqdn, available)
	}
	status := pstat.Status()
	if status != "DOWN" {
		t.Fatalf("expected status 'DOWN' got %s\n", status)
	}
	pstat.ManualReason = true
	available = pstat.available("active")
	if !available {
		t.Fatalf("expected parent status for %s to be true, got %v\n",
			pstat.Fqdn, available)
	}
	status = pstat.Status()
	if status != "UP" {
		t.Fatalf("expected status 'UP' got %s\n", status)
	}
	pstat.LocalReason = false
	available = pstat.available("local")
	if available {
		t.Fatalf("expected parent status for %s to be false, got %v\n",
			pstat.Fqdn, available)
	}
	status = pstat.Status()
	if status != "DOWN" {
		t.Fatalf("expected status 'DOWN' got %s\n", status)
	}

	pstat.LocalReason = true
	pstat.ActiveReason = false
	available = pstat.available("active")
	if available {
		t.Fatalf("expected parent status for %s to be false, got %v\n",
			pstat.Fqdn, available)
	}
	status = pstat.Status()
	if status != "DOWN" {
		t.Fatalf("expected status 'DOWN' got %s\n", status)
	}
	pstat.ActiveReason = true
	available = pstat.available("active")
	if !available {
		t.Fatalf("expected parent status for %s to be false, got %v\n",
			pstat.Fqdn, available)
	}
	status = pstat.Status()
	if status != "UP" {
		t.Fatalf("expected status 'UP' got %s\n", status)
	}
}

func TestStatusReason(t *testing.T) {
	var reason StatusReason

	reason = ACTIVE
	if reason.String() != "ACTIVE" {
		t.Fatalf("expected reason 'ACTIVE' got %s\n", reason.String())
	}
	reason = LOCAL
	if reason.String() != "LOCAL" {
		t.Fatalf("expected reason 'LOCAL' got %s\n", reason.String())
	}
	reason = MANUAL
	if reason.String() != "MANUAL" {
		t.Fatalf("expected reason 'MANUAL' got %s\n", reason.String())
	}
	reason = 9
	if reason.String() != "UNDEFINED" {
		t.Fatalf("expected reason 'UNDEFINED' got %s\n", reason.String())
	}
}
