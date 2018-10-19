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
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"

	"github.com/lib/pq"
)

func Get(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	feds, err := getUserFederations(inf.Tx.Tx, inf.User.UserName)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("federations.Get getting federations: "+err.Error()))
		return
	}
	fedsResolvers, err := getFederationResolvers(inf.Tx.Tx, fedInfoIDs(feds))
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("federations.Get getting federations resolvers: "+err.Error()))
		return
	}
	allFederations := addResolvers([]tc.IAllFederation{}, feds, fedsResolvers)
	api.WriteResp(w, r, allFederations)
}

func addResolvers(allFederations []tc.IAllFederation, feds []FedInfo, fedsResolvers map[int][]FedResolverInfo) []tc.IAllFederation {
	dsFeds := map[tc.DeliveryServiceName][]tc.AllFederationMapping{}
	for _, fed := range feds {
		mapping := tc.AllFederationMapping{}
		mapping.TTL = util.IntPtr(fed.TTL)
		mapping.CName = util.StrPtr(fed.CName)
		for _, resolver := range fedsResolvers[fed.ID] {
			switch resolver.Type {
			case tc.FederationResolverType4:
				mapping.Resolve4 = append(mapping.Resolve4, resolver.IP)
			case tc.FederationResolverType6:
				mapping.Resolve6 = append(mapping.Resolve6, resolver.IP)
			default:
				log.Warnf("federations addResolvers got invalid resolver type for federation '%v', skipping\n", fed.ID)
			}
		}
		dsFeds[fed.DS] = append(dsFeds[fed.DS], mapping)
	}

	for ds, mappings := range dsFeds {
		allFederations = append(allFederations, tc.AllFederation{DeliveryService: ds, Mappings: mappings})
	}
	return allFederations
}

func fedInfoIDs(feds []FedInfo) []int {
	ids := []int{}
	for _, fed := range feds {
		ids = append(ids, fed.ID)
	}
	return ids
}

type FedInfo struct {
	ID    int
	TTL   int
	CName string
	DS    tc.DeliveryServiceName
}

type FedResolverInfo struct {
	Type tc.FederationResolverType
	IP   string
}

// getFederationResolvers takes a slice of federation IDs, and returns a map[federationID]info.
func getFederationResolvers(tx *sql.Tx, fedIDs []int) (map[int][]FedResolverInfo, error) {
	qry := `
SELECT
  ffr.federation,
  frt.name as resolver_type,
  fr.ip_address
FROM
  federation_federation_resolver ffr
  JOIN federation_resolver fr ON ffr.federation_resolver = fr.id
  JOIN type frt on fr.type = frt.id
WHERE
  ffr.federation = ANY($1)
`
	rows, err := tx.Query(qry, pq.Array(fedIDs))
	if err != nil {
		return nil, errors.New("all federations resolvers querying: " + err.Error())
	}
	defer rows.Close()

	feds := map[int][]FedResolverInfo{}
	for rows.Next() {
		fedID := 0
		f := FedResolverInfo{}
		fType := ""
		if err := rows.Scan(&fedID, &fType, &f.IP); err != nil {
			return nil, errors.New("all federations resolvers scanning: " + err.Error())
		}
		f.Type = tc.FederationResolverTypeFromString(fType)
		feds[fedID] = append(feds[fedID], f)
	}
	return feds, nil
}

func getUserFederations(tx *sql.Tx, userName string) ([]FedInfo, error) {
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
  JOIN federation_tmuser fu on fu.federation = fd.id
  JOIN tm_user u on u.id = fu.tm_user
WHERE
  u.username = $1
ORDER BY
  ds.xml_id
`
	rows, err := tx.Query(qry, userName)
	if err != nil {
		return nil, errors.New("user federations querying: " + err.Error())
	}
	defer rows.Close()

	feds := []FedInfo{}
	for rows.Next() {
		f := FedInfo{}
		if err := rows.Scan(&f.ID, &f.TTL, &f.CName, &f.DS); err != nil {
			return nil, errors.New("user federations scanning: " + err.Error())
		}
		feds = append(feds, f)
	}
	return feds, nil
}
