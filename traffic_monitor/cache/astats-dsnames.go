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

// stats_type_astats_dsnames is a Stat format similar to the default Astats,
// but with Fully Qualified Domain Names replaced with Delivery Service XMLIDs.
// That is, stat names are of the form:
// `"plugin.remap_stats.delivery-service-xmlid.stat-name"`

import (
	"fmt"
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/dsdata"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/todata"
)

func init() {
	registerDecoder("astats-dsnames", astatsParse, astatsdsnamesPrecompute)
}

func astatsdsnamesPrecompute(cache string, toData todata.TOData, stats Statistics, rawStats map[string]interface{}) PrecomputedData {
	dsStats := make(map[string]*DSStat)
	var precomputed PrecomputedData
	precomputed.OutBytes = 0
	precomputed.MaxKbps = 0
	for _, iface := range stats.Interfaces {
		precomputed.OutBytes += iface.BytesOut
		if iface.Speed > precomputed.MaxKbps {
			precomputed.MaxKbps = iface.Speed
		}
	}
	precomputed.MaxKbps *= 1000

	for stat, value := range rawStats {
		var err error
		dsStats, err = astatsdsnamesProcessStat(cache, dsStats, toData, stat, value)
		if err != nil && err != dsdata.ErrNotProcessedStat {
			log.Infof("precomputing cache %v stat %v value %v error %v", cache, stat, value, err)
			precomputed.Errors = append(precomputed.Errors, err)
		}
	}
	precomputed.DeliveryServiceStats = dsStats
	return precomputed
}

// astatsdsnamesProcessStat and its subsidiary functions act as a State Machine,
// flowing the stat through states for each "." component of the stat name.
func astatsdsnamesProcessStat(server string, stats map[string]*DSStat, toData todata.TOData, stat string, value interface{}) (map[string]*DSStat, error) {
	parts := strings.Split(stat, ".")
	if len(parts) < 1 {
		return stats, fmt.Errorf("stat has no initial part")
	}

	switch parts[0] {
	case "plugin":
		return astatsdsnamesProcessStatPlugin(server, stats, toData, stat, parts[1:], value)
	case "proxy":
		return stats, dsdata.ErrNotProcessedStat
	case "server":
		return stats, dsdata.ErrNotProcessedStat
	default:
		return stats, fmt.Errorf("stat '%s' has unknown initial part '%s'", stat, parts[0])
	}
}

func astatsdsnamesProcessStatPlugin(server string, stats map[string]*DSStat, toData todata.TOData, stat string, statParts []string, value interface{}) (map[string]*DSStat, error) {
	if len(statParts) < 1 {
		return stats, fmt.Errorf("stat has no plugin part")
	}
	switch statParts[0] {
	case "remap_stats":
		return astatsdsnamesProcessStatPluginRemapStats(server, stats, toData, stat, statParts[1:], value)
	default:
		return stats, fmt.Errorf("stat has unknown plugin part '%s'", statParts[0])
	}
}

func astatsdsnamesProcessStatPluginRemapStats(server string, stats map[string]*DSStat, toData todata.TOData, stat string, statParts []string, value interface{}) (map[string]*DSStat, error) {
	if len(statParts) < 3 {
		return stats, fmt.Errorf("stat has no remap_stats deliveryservice and name parts")
	}

	ds := statParts[0]
	statName := statParts[len(statParts)-1]

	if _, ok := toData.DeliveryServiceTypes[tc.DeliveryServiceName(ds)]; !ok {
		return stats, fmt.Errorf("no delivery service match for name '%v' stat '%v'\n", ds, statName)
	}

	dsStat := stats[ds]
	if err := astatsdstypesAddCacheStat(dsStat, statName, value); err != nil {
		return stats, err
	}

	stats[ds] = dsStat
	return stats, nil
}

// astatsdstypesAddCacheStat adds the given stat to the existing stat. Note this
// adds, it doesn't overwrite. Numbers are summed, strings are concatenated.
func astatsdstypesAddCacheStat(stat *DSStat, name string, val interface{}) error {
	// TODO make this less duplicate code somehow.
	// NOTE this is superficially duplicated from astatsAddCacheStat, but they are conceptually different, because the `astats` format changing should not necessarily affect the `astats-dstypes` format. The MUST be kept separate, and code between them MUST NOT be de-duplicated.
	switch name {
	case "status_2xx":
		v, ok := val.(float64)
		if !ok {
			return fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.Status2xx += uint64(v)
	case "status_3xx":
		v, ok := val.(float64)
		if !ok {
			return fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.Status3xx += uint64(v)
	case "status_4xx":
		v, ok := val.(float64)
		if !ok {
			return fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.Status4xx += uint64(v)
	case "status_5xx":
		v, ok := val.(float64)
		if !ok {
			return fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.Status5xx += uint64(v)
	case "out_bytes":
		v, ok := val.(float64)
		if !ok {
			return fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.OutBytes += uint64(v)
	case "in_bytes":
		v, ok := val.(float64)
		if !ok {
			return fmt.Errorf("stat '%s' value expected int actual '%v' type %T", name, val, val)
		}
		stat.InBytes += uint64(v)
	case "status_unknown":
		return dsdata.ErrNotProcessedStat
	default:
		return fmt.Errorf("unknown stat '%s'", name)
	}
	return nil
}
