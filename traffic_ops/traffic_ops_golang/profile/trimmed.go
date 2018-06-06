package profile

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

	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/api"
)

func Trimmed(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		api.RespWriter(w, r)(getTrimmedProfiles(db))
	}
}

func getTrimmedProfiles(db *sql.DB) ([]tc.ProfileTrimmed, error) {
	rows, err := db.Query(`SELECT name FROM profile`)
	if err != nil {
		return nil, errors.New("querying trimmed profiles: " + err.Error())
	}
	defer rows.Close()
	profiles := []tc.ProfileTrimmed{}
	for rows.Next() {
		name := ""
		if err = rows.Scan(&name); err != nil {
			return nil, errors.New("scanning trimmed profiles: " + err.Error())
		}
		profiles = append(profiles, tc.ProfileTrimmed{Name: name})
	}
	return profiles, nil
}
