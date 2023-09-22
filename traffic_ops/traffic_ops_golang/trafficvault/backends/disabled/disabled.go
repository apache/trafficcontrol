// Package disabled provides a TrafficVault implementation that simply returns an
// error for every method stating that Traffic Vault is disabled. This is used instead
// of passing around a nil TrafficVault instance when Traffic Vault is not enabled, in
// order to reduce the likelihood of accidentally de-referencing a nil pointer.
package disabled

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

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

const disabledErr = Error("traffic vault is disabled")

type Disabled struct {
}

func (d *Disabled) GetDeliveryServiceSSLKeys(xmlID string, version string, tx *sql.Tx, ctx context.Context) (tc.DeliveryServiceSSLKeysV15, bool, error) {
	return tc.DeliveryServiceSSLKeysV15{}, false, disabledErr
}

func (d *Disabled) GetExpirationInformation(tx *sql.Tx, ctx context.Context, days int) ([]tc.SSLKeyExpirationInformation, error) {
	return []tc.SSLKeyExpirationInformation{}, disabledErr
}

func (d *Disabled) PutDeliveryServiceSSLKeys(key tc.DeliveryServiceSSLKeys, tx *sql.Tx, ctx context.Context) error {
	return disabledErr
}

func (d *Disabled) DeleteDeliveryServiceSSLKeys(xmlID string, version string, tx *sql.Tx, ctx context.Context) error {
	return disabledErr
}

func (d *Disabled) DeleteOldDeliveryServiceSSLKeys(existingXMLIDs map[string]struct{}, cdnName string, tx *sql.Tx, ctx context.Context) error {
	return disabledErr
}

func (d *Disabled) GetCDNSSLKeys(cdnName string, tx *sql.Tx, ctx context.Context) ([]tc.CDNSSLKey, error) {
	return nil, disabledErr
}

func (d *Disabled) GetDNSSECKeys(cdnName string, tx *sql.Tx, ctx context.Context) (tc.DNSSECKeysTrafficVault, bool, error) {
	return nil, false, disabledErr
}

func (d *Disabled) PutDNSSECKeys(cdnName string, keys tc.DNSSECKeysTrafficVault, tx *sql.Tx, ctx context.Context) error {
	return disabledErr
}

func (d *Disabled) DeleteDNSSECKeys(cdnName string, tx *sql.Tx, ctx context.Context) error {
	return disabledErr
}

func (d *Disabled) GetURLSigKeys(xmlID string, tx *sql.Tx, ctx context.Context) (tc.URLSigKeys, bool, error) {
	return nil, false, disabledErr
}

func (d *Disabled) PutURLSigKeys(xmlID string, keys tc.URLSigKeys, tx *sql.Tx, ctx context.Context) error {
	return disabledErr
}

func (d *Disabled) DeleteURLSigKeys(xmlID string, tx *sql.Tx, ctx context.Context) error {
	return disabledErr
}

func (d *Disabled) GetURISigningKeys(xmlID string, tx *sql.Tx, ctx context.Context) ([]byte, bool, error) {
	return nil, false, disabledErr
}

func (d *Disabled) PutURISigningKeys(xmlID string, keysJson []byte, tx *sql.Tx, ctx context.Context) error {
	return disabledErr
}

func (d *Disabled) DeleteURISigningKeys(xmlID string, tx *sql.Tx, ctx context.Context) error {
	return disabledErr
}
func (d *Disabled) Ping(tx *sql.Tx, ctx context.Context) (tc.TrafficVaultPing, error) {
	return tc.TrafficVaultPing{}, disabledErr
}

func (d *Disabled) GetBucketKey(bucket string, key string, tx *sql.Tx) ([]byte, bool, error) {
	return nil, false, disabledErr
}
