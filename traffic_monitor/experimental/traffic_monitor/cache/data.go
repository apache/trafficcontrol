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
	"github.com/apache/incubator-trafficcontrol/traffic_monitor/experimental/traffic_monitor/enum"
)

// CacheAvailableStatusReported is the status string returned by caches set to "reported" in Traffic Ops.
// TODO put somewhere more generic
const AvailableStatusReported = "REPORTED"

// CacheAvailableStatus is the available status of the given cache. It includes a boolean available/unavailable flag, and a descriptive string.
type AvailableStatus struct {
	Available bool
	Status    string
	Why       string
}

// CacheAvailableStatuses is the available status of each cache.
type AvailableStatuses map[enum.CacheName]AvailableStatus

// Copy copies this CacheAvailableStatuses. It does not modify, and thus is safe for multiple reader goroutines.
func (a AvailableStatuses) Copy() AvailableStatuses {
	b := AvailableStatuses(map[enum.CacheName]AvailableStatus{})
	for k, v := range a {
		b[k] = v
	}
	return b
}

// Event represents an event change in aggregated data. For example, a cache being marked as unavailable.
type Event struct {
	Index       uint64         `json:"index"`
	Time        int64          `json:"time"`
	Description string         `json:"description"`
	Name        enum.CacheName `json:"name"`
	Hostname    enum.CacheName `json:"hostname"`
	Type        string         `json:"type"`
	Available   bool           `json:"isAvailable"`
}

// ResultHistory is a map of cache names, to an array of result history from each cache.
type ResultHistory map[enum.CacheName][]Result

func copyResult(a []Result) []Result {
	b := make([]Result, len(a), len(a))
	copy(b, a)
	return b
}

// Copy copies returns a deep copy of this ResultHistory
func (a ResultHistory) Copy() ResultHistory {
	b := ResultHistory{}
	for k, v := range a {
		b[k] = copyResult(v)
	}
	return b
}
