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

	"github.com/apache/trafficcontrol/lib/go-log"
)

type astatsDataCsv struct {
	Ats map[string]interface{}
}

func astatsCsvParseCsv(cacheName string, data io.Reader, pollCTX interface{}) (Statistics, map[string]interface{}, error) {
	var stats Statistics
	var err error
	if data == nil {
		log.Warnf("Cannot read stats data for cache '%s' - nil data reader", cacheName)
		return stats, nil, errors.New("handler got nil reader")
	}

	var atsData astatsDataCsv
	atsData.Ats = make(map[string]interface{})
	scanner := bufio.NewScanner(data)

	for scanner.Scan() {

		line := scanner.Text()
		delim := strings.IndexByte(line, ',')

		// No delimiter found, skip this line as invalid
		if delim < 0 {
			continue
		}
		// Special cases where we just want the string value
		if strings.Contains(line[0:delim], "proc.") || strings.Contains(line[0:delim], "inf.name") {
			atsData.Ats[line[0:delim]] = line[delim+1 : len(line)]
		} else {
			value, err := strconv.ParseFloat(line[delim+1:len(line)], 64)

			// Skip values that dont parse
			if err != nil {
				continue
			}
			atsData.Ats[line[0:delim]] = value
		}
	}

	if len(atsData.Ats) < 1 {
		return stats, nil, errors.New("No 'global' data object found in stats_over_http payload")
	}

	statMap := atsData.Ats

	// Handle system specific values and remove them from the map for precomputing to not have issues
	if stats.Loadavg, err = LoadavgFromRawLine(statMap["proc.loadavg"].(string)); err != nil {
		return stats, nil, fmt.Errorf("Error parsing loadavg for cache '%s': %v", cacheName, err)
	} else {
		delete(statMap, "proc.loadavg")
	}

	if err := stats.AddInterfaceFromRawLine(statMap["proc.net.dev"].(string)); err != nil {
		return stats, nil, fmt.Errorf("Failed to parse interface line for cache '%s': %v", cacheName, err)
	} else {
		delete(statMap, "proc.net.dev")
	}

	if inf, ok := stats.Interfaces[statMap["inf.name"].(string)]; !ok {
		return stats, nil, errors.New("/proc/net/dev line didn't match reported interface line")
	} else {
		inf.Speed = int64(statMap["inf.speed"].(float64)) //strconv.ParseInt(statMap["inf.speed"].(string), 10, 64)
		stats.Interfaces[statMap["inf.name"].(string)] = inf
		delete(statMap, "inf.speed")
		delete(statMap, "inf.name")

	}

	// Clean up other non-stats entries
	delete(statMap, "astatsLoad")
	delete(statMap, "lastReloadRequest")
	delete(statMap, "version")
	delete(statMap, "something")
	delete(statMap, "lastReload")
	delete(statMap, "configReloadRequests")
	delete(statMap, "configReloads")
	if len(stats.Interfaces) < 1 {
		return stats, nil, fmt.Errorf("cache '%s' had no interfaces", cacheName)
	}

	return stats, statMap, nil
}
