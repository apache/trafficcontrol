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
	"cache_health_polling_interval_ms": 120000,
	"cache_stat_polling_interval_ms": 120000,
	"monitor_config_polling_interval_ms": 5000,
	"http_timeout_ms": 30000,
	"peer_polling_interval_ms": 120000,
	"peer_optimistic": true,
	"peer_optimistic_quorum_min": 0,
	"max_events": 200,
	"max_stat_history": 5,
	"max_health_history": 5,
	"health_flush_interval_ms": 1000,
	"stat_flush_interval_ms": 1000,
	"log_location_event": "event.log",
	"log_location_error": "error.log",
	"log_location_warning": "warning.log",
	"log_location_info": "info.log",
	"log_location_debug": "debug.log",
	"serve_read_timeout_ms": 10000,
	"serve_write_timeout_ms": 10000,
	"stat_buffer_interval_ms": 20000,
	"http_poll_no_sleep": false,
	"static_file_dir": "static/"
}
`

func TestLoggingConfig(t *testing.T) {
	c, err := LoadBytes([]byte(exampleTMConfig))
	if err != nil {
		t.Fatalf("loading config bytes - expected: no error, actual: %v", err)
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
}
