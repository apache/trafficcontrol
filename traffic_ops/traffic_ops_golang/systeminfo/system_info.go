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
	"errors"
	"net/http"
	"time"

	tc "github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/auth"

	"github.com/jmoiron/sqlx"
)

func Get(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	canReadSecureValue := false
	if (inf.Version.GreaterThanOrEqualTo(&api.Version{Major: 4}) && inf.Config.RoleBasedPermissions) || inf.Version.GreaterThanOrEqualTo(&api.Version{Major: 5}) {
		canReadSecureValue = inf.User.Can(tc.PermParameterSecureRead)
	} else {
		canReadSecureValue = inf.User.PrivLevel == auth.PrivLevelAdmin
	}
	api.RespWriter(w, r, inf.Tx.Tx)(getSystemInfo(inf.Tx, inf.User.PrivLevel, time.Duration(inf.Config.DBQueryTimeoutSeconds)*time.Second, canReadSecureValue))
}

func getSystemInfo(tx *sqlx.Tx, privLevel int, timeout time.Duration, canReadHiddenValue bool) (*tc.SystemInfo, error) {
	q := `
SELECT
  p.name,
  p.secure,
  p.last_updated,
  p.value
FROM
  parameter p
WHERE
  p.config_file = $1
`
	rows, err := tx.Queryx(q, tc.GlobalConfigFileName)
	if err != nil {
		return nil, errors.New("querying system info global parameters: " + err.Error())
	}
	defer rows.Close()
	info := map[string]string{}
	for rows.Next() {
		p := tc.ParameterNullable{}
		if err = rows.StructScan(&p); err != nil {
			return nil, errors.New("sqlx scanning system info global parameters: " + err.Error())
		}
		if p.Secure != nil && *p.Secure && !canReadHiddenValue {
			continue
		}
		if p.Name != nil && p.Value != nil {
			info[*p.Name] = *p.Value
		}
	}
	return &tc.SystemInfo{ParametersNullable: info}, nil
}
