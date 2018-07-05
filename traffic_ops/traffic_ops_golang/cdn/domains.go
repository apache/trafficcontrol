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

	"github.com/apache/trafficcontrol/lib/go-tc/v13"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/jmoiron/sqlx"
)

// GetDomainsList gathers a list of domains, except for the domain name.
// There seems to be an issue nesting queries (performing a query while
// rows.Next still needs to iterate). https://github.com/lib/pq/issues/81
func getDomainsList(tx *sqlx.Tx) ([]v13.Domain, []int, error) {

	var (
		cdn  int
		id   int
		name string
		desc string
	)

	domains := []v13.Domain{}
	cdn_ids := []int{}

	q := `SELECT cdn, id, name, description FROM Profile WHERE name LIKE 'CCR%'`
	rows, err := tx.Query(q)
	if err != nil {
		return nil, nil, fmt.Errorf("querying for profile: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&cdn, &id, &name, &desc); err != nil {
			return nil, nil, fmt.Errorf("getting profile: %s", err)
		}

		elem := v13.Domain{
			ProfileID:          id,
			ParameterID:        -1,
			ProfileName:        name,
			ProfileDescription: desc,
		}

		cdn_ids = append(cdn_ids, cdn)
		domains = append(domains, elem)
	}

	return domains, cdn_ids, nil
}

func DomainsHandler(w http.ResponseWriter, r *http.Request) {

	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	domains, cdn_ids, err := getDomainsList(inf.Tx)
	if err != nil {
		api.HandleErr(w, r, http.StatusInternalServerError, nil, err)
		return
	}

	for i, cdn := range cdn_ids {
		row := inf.Tx.QueryRow(`SELECT DOMAIN_NAME FROM CDN WHERE id = $1`, cdn)
		err := row.Scan(&domains[i].DomainName)
		if err != nil {
			api.HandleErr(w, r, http.StatusInternalServerError, nil, fmt.Errorf("getting domain name of cdn: %s", err))
			return
		}
	}

	api.WriteResp(w, r, domains)
}
