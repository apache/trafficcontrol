package dbhelpers

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
	"errors"
	"net/url"
	"strings"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/lib/pq"
)

func BuildQuery(v url.Values, selectStmt string, queryParamsToSQLCols map[string]string) (string, map[string]interface{}) {
	var sqlQuery string
	var criteria string
	var queryValues map[string]interface{}
	sqlQuery = selectStmt
	criteria, queryValues = parseCriteriaAndQueryValues(queryParamsToSQLCols, v)

	if len(queryValues) > 0 {
		sqlQuery += "\nWHERE " + criteria
	}

	if orderby, ok := v["orderby"]; ok {
		log.Debugln("orderby: ", orderby[0])
		if col, ok := queryParamsToSQLCols[orderby[0]]; ok {
			log.Debugln("orderby column ", col)
			sqlQuery += "\nORDER BY " + col
		} else {
			log.Debugln("Incorrect name for orderby: ", orderby[0])
		}
	}
	log.Debugln("\n--\n" + sqlQuery)
	return sqlQuery, queryValues
}

func parseCriteriaAndQueryValues(queryParamsToSQLCols map[string]string, v url.Values) (string, map[string]interface{}) {
	m := make(map[string]interface{})
	var criteria string

	var criteriaArgs []string
	queryValues := make(map[string]interface{})
	for key, val := range queryParamsToSQLCols {
		if urlValue, ok := v[key]; ok {
			m[key] = urlValue[0]
			criteria = val + "=:" + key
			criteriaArgs = append(criteriaArgs, criteria)
			queryValues[key] = urlValue[0]
		}
	}
	criteria = strings.Join(criteriaArgs, " AND ")

	return criteria, queryValues
}

//parses pq errors for uniqueness constraint violations
func ParsePQUniqueConstraintError(err *pq.Error) (error, tc.ApiErrorType) {
	if len(err.Constraint) > 0 && len(err.Detail) > 0 { //we only want to continue parsing if it is a constraint error with details
		detail := err.Detail
		if strings.HasPrefix(detail, "Key ") && strings.HasSuffix(detail, " already exists.") { //we only want to continue parsing if it is a uniqueness constraint error
			detail = strings.TrimPrefix(detail, "Key ")
			detail = strings.TrimSuffix(detail, " already exists.")
			//should look like "(column)=(dupe value)" at this point
			details := strings.Split(detail, "=")
			if len(details) == 2 {
				column := strings.Trim(details[0], "()")
				dupValue := strings.Trim(details[1], "()")
				return errors.New(column + " " + dupValue + " already exists."), tc.DataConflictError
			}
		}
	}
	log.Error.Printf("failed to parse unique constraint from pq error: %v", err)
	return tc.DBError, tc.SystemError
}
