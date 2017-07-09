package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

type Config struct {
	HTTPPort string   `json:"port"`
	DBUser   string   `json:"db_user"`
	DBPass   string   `json:"db_pass"`
	DBServer string   `json:"db_server"`
	DBDB     string   `json:"db_name"`
	DBSSL    bool     `json:"db_ssl"`
	TOSecret string   `json:"to_secret"`
	TOURLStr string   `json:"to_url"`
	TOURL    *url.URL `json:"-"`
	NoAuth   bool     `json:"no_auth"`
	CertPath string   `json:"cert_path"`
	KeyPath  string   `json:"key_path"`
}

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
	if cfg.HTTPPort == "" {
		return Config{}, fmt.Errorf("missing port")
	}
	if cfg.DBUser == "" {
		return Config{}, fmt.Errorf("missing database user")
	}
	if cfg.DBPass == "" {
		return Config{}, fmt.Errorf("missing database password")
	}
	if cfg.DBServer == "" {
		return Config{}, fmt.Errorf("missing database server")
	}
	if cfg.DBDB == "" {
		return Config{}, fmt.Errorf("missing database name")
	}
	if cfg.TOSecret == "" {
		return Config{}, fmt.Errorf("missing secret")
	}
	if cfg.CertPath == "" {
		return Config{}, fmt.Errorf("missing certificate path")
	}
	if cfg.KeyPath == "" {
		return Config{}, fmt.Errorf("missing certificate key path")
	}

	var err error
	if cfg.TOURL, err = url.Parse(cfg.TOURLStr); err != nil {
		return Config{}, fmt.Errorf("Invalid Traffic Ops URL '%v': err", cfg.TOURL, err)
	}

	return cfg, nil
}
