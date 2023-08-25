package main

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
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type app struct {
	mu            *sync.Mutex
	vstats        vstats
	checkInterval time.Duration
}

type vstats struct {
	ProcLoadavg  string `json:"proc.loadavg"`
	ProcNetDev   string `json:"proc.net.dev"`
	InfSpeed     int64  `json:"inf_speed"`
	NotAvailable bool   `json:"not_available"`
	// TODO: stats
}

func (a *app) getSystemData(ctx context.Context) {
	ticker := time.NewTicker(a.checkInterval)
	for {
		select {
		case <-ticker.C:
			var vstats vstats
			loadavg, err := os.ReadFile("/proc/loadavg")
			if err != nil {
				log.Printf("failed to read /proc/loadavg: %s\n", err)
			}
			vstats.ProcLoadavg = strings.TrimSpace(string(loadavg))

			procNetDev, err := os.ReadFile("/proc/net/dev")
			if err != nil {
				log.Printf("failed to read /proc/net/dev: %s\n", err)
			}
			parts := strings.Split(string(procNetDev), "\n")
			// 3 because first two are columns headers and 2 is loopback interface
			vstats.ProcNetDev = strings.TrimSpace(parts[3])

			infSpeedFile := fmt.Sprintf("/sys/class/net/%s/speed", strings.Split(vstats.ProcNetDev, ":")[0])
			speedStr, err := os.ReadFile(infSpeedFile)
			if err != nil {
				log.Printf("failed to read %s: %s\n", infSpeedFile, err)
			}
			speed, err := strconv.ParseInt(strings.TrimSpace(string(speedStr)), 10, 64)
			if err != nil {
				log.Printf("failed to convert speed '%s' to int: %s\n", speedStr, err)
			}
			vstats.InfSpeed = speed

			cmd := exec.Command("systemctl", "status", "varnish.service")
			err = cmd.Run()
			if err != nil {
				log.Printf("failed to run systemctl: %s\n", err)
			}
			if cmd.ProcessState.ExitCode() != 0 {
				vstats.NotAvailable = true
			}

			a.mu.Lock()
			a.vstats = vstats
			a.mu.Unlock()

		case <-ctx.Done():
			break
		}
	}
}

func (a *app) getStats(w http.ResponseWriter, r *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()
	encoder := json.NewEncoder(w)
	err := encoder.Encode(a.vstats)
	if err != nil {
		log.Printf("failed to write Varnish stats: %s", err)
	}
}

func main() {
	var checkInterval int
	flag.IntVar(&checkInterval, "check-interval", 1, "the duration in seconds to get system data and poll Varnish cache")

	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	app := app{
		mu:            &sync.Mutex{},
		checkInterval: time.Duration(checkInterval) * time.Second,
	}
	go app.getSystemData(ctx)

	http.HandleFunc("/", app.getStats)

	if err := http.ListenAndServe(":2000", nil); err != nil {
		log.Printf("server stopped %s", err)
	}
}
