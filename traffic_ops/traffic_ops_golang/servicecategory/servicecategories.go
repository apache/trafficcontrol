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
	"fmt"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"

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
func (v *TOServiceCategory) UpdateQuery() string { return updateQuery() }
func (v *TOServiceCategory) DeleteQuery() string { return deleteQuery() }

func (serviceCategory TOServiceCategory) GetAuditName() string {
	if serviceCategory.Name != nil {
		return *serviceCategory.Name
	}
	if serviceCategory.ID != nil {
		return strconv.Itoa(*serviceCategory.ID)
	}
	if serviceCategory.TenantID != nil {
		return strconv.Itoa(*serviceCategory.TenantID)
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
		"name": 		validation.Validate(serviceCategory.Name, validation.NotNil, validation.Required),
		"tenantId":		validation.Validate(serviceCategory.TenantID, validation.Min(1)),
	}
	return util.JoinErrs(tovalidate.ToErrors(errs))
}

func (serviceCategory *TOServiceCategory) Create() (error, error, int) { return api.GenericCreate(serviceCategory) }

func (serviceCategory *TOServiceCategory) Read() ([]interface{}, error, error, int) {
	returnable := []interface{}{}

	serviceCategories, userErr, sysErr, errCode := getServiceCategories(serviceCategory.ReqInfo.Params, serviceCategory.ReqInfo.Tx, serviceCategory.ReqInfo.User)
	if userErr != nil || sysErr != nil {
		return nil, userErr, sysErr, errCode
	}

	for _, serviceCategory := range serviceCategories {
		returnable = append(returnable, serviceCategory)
	}

	return returnable, nil, nil, http.StatusOK
}
func (serviceCategory *TOServiceCategory) Update() (error, error, int) { return api.GenericUpdate(serviceCategory) }
func (serviceCategory *TOServiceCategory) Delete() (error, error, int) { return api.GenericDelete(serviceCategory) }

func getServiceCategories(params map[string]string, tx *sqlx.Tx, user *auth.CurrentUser) ([]tc.ServiceCategory, error, error, int) {

	// Query Parameters to Database Query column mappings
	// see the fields mapped in the SQL query
	queryParamsToSQLCols := map[string]dbhelpers.WhereColumnInfo{
		"id":   dbhelpers.WhereColumnInfo{"sc.id", api.IsInt},
		"name": dbhelpers.WhereColumnInfo{"sc.name", nil},
		"tenantId": dbhelpers.WhereColumnInfo{"sc.tenant_id", api.IsInt},
	}

	where, orderBy, pagination, queryValues, errs := dbhelpers.BuildWhereAndOrderByAndPagination(params, queryParamsToSQLCols)
	if len(errs) > 0 {
		return nil, util.JoinErrs(errs), nil, http.StatusBadRequest
	}

	tenantIDs, err := tenant.GetUserTenantIDListTx(tx.Tx, user.TenantID)

	if err != nil {
		log.Errorln("received error querying for user's tenants: " + err.Error())
		return nil, nil, tc.DBError, http.StatusInternalServerError
	}

	where, queryValues = dbhelpers.AddTenancyCheck(where, queryValues, "tenant_id", tenantIDs)

	query := selectQuery() + where + orderBy + pagination

	log.Debugln("generated serviceCategory query: " + query)
	log.Debugf("executing with values: %++v\n", queryValues)

	rows, err := tx.NamedQuery(query, queryValues)
	if err != nil {
		return nil, nil, fmt.Errorf("querying: %v", err), http.StatusInternalServerError
	}
	defer rows.Close()

	serviceCategories := []tc.ServiceCategory{}

	for rows.Next() {
		var serviceCategory tc.ServiceCategory
		err := rows.Scan(&serviceCategory.ID,
			&serviceCategory.TenantID,
			&serviceCategory.TenantName,
			&serviceCategory.LastUpdated,
			&serviceCategory.Name)

		if err != nil {
			return nil, nil, fmt.Errorf("getting service categories: %v", err), http.StatusInternalServerError
		}

		serviceCategories = append(serviceCategories, serviceCategory)
	}
	return serviceCategories, nil, nil, http.StatusOK
}

func insertQuery() string {
	return `INSERT INTO service_category (name, tenant_id) VALUES (:name, :tenant_id) RETURNING id,last_updated`
}

func selectQuery() string {
	return `SELECT
sc.id,
sc.tenant_id,
t.name,
sc.last_updated,
sc.name
FROM service_category as sc
LEFT JOIN tenant t ON sc.tenant_id = t.id`
}

func updateQuery() string {
	return `UPDATE
service_category SET
name=:name,
tenant_id=:tenant_id
WHERE id=:id RETURNING last_updated`
}

func deleteQuery() string {
	return `DELETE FROM service_category WHERE id=:id`
}
