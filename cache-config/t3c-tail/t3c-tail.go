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
	"time"

	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/lib/go-log"

	"github.com/nxadm/tail"
	"github.com/pborman/getopt/v2"
)

const AppName = "t3c-tail"

// Version is the application version.
// This is overwritten by the build with the current project version.
var Version = "0.4"

// GitRevision is the git revision the application was built from.
// This is overwritten by the build with the current project version.
var GitRevision = "nogit"

// defaultTimeOutMs is 15000 milliseconds, if not included in input.
var defaultTimeOutMs = 15000

func main() {
	file := getopt.StringLong("file", 'f', "", "Path to file to watch")
	match := getopt.StringLong("match", 'm', ".*", "Regex pattern you want to match while running tail default is .*")
	endMatch := getopt.StringLong("end-match", 'e', "^timeout", "Regex pattern that will cause tail to exit before timeout")
	timeOutMs := getopt.Int64Long("timeout-ms", 't', int64(defaultTimeOutMs), "Timeout in milliseconds that will cause tail to exit default is 15000 MS")
	version := getopt.BoolLong("version", 'V', "Print version information and exit.")
	help := getopt.BoolLong("help", 'h', "Print usage information and exit")
	getopt.Parse()

	log.Init(os.Stderr, os.Stderr, os.Stderr, os.Stderr, os.Stderr)

	if *help {
		fmt.Println(usageStr())
		os.Exit(0)
	} else if *version {
		fmt.Println(t3cutil.VersionStr(AppName, Version, GitRevision))
		os.Exit(0)
	}

	if *file == "" || file == nil {
		fmt.Println("Please provide file path for t3c-tail")
		fmt.Println(usageStr())
		os.Exit(1)
	}

	logMatch := regexp.MustCompile(*match)
	tailStop := regexp.MustCompile(*endMatch)
	timeOut := *timeOutMs

	t, err := tail.TailFile(*file,
		tail.Config{
			MustExist: true,
			Follow:    true,
			Location: &tail.SeekInfo{
				Offset: 0,
				Whence: 2,
			},
		})
	if err != nil {
		log.Errorln("error running tail on ", file, err)
		os.Exit(1)
	}
	timer := time.NewTimer(time.Millisecond * time.Duration(timeOut))
	go func() {
		for line := range t.Lines {
			if logMatch.MatchString(line.Text) {
				fmt.Println(line.Text)
			}
			if tailStop.MatchString(line.Text) {
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(0)
				break
			}
		}
	}()

	<-timer.C
	t.Cleanup()
}

func usageStr() string {
	return `usage: t3c-tail [--help]
	-f <path to file> -m <regex to match> -e <regex match to exit> -t <timeout in ms>

	file is  path to the file you want to tail

	match is regex string you wish to match on, default is '.*'

	endMatch is a regex used to exit tail when it is found in the logs with out waiting for timeout

	timeOutMs is when tail will stop if endMatch isn't found default is 15000
	`
}
