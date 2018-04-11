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
	"strings"

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"

	"github.com/lib/pq"
)

type WhereColumnInfo struct {
	Column  string
	Checker func(string) error
}

func BuildWhereAndOrderBy(parameters map[string]string, queryParamsToSQLCols map[string]WhereColumnInfo) (string, string, map[string]interface{}, []error) {
	var whereClause string
	var orderBy string
	var criteria string
	var queryValues map[string]interface{}
	var errs []error
	criteria, queryValues, errs = parseCriteriaAndQueryValues(queryParamsToSQLCols, parameters)

	if len(queryValues) > 0 {
		whereClause += " " + criteria
	}
	if len(errs) > 0 {
		return "", "", queryValues, errs
	}

	if names, ok := parameters["orderby"]; ok {
		var columns []string
		log.Debugln("orderby: ", names)
		for _, n := range strings.Split(names, ",") {
			if colInfo, ok := queryParamsToSQLCols[n]; ok {
				log.Debugln("orderby column ", colInfo)
				columns = append(columns, colInfo.Column)
			} else {
				log.Debugln("Incorrect name for orderby: ", n)
			}
		}
		orderBy = strings.Join(columns, ",")
	}
	if whereClause != "" {
		whereClause = "\nWHERE " + whereClause
	}
	if orderBy != "" {
		whereClause = "\nORDER BY " + whereClause
	}
	log.Debugf("\n--\n Where: %s \n Order By: %s", whereClause, orderBy)
	return whereClause, orderBy, queryValues, errs
}

func parseCriteriaAndQueryValues(queryParamsToSQLCols map[string]WhereColumnInfo, parameters map[string]string) (string, map[string]interface{}, []error) {
	m := make(map[string]interface{})
	var criteria string

	var criteriaArgs []string
	errs := []error{}
	queryValues := make(map[string]interface{})
	for key, colInfo := range queryParamsToSQLCols {
		if urlValue, ok := parameters[key]; ok {
			var err error
			if colInfo.Checker != nil {
				err = colInfo.Checker(urlValue)
			}
			if err != nil {
				errs = append(errs, errors.New(key+" "+err.Error()))
			} else {
				m[key] = urlValue
				criteria = colInfo.Column + "=:" + key
				criteriaArgs = append(criteriaArgs, criteria)
				queryValues[key] = urlValue
			}
		}
	}
	criteria = strings.Join(criteriaArgs, " AND ")

	return criteria, queryValues, errs
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
