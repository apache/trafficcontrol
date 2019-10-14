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

	"github.com/apache/trafficcontrol/lib/go-log"
)

const ChkconfigFileName = `chkconfig`

const ChkconfigParamConfigFile = `chkconfig`

type ChkConfigEntry struct {
	Name string `json:"name"`
	Val  string `json:"value"`
}

// MakeChkconfig returns the 'chkconfig' ATS config file endpoint.
// This is a JSON object, and should be served with an 'application/json' Content-Type.
func MakeChkconfig(
	params map[string][]string, // map[name]value - config file should always be 'chkconfig'
) string {

	chkconfig := []ChkConfigEntry{}
	for name, vals := range params {
		for _, val := range vals {
			chkconfig = append(chkconfig, ChkConfigEntry{Name: name, Val: val})
		}
	}
	bts, err := json.Marshal(&chkconfig)
	if err != nil {
		// should never happen
		log.Errorln("marshalling chkconfig NameVals: " + err.Error())
		bts = []byte("error encoding params to json, see Traffic Ops log for details")
	}
	return string(bts)
}
