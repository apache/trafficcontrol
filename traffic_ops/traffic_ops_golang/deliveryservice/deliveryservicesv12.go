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
	"fmt"
	"net/http"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type TODeliveryServiceV12 struct {
	api.APIInfoImpl
	tc.DeliveryServiceNullableV12
}

func (v *TODeliveryServiceV12) DeleteQuery() string {
	return `DELETE FROM deliveryservice WHERE id = :id`
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

// GetDSTenantIDByIDTx returns the tenant ID, whether the delivery service exists, and any error.
func GetDSTenantIDByIDTx(tx *sql.Tx, id int) (*int, bool, error) {
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

// GetDSTenantIDByNameTx returns the tenant ID, whether the delivery service exists, and any error.
func GetDSTenantIDByNameTx(tx *sql.Tx, ds tc.DeliveryServiceName) (*int, bool, error) {
	tenantID := (*int)(nil)
	if err := tx.QueryRow(`SELECT tenant_id FROM deliveryservice where xml_id = $1`, ds).Scan(&tenantID); err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("querying tenant ID for delivery service name '%v': %v", ds, err)
	}
	return tenantID, true, nil
}

// GetXMLID loads the DeliveryService's xml_id from the database, from the ID. Returns whether the delivery service was found, and any error.

func (ds *TODeliveryServiceV12) GetXMLID(tx *sqlx.Tx) (string, bool, error) {
	if ds.ID == nil {
		return "", false, errors.New("missing ID")
	}
	return GetXMLID(tx.Tx, *ds.ID)
}

// GetXMLID loads the DeliveryService's xml_id from the database, from the ID. Returns whether the delivery service was found, and any error.
func GetXMLID(tx *sql.Tx, id int) (string, bool, error) {
	xmlID := ""
	if err := tx.QueryRow(`SELECT xml_id FROM deliveryservice where id = $1`, id).Scan(&xmlID); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, fmt.Errorf("querying xml_id for delivery service ID '%v': %v", id, err)
	}
	return xmlID, true, nil
}

// IsTenantAuthorized checks that the user is authorized for both the delivery service's existing tenant, and the new tenant they're changing it to (if different).

func (ds *TODeliveryServiceV12) IsTenantAuthorized(user *auth.CurrentUser) (bool, error) {
	return isTenantAuthorized(user, ds.ReqInfo.Tx.Tx, &ds.DeliveryServiceNullableV12)
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

func (ds *TODeliveryServiceV12) Validate() error {
	return ds.DeliveryServiceNullableV12.Validate(ds.ReqInfo.Tx.Tx)
}

// Create is unimplemented, needed to satisfy CRUDer, since the framework doesn't allow a create to return an array of one
func (ds *TODeliveryServiceV12) Create() (error, error, int) {
	return nil, nil, http.StatusNotImplemented
}

func CreateV12(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	ds := tc.DeliveryServiceNullableV12{}
	if err := api.Parse(r.Body, inf.Tx.Tx, &ds); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("decoding: "+err.Error()), nil)
		return
	}
	dsv13 := tc.NewDeliveryServiceNullableFromV12(ds)
	dsv13, errCode, userErr, sysErr = create(inf.Tx.Tx, *inf.Config, inf.User, dsv13)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Deliveryservice creation was successful.", []tc.DeliveryServiceNullableV12{dsv13.DeliveryServiceNullableV12})
}

func (ds *TODeliveryServiceV12) Read() ([]interface{}, error, error, int) {
	returnable := []interface{}{}
	dses, errs, _ := readGetDeliveryServices(ds.APIInfo().Params, ds.APIInfo().Tx, ds.APIInfo().User)
	if len(errs) > 0 {
		for _, err := range errs {
			if err.Error() == `id cannot parse to integer` {
				return nil, errors.New("Resource not found."), nil, http.StatusNotFound //matches perl response
			}
		}
		return nil, nil, errors.New("reading ds v12: " + util.JoinErrsStr(errs)), http.StatusInternalServerError
	}

	for _, ds := range dses {
		returnable = append(returnable, ds.DeliveryServiceNullableV12)
	}
	return returnable, nil, nil, http.StatusOK
}

//Delete is the DeliveryService implementation of the Deleter interface
//all implementations of Deleter should use transactions and return the proper errorType
func (ds *TODeliveryServiceV12) Delete() (error, error, int) {
	if ds.ID == nil {
		return errors.New("missing id"), nil, http.StatusBadRequest
	}
	xmlID, ok, err := ds.GetXMLID(ds.ReqInfo.Tx)
	if err != nil {
		return nil, errors.New("dsv12 delete: getting xmlid: " + err.Error()), http.StatusInternalServerError
	} else if !ok {
		return errors.New("delivery service not found"), nil, http.StatusNotFound
	}
	ds.XMLID = &xmlID

	// Note ds regexes MUST be deleted before the ds, because there's a ON DELETE CASCADE on deliveryservice_regex (but not on regex).
	// Likewise, it MUST happen in a transaction with the later DS delete, so they aren't deleted if the DS delete fails.
	if _, err := ds.ReqInfo.Tx.Tx.Exec(`DELETE FROM regex WHERE id IN (SELECT regex FROM deliveryservice_regex WHERE deliveryservice=$1)`, *ds.ID); err != nil {
		return nil, errors.New("TODeliveryServiceV12.Delete deleting regexes for delivery service: " + err.Error()), http.StatusInternalServerError
	}

	if _, err := ds.ReqInfo.Tx.Tx.Exec(`DELETE FROM deliveryservice_regex WHERE deliveryservice=$1`, *ds.ID); err != nil {
		return nil, errors.New("TODeliveryServiceV12.Delete deleting delivery service regexes: " + err.Error()), http.StatusInternalServerError
	}

	userErr, sysErr, errCode := api.GenericDelete(ds)
	if userErr != nil || sysErr != nil {
		return userErr, sysErr, errCode
	}

	paramConfigFilePrefixes := []string{"hdr_rw_", "hdr_rw_mid_", "regex_remap_", "cacheurl_"}
	configFiles := []string{}
	for _, prefix := range paramConfigFilePrefixes {
		configFiles = append(configFiles, prefix+*ds.XMLID+".config")
	}

	if _, err := ds.ReqInfo.Tx.Tx.Exec(`DELETE FROM parameter WHERE name = 'location' AND config_file = ANY($1)`, pq.Array(configFiles)); err != nil {
		return nil, errors.New("TODeliveryServiceV12.Delete deleting delivery service parameteres: " + err.Error()), http.StatusInternalServerError
	}

	return nil, nil, http.StatusOK
}

// Update is unimplemented, needed to satisfy CRUDer, since the framework doesn't allow an update to return an array of one.
func (ds *TODeliveryServiceV12) Update() (error, error, int) {
	return nil, nil, http.StatusNotImplemented
}

func UpdateV12(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"id"}, []string{"id"})
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	ds := tc.DeliveryServiceNullableV12{}
	ds.ID = util.IntPtr(inf.IntParams["id"])
	if err := api.Parse(r.Body, inf.Tx.Tx, &ds); err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusBadRequest, errors.New("decoding: "+err.Error()), nil)
		return
	}
	dsv13 := tc.NewDeliveryServiceNullableFromV12(ds)
	dsv13, errCode, userErr, sysErr = update(inf.Tx.Tx, *inf.Config, inf.User, &dsv13)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	api.WriteRespAlertObj(w, r, tc.SuccessLevel, "Deliveryservice update was successful.", []tc.DeliveryServiceNullableV12{dsv13.DeliveryServiceNullableV12})
}

// GetDeliveryServiceType returns the type of the deliveryservice.
func GetDeliveryServiceType(dsID int, tx *sql.Tx) (tc.DSType, error) {
	var dsType tc.DSType
	if err := tx.QueryRow(`SELECT t.name FROM deliveryservice as ds JOIN type t ON ds.type = t.id WHERE ds.id=$1`, dsID).Scan(&dsType); err != nil {
		if err == sql.ErrNoRows {
			return tc.DSTypeInvalid, errors.New("a deliveryservice with id '" + strconv.Itoa(dsID) + "' was not found")
		}
		return tc.DSTypeInvalid, errors.New("querying type from delivery service: " + err.Error())
	}
	return dsType, nil
}
