package cdn

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

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
)

func GetConfigs(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	api.RespWriter(w, r, inf.Tx.Tx)(getConfigs(inf.Tx.Tx))
}

func getConfigs(tx *sql.Tx) ([]tc.CDNConfig, error) {
	rows, err := tx.Query(`SELECT name, id FROM cdn`)
	if err != nil {
		return nil, errors.New("querying cdn configs: " + err.Error())
	}
	cdns := []tc.CDNConfig{}
	defer rows.Close()
	for rows.Next() {
		c := tc.CDNConfig{}
		if err := rows.Scan(&c.Name, &c.ID); err != nil {
			return nil, errors.New("scanning cdn config: " + err.Error())
		}
		cdns = append(cdns, c)
	}
	return cdns, nil
}
