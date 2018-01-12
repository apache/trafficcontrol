package auth

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

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/jmoiron/sqlx"
)

type Tenant struct {
	ID       int
	Name     string
	Active   bool
	ParentID int
}

// returns a Tenant list that the specified user has access too.
// NOTE: This method does not use the use_tenancy parameter and if this method is being used
// to control tenancy the parameter must be checked. The method IsResourceAuthorizedToUser checks the use_tenancy parameter
// and should be used for this purpose in most cases.
func GetUserTenantList(user CurrentUser, db *sqlx.DB) ([]Tenant, error) {
	query := `WITH RECURSIVE q AS (SELECT id, name, active, parent_id FROM tenant WHERE id = $1
	UNION SELECT t.id, t.name, t.active, t.parent_id  FROM tenant t JOIN q ON q.id = t.parent_id)
	SELECT id, name, active, parent_id FROM q;`

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
func IsResourceAuthorizedToUser(resourceTenantID int, user CurrentUser, db *sqlx.DB) (bool, error) {
	// $1 is the user tenant ID and $2 is the resource tenant ID
	query := `WITH RECURSIVE q AS (SELECT id, active FROM tenant WHERE id = $1
	UNION SELECT t.id, t.active FROM TENANT t JOIN q ON q.id = t.parent_id),
	tenancy AS (SELECT COALESCE(value::boolean,FALSE) AS value FROM parameter WHERE name = 'use_tenancy' AND config_file = 'global' UNION ALL SELECT FALSE FETCH FIRST 1 ROW ONLY)
	SELECT id, active, tenancy.value AS use_tenancy FROM tenancy, q WHERE id = $2 UNION ALL SELECT -1, false, tenancy.value AS use_tenancy FROM tenancy FETCH FIRST 1 ROW ONLY;`

	var tenantID int
	var active bool
	var useTenancy bool

	err := db.QueryRow(query, user.TenantID, resourceTenantID).Scan(&tenantID, &active, &useTenancy)

	switch {
	case err == sql.ErrNoRows:
		log.Errorf("checking user tenant %v access on resourceTenant %v: user has no access", user.TenantID, resourceTenantID)
		return false, nil
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
			return false, nil
		}
	}
}
