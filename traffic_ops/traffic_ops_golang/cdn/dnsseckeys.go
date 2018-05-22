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
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/riaksvc"
)

func GetDNSSECKeys(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"name"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	api.RespWriter(w, r)(getCDNDNSSECKeys(inf.Params["name"], inf.Tx.Tx, inf.Config))
}

func getCDNDNSSECKeys(cdnName string, tx *sql.Tx, cfg *config.Config) (tc.DNSSECKeys, error) {
	keys, ok, err := riaksvc.GetDNSSECKeys(cdnName, tx, cfg.RiakAuthOptions)
	if err != nil {
		return tc.DNSSECKeys{}, errors.New("getting DNSSec keys from Riak: " + err.Error())
	}
	if !ok {
		return tc.DNSSECKeys{}, nil
	}
	return keys, nil
}
