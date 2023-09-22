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
	"path/filepath"
	"testing"

	"github.com/apache/trafficcontrol/v8/tc-health-client/config"
	"github.com/apache/trafficcontrol/v8/tc-health-client/util"
)

const (
	test_config_file = "test_files/tc-health-client.json"
)

func TestReadParentDotConfig(t *testing.T) {
	cf := util.ConfigFile{
		Filename:       test_config_file,
		LastModifyTime: 0,
	}

	cfg := &config.Cfg{
		HealthClientConfigFile: cf,
		TrafficServerConfigDir: "test_files/etc/",
	}
	parents := util.ConfigFile{
		Filename:       filepath.Join(cfg.TrafficServerConfigDir, ParentsFile),
		LastModifyTime: 1,
	}
	strategies := util.ConfigFile{
		Filename:       filepath.Join(cfg.TrafficServerConfigDir, StrategiesFile),
		LastModifyTime: 1,
	}
	pi := ParentInfo{
		ParentDotConfig:        parents,
		StrategiesDotYaml:      strategies,
		TrafficServerBinDir:    cfg.TrafficServerBinDir,
		TrafficServerConfigDir: cfg.TrafficServerConfigDir,
	}

	if _, err := config.LoadConfig(cfg); err != nil {
		t.Fatalf("failed to load %s: %s\n", test_config_file, err.Error())
	}

	if err := pi.readParentConfig(); err != nil {
		t.Fatalf("failed readParentConfig(): %s\n", err.Error())
	}

	numParents := len(pi.GetParents())
	if numParents != 8 {
		t.Fatalf("failed readParentConfig(): expected 8 parents got %d\n", numParents)
	}
}

func TestReadStrategiesDotYaml(t *testing.T) {
	cf := util.ConfigFile{
		Filename:       test_config_file,
		LastModifyTime: 0,
	}

	cfg := &config.Cfg{
		HealthClientConfigFile: cf,
		TrafficServerConfigDir: "test_files/etc/",
	}
	parents := util.ConfigFile{
		Filename:       filepath.Join(cfg.TrafficServerConfigDir, ParentsFile),
		LastModifyTime: 1,
	}
	strategies := util.ConfigFile{
		Filename:       filepath.Join(cfg.TrafficServerConfigDir, StrategiesFile),
		LastModifyTime: 1,
	}

	if _, err := config.LoadConfig(cfg); err != nil {
		t.Fatalf("failed to load %s: %s\n", test_config_file, err.Error())
	}

	t.Run("Read Strategies using config file", func(t *testing.T) {
		pi := ParentInfo{
			ParentDotConfig:        parents,
			StrategiesDotYaml:      strategies,
			TrafficServerBinDir:    cfg.TrafficServerBinDir,
			TrafficServerConfigDir: cfg.TrafficServerConfigDir,
		}
		t.Logf("Monitoring peers value: %v", cfg.MonitorStrategiesPeers)

		if err := pi.readStrategies(cfg.MonitorStrategiesPeers); err != nil {
			t.Fatalf("failed readStrategies(): %s\n", err.Error())
		}

		numParents := len(pi.GetParents())
		if numParents != 2 && !cfg.MonitorStrategiesPeers {
			t.Fatalf("failed readStrategies(): expected 2 parents got %d\n", numParents)
		} else if numParents != 6 && cfg.MonitorStrategiesPeers {
			t.Fatalf("failed readStrategies(): expected 6 parents got %d\n", numParents)
		}
	})

	t.Run("Read Strategies with monitoring peers on", func(t *testing.T) {
		pi := ParentInfo{
			ParentDotConfig:        parents,
			StrategiesDotYaml:      strategies,
			TrafficServerBinDir:    cfg.TrafficServerBinDir,
			TrafficServerConfigDir: cfg.TrafficServerConfigDir,
		}
		if err := pi.readStrategies(true); err != nil {
			t.Fatalf("failed readStrategies(): %s\n", err.Error())
		}

		numParents := len(pi.GetParents())
		if numParents != 6 {
			t.Fatalf("failed readStrategies(): expected 6 parents got %d\n", numParents)
		}
	})

	t.Run("Read Strategies with monitoring peers off", func(t *testing.T) {
		pi := ParentInfo{
			ParentDotConfig:        parents,
			StrategiesDotYaml:      strategies,
			TrafficServerBinDir:    cfg.TrafficServerBinDir,
			TrafficServerConfigDir: cfg.TrafficServerConfigDir,
		}
		if err := pi.readStrategies(false); err != nil {
			t.Fatalf("failed readStrategies(): %s\n", err.Error())
		}

		numParents := len(pi.GetParents())
		if numParents != 2 {
			t.Log(pi.GetParents())
			t.Fatalf("failed readStrategies(): expected 2 parents got %d\n", numParents)
		}
	})
}

func TestReadHostStatus(t *testing.T) {
	cf := util.ConfigFile{
		Filename:       test_config_file,
		LastModifyTime: 0,
	}

	cfg := &config.Cfg{
		HealthClientConfigFile: cf,
		TrafficServerConfigDir: "test_files/etc/",
	}
	parents := util.ConfigFile{
		Filename:       filepath.Join(cfg.TrafficServerConfigDir, ParentsFile),
		LastModifyTime: 1,
	}
	strategies := util.ConfigFile{
		Filename:       filepath.Join(cfg.TrafficServerConfigDir, StrategiesFile),
		LastModifyTime: 1,
	}
	pi := ParentInfo{
		ParentDotConfig:        parents,
		StrategiesDotYaml:      strategies,
		TrafficServerBinDir:    cfg.TrafficServerBinDir,
		TrafficServerConfigDir: cfg.TrafficServerConfigDir,
	}

	if _, err := config.LoadConfig(cfg); err != nil {
		t.Fatalf("failed to load %s: %s\n", test_config_file, err.Error())
	}
	fmt.Println(cfg)

	if err := pi.readHostStatus(cfg); err != nil {
		t.Fatalf("failed readHostStatus(): %s\n", err.Error())
	}

	numParents := len(pi.GetParents())
	if numParents != 14 {
		t.Fatalf("failed readHostStatus(): expected 14 parents got %d\n", numParents)
	}
}

func TestFindATrafficMonitor(t *testing.T) {
	cf := util.ConfigFile{
		Filename:       test_config_file,
		LastModifyTime: 0,
	}

	cfg := &config.Cfg{
		HealthClientConfigFile: cf,
		TrafficServerConfigDir: "test_files/etc/",
	}
	cfgPtr := config.NewCfgPtr(cfg)

	pi, err := NewParentInfo(cfgPtr)
	if err != nil {
		t.Fatalf("failed to create parent info: %s\n", err.Error())
	}

	if _, err := config.LoadConfig(cfg); err != nil {
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
