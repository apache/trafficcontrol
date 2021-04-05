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
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var remapHeaderRewriteRegex = regexp.MustCompile(`\@plugin=header_rewrite\.so @pparam=([^\s]+)\s+.*\n`)

func main() {
	bts, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println("Error reading input: " + err.Error())
		os.Exit(1)
	}

	matches := remapHeaderRewriteRegex.FindAllSubmatch(bts, -1)

	headerRewriteFiles := []string{}

	for _, match := range matches {
		fileName := string(match[1])
		headerRewriteFiles = append(headerRewriteFiles, fileName)
	}

	anyFileMissing := false
	for _, fileName := range headerRewriteFiles {
		reStr := `\r\nPath\: .*\/` + strings.Replace(fileName, `.`, `\.`, -1) + "\r"
		re := regexp.MustCompile(reStr)
		if re.FindIndex(bts) == nil {
			fmt.Println(fileName)
			anyFileMissing = true
		}
	}

	if anyFileMissing {
		os.Exit(1)
	}
	os.Exit(0)
}
