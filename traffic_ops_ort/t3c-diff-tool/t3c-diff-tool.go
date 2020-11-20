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
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/apache/trafficcontrol/traffic_ops_ort/t3cutil"
	"github.com/kylelemons/godebug/diff"
	"github.com/pborman/getopt/v2"
)

func main() {
	different := false
	tropsFile := getopt.StringLong("trops-file", 't', "", "Required: Config file name in Traffic Ops")
	diskFile := getopt.StringLong("disk-file", 'd', "", "Required: Config file on disk")
	help := getopt.BoolLong("help", 'h', "Print usage info and exit")
	getopt.ParseV2()

	if *help {
		getopt.PrintUsage(os.Stdout)
		os.Exit(0)
	}
	if len(strings.TrimSpace(*tropsFile)) == 0 || len(strings.TrimSpace(*diskFile)) == 0 {
		getopt.PrintUsage(os.Stdout)
		os.Exit(1)
	}
	tropsInput := t3cutil.ReadFile(*tropsFile)
	diskInput := t3cutil.ReadFile(*diskFile)

	tropsData := strings.Split(string(tropsInput), "\n")
	tropsData = t3cutil.UnencodeFilter(tropsData)
	tropsData = t3cutil.CommentsFilter(tropsData)
	trops := strings.Join(tropsData, "\n")
	trops = t3cutil.NewLineFilter(trops)

	diskData := strings.Split(string(diskInput), "\n")
	diskData = t3cutil.UnencodeFilter(diskData)
	diskData = t3cutil.CommentsFilter(diskData)
	disk := strings.Join(diskData, "\n")
	disk = t3cutil.NewLineFilter(disk)

	if trops != disk {
		different = true
		match := regexp.MustCompile(`(?m)^\+.*|^\-.*`)
		changes := diff.Diff(disk, trops)
		for _, diff := range match.FindAllString(changes, -1) {
			fmt.Println(diff)
		}
	}
	fmt.Println(different)
}
