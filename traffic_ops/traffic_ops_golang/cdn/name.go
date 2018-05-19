package cdn

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
	"net/http"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	tcv13 "github.com/apache/incubator-trafficcontrol/lib/go-tc/v13"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
)

func GetName(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, _, userErr, sysErr, errCode := api.AllParams(r, nil)
		if userErr != nil || sysErr != nil {
			api.HandleErr(w, r, errCode, userErr, sysErr)
			return
		}
		cdnName, ok := params["name"]
		if !ok {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("CDN name route missing param"))
			return
		}
		api.RespWriter(w, r)(getCDNFromName(db, tc.CDNName(cdnName)))
	}
}

func getCDNFromName(db *sql.DB, name tc.CDNName) ([]tcv13.CDN, error) {
	rows, err := db.Query(`SELECT id, domain_name, last_updated, dnssec_enabled FROM cdn WHERE name = $1`, name)
	if err != nil {
		return nil, errors.New("querying cdns: " + err.Error())
	}
	cdns := []tcv13.CDN{}
	for rows.Next() {
		cdn := tcv13.CDN{Name: string(name)}
		if err := rows.Scan(&cdn.ID, &cdn.DomainName, &cdn.LastUpdated, &cdn.DNSSECEnabled); err != nil {
			return nil, errors.New("scanning cdns: " + err.Error())
		}
		cdns = append(cdns, cdn)
	}
	return cdns, nil
}
