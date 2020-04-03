// Package cache contains definitions for mechanisms used to extract health
// and statistics data from cache-server-provided data. The most commonly
// used format is the ``stats_over_http'' format provided by the plugin of the
// same name for Apache Traffic Server, followed closely by ``astats''  which
// is the legacy format used by older versions of Apache Traffic Control.
//
// Creating A New Stats Type
//
// To create a new Stats Type, for a custom caching proxy with its own stats
// format:
//
// 1. Create a file for your type in the traffic_monitor/cache directory and
//    package, `github.com/apache/trafficcontrol/traffic_monitor/cache/`
// 2. Create Parse and (optionally) Precompute functions in your file, with the
//     signature of `StatsTypeParser` and `StatsTypePrecomputer`
// 3. In your file, add
//    `func init(){AddStatsType(myTypeParser, myTypePrecomputer})`
//
// Your Parser should take the raw bytes from the `io.Reader` and populate the
// raw stats from them. For maximum compatibility, the names of these should be
// of the same form as Apache Traffic Server's `stats_over_http`, of the form
// "plugin.remap_stats.delivery-service-fqdn.com.in_bytes" et cetera. Traffic
// Control _may_ work with custom stat names, but we don't currently guarantee
// it.
//
// Your Precomputer should take the Stats and System information your Parser
// created, and populate the PrecomputedData. It is essential that all
// PrecomputedData fields are populated, especially `DeliveryServiceStats`,
// as they are used for cache and delivery service availability and threshold
// computation. If PrecomputedData is not properly and fully populated, the
// cache's availability will not be properly computed.
//
// Note the PrecomputedData `Reporting` and `Time` fields are the exception:
// they do not need to be set, and will be forcibly overridden by the Handler
// after your Precomputer function returns.
//
// Note these functions will not be called for Health polls, only Stat polls.
// Your Cache should have two separate stats endpoints: a small light endpoint
// returning only system stats and used to quickly verify reachability, and a
// large endpoint with all stats. If your cache does not have two stat
// endpoints, you may use your large stat endpoint for the Health poll, and
// configure the Health poll interval to be arbitrarily slow.
//
// Note your stats functions SHOULD NOT reuse functions from other stats types,
// even if they are similar, or have identical helper functions. This is a case
// where "duplicate" code is acceptable, because it's not conceptually
// duplicate. You don't want your stat parsers to break if the similar stats
// format you reuse code from changes.
//
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

import "encoding/json"
import "errors"
import "fmt"
import "io"
import "math"
import "strings"
import "strconv"

import "github.com/apache/trafficcontrol/lib/go-log"

func init() {
	// AddStatsType("stats_over_http", statsParse, statsPrecompute)
	registerDecoder("stats_over_http", parseStats)
}

type stats_over_httpData struct {
	Global map[string]interface{} `json:"global"`
}

func parseStats(cacheName string, data io.Reader) (Statistics, error) {
	var stats Statistics
	if (data == nil) {
		log.Warnf("Cannot read stats data for cache '%s' - nil data reader", cacheName)
		return stats, errors.New("handler got nil reader")
	}

	var sohData stats_over_httpData
	err := json.NewDecoder(data).Decode(&sohData)
	if err != nil {
		return stats, err
	}

	if len(sohData.Global) < 1 {
		return stats, errors.New("No 'global' data object found in stats_over_http payload")
	}

	statMap := sohData.Global

	if stats.Loadavg, err = parseLoadAvg(statMap); err != nil {
		return stats, fmt.Errorf("Error parsing loadavg for cache '%s': %v", cacheName, err)
	}

	stats.Interfaces = parseInterfaces(statMap)
	if len(stats.Interfaces) < 1 {
		return stats, fmt.Errorf("cache '%s' had no interfaces", cacheName)
	}

	stats.Miscellaneous = statMap
	return stats, nil
}

func parseLoadAvg(stats map[string]interface{}) (Loadavg, error) {
	var load Loadavg
	if stat, ok := stats["plugin.system_stats.loadavg.one"]; !ok {
		return load, errors.New("Data was missing 'plugin.system_stats.loadavg.one'")
	} else {
		switch t := stat.(type) {
		case float64:
			load.One = stat.(float64)
		case string:
			if statVal, err := strconv.ParseFloat(stat.(string), 64); err != nil {
				return load, fmt.Errorf("loadavg.one could not parse to float, was '%v' (%v)", stat, err)
			} else {
				load.One = statVal
			}
		default:
			return load, fmt.Errorf("loadavg.one had unrecognized type '%T'", t)
		}
	}
	delete(stats, "plugin.system_stats.loadavg.one")

	if stat, ok := stats["plugin.system_stats.loadavg.five"]; !ok {
		return load, errors.New("Data was missing 'plugin.system_stats.loadavg.five'")
	} else {
		switch t := stat.(type) {
		case float64:
			load.Five = stat.(float64)
		case string:
			if statVal, err := strconv.ParseFloat(stat.(string), 64); err != nil {
				return load, fmt.Errorf("loadavg.five could not parse to float, was '%v' (%v)", stat, err)
			} else {
				load.Five = statVal
			}
		default:
			return load, fmt.Errorf("loadavg.five had unrecognized type '%T'", t)
		}
	}
	delete(stats, "plugin.system_stats.loadavg.five")

	if stat, ok := stats["plugin.system_stats.loadavg.fifteen"]; !ok {
		return load, errors.New("Data was missing 'plugin.system_stats.loadavg.fifteen'")
	} else {
		switch t := stat.(type) {
		case float64:
			load.Fifteen = stat.(float64)
		case string:
			if statVal, err := strconv.ParseFloat(stat.(string), 64); err != nil {
				return load, fmt.Errorf("loadavg.fifteen could not parse to float, was '%v' (%v)", stat, err)
			} else {
				load.Fifteen = statVal
			}
		default:
			return load, fmt.Errorf("loadavg.fifteen had unrecognized type '%T'", t)
		}
	}
	delete(stats, "plugin.system_stats.loadavg.fifteen")

	if stat, ok := stats["plugin.system_stats.current_processes"]; !ok {
		return load, errors.New("Data was missing 'plugin.system_stats.current_processes'")
	} else {
		switch t := stat.(type) {
		case float64:
			if stat.(float64) > math.MaxUint64 {
				return load, fmt.Errorf("current number of processes cannot be represented as a uint64 - too big (%v)", stat)
			} else if stat.(float64) < 0 {
				return load, fmt.Errorf("current_processes cannot be negative, got %v", stat)
			}
			load.TotalProcesses = uint64(stat.(float64))
		case string:
			if statVal, err := strconv.ParseUint(stat.(string), 10, 64); err != nil {
				return load, fmt.Errorf("current_processes could not parse to uint64, was '%v' (%v)", stat, err)
			} else {
				load.TotalProcesses = statVal
			}
		default:
			return load, fmt.Errorf("current_processes had unrecognized type '%T'", t)
		}
	}
	delete(stats, "plugin.system_stats.current_processes")
	return load, nil
}

func parseInterfaces(stats map[string]interface{}) (map[string]Interface) {
	ifaces := make(map[string]Interface)
	for stat, value := range stats {
		// The form of the output isn't fully documented, so when something isn't
		// of the right form you don't KNOW something went wrong; it could just be
		// that you're not looking at what you think you're looking at. So when that
		// happens we issue a warning and continue.
		if strings.HasPrefix(stat, "plugin.system_stats.net.") {
			statParts := strings.SplitN(strings.TrimPrefix(stat, "plugin.system_stats.net."), ".", 2)
			if len(statParts) != 2 {
				log.Warnf("stat '%s' appears to be network related, but is not an interface", stat)
				continue
			}
			switch statParts[1] {
			case "rx_bytes":
				var rxBytes uint64
				switch t := value.(type) {
				case float64:
					if value.(float64) > math.MaxUint64 {
						log.Warnf("received bytes for interface '%s' cannot be represented as a uint64 - too big (%v)", statParts[0], value)
						continue
					} else if value.(float64) < 0 {
						log.Warnf("received bytes for interface '%s' was negative (%v)", statParts[0], value)
						continue
					}
					rxBytes = uint64(value.(float64))
				case string:
					if statVal, err := strconv.ParseUint(value.(string), 10, 64); err != nil {
						log.Warnf("received bytes for interface '%s' cannot parse as uint64, was '%v' (%v)", statParts[0], value, err)
						continue
					} else {
						rxBytes = statVal
					}
				default:
					log.Warnf("received bytes for interface '%s' had unrecognized type '%T'", statParts[0], t)
					continue
				}
				tmp := ifaces[statParts[0]]
				tmp.BytesIn = rxBytes
				ifaces[statParts[0]] = tmp
			case "tx_bytes":
				var txBytes uint64
				switch t := value.(type) {
				case float64:
					if value.(float64) > math.MaxUint64 {
						log.Warnf("transmitted bytes for interface '%s' cannot be represented as a uint64 - too big (%v)", statParts[0], value)
						continue
					} else if value.(float64) < 0 {
						log.Warnf("transmitted bytes for interface '%s' was negative (%v)", statParts[0], value)
						continue
					}
					txBytes = uint64(value.(float64))
				case string:
					if statVal, err := strconv.ParseUint(value.(string), 10, 64); err != nil {
						log.Warnf("transmitted bytes for interface '%s' cannot parse as uint64, was '%v' (%v)", statParts[0], value, err)
						continue
					} else {
						txBytes = statVal
					}
				default:
					log.Warnf("transmitted bytes for interface '%s' had unrecognized type '%T'", statParts[0], t)
					continue
				}
				tmp := ifaces[statParts[0]]
				tmp.BytesOut = txBytes
				ifaces[statParts[0]] = tmp
			case "speed":
				var speed int64
				switch t := value.(type) {
				case float64:
					if value.(float64) > math.MaxInt64 || value.(float64) < math.MinInt64 {
						log.Warnf("speed of interface '%s' outside of representable integer range: %v", statParts[0], value)
						continue
					}
					speed = int64(value.(float64))
				case string:
					if statVal, err := strconv.ParseInt(value.(string), 10, 64); err != nil {
						log.Warnf("speed of interface '%s' cannot parse to int64, was '%v': %v", statParts[0], value, err)
						continue
					} else {
						speed = statVal
					}
				default:
					log.Warnf("speed for interface '%s' had unrecognized type '%T'", statParts[0], t)
				}
				tmp := ifaces[statParts[0]]
				tmp.Speed = speed
				ifaces[statParts[0]] = tmp
			}
		}
	}
	return ifaces
}

func statsParse(cacheName string, data io.Reader) (error, map[string]interface{}, AstatsSystem) {
	var outStats AstatsSystem
	if (data == nil) {
		log.Warnf("Cannot read stats data for cache '%s' - nil data reader", cacheName)
		return errors.New("handler got nil reader"), nil, outStats
	}

	statsData := make(map[string]interface{})
	if err := json.NewDecoder(data).Decode(&statsData); err != nil {
		return err, nil, outStats
	}

	if stat, ok := statsData["plugin.system_stats.loadavg.one"]; !ok {
		return errors.New("Data was missing 'plugin.system_stats.loadavg.one'"), nil, outStats
	} else if statStr, ok := stat.(string); !ok {
		return errors.New("'plugin.system_stats.loadavg.one' was not a string"), nil, outStats
	} else {
		outStats.ProcLoadavg = statStr
	}

	if stat, ok := statsData["plugin.system_stats.loadavg.five"]; !ok {
		return errors.New("Data was missing 'plugin.system_stats.loadavg.five'"), nil, outStats
	} else if statStr, ok := stat.(string); !ok {
		return errors.New("'plugin.system_stats.loadavg.five' was not a string"), nil, outStats
	} else {
		outStats.ProcLoadavg += " " + statStr
	}

	if stat, ok := statsData["plugin.system_stats.loadavg.ten"]; !ok {
		return errors.New("Data was missing 'plugin.system_stats.loadavg.ten'"), nil, outStats
	} else if statStr, ok := stat.(string); !ok {
		return errors.New("'plugin.system_stats.loadavg.ten' was not a string"), nil, outStats
	} else {
		outStats.ProcLoadavg += " " + statStr + " 0/0 0" // Dummy end stats, since they aren't used by TM and aren't reported by stats_over_http
	}

	systemStatsData := make(map[string]interface{})
	for key, value := range statsData {
		if strings.HasPrefix(key, "plugin.system_stats.") {
			key = strings.TrimPrefix(key, "plugin.system_stats.")
			if strings.HasPrefix(key, "loadavg.") {
				key = strings.TrimPrefix(key, "loadavg.")
			}
			systemStatsData[strings.TrimPrefix(key, "plugin.system_stats")] = value
		}
	}

	return nil, nil, outStats

}

// func statsPrecompute(cacheName string, data todata.TOData, stats map[string]interface{}, systemStats AstatsSystem) PrecomputedData {

// }
