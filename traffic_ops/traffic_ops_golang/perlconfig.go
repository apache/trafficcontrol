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
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/basho/riak-go-client"
)

const OldAccessLogPath = "/var/log/traffic_ops/access.log"
const NewLogPath = "/var/log/traffic_ops/traffic_ops_golang.log"
const DefaultMaxDBConnections = 50
const DefaultProxyTimeout = 60
const DefaultProxyTLSTimeout = 60
const DefaultProxyReadHeaderTimeout = 60
const DefaultProxyKeepAlive = 60
const DefaultReadTimeout = 60
const DefaultReadHeaderTimeout = 60
const DefaultWriteTimeout = 60
const DefaultIdleTimeout = 60

func GetPerlConfigs(cdnConfPath string, dbConfPath string, riakConfPath string) (Config, error) {
	configBytes, err := ioutil.ReadFile(cdnConfPath)
	if err != nil {
		return Config{}, fmt.Errorf("reading CDN conf '%v': %v", cdnConfPath, err)
	}
	dbConfBytes, err := ioutil.ReadFile(dbConfPath)
	if err != nil {
		return Config{}, fmt.Errorf("reading db conf '%v': %v", dbConfPath, err)
	}
	riakConfBytes, err := ioutil.ReadFile(riakConfPath)
	if err != nil {
		return Config{}, fmt.Errorf("reading riak conf '%v': %v", riakConfPath, err)
	}
	return getPerlConfigsFromStrs(string(configBytes), string(dbConfBytes), string(riakConfBytes))
}

func getPerlConfigsFromStrs(cdnConfBytes string, dbConfBytes string, riakConfBytes string) (Config, error) {
	cfg, err := getCDNConf(cdnConfBytes)
	if err != nil {
		return Config{}, fmt.Errorf("parsing CDN conf '%v': %v", cdnConfBytes, err)
	}

	dbconf, err := getDbConf(string(dbConfBytes))
	if err != nil {
		return Config{}, fmt.Errorf("parsing db conf '%v': %v", dbConfBytes, err)
	}
	cfg.DBUser = dbconf.User
	cfg.DBPass = dbconf.Password
	cfg.DBServer = dbconf.Hostname
	cfg.DBDB = dbconf.DBName
	cfg.DBSSL = false // TODO fix
	if dbconf.Port != "" {
		cfg.DBServer += ":" + dbconf.Port
	}

	cfg.LogLocationInfo = NewLogPath
	cfg.LogLocationError = NewLogPath
	cfg.LogLocationWarning = NewLogPath
	cfg.LogLocationEvent = OldAccessLogPath
	cfg.LogLocationDebug = log.LogLocationNull

	riakconf, err := getRiakAuthOptions(riakConfBytes)
	if err != nil {
		return Config{}, fmt.Errorf("parsing riak conf '%v': %v", riakConfBytes, err)
	}
	cfg.RiakAuthOptions = riakconf

	return cfg, nil
}

func getCDNConf(s string) (Config, error) {
	cfg := Config{}
	obj, err := ParsePerlObj(s)
	if err != nil {
		return Config{}, fmt.Errorf("parsing Perl object: %v", err)
	}

	if cfg.HTTPPort, err = getPort(obj); err != nil {
		return Config{}, err
	}

	if cfg.TOSecret, err = getSecret(obj); err != nil {
		return Config{}, err
	}

	oldPort, err := getOldPort(obj)
	if err != nil {
		return Config{}, err
	}
	cfg.TOURLStr = "https://127.0.0.1:" + oldPort
	if cfg.TOURL, err = url.Parse(cfg.TOURLStr); err != nil {
		return Config{}, fmt.Errorf("Invalid Traffic Ops URL '%v': %s", cfg.TOURL, err)
	}

	cfg.CertPath, err = getConfigCert(obj)
	if err != nil {
		return Config{}, err
	}

	cfg.KeyPath, err = getConfigKey(obj)
	if err != nil {
		return Config{}, err
	}

	if proxyTimeout, err := getIntFromConfigurationKey("traffic_ops_golang_proxy_timeout", obj); err != nil {
		log.Warnf("failed to get proxy timeout from cdn.conf (%s), using default %d\n", err, DefaultProxyTimeout)
		cfg.ProxyTimeout = DefaultProxyTimeout
	} else {
		cfg.ProxyTimeout = proxyTimeout
	}

	if proxyKeepAlive, err := getIntFromConfigurationKey("traffic_ops_golang_proxy_keep_alive", obj); err != nil {
		log.Warnf("failed to get proxy keepAlive from cdn.conf (%s), using default %d\n", err, DefaultProxyKeepAlive)
		cfg.ProxyKeepAlive = DefaultProxyKeepAlive
	} else {
		cfg.ProxyKeepAlive = proxyKeepAlive
	}

	if proxyTLSTimeout, err := getIntFromConfigurationKey("traffic_ops_golang_proxy_tls_timeout", obj); err != nil {
		log.Warnf("failed to get proxy TLSTimeout from cdn.conf (%s), using default %d\n", err, DefaultProxyTLSTimeout)
		cfg.ProxyTLSTimeout = DefaultProxyTLSTimeout
	} else {
		cfg.ProxyTLSTimeout = proxyTLSTimeout
	}

	if proxyReadHeaderTimeout, err := getIntFromConfigurationKey("traffic_ops_golang_proxy_read_header_timeout", obj); err != nil {
		log.Warnf("failed to get proxy readHeaderTimeout from cdn.conf (%s), using default %d\n", err, DefaultProxyReadHeaderTimeout)
		cfg.ProxyReadHeaderTimeout = DefaultProxyReadHeaderTimeout
	} else {
		cfg.ProxyReadHeaderTimeout = proxyReadHeaderTimeout
	}

	if readTimeout, err := getIntFromConfigurationKey("traffic_ops_golang_read_timeout", obj); err != nil {
		log.Warnf("failed to get read timeout from cdn.conf (%s), using default %d\n", err, DefaultReadTimeout)
		cfg.ReadTimeout = DefaultReadTimeout
	} else {
		cfg.ReadTimeout = readTimeout
	}

	if readHeaderTimeout, err := getIntFromConfigurationKey("traffic_ops_golang_read_header_timeout", obj); err != nil {
		log.Warnf("failed to get readHeader timeout from cdn.conf (%s), using default %d\n", err, DefaultReadHeaderTimeout)
		cfg.ReadHeaderTimeout = DefaultReadHeaderTimeout
	} else {
		cfg.ReadHeaderTimeout = readHeaderTimeout
	}

	if writeTimeout, err := getIntFromConfigurationKey("traffic_ops_golang_write_timeout", obj); err != nil {
		log.Warnf("failed to get write timeout from cdn.conf (%s), using default %d\n", err, DefaultWriteTimeout)
		cfg.WriteTimeout = DefaultWriteTimeout
	} else {
		cfg.WriteTimeout = writeTimeout
	}

	if idleTimeout, err := getIntFromConfigurationKey("traffic_ops_golang_idle_timeout", obj); err != nil {
		log.Warnf("failed to get idle timeout from cdn.conf (%s), using default %d\n", err, DefaultIdleTimeout)
		cfg.IdleTimeout = DefaultIdleTimeout
	} else {
		cfg.IdleTimeout = idleTimeout
	}

	if dbMaxConns, err := getIntFromConfigurationKey("traffic_ops_golang_max_db_connections", obj); err != nil {
		log.Warnf("failed to get Max DB Connections from cdn.conf (%s), using default %d\n", err, DefaultMaxDBConnections)
		cfg.MaxDBConnections = DefaultMaxDBConnections
	} else {
		cfg.MaxDBConnections = dbMaxConns
	}

	return cfg, nil
}

func getPort(obj map[string]interface{}) (string, error) {
	portStrI, ok := obj["traffic_ops_golang_port"]
	if !ok {
		return "", fmt.Errorf("missing traffic_ops_golang_port key")
	}
	portStr, ok := portStrI.(string)
	if !ok {
		return "", fmt.Errorf("traffic_ops_golang_port key '%v' not a string", portStrI)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil || port < 0 || port > 65535 {
		return "", fmt.Errorf("invalid port '%s'", portStr)
	}
	return strconv.Itoa(port), nil
}

func getIntFromConfigurationKey(key string, obj map[string]interface{}) (int, error) {
	confNum, ok := obj[key]
	if !ok {
		return 0, fmt.Errorf("missing %s key", key)
	}
	confFloat, ok := confNum.(float64)
	if !ok {
		return 0, fmt.Errorf("key: %s's value: '%v' type %T not a number", key, confNum, confNum)
	}
	return int(confFloat), nil
}

func getOldPort(obj map[string]interface{}) (string, error) {
	hypnotoadI, ok := obj["hypnotoad"]
	if !ok {
		return "", fmt.Errorf("missing hypnotoad key")
	}
	hypnotoad, ok := hypnotoadI.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("hypnotoad key '%v' not an object", hypnotoadI)
	}

	listenArrI, ok := hypnotoad["listen"]
	if !ok {
		return "", fmt.Errorf("missing hypnotoad.listen key")
	}
	listenArr, ok := listenArrI.([]interface{})
	if !ok {
		return "", fmt.Errorf("listen key '%v' type %T not an array", listenArrI, listenArrI)
	}
	if len(listenArr) < 1 {
		return "", fmt.Errorf("empty hypnotoad.listen key")
	}
	listenI := listenArr[0]
	listen, ok := listenI.(string)
	if !ok {
		return "", fmt.Errorf("listen[0] key '%v' type %T not a string", listenI, listenI)
	}

	listenRe := regexp.MustCompile(`:(\d+)`)
	portMatch := listenRe.FindStringSubmatch(listen)
	if len(portMatch) < 2 {
		return "", fmt.Errorf("failed to find port in listen '%s'", listen)
	}
	portStr := portMatch[1]

	port, err := strconv.Atoi(portStr)
	if err != nil || port < 0 || port > 65535 {
		return "", fmt.Errorf("invalid port in listen '%s'", listen)
	}
	return strconv.Itoa(port), nil
}

func getConfigCert(obj map[string]interface{}) (string, error) {
	hypnotoadI, ok := obj["hypnotoad"]
	if !ok {
		return "", fmt.Errorf("missing hypnotoad key")
	}
	hypnotoad, ok := hypnotoadI.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("hypnotoad key '%v' not an object", hypnotoadI)
	}

	listenArrI, ok := hypnotoad["listen"]
	if !ok {
		return "", fmt.Errorf("missing hypnotoad.listen key")
	}
	listenArr, ok := listenArrI.([]interface{})
	if !ok {
		return "", fmt.Errorf("listen key '%v' type %T not an array", listenArrI, listenArrI)
	}
	if len(listenArr) < 1 {
		return "", fmt.Errorf("empty hypnotoad.listen key")
	}
	listenI := listenArr[0]
	listen, ok := listenI.(string)
	if !ok {
		return "", fmt.Errorf("listen[0] key '%v' type %T not a string", listenI, listenI)
	}

	keyStr := "cert="
	start := strings.Index(listen, keyStr)
	if start < 0 {
		return "", fmt.Errorf("failed to find key in listen '%s'", listen)
	}
	listen = listen[start+len(keyStr):]
	end := strings.Index(listen, "&")
	if end < 0 {
		return listen[start:], nil
	}
	return listen[:end], nil
}

func getConfigKey(obj map[string]interface{}) (string, error) {
	hypnotoadI, ok := obj["hypnotoad"]
	if !ok {
		return "", fmt.Errorf("missing hypnotoad key")
	}
	hypnotoad, ok := hypnotoadI.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("hypnotoad key '%v' not an object", hypnotoadI)
	}

	listenArrI, ok := hypnotoad["listen"]
	if !ok {
		return "", fmt.Errorf("missing hypnotoad.listen key")
	}
	listenArr, ok := listenArrI.([]interface{})
	if !ok {
		return "", fmt.Errorf("listen key '%v' type %T not an array", listenArrI, listenArrI)
	}
	if len(listenArr) < 1 {
		return "", fmt.Errorf("empty hypnotoad.listen key")
	}
	listenI := listenArr[0]
	listen, ok := listenI.(string)
	if !ok {
		return "", fmt.Errorf("listen[0] key '%v' type %T not a string", listenI, listenI)
	}

	keyStr := "key="
	start := strings.Index(listen, keyStr)
	if start < 0 {
		return "", fmt.Errorf("failed to find key in listen '%s'", listen)
	}
	listen = listen[start+len(keyStr):]
	end := strings.Index(listen, "&")
	if end < 0 {
		return listen[start:], nil
	}
	return listen[:end], nil
}

func getSecret(obj map[string]interface{}) (string, error) {
	secretsI, ok := obj["secrets"]
	if !ok {
		return "", fmt.Errorf("missing secrets key")
	}
	secrets, ok := secretsI.([]interface{})
	if !ok {
		return "", fmt.Errorf("secrets key '%v' not an array", secretsI)
	}

	if len(secrets) < 1 {
		return "", fmt.Errorf("empty secrets key")
	}
	secretI := secrets[0]
	secret, ok := secretI.(string)
	if !ok {
		return "", fmt.Errorf("secret '%v' not a string", secretI)
	}

	return secret, nil
}

type DatabaseConf struct {
	Description string `json:"description"`
	DBName      string `json:"dbname"`
	Hostname    string `json:"hostname"`
	User        string `json:"user"`
	Password    string `json:"password"`
	Port        string `json:"port"`
	Type        string `json:"type"`
}

func getDbConf(s string) (DatabaseConf, error) {
	dbc := DatabaseConf{}
	err := json.Unmarshal([]byte(s), &dbc)
	return dbc, err
}

func getRiakAuthOptions(s string) (*riak.AuthOptions, error) {
	rconf := &riak.AuthOptions{}
	rconf.TlsConfig = &tls.Config{}
	err := json.Unmarshal([]byte(s), &rconf)
	return rconf, err
}
