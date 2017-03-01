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
	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/traffic_monitor/tmcheck"
	to "github.com/apache/incubator-trafficcontrol/traffic_ops/client"
	"net/http"
	"sync"
	"time"
)

const UserAgent = "tm-offline-validator/0.1"

type Log struct {
	l *[]string
	m *sync.RWMutex
}

func (l *Log) Add(msg string) {
	l.m.Lock()
	defer l.m.Unlock()
	*l.l = append(*l.l, msg)
}

func (l *Log) Get() []string {
	l.m.RLock()
	defer l.m.RUnlock()
	return *l.l
}

func NewLog() Log {
	s := make([]string, 0)
	return Log{l: &s, m: &sync.RWMutex{}}
}

func main() {
	toURI := flag.String("to", "", "The Traffic Ops URI, whose CRConfig to validate")
	toUser := flag.String("touser", "", "The Traffic Ops user")
	toPass := flag.String("topass", "", "The Traffic Ops password")
	tmURI := flag.String("tm", "", "The Traffic Monitor URI whose CRStates to validate")
	interval := flag.Duration("interval", time.Second*time.Duration(5), "The interval to validate")
	grace := flag.Duration("grace", time.Second*time.Duration(30), "The grace period before invalid states are reported")
	help := flag.Bool("help", false, "Usage info")
	helpBrief := flag.Bool("h", false, "Usage info")
	flag.Parse()
	if *help || *helpBrief {
		fmt.Printf("Usage: go run validate-offline -to https://traffic-ops.example.net -touser bill -topass thelizard -tm http://traffic-monitor.example.net -interval 5s -grace 30s\n")
		return
	}

	toClient, err := to.LoginWithAgent(*toURI, *toUser, *toPass, true, UserAgent, false, tmcheck.RequestTimeout)
	if err != nil {
		fmt.Printf("Error logging in to Traffic Ops: %v\n", err)
		return
	}

	log := NewLog()

	onErr := func(err error) {
		log.Add(fmt.Sprintf("%v ERROR %v\n", time.Now(), err))
	}

	onResumeSuccess := func() {
		log.Add(fmt.Sprintf("%v INFO State Valid\n", time.Now()))
	}

	onCheck := func(err error) {
		if err != nil {
			log.Add(fmt.Sprintf("%v DEBUG invalid: %v\n", time.Now(), err))
		} else {
			log.Add(fmt.Sprintf("%v DEBUG valid\n", time.Now()))
		}
	}

	go tmcheck.Validator(*tmURI, toClient, *interval, *grace, onErr, onResumeSuccess, onCheck)

	if err := serve(log); err != nil {
		fmt.Printf("Serve error: %v\n", err)
	}
}

func serve(log Log) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<html><head><meta http-equiv="refresh" content="5"></head><body><pre>`)
		logCopy := log.Get()
		for i := len(logCopy) - 1; i >= 0; i-- {
			fmt.Fprintf(w, "%s\n", logCopy[i])
		}
		fmt.Fprintf(w, `</pre></body></html>`)
	})
	return http.ListenAndServe(":80", nil)
}
