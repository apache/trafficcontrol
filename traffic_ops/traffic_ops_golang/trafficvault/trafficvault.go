package trafficvault

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
	"fmt"

	"github.com/apache/trafficcontrol/lib/go-tc"
)

type TrafficVault interface {
	GetDeliveryServiceSSLKeys(xmlID string, version string, tx *sql.Tx) (tc.DeliveryServiceSSLKeysV15, bool, error)
	PutDeliveryServiceSSLKeys(key tc.DeliveryServiceSSLKeys, tx *sql.Tx) error
	DeleteDeliveryServiceSSLKeys(xmlID string, version string, tx *sql.Tx) error
	DeleteOldDeliveryServiceSSLKeys(existingXMLIDs map[string]struct{}, cdnName string, tx *sql.Tx) error
	GetCDNSSLKeys(cdnName string, tx *sql.Tx) ([]tc.CDNSSLKey, error)
	GetDNSSECKeys(cdnName string, tx *sql.Tx) (tc.DNSSECKeysTrafficVault, bool, error)
	PutDNSSECKeys(cdnName string, keys tc.DNSSECKeysTrafficVault, tx *sql.Tx) error
	DeleteDNSSECKeys(cdnName string, tx *sql.Tx) error
	GetURLSigKeys(xmlID string, tx *sql.Tx) (tc.URLSigKeys, bool, error)
	PutURLSigKeys(xmlID string, keys tc.URLSigKeys, tx *sql.Tx) error
	GetURISigningKeys(xmlID string, tx *sql.Tx) ([]byte, bool, error)
	PutURISigningKeys(xmlID string, keysJson []byte, tx *sql.Tx) error
	DeleteURISigningKeys(xmlID string, tx *sql.Tx) error
	Ping(tx *sql.Tx) (tc.TrafficVaultPingResponse, error)
	GetBucketKey(bucket string, key string, tx *sql.Tx) ([]byte, bool, error)
}
type Error string

func (e Error) Error() string {
	return string(e)
}

const disabledErr = Error("traffic vault is disabled")

type Disabled struct {
}

func (d *Disabled) GetDeliveryServiceSSLKeys(xmlID string, version string, tx *sql.Tx) (tc.DeliveryServiceSSLKeysV15, bool, error) {
	return tc.DeliveryServiceSSLKeysV15{}, false, disabledErr
}

func (d *Disabled) PutDeliveryServiceSSLKeys(key tc.DeliveryServiceSSLKeys, tx *sql.Tx) error {
	return disabledErr
}

func (d *Disabled) DeleteDeliveryServiceSSLKeys(xmlID string, version string, tx *sql.Tx) error {
	return disabledErr
}

func (d *Disabled) DeleteOldDeliveryServiceSSLKeys(existingXMLIDs map[string]struct{}, cdnName string, tx *sql.Tx) error {
	return disabledErr
}

func (d *Disabled) GetCDNSSLKeys(cdnName string, tx *sql.Tx) ([]tc.CDNSSLKey, error) {
	return nil, disabledErr
}

func (d *Disabled) GetDNSSECKeys(cdnName string, tx *sql.Tx) (tc.DNSSECKeysTrafficVault, bool, error) {
	return nil, false, disabledErr
}

func (d *Disabled) PutDNSSECKeys(cdnName string, keys tc.DNSSECKeysTrafficVault, tx *sql.Tx) error {
	return disabledErr
}

func (d *Disabled) DeleteDNSSECKeys(cdnName string, tx *sql.Tx) error {
	return disabledErr
}

func (d *Disabled) GetURLSigKeys(xmlID string, tx *sql.Tx) (tc.URLSigKeys, bool, error) {
	return nil, false, disabledErr
}

func (d *Disabled) PutURLSigKeys(xmlID string, keys tc.URLSigKeys, tx *sql.Tx) error {
	return disabledErr
}

func (d *Disabled) GetURISigningKeys(xmlID string, tx *sql.Tx) ([]byte, bool, error) {
	return nil, false, disabledErr
}

func (d *Disabled) PutURISigningKeys(xmlID string, keysJson []byte, tx *sql.Tx) error {
	return disabledErr
}

func (d *Disabled) DeleteURISigningKeys(xmlID string, tx *sql.Tx) error {
	return disabledErr
}
func (d *Disabled) Ping(tx *sql.Tx) (tc.TrafficVaultPingResponse, error) {
	return tc.TrafficVaultPingResponse{}, disabledErr
}

func (d *Disabled) GetBucketKey(bucket string, key string, tx *sql.Tx) ([]byte, bool, error) {
	return nil, false, disabledErr
}

var backends = make(map[string]LoadFunc)

type LoadFunc func(json.RawMessage) (TrafficVault, error)

func AddBackend(name string, loadConfig LoadFunc) {
	backends[name] = loadConfig
}

func GetBackend(name string, cfgJson json.RawMessage) (TrafficVault, error) {
	loader, ok := backends[name]
	if !ok {
		return nil, fmt.Errorf("no support Traffic Vault backend named '%s' was found", name)
	}
	backend, err := loader(cfgJson)
	if err != nil {
		return nil, fmt.Errorf("failed to load backend '%s': %s", name, err.Error())
	}
	return backend, nil
}
