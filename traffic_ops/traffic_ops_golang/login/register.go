package login

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
	"bytes"
	"encoding/json"
)
import "database/sql"
import "errors"
import "fmt"
import "html/template"
import "net/http"

import "github.com/apache/trafficcontrol/v8/lib/go-log"
import "github.com/apache/trafficcontrol/v8/lib/go-rfc"
import "github.com/apache/trafficcontrol/v8/lib/go-tc"

import "github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
import "github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/config"
import "github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/dbhelpers"

type registrationEmailFormatter struct {
	From         rfc.EmailAddress
	InstanceName string
	RegisterURL  string
	To           rfc.EmailAddress
	Token        string
}

const registerUserQuery = `
INSERT INTO tm_user (email,
                     new_user,
                     registration_sent,
                     role,
                     tenant_id,
                     token,
                     username)
VALUES ($1,
        TRUE,
        NOW(),
        $2,
        $3,
        $4,
        'registration_' || (SELECT md5(random()::text)))
RETURNING (
	SELECT role.name
	FROM role
	WHERE id=tm_user.role
) AS role,
(
	SELECT tenant.name
	FROM tenant
	WHERE tenant.id=tm_user.tenant_id
) AS tenant,
username
`

const renewRegistrationQuery = `
UPDATE tm_user
SET registration_sent = now(),
    role = $1,
    tenant_id = $2,
    token = $3
WHERE email = $4
RETURNING (
	SELECT role.name
	FROM role
	WHERE id=tm_user.role
) AS role,
(
	SELECT tenant.name
	FROM tenant
	WHERE tenant.id=tm_user.tenant_id
) AS tenant
`

var registrationEmailTemplate = template.Must(template.New("Registration Email").Parse("From: {{.From.Address.Address}}\r" + `
To: {{.To.Address.Address}}` + "\r" + `
Content-Type: text/html` + "\r" + `
Subject: {{.InstanceName}} New User Registration` + "\r\n\r" + `
<!DOCTYPE html>
<html lang="en">
<head>
	<title>{{.InstanceName}} New User Registration</title>
	<meta charset="utf-8"/>
	<style>
		.button_link {
			display: block;
			width: 130px;
			background: #2682AF;
			padding: 5px;
			text-align: center;
			border-radius: 5px;
			color: white;
			font-weight: bold;
			text-decoration: none;
			cursor: pointer;
		}
	</style>
</head>
<body>
	<main>
		<p>A new account has been created for you on the {{.InstanceName}} Portal. In the
		{{.InstanceName}} Portal, you'll find a dashboard that provides access to all of your
		Delivery Services.</p>
		<p><a class="button_link" href="{{.RegisterURL}}?token={{.Token}}" target="_blank">Click here to finish your registration</a></p>
	</main>
	<footer>
		<p>Thank you,<br/>
		The {{.InstanceName}} Team</p>
	</footer>
</body>
</html>
`))

func createRegistrationMsg(addr rfc.EmailAddress, t string, tx *sql.Tx, c config.ConfigPortal) ([]byte, error) {
	var instanceName string
	if err := tx.QueryRow(instanceNameQuery, tc.GlobalConfigFileName).Scan(&instanceName); err != nil {
		return nil, err
	}

	var f = registrationEmailFormatter{
		From:         c.EmailFrom,
		InstanceName: instanceName,
		RegisterURL:  c.BaseURL.String() + c.UserRegisterPath,
		To:           addr,
		Token:        t,
	}

	var tmpl bytes.Buffer
	if err := registrationEmailTemplate.Execute(&tmpl, &f); err != nil {
		return nil, err
	}
	return tmpl.Bytes(), nil
}

// RegisterUser is the handler for /users/register. It sends registration through Email.
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var tenantID uint
	var req tc.UserRegistrationRequest
	var reqV4 tc.UserRegistrationRequestV4
	var email rfc.EmailAddress

	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	var tx = inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()
	defer r.Body.Close()

	// ToDo: uncomment this once the perm based roles and config options are implemented
	if inf.Version.Major >= 4 {
		if err := json.NewDecoder(r.Body).Decode(&reqV4); err != nil {
			api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
			return
		}
		if err := reqV4.Validate(tx); err != nil {
			api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
			return
		}
		tenantID = reqV4.TenantID
		email = reqV4.Email
	} else {
		if userErr = api.Parse(r.Body, tx, &req); userErr != nil {
			api.HandleErr(w, r, tx, http.StatusBadRequest, userErr, nil)
			return
		}
		tenantID = req.TenantID
		email = req.Email
	}

	if ok, err := inf.IsResourceAuthorizedToCurrentUser(int(tenantID)); err != nil {
		sysErr = fmt.Errorf("Checking tenancy permissions of current user (%+v) on tenant #%d", inf.User, tenantID)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, tx, errCode, nil, sysErr)
		return
	} else if !ok {
		sysErr = fmt.Errorf("User %s requested unauthorized access to tenant #%d to register new user", inf.User.UserName, tenantID)
		userErr = errors.New("not authorized on this tenant")
		errCode = http.StatusForbidden
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	// ToDo: Add checks for permission based role checking here, if the version is >=5 and the config option is turned on.
	if inf.Version.Major < 4 {
		privLevel, ok, err := dbhelpers.GetPrivLevelFromRoleID(tx, int(req.Role))
		if err != nil {
			sysErr = fmt.Errorf("checking role #%d privilege level: %w", req.Role, err)
			errCode = http.StatusInternalServerError
			api.HandleErr(w, r, tx, errCode, nil, sysErr)
			return
		}
		if !ok {
			userErr = fmt.Errorf("No such role: %d", req.Role)
			errCode = http.StatusNotFound
			api.HandleErr(w, r, tx, errCode, userErr, nil)
			return
		}
		if privLevel > inf.User.PrivLevel {
			userErr = errors.New("Cannot register a user with a role with higher privileges than yourself")
			errCode = http.StatusForbidden
			api.HandleErr(w, r, tx, errCode, userErr, nil)
			return
		}
	} else {
		req.Email = reqV4.Email
		req.TenantID = reqV4.TenantID
		roleID, ok, err := dbhelpers.GetRoleIDFromName(tx, reqV4.Role)
		if err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, fmt.Errorf("error fetching ID from role name: %w", err))
			return
		} else if !ok {
			api.HandleErr(w, r, tx, http.StatusNotFound, errors.New("no such role"), nil)
			return
		}
		req.Role = uint(roleID)
	}

	t, err := generateToken()
	if err != nil {
		errCode = http.StatusInternalServerError
		sysErr = fmt.Errorf("Failed to generate token: %v", err)
		api.HandleErr(w, r, tx, errCode, nil, sysErr)
		return
	}

	var role string
	var tenant string
	var username string
	user, exists, err := dbhelpers.GetUserByEmail(email.Address.Address, inf.Tx.Tx)
	if err != nil {
		errCode = http.StatusInternalServerError
		sysErr = fmt.Errorf("Checking for existing user with email %s: %v", email, err)
		api.HandleErr(w, r, tx, errCode, nil, sysErr)
		return
	}
	if exists {
		if user.NewUser == nil || !*user.NewUser {
			userErr = errors.New("User already exists and has completed registration.")
			errCode = http.StatusConflict
			api.HandleErr(w, r, tx, errCode, userErr, nil)
			return
		}
		role, tenant, err = renewRegistration(tx, req, t, user)
	} else {
		role, tenant, username, err = newRegistration(tx, req, t)
	}

	if err != nil {
		log.Errorf("Bare error: %v", err)
		userErr, sysErr, errCode = api.ParseDBError(err)
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	if user.Username != nil {
		username = *user.Username
	}

	msg, err := createRegistrationMsg(email, t, tx, inf.Config.ConfigPortal)
	if err != nil {
		sysErr = fmt.Errorf("failed to create email message: %v", err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, tx, errCode, nil, sysErr)
		return
	}

	log.Debugf("Sending registration email to %s", email)

	if errCode, userErr, sysErr = inf.SendMail(email, msg); userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	var alert = "Sent user registration to %s with the following permissions [ role: %s | tenant: %s ]"
	alert = fmt.Sprintf(alert, email, role, tenant)
	api.WriteRespAlert(w, r, tc.SuccessLevel, alert)

	var changeLog = "USER: %s, EMAIL: %s, ACTION: registration sent with role %s and tenant %s"
	changeLog = fmt.Sprintf(changeLog, username, email, role, tenant)
	api.CreateChangeLogRawTx(api.ApiChange, changeLog, inf.User, tx)
}

func renewRegistration(tx *sql.Tx, req tc.UserRegistrationRequest, t string, u tc.User) (string, string, error) {
	var role string
	var tenant string

	var row = tx.QueryRow(renewRegistrationQuery, req.Role, req.TenantID, t, *u.Email)
	if err := row.Scan(&role, &tenant); err != nil {
		return "", "", err
	}

	return role, tenant, nil
}

func newRegistration(tx *sql.Tx, req tc.UserRegistrationRequest, t string) (string, string, string, error) {
	var role string
	var tenant string
	var username string
	var row = tx.QueryRow(registerUserQuery, req.Email.Address.Address, req.Role, req.TenantID, t)
	if err := row.Scan(&role, &tenant, &username); err != nil {
		return "", "", "", err
	}

	return role, tenant, username, nil
}
