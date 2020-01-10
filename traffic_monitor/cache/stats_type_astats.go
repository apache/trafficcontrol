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
// It is the Stats format produced by the `astats` plugin to Apache Traffic Server, included with Traffic Control.
//
// Stats are of the form `{"ats": {"name", number}}`,
// Where `name` is of the form:
//   `"plugin.remap_stats.fully-qualfiied-domain-name.example.net.stat-name"`
// Where `stat-name` is one of:
//   `in_bytes`, `out_bytes`, `status_2xx`, `status_3xx`, `status_4xx`, `status_5xx`

import (
	"errors"
	"fmt"
	"github.com/apache/trafficcontrol/lib/go-tc/enum"
	"io"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/traffic_monitor/dsdata"
	"github.com/apache/trafficcontrol/traffic_monitor/todata"
	"github.com/json-iterator/go"
)

func init() {
	AddStatsType("astats", astatsParse, astatsPrecompute)
}

func astatsParse(cache enum.CacheName, rdr io.Reader) (error, map[string]interface{}, AstatsSystem) {
	if rdr == nil {
		log.Warnln(string(cache) + " handle reader nil")
		return errors.New("handler got nil reader"), nil, AstatsSystem{}
	}
	astats := Astats{}

	json := jsoniter.ConfigFastest // TODo make configurable?
	err := json.NewDecoder(rdr).Decode(&astats)
	return err, astats.Ats, astats.System
}

func astatsPrecompute(cache enum.CacheName, toData todata.TOData, rawStats map[string]interface{}, system AstatsSystem) PrecomputedData {
	stats := map[enum.DeliveryServiceName]*AStat{}

	precomputed := PrecomputedData{}
	var err error
	if precomputed.OutBytes, err = astatsOutBytes(system.ProcNetDev, system.InfName); err != nil {
		precomputed.OutBytes = 0
		log.Errorf("precomputeAstats %s handle precomputing outbytes '%v'\n", cache, err)
	}

	kbpsInMbps := int64(1000)
	precomputed.MaxKbps = int64(system.InfSpeed) * kbpsInMbps

	for stat, value := range rawStats {
		stats, err = astatsProcessStat(cache, stats, toData, stat, value)
		if err != nil && err != dsdata.ErrNotProcessedStat {
			log.Infof("precomputing cache %v stat %v value %v error %v", cache, stat, value, err)
			precomputed.Errors = append(precomputed.Errors, err)
		}
	}
	precomputed.DeliveryServiceStats = stats
	return precomputed
}

// outBytes takes the proc.net.dev string, and the interface name, and returns the bytes field
func astatsOutBytes(procNetDev, iface string) (int64, error) {
	if procNetDev == "" {
		return 0, fmt.Errorf("procNetDev empty")
	}
	if iface == "" {
		return 0, fmt.Errorf("iface empty")
	}
	ifacePos := strings.Index(procNetDev, iface)
	if ifacePos == -1 {
		return 0, fmt.Errorf("interface '%s' not found in proc.net.dev '%s'", iface, procNetDev)
	}

	procNetDevIfaceBytes := procNetDev[ifacePos+len(iface)+1:]
	procNetDevIfaceBytesArr := strings.Fields(procNetDevIfaceBytes) // TODO test
	if len(procNetDevIfaceBytesArr) < 10 {
		return 0, fmt.Errorf("proc.net.dev iface '%v' unknown format '%s'", iface, procNetDev)
	}
	procNetDevIfaceBytes = procNetDevIfaceBytesArr[8]

	return strconv.ParseInt(procNetDevIfaceBytes, 10, 64)
}

// astatsProcessStat and its subsidiary functions act as a State Machine, flowing the stat thru states for each "." component of the stat name
func astatsProcessStat(server enum.CacheName, stats map[enum.DeliveryServiceName]*AStat, toData todata.TOData, stat string, value interface{}) (map[enum.DeliveryServiceName]*AStat, error) {
	parts := strings.Split(stat, ".")
	if len(parts) < 1 {
		return stats, fmt.Errorf("stat has no initial part")
	}

	switch parts[0] {
	case "plugin":
		return astatsProcessStatPlugin(server, stats, toData, stat, parts[1:], value)
	case "proxy":
		return stats, dsdata.ErrNotProcessedStat
	case "server":
		return stats, dsdata.ErrNotProcessedStat
	default:
		return stats, fmt.Errorf("stat '%s' has unknown initial part '%s'", stat, parts[0])
	}
}

func astatsProcessStatPlugin(server enum.CacheName, stats map[enum.DeliveryServiceName]*AStat, toData todata.TOData, stat string, statParts []string, value interface{}) (map[enum.DeliveryServiceName]*AStat, error) {
	if len(statParts) < 1 {
		return stats, fmt.Errorf("stat has no plugin part")
	}
	switch statParts[0] {
	case "remap_stats":
		return astatsProcessStatPluginRemapStats(server, stats, toData, stat, statParts[1:], value)
	default:
		return stats, fmt.Errorf("stat has unknown plugin part '%s'", statParts[0])
	}
}

func astatsProcessStatPluginRemapStats(server enum.CacheName, stats map[enum.DeliveryServiceName]*AStat, toData todata.TOData, stat string, statParts []string, value interface{}) (map[enum.DeliveryServiceName]*AStat, error) {
	if len(statParts) < 3 {
		return stats, fmt.Errorf("stat has no remap_stats deliveryservice and name parts")
	}

	// the FQDN is `subsubdomain`.`subdomain`.`domain`. For a HTTP delivery service, `subsubdomain` will be the cache hostname; for a DNS delivery service, it will be `edge`. Then, `subdomain` is the delivery service regex.
	subsubdomain := statParts[0]
	subdomain := statParts[1]
	domain := strings.Join(statParts[2:len(statParts)-1], ".")

	ds, ok := toData.DeliveryServiceRegexes.DeliveryService(domain, subdomain, subsubdomain)
	if !ok {
		fqdn := fmt.Sprintf("%s.%s.%s", subsubdomain, subdomain, domain)
		return stats, fmt.Errorf("ERROR no delivery service match for fqdn '%v' stat '%v'\n", fqdn, strings.Join(statParts, "."))
	}
	if ds == "" {
		fqdn := fmt.Sprintf("%s.%s.%s", subsubdomain, subdomain, domain)
		return stats, fmt.Errorf("ERROR EMPTY delivery service fqdn %v stat %v\n", fqdn, strings.Join(statParts, "."))
	}

	statName := statParts[len(statParts)-1]

	dsStat, ok := stats[ds]
	if !ok {
		dsStat = &AStat{}
		stats[ds] = dsStat
	}

	if err := astatsAddCacheStat(dsStat, statName, value); err != nil {
		return stats, err
	}
	stats[ds] = dsStat // TODO verify unnecessary, remove
	return stats, nil
}

// addCacheStat adds the given stat to the existing stat. Note this adds, it doesn't overwrite. Numbers are summed, strings are concatenated.
// TODO make this less duplicate code somehow.
func astatsAddCacheStat(stat *AStat, name string, val interface{}) error {
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
