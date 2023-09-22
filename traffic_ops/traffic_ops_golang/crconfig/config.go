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
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

func makeCRConfigConfig(cdn string, tx *sql.Tx, dnssecEnabled bool, domain string) (map[string]interface{}, error) {
	configParams, err := getConfigParams(cdn, tx)
	if err != nil {
		return nil, errors.New("Error getting router params: " + err.Error())
	}
	soa := map[string]string{}
	ttl := map[string]string{}
	maxmindDefaultOverrides := []CRConfigConfigMaxmindDefaultOverride{}

	const soaPrefix = "tld.soa."
	const ttlPrefix = "tld.ttls."
	const maxmindDefaultOverrideParameterName = "maxmind.default.override"
	const logRequestHeadersParameterName = "LogRequestHeaders"
	crConfigConfig := map[string]interface{}{}
	for _, param := range configParams {
		k := param.Name
		v := param.Value
		if strings.HasPrefix(k, soaPrefix) {
			soa[k[len(soaPrefix):]] = v
		} else if strings.HasPrefix(k, ttlPrefix) {
			ttl[k[len(ttlPrefix):]] = v
		} else if k == logRequestHeadersParameterName {
			hdrs := []string{}
			for _, hdr := range strings.Split(param.Value, `__RETURN__`) {
				hdrs = append(hdrs, strings.TrimSpace(hdr))
			}
			crConfigConfig["requestHeaders"] = hdrs
		} else if k == maxmindDefaultOverrideParameterName {
			overrideObj, err := createMaxmindDefaultOverrideObj(v)
			if err != nil {
				return nil, errors.New("Error parsing " + maxmindDefaultOverrideParameterName + " parameter: " + err.Error())
			}
			maxmindDefaultOverrides = append(maxmindDefaultOverrides, overrideObj)
		} else {
			crConfigConfig[k] = v
		}
	}
	crConfigConfig["domain_name"] = domain
	if len(soa) > 0 {
		crConfigConfig["soa"] = soa
	}
	if len(ttl) > 0 {
		crConfigConfig["ttls"] = ttl
	}
	if len(maxmindDefaultOverrides) > 0 {
		crConfigConfig["maxmindDefaultOverride"] = maxmindDefaultOverrides
	}
	dnssecStr := "false"
	if dnssecEnabled {
		dnssecStr = "true"
	}
	crConfigConfig["dnssec.enabled"] = dnssecStr

	return crConfigConfig, nil
}

type CRConfigConfigParameter struct {
	Name  string
	Value string
}

func getConfigParams(cdn string, tx *sql.Tx) ([]CRConfigConfigParameter, error) {
	// TODO change to []struct{string,string} ? Speed might matter.
	q := `
select name, value from parameter where id in (
  select parameter from profile_parameter where profile in (
  	select distinct profile from server where cdn_id = (
	    select id from cdn where name = $1
    ) and type = (
        select id from type where name = $2
    )
  )
)
and config_file = 'CRConfig.json'
`
	rows, err := tx.Query(q, cdn, tc.RouterTypeName)
	if err != nil {
		return nil, errors.New("Error querying router params: " + err.Error())
	}
	defer log.Close(rows, "closing crconfig parameter rows")

	params := []CRConfigConfigParameter{}
	for rows.Next() {
		name := ""
		val := ""
		if err := rows.Scan(&name, &val); err != nil {
			return nil, errors.New("Error scanning router param: " + err.Error())
		}
		params = append(params, CRConfigConfigParameter{Name: name, Value: val})
	}
	if err := rows.Err(); err != nil {
		return nil, errors.New("Error iterating router param rows: " + err.Error())
	}
	return params, nil
}

type CRConfigConfigMaxmindDefaultOverride struct {
	CountryCode string  `json:"countryCode"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"long"`
}

func createMaxmindDefaultOverrideObj(maxmindDefaultOverrideParamVal string) (CRConfigConfigMaxmindDefaultOverride, error) {
	countryCodeCoords := strings.Split(maxmindDefaultOverrideParamVal, ";")
	if len(countryCodeCoords) < 2 {
		return CRConfigConfigMaxmindDefaultOverride{}, errors.New("malformed maxmind.default.override parameter: '" + maxmindDefaultOverrideParamVal + "'")
	}
	countryCode := countryCodeCoords[0]
	coords := countryCodeCoords[1]
	latLon := strings.Split(coords, ",")
	if len(latLon) < 2 {
		return CRConfigConfigMaxmindDefaultOverride{}, errors.New("malformed maxmind.default.override parameter coordinates '" + maxmindDefaultOverrideParamVal + "'")
	}
	latStr := latLon[0]
	lonStr := latLon[1]
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return CRConfigConfigMaxmindDefaultOverride{}, errors.New("malformed maxmind.default.override parameter coordinates, latitude not a number: '" + maxmindDefaultOverrideParamVal + "'")
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		return CRConfigConfigMaxmindDefaultOverride{}, errors.New("malformed maxmind.default.override parameter coordinates, longitude not an number: '" + maxmindDefaultOverrideParamVal + "'")
	}
	return CRConfigConfigMaxmindDefaultOverride{
		CountryCode: countryCode,
		Lat:         lat,
		Lon:         lon,
	}, nil
}
