package systeminfo

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
	"tmp/prometheus/common/log"

	tc "github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"

	"github.com/jmoiron/sqlx"
)

//we need a type alias to define functions on
type TOParameter tc.ParameterNullable

//the refType is passed into the handlers where a copy of its type is used to decode the json.
var refType = TOParameter(tc.ParameterNullable{})

func GetRefType() *TOParameter {
	return &refType
}

func (pl *TOParameter) Read(db *sqlx.DB, parameters map[string]string, user auth.CurrentUser) ([]interface{}, []error, tc.ApiErrorType) {
	var rows *sqlx.Rows

	query := selectQuery()
	log.Debugln("Query is ", query)

	rows, err := db.Queryx(query)
	if err != nil {
		log.Errorf("Error querying Parameter: %v", err)
		return nil, []error{tc.DBError}, tc.SystemError
	}
	defer rows.Close()

	params := []interface{}{}
	for rows.Next() {
		var s tc.ParameterNullable
		if err = rows.StructScan(&s); err != nil {
			log.Errorf("error parsing Parameter rows: %v", err)
			return nil, []error{tc.DBError}, tc.SystemError
		}
		params = append(params, s)
	}

	return params, []error{}, tc.NoError

}

func selectQuery() string {

	// system info returns all global parameters
	query := `SELECT
p.config_file,
p.name,
p.secure,
p.value
FROM parameter p
WHERE p.config_file='global'`

	return query
}
