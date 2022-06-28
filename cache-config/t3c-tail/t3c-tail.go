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
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/cache-config/t3cutil"
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

//default time out is 15 seconds, if not included in json input.
var timeOutSeconds = 15



 func main() {
	version := getopt.BoolLong("version", 'V', "Print version information and exit.")
	help := getopt.BoolLong("help", 'h', "Print usage information and exit")
	getopt.Parse()

	if *help {
		fmt.Println(usageStr())
		os.Exit(0)
	} else if *version {
		fmt.Println(t3cutil.VersionStr(AppName, Version, GitRevision))
		os.Exit(0)
	}

	tailCfg := &t3cutil.TailCfg{}
	if err := json.NewDecoder(os.Stdin).Decode(tailCfg); err != nil {
		fmt.Println("Error reading json input", err)
	}

	if tailCfg.LogMatch == nil {
		fmt.Println("must provide a regex to match")
		fmt.Println(usageStr())
		os.Exit(1)
	}
	var endMatch string
	if tailCfg.EndMatch != nil {
		endMatch = *tailCfg.EndMatch
	} else {
		endMatch = "use timeout"
	}
	
	logMatch := regexp.MustCompile(*tailCfg.LogMatch)
	timeOut := timeOutSeconds
	if tailCfg.TimeOut != nil {
		timeOut = *tailCfg.TimeOut
	}
	
	file := tailCfg.File
	t, err := tail.TailFile(*file,
		tail.Config {
			MustExist: true,
			Follow: true ,
			Location: &tail.SeekInfo {
				Offset: 0,
				Whence: 2,
				},
			})
	if err != nil {
		fmt.Println("error running tail on ", file)
		os.Exit(1)
	}
	go func() {
		for line := range t.Lines {
			if logMatch.MatchString(line.Text) {
				fmt.Println(line.Text)
			}
			if strings.Contains(line.Text, endMatch)  {
				fmt.Println("Stopping on stop match")
				break
			}
		}
	}()
	
	time.Sleep(time.Second * time.Duration(timeOut))
	fmt.Println("stopping with timeout")
	err = t.Stop()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
	}
	t.Cleanup()
	
}

func usageStr() string {
	return `usage: t3c-tail [--help]
	accepts json input from stdin in the following format:
	file is file you want to tail
	match is regex string you wish to match on, if you want everything use '.*'
	stopMatch is a string used to exit tail when it is found in the logs
	timeOut is given in seconds the default is 15
	{"file":"diags.log", "match":"<regex to match>", "endMatch": "<string>", "timeOut": 4}
	`
}