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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/apache/trafficcontrol/cache-config/t3cutil"
	"github.com/kylelemons/godebug/diff"
	"github.com/pborman/getopt/v2"
)

func main() {
	help := getopt.BoolLong("help", 'h', "Print usage info and exit")
	getopt.ParseV2()
	if *help {
		fmt.Println(usageStr)
		os.Exit(0)
	}

	if len(os.Args) < 3 {
		fmt.Println(usageStr)
		os.Exit(3)
	}

	fileNameA := strings.TrimSpace(os.Args[1])
	fileNameB := strings.TrimSpace(os.Args[2])

	if len(fileNameA) == 0 || len(fileNameB) == 0 {
		fmt.Println(usageStr)
		os.Exit(4)
	}

	fileA, fileAExisted, err := readFileOrStdin(fileNameA)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading first: "+err.Error())
		os.Exit(5)
	}
	fileB, fileBExisted, err := readFileOrStdin(fileNameB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading second: "+err.Error())
		os.Exit(6)
	}

	fileALines := strings.Split(string(fileA), "\n")
	fileALines = t3cutil.UnencodeFilter(fileALines)
	fileALines = t3cutil.CommentsFilter(fileALines)
	fileA = strings.Join(fileALines, "\n")
	fileA = t3cutil.NewLineFilter(fileA)

	fileBLines := strings.Split(string(fileB), "\n")
	fileBLines = t3cutil.UnencodeFilter(fileBLines)
	fileBLines = t3cutil.CommentsFilter(fileBLines)
	fileB = strings.Join(fileBLines, "\n")
	fileB = t3cutil.NewLineFilter(fileB)

	if fileA != fileB {
		match := regexp.MustCompile(`(?m)^\+.*|^-.*`)
		changes := diff.Diff(fileA, fileB)
		for _, change := range match.FindAllString(changes, -1) {
			fmt.Println(change)
		}
		os.Exit(1)
	}
	if fileAExisted != fileBExisted {
		os.Exit(1)
	}
	os.Exit(0)

}

const usageStr = `usage: t3c-diff [--help]
       <file-a> <file-b>

Either file may be 'stdin', in which case that file is read from stdin.
Either file may not exist.

Prints the diff to stdout, and returns the exit code 0 if there was no diff, 1 if there was a diff.
If one file exists but the other doesn't, it will always be a diff.

Note this means there may be no diff text printed to stdout but still exit 1 indicating a diff
if the file being created or deleted is semantically empty.`

// readFileOrStdin reads the file, or if fileOrStdin is 'stdin', reads from stdin.
// Returns the file, whether it existed, and any error.
func readFileOrStdin(fileOrStdin string) (string, bool, error) {
	if strings.ToLower(fileOrStdin) == "stdin" {
		bts, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return "", false, errors.New("reading stdin: " + err.Error())
		}
		return string(bts), true, nil
	}
	bts, err := ioutil.ReadFile(fileOrStdin)
	if err != nil {
		if os.IsNotExist(err) {
			return "", false, nil
		}
		return "", false, errors.New("reading file: " + err.Error())
	}
	return string(bts), true, nil
}
