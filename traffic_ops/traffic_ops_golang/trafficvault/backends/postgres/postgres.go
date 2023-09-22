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

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/deliveryservice"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/trafficvault"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
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

	defaultHashiCorpVaultLoginPath  = "/v1/auth/approle/login"
	defaultHashiCorpVaultTimeoutSec = 30

	latestVersion = "latest"
)

type Config struct {
	DBName                 string          `json:"dbname"`
	Hostname               string          `json:"hostname"`
	User                   string          `json:"user"`
	Password               string          `json:"password"`
	Port                   int             `json:"port"`
	SSL                    bool            `json:"ssl"`
	MaxConnections         int             `json:"max_connections"`
	MaxIdleConnections     int             `json:"max_idle_connections"`
	ConnMaxLifetimeSeconds int             `json:"conn_max_lifetime_seconds"`
	QueryTimeoutSeconds    int             `json:"query_timeout_seconds"`
	AesKeyLocation         string          `json:"aes_key_location"`
	HashiCorpVault         *HashiCorpVault `json:"hashicorp_vault"`
}

type HashiCorpVault struct {
	Address    string `json:"address"`
	RoleID     string `json:"role_id"`
	SecretID   string `json:"secret_id"`
	LoginPath  string `json:"login_path"`
	SecretPath string `json:"secret_path"`
	TimeoutSec int    `json:"timeout_sec"`
	Insecure   bool   `json:"insecure"`
}

type Postgres struct {
	cfg    Config
	db     *sqlx.DB
	aesKey []byte
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
	var encryptedSslKeys []byte
	query := "SELECT data FROM sslkey WHERE deliveryservice=$1 AND version=$2"
	if version == "" {
		version = "latest"
	}
	err = tvTx.QueryRow(query, xmlID, version).Scan(&encryptedSslKeys)
	if err != nil {
		if err == sql.ErrNoRows {
			return tc.DeliveryServiceSSLKeysV15{}, false, err
		}
		e := checkErrWithContext("Traffic Vault PostgreSQL: executing SELECT SSL Keys query", err, ctx.Err())
		return tc.DeliveryServiceSSLKeysV15{}, false, e
	}

	jsonKeys, err := util.AESDecrypt(encryptedSslKeys, p.aesKey)
	if err != nil {
		return tc.DeliveryServiceSSLKeysV15{}, false, err
	}

	sslKey := tc.DeliveryServiceSSLKeysV15{}
	err = json.Unmarshal([]byte(jsonKeys), &sslKey)
	if err != nil {
		return tc.DeliveryServiceSSLKeysV15{}, false, errors.New("unmarshalling ssl keys: " + err.Error())
	}
	return sslKey, true, nil
}

// GetExpirationInformation returns the expiration information for all SSL Keys.
func (p *Postgres) GetExpirationInformation(tx *sql.Tx, ctx context.Context, days int) ([]tc.SSLKeyExpirationInformation, error) {
	tvTx, dbCtx, cancelFunc, err := p.beginTransaction(ctx)
	if err != nil {
		return []tc.SSLKeyExpirationInformation{}, err
	}
	defer p.commitTransaction(tvTx, dbCtx, cancelFunc)

	fedMap := map[string]bool{}
	fedRows, err := tx.Query("SELECT DISTINCT(ds.xml_id) FROM federation_deliveryservice AS fd JOIN deliveryservice AS ds ON ds.id = fd.deliveryservice")
	if err != nil {
		return []tc.SSLKeyExpirationInformation{}, err
	}
	defer fedRows.Close()

	for fedRows.Next() {
		var fedString string
		if err = fedRows.Scan(&fedString); err != nil {
			return []tc.SSLKeyExpirationInformation{}, err
		}
		fedMap[fedString] = true
	}

	inactiveQuery := "SELECT xml_id FROM deliveryservice WHERE active = 'INACTIVE' OR active = 'PRIMED'"
	iaRows, err := tx.Query(inactiveQuery)
	if err != nil {
		return []tc.SSLKeyExpirationInformation{}, err
	}
	defer iaRows.Close()

	inactiveList := map[string]bool{}
	for iaRows.Next() {
		var inactiveXmlId string
		if err = iaRows.Scan(&inactiveXmlId); err != nil {
			return []tc.SSLKeyExpirationInformation{}, err
		}
		inactiveList[inactiveXmlId] = true
	}

	query := "SELECT deliveryservice, cdn, provider, expiration FROM sslkey WHERE version='latest'"

	if days != 0 {
		query = query + fmt.Sprintf(" AND expiration <= (now() + '%d days'::interval)", days)
	}

	expirationInfos := []tc.SSLKeyExpirationInformation{}

	rows, err := tvTx.Query(query)
	if err != nil {
		return []tc.SSLKeyExpirationInformation{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var expirationInfo tc.SSLKeyExpirationInformation
		if err = rows.Scan(&expirationInfo.DeliveryService, &expirationInfo.CDN, &expirationInfo.Provider, &expirationInfo.Expiration); err != nil {
			return []tc.SSLKeyExpirationInformation{}, err
		}
		if inactiveList[expirationInfo.DeliveryService] {
			continue
		}
		expirationInfo.Federated = fedMap[expirationInfo.DeliveryService]

		expirationInfos = append(expirationInfos, expirationInfo)
	}

	return expirationInfos, nil
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

	encryptedKey, err := util.AESEncrypt(keyJSON, p.aesKey)
	if err != nil {
		return fmt.Errorf("encrypting keys: %w", err)
	}

	err = deliveryservice.Base64DecodeCertificate(&key.Certificate)
	if err != nil {
		return fmt.Errorf("decoding SSL keys, %w", err)
	}
	expiration, _, err := deliveryservice.ParseExpirationAndSansFromCert([]byte(key.Certificate.Crt), key.Hostname)
	if err != nil {
		return fmt.Errorf("parsing expiration from certificate: %w", err)
	}

	// insert the new ssl keys now
	res, err := tvTx.Exec("INSERT INTO sslkey (deliveryservice, data, cdn, version, provider, expiration) VALUES ($1, $2, $3, $4, $5, $6), ($7, $8, $9, $10, $11, $12)", key.DeliveryService, encryptedKey, key.CDN, strconv.FormatInt(int64(key.Version), 10), key.AuthType, expiration, key.DeliveryService, encryptedKey, key.CDN, latestVersion, key.AuthType, expiration)
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
		encryptedSslKeys := []byte{}
		if err := rows.Scan(&encryptedSslKeys); err != nil {
			e := checkErrWithContext("Traffic Vault PostgreSQL: scanning CDN SSL keys", err, ctx.Err())
			return keys, e
		}

		jsonKey, err := util.AESDecrypt(encryptedSslKeys, p.aesKey)
		if err != nil {
			log.Errorf("couldn't decrypt key: %v", err)
			continue
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
	var encryptedDnssecKey []byte
	if err := tvTx.QueryRow("SELECT data FROM dnssec WHERE cdn = $1", cdnName).Scan(&encryptedDnssecKey); err != nil {
		if err == sql.ErrNoRows {
			return tc.DNSSECKeysTrafficVault{}, false, nil
		}
		e := checkErrWithContext("Traffic Vault PostgreSQL: executing SELECT DNSSEC keys query", err, ctx.Err())
		return tc.DNSSECKeysTrafficVault{}, false, e
	}

	dnssecJSON, err := util.AESDecrypt(encryptedDnssecKey, p.aesKey)
	if err != nil {
		return tc.DNSSECKeysTrafficVault{}, false, err
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

	encryptedKey, err := util.AESEncrypt(dnssecJSON, p.aesKey)
	if err != nil {
		return errors.New("encrypting keys: " + err.Error())
	}

	res, err := tvTx.Exec("INSERT INTO dnssec (cdn, data) VALUES ($1, $2)", cdnName, encryptedKey)
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
	return getURLSigKeys(xmlID, tvTx, ctx, p.aesKey)
}

func (p *Postgres) PutURLSigKeys(xmlID string, keys tc.URLSigKeys, tx *sql.Tx, ctx context.Context) error {
	tvTx, dbCtx, cancelFunc, err := p.beginTransaction(ctx)
	if err != nil {
		return err
	}
	defer p.commitTransaction(tvTx, dbCtx, cancelFunc)

	return putURLSigKeys(xmlID, tvTx, keys, ctx, p.aesKey)
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
	return getURISigningKeys(xmlID, tvTx, ctx, p.aesKey)
}

func (p *Postgres) PutURISigningKeys(xmlID string, keysJson []byte, tx *sql.Tx, ctx context.Context) error {
	tvTx, dbCtx, cancelFunc, err := p.beginTransaction(ctx)
	if err != nil {
		return err
	}
	defer p.commitTransaction(tvTx, dbCtx, cancelFunc)

	return putURISigningKeys(xmlID, tvTx, keysJson, ctx, p.aesKey)
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
	if pgCfg.HashiCorpVault != nil {
		if pgCfg.HashiCorpVault.LoginPath == "" {
			pgCfg.HashiCorpVault.LoginPath = defaultHashiCorpVaultLoginPath
		}
		if pgCfg.HashiCorpVault.TimeoutSec == 0 {
			pgCfg.HashiCorpVault.TimeoutSec = defaultHashiCorpVaultTimeoutSec
		}
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

	aesKey, err := readKey(pgCfg)
	if err != nil {
		return nil, err
	}

	return &Postgres{cfg: pgCfg, db: db, aesKey: aesKey}, nil
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
	aesKeyLocSet := cfg.AesKeyLocation != ""
	hashiCorpVaultSet := cfg.HashiCorpVault != nil && *cfg.HashiCorpVault != HashiCorpVault{}
	if aesKeyLocSet && hashiCorpVaultSet {
		errs = append(errs, errors.New("aes_key_location and hashicorp_vault cannot both be set"))
	} else if hashiCorpVaultSet {
		hashiErrs := tovalidate.ToErrors(validation.Errors{
			"address":     validation.Validate(cfg.HashiCorpVault.Address, validation.Required, is.URL),
			"role_id":     validation.Validate(cfg.HashiCorpVault.RoleID, validation.Required),
			"secret_id":   validation.Validate(cfg.HashiCorpVault.SecretID, validation.Required),
			"secret_path": validation.Validate(cfg.HashiCorpVault.SecretPath, validation.Required),
		})
		errs = append(errs, hashiErrs...)
	} else if !aesKeyLocSet {
		errs = append(errs, errors.New("one of either aes_key_location or hashicorp_vault is required"))
	}
	if len(errs) == 0 {
		return nil
	}
	return util.JoinErrs(errs)
}
