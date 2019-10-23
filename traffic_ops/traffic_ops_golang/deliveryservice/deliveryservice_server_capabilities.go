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
	"errors"
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	deliveryServiceQueryParam  = "deliveryServiceID"
	serverCapabilityQueryParam = "serverCapability"
	xmlIDQueryParam            = "xmlID"
)

// ServerCapability provides a type to define methods on.
type ServerCapability struct {
	api.APIInfoImpl `json:"-"`
	tc.DeliveryServiceServerCapability
}

// SetLastUpdated implements the api.GenericCreator interfaces and
// sets the timestamp on insert.
func (sc *ServerCapability) SetLastUpdated(t tc.TimeNoMod) { sc.LastUpdated = &t }

// NewReadObj implements the api.GenericReader interfaces.
func (sc *ServerCapability) NewReadObj() interface{} {
	return &tc.DeliveryServiceServerCapability{}
}

// SelectQuery implements the api.GenericReader interface.
func (sc *ServerCapability) SelectQuery() string {
	return `SELECT
	sc.server_capability,
	sc.deliveryservice_id,
	ds.xml_id,
	sc.last_updated
	FROM deliveryservice_server_capability sc
	JOIN deliveryservice ds ON ds.id = sc.deliveryservice_id`
}

// ParamColumns implements the api.GenericReader interface.
func (sc *ServerCapability) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		deliveryServiceQueryParam: dbhelpers.WhereColumnInfo{
			Column:  "sc.deliveryservice_id",
			Checker: api.IsInt,
		},
		xmlIDQueryParam: dbhelpers.WhereColumnInfo{
			Column:  "ds.xml_id",
			Checker: nil,
		},
		serverCapabilityQueryParam: dbhelpers.WhereColumnInfo{
			Column:  "sc.server_capability",
			Checker: nil,
		},
	}
}

// DeleteQuery implements the api.GenericDeleter interface.
func (sc *ServerCapability) DeleteQuery() string {
	return `DELETE FROM deliveryservice_server_capability
	WHERE deliveryservice_id = :deliveryservice_id AND server_capability = :server_capability`
}

// GetKeyFieldsInfo implements the api.Identifier interface.
func (sc ServerCapability) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{
		{
			Field: deliveryServiceQueryParam,
			Func:  api.GetIntKey,
		},
		{
			Field: xmlIDQueryParam,
			Func:  api.GetStringKey,
		},
		{
			Field: serverCapabilityQueryParam,
			Func:  api.GetStringKey,
		},
	}
}

// GetKeys implements the api.Identifier interface and is not needed
// because Update is not available.
func (sc ServerCapability) GetKeys() (map[string]interface{}, bool) {
	return nil, false
}

// SetKeys implements the api.Identifier interface and allows the
// create handler to assign deliveryServiceID and serverCapability.
func (sc *ServerCapability) SetKeys(keys map[string]interface{}) {
	// this utilizes the non panicking type assertion, if the thrown
	// away ok variable is false it will be the zero of the type.
	id, _ := keys[deliveryServiceQueryParam].(int)
	sc.DeliveryServiceID = &id

	capability, _ := keys[serverCapabilityQueryParam].(string)
	sc.ServerCapability = &capability
}

// GetAuditName implements the api.Identifier interface and
// returns the name of the object.
func (sc *ServerCapability) GetAuditName() string {
	if sc.ServerCapability != nil {
		return *sc.ServerCapability
	}
	return "unknown"
}

// GetType implements the api.Identifier interface and
// returns the name of the struct.
func (sc *ServerCapability) GetType() string {
	return "deliveryservice.ServerCapability"
}

// Validate implements the api.Validator interface.
func (sc ServerCapability) Validate() error {
	errs := validation.Errors{
		deliveryServiceQueryParam:  validation.Validate(sc.DeliveryServiceID, validation.Required),
		serverCapabilityQueryParam: validation.Validate(sc.ServerCapability, validation.Required),
	}

	return util.JoinErrs(tovalidate.ToErrors(errs))
}

// Update implements the api.CRUDer interface.
func (sc *ServerCapability) Update() (error, error, int) {
	return nil, nil, http.StatusNotImplemented
}

// Read implements the api.CRUDer interface.
func (sc *ServerCapability) Read() ([]interface{}, error, error, int) {
	return api.GenericRead(sc)
}

// Delete implements the api.CRUDer interface.
func (sc *ServerCapability) Delete() (error, error, int) {
	return api.GenericDelete(sc)
}

// Create implements the api.CRUDer interface.
func (sc *ServerCapability) Create() (error, error, int) {
	resultRows, err := sc.APIInfo().Tx.NamedQuery(scInsertQuery(), sc)
	if err != nil {
		return api.ParseDBError(err)
	}
	defer resultRows.Close()

	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.StructScan(&sc); err != nil {
			return nil, errors.New(sc.GetType() + " create scanning: " + err.Error()), http.StatusInternalServerError
		}
	}
	if rowsAffected == 0 {
		return nil, errors.New(sc.GetType() + " create: no " + sc.GetType() + " was inserted, no rows was returned"), http.StatusInternalServerError
	} else if rowsAffected > 1 {
		return nil, errors.New("too many rows returned from " + sc.GetType() + " insert"), http.StatusInternalServerError
	}

	return nil, nil, http.StatusOK
}

func scInsertQuery() string {
	return `INSERT INTO deliveryservice_server_capability (
server_capability,
deliveryservice_id) VALUES (
:server_capability,
:deliveryservice_id) RETURNING deliveryservice_id, server_capability, last_updated`
}
