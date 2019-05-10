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
	"encoding/json"
	"errors"
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
)

//we need a type alias to define functions on

type TODeliveryServiceV13 struct {
	api.APIInfoImpl
	tc.DeliveryServiceNullableV13
}

func (ds TODeliveryServiceV13) MarshalJSON() ([]byte, error) {
	return json.Marshal(ds.DeliveryServiceNullableV13)
}

func (ds *TODeliveryServiceV13) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, ds.DeliveryServiceNullableV13)
}

func (ds *TODeliveryServiceV13) APIInfo() *api.APIInfo { return ds.ReqInfo }

func (ds *TODeliveryServiceV13) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	ds.ID = &i
}

func (ds *TODeliveryServiceV13) Validate() error {
	return ds.DeliveryServiceNullableV13.Validate(ds.APIInfo().Tx.Tx)
}

// 	TODO allow users to post names (type, cdn, etc) and get the IDs from the names. This isn't trivial to do in a single query, without dynamically building the entire insert query, and ideally inserting would be one query. But it'd be much more convenient for users. Alternatively, remove IDs from the database entirely and use real candidate keys.
func CreateV13(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	ds := tc.DeliveryServiceNullableV13{}
	if err := api.Parse(r.Body, inf.Tx.Tx, &ds); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("decoding: "+err.Error()), nil)
		return
	}

	if ds.RoutingName == nil || *ds.RoutingName == "" {
		ds.RoutingName = util.StrPtr("cdn")
	}
	if err := ds.Validate(inf.Tx.Tx); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("invalid request: "+err.Error()), nil)
		return
	}
	tcDS := tc.NewDeliveryServiceNullableFromV13(ds)
	tcDS, errCode, userErr, sysErr = create(inf, tcDS)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Deliveryservice creation was successful.", []tc.DeliveryServiceNullableV13{tcDS.DeliveryServiceNullableV13})
}

func (ds *TODeliveryServiceV13) Read() ([]interface{}, error, error, int) {
	returnable := []interface{}{}
	dses, errs, _ := readGetDeliveryServices(ds.APIInfo().Params, ds.APIInfo().Tx, ds.APIInfo().User)
	if len(errs) > 0 {
		for _, err := range errs {
			if err.Error() == `id cannot parse to integer` { // TODO create const for string
				return nil, errors.New("Resource not found."), nil, http.StatusNotFound //matches perl response
			}
		}
		return nil, nil, errors.New("reading dses: " + util.JoinErrsStr(errs)), http.StatusInternalServerError
	}

	for _, ds := range dses {
		returnable = append(returnable, ds.DeliveryServiceNullableV13)
	}
	return returnable, nil, nil, http.StatusOK
}

func UpdateV13(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	id := inf.IntParams["id"]

	ds := tc.DeliveryServiceNullable{}
	if err := json.NewDecoder(r.Body).Decode(&ds); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("malformed JSON: "+err.Error()), nil)
		return
	}
	ds.ID = &id

	if err := ds.Validate(inf.Tx.Tx); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("invalid request: "+err.Error()), nil)
		return
	}

	ds, errCode, userErr, sysErr = update(inf, &ds)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Deliveryservice update was successful.", []tc.DeliveryServiceNullable{ds})
}
