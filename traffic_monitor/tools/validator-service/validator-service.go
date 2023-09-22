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

// validate-offline is a utility HTTP service which polls the given Traffic Monitor and validates that no OFFLINE or ADMIN_DOWN caches in the Traffic Ops CRConfig are marked Available in Traffic Monitor's CRstates endpoint.

package main

import (
	"flag"
	"fmt"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_monitor/tmcheck"
	to "github.com/apache/trafficcontrol/v8/traffic_ops/v4-client"
	"io"
	"net/http"
	"sort"
	"sync"
	"time"
)

const UserAgent = "tm-offline-validator/0.1"

const LogLimit = 10

type Log struct {
	log       *[]string
	limit     int
	errored   *bool
	lastCheck *time.Time
	m         *sync.RWMutex
}

func (l *Log) Add(msg string) {
	l.m.Lock()
	defer l.m.Unlock()
	*l.log = append([]string{msg}, *l.log...)
	if len(*l.log) > l.limit {
		*l.log = (*l.log)[:l.limit]
	}
}

func (l *Log) Get() []string {
	l.m.RLock()
	defer l.m.RUnlock()
	return *l.log
}

func (l *Log) GetErrored() (bool, time.Time) {
	l.m.RLock()
	defer l.m.RUnlock()
	return *l.errored, *l.lastCheck
}

func (l *Log) SetErrored(e bool) {
	l.m.Lock()
	defer l.m.Unlock()
	*l.errored = e
	*l.lastCheck = time.Now()
}

func NewLog() Log {
	log := make([]string, 0, LogLimit+1)
	errored := false
	limit := LogLimit
	lastCheck := time.Time{}
	return Log{log: &log, errored: &errored, m: &sync.RWMutex{}, limit: limit, lastCheck: &lastCheck}
}

type Logs struct {
	logs map[tc.TrafficMonitorName]Log
	m    *sync.RWMutex
}

func NewLogs() Logs {
	return Logs{logs: map[tc.TrafficMonitorName]Log{}, m: &sync.RWMutex{}}
}

func (l Logs) Get(name tc.TrafficMonitorName) Log {
	l.m.Lock()
	defer l.m.Unlock()
	if _, ok := l.logs[name]; !ok {
		l.logs[name] = NewLog()
	}
	return l.logs[name]
}

func (l Logs) GetMonitors() []string {
	l.m.RLock()
	defer l.m.RUnlock()
	monitors := []string{}
	for name, _ := range l.logs {
		monitors = append(monitors, string(name))
	}
	return monitors
}

func startValidator(validator tmcheck.AllValidatorFunc, toClient *to.Session, interval time.Duration, includeOffline bool, grace time.Duration) Logs {
	logs := NewLogs()

	onErr := func(name tc.TrafficMonitorName, err error) {
		log := logs.Get(name)
		log.Add(fmt.Sprintf("%v ERROR %v\n", time.Now(), err))
		log.SetErrored(true)
	}

	onResumeSuccess := func(name tc.TrafficMonitorName) {
		log := logs.Get(name)
		log.Add(fmt.Sprintf("%v INFO State Valid\n", time.Now()))
		log.SetErrored(false)
	}

	onCheck := func(name tc.TrafficMonitorName, err error) {
		log := logs.Get(name)
		log.SetErrored(err != nil)
	}

	go validator(toClient, interval, includeOffline, grace, onErr, onResumeSuccess, onCheck)
	return logs
}

func main() {
	toURI := flag.String("to", "", "The Traffic Ops URI, whose CRConfig to validate")
	toUser := flag.String("touser", "", "The Traffic Ops user")
	toPass := flag.String("topass", "", "The Traffic Ops password")
	interval := flag.Duration("interval", time.Second*time.Duration(5), "The interval to validate")
	grace := flag.Duration("grace", time.Second*time.Duration(30), "The grace period before invalid states are reported")
	includeOffline := flag.Bool("includeOffline", false, "Whether to include Offline Monitors")
	help := flag.Bool("help", false, "Usage info")
	helpBrief := flag.Bool("h", false, "Usage info")
	flag.Parse()
	if *help || *helpBrief {
		fmt.Printf("Usage: go run validate-offline -to https://traffic-ops.example.net -touser bill -topass thelizard -tm http://traffic-monitor.example.net -interval 5s -grace 30s -includeOffline true\n")
		return
	}

	toClient, _, err := to.LoginWithAgent(*toURI, *toUser, *toPass, true, UserAgent, false, tmcheck.RequestTimeout)
	if err != nil {
		fmt.Printf("Error logging in to Traffic Ops: %v\n", err)
		return
	}

	crStatesOfflineLogs := startValidator(tmcheck.AllMonitorsCRStatesOfflineValidator, toClient, *interval, *includeOffline, *grace)
	peerPollerLogs := startValidator(tmcheck.PeerPollersAllValidator, toClient, *interval, *includeOffline, *grace)
	dsStatsLogs := startValidator(tmcheck.AllMonitorsDSStatsValidator, toClient, *interval, *includeOffline, *grace)
	queryIntervalLogs := startValidator(tmcheck.AllMonitorsQueryIntervalValidator, toClient, *interval, *includeOffline, *grace)

	if err := serve(*toURI, crStatesOfflineLogs, peerPollerLogs, dsStatsLogs, queryIntervalLogs); err != nil {
		fmt.Printf("Serve error: %v\n", err)
	}
}

func printLogs(logs Logs, w io.Writer) {
	fmt.Fprintf(w, `<table style="width:100%%">`)

	monitors := logs.GetMonitors()
	sort.Strings(monitors) // sort, so they're always in the same order in the webpage
	for _, monitor := range monitors {
		fmt.Fprintf(w, `</tr>`)

		log := logs.Get(tc.TrafficMonitorName(monitor))

		fmt.Fprintf(w, `<td><span>%s</span></td>`, monitor)
		errored, lastCheck := log.GetErrored()
		if errored {
			fmt.Fprintf(w, `<td><span style="color:red">Invalid</span></td>`)
		} else {
			fmt.Fprintf(w, `<td><span style="color:limegreen">Valid</span></td>`)
		}
		fmt.Fprintf(w, `<td><span>as of %v</span></td>`, lastCheck)

		if errored {
			fmt.Fprintf(w, `<td><span style="font-family:monospace">`)
			logCopy := log.Get()
			firstMsg := ""
			if len(logCopy) > 0 {
				firstMsg = logCopy[0]
			}
			fmt.Fprintf(w, "%s\n", firstMsg)
			fmt.Fprintf(w, `</span></td>`)
		}

		fmt.Fprintf(w, `</tr>`)
	}
	fmt.Fprintf(w, `</table>`)
}

func serve(toURI string, crStatesOfflineLogs Logs, peerPollerLogs Logs, dsStatsLogs Logs, queryIntervalLogs Logs) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<!DOCTYPE html>
<meta http-equiv="refresh" content="5">
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Traffic Monitor Offline Validator</title>
<style type="text/css">body{margin:40px auto;line-height:1.6;font-size:18px;color:#444;padding:0 8px 0 8px}h1,h2,h3{line-height:1.2}span{padding:0px 4px 0px 4px;}</style>`)

		fmt.Fprintf(w, `<h1>Traffic Monitor Validator</h1>`)

		fmt.Fprintf(w, `<p>%s`, toURI)
		fmt.Fprintf(w, `<p>%s`, time.Now())

		fmt.Fprintf(w, `<h2>CRStates Offline</h2>`)
		fmt.Fprintf(w, `<h3>validates all OFFLINE and ADMIN_DOWN caches in the CRConfig are Unavailable</h3>`)
		printLogs(crStatesOfflineLogs, w)

		fmt.Fprintf(w, `<h2>Peer Poller</h2>`)
		fmt.Fprintf(w, `<h3>validates all peers in the CRConfig have been polled within the last %v</h3>`, tmcheck.PeerPollMax)
		printLogs(peerPollerLogs, w)

		fmt.Fprintf(w, `<h2>Delivery Services</h2>`)
		fmt.Fprintf(w, `<h3>validates all Delivery Services in the CRConfig exist in DsStats</h3>`)
		printLogs(dsStatsLogs, w)

		fmt.Fprintf(w, `<h2>Query Interval</h2>`)
		fmt.Fprintf(w, `<h3>validates all Monitors' Query Interval (95th percentile) is less than %v</h3>`, tmcheck.QueryIntervalMax)
		printLogs(queryIntervalLogs, w)
	})
	return http.ListenAndServe(":80", nil)
}
