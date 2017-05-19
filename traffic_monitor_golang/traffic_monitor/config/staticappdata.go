package config

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
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// StaticAppData encapsulates data about the app available at startup
type StaticAppData struct {
	StartTime      time.Time
	GitRevision    string
	FreeMemoryMB   uint64
	Version        string
	WorkingDir     string
	Name           string
	BuildTimestamp string
	Hostname       string
	UserAgent      string
}

// getStaticAppData returns app data available at start time.
// This should be called immediately, as it includes calculating when the app was started.
func GetStaticAppData(version, gitRevision, buildTimestamp string) (StaticAppData, error) {
	var d StaticAppData
	var err error
	d.StartTime = time.Now()
	d.GitRevision = gitRevision
	d.FreeMemoryMB = math.MaxUint64 // TODO remove if/when nothing needs this
	d.Version = version
	if d.WorkingDir, err = os.Getwd(); err != nil {
		return StaticAppData{}, err
	}
	d.Name = os.Args[0]
	d.BuildTimestamp = buildTimestamp
	if d.Hostname, err = getHostNameWithoutDomain(); err != nil {
		return StaticAppData{}, err
	}

	d.UserAgent = fmt.Sprintf("%s/%s", filepath.Base(d.Name), d.Version)
	return d, nil
}

// getHostNameWithoutDomain returns the machine hostname, without domain information.
// Modified from http://stackoverflow.com/a/34331660/292623
func getHostNameWithoutDomain() (string, error) {
	cmd := exec.Command("/bin/hostname", "-s")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	hostname := out.String()
	if len(hostname) < 1 {
		return "", fmt.Errorf("OS returned empty hostname")
	}
	hostname = hostname[:len(hostname)-1] // removing EOL
	return hostname, nil
}
