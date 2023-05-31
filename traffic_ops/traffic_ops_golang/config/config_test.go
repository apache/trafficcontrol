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
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const (
	logError   = "/var/log/traffic_ops/error.log"
	logWarning = "/var/log/traffic_ops/warning.log"
	logInfo    = "/var/log/traffic_ops/info.log"
	logDebug   = "/var/log/traffic_ops/debug.log"
	logEvent   = "/var/log/traffic_ops/event.log"
)

var debugLogging = flag.Bool("debug", false, "enable debug logging in test")

var cfg = Config{
	URL: nil,
	ConfigTrafficOpsGolang: ConfigTrafficOpsGolang{
		LogLocationError:   logError,
		LogLocationWarning: logWarning,
		LogLocationInfo:    logInfo,
		LogLocationDebug:   logDebug,
		LogLocationEvent:   logEvent,
	},
	DB:      ConfigDatabase{},
	Secrets: []string{},
}

func TestLogLocation(t *testing.T) {
	if cfg.ErrorLog() != logError {
		t.Error("ErrorLog should be ", logError)
	}
	if cfg.WarningLog() != logWarning {
		t.Error("WarningLog should be ", logWarning)
	}
	if cfg.InfoLog() != logInfo {
		t.Error("InfoLog should be ", logInfo)
	}
	if cfg.DebugLog() != logDebug {
		t.Error("DebugLog should be ", logDebug)
	}
	if cfg.EventLog() != logEvent {
		t.Error("EventLog should be ", logEvent)
	}
}

func tempFileWith(content []byte) (string, error) {
	tmpfile, err := ioutil.TempFile("", "badcdn.conf")
	if err != nil {
		return "", err
	}
	if _, err := tmpfile.Write(content); err != nil {
		return "", err
	}
	if err := tmpfile.Close(); err != nil {
		return "", err
	}
	return tmpfile.Name(), nil
}

const (
	goodConfig = `
{
	"user_cache_refresh_interval_sec": 30,
	"server_update_status_cache_refresh_interval_sec": 15,
	"disable_auto_cert_deletion": true,
	"traffic_ops_golang" : {
		"cert" : "/etc/pki/tls/certs/localhost.crt",
		"key" : "/etc/pki/tls/private/localhost.key",
		"port" : "443",
		"proxy_timeout" : 60,
		"proxy_keep_alive" : 60,
		"proxy_tls_timeout" : 60,
		"proxy_read_header_timeout" : 60,
		"read_timeout" : 60,
		"read_header_timeout" : 60,
		"write_timeout" : 60,
		"idle_timeout" : 60,
		"routing_blacklist": {
			"ignore_unknown_routes": true,
			"disabled_routes": [4, 5, 6]
		},
		"traffic_vault_backend": "something",
		"traffic_vault_config": {
			"foo": "bar"
		},
		"log_location_error": "stderr",
		"log_location_warning": "stdout",
		"log_location_info": "stdout",
		"log_location_debug": "stdout",
		"log_location_event": "access.log"
	},
	"cors" : {
		"access_control_allow_origin" : "*"
	},
	"to" : {
		"base_url" : "http://localhost:3000",
		"email_from" : "no-reply@traffic-ops-domain.com",
		"no_account_found_msg" : "A Traffic Ops user account is required for access. Please contact your Traffic Ops user administrator."
	},
	"portal" : {
		"base_url" : "http://localhost:8080/!#/",
		"email_from" : "no-reply@traffic-portal-domain.com",
		"pass_reset_path" : "user",
		"user_register_path" : "user"
	},
	"secrets" : [
		"mONKEYDOmONKEYSEE."
	],
	"inactivity_timeout" : 60,
	"client_cert_auth" : {
		"root_certs_dir" : "/etc/pki/tls/certs/"
	}
}
`

	goodDbConfig = `
{
	"description": "Local PostgreSQL database on port 5432",
	"dbname": "traffic_ops",
	"hostname": "localhost",
	"user": "traffic_ops",
	"password": "password",
	"port": "5432",
	"type": "Pg"
}
`
)

func TestLoadConfig(t *testing.T) {
	var err error
	var exp string
	version := "Test Version"

	// set up config paths
	badPath := "/invalid-path/no-file-exists-here"
	badCfg, err := tempFileWith([]byte("no way this is valid json..."))
	if err != nil {
		t.Errorf("cannot create temp file: %v", err)
	}
	defer os.Remove(badCfg) // clean up

	goodCfg, err := tempFileWith([]byte(goodConfig))
	if err != nil {
		t.Errorf("cannot create temp file: %v", err)
	}
	defer os.Remove(goodCfg) // clean up

	goodDbCfg, err := tempFileWith([]byte(goodDbConfig))
	if err != nil {
		t.Errorf("cannot create temp file: %v", err)
	}
	defer os.Remove(goodDbCfg) // clean up

	// test bad paths
	_, errs, blockStartup := LoadConfig(badPath, badPath, version)
	exp = fmt.Sprintf("Loading cdn config from '%s'", badPath)
	if !strings.HasPrefix(errs[0].Error(), exp) {
		t.Error("expected", exp, "got", errs[0].Error())
	}
	if blockStartup != true {
		t.Error("expected blockStartup to be true but it was ", blockStartup)
	}

	// bad json in cdn.conf
	_, errs, blockStartup = LoadConfig(badCfg, badCfg, version)
	exp = fmt.Sprintf("Loading cdn config from '%s': unmarshalling '%s'", badCfg, badCfg)
	if !strings.HasPrefix(errs[0].Error(), exp) {
		t.Error("expected", exp, "got", errs[0].Error())
	}
	if blockStartup != true {
		t.Error("expected blockStartup to be true but it was ", blockStartup)
	}

	// good cdn.conf, bad db conf
	_, errs, blockStartup = LoadConfig(goodCfg, badPath, version)
	exp = fmt.Sprintf("reading db conf '%s'", badPath)
	if !strings.HasPrefix(errs[0].Error(), exp) {
		t.Error("expected", exp, "got", errs[0].Error())
	}
	if blockStartup != true {
		t.Error("expected blockStartup to be true but it was ", blockStartup)
	}

	// good cdn.conf,  bad json in database.conf
	_, errs, blockStartup = LoadConfig(goodCfg, badCfg, version)
	exp = fmt.Sprintf("unmarshalling '%s'", badCfg)
	if !strings.HasPrefix(errs[0].Error(), exp) {
		t.Error("expected", exp, "got", errs[0].Error())
	}
	if blockStartup != true {
		t.Error("expected blockStartup to be true but it was ", blockStartup)
	}

	// good cdn.conf,  good database.conf
	cfg, errs, blockStartup = LoadConfig(goodCfg, goodDbCfg, version)
	if len(errs) != 0 {
		t.Error("Good config -- unexpected errors: ", errs)
	}
	if blockStartup != false {
		t.Error("expected blockStartup to be false but it was ", blockStartup)
	}
	if !cfg.DisableAutoCertDeletion {
		t.Errorf("expected disable_auto_cert_deletion to be true, actual: false")
	}

	if cfg.TrafficVaultBackend != "something" {
		t.Errorf("expected traffic_vault_backend to be 'something', actual: '%s'", cfg.TrafficVaultBackend)
	}
	if cfg.UserCacheRefreshIntervalSec != 30 {
		t.Errorf("expected user_cache_refresh_interval_sec: 30, actual: %d", cfg.UserCacheRefreshIntervalSec)
	}
	if cfg.ServerUpdateStatusCacheRefreshIntervalSec != 15 {
		t.Errorf("expected server_update_status_cache_refresh_interval_sec: 15, actual: %d", cfg.ServerUpdateStatusCacheRefreshIntervalSec)
	}
	tvConfig := make(map[string]string)
	err = json.Unmarshal(cfg.TrafficVaultConfig, &tvConfig)
	if err != nil {
		t.Errorf("unmarshalling traffic_vault_config - expected: no error, actual: %s", err.Error())
	}
	if tvConfig["foo"] != "bar" {
		t.Errorf("unmarshalling traffic_vault_config - expected: foo = bar, actual: foo = %s", tvConfig["foo"])
	}

	if *debugLogging {
		fmt.Printf("Cfg: %+v\n", cfg)
	}

	if cfg.CertPath != "/etc/pki/tls/certs/localhost.crt" {
		t.Error("expected CertPath() == /etc/pki/tls/private/localhost.crt")
	}

	if cfg.KeyPath != "/etc/pki/tls/private/localhost.key" {
		t.Error("expected KeyPath() == /etc/pki/tls/private/localhost.key")
	}
}

func TestValidateRoutingBlacklist(t *testing.T) {
	type testCase struct {
		Input     RoutingBlacklist
		ExpectErr bool
	}
	testCases := []testCase{
		{
			Input: RoutingBlacklist{
				DisabledRoutes: nil,
			},
			ExpectErr: false,
		},
		{
			Input: RoutingBlacklist{
				DisabledRoutes: []int{4, 5, 6},
			},
			ExpectErr: false,
		},
		{
			Input: RoutingBlacklist{
				DisabledRoutes: []int{4, 5, 6},
			},
			ExpectErr: false,
		},
		{
			Input: RoutingBlacklist{
				DisabledRoutes: nil,
			},
			ExpectErr: false,
		},
		{
			Input: RoutingBlacklist{
				DisabledRoutes: []int{2, 2, 4},
			},
			ExpectErr: true,
		},
		{
			Input: RoutingBlacklist{
				DisabledRoutes: []int{4, 4, 6},
			},
			ExpectErr: true,
		},
	}
	for _, tc := range testCases {
		if err := ValidateRoutingBlacklist(tc.Input); err != nil && !tc.ExpectErr {
			t.Errorf("Expected: no error, actual: %v", err)
		} else if err == nil && tc.ExpectErr {
			t.Errorf("Expected: non-nil error, actual: nil")
		}
	}
}
