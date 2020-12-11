package acme

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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/go-ozzo/ozzo-validation"
	"net/http"
)

type AcmeAccount struct {
	Email      *string `json:"email" db:"email"`
	PrivateKey *string `json:"private_key" db:"private_key"`
	Uri        *string `json:"uri" db:"uri"`
	Provider   *string `json:"provider" db:"provider"`
}

func (aa *AcmeAccount) Validate(tx *sql.Tx) error {

	errs := validation.Errors{
		"email":       validation.Validate(aa.Email, validation.Required),
		"private_key": validation.Validate(aa.PrivateKey, validation.Required),
		"uri":         validation.Validate(aa.Uri, validation.Required),
		"provider":    validation.Validate(aa.Provider, validation.Required),
	}

	return util.JoinErrs(tovalidate.ToErrors(errs))
}

func (aa *AcmeAccount) ValidateUpdate(tx *sql.Tx) error {

	errs := validation.Errors{
		"email":    validation.Validate(aa.Email, validation.Required),
		"provider": validation.Validate(aa.Provider, validation.Required),
	}

	return util.JoinErrs(tovalidate.ToErrors(errs))
}

const readQuery = `SELECT email, private_key, uri, provider FROM acme_account`
const createQuery = `INSERT INTO acme_account (email, private_key, uri, provider) VALUES (:email, :private_key, :uri, :provider) RETURNING email, provider`
const updateQuery = `UPDATE acme_account SET private_key=:private_key, uri=:uri WHERE email=:email and provider=:provider RETURNING email, provider`
const deleteQuery = `DELETE FROM acme_account WHERE email=$1 and provider=$2`

func Read(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	var acmeAccounts []AcmeAccount
	rows, err := tx.Query(readQuery)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("querying acme accounts: "+err.Error()))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var acct AcmeAccount
		if err = rows.Scan(&acct.Email, &acct.PrivateKey, &acct.Uri, &acct.Provider); err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("scanning acme accounts: "+err.Error()))
			return
		}
		acmeAccounts = append(acmeAccounts, acct)
	}

	api.WriteResp(w, r, acmeAccounts)
}

func Create(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	tx := inf.Tx.Tx

	var acmeAccount AcmeAccount
	if err := json.NewDecoder(r.Body).Decode(&acmeAccount); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	if validErr := acmeAccount.Validate(tx); validErr != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, validErr, nil)
		return
	}

	if acmeAccount.Email != nil && acmeAccount.Provider != nil {
		var prevEmail string
		var prevProvider string
		err := tx.QueryRow("SELECT email, provider from acme_account where email = $1 and provider = $2", acmeAccount.Email, acmeAccount.Provider).Scan(&prevEmail, &prevProvider)
		if err != nil && err != sql.ErrNoRows {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New(fmt.Sprintf("checking if acme account with email %s and provider %s exists: %v", *acmeAccount.Email, *acmeAccount.Provider, err.Error())))
			return
		}

		if prevEmail != "" && prevProvider != "" {
			api.HandleErr(w, r, tx, http.StatusBadRequest, errors.New("acme account already exists"), nil)
			return
		}
	}

	resultRows, err := inf.Tx.NamedQuery(createQuery, acmeAccount)
	if err != nil {
		userErr, sysErr, errCode := api.ParseDBError(err)
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer resultRows.Close()

	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
	}
	if rowsAffected == 0 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("acme account create: no account was inserted"))
		return
	} else if rowsAffected > 1 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("too many rows returned from acme account insert"))
		return
	}

	alerts := tc.CreateAlerts(tc.SuccessLevel, "Acme account created")
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, acmeAccount)

	changeLogMsg := fmt.Sprintf("ACME ACCOUNT: %s %s, ACTION: created", *acmeAccount.Email, *acmeAccount.Provider)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}

func Update(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	tx := inf.Tx.Tx

	var acmeAccount AcmeAccount
	if err := json.NewDecoder(r.Body).Decode(&acmeAccount); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	if validErr := acmeAccount.ValidateUpdate(tx); validErr != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, validErr, nil)
		return
	}

	var prevAccount AcmeAccount
	err := tx.QueryRow("SELECT email, private_key, uri, provider from acme_account where email = $1 and provider = $2", acmeAccount.Email, acmeAccount.Provider).Scan(&prevAccount.Email, &prevAccount.PrivateKey, &prevAccount.Uri, &prevAccount.Provider)
	if err == sql.ErrNoRows {
		api.HandleErr(w, r, tx, http.StatusBadRequest, errors.New(fmt.Sprintf("acme account not found")), nil)
		return
	}
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New(fmt.Sprintf("checking if acme account with email %s and provider %s exists: %v", *acmeAccount.Email, *acmeAccount.Provider, err.Error())))
		return
	}

	if acmeAccount.Uri == nil {
		acmeAccount.Uri = prevAccount.Uri
	}

	if acmeAccount.PrivateKey == nil {
		acmeAccount.PrivateKey = prevAccount.PrivateKey
	}

	resultRows, err := inf.Tx.NamedQuery(updateQuery, acmeAccount)
	if err != nil {
		userErr, sysErr, errCode := api.ParseDBError(err)
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer resultRows.Close()

	rowsAffected := 0
	for resultRows.Next() {
		rowsAffected++
	}
	if rowsAffected == 0 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("acme account update: no account was updated"))
		return
	} else if rowsAffected > 1 {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("too many rows returned from acme account update"))
		return
	}

	alerts := tc.CreateAlerts(tc.SuccessLevel, "Acme account updated")
	api.WriteAlertsObj(w, r, http.StatusCreated, alerts, acmeAccount)

	changeLogMsg := fmt.Sprintf("ACME ACCOUNT: %s %s, ACTION: updated", *acmeAccount.Email, *acmeAccount.Provider)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, []string{"provider", "email"}, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	provider := inf.Params["provider"]
	email := inf.Params["email"]

	tx := inf.Tx.Tx

	var prevAccount AcmeAccount
	err := tx.QueryRow("SELECT email, private_key, uri, provider from acme_account where email = $1 and provider = $2", email, provider).Scan(&prevAccount.Email, &prevAccount.PrivateKey, &prevAccount.Uri, &prevAccount.Provider)
	if err == sql.ErrNoRows {
		api.HandleErr(w, r, tx, http.StatusBadRequest, errors.New(fmt.Sprintf("acme account not found")), nil)
		return
	}
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New(fmt.Sprintf("checking if acme account with email %s and provider %s exists: %v", email, provider, err.Error())))
		return
	}

	if _, err := tx.Exec(deleteQuery, email, provider); err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New(fmt.Sprintf("deleting acme account with email %s and provider %s: %v", email, provider, err.Error())))
		return
	}

	api.WriteRespAlert(w, r, tc.SuccessLevel, "Acme account deleted")

	changeLogMsg := fmt.Sprintf("ACME ACCOUNT: %s %s, ACTION: deleted", email, provider)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}
