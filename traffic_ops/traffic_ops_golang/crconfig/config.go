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
	"strings"
)

func makeCRConfigConfig(cdn string, db *sql.DB, dnssecEnabled bool) (map[string]interface{}, error) {
	configParams, err := getConfigParams(cdn, db)
	if err != nil {
		return nil, errors.New("Error getting router params: " + err.Error())
	}
	soa := map[string]string{}
	ttl := map[string]string{}
	const soaPrefix = "tld.soa."
	ttlPrefix := "tld.ttls."
	crConfigConfig := map[string]interface{}{}
	for k, v := range configParams {
		if strings.HasPrefix(k, soaPrefix) {
			soa[k[len(soaPrefix):]] = v
		} else if strings.HasPrefix(k, ttlPrefix) {
			ttl[k[len(ttlPrefix):]] = v
		} else {
			crConfigConfig[k] = v
		}
	}
	if len(soa) > 0 {
		crConfigConfig["soa"] = soa
	}
	if len(ttl) > 0 {
		crConfigConfig["ttls"] = ttl
	}

	dnssecStr := "false"
	if dnssecEnabled {
		dnssecStr = "true"
	}
	crConfigConfig["dnssec.enabled"] = dnssecStr

	return crConfigConfig, nil
}

func getConfigParams(cdn string, db *sql.DB) (map[string]string, error) {
	// TODO change to []struct{string,string} ? Speed might matter.
	q := `
select name, value from parameter where id in (
  select parameter from profile_parameter where profile in (
  	select distinct profile from server where cdn_id = (
	    select id from cdn where name = $1
    )
  )
)
and config_file = 'CRConfig.json'
`
	rows, err := db.Query(q, cdn)
	if err != nil {
		return nil, errors.New("Error querying router params: " + err.Error())
	}
	defer rows.Close()

	params := map[string]string{}
	for rows.Next() {
		name := ""
		val := ""
		if err := rows.Scan(&name, &val); err != nil {
			return nil, errors.New("Error scanning router param: " + err.Error())
		}
		params[name] = val
	}
	if err := rows.Err(); err != nil {
		return nil, errors.New("Error iterating router param rows: " + err.Error())
	}
	return params, nil
}
