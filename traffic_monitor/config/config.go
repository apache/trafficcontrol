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
	"io/ioutil"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"

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

// Config is the configuration for the application. It includes myriad data, such as polling intervals and log locations.
type Config struct {
	CacheHealthPollingInterval   time.Duration   `json:"-"`
	CacheStatPollingInterval     time.Duration   `json:"-"`
	MonitorConfigPollingInterval time.Duration   `json:"-"`
	HTTPTimeout                  time.Duration   `json:"-"`
	PeerPollingInterval          time.Duration   `json:"-"`
	PeerOptimistic               bool            `json:"peer_optimistic"`
	PeerOptimisticQuorumMin      int             `json:"peer_optimistic_quorum_min"`
	MaxEvents                    uint64          `json:"max_events"`
	MaxStatHistory               uint64          `json:"max_stat_history"`
	MaxHealthHistory             uint64          `json:"max_health_history"`
	HealthFlushInterval          time.Duration   `json:"-"`
	StatFlushInterval            time.Duration   `json:"-"`
	StatBufferInterval           time.Duration   `json:"-"`
	LogLocationError             string          `json:"log_location_error"`
	LogLocationWarning           string          `json:"log_location_warning"`
	LogLocationInfo              string          `json:"log_location_info"`
	LogLocationDebug             string          `json:"log_location_debug"`
	LogLocationEvent             string          `json:"log_location_event"`
	ServeReadTimeout             time.Duration   `json:"-"`
	ServeWriteTimeout            time.Duration   `json:"-"`
	HealthToStatRatio            uint64          `json:"health_to_stat_ratio"`
	StaticFileDir                string          `json:"static_file_dir"`
	CRConfigHistoryCount         uint64          `json:"crconfig_history_count"`
	TrafficOpsMinRetryInterval   time.Duration   `json:"-"`
	TrafficOpsMaxRetryInterval   time.Duration   `json:"-"`
	CRConfigBackupFile           string          `json:"crconfig_backup_file"`
	TMConfigBackupFile           string          `json:"tmconfig_backup_file"`
	TrafficOpsDiskRetryMax       uint64          `json:"-"`
	CachePollingProtocol         PollingProtocol `json:"cache_polling_protocol"`
	PeerPollingProtocol          PollingProtocol `json:"peer_polling_protocol"`
	HTTPPollingFormat            string          `json:"http_polling_format"`
}

func (c Config) ErrorLog() log.LogLocation   { return log.LogLocation(c.LogLocationError) }
func (c Config) WarningLog() log.LogLocation { return log.LogLocation(c.LogLocationWarning) }
func (c Config) InfoLog() log.LogLocation    { return log.LogLocation(c.LogLocationInfo) }
func (c Config) DebugLog() log.LogLocation   { return log.LogLocation(c.LogLocationDebug) }
func (c Config) EventLog() log.LogLocation   { return log.LogLocation(c.LogLocationEvent) }

// DefaultConfig is the default configuration for the application, if no configuration file is given, or if a given config setting doesn't exist in the config file.
var DefaultConfig = Config{
	CacheHealthPollingInterval:   6 * time.Second,
	CacheStatPollingInterval:     6 * time.Second,
	MonitorConfigPollingInterval: 5 * time.Second,
	HTTPTimeout:                  2 * time.Second,
	PeerPollingInterval:          5 * time.Second,
	PeerOptimistic:               true,
	PeerOptimisticQuorumMin:      0,
	MaxEvents:                    200,
	MaxStatHistory:               5,
	MaxHealthHistory:             5,
	HealthFlushInterval:          200 * time.Millisecond,
	StatFlushInterval:            200 * time.Millisecond,
	StatBufferInterval:           0,
	LogLocationError:             LogLocationStderr,
	LogLocationWarning:           LogLocationStdout,
	LogLocationInfo:              LogLocationNull,
	LogLocationDebug:             LogLocationNull,
	LogLocationEvent:             LogLocationStdout,
	ServeReadTimeout:             10 * time.Second,
	ServeWriteTimeout:            10 * time.Second,
	HealthToStatRatio:            4,
	StaticFileDir:                StaticFileDir,
	CRConfigHistoryCount:         20000,
	TrafficOpsMinRetryInterval:   100 * time.Millisecond,
	TrafficOpsMaxRetryInterval:   60000 * time.Millisecond,
	CRConfigBackupFile:           CRConfigBackupFile,
	TMConfigBackupFile:           TMConfigBackupFile,
	TrafficOpsDiskRetryMax:       2,
	CachePollingProtocol:         Both,
	PeerPollingProtocol:          Both,
	HTTPPollingFormat:            HTTPPollingFormat,
}

// MarshalJSON marshals custom millisecond durations. Aliasing inspired by http://choly.ca/post/go-json-marshalling/
func (c *Config) MarshalJSON() ([]byte, error) {
	type Alias Config
	json := jsoniter.ConfigFastest // TODO make configurable?
	return json.Marshal(&struct {
		CacheHealthPollingIntervalMs   uint64 `json:"cache_health_polling_interval_ms"`
		CacheStatPollingIntervalMs     uint64 `json:"cache_stat_polling_interval_ms"`
		MonitorConfigPollingIntervalMs uint64 `json:"monitor_config_polling_interval_ms"`
		HTTPTimeoutMS                  uint64 `json:"http_timeout_ms"`
		PeerPollingIntervalMs          uint64 `json:"peer_polling_interval_ms"`
		PeerOptimistic                 bool   `json:"peer_optimistic"`
		PeerOptimisticQuorumMin        int    `json:"peer_optimistic_quorum_min"`
		HealthFlushIntervalMs          uint64 `json:"health_flush_interval_ms"`
		StatFlushIntervalMs            uint64 `json:"stat_flush_interval_ms"`
		StatBufferIntervalMs           uint64 `json:"stat_buffer_interval_ms"`
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
		PeerOptimisticQuorumMin:        int(c.PeerOptimisticQuorumMin),
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
		CacheHealthPollingIntervalMs   *uint64 `json:"cache_health_polling_interval_ms"`
		CacheStatPollingIntervalMs     *uint64 `json:"cache_stat_polling_interval_ms"`
		MonitorConfigPollingIntervalMs *uint64 `json:"monitor_config_polling_interval_ms"`
		HTTPTimeoutMS                  *uint64 `json:"http_timeout_ms"`
		PeerPollingIntervalMs          *uint64 `json:"peer_polling_interval_ms"`
		PeerOptimistic                 *bool   `json:"peer_optimistic"`
		PeerOptimisticQuorumMin        *int    `json:"peer_optimistic_quorum_min"`
		HealthFlushIntervalMs          *uint64 `json:"health_flush_interval_ms"`
		StatFlushIntervalMs            *uint64 `json:"stat_flush_interval_ms"`
		StatBufferIntervalMs           *uint64 `json:"stat_buffer_interval_ms"`
		ServeReadTimeoutMs             *uint64 `json:"serve_read_timeout_ms"`
		ServeWriteTimeoutMs            *uint64 `json:"serve_write_timeout_ms"`
		TrafficOpsMinRetryIntervalMs   *uint64 `json:"traffic_ops_min_retry_interval_ms"`
		TrafficOpsMaxRetryIntervalMs   *uint64 `json:"traffic_ops_max_retry_interval_ms"`
		TrafficOpsDiskRetryMax         *uint64 `json:"traffic_ops_disk_retry_max"`
		CRConfigBackupFile             *string `json:"crconfig_backup_file"`
		TMConfigBackupFile             *string `json:"tmconfig_backup_file"`
		HTTPPollingFormat              *string `json:"http_polling_format"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	json := jsoniter.ConfigFastest // TODO make configurable?
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
	if aux.StatBufferIntervalMs != nil {
		c.StatBufferInterval = time.Duration(*aux.StatBufferIntervalMs) * time.Millisecond
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
	if aux.PeerOptimisticQuorumMin != nil {
		c.PeerOptimisticQuorumMin = *aux.PeerOptimisticQuorumMin
	}
	if aux.TrafficOpsMinRetryIntervalMs != nil {
		c.TrafficOpsMinRetryInterval = time.Duration(*aux.TrafficOpsMinRetryIntervalMs) * time.Millisecond
	}
	if aux.TrafficOpsMaxRetryIntervalMs != nil {
		c.TrafficOpsMaxRetryInterval = time.Duration(*aux.TrafficOpsMaxRetryIntervalMs) * time.Millisecond
	}
	if aux.TrafficOpsDiskRetryMax != nil {
		c.TrafficOpsDiskRetryMax = *aux.TrafficOpsDiskRetryMax
	}
	if aux.CRConfigBackupFile != nil {
		c.CRConfigBackupFile = *aux.CRConfigBackupFile
	}
	if aux.TMConfigBackupFile != nil {
		c.TMConfigBackupFile = *aux.TMConfigBackupFile
	}
	if aux.HTTPPollingFormat != nil {
		c.HTTPPollingFormat = *aux.HTTPPollingFormat
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
