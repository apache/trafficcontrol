package postgres

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
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util"

	"github.com/jmoiron/sqlx"
)

func getURLSigKeys(xmlID string, tvTx *sqlx.Tx, ctx context.Context, aesKey []byte) (tc.URLSigKeys, bool, error) {
	var encryptedUrlSigKey []byte
	if err := tvTx.QueryRow("SELECT data FROM url_sig_key WHERE deliveryservice = $1", xmlID).Scan(&encryptedUrlSigKey); err != nil {
		if err == sql.ErrNoRows {
			return tc.URLSigKeys{}, false, nil
		}
		e := checkErrWithContext("Traffic Vault PostgreSQL: executing SELECT URL Sig Keys query", err, ctx.Err())
		return tc.URLSigKeys{}, false, e
	}

	jsonUrlKeys, err := util.AESDecrypt(encryptedUrlSigKey, aesKey)
	if err != nil {
		return tc.URLSigKeys{}, false, err
	}

	urlSignKey := tc.URLSigKeys{}
	err = json.Unmarshal(jsonUrlKeys, &urlSignKey)
	if err != nil {
		return tc.URLSigKeys{}, false, errors.New("unmarshalling keys: " + err.Error())
	}

	return urlSignKey, true, nil
}

func putURLSigKeys(xmlID string, tvTx *sqlx.Tx, keys tc.URLSigKeys, ctx context.Context, aesKey []byte) error {
	keyJSON, err := json.Marshal(&keys)
	if err != nil {
		return errors.New("marshalling keys: " + err.Error())
	}

	// Delete old keys first if they exist
	if err = deleteURLSigKeys(xmlID, tvTx, ctx); err != nil {
		return err
	}

	encryptedKey, err := util.AESEncrypt(keyJSON, aesKey)
	if err != nil {
		return errors.New("encrypting keys: " + err.Error())
	}

	res, err := tvTx.Exec("INSERT INTO url_sig_key (deliveryservice, data) VALUES ($1, $2)", xmlID, encryptedKey)
	if err != nil {
		e := checkErrWithContext("Traffic Vault PostgreSQL: executing INSERT URL Sig Keys query", err, ctx.Err())
		return e
	}
	if rowsAffected, err := res.RowsAffected(); err != nil {
		return err
	} else if rowsAffected == 0 {
		return errors.New("URL Sign Keys: no keys were inserted")
	}
	return nil
}

func deleteURLSigKeys(xmlID string, tvTx *sqlx.Tx, ctx context.Context) error {
	if _, err := tvTx.Exec("DELETE FROM url_sig_key WHERE deliveryservice = $1", xmlID); err != nil {
		e := checkErrWithContext("Traffic Vault PostgreSQL: executing DELETE URL Sig Keys query", err, ctx.Err())
		return e
	}
	return nil
}
