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
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/tenant"
	"github.com/go-ozzo/ozzo-validation"
)

type TOServiceCategory struct {
	api.APIInfoImpl `json:"-"`
	tc.ServiceCategory
}

func (v *TOServiceCategory) SetLastUpdated(t tc.TimeNoMod) { v.LastUpdated = t }
func (v *TOServiceCategory) InsertQuery() string           { return insertQuery() }
func (v *TOServiceCategory) NewReadObj() interface{}       { return &tc.ServiceCategory{} }
func (v *TOServiceCategory) SelectQuery() string           { return selectQuery() }
func (v *TOServiceCategory) UpdateQuery() string           { return updateQuery() }
func (v *TOServiceCategory) DeleteQuery() string           { return deleteQuery() }

func (serviceCategory TOServiceCategory) GetAuditName() string {
	if serviceCategory.Name != "" {
		return serviceCategory.Name
	}
	if serviceCategory.ID != 0 {
		return strconv.Itoa(serviceCategory.ID)
	}
	if serviceCategory.TenantID != 0 {
		return strconv.Itoa(serviceCategory.TenantID)
	}
	return "unknown"
}

func (serviceCategory TOServiceCategory) GetKeyFieldsInfo() []api.KeyFieldInfo {
	return []api.KeyFieldInfo{{"id", api.GetIntKey}}
}

//Implementation of the Identifier, Validator interface functions
func (serviceCategory TOServiceCategory) GetKeys() (map[string]interface{}, bool) {
	if serviceCategory.ID == 0 {
		return map[string]interface{}{"id": 0}, false
	}
	return map[string]interface{}{"id": serviceCategory.ID}, true
}

func (serviceCategory *TOServiceCategory) SetKeys(keys map[string]interface{}) {
	i, _ := keys["id"].(int) //this utilizes the non panicking type assertion, if the thrown away ok variable is false i will be the zero of the type, 0 here.
	serviceCategory.ID = i
}

func (serviceCategory TOServiceCategory) GetType() string {
	return "serviceCategory"
}

func (serviceCategory *TOServiceCategory) ParamColumns() map[string]dbhelpers.WhereColumnInfo {
	return map[string]dbhelpers.WhereColumnInfo{
		"id":         dbhelpers.WhereColumnInfo{"sc.id", api.IsInt},
		"name":       dbhelpers.WhereColumnInfo{"sc.name", nil},
		"tenantId":   dbhelpers.WhereColumnInfo{"sc.tenant_id", api.IsInt},
		"tenantName": dbhelpers.WhereColumnInfo{"sc.tenant", nil},
	}
}

func (serviceCategory *TOServiceCategory) SelectMaxLastUpdatedQuery(where, orderBy, pagination, tableName string) string {
	return `SELECT max(t) from (
		SELECT max(last_updated) as t from ` + tableName + ` sc ` + where + orderBy + pagination +
		` UNION ALL
	select max(last_updated) as t from last_deleted l where l.table_name='` + tableName + `') as res`
}

func (serviceCategory TOServiceCategory) Validate() error {
	errs := validation.Errors{
		"name":     validation.Validate(serviceCategory.Name, validation.NotNil, validation.Required),
		"tenantId": validation.Validate(serviceCategory.TenantID, validation.Min(1)),
	}
	return util.JoinErrs(tovalidate.ToErrors(errs))
}

func (serviceCategory *TOServiceCategory) Create() (error, error, int) {
	return api.GenericCreate(serviceCategory)
}

func (serviceCategory *TOServiceCategory) Read(h http.Header, useIMS bool) ([]interface{}, error, error, int, *time.Time) {
	tenantIDs, err := tenant.GetUserTenantIDListTx(serviceCategory.APIInfo().Tx.Tx, serviceCategory.APIInfo().User.TenantID)
	if err != nil {
		return nil, nil, errors.New("getting tenant list for user: " + err.Error()), http.StatusInternalServerError, nil
	}

	serviceCategories, userErr, sysErr, errCode, maxTime := api.GenericRead(h, serviceCategory, useIMS)
	if userErr != nil || sysErr != nil {
		return nil, userErr, sysErr, errCode, nil
	}

	filteredServiceCategories := []interface{}{}
	for _, sc := range serviceCategories {
		sc1 := sc.(*tc.ServiceCategory)
		if checkTenancy(sc1, tenantIDs) {
			filteredServiceCategories = append(filteredServiceCategories, sc1)
		}
	}

	return filteredServiceCategories, nil, nil, errCode, maxTime
}

func checkTenancy(category *tc.ServiceCategory, tenantIDs []int) bool {
	for _, tenantID := range tenantIDs {
		if tenantID == category.TenantID {
			return true
		}
	}
	return false
}

func (serviceCategory *TOServiceCategory) Update() (error, error, int) {
	return api.GenericUpdate(serviceCategory)
}
func (serviceCategory *TOServiceCategory) Delete() (error, error, int) {
	return api.GenericDelete(serviceCategory)
}

func insertQuery() string {
	return `INSERT INTO service_category (name, tenant_id) VALUES (:name, :tenant_id) RETURNING id,last_updated`
}

func selectQuery() string {
	return `SELECT
sc.id,
sc.tenant_id,
t.name as tenant,
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
