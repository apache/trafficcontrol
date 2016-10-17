package config

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type LogLocation string

const LogLocationStdout = "stdout"
const LogLocationStderr = "stderr"
const LogLocationNull = "null"

type Config struct {
	CacheHealthPollingInterval   time.Duration `json:"-"`
	CacheStatPollingInterval     time.Duration `json:"-"`
	MonitorConfigPollingInterval time.Duration `json:"-"`
	HttpTimeout                  time.Duration `json:"-"`
	PeerPollingInterval          time.Duration `json:"-"`
	MaxEvents                    uint64        `json:"max_events"`
	MaxStatHistory               uint64        `json:"max_stat_history"`
	MaxHealthHistory             uint64        `json:"max_health_history"`
	HealthFlushInterval          time.Duration `json:"-"`
	StatFlushInterval            time.Duration `json:"-"`
	LogLocationError             string        `json:"log_location_error"`
	LogLocationWarning           string        `json:"log_location_warning"`
	LogLocationInfo              string        `json:"log_location_info"`
	LogLocationDebug             string        `json:"log_location_debug"`
	ServeReadTimeout             time.Duration `json:"-"`
	ServeWriteTimeout            time.Duration `json:"-"`
}

var DefaultConfig = Config{
	CacheHealthPollingInterval:   6 * time.Second,
	CacheStatPollingInterval:     6 * time.Second,
	MonitorConfigPollingInterval: 5 * time.Second,
	HttpTimeout:                  2 * time.Second,
	PeerPollingInterval:          5 * time.Second,
	MaxEvents:                    200,
	MaxStatHistory:               5,
	MaxHealthHistory:             5,
	HealthFlushInterval:          200 * time.Millisecond,
	StatFlushInterval:            200 * time.Millisecond,
	LogLocationError:             LogLocationStderr,
	LogLocationWarning:           LogLocationStdout,
	LogLocationInfo:              LogLocationNull,
	LogLocationDebug:             LogLocationNull,
	ServeReadTimeout:             10 * time.Second,
	ServeWriteTimeout:            10 * time.Second,
}

// MarshalJSON marshals custom millisecond durations. Aliasing inspired by http://choly.ca/post/go-json-marshalling/
func (c *Config) MarshalJSON() ([]byte, error) {
	type Alias Config
	return json.Marshal(&struct {
		CacheHealthPollingIntervalMs   uint64 `json:"cache_health_polling_interval_ms"`
		CacheStatPollingIntervalMs     uint64 `json:"cache_stat_polling_interval_ms"`
		MonitorConfigPollingIntervalMs uint64 `json:"monitor_config_polling_interval_ms"`
		HttpTimeoutMs                  uint64 `json:"http_timeout_ms"`
		PeerPollingIntervalMs          uint64 `json:"peer_polling_interval_ms"`
		HealthFlushIntervalMs          uint64 `json:"health_flush_interval_ms"`
		StatFlushIntervalMs            uint64 `json:"stat_flush_interval_ms"`
		ServeReadTimeoutMs             uint64 `json:"serve_read_timeout_ms"`
		ServeWriteTimeoutMs            uint64 `json:"serve_write_timeout_ms"`
		*Alias
	}{
		CacheHealthPollingIntervalMs:   uint64(c.CacheHealthPollingInterval / time.Millisecond),
		CacheStatPollingIntervalMs:     uint64(c.CacheStatPollingInterval / time.Millisecond),
		MonitorConfigPollingIntervalMs: uint64(c.MonitorConfigPollingInterval / time.Millisecond),
		HttpTimeoutMs:                  uint64(c.HttpTimeout / time.Millisecond),
		PeerPollingIntervalMs:          uint64(c.PeerPollingInterval / time.Millisecond),
		HealthFlushIntervalMs:          uint64(c.HealthFlushInterval / time.Millisecond),
		StatFlushIntervalMs:            uint64(c.StatFlushInterval / time.Millisecond),
		Alias:                          (*Alias)(c),
	})
}

func (c *Config) UnmarshalJSON(data []byte) error {
	type Alias Config
	aux := &struct {
		CacheHealthPollingIntervalMs   *uint64 `json:"cache_health_polling_interval_ms"`
		CacheStatPollingIntervalMs     *uint64 `json:"cache_stat_polling_interval_ms"`
		MonitorConfigPollingIntervalMs *uint64 `json:"monitor_config_polling_interval_ms"`
		HttpTimeoutMs                  *uint64 `json:"http_timeout_ms"`
		PeerPollingIntervalMs          *uint64 `json:"peer_polling_interval_ms"`
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
	if aux.HttpTimeoutMs != nil {
		c.HttpTimeout = time.Duration(*aux.HttpTimeoutMs) * time.Millisecond
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
	return nil
}

// Load loads the given config file. If an empty string is passed, the default config is returned.
func Load(fileName string) (Config, error) {
	cfg := DefaultConfig
	if fileName == "" {
		return cfg, nil
	}
	configBytes, err := ioutil.ReadFile(fileName)
	if err == nil {
		err = json.Unmarshal(configBytes, &cfg)
	}
	return cfg, err
}
