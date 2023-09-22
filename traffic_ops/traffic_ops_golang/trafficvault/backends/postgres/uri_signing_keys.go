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
	"errors"

	"github.com/apache/trafficcontrol/v8/lib/go-util"

	"github.com/jmoiron/sqlx"
)

func getURISigningKeys(xmlID string, tvTx *sqlx.Tx, ctx context.Context, aesKey []byte) ([]byte, bool, error) {
	var encryptedUriSigningKey []byte
	if err := tvTx.QueryRow("SELECT data FROM uri_signing_key WHERE deliveryservice = $1", xmlID).Scan(&encryptedUriSigningKey); err != nil {
		if err == sql.ErrNoRows {
			return []byte{}, false, nil
		}
		e := checkErrWithContext("Traffic Vault PostgreSQL: executing SELECT URI Sig Keys query", err, ctx.Err())
		return []byte{}, false, e
	}

	jsonUriKeys, err := util.AESDecrypt(encryptedUriSigningKey, aesKey)
	if err != nil {
		return []byte{}, false, err
	}

	return jsonUriKeys, true, nil
}

func putURISigningKeys(xmlID string, tvTx *sqlx.Tx, keys []byte, ctx context.Context, aesKey []byte) error {
	// Delete old keys first if they exist
	if err := deleteURISigningKeys(xmlID, tvTx, ctx); err != nil {
		return err
	}

	encryptedKey, err := util.AESEncrypt(keys, aesKey)
	if err != nil {
		return errors.New("encrypting keys: " + err.Error())
	}

	res, err := tvTx.Exec("INSERT INTO uri_signing_key (deliveryservice, data) VALUES ($1, $2)", xmlID, encryptedKey)
	if err != nil {
		e := checkErrWithContext("Traffic Vault PostgreSQL: executing INSERT URI Sig Keys query", err, ctx.Err())
		return e
	}
	if rowsAffected, err := res.RowsAffected(); err != nil {
		return err
	} else if rowsAffected == 0 {
		return errors.New("URI Sign Keys: no keys were inserted")
	}
	return nil
}

func deleteURISigningKeys(xmlID string, tvTx *sqlx.Tx, ctx context.Context) error {
	if _, err := tvTx.Exec("DELETE FROM uri_signing_key WHERE deliveryservice = $1", xmlID); err != nil {
		e := checkErrWithContext("Traffic Vault PostgreSQL: executing DELETE URI Sig Keys query", err, ctx.Err())
		return e
	}
	return nil
}
