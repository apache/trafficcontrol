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
	"fmt"
	"html"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

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
