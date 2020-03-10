package servicecategory
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
"strconv"
"strings"

"github.com/apache/trafficcontrol/lib/go-tc"
"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
"github.com/apache/trafficcontrol/lib/go-util"
"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"

validation "github.com/go-ozzo/ozzo-validation"
)

//we need a type alias to define functions on
type TOServiceCategory struct {
	api.APIInfoImpl `json:"-"`
	tc.ServiceCategoryNullable
}

func (v *TOServiceCategory) SetLastUpdated(t tc.TimeNoMod) { v.LastUpdated = &t }
func (v *TOServiceCategory) InsertQuery() string           { return insertQuery() }
func (v *TOServiceCategory) NewReadObj() interface{}       { return &tc.ServiceCategory{} }
func (v *TOServiceCategory) SelectQuery() string           { return selectQuery() }
func (v *TOServiceCategory) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		"id":   dbhelpers.WhereColumnInfo{"id", api.IsInt},
		"name": dbhelpers.WhereColumnInfo{"name", nil},
	}
}
func (v *TOServiceCategory) UpdateQuery() string { return updateQuery() }
func (v *TOServiceCategory) DeleteQuery() string { return deleteQuery() }

func (serviceCategory TOServiceCategory) GetAuditName() string {
	if serviceCategory.Name != nil {
		return *serviceCategory.Name
	}
	if serviceCategory.ID != nil {
		return strconv.Itoa(*serviceCategory.ID)
	}
	return "unknown"
}

func (serviceCategory TOServiceCategory) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{"id", api.GetIntKey}}
}

//Implementation of the Identifier, Validator interface functions
func (serviceCategory TOServiceCategory) GetKeys() (map[string]interface{}, bool) {
	if serviceCategory.ID == nil {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": *serviceCategory.ID}, true
}

func (serviceCategory *TOServiceCategory) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	serviceCategory.ID = &i
}

func (serviceCategory TOServiceCategory) GetType() string {
	return "serviceCategory"
}

func (serviceCategory TOServiceCategory) Validate() error {
	errs := validation.Errors{
		"name": validation.Validate(serviceCategory.Name, validation.NotNil, validation.Required),
	}
	return util.JoinErrs(tovalidate.ToErrors(errs))
}

func (dv *TOServiceCategory) Create() (error, error, int) { return api.GenericCreate(dv) }
func (dv *TOServiceCategory) Read() ([]interface{}, error, error, int) {
	params := dv.APIInfo().Params
	// TODO move to router, and do for all endpoints
	if strings.HasSuffix(params["name"], ".json") {
		params["name"] = params["name"][:len(params["name"])-len(".json")]
	}
	return api.GenericRead(dv)
}
func (dv *TOServiceCategory) Update() (error, error, int) { return api.GenericUpdate(dv) }
func (dv *TOServiceCategory) Delete() (error, error, int) { return api.GenericDelete(dv) }

func insertQuery() string {
	return `INSERT INTO service_category (name) VALUES (:name) RETURNING id,last_updated`
}

func selectQuery() string {
	return `SELECT
id,
last_updated,
name
FROM service_category`
}

func updateQuery() string {
	return `UPDATE
service_category SET
name=:name
WHERE id=:id RETURNING last_updated`
}

func deleteQuery() string {
	return `DELETE FROM service_category WHERE id=:id`
}
