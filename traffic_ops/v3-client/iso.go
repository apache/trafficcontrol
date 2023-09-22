package client

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
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/toclientlib"
)

const (
	// API_OSVERSIONS is Deprecated: will be removed in the next major version. Be aware this may not be the URI being requested, for clients created with Login and ClientOps.ForceLatestAPI false.
	API_OSVERSIONS = apiBase + "/osversions"

	APIOSVersions = "/osversions"
)

// GetOSVersions GETs all available Operating System (OS) versions for ISO generation,
// as well as the name of the directory where the "kickstarter" files are found.
// Structure of returned map:
//
//	key:   Name of OS
//	value: Directory where the ISO source can be found
func (to *Session) GetOSVersions() (map[string]string, toclientlib.ReqInf, error) {
	var data struct {
		Versions tc.OSVersionsResponse `json:"response"`
	}
	reqInf, err := to.get(APIOSVersions, nil, &data)
	return data.Versions, reqInf, err
}
