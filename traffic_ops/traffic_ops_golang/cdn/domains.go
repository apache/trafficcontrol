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
	"fmt"
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"

	"github.com/jmoiron/sqlx"
)

const RouterProfilePrefix = "CCR"

func getDomainsList(tx *sqlx.Tx) ([]tc.Domain, error) {

	domains := []tc.Domain{}

	q := `SELECT p.id, p.name, p.description, domain_name FROM profile AS p
	JOIN cdn ON p.cdn = cdn.id WHERE p.name LIKE '` + RouterProfilePrefix + `%'`

	rows, err := tx.Query(q)
	if err != nil {
		return nil, fmt.Errorf("querying for profile: %s", err)
	}
	defer rows.Close()

	for rows.Next() {

		d := tc.Domain{ParameterID: -1}
		err := rows.Scan(&d.ProfileID, &d.ProfileName, &d.ProfileDescription, &d.DomainName)
		if err != nil {
			return nil, fmt.Errorf("getting profile: %s", err)
		}
		domains = append(domains, d)
	}

	return domains, nil
}

func DomainsHandler(w http.ResponseWriter, r *http.Request) {

	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	domains, err := getDomainsList(inf.Tx)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, err, err)
		return
	}

	api.WriteResp(w, r, domains)
}
