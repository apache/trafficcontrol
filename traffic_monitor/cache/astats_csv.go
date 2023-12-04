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
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
)

type astatsDataCsv struct {
	Ats map[string]interface{}
}

func astatsCsvParseCsv(cacheName string, data io.Reader) (Statistics, map[string]interface{}, error) {
	var stats Statistics
	var err error
	if data == nil {
		log.Warnf("Cannot read stats data for cache '%s' - nil data reader", cacheName)
		return stats, nil, errors.New("handler got nil reader")
	}

	var atsData astatsDataCsv
	var allData []string
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		allData = append(allData, scanner.Text())
	}

	atsData.Ats = make(map[string]interface{}, len(allData))

	for _, line := range allData {
		delim := strings.IndexByte(line, ',')

		// No delimiter found, skip this line as invalid
		if delim < 0 {
			continue
		}
		// Special cases where we just want the string value
		if strings.Contains(line[0:delim], "proc.") || strings.Contains(line[0:delim], "inf.name") {
			atsData.Ats[line[0:delim]] = line[delim+1:]
		} else {
			value, err := strconv.ParseFloat(line[delim+1:], 64)

			// Skip values that dont parse
			if err != nil {
				continue
			}
			atsData.Ats[line[0:delim]] = value
		}
	}

	if len(atsData.Ats) < 1 {
		return stats, nil, errors.New("no 'global' data object found in stats_over_http payload")
	}

	statMap := atsData.Ats

	// Handle system specific values and remove them from the map for precomputing to not have issues
	if loadAvg, hadLoadAvg := statMap["proc.loadavg"].(string); !hadLoadAvg {
		return stats, nil, fmt.Errorf("failed to parse loadavg line for cache '%s': %v", cacheName, "no proc.loadavg in Astats response")
	} else if stats.Loadavg, err = LoadavgFromRawLine(loadAvg); err != nil {
		return stats, nil, fmt.Errorf("parsing loadavg for cache '%s': %v", cacheName, err)
	} else {
		delete(statMap, "proc.loadavg")
	}

	if procNetDev, hadProcNetDev := statMap["proc.net.dev"].(string); !hadProcNetDev {
		return stats, nil, fmt.Errorf("failed to parse interface line for cache '%s': %v", cacheName, "no proc.net.dev in Astats response")
	} else if err := stats.AddInterfaceFromRawLine(procNetDev); err != nil {
		return stats, nil, fmt.Errorf("failed to parse interface line for cache '%s': %v", cacheName, err)
	} else {
		delete(statMap, "proc.net.dev")
	}

	if infName, hadinf := statMap["inf.name"].(string); !hadinf {
		return stats, nil, fmt.Errorf("failed to parse inf line for cache '%s': %v", cacheName, "no inf.name in Astats response")
	} else if inf, err := stats.Interfaces[infName]; !err {
		return stats, nil, errors.New("/proc/net/dev line didn't match reported interface line")
	} else {
		if infSpeed, hadspeed := statMap["inf.speed"].(float64); !hadspeed {
			return stats, nil, fmt.Errorf("failed to parse interface speed line for cache '%s': %v", cacheName, "no inf.speed in Astats response")
		} else {
			inf.Speed = int64(infSpeed)
			stats.Interfaces[infName] = inf
			delete(statMap, "inf.speed")
			delete(statMap, "inf.name")
		}
	}

	// Clean up other non-stats entries
	nonStats := []string{
		"astatsLoad",
		"lastReloadRequest",
		"version",
		"something",
		"lastReload",
		"configReloadRequests",
		"configReloads",
	}
	for _, nonStat := range nonStats {
		delete(statMap, nonStat)
	}

	if len(stats.Interfaces) < 1 {
		return stats, nil, fmt.Errorf("cache '%s' had no interfaces", cacheName)
	}

	return stats, statMap, nil
}
