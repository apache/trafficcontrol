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
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"

	jsoniter "github.com/json-iterator/go"
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
	//CrConfigBackupFile is the default file name to store the last crconfig
	CRConfigBackupFile = "/opt/traffic_monitor/crconfig.backup"
	//TmConfigBackupFile is the default file name to store the last tmconfig
	TMConfigBackupFile = "/opt/traffic_monitor/tmconfig.backup"
	//HTTPPollingFormat is the default accept encoding for stats from caches
	HTTPPollingFormat = "text/json"
)

// PollingProtocol is a string value indicating whether to use IPv4, IPv6, or both.
type PollingProtocol string

const (
	IPv4Only               = PollingProtocol("ipv4only")
	IPv6Only               = PollingProtocol("ipv6only")
	Both                   = PollingProtocol("both")
	InvalidPollingProtocol = PollingProtocol("invalid_polling_protocol")
)

// String returns a string representation of this PollingProtocol.
func (t PollingProtocol) String() string {
	return string(t)
}

// PollingProtocolFromString returns a PollingProtocol based on the string input.
func PollingProtocolFromString(s string) PollingProtocol {
	s = strings.ToLower(s)
	switch s {
	case IPv4Only.String():
		return IPv4Only
	case IPv6Only.String():
		return IPv6Only
	case Both.String():
		return Both
	default:
		return InvalidPollingProtocol
	}
}

// UnmarshalJSON implements the json.Unmarshaller interface
func (t *PollingProtocol) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*t = PollingProtocolFromString(s)
	if *t == InvalidPollingProtocol {
		return errors.New("parsed invalid PollingProtocol: " + s)
	}
	return nil
}

// Config is the configuration for the application. It includes myriad data,
// such as polling intervals and log locations.
type Config struct {
	// Sets the Internet Protocol version used for polling cache servers.
	CachePollingProtocol PollingProtocol `json:"cache_polling_protocol"`
	// A path to a file where CDN Snapshot backups are written.
	CRConfigBackupFile string `json:"crconfig_backup_file"`
	// The number of historical CDN Snapshots to store.
	CRConfigHistoryCount uint64 `json:"crconfig_history_count"`
	// Controls whether Distributed Polling is enabled.
	DistributedPolling bool `json:"distributed_polling"`
	// Defines an interval on which Traffic Monitor will flush its collected
	// health data such that it is made available through the API.
	HealthFlushInterval time.Duration `json:"-"`
	// A MIME-Type that will be sent in the Accept HTTP header in requests to
	// cache servers for health and stats data.
	HTTPPollingFormat string `json:"http_polling_format"`
	// Sets the timeout duration for all HTTP operations - peer-polling and
	// health data polling.
	HTTPTimeout time.Duration `json:"-"`
	// A "lib/go-log"-compliant string defining the logging of Access-level
	// logs.
	LogLocationAccess string `json:"log_location_access"`
	// A "lib/go-log"-compliant string defining the logging of Debug-level logs.
	LogLocationDebug string `json:"log_location_debug"`
	// A "lib/go-log"-compliant string defining the logging of Error-level logs.
	LogLocationError string `json:"log_location_error"`
	// A "lib/go-log"-compliant string defining the logging of Event-level logs.
	LogLocationEvent string `json:"log_location_event"`
	// A "lib/go-log"-compliant string defining the logging of Info-level logs.
	LogLocationInfo string `json:"log_location_info"`
	// A "lib/go-log"-compliant string defining the logging of Warning-level
	// logs.
	LogLocationWarning string `json:"log_location_warning"`
	// The maximum number of events to keep in TM's buffer to be served via its
	// API.
	MaxEvents uint64 `json:"max_events"`
	// The interval on which to poll for this TM's CDN's "monitoring config".
	MonitorConfigPollingInterval time.Duration `json:"-"`
	// Specifies the minimum number of peers that must be available in order to
	// participate in the optimistic health protocol.
	PeerOptimisticQuorumMin int `json:"peer_optimistic_quorum_min"`
	// The timeout for the API server for reading requests.
	ServeReadTimeout time.Duration `json:"-"`
	// The timeout for the API server for writing responses.
	ServeWriteTimeout time.Duration `json:"-"`
	// ShortHostnameOverride is for explicitly setting a hostname rather than using the output of `hostname -s`.
	ShortHostnameOverride string `json:"short_hostname_override"`
	// The interval for which to buffer stats data before processing it.
	StatBufferInterval time.Duration `json:"-"`
	// The interval on which Traffic Monitor will flush its collected stats data
	// such that it is made available through the API.
	StatFlushInterval time.Duration `json:"-"`
	// The directory in which Traffic Monitor will look for the static files for
	// its web UI.
	StaticFileDir string `json:"static_file_dir"`
	// Controls whether stats data is polled.
	StatPolling bool `json:"stat_polling"`
	// A file location to which a backup of the "monitoring configuration"
	// currently in use by Traffic Monitor will be written.
	TMConfigBackupFile string `json:"tmconfig_backup_file"`
	// The number of times Traffic Monitor should attempt to log in to Traffic
	// Ops before using its backup monitoring configuration and CDN Snapshot (if
	// those exist).
	TrafficOpsDiskRetryMax uint64 `json:"traffic_ops_disk_retry_max"`
	// The maximum exponential backoff duration for logging in to Traffic Ops.
	TrafficOpsMaxRetryInterval time.Duration `json:"-"`
	// The minimum exponential backoff duration for logging in to Traffic Ops.
	TrafficOpsMinRetryInterval time.Duration `json:"-"`
}

func (c Config) ErrorLog() log.LogLocation   { return log.LogLocation(c.LogLocationError) }
func (c Config) WarningLog() log.LogLocation { return log.LogLocation(c.LogLocationWarning) }
func (c Config) InfoLog() log.LogLocation    { return log.LogLocation(c.LogLocationInfo) }
func (c Config) DebugLog() log.LogLocation   { return log.LogLocation(c.LogLocationDebug) }
func (c Config) EventLog() log.LogLocation   { return log.LogLocation(c.LogLocationEvent) }
func (c Config) AccessLog() log.LogLocation  { return log.LogLocation(c.LogLocationAccess) }

func GetAccessLogWriter(cfg Config) (io.WriteCloser, error) {
	accessLoc := cfg.AccessLog()

	accessW, err := log.GetLogWriter(accessLoc)
	if err != nil {
		return nil, fmt.Errorf("getting log access writer %v: %v", accessLoc, err)
	}
	return accessW, nil
}

// DefaultConfig is the default configuration for the application, if no configuration file is given, or if a given config setting doesn't exist in the config file.
var DefaultConfig = Config{
	CachePollingProtocol:         Both,
	CRConfigBackupFile:           CRConfigBackupFile,
	CRConfigHistoryCount:         100,
	HealthFlushInterval:          200 * time.Millisecond,
	HTTPPollingFormat:            HTTPPollingFormat,
	HTTPTimeout:                  2 * time.Second,
	LogLocationAccess:            LogLocationNull,
	LogLocationDebug:             LogLocationNull,
	LogLocationError:             LogLocationStderr,
	LogLocationEvent:             LogLocationStdout,
	LogLocationInfo:              LogLocationNull,
	LogLocationWarning:           LogLocationStdout,
	MaxEvents:                    200,
	MonitorConfigPollingInterval: 5 * time.Second,
	PeerOptimisticQuorumMin:      0,
	ServeReadTimeout:             10 * time.Second,
	ServeWriteTimeout:            10 * time.Second,
	ShortHostnameOverride:        "",
	StatBufferInterval:           0,
	StatFlushInterval:            200 * time.Millisecond,
	StaticFileDir:                StaticFileDir,
	StatPolling:                  true,
	TMConfigBackupFile:           TMConfigBackupFile,
	TrafficOpsDiskRetryMax:       2,
	TrafficOpsMaxRetryInterval:   60000 * time.Millisecond,
	TrafficOpsMinRetryInterval:   100 * time.Millisecond,
}

// MarshalJSON marshals custom millisecond durations. Aliasing inspired by http://choly.ca/post/go-json-marshalling/
func (c *Config) MarshalJSON() ([]byte, error) {
	type Alias Config
	json := jsoniter.ConfigFastest // TODO make configurable?
	return json.Marshal(&struct {
		MonitorConfigPollingIntervalMs uint64 `json:"monitor_config_polling_interval_ms"`
		HTTPTimeoutMS                  uint64 `json:"http_timeout_ms"`
		HealthFlushIntervalMs          uint64 `json:"health_flush_interval_ms"`
		StatFlushIntervalMs            uint64 `json:"stat_flush_interval_ms"`
		StatBufferIntervalMs           uint64 `json:"stat_buffer_interval_ms"`
		ServeReadTimeoutMs             uint64 `json:"serve_read_timeout_ms"`
		ServeWriteTimeoutMs            uint64 `json:"serve_write_timeout_ms"`
		*Alias
	}{
		MonitorConfigPollingIntervalMs: uint64(c.MonitorConfigPollingInterval / time.Millisecond),
		HTTPTimeoutMS:                  uint64(c.HTTPTimeout / time.Millisecond),
		HealthFlushIntervalMs:          uint64(c.HealthFlushInterval / time.Millisecond),
		StatFlushIntervalMs:            uint64(c.StatFlushInterval / time.Millisecond),
		StatBufferIntervalMs:           uint64(c.StatBufferInterval / time.Millisecond),
		Alias:                          (*Alias)(c),
	})
}

// UnmarshalJSON populates this config object from given JSON bytes.
func (c *Config) UnmarshalJSON(data []byte) error {
	type Alias Config
	aux := &struct {
		MonitorConfigPollingIntervalMs *uint64 `json:"monitor_config_polling_interval_ms"`
		HTTPTimeoutMS                  *uint64 `json:"http_timeout_ms"`
		HealthFlushIntervalMs          *uint64 `json:"health_flush_interval_ms"`
		StatFlushIntervalMs            *uint64 `json:"stat_flush_interval_ms"`
		StatBufferIntervalMs           *uint64 `json:"stat_buffer_interval_ms"`
		ServeReadTimeoutMs             *uint64 `json:"serve_read_timeout_ms"`
		ServeWriteTimeoutMs            *uint64 `json:"serve_write_timeout_ms"`
		TrafficOpsMinRetryIntervalMs   *uint64 `json:"traffic_ops_min_retry_interval_ms"`
		TrafficOpsMaxRetryIntervalMs   *uint64 `json:"traffic_ops_max_retry_interval_ms"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	json := jsoniter.ConfigFastest // TODO make configurable?
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.MonitorConfigPollingIntervalMs != nil {
		c.MonitorConfigPollingInterval = time.Duration(*aux.MonitorConfigPollingIntervalMs) * time.Millisecond
	}
	if aux.HTTPTimeoutMS != nil {
		c.HTTPTimeout = time.Duration(*aux.HTTPTimeoutMS) * time.Millisecond
	}
	if aux.HealthFlushIntervalMs != nil {
		c.HealthFlushInterval = time.Duration(*aux.HealthFlushIntervalMs) * time.Millisecond
	}
	if aux.StatFlushIntervalMs != nil {
		c.StatFlushInterval = time.Duration(*aux.StatFlushIntervalMs) * time.Millisecond
	}
	if aux.StatBufferIntervalMs != nil {
		c.StatBufferInterval = time.Duration(*aux.StatBufferIntervalMs) * time.Millisecond
	}
	if aux.ServeReadTimeoutMs != nil {
		c.ServeReadTimeout = time.Duration(*aux.ServeReadTimeoutMs) * time.Millisecond
	}
	if aux.ServeWriteTimeoutMs != nil {
		c.ServeWriteTimeout = time.Duration(*aux.ServeWriteTimeoutMs) * time.Millisecond
	}
	if aux.TrafficOpsMinRetryIntervalMs != nil {
		c.TrafficOpsMinRetryInterval = time.Duration(*aux.TrafficOpsMinRetryIntervalMs) * time.Millisecond
	}
	if aux.TrafficOpsMaxRetryIntervalMs != nil {
		c.TrafficOpsMaxRetryInterval = time.Duration(*aux.TrafficOpsMaxRetryIntervalMs) * time.Millisecond
	}
	if c.StatPolling && c.DistributedPolling {
		return errors.New("invalid configuration: stat_polling cannot be enabled if distributed_polling is also enabled")
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
	json := jsoniter.ConfigFastest // TODO make configurable?
	err := json.Unmarshal(bytes, &cfg)
	return cfg, err
}
