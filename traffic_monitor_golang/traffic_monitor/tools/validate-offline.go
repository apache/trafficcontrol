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
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/enum"
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/tmcheck"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	"net/http"
	"sync"
	"time"
)

const UserAgent = "tm-offline-validator/0.1"

const LogLimit = 10000

type Log struct {
	log     *[]string
	limit   int
	errored *bool
	m       *sync.RWMutex
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

func (l *Log) GetErrored() bool {
	l.m.RLock()
	defer l.m.RUnlock()
	return *l.errored
}

func (l *Log) SetErrored(e bool) {
	l.m.Lock()
	defer l.m.Unlock()
	*l.errored = e
}

func NewLog() Log {
	log := make([]string, 0, LogLimit+1)
	errored := false
	limit := LogLimit
	return Log{log: &log, errored: &errored, m: &sync.RWMutex{}, limit: limit}
}

type Logs struct {
	logs map[enum.TrafficMonitorName]Log
	m    *sync.RWMutex
}

func NewLogs() Logs {
	return Logs{logs: map[enum.TrafficMonitorName]Log{}, m: &sync.RWMutex{}}
}

func (l Logs) Get(name enum.TrafficMonitorName) Log {
	l.m.Lock()
	defer l.m.Unlock()
	if _, ok := l.logs[name]; !ok {
		l.logs[name] = NewLog()
	}
	return l.logs[name]
}

func (l Logs) GetMonitors() []enum.TrafficMonitorName {
	l.m.RLock()
	defer l.m.RUnlock()
	monitors := []enum.TrafficMonitorName{}
	for name, _ := range l.logs {
		monitors = append(monitors, name)
	}
	return monitors
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

	toClient, err := to.LoginWithAgent(*toURI, *toUser, *toPass, true, UserAgent, false, tmcheck.RequestTimeout)
	if err != nil {
		fmt.Printf("Error logging in to Traffic Ops: %v\n", err)
		return
	}

	logs := NewLogs()

	onErr := func(name enum.TrafficMonitorName, err error) {
		log := logs.Get(name)
		log.Add(fmt.Sprintf("%v ERROR %v\n", time.Now(), err))
		log.SetErrored(true)
	}

	onResumeSuccess := func(name enum.TrafficMonitorName) {
		log := logs.Get(name)
		log.Add(fmt.Sprintf("%v INFO State Valid\n", time.Now()))
		log.SetErrored(false)
	}

	onCheck := func(name enum.TrafficMonitorName, err error) {
		log := logs.Get(name)
		if err != nil {
			log.Add(fmt.Sprintf("%v DEBUG invalid: %v\n", time.Now(), err))
		} else {
			log.Add(fmt.Sprintf("%v DEBUG valid\n", time.Now()))
		}
	}

	go tmcheck.AllMonitorsCRStatesOfflineValidator(toClient, *interval, *includeOffline, *grace, onErr, onResumeSuccess, onCheck)

	if err := serve(logs, *toURI); err != nil {
		fmt.Printf("Serve error: %v\n", err)
	}
}

func serve(logs Logs, toURI string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<!DOCTYPE html>
<meta http-equiv="refresh" content="5">
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Traffic Monitor Offline Validator</title>
<style type="text/css">body{margin:40px auto;line-height:1.6;font-size:18px;color:#444;padding:0 2px}h1,h2,h3{line-height:1.2}</style>`)

		fmt.Fprintf(w, `<p>%s`, toURI)

		fmt.Fprintf(w, `<table style="width:100%%"><tr>`)

		monitors := logs.GetMonitors()
		for _, monitor := range monitors {
			fmt.Fprintf(w, `<td>`)

			log := logs.Get(monitor)

			fmt.Fprintf(w, `<p>%s`, monitor)
			if log.GetErrored() {
				fmt.Fprintf(w, `<h1 style="color:red">Invalid</h1>`)
			} else {
				fmt.Fprintf(w, `<h1 style="color:limegreen">Valid</h1>`)
			}

			fmt.Fprintf(w, `<pre>`)
			logCopy := log.Get()
			for _, msg := range logCopy {
				fmt.Fprintf(w, "%s\n", msg)
			}
			fmt.Fprintf(w, `</pre>`)

			fmt.Fprintf(w, `</td>`)
		}

		fmt.Fprintf(w, `</tr></table>`)
	})
	return http.ListenAndServe(":80", nil)
}
