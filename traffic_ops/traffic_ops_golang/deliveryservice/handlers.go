package deliveryservice

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
	"net/url"

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	tcapi "github.com/apache/incubator-trafficcontrol/lib/go-tc/v13"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/jmoiron/sqlx"
)

const DeliveryServicsPrivLevel = 10

func Handler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)

		ctx := r.Context()
		pathParams, err := api.GetPathParams(ctx)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		q := r.URL.Query()
		for k, v := range pathParams {
			q.Set(k, v)
		}

		resp, err := getDeliveryServicesResponse(q, db)

		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		respBts, err := json.Marshal(resp)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", respBts)
	}
}

func getDeliveryServicesResponse(q url.Values, db *sqlx.DB) (*tcapi.DeliveryServicesResponse, error) {
	dses, err := getDeliveryServices(q, db)
	if err != nil {
		return nil, fmt.Errorf("getting DeliveryServices response: %v", err)
	}

	resp := tcapi.DeliveryServicesResponse{
		Response: dses,
	}
	return &resp, nil
}

func getDeliveryServices(v url.Values, db *sqlx.DB) ([]tcapi.DeliveryService, error) {
	var rows *sqlx.Rows
	var err error

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToQueryCols := map[string]string{
		"xmlId": "xml_id",
	}

	query, queryValues := dbhelpers.BuildQuery(v, selectDSesQuery(), queryParamsToQueryCols)

	rows, err = db.NamedQuery(query, queryValues)
	fmt.Printf("rows ---> %v\n", rows)
	fmt.Printf("err ---> %v\n", err)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dses := []tcapi.DeliveryService{}
	for rows.Next() {
		var s tcapi.DeliveryService
		if err = rows.StructScan(&s); err != nil {
			return nil, fmt.Errorf("getting Delivery Services: %v", err)
		}
		dses = append(dses, s)
	}
	return dses, nil
}

func selectDSesQuery() string {
	query := `SELECT
 active,
 ccr_dns_ttl,
 cdn_id,
 cacheurl,
 check_path,
 dns_bypass_cname,
 dns_bypass_ip,
 dns_bypass_ip6,
 dns_bypass_ttl,
 dscp,
 display_name,
 edge_header_rewrite,
 geo_limit,
 geo_limit_countries,
 geolimit_redirect_url,
 geo_provider,
 global_max_mbps,
 global_max_tps,
 http_bypass_fqdn,
 id,
 ipv6_routing_enabled,
 info_url,
 initial_dispersion,
 last_updated,
 logs_enabled,
 long_desc,
 long_desc_1,
 long_desc_2,
 max_dns_answers,
 mid_header_rewrite,
 miss_lat,
 miss_long,
 multi_site_origin,
 multi_site_origin_algorithm,
 org_server_fqdn,
 origin_shield,
 profile,
 protocol,
 qstring_ignore,
 range_request_handling,
 regex_remap,
 regional_geo_blocking,
 remap_text,
 routing_name,
 ssl_key_version,
 signing_algorithm,
 tr_request_headers,
 tr_response_headers,
 tenant_id,
 type,
 xml_id

FROM deliveryservice d`
	return query
}
