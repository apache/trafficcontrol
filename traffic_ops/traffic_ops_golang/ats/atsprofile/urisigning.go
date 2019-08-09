package atsprofile

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
	"strings"

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/ats"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/riaksvc"
)

func GetURISigning(w http.ResponseWriter, r *http.Request) {
	WithProfileData(w, r, tc.ApplicationJson, uriSigningDotConfig)
}

func uriSigningDotConfig(tx *sql.Tx, cfg *config.Config, _ ats.ProfileData, fileName string) (string, error) {
	riakKey := strings.TrimSuffix(strings.TrimPrefix(fileName, "uri_signing_"), ".config")
	keys, hasKeys, err := riaksvc.GetURISigningKeysRaw(tx, cfg.RiakAuthOptions, cfg.RiakPort, riakKey)
	if err != nil {
		return "", errors.New("getting uri signing keys from Riak: " + err.Error())
	}
	if !hasKeys {
		keys = []byte{} // TODO verify? Perl seems to return without returning its $text
	}
	return atscfg.MakeURISigningConfig(keys), nil
}
