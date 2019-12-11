package iso

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
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// environment variables used by the mock command
const (
	mockCmdEnvInvoke = "GO_HELPER_CMD"
	mockCmdEnvError  = "GO_HELPER_CMD_FORCE_ERROR"
	mockCmdEnvOutput = "GO_HELPER_CMD_OUTPUT"
)

// mockISOCmd returns a modified version of the given Cmd
// so that when run, the command actually invokes the
// TestHelperMockCmd test. See TestHelperMockCmd for
// more details on its behavior.
//
// - If forceError is true, the command will exit a non-0 code and write to STDERR.
// - If cmdOutput is blank, the command will write to STDOUT, otherwise
// it will write its output to the file specified by cmdOutput.
func mockISOCmd(cmd *exec.Cmd, forceError bool, cmdOutput string) *exec.Cmd {
	args := []string{
		"-test.run=TestHelperMockCmd",
		"--",
	}
	args = append(args, cmd.Args...)

	// os.Args[0] is the invokation of this test binary
	mocked := exec.Command(os.Args[0], args...)

	env := cmd.Env
	env = append(cmd.Env, fmt.Sprintf("%s=1", mockCmdEnvInvoke))
	if forceError {
		env = append(env, fmt.Sprintf("%s=1", mockCmdEnvError))
	}
	if cmdOutput != "" {
		env = append(env, fmt.Sprintf("%s=%s", mockCmdEnvOutput, cmdOutput))
	}
	mocked.Env = env

	return mocked
}

// TestHelperMockCmd is a special test case that is meant to be invoked
// by a subprocess, e.g. go test -run=TestHelperMockCmd.
//
// Described in detail at: https://npf.io/2015/06/testing-exec-command/
//
// In order for the test to act like a subprocess, the GO_HELPER_CMD environment
// variable must be set, otherwise the test is skipped. Use the mockISOCmd to
// assist with modifying an existing exec.Cmd to use this helper.
//
// The test writes to either STDOUT/STDERR/or a file all arguments it receives
// after the '--' argument. In practice, the mockISOCmd passes the original
// command's arguments in this position, essentially echoing back what the
// original command was.
func TestHelperMockCmd(t *testing.T) {
	if os.Getenv(mockCmdEnvInvoke) != "1" {
		return
	}

	var (
		respCode int
		dest     io.Writer
	)

	switch {
	case os.Getenv(mockCmdEnvError) == "1":
		respCode = 1
		dest = os.Stderr

	case os.Getenv(mockCmdEnvOutput) != "":
		fd, err := os.Create(os.Getenv(mockCmdEnvOutput))
		if err == nil {
			defer fd.Close()
			dest = fd
		} else {
			respCode = 100
			dest = os.Stderr
		}

	default:
		dest = os.Stdout
	}

	// Set args to all arguments past '--'.
	var args []string
	for i, v := range os.Args {
		if v == "--" {
			args = os.Args[i+1:]
			break
		}
	}

	fmt.Fprintf(dest, strings.Join(args, " "))
	os.Exit(respCode)
}
