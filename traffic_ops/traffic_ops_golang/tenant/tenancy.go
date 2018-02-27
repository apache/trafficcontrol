package tenant

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

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/jmoiron/sqlx"
)

type Tenant struct {
	ID       int
	Name     string
	Active   bool
	ParentID int
}

type DeliveryServiceTenantInfo tc.DeliveryServiceNullable

// returns true if the user has tenant access on this deliveryservice
func (dsInfo DeliveryServiceTenantInfo) IsTenantAuthorized(user auth.CurrentUser, db *sqlx.DB) (bool, error) {
	if dsInfo.TenantID == nil {
		return false, errors.New("TenantID is nil")
	}
	return IsResourceAuthorizedToUser(*dsInfo.TenantID, user, db)
}

// returns tenant information for a deliveryservice
func GetDeliveryServiceTenantInfo(xmlId string, db *sqlx.DB) (*DeliveryServiceTenantInfo, error) {
	ds := DeliveryServiceTenantInfo{}
	query := "SELECT xml_id,tenant_id FROM deliveryservice where xml_id = $1"

	err := db.Get(&ds, query, xmlId)
	switch {
	case err == sql.ErrNoRows:
		ds = DeliveryServiceTenantInfo{}
		return &ds, fmt.Errorf("a deliveryservice with xml_id '%s' was not found", xmlId)
	case err != nil:
		return nil, err
	default:
		return &ds, nil
	}
}

// tenancy check wrapper for deliveryservice
func HasTenant(user auth.CurrentUser, XMLID string, db *sqlx.DB) (bool, error, tc.ApiErrorType) {
	dsInfo, err := GetDeliveryServiceTenantInfo(XMLID, db)
	if err != nil {
		if dsInfo == nil {
			return false, fmt.Errorf("deliveryservice lookup failure: %v", err), tc.SystemError
		} else {
			return false, fmt.Errorf("no such deliveryservice: '%s'", XMLID), tc.DataMissingError
		}
	}
	hasAccess, err := dsInfo.IsTenantAuthorized(user, db)
	if err != nil {
		return false, fmt.Errorf("user tenancy check failure: %v", err), tc.SystemError
	}
	if !hasAccess {
		return false, fmt.Errorf("Access to this resource is not authorized"), tc.ForbiddenError
	}
	return true, nil, tc.NoError
}

// returns a Tenant list that the specified user has access too.
// NOTE: This method does not use the use_tenancy parameter and if this method is being used
// to control tenancy the parameter must be checked. The method IsResourceAuthorizedToUser checks the use_tenancy parameter
// and should be used for this purpose in most cases.
func GetUserTenantList(user auth.CurrentUser, db *sqlx.DB) ([]Tenant, error) {
	query := `WITH RECURSIVE q AS (SELECT id, name, active, parent_id FROM tenant WHERE id = $1
	UNION SELECT t.id, t.name, t.active, t.parent_id  FROM tenant t JOIN q ON q.id = t.parent_id)
	SELECT id, name, active, parent_id FROM q;`

	log.Debugln("\nQuery: ", query)

	var tenantID int
	var name string
	var active bool
	var parentID int

	rows, err := db.Query(query, user.TenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tenants := []Tenant{}

	for rows.Next() {
		if err := rows.Scan(&tenantID, &name, &active, &parentID); err != nil {
			return nil, err
		}
		tenants = append(tenants, Tenant{ID: tenantID, Name: name, Active: active, ParentID: parentID})
	}

	return tenants, nil
}

// returns a boolean value describing if the user has access to the provided resource tenant id and an error
// if use_tenancy is set to false (0 in the db) this method will return true allowing access.
func IsResourceAuthorizedToUser(resourceTenantID int, user auth.CurrentUser, db *sqlx.DB) (bool, error) {
	// $1 is the user tenant ID and $2 is the resource tenant ID
	query := `WITH RECURSIVE q AS (SELECT id, active FROM tenant WHERE id = $1
	UNION SELECT t.id, t.active FROM TENANT t JOIN q ON q.id = t.parent_id),
	tenancy AS (SELECT COALESCE(value::boolean,FALSE) AS value FROM parameter WHERE name = 'use_tenancy' AND config_file = 'global' UNION ALL SELECT FALSE FETCH FIRST 1 ROW ONLY)
	SELECT id, active, tenancy.value AS use_tenancy FROM tenancy, q WHERE id = $2 UNION ALL SELECT -1, false, tenancy.value AS use_tenancy FROM tenancy FETCH FIRST 1 ROW ONLY;`

	var tenantID int
	var active bool
	var useTenancy bool

	log.Debugln("\nQuery: ", query)
	err := db.QueryRow(query, user.TenantID, resourceTenantID).Scan(&tenantID, &active, &useTenancy)

	switch {
	case err != nil:
		log.Errorf("Error checking user tenant %v access on resourceTenant  %v: %v", user.TenantID, resourceTenantID, err.Error())
		return false, err
	default:
		if !useTenancy {
			return true, nil
		}
		if active && tenantID == resourceTenantID {
			return true, nil
		} else {
			fmt.Printf("default")
			return false, nil
		}
	}
}
