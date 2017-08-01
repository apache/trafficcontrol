package main

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
	"fmt"
	"net/http"
	"strings"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	toclient "github.com/apache/incubator-trafficcontrol/traffic_ops/client"

	"github.com/jmoiron/sqlx"
)

const CacheConfigPrivLevel = 10

func cacheconfigHandler(db *sqlx.DB) RegexHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p PathParams) {
		handleErr := func(err error, status int) {
			log.Errorf("%v %v\n", r.RemoteAddr, err)
			w.WriteHeader(status)
			fmt.Fprintf(w, http.StatusText(status))
		}

		hostname := p["cache"]

		resp, err := getCacheConfigJSON(hostname, db)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		respBts, err := json.Marshal(resp)
		if err != nil {
			handleErr(err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", respBts)
	}
}

const CacheConfigCDNDomainQuery = `
SELECT name AS cdn, domain_name AS domain
FROM cdn
WHERE id = (
  SELECT cdn_id
  FROM server
  WHERE host_name = $1
)
`

const CacheConfigDeliveryServicesQuery = `
SELECT d.protocol, d.qstring_ignore AS query_string, d.xml_id AS name, t.name AS type, d.org_server_fqdn AS origin_uri, d.dscp
FROM deliveryservice AS d
INNER JOIN type AS t on t.id = d.type
WHERE d.id IN (
  SELECT deliveryservice
  FROM deliveryservice_server
  WHERE server = (
    SELECT id
    FROM server
    WHERE host_name = $1
  )
)
`

const CacheConfigParentsQuery = `
WITH child AS (
  SELECT cachegroup, cdn_id
  FROM server
  WHERE host_name = $1
)
SELECT s.host_name AS host, s.domain_name AS domain, s.tcp_port AS port
FROM server s
inner join status AS st on st.id = s.status
WHERE cachegroup = (
  SELECT parent_cachegroup_id FROM cachegroup
  WHERE id = (SELECT cachegroup FROM child)
)
and st.name IN ('ONLINE', 'REPORTED')
and s.cdn_id = (SELECT cdn_id FROM child);`

const CacheConfigIPAllowQuery = `
SELECT value
FROM parameter
WHERE id IN (
  SELECT parameter
  FROM profile_parameter
  WHERE profile = (
    SELECT id
    FROM profile
    WHERE id = (
      SELECT profile
      FROM server
      WHERE host_name = $1
    )
  )
)
AND (name = 'allow_ip' or name = 'allow_ip6')
`

const CacheConfigDeliveryServiceRegexQuery = `
SELECT d.xml_id, r.pattern
FROM deliveryservice_regex AS dr
INNER JOIN regex AS r on r.id = dr.regex
INNER JOIN deliveryservice AS d on d.id = dr.deliveryservice
WHERE dr.deliveryservice IN (
  SELECT deliveryservice
  FROM deliveryservice_server
  WHERE server = (
    SELECT id
    FROM server
    WHERE host_name = $1
  )
)
`

func getCacheConfigJSON(cache string, db *sqlx.DB) (*toclient.CacheConfigResponse, error) {
	cdn, domain, err := getCacheConfigCDNDomain(cache, db)
	if err != nil {
		return nil, fmt.Errorf("getting domain: %v", err)
	}
	deliveryServices, err := getCacheConfigDeliveryServices(cache, db)
	if err != nil {
		return nil, fmt.Errorf("getting delivery services: %v", err)
	}
	parents, err := getCacheConfigParents(cache, db)
	if err != nil {
		return nil, fmt.Errorf("getting parents: %v", err)
	}
	ipallow, err := getCacheConfigIPAllow(cache, db)
	if err != nil {
		return nil, fmt.Errorf("getting IP allow parameters: %v", err)
	}

	deliveryServiceRegexes, err := getCacheConfigDeliveryServiceRegex(cache, db)
	if err != nil {
		return nil, fmt.Errorf("getting delivery service regexes: %v", err)
	}
	for i, ds := range deliveryServices {
		ds.Regexes = deliveryServiceRegexes[ds.XMLID]
		deliveryServices[i] = ds
	}

	return &toclient.CacheConfigResponse{
		Response: toclient.CacheConfig{
			CDN:              cdn,
			Domain:           domain,
			DeliveryServices: deliveryServices,
			Parents:          parents,
			AllowIP:          ipallow,
		},
	}, nil
}

func getCacheConfigCDNDomain(cache string, db *sqlx.DB) (string, string, error) {
	cdn := ""
	domain := ""
	if err := db.QueryRow(CacheConfigCDNDomainQuery, cache).Scan(&cdn, &domain); err != nil {
		return "", "", fmt.Errorf("querying %s domain: %v", cache, err)
	}
	return cdn, domain, nil
}

func getCacheConfigDeliveryServices(cache string, db *sqlx.DB) ([]toclient.CacheConfigDeliveryService, error) {
	dses := []toclient.CacheConfigDeliveryService{}
	rows, err := db.Query(CacheConfigDeliveryServicesQuery, cache)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		ds := toclient.CacheConfigDeliveryService{}
		if err := rows.Scan(&ds.Protocol, &ds.QueryStringIgnore, &ds.XMLID, &ds.Type, &ds.OriginFQDN, &ds.DSCP); err != nil {
			return nil, fmt.Errorf("row error: %v", err)
		}
		dses = append(dses, ds)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}
	return dses, nil
}

func getCacheConfigDeliveryServiceRegex(cache string, db *sqlx.DB) (map[string][]string, error) {
	dsr := map[string][]string{}
	rows, err := db.Query(CacheConfigDeliveryServiceRegexQuery, cache)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		dsid := ""
		rgx := ""
		if err := rows.Scan(&dsid, &rgx); err != nil {
			return nil, fmt.Errorf("row error: %v", err)
		}
		dsr[dsid] = append(dsr[dsid], rgx)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}
	return dsr, nil
}

func getCacheConfigParents(cache string, db *sqlx.DB) ([]toclient.CacheConfigParent, error) {
	parents := []toclient.CacheConfigParent{}
	rows, err := db.Query(CacheConfigParentsQuery, cache)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		p := toclient.CacheConfigParent{}
		if err := rows.Scan(&p.Host, &p.Domain, &p.Port); err != nil {
			return nil, fmt.Errorf("row error: %v", err)
		}
		parents = append(parents, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}
	return parents, nil
}

func getCacheConfigIPAllow(cache string, db *sqlx.DB) ([]string, error) {
	ips := []string{}
	rows, err := db.Query(CacheConfigIPAllowQuery, cache)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		ipsStr := ""
		if err := rows.Scan(&ipsStr); err != nil {
			return nil, fmt.Errorf("row error: %v", err)
		}
		ipsArr := strings.Split(ipsStr, ",")
		for i := 0; i < len(ipsArr); i++ {
			ipsArr[i] = strings.TrimSpace(ipsArr[i])
		}
		ips = append(ips, ipsArr...)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}
	return ips, nil
}
