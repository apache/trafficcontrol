package crconfig

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
	"database/sql"
	"errors"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

func Make(tx *sql.Tx, cdn, user, toHost, reqPath, toVersion string) (*tc.CRConfig, error) {
	crc := tc.CRConfig{APIVersion: 1.4}
	err := error(nil)

	cdnDomain, dnssecEnabled, err := getCDNInfo(cdn, tx)
	if err != nil {
		return nil, errors.New("Error getting CDN info: " + err.Error())
	}
	if crc.Config, err = makeCRConfigConfig(cdn, tx, dnssecEnabled, cdnDomain); err != nil {
		return nil, errors.New("Error getting Config: " + err.Error())
	}
	serverDSNames, err := getServerDSNames(cdn, tx)
	if err != nil {
		return nil, errors.New("Error getting server delivery services: " + err.Error())
	}
	if crc.ContentServers, crc.ContentRouters, crc.Monitors, err = makeCRConfigServers(cdn, tx, cdnDomain, serverDSNames); err != nil {
		return nil, errors.New("Error getting Servers: " + err.Error())
	}
	if crc.EdgeLocations, crc.RouterLocations, err = makeLocations(cdn, tx); err != nil {
		return nil, errors.New("Error getting Edge Locations: " + err.Error())
	}
	if crc.DeliveryServices, err = makeDSes(cdn, cdnDomain, serverDSNames, tx); err != nil {
		return nil, errors.New("Error getting Delivery Services: " + err.Error())
	}

	// TODO change to real reqPath, and verify everything works. Currently emulates the existing TO, in case anything relies on it
	emulateOldPath := "/tools/write_crconfig/" + cdn
	crc.Stats = makeStats(cdn, user, toHost, emulateOldPath, toVersion)
	return &crc, nil
}
