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
	"path/filepath"
	"syscall" // TODO change to x/unix ?

	"github.com/apache/trafficcontrol/v8/cache-config/t3cutil"
	"github.com/apache/trafficcontrol/v8/lib/go-log"

	"github.com/pborman/getopt/v2"
)

const AppName = "t3c-check"

// Version is the application version.
// This is overwritten by the build with the current project version.
var Version = "0.4"

// GitRevision is the git revision the application was built from.
// This is overwritten by the build with the current project version.
var GitRevision = "nogit"

var commands = map[string]struct{}{
	"refs":   {},
	"reload": {},
}

const ExitCodeSuccess = 0
const ExitCodeNoCommand = 1
const ExitCodeUnknownCommand = 2
const ExitCodeCommandErr = 3
const ExitCodeExeErr = 4
const ExitCodeCommandLookupErr = 5

func main() {
	flagHelp := getopt.BoolLong("help", 'h', "Print usage information and exit")
	flagVersion := getopt.BoolLong("version", 'V', "Print version information and exit.")
	getopt.Parse()
	log.Init(os.Stderr, os.Stderr, os.Stderr, os.Stderr, os.Stderr)
	if *flagHelp {
		log.Errorln(usageStr())
		os.Exit(ExitCodeSuccess)
	} else if *flagVersion {
		fmt.Println(t3cutil.VersionStr(AppName, Version, GitRevision))
		os.Exit(ExitCodeSuccess)
	}

	if len(os.Args) < 2 {
		log.Errorf("no command\n\n" + usageStr())
		os.Exit(ExitCodeNoCommand)
	}

	cmd := os.Args[1]
	if _, ok := commands[cmd]; !ok {
		log.Errorf("unknown command\n%s", usageStr())
		os.Exit(ExitCodeUnknownCommand)
	}

	app := "t3c-check-" + cmd
	appPath := filepath.Join(t3cutil.InstallDir(), app)
	_, err := os.Stat(appPath)
	if err != nil {
		log.Errorf("error finding path to '%s': %s\n", app, err.Error())
		os.Exit(ExitCodeCommandLookupErr)
	}

	args := append([]string{app}, os.Args[2:]...)

	env := os.Environ()

	if err := syscall.Exec(appPath, args, env); err != nil {
		log.Errorf("error executing sub-command: %s\n", err.Error())
		os.Exit(ExitCodeCommandErr)
	}
}

func usageStr() string {
	return `usage: t3c-check [--help]
       <command> [<args>]

t3c-check has commands for checking things about new config files, such as
whether they can be safely applied or if a service reload or restart will
be required.

For the arguments of a command, see 't3c-check <command> --help'.

These are the available commands:

  reload  if a reload or restart is needed
  refs    if a config file's referenced plugins and files are valid
`
}
