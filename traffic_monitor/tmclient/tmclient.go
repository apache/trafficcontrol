package tmclient

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
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/datareq"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/dsdata"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/handler"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/towrap"
)

type TMClient struct {
	url       string
	timeout   time.Duration
	Transport *http.Transport // optional http Transport
}

func New(url string, timeout time.Duration) *TMClient {
	return &TMClient{url: strings.TrimSuffix(url, "/"), timeout: timeout}
}

func (c *TMClient) CacheCount() (int, error) { return c.getInt("/api/cache-count") }

func (c *TMClient) CacheAvailableCount() (int, error) { return c.getInt("/api/cache-available-count") }

func (c *TMClient) CacheDownCount() (int, error) { return c.getInt("/api/cache-down-count") }

func (c *TMClient) Version() (string, error) { return c.getStr("/api/version") }

func (c *TMClient) TrafficOpsURI() (string, error) { return c.getStr("/api/traffic-ops-uri") }

func (c *TMClient) BandwidthKBPS() (float64, error) { return c.getFloat("/api/bandwidth-kbps") }

func (c *TMClient) BandwidthCapacityKBPS() (float64, error) {
	return c.getFloat("/api/bandwidth-capacity-kbps")
}

func (c *TMClient) CacheStatuses() (map[tc.CacheName]datareq.CacheStatus, error) {
	path := "/api/cache-statuses"
	obj := map[tc.CacheName]datareq.CacheStatus{}
	if err := c.GetJSON(path, &obj); err != nil {
		return nil, err // GetJSON adds context
	}
	return obj, nil
}

func (c *TMClient) MonitorConfig() (tc.TrafficMonitorConfigMap, error) {
	path := "/api/monitor-config"
	obj := tc.TrafficMonitorConfigMap{}
	if err := c.GetJSON(path, &obj); err != nil {
		return tc.TrafficMonitorConfigMap{}, err // GetJSON adds context
	}
	return obj, nil
}

func (c *TMClient) CRConfigHistory() ([]towrap.CRConfigStat, error) {
	path := "/api/crconfig-history"
	obj := []towrap.CRConfigStat{}
	if err := c.GetJSON(path, &obj); err != nil {
		return nil, err // GetJSON adds context
	}
	return obj, nil
}

func (c *TMClient) EventLog() (datareq.JSONEvents, error) {
	path := "/publish/EventLog"
	obj := datareq.JSONEvents{}
	if err := c.GetJSON(path, &obj); err != nil {
		return datareq.JSONEvents{}, err // GetJSON adds context
	}
	return obj, nil
}

func (c *TMClient) CacheStatsNew() (tc.Stats, error) {
	path := "/publish/CacheStats"
	obj := tc.Stats{}
	if err := c.GetJSON(path, &obj); err != nil {
		return tc.Stats{}, err // GetJSON adds context
	}
	return obj, nil
}

func (c *TMClient) CacheStats() (tc.LegacyStats, error) {
	path := "/publish/CacheStats"
	obj := tc.LegacyStats{}
	if err := c.GetJSON(path, &obj); err != nil {
		return tc.LegacyStats{}, err // GetJSON adds context
	}
	return obj, nil
}

func (c *TMClient) DSStats() (dsdata.Stats, error) {
	path := "/publish/DsStats"
	obj := dsdata.Stats{}
	if err := c.GetJSON(path, &obj); err != nil {
		return dsdata.Stats{}, err // GetJSON adds context
	}
	return obj, nil
}

func (c *TMClient) CRStates(raw bool) (tc.CRStates, error) {
	path := "/publish/CrStates"
	if raw {
		path += "?raw"
	}
	obj := tc.CRStates{}
	if err := c.GetJSON(path, &obj); err != nil {
		return tc.CRStates{}, err // GetJSON adds context
	}
	return obj, nil
}

func (c *TMClient) CRConfig() (tc.CRConfig, error) {
	path := "/publish/CrConfig"
	obj := tc.CRConfig{}
	if err := c.GetJSON(path, &obj); err != nil {
		return tc.CRConfig{}, err // GetJSON adds context
	}
	return obj, nil
}

// CRConfigBytes returns the raw bytes of the Monitor's CRConfig.
//
// If you need a deserialized object, use TMClient.CRConfig() instead.
//
// This function exists because the Monitor very intentionally serves the CRConfig bytes as
// published by Traffic Ops, without deserializing or reserializing it.
//
// This can be useful to check for serialization or versioning issues, in case the Go object
// is missing values sent by Traffic Ops, or has other serialization issues.
func (c *TMClient) CRConfigBytes() ([]byte, error) { return c.getBytes("publish/CrConfig") }

func (c *TMClient) PeerStates() (datareq.APIPeerStates, error) {
	path := "/publish/PeerStates"
	obj := datareq.APIPeerStates{}
	if err := c.GetJSON(path, &obj); err != nil {
		return datareq.APIPeerStates{}, err // GetJSON adds context
	}
	return obj, nil
}

func (c *TMClient) Stats() (datareq.Stats, error) {
	path := "/publish/Stats"
	obj := datareq.Stats{}
	if err := c.GetJSON(path, &obj); err != nil {
		return datareq.Stats{}, err // GetJSON adds context
	}
	return obj, nil
}

func (c *TMClient) StatSummary() (datareq.StatSummary, error) {
	path := "/publish/StatSummary"
	obj := datareq.StatSummary{}
	if err := c.GetJSON(path, &obj); err != nil {
		return datareq.StatSummary{}, err // GetJSON adds context
	}
	return obj, nil
}

func (c *TMClient) ConfigDoc() (handler.OpsConfig, error) {
	path := "/publish/ConfigDoc"
	obj := handler.OpsConfig{}
	if err := c.GetJSON(path, &obj); err != nil {
		return handler.OpsConfig{}, err // GetJSON adds context
	}
	return obj, nil
}

func (c *TMClient) getBytes(path string) ([]byte, error) {
	url := c.url + path
	httpClient := http.Client{Timeout: c.timeout}
	if c.Transport != nil {
		httpClient.Transport = c.Transport
	}
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, errors.New("getting from '" + url + "': " + err.Error())
	}
	defer log.Close(resp.Body, "Unable to close http client "+url)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("monitor='"+url+"' monitor_status=%v event=\"error in TrafficMonitor polling returned bad status\"", resp.StatusCode)
	}

	respBts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("reading body from '" + url + "': " + err.Error())
	}
	return respBts, nil
}

func (c *TMClient) GetJSON(path string, obj interface{}) error {
	bts, err := c.getBytes(path)
	if err != nil {
		return err // getBytes already adds context
	}
	if err := json.Unmarshal(bts, obj); err != nil {
		return errors.New("unmarshalling response '" + string(bts) + "' json: " + err.Error())
	}
	return nil
}

func (c *TMClient) getStr(path string) (string, error) {
	respBts, err := c.getBytes(path)
	if err != nil {
		return "", err // getBytes already adds context
	}
	return string(respBts), nil
}

func (c *TMClient) getInt(path string) (int, error) {
	respStr, err := c.getStr(path)
	if err != nil {
		return 0, err // getStr already adds context
	}

	respInt, err := strconv.Atoi(respStr)
	if err != nil {
		return 0, errors.New("parsing response '" + respStr + "': " + err.Error())
	}
	return respInt, nil
}

func (c *TMClient) getFloat(path string) (float64, error) {
	respStr, err := c.getStr(path)
	if err != nil {
		return 0, err // getStr already adds context
	}

	respFloat, err := strconv.ParseFloat(respStr, 64)
	if err != nil {
		return 0, errors.New("parsing response '" + respStr + "': " + err.Error())
	}
	return respFloat, nil
}
