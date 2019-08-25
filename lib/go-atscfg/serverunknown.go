package atscfg

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
	"sort"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

func MakeServerUnknown(
	serverName tc.CacheName,
	serverDomain string,
	toToolName string, // tm.toolname global parameter (TODO: cache itself?)
	toURL string, // tm.url global parameter (TODO: cache itself?)
	params map[string][]string, // map[name]value - for the requested unknown 'take-and-bake' config_file
) string {

	hdr := GenericHeaderComment(string(serverName), toToolName, toURL)

	txt := ""

	sortedParams := SortParams(params)
	for _, pa := range sortedParams {
		if pa.Name == "location" {
			continue
		}
		if pa.Name == "header" {
			if pa.Val == "none" {
				hdr = ""
			} else {
				hdr = pa.Val + "\n"
			}
			continue
		}
		txt += pa.Val + "\n"
	}

	txt = strings.Replace(txt, `__HOSTNAME__`, string(serverName)+`.`+serverDomain, -1)
	txt = strings.Replace(txt, `__RETURN__`, "\n", -1)

	return hdr + txt
}

type Param struct {
	Name string
	Val  string
}

type Params []Param

func (a Params) Len() int           { return len(a) }
func (a Params) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Params) Less(i, j int) bool { return a[i].Name < a[j].Name }

func SortParams(params map[string][]string) []Param {
	sortedParams := []Param{}
	for name, vals := range params {
		for _, val := range vals {
			sortedParams = append(sortedParams, Param{Name: name, Val: val})
		}
	}
	sort.Sort(Params(sortedParams))
	return sortedParams
}
