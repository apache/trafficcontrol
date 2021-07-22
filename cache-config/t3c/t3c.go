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
	"github.com/apache/trafficcontrol/lib/go-log"
	"os"
	"path/filepath"
	"syscall" // TODO change to x/unix ?

	"github.com/pborman/getopt/v2"
)

var commands = map[string]struct{}{
	"apply":      struct{}{},
	"check":      struct{}{},
	"diff":       struct{}{},
	"generate":   struct{}{},
	"preprocess": struct{}{},
	"request":    struct{}{},
	"update":     struct{}{},
}

const ExitCodeSuccess = 0
const ExitCodeNoCommand = 1
const ExitCodeUnknownCommand = 2
const ExitCodeCommandErr = 3
const ExitCodeExeErr = 4

func main() {
	flagHelp := getopt.BoolLong("help", 'h', "Print usage information and exit")
	getopt.Parse()
	if *flagHelp {
		log.Errorln(usageStr())
		os.Exit(ExitCodeSuccess)
	}

	if len(os.Args) < 2 {
		log.Errorf("no command\n\n%s", usageStr())
		os.Exit(ExitCodeNoCommand)
	}

	cmd := os.Args[1]
	if _, ok := commands[cmd]; !ok {
		log.Errorf("unknown command\n%s", usageStr())
		os.Exit(ExitCodeUnknownCommand)
	}

	app := "t3c-" + cmd
	args := append([]string{app}, os.Args[2:]...)

	ex, err := os.Executable()
	if err != nil {
		log.Errorf("error getting application information: %s\n", err.Error())
		os.Exit(ExitCodeExeErr)
	}
	dir := filepath.Dir(ex)
	appDir := filepath.Join(dir, app) // TODO use path, not exact dir of this exe

	env := os.Environ()

	if err := syscall.Exec(appDir, args, env); err != nil {
		log.Errorf("error executing sub-command: %s\n", err.Error())
		os.Exit(ExitCodeCommandErr)
	}
}

func usageStr() string {
	return `usage: t3c [--help]
       <command> [<args>]

For the arguments of a command, see 't3c <command> --help'.

These are the available commands:

  apply      generate and apply configuration

  check      check that new config can be applied
  diff       diff config files, with logic like ignoring comments
  generate   generate configuration from Traffic Ops data
  preprocess preprocess generated config files
  request    request Traffic Ops data
  update     update a cache's queue and reval status in Traffic Ops
`
}
