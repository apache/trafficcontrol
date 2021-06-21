package t3cutil

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
	"html"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type ATSConfigFile struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	ContentType string `json:"content_type"`
	LineComment string `json:"line_comment"`
	Text        string `json:"text"`
}

// ATSConfigFiles implements sort.Interface and sorts by the Location and then FileNameOnDisk, i.e. the full file path.
type ATSConfigFiles []ATSConfigFile

func (fs ATSConfigFiles) Len() int { return len(fs) }
func (fs ATSConfigFiles) Less(i, j int) bool {
	if fs[i].Path != fs[j].Path {
		return fs[i].Path < fs[j].Path
	}
	return fs[i].Name < fs[j].Name
}
func (fs ATSConfigFiles) Swap(i, j int) { fs[i], fs[j] = fs[j], fs[i] }

// commentsFilter is used to remove comment
// lines from config files while making
// comparisons.
func CommentsFilter(body []string) []string {
	var newlines []string

	newlines = make([]string, 0)

	for ii := range body {
		line := body[ii]
		if strings.HasPrefix(line, "#") {
			continue
		}
		newlines = append(newlines, line)
	}

	return newlines
}

// NewLineFilter removes carriage returns
// from config files while making comparisons.
func NewLineFilter(str string) string {
	str = strings.ReplaceAll(str, "\r\n", "\n")
	return strings.TrimSpace(str)
}

// ReadFile reads a file and returns the
// file contents.
func ReadFile(f string) []byte {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		fmt.Println("Error reading file ", f)
		os.Exit(1)
	}
	return data
}

// UnencodeFilter translates HTML escape
// sequences while making config file comparisons.
func UnencodeFilter(body []string) []string {
	var newlines []string

	newlines = make([]string, 0)
	sp := regexp.MustCompile(`\s+`)
	el := regexp.MustCompile(`^\s+|\s+$`)

	for ii := range body {
		s := body[ii]
		s = sp.ReplaceAllString(s, " ")
		s = el.ReplaceAllString(s, "")
		s = html.UnescapeString(s)
		s = strings.TrimSpace(s)
		newlines = append(newlines, s)
	}

	return newlines
}

// Do executes the given command and returns the stdout, stderr, and exit code.

// This is a convenience wrapper around os/exec.
// Since t3c only needs to make simple calls and get the stdout, stderr, and code, this provides a simpler and terser interface.
//
// If you need anything more complex, or don't find this simpler, you should probably use os/exec directly.
//
// Each arg must be passed as its own string. Unfortunately, Go doesn't have a way to pass multiple args as a single string, and splitting on spaces would require complex quote parsing.
//
// Note each arg must be passed without quotes. Go calls the app with args as if they were quoted. if you add quotes, they'll be passed to the command literally, as if you called 'mycommand "\"escaped-quotes\""`.
//
// Note if Go fails to run the command, the error from Go will be returned as the stderr and the code -1,
// which will differ from what would have been returned by a command line.
//
func Do(cmdStr string, args ...string) ([]byte, []byte, int) {
	cmd := exec.Command(cmdStr, args...)

	var outbuf bytes.Buffer
	var errbuf bytes.Buffer

	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	code := 0
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); !ok {
			return nil, []byte(err.Error()), -1
		} else {
			code = exitErr.ExitCode()
		}
	}

	return outbuf.Bytes(), errbuf.Bytes(), code
}

// DoInput is like Do but takes the stdin to pass to the command.
func DoInput(input []byte, cmdStr string, args ...string) ([]byte, []byte, int) {
	cmd := exec.Command(cmdStr, args...)

	var outbuf bytes.Buffer
	var errbuf bytes.Buffer

	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	cmd.Stdin = bytes.NewBuffer(input)

	code := 0
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); !ok {
			return nil, []byte(err.Error()), -1
		} else {
			code = exitErr.ExitCode()
		}
	}

	return outbuf.Bytes(), errbuf.Bytes(), code
}
