package consistenthash

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
	"net/http"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"strconv"
	"io/ioutil"
)

func Get(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"regex", "requestpath", "cdnid"}, []string{"cdnid"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	regex := inf.Params["regex"]
	requestPath := inf.Params["requestpath"]
	cdnId := int64(inf.IntParams["cdnid"])

	resultPath, err := getPatternBasedConsistentHash(inf.Tx.Tx, regex, requestPath, cdnId)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting pattern based consistent hash from Traffic Router: " + err.Error()))
		return
	}

	api.WriteResp(w, r, resultPath)
}

func getPatternBasedConsistentHash(tx *sql.Tx, regex string, requestPath string, cdnId int64) ([]byte, error){
	q := `
SELECT concat(server.host_name, '.', server.domain_name) AS fqdn,
   parameter.value AS apiport
  FROM (((server
    JOIN profile ON ((profile.id = server.profile)))
    JOIN profile_parameter ON ((profile_parameter.profile = profile.id)))
    JOIN parameter ON ((parameter.id = profile_parameter.parameter)))
    JOIN status ON ((server.status = status.id))
 WHERE ((server.type = (select id from type where name = 'CCR')) AND
 (parameter.name = 'api.port'::text) AND
 (status.name = 'ONLINE') AND
 (server.cdn_id = $1))
`
	rows, err := tx.Query(q, cdnId)
	defer rows.Close()
	if err != nil {
		return nil, errors.New("querying for eligible Traffic Router to test pattern based consistent hashing: " + err.Error())
	}

	trafficRouter := ""
	apiPort := ""
	for rows.Next() {
		if err := rows.Scan(&trafficRouter, &apiPort); err != nil {
			return nil, errors.New("scanning eligible Traffic Routers for pattern based consistent hashing: " + err.Error())
		}
	}

	if trafficRouter == "" {
		return nil, errors.New("no eligible Traffic Router found for pattern based consistent hashing with cdn Id: " + strconv.FormatInt(cdnId, 10))
	}
	if apiPort == "" {
		return nil, errors.New("no parameter 'api.port' found for pattern based consistent hashing with cdn Id: " + strconv.FormatInt(cdnId, 10))
	}

	trafficRouterAPI := "http://" + trafficRouter + ":" + apiPort + "/crs/consistenthash/patternbased/regex?regex=" + regex + "&requestPath=" + requestPath

	r, err := http.Get(trafficRouterAPI)
	if err != nil {
		return nil, errors.New("Error creating request to Traffic Router: " + err.Error())
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)

	return body, nil
}
