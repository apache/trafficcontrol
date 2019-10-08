package server

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
	ServerCapabilityQueryParam = "serverCapability"
	ServerQueryParam           = "serverId"
	ServerHostNameQueryParam   = "serverHostName"
)

//we need a type alias to define functions on
type TOServerServerCapability struct {
	api.APIInfoImpl `json:"-"`
	tc.ServerServerCapability
}

func (ssc *TOServerServerCapability) SetLastUpdated(t tc.TimeNoMod) { ssc.LastUpdated = &t }
func (ssc *TOServerServerCapability) NewReadObj() interface{} {
	return &tc.ServerServerCapability{}
}
func (ssc *TOServerServerCapability) SelectQuery() string { return scSelectQuery() }
func (ssc *TOServerServerCapability) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		ServerCapabilityQueryParam: dbhelpers.WhereColumnInfo{"sc.server_capability", nil},
		ServerQueryParam:           dbhelpers.WhereColumnInfo{"s.id", api.IsInt},
		ServerHostNameQueryParam:   dbhelpers.WhereColumnInfo{"s.host_name", nil},
	}

}
func (ssc *TOServerServerCapability) DeleteQuery() string { return scDeleteQuery() }

func (ssc TOServerServerCapability) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{
		{ServerQueryParam, api.GetIntKey},
		{ServerCapabilityQueryParam, api.GetStringKey},
	}
}

// Need to satisfy Identifier interface but is a no-op as path does not have Update
func (ssc TOServerServerCapability) GetKeys() (map[string]interface{}, bool) {
	return map[string]interface{}{}, true
}

func (ssc *TOServerServerCapability) SetKeys(keys map[string]interface{}) {
	sID, _ := keys[ServerQueryParam].(int)
	ssc.ServerID = &sID

	sc, _ := keys[ServerCapabilityQueryParam].(string)
	ssc.ServerCapability = &sc
}

func (ssc *TOServerServerCapability) GetAuditName() string {
	if ssc.ServerCapability != nil {
		return *ssc.ServerCapability
	}
	return "unknown"
}

func (ssc *TOServerServerCapability) GetType() string {
	return "server server_capability"
}

// Validate fulfills the api.Validator interface
func (ssc TOServerServerCapability) Validate() error {
	errs := validation.Errors{
		ServerQueryParam:           validation.Validate(ssc.ServerID, validation.Required),
		ServerCapabilityQueryParam: validation.Validate(ssc.ServerCapability, validation.Required),
	}

	return util.JoinErrs(tovalidate.ToErrors(errs))
}

func (ssc *TOServerServerCapability) Update() (error, error, int) {
	return nil, nil, http.StatusNotImplemented
}

func (ssc *TOServerServerCapability) Read() ([]interface{}, error, error, int) {
	return api.GenericRead(ssc)
}

func (ssc *TOServerServerCapability) Delete() (error, error, int) {
	return api.GenericDelete(ssc)
}

func (ssc *TOServerServerCapability) Create() (error, error, int) {
	resultRows, err := ssc.APIInfo().Tx.NamedQuery(scInsertQuery(), ssc)
	if err != nil {
		return api.ParseDBError(err)
	}
	defer resultRows.Close()

	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
		if err := resultRows.StructScan(&ssc); err != nil {
			return nil, errors.New(ssc.GetType() + " create scanning: " + err.Error()), http.StatusInternalServerError
		}
	}
	if rowsAffected == 0 {
		return nil, errors.New(ssc.GetType() + " create: no " + ssc.GetType() + " was inserted, no rows was returned"), http.StatusInternalServerError
	} else if rowsAffected > 1 {
		return nil, errors.New("too many rows returned from " + ssc.GetType() + " insert"), http.StatusInternalServerError
	}

	return nil, nil, http.StatusOK
}

func scSelectQuery() string {
	return `SELECT
sc.server_capability,
sc.server,
sc.last_updated,
s.host_name as host_name
FROM server_server_capability sc
JOIN server s ON sc.server = s.id`
}

func scDeleteQuery() string {
	return `DELETE FROM server_server_capability
WHERE server = :server AND server_capability = :server_capability`
}

func scInsertQuery() string {
	return `INSERT INTO server_server_capability (
server_capability,
server) VALUES (
:server_capability,
:server) RETURNING server, server_capability, last_updated`
}
