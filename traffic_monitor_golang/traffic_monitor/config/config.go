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
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
)

// LogLocation is a location to log to. This may be stdout, stderr, null (/dev/null), or a valid file path.
type LogLocation string

const (
	// LogLocationStdout indicates the stdout IO stream
	LogLocationStdout = "stdout"
	// LogLocationStderr indicates the stderr IO stream
	LogLocationStderr = "stderr"
	// LogLocationNull indicates the null IO stream (/dev/null)
	LogLocationNull = "null"
	//StaticFileDir is the directory that contains static html and js files.
	StaticFileDir = "/opt/traffic_monitor/static/"
)

// Config is the configuration for the application. It includes myriad data, such as polling intervals and log locations.
type Config struct {
	CacheHealthPollingInterval   time.Duration `json:"-"`
	CacheStatPollingInterval     time.Duration `json:"-"`
	MonitorConfigPollingInterval time.Duration `json:"-"`
	HTTPTimeout                  time.Duration `json:"-"`
	PeerPollingInterval          time.Duration `json:"-"`
	PeerOptimistic               bool          `json:"peer_optimistic"`
	MaxEvents                    uint64        `json:"max_events"`
	MaxStatHistory               uint64        `json:"max_stat_history"`
	MaxHealthHistory             uint64        `json:"max_health_history"`
	HealthFlushInterval          time.Duration `json:"-"`
	StatFlushInterval            time.Duration `json:"-"`
	LogLocationError             string        `json:"log_location_error"`
	LogLocationWarning           string        `json:"log_location_warning"`
	LogLocationInfo              string        `json:"log_location_info"`
	LogLocationDebug             string        `json:"log_location_debug"`
	LogLocationEvent             string        `json:"log_location_event"`
	ServeReadTimeout             time.Duration `json:"-"`
	ServeWriteTimeout            time.Duration `json:"-"`
	HealthToStatRatio            uint64        `json:"health_to_stat_ratio"`
	HTTPPollNoSleep              bool          `json:"http_poll_no_sleep"`
	StaticFileDir                string        `json:"static_file_dir"`
}

// DefaultConfig is the default configuration for the application, if no configuration file is given, or if a given config setting doesn't exist in the config file.
var DefaultConfig = Config{
	CacheHealthPollingInterval:   6 * time.Second,
	CacheStatPollingInterval:     6 * time.Second,
	MonitorConfigPollingInterval: 5 * time.Second,
	HTTPTimeout:                  2 * time.Second,
	PeerPollingInterval:          5 * time.Second,
	PeerOptimistic:               true,
	MaxEvents:                    200,
	MaxStatHistory:               5,
	MaxHealthHistory:             5,
	HealthFlushInterval:          200 * time.Millisecond,
	StatFlushInterval:            200 * time.Millisecond,
	LogLocationError:             LogLocationStderr,
	LogLocationWarning:           LogLocationStdout,
	LogLocationInfo:              LogLocationNull,
	LogLocationDebug:             LogLocationNull,
	LogLocationEvent:             LogLocationStdout,
	ServeReadTimeout:             10 * time.Second,
	ServeWriteTimeout:            10 * time.Second,
	HealthToStatRatio:            4,
	HTTPPollNoSleep:              false,
	StaticFileDir:                StaticFileDir,
}

// MarshalJSON marshals custom millisecond durations. Aliasing inspired by http://choly.ca/post/go-json-marshalling/
func (c *Config) MarshalJSON() ([]byte, error) {
	type Alias Config
	return json.Marshal(&struct {
		CacheHealthPollingIntervalMs   uint64 `json:"cache_health_polling_interval_ms"`
		CacheStatPollingIntervalMs     uint64 `json:"cache_stat_polling_interval_ms"`
		MonitorConfigPollingIntervalMs uint64 `json:"monitor_config_polling_interval_ms"`
		HTTPTimeoutMS                  uint64 `json:"http_timeout_ms"`
		PeerPollingIntervalMs          uint64 `json:"peer_polling_interval_ms"`
		PeerOptimistic                 bool   `json:"peer_optimistic"`
		HealthFlushIntervalMs          uint64 `json:"health_flush_interval_ms"`
		StatFlushIntervalMs            uint64 `json:"stat_flush_interval_ms"`
		ServeReadTimeoutMs             uint64 `json:"serve_read_timeout_ms"`
		ServeWriteTimeoutMs            uint64 `json:"serve_write_timeout_ms"`
		*Alias
	}{
		CacheHealthPollingIntervalMs:   uint64(c.CacheHealthPollingInterval / time.Millisecond),
		CacheStatPollingIntervalMs:     uint64(c.CacheStatPollingInterval / time.Millisecond),
		MonitorConfigPollingIntervalMs: uint64(c.MonitorConfigPollingInterval / time.Millisecond),
		HTTPTimeoutMS:                  uint64(c.HTTPTimeout / time.Millisecond),
		PeerPollingIntervalMs:          uint64(c.PeerPollingInterval / time.Millisecond),
		PeerOptimistic:                 bool(true),
		HealthFlushIntervalMs:          uint64(c.HealthFlushInterval / time.Millisecond),
		StatFlushIntervalMs:            uint64(c.StatFlushInterval / time.Millisecond),
		Alias:                          (*Alias)(c),
	})
}

// UnmarshalJSON populates this config object from given JSON bytes.
func (c *Config) UnmarshalJSON(data []byte) error {
	type Alias Config
	aux := &struct {
		CacheHealthPollingIntervalMs   *uint64 `json:"cache_health_polling_interval_ms"`
		CacheStatPollingIntervalMs     *uint64 `json:"cache_stat_polling_interval_ms"`
		MonitorConfigPollingIntervalMs *uint64 `json:"monitor_config_polling_interval_ms"`
		HTTPTimeoutMS                  *uint64 `json:"http_timeout_ms"`
		PeerPollingIntervalMs          *uint64 `json:"peer_polling_interval_ms"`
		PeerOptimistic                 *bool   `json:"peer_optimistic"`
		HealthFlushIntervalMs          *uint64 `json:"health_flush_interval_ms"`
		StatFlushIntervalMs            *uint64 `json:"stat_flush_interval_ms"`
		ServeReadTimeoutMs             *uint64 `json:"serve_read_timeout_ms"`
		ServeWriteTimeoutMs            *uint64 `json:"serve_write_timeout_ms"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.CacheHealthPollingIntervalMs != nil {
		c.CacheHealthPollingInterval = time.Duration(*aux.CacheHealthPollingIntervalMs) * time.Millisecond
	}
	if aux.CacheStatPollingIntervalMs != nil {
		c.CacheStatPollingInterval = time.Duration(*aux.CacheStatPollingIntervalMs) * time.Millisecond
	}
	if aux.MonitorConfigPollingIntervalMs != nil {
		c.MonitorConfigPollingInterval = time.Duration(*aux.MonitorConfigPollingIntervalMs) * time.Millisecond
	}
	if aux.HTTPTimeoutMS != nil {
		c.HTTPTimeout = time.Duration(*aux.HTTPTimeoutMS) * time.Millisecond
	}
	if aux.PeerPollingIntervalMs != nil {
		c.PeerPollingInterval = time.Duration(*aux.PeerPollingIntervalMs) * time.Millisecond
	}
	if aux.HealthFlushIntervalMs != nil {
		c.HealthFlushInterval = time.Duration(*aux.HealthFlushIntervalMs) * time.Millisecond
	}
	if aux.StatFlushIntervalMs != nil {
		c.StatFlushInterval = time.Duration(*aux.StatFlushIntervalMs) * time.Millisecond
	}
	if aux.ServeReadTimeoutMs != nil {
		c.ServeReadTimeout = time.Duration(*aux.ServeReadTimeoutMs) * time.Millisecond
	}
	if aux.ServeWriteTimeoutMs != nil {
		c.ServeWriteTimeout = time.Duration(*aux.ServeWriteTimeoutMs) * time.Millisecond
	}
	if aux.PeerOptimistic != nil {
		c.PeerOptimistic = *aux.PeerOptimistic
	}
	return nil
}

// Load loads the given config file. If an empty string is passed, the default config is returned.
func Load(fileName string) (Config, error) {
	cfg := DefaultConfig
	if fileName == "" {
		return cfg, nil
	}
	configBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return DefaultConfig, err
	}
	return LoadBytes(configBytes)
}

// LoadBytes loads the given file bytes.
func LoadBytes(bytes []byte) (Config, error) {
	cfg := DefaultConfig
	err := json.Unmarshal(bytes, &cfg)
	return cfg, err
}

func getLogWriter(location string) (io.WriteCloser, error) {
	switch location {
	case LogLocationStdout:
		return log.NopCloser(os.Stdout), nil
	case LogLocationStderr:
		return log.NopCloser(os.Stderr), nil
	case LogLocationNull:
		return log.NopCloser(ioutil.Discard), nil
	default:
		return os.OpenFile(location, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	}
}

func GetLogWriters(cfg Config) (io.WriteCloser, io.WriteCloser, io.WriteCloser, io.WriteCloser, io.WriteCloser, error) {
	eventLoc := cfg.LogLocationEvent
	errLoc := cfg.LogLocationError
	warnLoc := cfg.LogLocationWarning
	infoLoc := cfg.LogLocationInfo
	debugLoc := cfg.LogLocationDebug

	eventW, err := getLogWriter(eventLoc)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("getting log event writer %v: %v", eventLoc, err)
	}
	errW, err := getLogWriter(errLoc)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("getting log error writer %v: %v", errLoc, err)
	}
	warnW, err := getLogWriter(warnLoc)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("getting log warning writer %v: %v", warnLoc, err)
	}
	infoW, err := getLogWriter(infoLoc)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("getting log info writer %v: %v", infoLoc, err)
	}
	debugW, err := getLogWriter(debugLoc)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("getting log debug writer %v: %v", debugLoc, err)
	}
	return eventW, errW, warnW, infoW, debugW, nil
}
