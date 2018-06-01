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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/config"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"

	"github.com/jmoiron/sqlx"
)

type TODeliveryServiceV12 struct {
	tc.DeliveryServiceNullableV12
	Cfg config.Config
	DB  *sqlx.DB
}

func (ds TODeliveryServiceV12) MarshalJSON() ([]byte, error) {
	return json.Marshal(ds.DeliveryServiceNullableV12)
}

func (ds *TODeliveryServiceV12) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, ds.DeliveryServiceNullableV12)
}

func GetRefTypeV12(cfg config.Config, db *sqlx.DB) *TODeliveryServiceV12 {
	return &TODeliveryServiceV12{Cfg: cfg, DB: db}
}

func (ds TODeliveryServiceV12) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{"id", api.GetIntKey}}
}

func (ds TODeliveryServiceV12) GetKeys() (map[string]interface{}, bool) {
	if ds.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *ds.ID}, true
}

func (ds *TODeliveryServiceV12) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	ds.ID = &i
}

func (ds *TODeliveryServiceV12) GetAuditName() string {
	if ds.XMLID != nil {
		return *ds.XMLID
	}
	return ""
}

func (ds *TODeliveryServiceV12) GetType() string {
	return "ds"
}

// getDSTenantIDByID returns the tenant ID, whether the delivery service exists, and any error.
// Note the id may be nil, even if true is returned, if the delivery service exists but its tenant_id field is null.
func getDSTenantIDByID(tx *sql.Tx, id int) (*int, bool, error) {
	tenantID := (*int)(nil)
	if err := tx.QueryRow(`SELECT tenant_id FROM deliveryservice where id = $1`, id).Scan(&tenantID); err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("querying tenant ID for delivery service ID '%v': %v", id, err)
	}
	return tenantID, true, nil
}

// getDSTenantIDByName returns the tenant ID, whether the delivery service exists, and any error.
// Note the id may be nil, even if true is returned, if the delivery service exists but its tenant_id field is null.
func getDSTenantIDByName(tx *sql.Tx, name string) (*int, bool, error) {
	tenantID := (*int)(nil)
	if err := tx.QueryRow(`SELECT tenant_id FROM deliveryservice where xml_id = $1`, name).Scan(&tenantID); err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("querying tenant ID for delivery service name '%v': %v", name, err)
	}
	return tenantID, true, nil
}

// GetXMLID loads the DeliveryService's xml_id from the database, from the ID. Returns whether the delivery service was found, and any error.
func (ds *TODeliveryServiceV12) GetXMLID(tx *sql.Tx) (string, bool, error) {
	if ds.ID == nil {
		return "", false, errors.New("missing ID")
	}
	xmlID := ""
	if err := tx.QueryRow(`SELECT xml_id FROM deliveryservice where id = $1`, ds.ID).Scan(&xmlID); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, fmt.Errorf("querying xml_id for delivery service ID '%v': %v", *ds.ID, err)
	}
	return xmlID, true, nil
}

// IsTenantAuthorized checks that the user is authorized for both the delivery service's existing tenant, and the new tenant they're changing it to (if different).
func (ds *TODeliveryServiceV12) IsTenantAuthorized(user *auth.CurrentUser, db *sqlx.DB) (bool, error) {
	tx, err := db.DB.Begin() // must be last, MUST not return an error if this suceeds, without closing the tx
	if err != nil {
		return false, errors.New("beginning transaction: " + err.Error())
	}
	defer dbhelpers.FinishTx(tx, util.BoolPtr(true))
	return isTenantAuthorized(user, tx, &ds.DeliveryServiceNullableV12)
}

// getTenantID returns the tenant Id of the given delivery service. Note it may return a nil id and nil error, if the tenant ID in the database is nil.
func getTenantID(tx *sql.Tx, ds *tc.DeliveryServiceNullableV12) (*int, error) {
	if ds.ID == nil && ds.XMLID == nil {
		return nil, errors.New("delivery service has no ID or XMLID")
	}
	if ds.ID != nil {
		existingID, _, err := getDSTenantIDByID(tx, *ds.ID) // ignore exists return - if the DS is new, we only need to check the user input tenant
		return existingID, err
	}
	existingID, _, err := getDSTenantIDByName(tx, *ds.XMLID) // ignore exists return - if the DS is new, we only need to check the user input tenant
	return existingID, err
}

func isTenantAuthorized(user *auth.CurrentUser, tx *sql.Tx, ds *tc.DeliveryServiceNullableV12) (bool, error) {
	existingID, err := getTenantID(tx, ds)
	if err != nil {
		return false, errors.New("getting tenant ID: " + err.Error())
	}
	if ds.TenantID == nil {
		ds.TenantID = existingID
	}
	if existingID != nil && existingID != ds.TenantID {
		userAuthorizedForExistingDSTenant, err := tenant.IsResourceAuthorizedToUserTx(*existingID, user, tx)
		if err != nil {
			return false, errors.New("checking authorization for existing DS ID: " + err.Error())
		}
		if !userAuthorizedForExistingDSTenant {
			return false, nil
		}
	}
	if ds.TenantID != nil {
		userAuthorizedForNewDSTenant, err := tenant.IsResourceAuthorizedToUserTx(*ds.TenantID, user, tx)
		if err != nil {
			return false, errors.New("checking authorization for new DS ID: " + err.Error())
		}
		if !userAuthorizedForNewDSTenant {
			return false, nil
		}
	}
	return true, nil
}

func (ds *TODeliveryServiceV12) Validate(db *sqlx.DB) []error {
	tx, err := db.DB.Begin()
	if err != nil {
		return []error{errors.New("beginning transaction: " + err.Error())}
	}
	defer dbhelpers.FinishTx(tx, util.BoolPtr(true))
	return []error{ds.DeliveryServiceNullableV12.Validate(tx)}
}

func CreateV12(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	ds := tc.DeliveryServiceNullableV12{}
	if err := api.Parse(r.Body, inf.Tx.Tx, &ds); err != nil {
		api.HandleErr(w, r, http.StatusBadRequest, errors.New("decoding: "+err.Error()), nil)
		return
	}
	dsv13 := tc.NewDeliveryServiceNullableV13FromV12(ds)
	if authorized, err := isTenantAuthorized(inf.User, inf.Tx.Tx, &ds); err != nil {
		api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("checking tenant: "+err.Error()))
		return
	} else if !authorized {
		api.HandleErr(w, r, http.StatusForbidden, errors.New("not authorized on this tenant"), nil)
		return
	}
	dsv13, errCode, userErr, sysErr = create(inf.Tx.Tx, inf.Config, inf.User, dsv13)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, errCode, userErr, sysErr)
		return
	}
	*inf.CommitTx = true
	api.WriteResp(w, r, []tc.DeliveryServiceNullableV12{dsv13.DeliveryServiceNullableV12})
}

func (ds *TODeliveryServiceV12) Read(db *sqlx.DB, params map[string]string, user auth.CurrentUser) ([]interface{}, []error, tc.ApiErrorType) {
	returnable := []interface{}{}
	dses, errs, errType := readGetDeliveryServices(params, db, user)
	if len(errs) > 0 {
		for _, err := range errs {
			if err.Error() == `id cannot parse to integer` {
				return nil, []error{errors.New("Resource not found.")}, tc.DataMissingError //matches perl response
			}
		}
		return nil, errs, errType
	}

	for _, ds := range dses {
		returnable = append(returnable, ds.DeliveryServiceNullableV12)
	}
	return returnable, nil, tc.NoError
}

func (ds *TODeliveryServiceV12) Delete(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType) {
	v13 := &TODeliveryServiceV13{
		Cfg: ds.Cfg,
		DB:  ds.DB,
		DeliveryServiceNullableV13: tc.DeliveryServiceNullableV13{
			DeliveryServiceNullableV12: ds.DeliveryServiceNullableV12,
		},
	}
	err, errType := v13.Delete(db, user)
	ds.DeliveryServiceNullableV12 = v13.DeliveryServiceNullableV12 // TODO avoid copy
	return err, errType
}

func UpdateV12(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	ds := tc.DeliveryServiceNullableV12{}
	ds.ID = util.IntPtr(inf.IntParams["id"])
	if err := api.Parse(r.Body, inf.Tx.Tx, &ds); err != nil {
		api.HandleErr(w, r, http.StatusBadRequest, errors.New("decoding: "+err.Error()), nil)
		return
	}
	dsv13 := tc.NewDeliveryServiceNullableV13FromV12(ds)
	if authorized, err := isTenantAuthorized(inf.User, inf.Tx.Tx, &ds); err != nil {
		api.HandleErr(w, r, http.StatusInternalServerError, nil, errors.New("checking tenant: "+err.Error()))
		return
	} else if !authorized {
		api.HandleErr(w, r, http.StatusForbidden, errors.New("not authorized on this tenant"), nil)
		return
	}
	dsv13, errCode, userErr, sysErr = update(inf.Tx.Tx, inf.Config, inf.User, &dsv13)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, errCode, userErr, sysErr)
		return
	}
	api.WriteResp(w, r, []tc.DeliveryServiceNullableV12{dsv13.DeliveryServiceNullableV12})
}
