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
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/util/ims"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/lib/pq"
)

const (
	deliveryServiceQueryParam    = "deliveryServiceID"
	requiredCapabilityQueryParam = "requiredCapability"
	xmlIDQueryParam              = "xmlID"
)

// RequiredCapability provides a type to define methods on.
type RequiredCapability struct {
	api.APIInfoImpl `json:"-"`
	tc.DeliveryServicesRequiredCapability
}

// SetLastUpdated implements the api.GenericCreator interfaces and
// sets the timestamp on insert.
func (rc *RequiredCapability) SetLastUpdated(t tc.TimeNoMod) { rc.LastUpdated = &t }

// NewReadObj implements the api.GenericReader interfaces.
func (rc *RequiredCapability) NewReadObj() interface{} {
	return &tc.DeliveryServicesRequiredCapability{}
}

// SelectQuery implements the api.GenericReader interface.
func (rc *RequiredCapability) SelectQuery() string {
	return `SELECT
	rc.required_capability,
	rc.deliveryservice_id,
	ds.xml_id,
	rc.last_updated
	FROM deliveryservices_required_capability rc
	JOIN deliveryservice ds ON ds.id = rc.deliveryservice_id`
}

// ParamColumns implements the api.GenericReader interface.
func (rc *RequiredCapability) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		deliveryServiceQueryParam: dbhelpers.WhereColumnInfo{
			Column:  "rc.deliveryservice_id",
			Checker: api.IsInt,
		},
		xmlIDQueryParam: dbhelpers.WhereColumnInfo{
			Column:  "ds.xml_id",
			Checker: nil,
		},
		requiredCapabilityQueryParam: dbhelpers.WhereColumnInfo{
			Column:  "rc.required_capability",
			Checker: nil,
		},
	}
}

// DeleteQuery implements the api.GenericDeleter interface.
func (rc *RequiredCapability) DeleteQuery() string {
	return `DELETE FROM deliveryservices_required_capability
	WHERE deliveryservice_id = :deliveryservice_id AND required_capability = :required_capability`
}

// GetKeyFieldsInfo implements the api.Identifier interface.
func (rc RequiredCapability) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{
		{
			Field: deliveryServiceQueryParam,
			Func:  api.GetIntKey,
		},
		{
			Field: requiredCapabilityQueryParam,
			Func:  api.GetStringKey,
		},
	}
}

// GetKeys implements the api.Identifier interface and is not needed
// because Update is not available.
func (rc RequiredCapability) GetKeys() (map[string]interface{}, bool) {
	if rc.DeliveryServiceID == nil {
		return map[string]interface{}{deliveryServiceQueryParam: 0}, false
	}
	if rc.RequiredCapability == nil {
		return map[string]interface{}{requiredCapabilityQueryParam: 0}, false
	}
	return map[string]interface{}{
		deliveryServiceQueryParam:    *rc.DeliveryServiceID,
		requiredCapabilityQueryParam: *rc.RequiredCapability,
	}, true
}

// SetKeys implements the api.Identifier interface and allows the
// create handler to assign deliveryServiceID and requiredCapability.
func (rc *RequiredCapability) SetKeys(keys map[string]interface{}) {
	// this utilizes the non panicking type assertion, if the thrown
	// away ok variable is false it will be the zero of the type.
	id, _ := keys[deliveryServiceQueryParam].(int)
	rc.DeliveryServiceID = &id

	capability, _ := keys[requiredCapabilityQueryParam].(string)
	rc.RequiredCapability = &capability
}

// GetAuditName implements the api.Identifier interface and
// returns the name of the object.
func (rc *RequiredCapability) GetAuditName() string {
	if rc.RequiredCapability != nil {
		return *rc.RequiredCapability
	}

	return "unknown"
}

// GetType implements the api.Identifier interface and
// returns the name of the struct.
func (rc *RequiredCapability) GetType() string {
	return "deliveryservice.RequiredCapability"
}

// Validate implements the api.Validator interface.
func (rc RequiredCapability) Validate() error {
	errs := validation.Errors{
		deliveryServiceQueryParam:    validation.Validate(rc.DeliveryServiceID, validation.Required),
		requiredCapabilityQueryParam: validation.Validate(rc.RequiredCapability, validation.Required),
	}

	return util.JoinErrs(tovalidate.ToErrors(errs))
}

// Update implements the api.CRUDer interface.
func (rc *RequiredCapability) Update() (error, error, int) {
	return nil, nil, http.StatusNotImplemented
}

// Read implements the api.CRUDer interface.
func (rc *RequiredCapability) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	tenantIDs, err := rc.getTenantIDs()
	if err != nil {
		return nil, nil, err, http.StatusInternalServerError, nil
	}

	capabilities, userErr, sysErr, errCode, maxTime := rc.getCapabilities(h, tenantIDs, useIMS)
	if userErr != nil || sysErr != nil {
		return nil, userErr, sysErr, errCode, nil
	}

	results := []interface{}{}
	for _, capability := range capabilities {
		results = append(results, capability)
	}

	return results, nil, nil, errCode, maxTime
}

func (rc *RequiredCapability) getTenantIDs() ([]int, error) {
	tenantIDs, err := tenant.GetUserTenantIDListTx(rc.APIInfo().Tx.Tx, rc.APIInfo().User.TenantID)
	if err != nil {
		return nil, err
	}
	return tenantIDs, nil
}

func (rc *RequiredCapability) getCapabilities(h http.Header, tenantIDs []int, useIMS bool) ([]tc.DeliveryServicesRequiredCapability, error, error, int, *time.Time) {
	var maxTime time.Time
	var runSecond bool
	var results []tc.DeliveryServicesRequiredCapability
	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(rc.APIInfo().Params, rc.ParamColumns())
	if len(errs) > 0 {
		return nil, util.JoinErrs(errs), nil, http.StatusBadRequest, nil
	}

	where, queryValues = dbhelpers.AddTenancyCheck(where, queryValues, "ds.tenant_id", tenantIDs)
	if useIMS {
		runSecond, maxTime = ims.TryIfModifiedSinceQuery(rc.APIInfo().Tx, h, queryValues, selectMaxLastUpdatedQueryRC(where, orderBy, pagination))
		if !runSecond {
			log.Debugln("IMS HIT")
			return results, nil, nil, http.StatusNotModified, &maxTime
		}
		log.Debugln("IMS MISS")
	} else {
		log.Debugln("Non IMS request")
	}
	query := rc.SelectQuery() + where + orderBy + pagination

	rows, err := rc.APIInfo().Tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, nil, err, http.StatusInternalServerError, nil
	}
	defer rows.Close()

	for rows.Next() {
		var result tc.DeliveryServicesRequiredCapability
		if err := rows.StructScan(&result); err != nil {
			return nil, nil, fmt.Errorf("%s get scanning: %s", rc.GetType(), err.Error()), http.StatusInternalServerError, nil
		}
		results = append(results, result)
	}

	return results, nil, nil, http.StatusOK, &maxTime
}

func selectMaxLastUpdatedQueryRC(where string, orderBy string, pagination string) string {
	return `SELECT max(t) from (
		SELECT max(rc.last_updated) as t FROM deliveryservices_required_capability rc
	JOIN deliveryservice ds ON ds.id = rc.deliveryservice_id ` + where + orderBy + pagination +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='deliveryservices_required_capability') as res`
}

// Delete implements the api.CRUDer interface.
func (rc *RequiredCapability) Delete() (error, error, int) {
	authorized, err := rc.isTenantAuthorized()
	if !authorized {
		return errors.New("not authorized on this tenant"), nil, http.StatusForbidden
	} else if err != nil {
		return nil, fmt.Errorf("checking authorization for existing DS ID: %s" + err.Error()), http.StatusInternalServerError
	}

	return api.GenericDelete(rc)
}

// Create implements the api.CRUDer interface.
func (rc *RequiredCapability) Create() (error, error, int) {
	authorized, err := rc.isTenantAuthorized()
	if !authorized {
		return errors.New("not authorized on this tenant"), nil, http.StatusForbidden
	} else if err != nil {
		return nil, fmt.Errorf("checking authorization for existing DS ID: %s" + err.Error()), http.StatusInternalServerError
	}

	// Ensure DS type is only of HTTP*, DNS* types
	dsType, dsExists, err := dbhelpers.GetDeliveryServiceType(*rc.DeliveryServiceID, rc.APIInfo().Tx.Tx)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	if !dsExists {
		return errors.New("a deliveryservice with id '" + strconv.Itoa(*rc.DeliveryServiceID) + "' was not found"), nil, http.StatusNotFound
	}

	if !dsType.IsHTTP() && !dsType.IsDNS() {
		return errors.New("Invalid DS type. Only DNS and HTTP delivery services can have required capabilities"), nil, http.StatusBadRequest
	}

	usrErr, sysErr, rCode := rc.checkServerCap()
	if usrErr != nil || sysErr != nil {
		return usrErr, sysErr, rCode
	}

	usrErr, sysErr, rCode = rc.ensureDSServerCap()
	if usrErr != nil || sysErr != nil {
		return usrErr, sysErr, rCode
	}

	rows, err := rc.APIInfo().Tx.NamedQuery(rcInsertQuery(), rc)
	if err != nil {
		return api.ParseDBError(err)
	}
	defer rows.Close()

	rowsAffected := 0
	for rows.Next() {
		rowsAffected++
		if err := rows.StructScan(&rc); err != nil {
			return nil, fmt.Errorf("%s create scanning: %s", rc.GetType(), err.Error()), http.StatusInternalServerError
		}
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("%s create: no %s was inserted, no rows was returned", rc.GetType(), rc.GetType()), http.StatusInternalServerError
	} else if rowsAffected > 1 {
		return nil, fmt.Errorf("too many rows returned from %s insert", rc.GetType()), http.StatusInternalServerError
	}

	return nil, nil, http.StatusOK
}

func (rc *RequiredCapability) checkServerCap() (error, error, int) {
	tx := rc.APIInfo().Tx

	// Get server capability name
	name := ""
	if err := tx.QueryRow(`
		SELECT name 
		FROM server_capability 
		WHERE name = $1`, rc.RequiredCapability).Scan(&name); err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("querying server capability for name '%v': %v", rc.RequiredCapability, err), http.StatusInternalServerError
	}

	if len(name) == 0 {
		return fmt.Errorf("server_capability not found"), nil, http.StatusNotFound
	}

	return nil, nil, http.StatusOK
}

func (rc *RequiredCapability) ensureDSServerCap() (error, error, int) {
	tx := rc.APIInfo().Tx

	// Get assigned DS server IDs
	dsServerIDs := []int64{}
	if err := tx.Tx.QueryRow(`
	SELECT ARRAY(
		SELECT ds.server 
		FROM deliveryservice_server ds
		JOIN server s ON ds.server = s.id
		JOIN type t ON s.type = t.id
		WHERE ds.deliveryservice=$1
		AND NOT t.name LIKE 'ORG%'
	)`, rc.DeliveryServiceID).Scan(pq.Array(&dsServerIDs)); err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("reading delivery service %v servers: %v", *rc.DeliveryServiceID, err), http.StatusInternalServerError
	}

	if len(dsServerIDs) == 0 { // no attached servers can return success right away
		return nil, nil, http.StatusOK
	}

	// Get servers IDs that have the new capability
	capServerIDs := []int64{}
	if err := tx.QueryRow(`
	SELECT ARRAY(
		SELECT server
		FROM server_server_capability 
		WHERE server = ANY($1)
		AND server_capability=$2
	)`, pq.Array(dsServerIDs), rc.RequiredCapability).Scan(pq.Array(&capServerIDs)); err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("reading servers that have server capability %v attached: %v", *rc.RequiredCapability, err), http.StatusInternalServerError
	}

	vIDs := getViolatingServerIDs(dsServerIDs, capServerIDs)
	if len(vIDs) > 0 {
		return fmt.Errorf("capability %v cannot be made required on the delivery service %v as it has the associated servers %v that do not have the capability assigned", *rc.RequiredCapability, *rc.DeliveryServiceID, strings.Join(vIDs, ",")), nil, http.StatusBadRequest
	}

	return nil, nil, http.StatusOK
}

func getViolatingServerIDs(dsServerIDs, capServerIDs []int64) []string {
	capServerIDsMap := map[int64]struct{}{}
	for _, id := range capServerIDs {
		capServerIDsMap[id] = struct{}{}
	}
	vIDs := []string{}
	for _, id := range dsServerIDs {
		if _, found := capServerIDsMap[id]; !found {
			vIDs = append(vIDs, strconv.FormatInt(id, 10))
		}
	}
	return vIDs
}

func (rc *RequiredCapability) isTenantAuthorized() (bool, error) {
	if rc.DeliveryServiceID == nil && rc.XMLID == nil {
		return false, errors.New("delivery service has no ID or XMLID")
	}

	var existingID *int
	var err error

	switch {
	case rc.DeliveryServiceID != nil:
		existingID, _, err = getDSTenantIDByID(rc.APIInfo().Tx.Tx, *rc.DeliveryServiceID)
		if err != nil {
			return false, err
		}
	case rc.XMLID != nil:
		existingID, _, err = getDSTenantIDByName(rc.APIInfo().Tx.Tx, tc.DeliveryServiceName(*rc.XMLID))
		if err != nil {
			return false, err
		}
	}

	if existingID != nil {
		authorized, err := tenant.IsResourceAuthorizedToUserTx(*existingID, rc.APIInfo().User, rc.APIInfo().Tx.Tx)
		if !authorized {
			return false, errors.New("not authorized on this tenant")
		} else if err != nil {
			return false, fmt.Errorf("checking authorization for existing DS ID: %s" + err.Error())
		}
	}

	return true, err
}

func rcInsertQuery() string {
	return `INSERT INTO deliveryservices_required_capability (
required_capability,
deliveryservice_id) VALUES (
:required_capability,
:deliveryservice_id) RETURNING deliveryservice_id, required_capability, last_updated`
}

// language=SQL
const HasRequiredCapabilitiesQuery = `
SELECT EXISTS(
	SELECT drc.required_capability
	FROM deliveryservices_required_capability drc
	WHERE drc.deliveryservice_id = $1
)
`
