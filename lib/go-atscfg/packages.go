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

const PackagesFileName = `packages`
const PackagesParamConfigFile = `package`

type Package struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// MakePackages returns the 'packages' ATS config file endpoint.
// This is a JSON object, and should be served with an 'application/json' Content-Type.
func MakePackages(
	params map[string][]string, // map[name]value - config file should always be 'package'
) string {
	packages := []Package{}
	for name, versions := range params {
		for _, version := range versions {
			packages = append(packages, Package{Name: name, Version: version})
		}
	}
	bts, err := json.Marshal(&packages)
	if err != nil {
		// should never happen
		log.Errorln("marshalling chkconfig NameVersions: " + err.Error())
		bts = []byte("error encoding params to json, see Traffic Ops log for details")
	}
	return string(bts)
}
