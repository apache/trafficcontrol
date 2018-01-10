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
func GetUserTenantList(user CurrentUser, db *sqlx.DB) ([]Tenant, error) {
	query := `WITH RECURSIVE q AS (SELECT id, name, active, parent_id FROM tenant where id = $1
	UNION SELECT t.id, t.name, t.active, t.parent_id  FROM TENANT t JOIN q on q.id = t.parent_id)
	SELECT id, name, active, parent_id from q;`

	var tenantID int
	var name string
	var active bool
	var parentID int

	rows, err := db.Query(query, user.TenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	Tenants := make([]Tenant, 0, 3)

	for rows.Next() {
		if err := rows.Scan(&tenantID, &name, &active, &parentID); err != nil {
			return nil, err
		}
		Tenants = append(Tenants, Tenant{ID: tenantID, Name: name, Active: active, ParentID: parentID})
	}

	return Tenants, nil
}

func IsResourceAuthorizedToUser(resourceTenantID int, user CurrentUser, db *sqlx.DB) (bool, error) {
	// $1 is the user tenant ID and $2 is the resource tenant ID
	query := `WITH RECURSIVE q AS (SELECT id, active FROM tenant where id = $1
	UNION SELECT t.id, t.active FROM TENANT t JOIN q on q.id = t.parent_id)
	SELECT id, active from q where id = $2;`

	var tenantId int
	var active bool

	err := db.QueryRow(query, user.TenantID, resourceTenantID).Scan(&tenantId, &active)

	switch {
	case err == sql.ErrNoRows:
		log.Errorf("checking user tenant %v access on resourceTenant %v: user has no access", user.TenantID, resourceTenantID)
		return false, nil
	case err != nil:
		log.Errorf("Error checking user tenant %v access on resourceTenant  %v: %v", user.TenantID, resourceTenantID, err.Error())
		return false, err
	default:
		if tenantId == resourceTenantID && active == true {
			return true, nil
		} else {
			return false, nil
		}
	}
}
