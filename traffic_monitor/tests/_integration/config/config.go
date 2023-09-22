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
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/kelseyhightower/envconfig"
)

// Config reflects the structure of the test-to-api.conf file
type Config struct {
	TrafficMonitor TrafficMonitor `json:"trafficMonitor"`
	Default        Default        `json:"default"`
}

// TrafficMonitor is the monitor config section.
type TrafficMonitor struct {
	// URL points to the Traffic Monitor instance being tested
	URL string `json:"url" envconfig:"TM_URL"`
}

// Default - config section
type Default struct {
	Session Session   `json:"session"`
	Log     Locations `json:"logLocations"`
}

// Session - config section
type Session struct {
	TimeoutInSecs int `json:"timeoutInSecs" envconfig:"SESSION_TIMEOUT_IN_SECS"`
}

// Locations - reflects the structure of the database.conf file
type Locations struct {
	Debug   string `json:"debug"`
	Event   string `json:"event"`
	Error   string `json:"error"`
	Info    string `json:"info"`
	Warning string `json:"warning"`
}

// LoadConfig - reads the config file into the Config struct
func LoadConfig(confPath string) (Config, error) {
	var cfg Config

	if _, err := os.Stat(confPath); !os.IsNotExist(err) {
		confBytes, err := ioutil.ReadFile(confPath)
		if err != nil {
			return Config{}, fmt.Errorf("reading CDN conf '%s': %v", confPath, err)
		}

		err = json.Unmarshal(confBytes, &cfg)
		if err != nil {
			return Config{}, fmt.Errorf("unmarshalling '%s': %v", confPath, err)
		}
	}
	errs := validate(confPath, cfg)
	if len(errs) > 0 {
		fmt.Printf("configuration error:\n")
		for _, e := range errs {
			fmt.Printf("%v\n", e)
		}
		os.Exit(0)
	}
	err := envconfig.Process("traffic-ops-client-tests", &cfg)
	if err != nil {
		log.Errorln(fmt.Errorf("cannot parse config: %v\n", err))
		os.Exit(0)
	}

	return cfg, err
}

// validate all required fields in the config.
func validate(confPath string, config Config) []error {

	errs := []error{}

	var f string
	f = "TrafficMonitor"
	toTag, ok := getStructTag(config, f)
	if !ok {
		errs = append(errs, fmt.Errorf("'%s' must be configured in %s", toTag, confPath))
	}

	if config.TrafficMonitor.URL == "" {
		f = "URL"
		tag, ok := getStructTag(config.TrafficMonitor, f)
		if !ok {
			errs = append(errs, fmt.Errorf("cannot lookup structTag: %s", f))
		}
		errs = append(errs, fmt.Errorf("'%s.%s' must be configured in %s", toTag, tag, confPath))
	}

	return errs
}

func getStructTag(thing interface{}, fieldName string) (string, bool) {
	var tag string
	var ok bool
	t := reflect.TypeOf(thing)
	if t != nil {
		if f, ok := t.FieldByName(fieldName); ok {
			tag = f.Tag.Get("json")
			return tag, ok
		}
	}
	return tag, ok
}

// ErrorLog - critical messages
func (c Config) ErrorLog() log.LogLocation {
	return log.LogLocation(c.Default.Log.Error)
}

// WarningLog - warning messages
func (c Config) WarningLog() log.LogLocation {
	return log.LogLocation(c.Default.Log.Warning)
}

// InfoLog - information messages
func (c Config) InfoLog() log.LogLocation {
	return log.LogLocation(c.Default.Log.Info)
}

// DebugLog - troubleshooting messages
func (c Config) DebugLog() log.LogLocation {
	return log.LogLocation(c.Default.Log.Debug)
}

// EventLog - access.log high level transactions
func (c Config) EventLog() log.LogLocation {
	return log.LogLocation(c.Default.Log.Event)
}
