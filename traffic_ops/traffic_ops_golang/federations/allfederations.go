package federations

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
	"github.com/apache/trafficcontrol/lib/go-tc/enum"
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
)

func GetAll(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	feds := []FedInfo{}
	err := error(nil)
	allFederations := []tc.IAllFederation{}

	if cdnParam, ok := inf.Params["cdnName"]; ok {
		cdnName := enum.CDNName(cdnParam)
		feds, err = getAllFederationsForCDN(inf.Tx.Tx, cdnName)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("federations.GetAll getting all federations: "+err.Error()))
			return
		}
		allFederations = append(allFederations, tc.AllFederationCDN{CDNName: &cdnName})
	} else {
		feds, err = getAllFederations(inf.Tx.Tx)
		if err != nil {
			api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("federations.GetAll getting all federations by CDN: "+err.Error()))
			return
		}
	}

	fedsResolvers, err := getFederationResolvers(inf.Tx.Tx, fedInfoIDs(feds))
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("federations.Get getting federations resolvers: "+err.Error()))
		return
	}
	allFederations = addResolvers(allFederations, feds, fedsResolvers)

	api.WriteResp(w, r, allFederations)
}

func getAllFederations(tx *sql.Tx) ([]FedInfo, error) {
	qry := `
SELECT
  fds.federation,
  fd.ttl,
  fd.cname,
  ds.xml_id
FROM
  federation_deliveryservice fds
  JOIN deliveryservice ds ON ds.id = fds.deliveryservice
  JOIN federation fd ON fd.id = fds.federation
ORDER BY
  ds.xml_id
`
	rows, err := tx.Query(qry)
	if err != nil {
		return nil, errors.New("all federations querying: " + err.Error())
	}
	defer rows.Close()

	feds := []FedInfo{}
	for rows.Next() {
		f := FedInfo{}
		if err := rows.Scan(&f.ID, &f.TTL, &f.CName, &f.DS); err != nil {
			return nil, errors.New("all federations scanning: " + err.Error())
		}
		feds = append(feds, f)
	}
	return feds, nil
}

func getAllFederationsForCDN(tx *sql.Tx, cdn enum.CDNName) ([]FedInfo, error) {
	qry := `
SELECT
  fds.federation,
  fd.ttl,
  fd.cname,
  ds.xml_id
FROM
  federation_deliveryservice fds
  JOIN deliveryservice ds ON ds.id = fds.deliveryservice
  JOIN federation fd ON fd.id = fds.federation
  JOIN cdn on cdn.id = ds.cdn_id
WHERE
  cdn.name = $1
ORDER BY
  ds.xml_id
`
	rows, err := tx.Query(qry, cdn)
	if err != nil {
		return nil, errors.New("all federations for cdn querying: " + err.Error())
	}
	defer rows.Close()

	feds := []FedInfo{}
	for rows.Next() {
		f := FedInfo{}
		if err := rows.Scan(&f.ID, &f.TTL, &f.CName, &f.DS); err != nil {
			return nil, errors.New("all federations for cdn scanning: " + err.Error())
		}
		feds = append(feds, f)
	}
	return feds, nil
}
