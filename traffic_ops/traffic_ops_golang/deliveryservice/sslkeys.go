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
	"database/sql"
	"errors"
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/riaksvc"
)

func GenerateSSLKeys(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	req := tc.DeliveryServiceSSLKeysReq{}
	if err := api.Parse(r.Body, inf.Tx.Tx, &req); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("parsing request: "+err.Error()), nil)
		return
	}

	if err := generatePutRiakKeys(req, inf.Tx.Tx, inf.Config); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("generating and putting SSL keys: "+err.Error()))
		return
	}
	api.WriteResp(w, r, "Successfully created ssl keys for "+*req.DeliveryService)
}

// generatePutRiakKeys generates a certificate, csr, and key from the given request, and insert it into the Riak key database.
// The req MUST be validated, ensuring required fields exist.
func generatePutRiakKeys(req tc.DeliveryServiceSSLKeysReq, tx *sql.Tx, cfg *config.Config) error {
	dsSSLKeys := tc.DeliveryServiceSSLKeys{
		CDN:             *req.CDN,
		DeliveryService: *req.DeliveryService,
		BusinessUnit:    *req.BusinessUnit,
		City:            *req.City,
		Organization:    *req.Organization,
		Hostname:        *req.HostName,
		Country:         *req.Country,
		State:           *req.State,
		Key:             *req.Key,
		Version:         *req.Version,
	}
	if req.Certificate != nil {
		dsSSLKeys.Certificate = *req.Certificate
	} else {
		csr, crt, key, err := GenerateCert(*req.HostName, *req.Country, *req.City, *req.State, *req.Organization, *req.BusinessUnit)
		if err != nil {
			return errors.New("generating certificate: " + err.Error())
		}
		dsSSLKeys.Certificate = tc.DeliveryServiceSSLKeysCertificate{Crt: string(crt), Key: string(key), CSR: string(csr)}
	}
	if err := riaksvc.PutDeliveryServiceSSLKeysObjTx(dsSSLKeys, tx, cfg.RiakAuthOptions); err != nil {
		return errors.New("putting riak keys: " + err.Error())
	}

	dsSSLKeys.Version = riaksvc.DSSSLKeyVersionLatest
	if err := riaksvc.PutDeliveryServiceSSLKeysObjTx(dsSSLKeys, tx, cfg.RiakAuthOptions); err != nil {
		return errors.New("putting latest riak keys: " + err.Error())
	}
	return nil
}
