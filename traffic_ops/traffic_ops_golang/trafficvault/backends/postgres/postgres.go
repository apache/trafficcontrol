// Package postgres provides a TrafficVault implementation which uses PostgreSQL as the backend.
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
	"fmt"
	"strconv"
	"time"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/trafficvault"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	notImplementedErr = Error("this Traffic Vault functionality is not implemented for the postgres backend")

	postgresBackendName = "postgres"

	defaultMaxIdleConnections     = 10 // if this is higher than MaxDBConnections it will be automatically adjusted below it by the db/sql library
	defaultConnMaxLifetimeSeconds = 60
	defaultDBQueryTimeoutSecs     = 30

	latestVersion = "latest"
)

type Config struct {
	DBName                 string `json:"dbname"`
	Hostname               string `json:"hostname"`
	User                   string `json:"user"`
	Password               string `json:"password"`
	Port                   int    `json:"port"`
	SSL                    bool   `json:"ssl"`
	MaxConnections         int    `json:"max_connections"`
	MaxIdleConnections     int    `json:"max_idle_connections"`
	ConnMaxLifetimeSeconds int    `json:"conn_max_lifetime_seconds"`
	QueryTimeoutSeconds    int    `json:"query_timeout_seconds"`
}

type Postgres struct {
	cfg Config
	db  *sqlx.DB
}

func checkErrWithContext(prefix string, err error, ctxErr error) error {
	e := prefix + err.Error()
	if ctxErr != nil {
		e = fmt.Sprintf("%s: %s: %s", prefix, ctxErr.Error(), err.Error())
	}
	return errors.New(e)
}

func (p *Postgres) beginTransaction(ctx context.Context) (*sqlx.Tx, context.Context, context.CancelFunc, error) {
	dbCtx, cancelFunc := context.WithTimeout(ctx, time.Duration(p.cfg.QueryTimeoutSeconds)*time.Second)
	tx, err := p.db.BeginTxx(dbCtx, nil)
	if err != nil {
		e := checkErrWithContext("could not begin Traffic Vault PostgreSQL transaction", err, ctx.Err())
		cancelFunc()
		return nil, nil, nil, e
	}
	return tx, dbCtx, cancelFunc, nil
}

func (p *Postgres) commitTransaction(tx *sqlx.Tx, ctx context.Context, cancelFunc context.CancelFunc) {
	if err := tx.Commit(); err != nil {
		e := checkErrWithContext("Traffic Vault PostgreSQL: committing transaction", err, ctx.Err())
		log.Errorln(e)
	}
	cancelFunc()
}

// GetDeliveryServiceSSLKeys retrieves the SSL keys of the given version for
// the delivery service identified by the given xmlID. If version is empty,
// the implementation should return the latest version.
func (p *Postgres) GetDeliveryServiceSSLKeys(xmlID string, version string, tx *sql.Tx, ctx context.Context) (tc.DeliveryServiceSSLKeysV15, bool, error) {
	tvTx, dbCtx, cancelFunc, err := p.beginTransaction(ctx)
	if err != nil {
		return tc.DeliveryServiceSSLKeysV15{}, false, err
	}
	defer p.commitTransaction(tvTx, dbCtx, cancelFunc)
	var jsonKeys string
	query := "SELECT data FROM sslkey WHERE deliveryservice=$1 AND version=$2"
	if version == "" {
		version = "latest"
	}
	err = tvTx.QueryRow(query, xmlID, version).Scan(&jsonKeys)
	if err != nil {
		if err == sql.ErrNoRows {
			return tc.DeliveryServiceSSLKeysV15{}, false, nil
		}
		e := checkErrWithContext("Traffic Vault PostgreSQL: executing SELECT SSL Keys query", err, ctx.Err())
		return tc.DeliveryServiceSSLKeysV15{}, false, e
	}
	sslKey := tc.DeliveryServiceSSLKeysV15{}
	err = json.Unmarshal([]byte(jsonKeys), &sslKey)
	if err != nil {
		return tc.DeliveryServiceSSLKeysV15{}, false, errors.New("unmarshalling ssl keys: " + err.Error())
	}
	return sslKey, true, nil
}

// PutDeliveryServiceSSLKeys stores the given SSL keys for a delivery service.
func (p *Postgres) PutDeliveryServiceSSLKeys(key tc.DeliveryServiceSSLKeys, tx *sql.Tx, ctx context.Context) error {
	tvTx, dbCtx, cancelFunc, err := p.beginTransaction(ctx)
	if err != nil {
		return err
	}
	defer p.commitTransaction(tvTx, dbCtx, cancelFunc)
	keyJSON, err := json.Marshal(&key)
	if err != nil {
		return errors.New("marshalling keys: " + err.Error())
	}

	// delete the old ssl keys first
	oldVersions := []string{strconv.FormatInt(int64(key.Version), 10), latestVersion}
	_, err = tvTx.Exec("DELETE FROM sslkey WHERE deliveryservice=$1 and version=ANY($2)", key.DeliveryService, pq.Array(oldVersions))
	if err != nil {
		e := checkErrWithContext("Traffic Vault PostgreSQL: executing DELETE SSL Key query for INSERT", err, ctx.Err())
		return e
	}

	// insert the new ssl keys now
	res, err := tvTx.Exec("INSERT INTO sslkey (deliveryservice, data, cdn, version) VALUES ($1, $2, $3, $4), ($5, $6, $7, $8)", key.DeliveryService, keyJSON, key.CDN, strconv.FormatInt(int64(key.Version), 10), key.DeliveryService, keyJSON, key.CDN, latestVersion)
	if err != nil {
		e := checkErrWithContext("Traffic Vault PostgreSQL: executing INSERT SSL Key query", err, ctx.Err())
		return e
	}
	if rowsAffected, err := res.RowsAffected(); err != nil {
		return err
	} else if rowsAffected == 0 {
		return errors.New("SSL Key: no keys were inserted")
	}
	return nil
}

// DeleteDeliveryServiceSSLKeys removes the SSL keys of the given version (or latest
// if version is empty) for the delivery service identified by the given xmlID.
func (p *Postgres) DeleteDeliveryServiceSSLKeys(xmlID string, version string, tx *sql.Tx, ctx context.Context) error {
	tvTx, dbCtx, cancelFunc, err := p.beginTransaction(ctx)
	if err != nil {
		return err
	}
	defer p.commitTransaction(tvTx, dbCtx, cancelFunc)
	if version == "" {
		version = latestVersion
	}
	query := "DELETE FROM sslkey WHERE deliveryservice=$1 AND version=$2"
	_, err = tvTx.Exec(query, xmlID, version)
	if err != nil {
		e := checkErrWithContext("Traffic Vault PostgreSQL: executing DELETE SSL Key query", err, ctx.Err())
		return e
	}
	return nil
}

// DeleteOldDeliveryServiceSSLKeys takes a set of existingXMLIDs as input and will remove
// all SSL keys for delivery services in the CDN identified by the given cdnName that
// do not contain an xmlID in the given set of existingXMLIDs. This method is called
// during a snapshot operation in order to delete SSL keys for delivery services that
// no longer exist.
func (p *Postgres) DeleteOldDeliveryServiceSSLKeys(existingXMLIDs map[string]struct{}, cdnName string, tx *sql.Tx, ctx context.Context) error {
	var keys []string
	tvTx, dbCtx, cancelFunc, err := p.beginTransaction(ctx)
	if err != nil {
		return err
	}
	defer p.commitTransaction(tvTx, dbCtx, cancelFunc)

	if len(existingXMLIDs) == 0 {
		keys = append(keys, "")
	}
	for k, _ := range existingXMLIDs {
		keys = append(keys, k)
	}
	_, err = tvTx.Exec("DELETE FROM sslkey WHERE cdn=$1 AND deliveryservice <> ALL ($2)", cdnName, pq.Array(keys))
	if err != nil {
		e := checkErrWithContext("Traffic Vault PostgreSQL: executing DELETE OLD SSL Key query", err, ctx.Err())
		return e
	}
	return nil
}

// GetCDNSSLKeys retrieves all the SSL keys for delivery services in the CDN identified
// by the given cdnName.
func (p *Postgres) GetCDNSSLKeys(cdnName string, tx *sql.Tx, ctx context.Context) ([]tc.CDNSSLKey, error) {
	var keys []tc.CDNSSLKey
	var key tc.CDNSSLKey

	tvTx, dbCtx, cancelFunc, err := p.beginTransaction(ctx)
	if err != nil {
		return keys, err
	}
	defer p.commitTransaction(tvTx, dbCtx, cancelFunc)

	rows, err := tvTx.Query("SELECT data from sslkey WHERE cdn=$1 AND version=$2", cdnName, latestVersion)
	if err != nil {
		e := checkErrWithContext("Traffic Vault PostgreSQL: executing GET SSL Keys for CDN query", err, ctx.Err())
		return keys, e
	}
	defer rows.Close()
	for rows.Next() {
		jsonKey := ""
		if err := rows.Scan(&jsonKey); err != nil {
			e := checkErrWithContext("Traffic Vault PostgreSQL: scanning CDN SSL keys", err, ctx.Err())
			return keys, e
		}
		err = json.Unmarshal([]byte(jsonKey), &key)
		if err != nil {
			log.Errorf("couldn't unmarshal json key: %v", err)
			continue
		}
		keys = append(keys, key)
	}
	return keys, nil
}

func (p *Postgres) GetDNSSECKeys(cdnName string, tx *sql.Tx, ctx context.Context) (tc.DNSSECKeysTrafficVault, bool, error) {
	tvTx, dbCtx, cancelFunc, err := p.beginTransaction(ctx)
	if err != nil {
		return tc.DNSSECKeysTrafficVault{}, false, err
	}
	defer p.commitTransaction(tvTx, dbCtx, cancelFunc)
	var dnssecJSON string
	if err := tvTx.QueryRow("SELECT data FROM dnssec WHERE cdn = $1", cdnName).Scan(&dnssecJSON); err != nil {
		if err == sql.ErrNoRows {
			return tc.DNSSECKeysTrafficVault{}, false, nil
		}
		e := checkErrWithContext("Traffic Vault PostgreSQL: executing SELECT DNSSEC keys query", err, ctx.Err())
		return tc.DNSSECKeysTrafficVault{}, false, e
	}
	dnssecKeys := tc.DNSSECKeysTrafficVault{}
	if err := json.Unmarshal([]byte(dnssecJSON), &dnssecKeys); err != nil {
		return tc.DNSSECKeysTrafficVault{}, false, errors.New("unmarshalling DNSSEC keys: " + err.Error())
	}
	return dnssecKeys, true, nil
}

func (p *Postgres) PutDNSSECKeys(cdnName string, keys tc.DNSSECKeysTrafficVault, tx *sql.Tx, ctx context.Context) error {
	tvTx, dbCtx, cancelFunc, err := p.beginTransaction(ctx)
	if err != nil {
		return err
	}
	defer p.commitTransaction(tvTx, dbCtx, cancelFunc)
	dnssecJSON, err := json.Marshal(&keys)
	if err != nil {
		return errors.New("marshalling DNSSEC keys: " + err.Error())
	}
	_, err = tvTx.Exec("DELETE FROM dnssec WHERE cdn = $1", cdnName)
	if err != nil {
		e := checkErrWithContext("Traffic Vault PostgreSQL: executing DELETE DNSSEC keys query prior to INSERT", err, ctx.Err())
		return e
	}
	res, err := tvTx.Exec("INSERT INTO dnssec (cdn, data) VALUES ($1, $2)", cdnName, dnssecJSON)
	if err != nil {
		e := checkErrWithContext("Traffic Vault PostgreSQL: executing INSERT DNSSEC keys query", err, ctx.Err())
		return e
	}
	if rowsAffected, err := res.RowsAffected(); err != nil {
		return err
	} else if rowsAffected == 0 {
		return errors.New("Traffic Vault PostgreSQL: executing INSERT DNSSEC keys query: no rows were inserted")
	}
	return nil
}

func (p *Postgres) DeleteDNSSECKeys(cdnName string, tx *sql.Tx, ctx context.Context) error {
	tvTx, dbCtx, cancelFunc, err := p.beginTransaction(ctx)
	if err != nil {
		return err
	}
	defer p.commitTransaction(tvTx, dbCtx, cancelFunc)
	_, err = tvTx.Exec("DELETE FROM dnssec WHERE cdn = $1", cdnName)
	if err != nil {
		e := checkErrWithContext("Traffic Vault PostgreSQL: executing DELETE DNSSEC keys query", err, ctx.Err())
		return e
	}
	return nil
}

func (p *Postgres) GetURLSigKeys(xmlID string, tx *sql.Tx, ctx context.Context) (tc.URLSigKeys, bool, error) {
	tvTx, dbCtx, cancelFunc, err := p.beginTransaction(ctx)
	if err != nil {
		return tc.URLSigKeys{}, false, err
	}
	defer p.commitTransaction(tvTx, dbCtx, cancelFunc)
	return getURLSigKeys(xmlID, tvTx, ctx)
}

func (p *Postgres) PutURLSigKeys(xmlID string, keys tc.URLSigKeys, tx *sql.Tx, ctx context.Context) error {
	tvTx, dbCtx, cancelFunc, err := p.beginTransaction(ctx)
	if err != nil {
		return err
	}
	defer p.commitTransaction(tvTx, dbCtx, cancelFunc)

	return putURLSigKeys(xmlID, tvTx, keys, ctx)
}

func (p *Postgres) DeleteURLSigKeys(xmlID string, tx *sql.Tx, ctx context.Context) error {
	tvTx, dbCtx, cancelFunc, err := p.beginTransaction(ctx)
	if err != nil {
		return err
	}
	defer p.commitTransaction(tvTx, dbCtx, cancelFunc)

	return deleteURLSigKeys(xmlID, tvTx, ctx)
}

func (p *Postgres) GetURISigningKeys(xmlID string, tx *sql.Tx, ctx context.Context) ([]byte, bool, error) {
	tvTx, dbCtx, cancelFunc, err := p.beginTransaction(ctx)
	if err != nil {
		return []byte{}, false, err
	}
	defer p.commitTransaction(tvTx, dbCtx, cancelFunc)
	return getURISigningKeys(xmlID, tvTx, ctx)
}

func (p *Postgres) PutURISigningKeys(xmlID string, keysJson []byte, tx *sql.Tx, ctx context.Context) error {
	tvTx, dbCtx, cancelFunc, err := p.beginTransaction(ctx)
	if err != nil {
		return err
	}
	defer p.commitTransaction(tvTx, dbCtx, cancelFunc)

	return putURISigningKeys(xmlID, tvTx, keysJson, ctx)
}

func (p *Postgres) DeleteURISigningKeys(xmlID string, tx *sql.Tx, ctx context.Context) error {
	tvTx, dbCtx, cancelFunc, err := p.beginTransaction(ctx)
	if err != nil {
		return err
	}
	defer p.commitTransaction(tvTx, dbCtx, cancelFunc)

	return deleteURISigningKeys(xmlID, tvTx, ctx)
}

func (p *Postgres) Ping(tx *sql.Tx, ctx context.Context) (tc.TrafficVaultPing, error) {
	tvTx, dbCtx, cancelFunc, err := p.beginTransaction(ctx)
	if err != nil {
		return tc.TrafficVaultPing{}, err
	}
	defer p.commitTransaction(tvTx, dbCtx, cancelFunc)
	n := 0
	if err := tvTx.QueryRow("SELECT 1").Scan(&n); err != nil {
		e := checkErrWithContext("Traffic Vault PostgreSQL: executing ping query", err, dbCtx.Err())
		return tc.TrafficVaultPing{}, e
	}
	if n != 1 {
		return tc.TrafficVaultPing{}, fmt.Errorf("Traffic Vault PostgreSQL: executing ping query: expected scanned value 1 but got %d instead", n)
	}
	return tc.TrafficVaultPing{Status: "OK", Server: p.cfg.Hostname + ":" + strconv.Itoa(p.cfg.Port)}, nil
}

func (p *Postgres) GetBucketKey(bucket string, key string, tx *sql.Tx) ([]byte, bool, error) {
	return nil, false, notImplementedErr
}

func init() {
	trafficvault.AddBackend(postgresBackendName, postgresLoad)
}

func postgresLoad(b json.RawMessage) (trafficvault.TrafficVault, error) {
	pgCfg := Config{}
	if err := json.Unmarshal(b, &pgCfg); err != nil {
		return nil, errors.New("unmarshalling Postgres config: " + err.Error())
	}
	if err := validateConfig(pgCfg); err != nil {
		return nil, errors.New("validating Postgres config: " + err.Error())
	}
	if pgCfg.MaxIdleConnections == 0 {
		pgCfg.MaxIdleConnections = defaultMaxIdleConnections
	}
	if pgCfg.ConnMaxLifetimeSeconds == 0 {
		pgCfg.ConnMaxLifetimeSeconds = defaultConnMaxLifetimeSeconds
	}
	if pgCfg.QueryTimeoutSeconds == 0 {
		pgCfg.QueryTimeoutSeconds = defaultDBQueryTimeoutSecs
	}

	sslStr := "require"
	if !pgCfg.SSL {
		sslStr = "disable"
	}
	db, err := sqlx.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s&fallback_application_name=trafficvault", pgCfg.User, pgCfg.Password, pgCfg.Hostname, pgCfg.Port, pgCfg.DBName, sslStr))
	if err != nil {
		return nil, errors.New("opening database: " + err.Error())
	}
	db.SetMaxOpenConns(pgCfg.MaxConnections)
	db.SetMaxIdleConns(pgCfg.MaxIdleConnections)
	db.SetConnMaxLifetime(time.Duration(pgCfg.ConnMaxLifetimeSeconds) * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pgCfg.QueryTimeoutSeconds)*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		// NOTE: not fatal since Traffic Vault not being available at startup shouldn't be fatal
		log.Errorln("pinging the Traffic Vault database: " + err.Error())
	} else {
		log.Infoln("successfully pinged the Traffic Vault database")
	}

	return &Postgres{cfg: pgCfg, db: db}, nil
}

func validateConfig(cfg Config) error {
	errs := tovalidate.ToErrors(validation.Errors{
		"user":                  validation.Validate(cfg.User, validation.Required),
		"password":              validation.Validate(cfg.Password, validation.Required),
		"hostname":              validation.Validate(cfg.Hostname, validation.Required),
		"dbname":                validation.Validate(cfg.DBName, validation.Required),
		"port":                  validation.Validate(cfg.Port, validation.By(tovalidate.IsValidPortNumber)),
		"max_connections":       validation.Validate(cfg.MaxConnections, validation.Min(0)),
		"query_timeout_seconds": validation.Validate(cfg.QueryTimeoutSeconds, validation.Min(0)),
	})
	if len(errs) == 0 {
		return nil
	}
	return util.JoinErrs(errs)
}
