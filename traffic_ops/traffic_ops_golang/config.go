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

	"crypto/tls"
	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/basho/riak-go-client"
)

// Config reflects the structure of the cdn.conf file
type Config struct {
	URL                    *url.URL `json:"-"`
	ConfigHypnotoad        `json:"hypnotoad"`
	ConfigTrafficOpsGolang `json:"traffic_ops_golang"`
	DB                     ConfigDatabase `json:"db"`
	Secrets                []string       `json:"secrets"`
	// NOTE: don't care about any other fields for now..
	RiakAuthOptions *riak.AuthOptions
}

// ConfigHypnotoad carries http setting for hypnotoad (mojolicious) server
type ConfigHypnotoad struct {
	Listen []string `json:"listen"`
	// NOTE: don't care about any other fields for now..
}

// ConfigTrafficOpsGolang carries settings specific to traffic_ops_golang server
type ConfigTrafficOpsGolang struct {
	Port                   string `json:"port"`
	ProxyTimeout           int    `json:"proxy_timeout"`
	ProxyKeepAlive         int    `json:"proxy_keep_alive"`
	ProxyTLSTimeout        int    `json:"proxy_tls_timeout"`
	ProxyReadHeaderTimeout int    `json:"proxy_read_header_timeout"`
	ReadTimeout            int    `json:"read_timeout"`
	ReadHeaderTimeout      int    `json:"read_header_timeout"`
	WriteTimeout           int    `json:"write_timeout"`
	IdleTimeout            int    `json:"idle_timeout"`
	LogLocationError       string `json:"log_location_error"`
	LogLocationWarning     string `json:"log_location_warning"`
	LogLocationInfo        string `json:"log_location_info"`
	LogLocationDebug       string `json:"log_location_debug"`
	LogLocationEvent       string `json:"log_location_event"`
	Insecure               bool   `json:"insecure"`
	MaxDBConnections       int    `json:"max_db_connections"`
}

// ConfigDatabase reflects the structure of the database.conf file
type ConfigDatabase struct {
	Description string `json:"description"`
	DBName      string `json:"dbname"`
	Hostname    string `json:"hostname"`
	User        string `json:"user"`
	Password    string `json:"password"`
	Port        string `json:"port"`
	Type        string `json:"type"`
	SSL         bool   `json:"ssl"`
}

// ErrorLog - critical messages
func (c Config) ErrorLog() log.LogLocation {
	return log.LogLocation(c.LogLocationError)
}

// WarningLog - warning messages
func (c Config) WarningLog() log.LogLocation {
	return log.LogLocation(c.LogLocationWarning)
}

// InfoLog - information messages
func (c Config) InfoLog() log.LogLocation { return log.LogLocation(c.LogLocationInfo) }

// DebugLog - troubleshooting messages
func (c Config) DebugLog() log.LogLocation {
	return log.LogLocation(c.LogLocationDebug)
}

// EventLog - access.log high level transactions
func (c Config) EventLog() log.LogLocation {
	return log.LogLocation(c.LogLocationEvent)
}

// LoadConfig - reads the config file into the Config struct
func LoadConfig(cdnConfPath string, dbConfPath string, riakConfPath string) (Config, error) {
	// load json from cdn.conf
	confBytes, err := ioutil.ReadFile(cdnConfPath)
	if err != nil {
		return Config{}, fmt.Errorf("reading CDN conf '%s': %v", cdnConfPath, err)
	}

	var cfg Config
	err = json.Unmarshal(confBytes, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("unmarshalling '%s': %v", cdnConfPath, err)
	}

	// load json from database.conf
	dbConfBytes, err := ioutil.ReadFile(dbConfPath)
	if err != nil {
		return Config{}, fmt.Errorf("reading db conf '%s': %v", dbConfPath, err)
	}
	err = json.Unmarshal(dbConfBytes, &cfg.DB)
	if err != nil {
		return Config{}, fmt.Errorf("unmarshalling '%s': %v", dbConfPath, err)
	}
	cfg, err = ParseConfig(cfg)

	riakConfBytes, err := ioutil.ReadFile(riakConfPath)
	if err != nil {
		return cfg, fmt.Errorf("reading riak conf '%v': %v", riakConfPath, err)
	}
	riakconf, err := getRiakAuthOptions(string(riakConfBytes))
	if err != nil {
		return cfg, fmt.Errorf("parsing riak conf '%v': %v", riakConfBytes, err)
	}
	cfg.RiakAuthOptions = riakconf

	return cfg, err
}

// CertPath extracts path to cert .cert file
func (c Config) CertPath() string {
	v, ok := c.URL.Query()["cert"]
	if ok {
		return v[0]
	}
	return ""
}

// KeyPath extracts path to cert .key file
func (c Config) KeyPath() string {
	v, ok := c.URL.Query()["key"]
	if ok {
		return v[0]
	}
	return ""
}

func getRiakAuthOptions(s string) (*riak.AuthOptions, error) {
	rconf := &riak.AuthOptions{}
	rconf.TlsConfig = &tls.Config{}
	err := json.Unmarshal([]byte(s), &rconf)
	return rconf, err
}

// ParseConfig validates required fields, and parses non-JSON types
func ParseConfig(cfg Config) (Config, error) {
	missings := ""
	if cfg.Port == "" {
		missings += "port, "
	}
	if len(cfg.Secrets) == 0 {
		missings += "secrets, "
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
	listen := cfg.Listen[0]
	if cfg.URL, err = url.Parse(listen); err != nil {
		invalidTOURLStr = fmt.Sprintf("invalid Traffic Ops URL '%s': %v", listen, err)
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
