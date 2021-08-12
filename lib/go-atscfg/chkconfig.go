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
	"encoding/json"
	"sort"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

// ChkconfigFileName is the name of the ckconfig configuration file.
const ChkconfigFileName = `chkconfig`

// ChkconfigParamConfigFile is the ConfigFile value of Parameters that affect
// the generation of the chkconfig configuration file.
const ChkconfigParamConfigFile = `chkconfig`

// ContentTypeChkconfig is the MIME content type of the contents of the
// chkconfig configuration file.
//
// Note that the GoDoc for MakeChkconfig says "This is a JSON object, and should
// be served with an 'application/json' Content-Type." but actually the file
// contents on disk are not JSON-encoded.
const ContentTypeChkconfig = ContentTypeTextASCII

// LineCommentChkconfig is the string that signifies the start of a line comment
// in the grammar of a chkconfig configuration file.
const LineCommentChkconfig = LineCommentHash

// ChkconfigOpts contains settings to configure generation options.
type ChkconfigOpts struct {
}

// MakeChkconfig returns the 'chkconfig' ATS config file endpoint.
//
// This is a JSON object, and should be served with an 'application/json'
// Content-Type.
//
// TODO: rename/rework? We systemd now, after all. Also, this may be unused as
// t3c now generates the contents of the file specially without calling into
// this function, possibly.
func MakeChkconfig(
	serverParams []tc.ParameterV5,
	opt *ChkconfigOpts,
) (Cfg, error) {
	if opt == nil {
		opt = &ChkconfigOpts{}
	}
	warnings := []string{}

	serverParams = filterParams(serverParams, ChkconfigParamConfigFile, "", "", "")

	chkconfig := []chkConfigEntry{}
	for _, param := range serverParams {
		chkconfig = append(chkconfig, chkConfigEntry{Name: param.Name, Val: param.Value})
	}

	sort.Sort(chkConfigEntries(chkconfig))

	bts, err := json.Marshal(&chkconfig)
	if err != nil {
		return Cfg{}, makeErr(warnings, "marshalling chkconfig NameVals: "+err.Error())
	}

	return Cfg{
		Text:        string(bts),
		ContentType: ContentTypeChkconfig,
		LineComment: LineCommentChkconfig,
		Warnings:    warnings,
	}, nil
}

type chkConfigEntry struct {
	Name string
	Val  string
}

type chkConfigEntries []chkConfigEntry

func (e chkConfigEntries) Len() int { return len(e) }
func (e chkConfigEntries) Less(i, j int) bool {
	if e[i].Name != e[j].Name {
		return e[i].Name < e[j].Name
	}
	return e[i].Val < e[j].Val
}
func (e chkConfigEntries) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
