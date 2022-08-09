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
	"fmt"
	"strconv"
	"strings"
)

// DSStat is a single Delivery Service statistic, which is associated with
// a particular cache server.
// DSStats are referenced by name in a map associated with cache server data,
// and so that name (XMLID) is not reiterated here.
type DSStat struct {
	// InBytes is the total number of bytes received by the cache server which
	// were for the Delivery Service.
	InBytes uint64
	// OutBytes is the total number of bytes transmitted by the cache server
	// in service of the Delivery Service.
	OutBytes uint64
	// Status2xx is the number of requests for the Delivery Service's content
	// that were served with responses having response codes on the interval
	// [200, 300).
	Status2xx uint64
	// Status3xx is the number of requests for the Delivery Service's content
	// that were served with responses having response codes on the interval
	// [300, 400).
	Status3xx uint64
	// Status4xx is the number of requests for the Delivery Service's content
	// that were served with responses having response codes on the interval
	// [400, 500).
	Status4xx uint64
	// Status2xx is the number of requests for the Delivery Service's content
	// that were served with responses having response codes on the interval
	// [500, 600).
	Status5xx uint64
}

// Loadavg contains the parsed "loadavg" data for a polled cache server.
// Specifically, it contains all of the data stored that can be found in
// /proc/loadavg on a Linux/Unix system.
//
// For more information on what a "loadavg" is, consult the “proc(5)” man page
// (web-hosted: https://linux.die.net/man/5/proc).
type Loadavg struct {
	// One is the cache server's "loadavg" in the past minute from the time it was
	// polled.
	One float64
	// Five is the cache server's "loadavg" in the past five minutes from the time
	// it was polled.
	Five float64
	// Fifteen is the cache server's "loadavg" in the past fifteen minutes from the
	// time it was polled.
	Fifteen float64
	// CurrentProcesses is the number of currently executing processes (or threads)
	// on the cache server.
	// Note that stats_over_http doesn't provide this, so in general it can't be
	// relied on to be set properly.
	CurrentProcesses uint64
	// TotalProcesses is the number of total processes (or threads) that exist on
	// the cache server.
	TotalProcesses uint64
	// LatestPID is the process ID of the most recently created process on the
	// cache server at the time of polling.
	// Note that stats_over_http doesn't provide this, so in general it can't be
	// relied on to be set properly - which is fine because what use could that
	// information actually have??
	LatestPID int64
}

func splitLoadavgOn(char rune) bool {
	return char == ' ' || char == '/' || char == '\t'
}

// LoadavgFromRawLine parses a raw line - presumably read from /proc/loadavg -
// and returns a Loadavg containing all of the same information, as well as
// any encountered error.
func LoadavgFromRawLine(line string) (Loadavg, error) {
	var load Loadavg
	fields := strings.FieldsFunc(line, splitLoadavgOn)
	if len(fields) != 6 {
		return load, fmt.Errorf("Expected 6 fields in a loadavg line, got %d", len(fields))
	}

	if loadStat, err := strconv.ParseFloat(fields[0], 64); err != nil {
		return load, fmt.Errorf("Error parsing one-minute loadavg: %v", err)
	} else {
		load.One = loadStat
	}
	if loadStat, err := strconv.ParseFloat(fields[1], 64); err != nil {
		return load, fmt.Errorf("Error parsing five-minute loadavg: %v", err)
	} else {
		load.Five = loadStat
	}
	if loadStat, err := strconv.ParseFloat(fields[2], 64); err != nil {
		return load, fmt.Errorf("Error parsing fifteen-minute loadavg: %v", err)
	} else {
		load.Fifteen = loadStat
	}

	if loadStat, err := strconv.ParseUint(fields[3], 10, 64); err != nil {
		return load, fmt.Errorf("Error parsing currently executing processes: %v", err)
	} else {
		load.CurrentProcesses = loadStat
	}

	if loadStat, err := strconv.ParseUint(fields[4], 10, 64); err != nil {
		return load, fmt.Errorf("Error parsing total processes: %v", err)
	} else {
		load.TotalProcesses = loadStat
	}

	if loadStat, err := strconv.ParseInt(fields[5], 10, 64); err != nil {
		return load, fmt.Errorf("Error parsing latest process ID: %v", err)
	} else {
		load.LatestPID = loadStat
	}

	return load, nil
}

// Interface represents a network interface. The name of the interface is
// used to access it within a Statistics object, and so is not stored here.
type Interface struct {
	// Speed is the "speed" of the interface, which is of unknown - but vitally
	// important - meaning.
	Speed int64
	// BytesOut is the total number of bytes transmitted by this interface.
	BytesOut uint64
	// BytesIn is the total number of bytes received by this interface.
	BytesIn uint64
}

// Statistics is a structure containing, most generally, the statistics of a
// cache server.
type Statistics struct {
	// Loadavg contains the Unix/Linux "loadavg" values for the cache server.
	Loadavg Loadavg
	// Interfaces is a map of network interface names to statistic data about
	// those interfaces.
	Interfaces map[string]Interface
	// NotAvailable reports whether or not the cache server is unavailable.
	// Sometimes caches can directly report this, but it's not supported by
	// stats_over_http (afaik), so it always just uses ``false''
	NotAvailable bool
}

// AddInterfaceFromRawLine parses the raw line - presumably read from
// /proc/net/dev - and inserts into the Statistics a new Interface
// containing the data provided.
//
// This will initialize s.Interfaces if that has not already been done (if the
// parse is successful).
//
// If the line cannot be parsed, s.Interfaces is unchanged and an error
// describing the problem is returned.
//
// Note that this does *not* set the interface's Speed.
func (s *Statistics) AddInterfaceFromRawLine(line string) error {
	var iface Interface
	parts := strings.Split(line, ":")
	if len(parts) != 2 {
		return fmt.Errorf("Expected /proc/net/dev line to be in the format '{{name}}:{{info}}', but got '%s'", line)
	}

	name := parts[0]
	parts = strings.Fields(parts[1])
	if len(parts) < 9 {
		return fmt.Errorf("Expected at least 9 /proc/net/dev fields, got %d", len(parts))
	}

	var err error
	if iface.BytesIn, err = strconv.ParseUint(parts[0], 10, 64); err != nil {
		return fmt.Errorf("Error parsing BytesIn: %v", err)
	}
	if iface.BytesOut, err = strconv.ParseUint(parts[8], 10, 64); err != nil {
		return fmt.Errorf("Error parsing BytesOut: %v", err)
	}

	if s.Interfaces == nil {
		s.Interfaces = map[string]Interface{
			name: iface,
		}
	} else {
		s.Interfaces[name] = iface
	}
	return nil
}
