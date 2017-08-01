package main

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
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
)

type Config struct {
	HTTPPort           string   `json:"port"`
	DBUser             string   `json:"db_user"`
	DBPass             string   `json:"db_pass"`
	DBServer           string   `json:"db_server"`
	DBDB               string   `json:"db_name"`
	DBSSL              bool     `json:"db_ssl"`
	TOSecret           string   `json:"to_secret"`
	TOURLStr           string   `json:"to_url"`
	TOURL              *url.URL `json:"-"`
	Insecure           bool     `json:"insecure"`
	CertPath           string   `json:"cert_path"`
	KeyPath            string   `json:"key_path"`
	LogLocationError   string   `json:"log_location_error"`
	LogLocationWarning string   `json:"log_location_warning"`
	LogLocationInfo    string   `json:"log_location_info"`
	LogLocationDebug   string   `json:"log_location_debug"`
	LogLocationEvent   string   `json:"log_location_event"`
}

func (c Config) ErrorLog() log.LogLocation   { return log.LogLocation(c.LogLocationError) }
func (c Config) WarningLog() log.LogLocation { return log.LogLocation(c.LogLocationWarning) }
func (c Config) InfoLog() log.LogLocation    { return log.LogLocation(c.LogLocationInfo) }
func (c Config) DebugLog() log.LogLocation   { return log.LogLocation(c.LogLocationDebug) }
func (c Config) EventLog() log.LogLocation   { return log.LogLocation(c.LogLocationEvent) }

func LoadConfig(fileName string) (Config, error) {
	if fileName == "" {
		return Config{}, fmt.Errorf("no filename")
	}

	configBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{}
	if err := json.Unmarshal(configBytes, &cfg); err != nil {
		return Config{}, err
	}

	if cfg, err = ParseConfig(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// ParseConfig validates required fields, and parses non-JSON types
func ParseConfig(cfg Config) (Config, error) {
	missings := ""
	if cfg.HTTPPort == "" {
		missings += "port, "
	}
	if cfg.DBUser == "" {
		missings += "db_user, "
	}
	if cfg.DBPass == "" {
		missings += "db_pass, "
	}
	if cfg.DBServer == "" {
		missings += "db_server, "
	}
	if cfg.DBDB == "" {
		missings += "db_name, "
	}
	if cfg.TOSecret == "" {
		missings += "to_secret, "
	}
	if cfg.CertPath == "" {
		missings += "cert_path, "
	}
	if cfg.KeyPath == "" {
		missings += "key_path, "
	}

	if cfg.LogLocationError == "" {
		cfg.LogLocationError = log.LogLocationNull
	}
	if cfg.LogLocationWarning == "" {
		cfg.LogLocationWarning = log.LogLocationNull
	}
	if cfg.LogLocationInfo == "" {
		cfg.LogLocationInfo = log.LogLocationNull
	}
	if cfg.LogLocationDebug == "" {
		cfg.LogLocationDebug = log.LogLocationNull
	}
	if cfg.LogLocationEvent == "" {
		cfg.LogLocationEvent = log.LogLocationNull
	}

	invalidTOURLStr := ""
	var err error
	if cfg.TOURL, err = url.Parse(cfg.TOURLStr); err != nil {
		invalidTOURLStr = fmt.Sprintf("invalid Traffic Ops URL '%v': %v", cfg.TOURLStr, err)
	}

	if len(missings) > 0 {
		missings = "missing fields: " + missings[:len(missings)-2] // strip final `, `
	}

	errStr := missings
	if errStr != "" && invalidTOURLStr != "" {
		errStr += "; "
	}
	errStr += invalidTOURLStr
	if errStr != "" {
		return Config{}, fmt.Errorf(errStr)
	}

	return cfg, nil
}
