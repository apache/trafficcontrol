package request

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
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"

	"github.com/jmoiron/sqlx"
)

// DeliveryServiceRequestPrivLevel ...
const DeliveryServiceRequestPrivLevel = 20

// Handler returns a func to handle GET requests
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

		resp, err := getDeliveryServiceRequestsResponse(q, db)

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

func getDeliveryServiceRequestsResponse(q url.Values, db *sqlx.DB) (*tc.DeliveryServiceRequestsResponse, error) {
	deliveryServiceRequests, err := getDeliveryServiceRequests(q, db)
	if err != nil {
		return nil, fmt.Errorf("getting DeliveryServiceRequests response: %v", err)
	}

	resp := tc.DeliveryServiceRequestsResponse{
		Response: deliveryServiceRequests,
	}
	return &resp, nil
}

func getDeliveryServiceRequests(v url.Values, db *sqlx.DB) ([]tc.DeliveryServiceRequest, error) {
	var rows *sqlx.Rows
	var err error

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToQueryCols := map[string]string{
		"assignee":   "s.username",
		"assigneeId": "r.assignee_id",
		"author":     "a.username",
		"authorId":   "r.author_id",
		"changeType": "r.change_type",
		"id":         "r.id",
		"status":     "r.status",
	}

	query, queryValues := dbhelpers.BuildQuery(v, selectDeliveryServiceRequestsQuery(), queryParamsToQueryCols)

	rows, err = db.NamedQuery(query, queryValues)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deliveryServiceRequests := []tc.DeliveryServiceRequest{}
	for rows.Next() {
		var s tc.DeliveryServiceRequest
		if err = rows.StructScan(&s); err != nil {
			return nil, fmt.Errorf("getting DeliveryServiceRequests: %v", err)
		}
		deliveryServiceRequests = append(deliveryServiceRequests, s)
	}
	return deliveryServiceRequests, nil
}

func selectDeliveryServiceRequestsQuery() string {

	query := `SELECT
a.username AS author,
r.assignee_id,
r.author_id,
r.change_type,
r.created_at,
r.id,
r.last_updated,
r.request,
r.status,
s.username AS assignee

FROM deliveryservice_request r
JOIN tm_user a ON r.author_id = a.id
LEFT OUTER JOIN tm_user s ON r.assignee_id = s.id`

	return query
}
