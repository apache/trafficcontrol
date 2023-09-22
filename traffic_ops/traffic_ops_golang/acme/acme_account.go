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
	"errors"
	"fmt"
	"net/http"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/api"
)

const readQuery = `SELECT email, private_key, uri, provider FROM acme_account`
const readProvidersQuery = `SELECT DISTINCT provider FROM acme_account`
const createQuery = `INSERT INTO acme_account (email, private_key, uri, provider) VALUES (:email, :private_key, :uri, :provider) RETURNING email, provider`
const updateQuery = `UPDATE acme_account SET private_key=:private_key, uri=:uri WHERE email=:email and provider=:provider RETURNING email, provider`
const deleteQuery = `DELETE FROM acme_account WHERE email=$1 and provider=$2`
const selectByProviderAndEmailQuery = `SELECT email, private_key, uri, provider from acme_account where email = $1 and provider = $2`
const selectLimitedQuery = `SELECT email, provider from acme_account where email = $1 and provider = $2`

// Read handles GET requests for all information about the ACME accounts.
func Read(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	acmeAccounts := []tc.AcmeAccount{}
	rows, err := tx.Query(readQuery)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("querying acme accounts: "+err.Error()))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var acct tc.AcmeAccount
		if err = rows.Scan(&acct.Email, &acct.PrivateKey, &acct.Uri, &acct.Provider); err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("scanning acme accounts: "+err.Error()))
			return
		}
		acmeAccounts = append(acmeAccounts, acct)
	}

	api.WriteResp(w, r, acmeAccounts)
}

// ReadProviders returns a list of unique ACME provider both from the database and cdn.conf
func ReadProviders(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	acmeProviders := []string{}
	rows, err := tx.Query(readProvidersQuery)
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("querying acme account providers: "+err.Error()))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var provider string
		if err = rows.Scan(&provider); err != nil {
			api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New("scanning acme account providers: "+err.Error()))
			return
		}
		acmeProviders = append(acmeProviders, provider)
	}

	for _, acmeCfg := range inf.Config.AcmeAccounts {
		alreadyInList := false
		for _, acmeProvider := range acmeProviders {
			if acmeCfg.AcmeProvider == acmeProvider {
				alreadyInList = true
			}
		}
		if !alreadyInList {
			acmeProviders = append(acmeProviders, acmeCfg.AcmeProvider)
		}
	}

	api.WriteResp(w, r, acmeProviders)
}

// Create handles POST requests to add a new ACME provider.
func Create(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, inf.Tx.Tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	tx := inf.Tx.Tx

	var acmeAccount tc.AcmeAccount
	if err := api.Parse(r.Body, tx, &acmeAccount); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	var prevEmail string
	var prevProvider string
	err := tx.QueryRow(selectLimitedQuery, acmeAccount.Email, acmeAccount.Provider).Scan(&prevEmail, &prevProvider)
	if err != nil && err != sql.ErrNoRows {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New(fmt.Sprintf("checking if acme account with email %s and provider %s exists: %v", *acmeAccount.Email, *acmeAccount.Provider, err.Error())))
		return
	}

	if prevEmail != "" && prevProvider != "" {
		api.HandleErr(w, r, tx, http.StatusBadRequest, errors.New("acme account already exists"), nil)
		return
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

	var acmeAccount tc.AcmeAccount
	if err := api.Parse(r.Body, tx, &acmeAccount); err != nil {
		api.HandleErr(w, r, tx, http.StatusBadRequest, err, nil)
		return
	}

	var prevAccount tc.AcmeAccount
	err := tx.QueryRow(selectByProviderAndEmailQuery, acmeAccount.Email, acmeAccount.Provider).Scan(&prevAccount.Email, &prevAccount.PrivateKey, &prevAccount.Uri, &prevAccount.Provider)
	if err == sql.ErrNoRows {
		api.HandleErr(w, r, tx, http.StatusBadRequest, errors.New(fmt.Sprintf("acme account not found")), nil)
		return
	}
	if err != nil {
		api.HandleErr(w, r, tx, http.StatusInternalServerError, nil, errors.New(fmt.Sprintf("checking if acme account with email %s and provider %s exists: %v", *acmeAccount.Email, *acmeAccount.Provider, err.Error())))
		return
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
	}

	alerts := tc.CreateAlerts(tc.SuccessLevel, "Acme account updated")
	api.WriteAlertsObj(w, r, http.StatusOK, alerts, acmeAccount)

	changeLogMsg := fmt.Sprintf("ACME ACCOUNT: %s %s, ACTION: updated", *acmeAccount.Email, *acmeAccount.Provider)
	api.CreateChangeLogRawTx(api.ApiChange, changeLogMsg, inf.User, tx)
}

// Delete removes the information about an ACME account.
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

	var prevAccount tc.AcmeAccount
	err := tx.QueryRow(selectByProviderAndEmailQuery, email, provider).Scan(&prevAccount.Email, &prevAccount.PrivateKey, &prevAccount.Uri, &prevAccount.Provider)
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
