// Package riaksvc provides a TrafficVault implementation which uses Riak as the backend.
package riaksvc

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
	"github.com/apache/trafficcontrol/v8/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/trafficvault"

	validation "github.com/go-ozzo/ozzo-validation"
)

const RiakBackendName = "riak"

type Riak struct {
	cfg Config
}

func (r *Riak) GetDeliveryServiceSSLKeys(xmlID string, version string, tx *sql.Tx, ctx context.Context) (tc.DeliveryServiceSSLKeysV15, bool, error) {
	return getDeliveryServiceSSLKeysObjV15(xmlID, version, tx, &r.cfg.AuthOptions, &r.cfg.Port)
}

func (r *Riak) GetExpirationInformation(tx *sql.Tx, ctx context.Context, days int) ([]tc.SSLKeyExpirationInformation, error) {
	return []tc.SSLKeyExpirationInformation{}, errors.New("Not implemented for this Traffic Vault backend.")
}

func (r *Riak) PutDeliveryServiceSSLKeys(key tc.DeliveryServiceSSLKeys, tx *sql.Tx, ctx context.Context) error {
	return putDeliveryServiceSSLKeysObj(key, tx, &r.cfg.AuthOptions, &r.cfg.Port)
}

func (r *Riak) DeleteDeliveryServiceSSLKeys(xmlID string, version string, tx *sql.Tx, ctx context.Context) error {
	return deleteDSSSLKeys(tx, &r.cfg.AuthOptions, &r.cfg.Port, xmlID, version)
}

func (r *Riak) DeleteOldDeliveryServiceSSLKeys(existingXMLIDs map[string]struct{}, cdnName string, tx *sql.Tx, ctx context.Context) error {
	return deleteOldDeliveryServiceSSLKeys(tx, &r.cfg.AuthOptions, &r.cfg.Port, tc.CDNName(cdnName), existingXMLIDs)
}

func (r *Riak) GetCDNSSLKeys(cdnName string, tx *sql.Tx, ctx context.Context) ([]tc.CDNSSLKey, error) {
	return getCDNSSLKeysObj(tx, &r.cfg.AuthOptions, &r.cfg.Port, cdnName)
}

func (r *Riak) GetDNSSECKeys(cdnName string, tx *sql.Tx, ctx context.Context) (tc.DNSSECKeysTrafficVault, bool, error) {
	keys, exists, err := getDNSSECKeys(cdnName, tx, &r.cfg.AuthOptions, &r.cfg.Port)
	return tc.DNSSECKeysTrafficVault(keys), exists, err
}

func (r *Riak) PutDNSSECKeys(cdnName string, keys tc.DNSSECKeysTrafficVault, tx *sql.Tx, ctx context.Context) error {
	return putDNSSECKeys(tc.DNSSECKeysRiak(keys), cdnName, tx, &r.cfg.AuthOptions, &r.cfg.Port)
}

func (r *Riak) DeleteDNSSECKeys(cdnName string, tx *sql.Tx, ctx context.Context) error {
	return deleteDNSSECKeys(cdnName, tx, &r.cfg.AuthOptions, &r.cfg.Port)
}

func (r *Riak) GetURLSigKeys(xmlID string, tx *sql.Tx, ctx context.Context) (tc.URLSigKeys, bool, error) {
	return getURLSigKeys(tx, &r.cfg.AuthOptions, &r.cfg.Port, tc.DeliveryServiceName(xmlID))
}

func (r *Riak) PutURLSigKeys(xmlID string, keys tc.URLSigKeys, tx *sql.Tx, ctx context.Context) error {
	return putURLSigKeys(tx, &r.cfg.AuthOptions, &r.cfg.Port, tc.DeliveryServiceName(xmlID), keys)
}

func (r *Riak) DeleteURLSigKeys(xmlID string, tx *sql.Tx, ctx context.Context) error {
	return deleteURLSigningKeys(tx, &r.cfg.AuthOptions, &r.cfg.Port, tc.DeliveryServiceName(xmlID))
}

func (r *Riak) GetURISigningKeys(xmlID string, tx *sql.Tx, ctx context.Context) ([]byte, bool, error) {
	return getURISigningKeys(tx, &r.cfg.AuthOptions, &r.cfg.Port, xmlID)
}

func (r *Riak) PutURISigningKeys(xmlID string, keysJson []byte, tx *sql.Tx, ctx context.Context) error {
	return putURISigningKeys(tx, &r.cfg.AuthOptions, &r.cfg.Port, xmlID, keysJson)
}

func (r *Riak) DeleteURISigningKeys(xmlID string, tx *sql.Tx, ctx context.Context) error {
	return deleteURISigningKeys(tx, &r.cfg.AuthOptions, &r.cfg.Port, xmlID)
}

func (r *Riak) Ping(tx *sql.Tx, ctx context.Context) (tc.TrafficVaultPing, error) {
	resp, err := ping(tx, &r.cfg.AuthOptions, &r.cfg.Port)
	return tc.TrafficVaultPing(resp), err
}

func (r *Riak) GetBucketKey(bucket string, key string, tx *sql.Tx) ([]byte, bool, error) {
	return getBucketKey(tx, &r.cfg.AuthOptions, &r.cfg.Port, bucket, key)
}

func init() {
	trafficvault.AddBackend(RiakBackendName, riakConfigLoad)
}

func riakConfigLoad(b json.RawMessage) (trafficvault.TrafficVault, error) {
	riakCfg, err := unmarshalRiakConfig(b)
	if err != nil {
		return nil, err
	}
	if err := validateConfig(riakCfg); err != nil {
		return nil, errors.New("validating Riak config: " + err.Error())
	}
	return &Riak{cfg: riakCfg}, nil
}

func validateConfig(cfg Config) error {
	errs := tovalidate.ToErrors(validation.Errors{
		"user":     validation.Validate(cfg.User, validation.Required),
		"password": validation.Validate(cfg.Password, validation.Required),
	})
	if len(errs) == 0 {
		return nil
	}
	return util.JoinErrs(errs)
}
