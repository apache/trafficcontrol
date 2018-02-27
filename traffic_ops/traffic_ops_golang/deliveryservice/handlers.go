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

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/jmoiron/sqlx"
)

const DeliveryServicsPrivLevel = 10

func Handler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleErrs := tc.GetHandleErrorsFunc(w, r)

		// Load the the query and path params with path params overriding query params
		params, err := api.GetCombinedParams(r)
		if err != nil {
			log.Errorf("unable to get parameters from request: %s", err)
			handleErrs(http.StatusInternalServerError, err)
		}

		resp, errs, errType := getDeliveryServicesResponse(params, db)
		tc.HandleErrorsWithType(errs, errType, handleErrs)

		respBts, err := json.Marshal(resp)
		if err != nil {
			handleErrs(http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", respBts)
	}
}

func getDeliveryServicesResponse(parameters map[string]string, db *sqlx.DB) (*tc.DeliveryServicesResponse, []error, tc.ApiErrorType) {
	dses, errs, errType := getDeliveryServices(parameters, db)
	if len(errs) > 0 {
		return nil, errs, errType
	}

	resp := tc.DeliveryServicesResponse{
		Response: dses,
	}
	return &resp, nil, tc.NoError
}

func getDeliveryServices(parameters map[string]string, db *sqlx.DB) ([]tc.DeliveryService, []error, tc.ApiErrorType) {
	var rows *sqlx.Rows
	var err error

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToQueryCols := map[string]dbhelpers.WhereColumnInfo{
		"xmlId": dbhelpers.WhereColumnInfo{"xml_id", nil},
	}

	where, orderBy, queryValues, errs := dbhelpers.BuildWhereAndOrderBy(parameters, queryParamsToQueryCols)
	if len(errs) > 0 {
		return nil, errs, tc.DataConflictError
	}
	query := selectDSesQuery() + where + orderBy

	rows, err = db.NamedQuery(query, queryValues)
	fmt.Printf("rows ---> %v\n", rows)
	fmt.Printf("err ---> %v\n", err)
	if err != nil {
		return nil, []error{err}, tc.SystemError
	}
	defer rows.Close()

	dses := []tc.DeliveryService{}
	for rows.Next() {
		var s tc.DeliveryService
		if err = rows.StructScan(&s); err != nil {
			return nil, []error{fmt.Errorf("getting Delivery Services: %v", err)}, tc.SystemError
		}
		dses = append(dses, s)
	}
	return dses, nil, tc.NoError
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
