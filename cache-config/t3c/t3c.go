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
)

var commands = map[string]struct{}{
	"apply":  struct{}{},
	"update": struct{}{},
}

const ExitCodeNoCommand = 1
const ExitCodeUnknownCommand = 2
const ExitCodeCommandErr = 3
const ExitCodeExeErr = 4

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "no command\n") // TODO print usage
		os.Exit(ExitCodeNoCommand)
	}

	cmd := os.Args[1]
	if _, ok := commands[cmd]; !ok {
		fmt.Fprintf(os.Stderr, "unknown command\n") // TODO print usage
		os.Exit(ExitCodeUnknownCommand)
	}

	app := "t3c-" + cmd
	args := append([]string{app}, os.Args[2:]...)

	ex, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting application information: "+err.Error()+"\n")
		os.Exit(ExitCodeExeErr)
	}
	dir := filepath.Dir(ex)
	appDir := filepath.Join(dir, app) // TODO use path, not exact dir of this exe

	env := os.Environ()

	if err := syscall.Exec(appDir, args, env); err != nil {
		fmt.Fprintf(os.Stderr, "error executing sub-command: "+err.Error()+"\n")
		os.Exit(ExitCodeCommandErr)
	}
}
