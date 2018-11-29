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
	"time"
)

// makeCRConfigConfig returns the CRConfig Config object, the last updated time of any parameter, and any error.
func makeCRConfigConfig(cdn string, tx *sql.Tx, dnssecEnabled bool, domain string) (map[string]interface{}, time.Time, error) {
	configParams, lastUpdated, err := getConfigParams(cdn, tx)
	if err != nil {
		return nil, time.Time{}, errors.New("Error getting router params: " + err.Error())
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
				return nil, time.Time{}, errors.New("Error parsing " + maxmindDefaultOverrideParameterName + " parameter: " + err.Error())
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

	return crConfigConfig, lastUpdated, nil
}

type CRConfigConfigParameter struct {
	Name  string
	Value string
}

// getConfigParams returns the list of parameters, the last updated time of any parameter, and any error
func getConfigParams(cdn string, tx *sql.Tx) ([]CRConfigConfigParameter, time.Time, error) {
	// TODO change to []struct{string,string} ? Speed might matter.
	// TODO verify MAX(GREATEST(last_updated) shouldn't include profile, server, or cdn
	qry := `
SELECT
  pa.name,
  pa.value,
  MAX(GREATEST(pa.last_updated, pp.last_updated)) as last_updated
FROM
  parameter pa
  JOIN profile_parameter pp ON pp.parameter = pa.id
  JOIN server s ON s.profile = pp.profile
  JOIN cdn ON cdn.id = s.cdn_id
WHERE
  cdn.name = $1
  AND pa.config_file = 'CRConfig.json'
GROUP BY pa.name, pa.value
`
	rows, err := tx.Query(qry, cdn)
	if err != nil {
		return nil, time.Time{}, errors.New("Error querying router params: " + err.Error())
	}
	defer rows.Close()

	lastUpdated := time.Time{}

	params := []CRConfigConfigParameter{}
	for rows.Next() {
		name := ""
		val := ""
		paLastUpdated := time.Time{}
		if err := rows.Scan(&name, &val, &paLastUpdated); err != nil {
			return nil, time.Time{}, errors.New("Error scanning router param: " + err.Error())
		}
		params = append(params, CRConfigConfigParameter{Name: name, Value: val})
		if paLastUpdated.After(lastUpdated) {
			lastUpdated = paLastUpdated
		}
	}
	if err := rows.Err(); err != nil {
		return nil, time.Time{}, errors.New("Error iterating router param rows: " + err.Error())
	}
	return params, lastUpdated, nil
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
