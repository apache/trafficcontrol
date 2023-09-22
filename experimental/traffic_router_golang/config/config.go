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
 *
 */

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
)

type Cfg struct {
	Port                  uint     `json:"port"`
	Monitors              []*URL   `json:"monitors"`
	ReqTimeout            Duration `json:"request_timeout_ms"`
	CRConfigInterval      Duration `json:"crconfig_poll_interval_ms"`
	CRStatesInterval      Duration `json:"crstates_poll_interval_ms"`
	CDN                   string   `json:"cdn"`
	TrafficOpsURI         *URL     `json:"traffic_ops_uri"`
	TrafficOpsUser        string   `json:"traffic_ops_user"`
	TrafficOpsPass        string   `json:"traffic_ops_pass"`
	TrafficOpsInsecure    bool     `json:"traffic_ops_insecure"`
	TrafficOpsClientCache bool     `json:"traffic_ops_client_cache"`
	TrafficOpsTimeout     Duration `json:"traffic_ops_timeout_ms"`
	CoverageZoneFile      string   `json:"coverage_zone_file"`
	LogLocations
}

func Load(filename string) (Cfg, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return Cfg{}, errors.New("loading " + filename + ": " + err.Error())
	}

	cfg := Cfg{}
	if err := json.Unmarshal(b, &cfg); err != nil {
		return Cfg{}, errors.New("parsing " + filename + ": " + err.Error())
	}
	return cfg, nil
}

type LogLocations struct {
	LogLocationErr   string `json:"log_location_error"`
	LogLocationWarn  string `json:"log_location_warning"`
	LogLocationInfo  string `json:"log_location_info"`
	LogLocationDebug string `json:"log_location_debug"`
	LogLocationEvent string `json:"log_location_event"`
}

// Fulfill the log.Logger interface

func (c Cfg) ErrorLog() log.LogLocation   { return log.LogLocation(c.LogLocationErr) }
func (c Cfg) WarningLog() log.LogLocation { return log.LogLocation(c.LogLocationWarn) }
func (c Cfg) InfoLog() log.LogLocation    { return log.LogLocation(c.LogLocationInfo) }
func (c Cfg) DebugLog() log.LogLocation   { return log.LogLocation(c.LogLocationDebug) }
func (c Cfg) EventLog() log.LogLocation   { return log.LogLocation(c.LogLocationEvent) }

type URL url.URL

func (u URL) MarshalJSON() ([]byte, error) {
	return []byte((*url.URL)(&u).String()), nil
}

func (u *URL) UnmarshalJSON(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return errors.New("unquoting string: " + err.Error())
	}
	newU, err := url.ParseRequestURI(s) // ParseRequestURI is absolute, Parse may be relative; we want absolute.
	if err != nil {
		return errors.New("parsing URL: " + err.Error())
	}
	*u = *(*URL)(newU)
	return nil
}

// Duration is a JSON config time.Duration, which is marshalled and unmarshalled as milliseconds. Therefore, JSON keys of this type should be suffixed with '_ms'.
type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(time.Duration(d)/time.Millisecond), 10)), nil
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	i, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return errors.New("converting string to number: " + err.Error())
	}
	*d = Duration(time.Duration(i) * time.Millisecond)
	return nil
}
