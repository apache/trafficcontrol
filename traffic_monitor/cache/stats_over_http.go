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

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/poller"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/todata"

	jsoniter "github.com/json-iterator/go"
)

// LOADAVG_SHIFT is the amount by which "loadavg" values returned by
// stats_over_http need to be divided to obtain the values with which ATC
// operators are more familiar.
//
// The reason for this is that the Linux kernel stores loadavg values as
// integral types internally, and performs conversions to floating-point
// numbers on-the-fly when the contents of /proc/loadavg are read. Since
// astats_over_http used to always just read that file to get the numbers,
// everyone's used to the floating-point form. But stats_over_http gets
// the numbers directly from a syscall, so they aren't pre-converted for us.
//
// Dividing by this number is kind of a shortcut, for the actual transformation
// used by the kernel itself, refer to the source:
// https://github.com/torvalds/linux/blob/master/fs/proc/loadavg.c
const LOADAVG_SHIFT = 65536

func init() {
	registerDecoder("stats_over_http", statsOverHTTPParse, statsOverHTTPPrecompute)
}

type stats_over_httpData struct {
	Global map[string]interface{} `json:"global"`
}

func statsOverHTTPParse(cacheName string, data io.Reader, pollCTX interface{}) (Statistics, map[string]interface{}, error) {
	var stats Statistics
	if data == nil {
		log.Warnf("Cannot read stats data for cache '%s' - nil data reader", cacheName)
		return stats, nil, errors.New("handler got nil reader")
	}

	var sohData stats_over_httpData
	var err error

	ctx := pollCTX.(*poller.HTTPPollCtx)

	ctype := ctx.HTTPHeader.Get("Content-Type")

	if ctype == "text/json" || ctype == "text/javascript" || ctype == "application/json" || ctype == "" {
		json := jsoniter.ConfigFastest
		err := json.NewDecoder(data).Decode(&sohData)
		if err != nil {
			return stats, nil, err
		}
	} else if ctype == "text/csv" {
		sohData.Global, err = statsOverHTTPParseCSV(cacheName, data)
		if err != nil {
			return stats, nil, err
		}
	} else {
		return stats, nil, fmt.Errorf("stats Content-Type (%s) can not be parsed by statsOverHTTP", ctype)
	}

	if len(sohData.Global) < 1 {
		return stats, nil, errors.New("No 'global' data object found in stats_over_http payload")
	}

	statMap := sohData.Global

	if stats.Loadavg, err = parseLoadAvg(statMap); err != nil {
		return stats, nil, fmt.Errorf("Error parsing loadavg for cache '%s': %v", cacheName, err)
	}

	stats.Interfaces = parseInterfaces(statMap)
	if len(stats.Interfaces) < 1 {
		return stats, nil, fmt.Errorf("cache '%s' had no interfaces", cacheName)
	}

	return stats, statMap, nil
}

func statsOverHTTPParseCSV(cacheName string, data io.Reader) (map[string]interface{}, error) {

	if data == nil {
		log.Warnf("Cannot read stats data for cache '%s' - nil data reader", cacheName)
		return nil, errors.New("handler got nil reader")
	}

	var allData []string
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		allData = append(allData, scanner.Text())
	}

	globalData := make(map[string]interface{}, len(allData))

	for _, line := range allData {
		delim := strings.IndexByte(line, ',')

		// No delimiter found, skip this line as invalid
		if delim < 0 {
			continue
		}

		value, err := strconv.ParseFloat(line[delim+1:], 64)

		// Skip values that dont parse
		if err != nil {
			continue
		}
		globalData[line[0:delim]] = value
	}

	if len(globalData) < 1 {
		return nil, errors.New("no valid data found in stats_over_http payload with csv format")
	}

	return globalData, nil
}

func parseLoadAvg(stats map[string]interface{}) (Loadavg, error) {
	var load Loadavg
	if stat, ok := stats["plugin.system_stats.loadavg.one"]; !ok {
		return load, errors.New("Data was missing 'plugin.system_stats.loadavg.one'")
	} else {
		switch t := stat.(type) {
		case float64:
			load.One = stat.(float64) / LOADAVG_SHIFT
		case string:
			if statVal, err := strconv.ParseFloat(stat.(string), 64); err != nil {
				return load, fmt.Errorf("loadavg.one could not parse to float, was '%v' (%v)", stat, err)
			} else {
				load.One = statVal / LOADAVG_SHIFT
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
			load.Five = stat.(float64) / LOADAVG_SHIFT
		case string:
			if statVal, err := strconv.ParseFloat(stat.(string), 64); err != nil {
				return load, fmt.Errorf("loadavg.five could not parse to float, was '%v' (%v)", stat, err)
			} else {
				load.Five = statVal / LOADAVG_SHIFT
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
			load.Fifteen = stat.(float64) / LOADAVG_SHIFT
		case string:
			if statVal, err := strconv.ParseFloat(stat.(string), 64); err != nil {
				return load, fmt.Errorf("loadavg.fifteen could not parse to float, was '%v' (%v)", stat, err)
			} else {
				load.Fifteen = statVal / LOADAVG_SHIFT
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

func parseInterfaces(stats map[string]interface{}) map[string]Interface {
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

func parseNumericStat(value interface{}) (uint64, error) {
	switch t := value.(type) {
	case uint:
		return uint64(value.(uint)), nil
	case uint32:
		return uint64(value.(uint32)), nil
	case uint64:
		return value.(uint64), nil
	case int:
		if value.(int) < 0 {
			return 0, errors.New("value was negative")
		}
		return uint64(value.(int)), nil
	case int32:
		if value.(int32) < 0 {
			return 0, errors.New("value was negative")
		}
		return uint64(value.(int32)), nil
	case int64:
		if value.(int64) < 0 {
			return 0, errors.New("value was negative")
		}
		return value.(uint64), nil
	case float64:
		if value.(float64) > math.MaxUint64 || value.(float64) < 0 {
			return 0, errors.New("value out of range for uint64")
		}
		return uint64(value.(float64)), nil
	case float32:
		if value.(float32) > math.MaxUint64 || value.(float32) < 0 {
			return 0, errors.New("value out of range for uint64")
		}
		return uint64(value.(float64)), nil
	case string:
		if statVal, err := strconv.ParseUint(value.(string), 10, 64); err != nil {
			return 0, fmt.Errorf("could not parse '%v' to uint64: %v", value, err)
		} else {
			return statVal, nil
		}
	default:
		return 0, fmt.Errorf("value '%v' is of unrecognized type %T", value, t)
	}
}

func statsOverHTTPPrecompute(cacheName string, data todata.TOData, stats Statistics, miscStats map[string]interface{}) PrecomputedData {
	var precomputed PrecomputedData
	precomputed.DeliveryServiceStats = make(map[string]*DSStat)

	precomputed.OutBytes = 0
	precomputed.MaxKbps = 0
	for _, iface := range stats.Interfaces {
		precomputed.OutBytes += iface.BytesOut
		if iface.Speed > precomputed.MaxKbps {
			precomputed.MaxKbps = iface.Speed
		}
	}
	precomputed.MaxKbps *= 1000

	for stat, value := range miscStats {
		if strings.HasPrefix(stat, "plugin.remap_stats.") {
			trimmedStat := strings.TrimPrefix(stat, "plugin.remap_stats.")
			statParts := strings.Split(trimmedStat, ".")
			if len(statParts) < 3 {
				err := errors.New("stat has no remap_stats deliveryservice and name parts")
				log.Infof("precomputing cache %s stat %s value %v error %v", cacheName, stat, value, err)
				precomputed.Errors = append(precomputed.Errors, err)
				continue
			}
			subsubdomain := statParts[0]
			subdomain := statParts[1]
			domain := strings.Join(statParts[2:len(statParts)-1], ".")
			ds, ok := data.DeliveryServiceRegexes.DeliveryService(domain, subdomain, subsubdomain)
			if !ok {
				err := errors.New("No Delivery Service match for stat")
				log.Infof("precomputing cache %s stat %s value %v error %v", cacheName, stat, value, err)
				precomputed.Errors = append(precomputed.Errors, err)
				continue
			}
			if ds == "" {
				err := errors.New("Empty Delivery Service FQDN")
				log.Infof("precomputing cache %s stat %s value %v error %v", cacheName, stat, value, err)
				precomputed.Errors = append(precomputed.Errors, err)
				continue
			}
			dsName := string(ds)
			dsStat, ok := precomputed.DeliveryServiceStats[dsName]
			if !ok || dsStat == nil {
				dsStat = new(DSStat)
			}

			parsedStat, err := parseNumericStat(value)
			if err != nil {
				err = fmt.Errorf("couldn't parse numeric stat: %v", err)
				log.Infof("precomputing cache %s stat %s value %v error %v", cacheName, stat, value, err)
				precomputed.Errors = append(precomputed.Errors, err)
				continue
			}

			switch statParts[len(statParts)-1] {
			case "status_2xx":
				dsStat.Status2xx += parsedStat
			case "status_3xx":
				dsStat.Status3xx += parsedStat
			case "status_4xx":
				dsStat.Status4xx += parsedStat
			case "status_5xx":
				dsStat.Status5xx += parsedStat
			case "out_bytes":
				dsStat.OutBytes += parsedStat
			case "in_bytes":
				dsStat.InBytes += parsedStat
			default:
				err = fmt.Errorf("Unknown stat '%s'", statParts[len(statParts)-1])
				log.Infof("precomputing cache %s stat %s value %v error %v", cacheName, stat, value, err)
				precomputed.Errors = append(precomputed.Errors, err)
				continue
			}
			precomputed.DeliveryServiceStats[dsName] = dsStat
		}
	}
	return precomputed
}
