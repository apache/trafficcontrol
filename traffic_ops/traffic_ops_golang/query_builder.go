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
	"net/url"

	"github.com/apache/incubator-trafficcontrol/traffic_monitor_golang/common/log"
)

func BuildQuery(v url.Values, selectStmt string, queryParamsToSQLCols map[string]string) (string, map[string]interface{}) {
	var sqlQuery string
	var criteria string
	var queryValues map[string]interface{}
	if len(v) > 0 {
		criteria, queryValues = parseCriteriaAndQueryValues(queryParamsToSQLCols, v)
		if len(queryValues) > 0 {
			sqlQuery = selectStmt + "\nWHERE " + criteria
		} else {
			sqlQuery = selectStmt
		}
	} else {
		sqlQuery = selectStmt
	}
	log.Debugln("\n--\n" + sqlQuery)
	return sqlQuery, queryValues
}

func parseCriteriaAndQueryValues(queryParamsToSQLCols map[string]string, v url.Values) (string, map[string]interface{}) {
	m := make(map[string]interface{})
	var criteria string
	queryValues := make(map[string]interface{})
	for key, val := range queryParamsToSQLCols {
		if urlValue, ok := v[key]; ok {
			m[key] = urlValue[0]
			criteria = val + "=:" + key
			queryValues[key] = urlValue[0]
			break
		}
	}
	return criteria, queryValues
}
