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

func (p *Postgres) GetDeliveryServiceSSLKeys(xmlID string, version string, tx *sql.Tx, ctx context.Context) (tc.DeliveryServiceSSLKeysV15, bool, error) {
	return tc.DeliveryServiceSSLKeysV15{}, false, notImplementedErr
}

func (p *Postgres) PutDeliveryServiceSSLKeys(key tc.DeliveryServiceSSLKeys, tx *sql.Tx, ctx context.Context) error {
	return notImplementedErr
}

func (p *Postgres) DeleteDeliveryServiceSSLKeys(xmlID string, version string, tx *sql.Tx, ctx context.Context) error {
	return notImplementedErr
}

func (p *Postgres) DeleteOldDeliveryServiceSSLKeys(existingXMLIDs map[string]struct{}, cdnName string, tx *sql.Tx, ctx context.Context) error {
	return notImplementedErr
}

func (p *Postgres) GetCDNSSLKeys(cdnName string, tx *sql.Tx, ctx context.Context) ([]tc.CDNSSLKey, error) {
	return nil, notImplementedErr
}

func (p *Postgres) GetDNSSECKeys(cdnName string, tx *sql.Tx, ctx context.Context) (tc.DNSSECKeysTrafficVault, bool, error) {
	return tc.DNSSECKeysTrafficVault{}, false, notImplementedErr
}

func (p *Postgres) PutDNSSECKeys(cdnName string, keys tc.DNSSECKeysTrafficVault, tx *sql.Tx, ctx context.Context) error {
	return notImplementedErr
}

func (p *Postgres) DeleteDNSSECKeys(cdnName string, tx *sql.Tx, ctx context.Context) error {
	return notImplementedErr
}

func (p *Postgres) GetURLSigKeys(xmlID string, tx *sql.Tx, ctx context.Context) (tc.URLSigKeys, bool, error) {
	tvTx, dbCtx, cancelFunc, err := p.beginTransaction(ctx)
	if err != nil {
		return tc.URLSigKeys{}, false, err
	}
	defer p.commitTransaction(tvTx, dbCtx, cancelFunc)
	return getURLSigKeys(xmlID, tvTx)
}

func (p *Postgres) PutURLSigKeys(xmlID string, keys tc.URLSigKeys, tx *sql.Tx, ctx context.Context) error {
	tvTx, dbCtx, cancelFunc, err := p.beginTransaction(ctx)
	if err != nil {
		return err
	}
	defer p.commitTransaction(tvTx, dbCtx, cancelFunc)

	return putURLSigKeys(xmlID, tvTx, keys)
}

func (p *Postgres) DeleteURLSigKeys(xmlID string, tx *sql.Tx, ctx context.Context) error {
	tvTx, dbCtx, cancelFunc, err := p.beginTransaction(ctx)
	if err != nil {
		return err
	}
	defer p.commitTransaction(tvTx, dbCtx, cancelFunc)

	return deleteURLSigKeys(xmlID, tvTx)
}

func (p *Postgres) GetURISigningKeys(xmlID string, tx *sql.Tx, ctx context.Context) ([]byte, bool, error) {
	return nil, false, notImplementedErr
}

func (p *Postgres) PutURISigningKeys(xmlID string, keysJson []byte, tx *sql.Tx, ctx context.Context) error {
	return notImplementedErr
}

func (p *Postgres) DeleteURISigningKeys(xmlID string, tx *sql.Tx, ctx context.Context) error {
	return notImplementedErr
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
