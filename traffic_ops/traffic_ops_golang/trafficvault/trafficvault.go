// Package trafficvault provides the interfaces and types necessary to support various
// Traffic Vault backend data stores.
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
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

// TrafficVault defines the methods necessary for a struct to implement in order to
// provide all the necessary functionality required of a Traffic Vault backend.
type TrafficVault interface {
	// NOTE: the ctx context.Context in these methods is for the HTTP request context in order to cancel the request
	// if the HTTP connection is closed. If the method is called asynchronously in a goroutine that is spawned while
	// handling the original HTTP request, you should use context.Background() so that the context isn't cancelled
	// when the original HTTP connection is closed.

	// GetDeliveryServiceSSLKeys retrieves the SSL keys of the given version for
	// the delivery service identified by the given xmlID. If version is empty,
	// the implementation should return the latest version.
	GetDeliveryServiceSSLKeys(xmlID string, version string, tx *sql.Tx, ctx context.Context) (tc.DeliveryServiceSSLKeysV15, bool, error)
	// GetExpirationInformation retrieves the SSL key expiration information for all delivery services.
	GetExpirationInformation(tx *sql.Tx, ctx context.Context, days int) ([]tc.SSLKeyExpirationInformation, error)
	// PutDeliveryServiceSSLKeys stores the given SSL keys for a delivery service.
	PutDeliveryServiceSSLKeys(key tc.DeliveryServiceSSLKeys, tx *sql.Tx, ctx context.Context) error
	// DeleteDeliveryServiceSSLKeys removes the SSL keys of the given version (or latest
	// if version is empty) for the delivery service identified by the given xmlID.
	DeleteDeliveryServiceSSLKeys(xmlID string, version string, tx *sql.Tx, ctx context.Context) error
	// DeleteOldDeliveryServiceSSLKeys takes a set of existingXMLIDs as input and will remove
	// all SSL keys for delivery services in the CDN identified by the given cdnName that
	// do not contain an xmlID in the given set of existingXMLIDs. This method is called
	// during a snapshot operation in order to delete SSL keys for delivery services that
	// no longer exist.
	DeleteOldDeliveryServiceSSLKeys(existingXMLIDs map[string]struct{}, cdnName string, tx *sql.Tx, ctx context.Context) error
	// GetCDNSSLKeys retrieves all the SSL keys for delivery services in the CDN identified
	// by the given cdnName.
	GetCDNSSLKeys(cdnName string, tx *sql.Tx, ctx context.Context) ([]tc.CDNSSLKey, error)
	// GetDNSSECKeys retrieves all the DNSSEC keys associated with the CDN identified by the
	// given cdnName.
	GetDNSSECKeys(cdnName string, tx *sql.Tx, ctx context.Context) (tc.DNSSECKeysTrafficVault, bool, error)
	// PutDNSSECKeys stores all the DNSSEC keys for the CDN identified by the given cdnName.
	PutDNSSECKeys(cdnName string, keys tc.DNSSECKeysTrafficVault, tx *sql.Tx, ctx context.Context) error
	// DeleteDNSSECKeys removes all the DNSSEC keys for the CDN identified by the given cdnName.
	DeleteDNSSECKeys(cdnName string, tx *sql.Tx, ctx context.Context) error
	// GetURLSigKeys retrieves the URL sig keys for the delivery service identified by the
	// given xmlID.
	GetURLSigKeys(xmlID string, tx *sql.Tx, ctx context.Context) (tc.URLSigKeys, bool, error)
	// PutURLSigKeys stores the given URL sig keys for the delivery service identified by
	// the given xmlID.
	PutURLSigKeys(xmlID string, keys tc.URLSigKeys, tx *sql.Tx, ctx context.Context) error
	// DeleteURLSigKeys deletes the URL sig keys for the delivery service identified
	// by the given xmlID.
	DeleteURLSigKeys(xmlID string, tx *sql.Tx, ctx context.Context) error
	// GetURISigningKeys retrieves the URI signing keys (as raw JSON bytes) for the delivery
	// service identified by the given xmlID.
	GetURISigningKeys(xmlID string, tx *sql.Tx, ctx context.Context) ([]byte, bool, error)
	// PutURISigningKeys stores the given URI signing keys (as raw JSON bytes) for the delivery
	// service identified by the given xmlID.
	PutURISigningKeys(xmlID string, keysJson []byte, tx *sql.Tx, ctx context.Context) error
	// DeleteURISigningKeys removes the URI signing keys for the delivery service identified by
	// the given xmlID.
	DeleteURISigningKeys(xmlID string, tx *sql.Tx, ctx context.Context) error
	// Ping simply checks the health of the Traffic Vault backend, returning a status and which
	// server hostname the status was returned by.
	Ping(tx *sql.Tx, ctx context.Context) (tc.TrafficVaultPing, error)
	// GetBucketKey returns the raw bytes identified by the given bucket and key. This may not
	// apply to every Traffic Vault backend implementation.
	// Deprecated: this method and associated API routes will be removed in the future.
	GetBucketKey(bucket string, key string, tx *sql.Tx) ([]byte, bool, error)
}

var backends = make(map[string]LoadFunc)

// A LoadFunc is a function that takes a json.RawMessage as input (the contents of
// traffic_vault_config in cdn.conf) and returns a valid TrafficVault as output. Each
// TrafficVault implementation should define its own LoadFunc which is responsible for
// parsing the given configuration and returning a valid TrafficVault implementation that
// may be used by request handlers.
type LoadFunc func(json.RawMessage) (TrafficVault, error)

// AddBackend should be called by each TrafficVault backend package's init() function in order
// to register its name and LoadFunc. This name corresponds to the traffic_vault_backend option
// in cdn.conf.
func AddBackend(name string, loadConfig LoadFunc) {
	backends[name] = loadConfig
}

// GetBackend is called with the contents of the traffic_vault_backend and traffic_vault_config
// options in cdn.conf, respectively, in order to lookup and load the chosen Traffic Vault
// backend to use.
func GetBackend(name string, cfgJson json.RawMessage) (TrafficVault, error) {
	loader, ok := backends[name]
	if !ok {
		return nil, fmt.Errorf("no supported Traffic Vault backend named '%s' was found", name)
	}
	backend, err := loader(cfgJson)
	if err != nil {
		return nil, fmt.Errorf("failed to load backend '%s': %s", name, err.Error())
	}
	return backend, nil
}
