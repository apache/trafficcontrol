/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package main

import (
	"bufio"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Documented []string

var documented Documented

func init() {
	f, err := os.Open(`LICENSE`)
	if err != nil {
		panic(err)
	}

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if len(line) != 0 && line[0] == '@' {
			documented = append(documented, line[1:])
		}
	}
}

func (d Documented) Documents(name string) bool {
	for _, re := range d {
		if ok, err := path.Match(re, name); ok && err == nil {
			return true
		}
	}
	dir := path.Dir(name)
	if dir != `` && dir != name {
		return d.Documents(dir)
	}
	return false
}

func (d Documented) Extra() []string {
	extra := make(map[string]struct{})
	for _, s := range d {
		extra[s] = struct{}{}
	}

	filepath.Walk(`.`, func(name string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Base(name) == `.git` {
			return filepath.SkipDir
		}

		if info.IsDir() {
			return nil
		}

		for re := range extra {
			if ok, err := path.Match(re, name); ok && err == nil {
				delete(extra, re)
			}
		}
		return nil
	})

	var extraDoc []string
	for re := range extra {
		extraDoc = append(extraDoc, re)
	}
	return extraDoc
}
