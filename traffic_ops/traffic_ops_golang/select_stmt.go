package main

import (
	"net/url"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
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

func SelectStmt(v url.Values, selectStmt string, queryParamsToQueryCols map[string]string) (string, map[string]interface{}) {

	var sqlQuery string
	if len(v) > 0 {
		sqlQuery = selectStmt + "\nWHERE " + appendCriteria(queryParamsToQueryCols, v)
	} else {
		sqlQuery = selectStmt
	}
	log.Debugln("\n--\n" + sqlQuery)
	queryValues := queryValues(queryParamsToQueryCols, v)
	return sqlQuery, queryValues
}

func appendCriteria(queryParamsToQueryCols map[string]string, v url.Values) string {

	m := make(map[string]interface{})
	var criteria string
	for key, val := range queryParamsToQueryCols {
		if urlValue, ok := v[key]; ok {
			m[key] = urlValue[0]
			criteria = val + "=:" + key
			break
		}
	}
	return criteria
}

func queryValues(queryParamsToQueryCols map[string]string, v url.Values) map[string]interface{} {

	queryValues := make(map[string]interface{})
	for key, _ := range queryParamsToQueryCols {
		if urlValue, ok := v[key]; ok {
			queryValues[key] = urlValue[0]
			break
		}
	}
	return queryValues
}
