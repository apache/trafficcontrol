package cache

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

// stats_type_astats is the default Stats format for Traffic Control.
// It is the Stats format produced by the `astats` plugin to Apache Traffic
// Server, included with Traffic Control.
//
// Stats are of the form `{"ats": {"name", number}}`,
// Where `name` is of the form:
//   `"plugin.remap_stats.fully-qualfiied-domain-name.example.net.stat-name"`
// Where `stat-name` is one of:
//   `in_bytes`, `out_bytes`, `status_2xx`, `status_3xx`, `status_4xx`, `status_5xx`

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/dsdata"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/poller"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/todata"
	jsoniter "github.com/json-iterator/go"
)

func init() {
	registerDecoder("astats", astatsParse, astatsPrecompute)
}

// AstatsSystem represents fixed system stats returned from the
// 'astats_over_http' ATS plugin.
type AstatsSystem struct {
	InfName           string `json:"inf.name"`
	InfSpeed          int    `json:"inf.speed"`
	ProcNetDev        string `json:"proc.net.dev"`
	ProcLoadavg       string `json:"proc.loadavg"`
	ConfigLoadRequest int    `json:"configReloadRequests"`
	LastReloadRequest int    `json:"lastReloadRequest"`
	ConfigReloads     int    `json:"configReloads"`
	LastReload        int    `json:"lastReload"`
	AstatsLoad        int    `json:"astatsLoad"`
	NotAvailable      bool   `json:"notAvailable,omitempty"`
}

// Astats contains ATS data returned from the Astats ATS plugin.
// This includes generic stats, as well as fixed system stats.
type Astats struct {
	Ats    map[string]interface{} `json:"ats"`
	System AstatsSystem           `json:"system"`
}

func astatsParse(cacheName string, rdr io.Reader, pollCTX interface{}) (Statistics, map[string]interface{}, error) {
	var stats Statistics
	if rdr == nil {
		log.Warnf("%s handle reader nil", cacheName)
		return stats, nil, errors.New("handler got nil reader")
	}

	ctx := pollCTX.(*poller.HTTPPollCtx)

	ctype := ctx.HTTPHeader.Get("Content-Type")

	if ctype == "text/json" || ctype == "text/javascript" || ctype == "application/json" || ctype == "" {
		var astats Astats
		json := jsoniter.ConfigFastest
		if err := json.NewDecoder(rdr).Decode(&astats); err != nil {
			return stats, nil, err
		}

		if err := stats.AddInterfaceFromRawLine(astats.System.ProcNetDev); err != nil {
			return stats, nil, fmt.Errorf("failed to parse interface line for cache '%s': %v", cacheName, err)
		}
		if inf, ok := stats.Interfaces[astats.System.InfName]; !ok {
			return stats, nil, errors.New("/proc/net/dev line didn't match reported interface line")
		} else {
			inf.Speed = int64(astats.System.InfSpeed)
			stats.Interfaces[astats.System.InfName] = inf
		}

		if load, err := LoadavgFromRawLine(astats.System.ProcLoadavg); err != nil {
			return stats, nil, fmt.Errorf("failed to parse loadavg line for cache '%s': %v", cacheName, err)
		} else {
			stats.Loadavg = load
		}

		stats.NotAvailable = astats.System.NotAvailable

		// TODO: what's using these?? Can we get rid of them?
		astats.Ats["system.astatsLoad"] = float64(astats.System.AstatsLoad)
		astats.Ats["system.configReloadRequests"] = float64(astats.System.ConfigLoadRequest)
		astats.Ats["system.configReloads"] = float64(astats.System.ConfigReloads)
		astats.Ats["system.inf.name"] = astats.System.InfName
		astats.Ats["system.inf.speed"] = float64(astats.System.InfSpeed)
		astats.Ats["system.lastReload"] = float64(astats.System.LastReload)
		astats.Ats["system.lastReloadRequest"] = float64(astats.System.LastReloadRequest)
		astats.Ats["system.notAvailable"] = stats.NotAvailable
		astats.Ats["system.proc.loadavg"] = astats.System.ProcLoadavg
		astats.Ats["system.proc.net.dev"] = astats.System.ProcNetDev

		return stats, astats.Ats, nil
	} else if ctype == "text/csv" {
		return astatsCsvParseCsv(cacheName, rdr)
	} else {
		return stats, nil, fmt.Errorf("stats Content-Type (%s) can not be parsed by astats", ctype)
	}
}

func astatsPrecompute(cacheName string, data todata.TOData, stats Statistics, miscStats map[string]interface{}) PrecomputedData {
	dsStats := make(map[string]*DSStat)
	var precomputed PrecomputedData
	precomputed.OutBytes = 0
	precomputed.MaxKbps = 0
	for _, iface := range stats.Interfaces {
		precomputed.OutBytes += iface.BytesOut
		kbps := iface.Speed * 1000
		if kbps > precomputed.MaxKbps {
			precomputed.MaxKbps = kbps
		}
	}

	var err error
	for stat, value := range miscStats {
		dsStats, err = astatsProcessStat(dsStats, data, stat, value)
		if err != nil && err != dsdata.ErrNotProcessedStat {
			log.Infof("precomputing cache %s stat %s value %v error %v", cacheName, stat, value, err)
			precomputed.Errors = append(precomputed.Errors, err)
			err = nil
		}
	}

	precomputed.DeliveryServiceStats = dsStats
	return precomputed
}

// astatsProcessStat and its subsidiary functions act as a State Machine,
// flowing the stat through states for each "." component of the stat name.
func astatsProcessStat(stats map[string]*DSStat, toData todata.TOData, stat string, value interface{}) (map[string]*DSStat, error) {
	parts := strings.Split(stat, ".")
	if len(parts) < 1 {
		return stats, fmt.Errorf("stat has no initial part")
	}

	switch parts[0] {
	case "plugin":
		return astatsProcessStatPlugin(stats, toData, parts[1:], value)
	case "proxy":
		fallthrough
	case "server":
		fallthrough
	case "system":
		return stats, dsdata.ErrNotProcessedStat
	default:
		return stats, fmt.Errorf("stat '%s' has unknown initial part '%s'", stat, parts[0])
	}
}

func astatsProcessStatPlugin(stats map[string]*DSStat, toData todata.TOData, statParts []string, value interface{}) (map[string]*DSStat, error) {
	if len(statParts) < 1 {
		return stats, fmt.Errorf("stat has no plugin part")
	}
	switch statParts[0] {
	case "remap_stats":
		return astatsProcessStatPluginRemapStats(stats, toData, statParts[1:], value)
	default:
		return stats, fmt.Errorf("stat has unknown plugin part '%s'", statParts[0])
	}
}

func astatsProcessStatPluginRemapStats(stats map[string]*DSStat, toData todata.TOData, statParts []string, value interface{}) (map[string]*DSStat, error) {
	if len(statParts) < 3 {
		return stats, fmt.Errorf("stat has no remap_stats deliveryservice and name parts")
	}

	// the FQDN is `subsubdomain`.`subdomain`.`domain`. For a HTTP Delivery
	// Service, `subsubdomain` will be the cache hostname; for a DNS Delivery
	// Service, it will be `edge`. Then, `subdomain` is the Delivery Service
	// regex.
	subsubdomain := statParts[0]
	subdomain := statParts[1]
	domain := strings.Join(statParts[2:len(statParts)-1], ".")

	ds, ok := toData.DeliveryServiceRegexes.DeliveryService(domain, subdomain, subsubdomain)
	if !ok {
		return stats, fmt.Errorf("no Delivery Service match for '%s.%s.%s' stat '%v'", subsubdomain, subdomain, domain, strings.Join(statParts, "."))
	}
	if ds == "" {
		return stats, fmt.Errorf("empty Delivery Service fqdn '%s.%s.%s' stat %v", subsubdomain, subdomain, domain, strings.Join(statParts, "."))
	}

	dsName := string(ds)

	statName := statParts[len(statParts)-1]
	if _, ok := stats[dsName]; !ok {
		stats[dsName] = new(DSStat)
	}

	dsStat := stats[dsName]

	if err := astatsAddCacheStat(dsStat, statName, value); err != nil {
		return stats, err
	}
	return stats, nil
}

// astatsAddCacheStat adds the given stat to the existing stat.
// Note this adds, it doesn't overwrite. Numbers are summed, strings are
// concatenated.
// TODO make this less duplicate code somehow.
func astatsAddCacheStat(stat *DSStat, name string, val interface{}) error {
	switch name {
	case "status_2xx":
		v, ok := val.(float64)
		if !ok {
			return fmt.Errorf("stat '%s' value expected float64 actual '%v' type %T", name, val, val)
		}
		stat.Status2xx += uint64(v)
	case "status_3xx":
		v, ok := val.(float64)
		if !ok {
			return fmt.Errorf("stat '%s' value expected float64 actual '%v' type %T", name, val, val)
		}
		stat.Status3xx += uint64(v)
	case "status_4xx":
		v, ok := val.(float64)
		if !ok {
			return fmt.Errorf("stat '%s' value expected float64 actual '%v' type %T", name, val, val)
		}
		stat.Status4xx += uint64(v)
	case "status_5xx":
		v, ok := val.(float64)
		if !ok {
			return fmt.Errorf("stat '%s' value expected float64 actual '%v' type %T", name, val, val)
		}
		stat.Status5xx += uint64(v)
	case "out_bytes":
		v, ok := val.(float64)
		if !ok {
			return fmt.Errorf("stat '%s' value expected float64 actual '%v' type %T", name, val, val)
		}
		stat.OutBytes += uint64(v)
	case "in_bytes":
		v, ok := val.(float64)
		if !ok {
			return fmt.Errorf("stat '%s' value expected float64 actual '%v' type %T", name, val, val)
		}
		stat.InBytes += uint64(v)
	case "status_unknown":
		return dsdata.ErrNotProcessedStat
	default:
		return fmt.Errorf("unknown stat '%s'", name)
	}
	return nil
}
