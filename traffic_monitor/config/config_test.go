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
)

const exampleTMConfig = `
{
	"monitor_config_polling_interval_ms": 5000,
	"http_timeout_ms": 30000,
	"peer_optimistic_quorum_min": 3,
	"max_events": 200,
	"health_flush_interval_ms": 1000,
	"stat_flush_interval_ms": 1000,
	"stat_polling": false,
	"distributed_polling": true,
	"log_location_access": "access.log",
	"log_location_event": "event.log",
	"log_location_error": "error.log",
	"log_location_warning": "warning.log",
	"log_location_info": "info.log",
	"log_location_debug": "debug.log",
	"serve_read_timeout_ms": 10000,
	"serve_write_timeout_ms": 10000,
	"stat_buffer_interval_ms": 20000,
	"short_hostname_override": "foobar",
	"traffic_ops_disk_retry_max": 35,
	"crconfig_backup_file": "crconfig.asdf",
	"tmconfig_backup_file": "tmconfig.asdf",
	"http_polling_format": "thisformatdoesnotexist",
	"static_file_dir": "static/"
}
`

const exampleBadTMConfig = `
{
	"monitor_config_polling_interval_ms": 5000,
	"http_timeout_ms": 30000,
	"peer_optimistic_quorum_min": 3,
	"max_events": 200,
	"health_flush_interval_ms": 1000,
	"stat_flush_interval_ms": 1000,
	"stat_polling": true,
	"distributed_polling": true,
	"log_location_access": "access.log",
	"log_location_event": "event.log",
	"log_location_error": "error.log",
	"log_location_warning": "warning.log",
	"log_location_info": "info.log",
	"log_location_debug": "debug.log",
	"serve_read_timeout_ms": 10000,
	"serve_write_timeout_ms": 10000,
	"stat_buffer_interval_ms": 20000,
	"short_hostname_override": "foobar",
	"traffic_ops_disk_retry_max": 35,
	"crconfig_backup_file": "crconfig.asdf",
	"tmconfig_backup_file": "tmconfig.asdf",
	"http_polling_format": "thisformatdoesnotexist",
	"static_file_dir": "static/"
}
`

func TestConfigLoad(t *testing.T) {
	c, err := LoadBytes([]byte(exampleTMConfig))
	if err != nil {
		t.Fatalf("loading config bytes - expected: no error, actual: %v", err)
	}
	if c.StatPolling != false {
		t.Errorf("StatPolling - expected: false, actual: %t", c.StatPolling)
	}
	if c.DistributedPolling != true {
		t.Errorf("DistributedPolling - expected: true, actual: %t", c.DistributedPolling)
	}
	if string(c.WarningLog()) != c.LogLocationWarning {
		t.Errorf("warning log location - expected: %s, actual: %s\n", c.LogLocationWarning, string(c.WarningLog()))
	}
	if string(c.ErrorLog()) != c.LogLocationError {
		t.Errorf("error log location - expected: %s, actual: %s\n", c.LogLocationError, string(c.ErrorLog()))
	}
	if string(c.EventLog()) != c.LogLocationEvent {
		t.Errorf("event log location - expected: %s, actual: %s\n", c.LogLocationEvent, string(c.EventLog()))
	}
	if string(c.InfoLog()) != c.LogLocationInfo {
		t.Errorf("info log location - expected: %s, actual: %s\n", c.LogLocationInfo, string(c.InfoLog()))
	}
	if string(c.DebugLog()) != c.LogLocationDebug {
		t.Errorf("debug log location - expected: %s, actual: %s\n", c.LogLocationDebug, string(c.DebugLog()))
	}
	if c.ShortHostnameOverride != "foobar" {
		t.Errorf("ShortHostnameOverride - expected: foobar, actual: %s", c.ShortHostnameOverride)
	}
	if c.PeerOptimisticQuorumMin != 3 {
		t.Errorf("PeerOmptimisticQuorumMin - expected: 3, actual: %d", c.PeerOptimisticQuorumMin)
	}
	if c.TrafficOpsDiskRetryMax != 35 {
		t.Errorf("TrafficOpsDiskRetryMax - expected: 35, actual: %d", c.TrafficOpsDiskRetryMax)
	}
	if c.CRConfigBackupFile != "crconfig.asdf" {
		t.Errorf("CRConfigBackupFile - expected: crconfig.asdf, actual: %s", c.CRConfigBackupFile)
	}
	if c.TMConfigBackupFile != "tmconfig.asdf" {
		t.Errorf("TMConfigBackupFile - expected: tmconfig.asdf, actual: %s", c.TMConfigBackupFile)
	}
	if c.HTTPPollingFormat != "thisformatdoesnotexist" {
		t.Errorf("HTTPPollingFormat - expected: thisformatdoesnotexist, actual: %s", c.HTTPPollingFormat)
	}
}

func TestBadConfigLoad(t *testing.T) {
	_, err := LoadBytes([]byte(exampleBadTMConfig))
	if err == nil {
		t.Errorf("loading bad config file (stat_polling and distributed_polling both enabled) -- expected: error, actual: nil")
	}
}

func TestConfigLoadDefaults(t *testing.T) {
	c, err := LoadBytes([]byte(`{}`))
	if err != nil {
		t.Fatalf("loading empty config bytes - expected: no error, actual: %v", err)
	}
	if c.StatPolling != true {
		t.Errorf("StatPolling default - expected: true, actual: %t", c.StatPolling)
	}
	if c.DistributedPolling != false {
		t.Errorf("DistributedPolling default - expected: false, actual: %t", c.DistributedPolling)
	}
}
