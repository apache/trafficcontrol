package user

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
	"net/http"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
)

func Current(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	currentUser, err := getUser(inf.Tx.Tx, inf.User.ID)
	if err != nil {
		api.HandleErr(w, r, inf.Tx.Tx, http.StatusInternalServerError, nil, errors.New("getting current user: "+err.Error()))
		return
	}
	api.WriteResp(w, r, currentUser)
}

func getUser(tx *sql.Tx, id int) (tc.UserCurrent, error) {
	q := `
SELECT
u.address_line1,
u.address_line2,
u.city,
u.company,
u.country,
u.email,
u.full_name,
u.id,
u.last_updated,
u.local_passwd,
u.new_user,
u.phone_number,
u.postal_code,
u.public_ssh_key,
u.role,
r.name as role_name,
u.state_or_province,
t.name as tenant,
u.tenant_id,
u.username
FROM tm_user as u
LEFT JOIN role as r ON r.id = u.role
LEFT JOIN tenant as t ON t.id = u.tenant_id
WHERE u.id=$1
`
	u := tc.UserCurrent{}
	localPassword := sql.NullString{}
	if err := tx.QueryRow(q, id).Scan(&u.AddressLine1, &u.AddressLine2, &u.City, &u.Company, &u.Country, &u.Email, &u.FullName, &u.ID, &u.LastUpdated, &localPassword, &u.NewUser, &u.PhoneNumber, &u.PostalCode, &u.PublicSSHKey, &u.Role, &u.RoleName, &u.StateOrProvince, &u.Tenant, &u.TenantID, &u.UserName); err != nil {
		return tc.UserCurrent{}, errors.New("querying current user: " + err.Error())
	}
	u.LocalUser = util.BoolPtr(localPassword.Valid)
	return u, nil
}
