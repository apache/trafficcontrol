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
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type vstats struct {
	ProcLoadavg  string `json:"proc.loadavg"`
	ProcNetDev   string `json:"proc.net.dev"`
	InfSpeed     int64  `json:"inf_speed"`
	NotAvailable bool   `json:"not_available"`
	// TODO: stats
}

func getSystemData(inf string) vstats {
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
	for _, line := range parts {
		if strings.HasPrefix(strings.TrimSpace(line), inf) {
			vstats.ProcNetDev = strings.TrimSpace(line)
			break
		}
	}

	infSpeedFile := fmt.Sprintf("/sys/class/net/%s/speed", inf)
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
	return vstats
}

func getStats(w http.ResponseWriter, r *http.Request) {
	inf := r.URL.Query().Get("inf.name")
	if inf == "" {
		// assume default eth0?
		inf = "eth0"
	}
	inf = strings.ReplaceAll(inf, ".", "")
	inf = strings.ReplaceAll(inf, "/", "")
	vstats := getSystemData(inf)
	encoder := json.NewEncoder(w)
	err := encoder.Encode(vstats)
	if err != nil {
		log.Printf("failed to write Varnish stats: %s", err)
	}
}

func main() {
	var port int
	flag.IntVar(&port, "port", 2000, "port to run vstats on")

	flag.Parse()
	http.HandleFunc("/", getStats)

	listenAddress := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(listenAddress, nil); err != nil {
		log.Printf("server stopped %s", err)
	}
}
