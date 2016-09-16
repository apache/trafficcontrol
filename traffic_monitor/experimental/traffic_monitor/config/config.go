package config

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type Config struct {
	CacheHealthPollingInterval   time.Duration `json:"-"`
	CacheStatPollingInterval     time.Duration `json:"-"`
	MonitorConfigPollingInterval time.Duration `json:"-"`
	HttpTimeout                  time.Duration `json:"-"`
	PeerPollingInterval          time.Duration `json:"-"`
	MaxEvents                    uint64        `json:"max_events"`
	MaxStatHistory               uint64        `json:"max_stat_history"`
	MaxHealthHistory             uint64        `json:"max_health_history"`
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
		*Alias
	}{
		CacheHealthPollingIntervalMs:   uint64(c.CacheHealthPollingInterval / time.Millisecond),
		CacheStatPollingIntervalMs:     uint64(c.CacheStatPollingInterval / time.Millisecond),
		MonitorConfigPollingIntervalMs: uint64(c.MonitorConfigPollingInterval / time.Millisecond),
		HttpTimeoutMs:                  uint64(c.HttpTimeout / time.Millisecond),
		PeerPollingIntervalMs:          uint64(c.PeerPollingInterval / time.Millisecond),
		Alias: (*Alias)(c),
	})
}

func (c *Config) UnmarshalJSON(data []byte) error {
	type Alias Config
	aux := &struct {
		CacheHealthPollingIntervalMs   uint64 `json:"cache_health_polling_interval_ms"`
		CacheStatPollingIntervalMs     uint64 `json:"cache_stat_polling_interval_ms"`
		MonitorConfigPollingIntervalMs uint64 `json:"monitor_config_polling_interval_ms"`
		HttpTimeoutMs                  uint64 `json:"http_timeout_ms"`
		PeerPollingIntervalMs          uint64 `json:"peer_polling_interval_ms"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	c.CacheHealthPollingInterval = time.Duration(aux.CacheHealthPollingIntervalMs) * time.Millisecond
	c.CacheStatPollingInterval = time.Duration(aux.CacheStatPollingIntervalMs) * time.Millisecond
	c.MonitorConfigPollingInterval = time.Duration(aux.MonitorConfigPollingIntervalMs) * time.Millisecond
	c.HttpTimeout = time.Duration(aux.HttpTimeoutMs) * time.Millisecond
	c.PeerPollingInterval = time.Duration(aux.PeerPollingIntervalMs) * time.Millisecond
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
