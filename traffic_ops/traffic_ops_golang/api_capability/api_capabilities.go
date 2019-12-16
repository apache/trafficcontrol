package api_capability

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
)

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

func GetAPICapabilitiesHandler(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	selectQuery := ` SELECT * FROM api_capability `

	capability, ok := inf.Params["capability"]
	if ok {
		selectQuery = fmt.Sprintf("%s WHERE capability = %s", selectQuery, capability)
	}

	order, ok := inf.Params["orderby"]
	if ok {
		selectQuery = fmt.Sprintf("%s ORDER BY %s", selectQuery, order)
	}

	rows, err := inf.Tx.Query(selectQuery)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New(fmt.Sprintf("db exception: could not query api_capbility with params: %v", inf.Params)), nil)
		return
	}

	fmt.Println("*** TEST AM I HERE ***")
	var apiCaps []tc.APICapability
	for rows.Next() {
		var ac tc.APICapability
		err = rows.Scan(
			&ac.ID,
			&ac.HTTPMethod,
			&ac.Route,
			&ac.Capability,
			&ac.LastUpdated,
		)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, nil, errors.New(fmt.Sprintf("api capability read: scanning: %s", err.Error())))
			return
		}
		apiCaps = append(apiCaps, ac)
	}

	api.WriteResp(w, r, tc.APICapabilityResponse{
		Response: apiCaps,
	})
	return
}
