package deliveryservice

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
	"context"
	"database/sql"
	"errors"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/riaksvc"
	"net/http"
	"strconv"
	"time"
)

type DsKey struct {
	XmlId   string
	Version sql.NullInt64
}

type DsExpirationInfo struct {
	XmlId      string
	Version    util.JSONIntStr
	Expiration time.Time
	AuthType   string
}

type ExpirationSummary struct {
	LetsEncryptExpirations []DsExpirationInfo
	SelfSignedExpirations  []DsExpirationInfo
	OtherExpirations       []DsExpirationInfo
}

func RenewCertificates(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	if inf.Config.RiakEnabled == false {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusServiceUnavailable, errors.New("the Riak service is unavailable"), errors.New("getting SSL keys from Riak by xml id: Riak is not configured"))
		return
	}

	rows, err := inf.Tx.Tx.Query(`SELECT xml_id, ssl_key_version FROM deliveryservice`)
	if err != nil {
		log.Errorf("querying: %v", err)
		return
	}
	defer rows.Close()

	dses := []DsKey{}
	for rows.Next() {
		ds := DsKey{}
		err := rows.Scan(&ds.XmlId, &ds.Version)
		if err != nil {
			log.Errorf("getting delivery services: %v", err)
			continue
		}
		if ds.Version.Valid && int(ds.Version.Int64) != 0 {
			dses = append(dses, ds)
		}
	}

	keysFound := ExpirationSummary{}
	for _, ds := range dses {
		keyObj, ok, err := riaksvc.GetDeliveryServiceSSLKeysObj(ds.XmlId, strconv.Itoa(int(ds.Version.Int64)), inf.Tx.Tx, inf.Config.RiakAuthOptions, inf.Config.RiakPort)
		if err != nil {
			log.Errorf("getting ssl keys: " + err.Error())
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting ssl keys: "+err.Error()))
			return
		}
		if !ok {
			log.Errorf("no object found for the specified key")
			api.WriteRespAlertObj(w, r, tc.InfoLevel, "no object found for the specified key", struct{}{}) // empty response object because Perl
			return
		}
		newVersion := util.JSONIntStr(keyObj.Version.ToInt64() + 1)

		dsExpInfo := DsExpirationInfo{
			XmlId:      keyObj.DeliveryService,
			Version:    newVersion,
			Expiration: keyObj.Expiration,
			AuthType:   keyObj.AuthType,
		}
		if keyObj.AuthType == tc.LetsEncryptAuthType {
			keysFound.LetsEncryptExpirations = append(keysFound.LetsEncryptExpirations, dsExpInfo)
			req := tc.DeliveryServiceLetsEncryptSSLKeysReq{
				DeliveryServiceSSLKeysReq: tc.DeliveryServiceSSLKeysReq{
					HostName:        &keyObj.Hostname,
					DeliveryService: &keyObj.DeliveryService,
					CDN:             &keyObj.CDN,
					Version:         &newVersion,
				},
			}
			ctx, _ := context.WithTimeout(r.Context(), time.Minute*10)

			if error := GetLetsEncryptCertificates(inf, req, ctx); error != nil {
				api.HandleErr(w, r, nil, http.StatusInternalServerError, error, nil)
				return
			}

		} else if keyObj.AuthType == tc.SelfSignedCertAuthType {
			keysFound.SelfSignedExpirations = append(keysFound.SelfSignedExpirations, dsExpInfo)
		} else {
			keysFound.OtherExpirations = append(keysFound.OtherExpirations, dsExpInfo)
		}

	}

	api.WriteResp(w, r, keysFound)

}
