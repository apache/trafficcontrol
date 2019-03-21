package config

/*
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

import (
	"encoding/json"
	"io/ioutil"

	"github.com/apache/trafficcontrol/lib/go-log"
)

const bytesPerGibibyte = 1024 * 1024 * 1024
const bytesPerMebibyte = 1024 * 1024

type Config struct {
	// CertFile is the path to the certificate to use for HTTPS requests which don't match any remap rule.
	CertFile string `json:"cert_file"`

	// ConcurrentRuleRequests is the number of concurrent requests permitted to a remap rule, that is, to an origin. Note this is overridden by any per-rule settings in the remap rules.
	ConcurrentRuleRequests int `json:"concurrent_rule_requests"`

	// ConnectionClose determines whether to send a `Connection: close` header. This is primarily designed for maintenance, to drain the cache of incoming requestors. This overrides rule-specific `connection-close: false` configuration, under the assumption that draining a cache is a temporary maintenance operation, and if connectionClose is true on the service and false on some rules, those rules' configuration is probably a permament setting whereas the operator probably wants to drain all connections if the global setting is true. If it's necessary to leave connection close false on some rules, set all other rules' connectionClose to true and leave the global connectionClose unset.
	ConnectionClose bool `json:"connection_close"`

	// DefaultCache is the default cache to use for remap rules which don't specify.
	DefaultCache string `json:"default_cache"`

	DisableHTTP2 bool `json:"disable_http2"`

	HTTPSPort int `json:"https_port"`

	InterfaceName string `json:"interface_name"`

	// KeyFile is the path to the key for the CertFile certificate.
	KeyFile string `json:"key_file"`

	LogLocationDebug   string `json:"log_location_debug"`
	LogLocationError   string `json:"log_location_error"`
	LogLocationEvent   string `json:"log_location_event"`
	LogLocationInfo    string `json:"log_location_info"`
	LogLocationWarning string `json:"log_location_warning"`

	Plugins []string `json:"plugins"`

	// Port is the HTTP port to serve on.
	Port int `json:"port"`

	RemapRulesFile string `json:"remap_rules_file"`

	// RFCCompliant determines whether `Cache-Control: no-cache` requests are honored. The ability to ignore `no-cache` is necessary to protect origin servers from DDOS attacks. In general, CDNs and caching proxies with the goal of origin protection should set RFCComplaint false. Cache with other goals (performance, load balancing, etc) should set RFCCompliant true.
	RFCCompliant bool `json:"rfc_compliant"`

	ReqIdleConnTimeoutMS int `json:"parent_request_idle_connection_timeout_ms"`
	ReqKeepAliveMS       int `json:"parent_request_keep_alive_ms"`
	ReqMaxIdleConns      int `json:"parent_request_max_idle_connections"`
	ReqTimeoutMS         int `json:"parent_request_timeout_ms"` // TODO rename "parent_request" to distinguish from client requests

	ServerIdleTimeoutMS  int `json:"server_idle_timeout_ms"`
	ServerReadTimeoutMS  int `json:"server_read_timeout_ms"`
	ServerWriteTimeoutMS int `json:"server_write_timeout_ms"`

	// CacheSizeBytes is the size of the memory cache, in bytes.
	// To disable the default lru mem cache, and only use other cache types, set this explicitly to 0.
	CacheSizeBytes int `json:"cache_size_bytes"`

	// CacheFiles is the name of each disk cache, and the files to use for it.
	CacheFiles map[string][]CacheFile `json:"cache_files"`

	// FileMemBytes is the amount of memory to use as an LRU in front of each name in CacheFiles, that is, each named group of files. E.g. if there are 10 files, the amount of memory used will be 10*FileMemBytes+CacheSizeBytes.
	FileMemBytes int `json:"file_mem_bytes"`
}

type CacheFile struct {
	Path  string `json:"path"`
	Bytes uint64 `json:"size_bytes"`
}

func (c Config) ErrorLog() log.LogLocation {
	return log.LogLocation(c.LogLocationError)
}
func (c Config) WarningLog() log.LogLocation {
	return log.LogLocation(c.LogLocationWarning)
}
func (c Config) InfoLog() log.LogLocation {
	return log.LogLocation(c.LogLocationInfo)
}
func (c Config) DebugLog() log.LogLocation {
	return log.LogLocation(c.LogLocationDebug)
}
func (c Config) EventLog() log.LogLocation {
	return log.LogLocation(c.LogLocationEvent)
}

const MSPerSec = 1000

// DefaultConfig is the default configuration for the application, if no configuration file is given, or if a given config setting doesn't exist in the config file.
var DefaultConfig = Config{
	RFCCompliant:           true,
	Port:                   80,
	DisableHTTP2:           false,
	HTTPSPort:              443,
	CacheSizeBytes:         bytesPerGibibyte,
	RemapRulesFile:         "remap.config",
	ConcurrentRuleRequests: 100000,
	ConnectionClose:        false,
	LogLocationError:       log.LogLocationStderr,
	LogLocationWarning:     log.LogLocationStdout,
	LogLocationInfo:        log.LogLocationNull,
	LogLocationDebug:       log.LogLocationNull,
	LogLocationEvent:       log.LogLocationStdout,
	ReqTimeoutMS:           30 * MSPerSec,
	ReqKeepAliveMS:         30 * MSPerSec,
	ReqMaxIdleConns:        100,
	ReqIdleConnTimeoutMS:   90 * MSPerSec,
	ServerIdleTimeoutMS:    10 * MSPerSec,
	ServerWriteTimeoutMS:   3 * MSPerSec,
	ServerReadTimeoutMS:    3 * MSPerSec,
	FileMemBytes:           bytesPerMebibyte * 100,
	DefaultCache:           "cache_mem_lru",
}

// LoadConfig loads the given config file. If an empty string is passed, the default config is returned.
func LoadConfig(fileName string) (Config, error) {
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
