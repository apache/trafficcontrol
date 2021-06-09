package main

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
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/apache/trafficcontrol/lib/go-tc"
	util "github.com/apache/trafficcontrol/lib/go-util"

	_ "github.com/lib/pq"
)

// PGConfig represents the configuration options available to the PG backend
type PGConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	SSLMode  string `json:"sslmode"`
	Database string `json:"database"`
	Key      string `json:"aesKey"`
	AESKey   []byte
}

// PGBackend is the Postgres implementation of TVBackend
type PGBackend struct {
	sslKey pgSSLKeyTable
	dnssec pgDNSSecTable
	uri    pgURISignKeyTable
	url    pgURLSigKeyTable
	cfg    PGConfig
	db     *sql.DB
}

// String returns a high level overview of the backend and its keys
func (pg *PGBackend) String() string {
	data := fmt.Sprintf("PG server %v@%v:%v\n", pg.cfg.User, pg.cfg.Host, pg.cfg.Port)
	data += fmt.Sprintf("\tSSL Keys: %v\n", len(pg.sslKey.Records))
	data += fmt.Sprintf("\tDNSSec Keys: %v\n", len(pg.dnssec.Records))
	data += fmt.Sprintf("\tURI Keys: %v\n", len(pg.uri.Records))
	data += fmt.Sprintf("\tURL Keys: %v\n", len(pg.url.Records))
	return data
}

// Name returns the name for this backend
func (pg *PGBackend) Name() string {
	return "PG"
}

// ReadConfig takes in a filename and will read it into the backends config
func (pg *PGBackend) ReadConfig(s string) error {
	err := UnmarshalConfig(s, &pg.cfg)
	if err != nil {
		return err
	}

	pg.cfg.AESKey, err = base64.StdEncoding.DecodeString(pg.cfg.Key)
	if err != nil {
		return fmt.Errorf("unable to decode PG AESKey: %w", err)
	}
	return nil
}

// Insert takes the current keys and inserts them into the backend DB
func (pg *PGBackend) Insert() error {
	if err := pg.sslKey.insertKeys(pg.db); err != nil {
		return err
	}
	if err := pg.dnssec.insertKeys(pg.db); err != nil {
		return err
	}
	if err := pg.url.insertKeys(pg.db); err != nil {
		return err
	}
	if err := pg.uri.insertKeys(pg.db); err != nil {
		return err
	}
	return nil
}

// Start initiates the connection to the backend DB
func (pg *PGBackend) Start() error {
	sqlStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", pg.cfg.User, pg.cfg.Password, pg.cfg.Host, pg.cfg.Port, pg.cfg.Database, pg.cfg.SSLMode)
	db, err := sql.Open("postgres", sqlStr)
	if err != nil {
		return fmt.Errorf("unable to start PG client: %w", err)
	}

	pg.db = db
	pg.sslKey = pgSSLKeyTable{}
	pg.dnssec = pgDNSSecTable{}
	pg.url = pgURLSigKeyTable{}
	pg.uri = pgURISignKeyTable{}

	return nil
}

// ValidateKey validates that the keys are valid (in most cases, certain fields are not null)
func (pg *PGBackend) ValidateKey() []string {
	var errors []string
	if errs := pg.sslKey.validate(); errs != nil {
		errors = append(errors, errs...)
	}
	if errs := pg.dnssec.validate(); errs != nil {
		errors = append(errors, errs...)
	}
	if errs := pg.uri.validate(); errs != nil {
		errors = append(errors, errs...)
	}
	if errs := pg.url.validate(); errs != nil {
		errors = append(errors, errs...)
	}
	return errors
}

// Stop terminates the connection to the backend DB
func (pg *PGBackend) Stop() error {
	return pg.db.Close()
}

// Ping checks the connection to the backend DB
func (pg *PGBackend) Ping() error {
	return pg.db.Ping()
}

// Fetch gets all of the keys from the backend DB
func (pg *PGBackend) Fetch() error {
	if err := pg.sslKey.gatherKeys(pg.db); err != nil {
		return err
	}

	if err := pg.dnssec.gatherKeys(pg.db); err != nil {
		return err
	}

	if err := pg.url.gatherKeys(pg.db); err != nil {
		return err
	}

	if err := pg.uri.gatherKeys(pg.db); err != nil {
		return err
	}

	return nil
}

// GetSSLKeys converts the backends internal key representation into the common representation (SSLKey)
func (pg *PGBackend) GetSSLKeys() ([]SSLKey, error) {
	if err := pg.sslKey.decrypt(pg.cfg.AESKey); err != nil {
		return nil, err
	}
	return pg.sslKey.toGeneric(), nil
}

// SetSSLKeys takes in keys and converts & encrypts the data into the backends internal format
func (pg *PGBackend) SetSSLKeys(keys []SSLKey) error {
	pg.sslKey.fromGeneric(keys)
	return pg.sslKey.encrypt(pg.cfg.AESKey)
}

// GetDNSSecKeys converts the backends internal key representation into the common representation (DNSSecKey)
func (pg *PGBackend) GetDNSSecKeys() ([]DNSSecKey, error) {
	if err := pg.dnssec.decrypt(pg.cfg.AESKey); err != nil {
		return nil, err
	}
	return pg.dnssec.toGeneric(), nil
}

// SetDNSSecKeys takes in keys and converts & encrypts the data into the backends internal format
func (pg *PGBackend) SetDNSSecKeys(keys []DNSSecKey) error {
	pg.dnssec.fromGeneric(keys)
	return pg.dnssec.encrypt(pg.cfg.AESKey)
}

// GetURISignKeys converts the pg internal key representation into the common representation (URISignKey)
func (pg *PGBackend) GetURISignKeys() ([]URISignKey, error) {
	if err := pg.uri.decrypt(pg.cfg.AESKey); err != nil {
		return nil, err
	}
	return pg.uri.toGeneric(), nil
}

// SetURISignKeys takes in keys and converts & encrypts the data into the backends internal format
func (pg *PGBackend) SetURISignKeys(keys []URISignKey) error {
	pg.uri.fromGeneric(keys)
	return pg.uri.encrypt(pg.cfg.AESKey)
}

// GetURLSigKeys converts the backends internal key representation into the common representation (URLSigKey)
func (pg *PGBackend) GetURLSigKeys() ([]URLSigKey, error) {
	if err := pg.url.decrypt(pg.cfg.AESKey); err != nil {
		return nil, err
	}
	return pg.url.toGeneric(), nil
}

// SetURLSigKeys takes in keys and converts & encrypts the data into the backends internal format
func (pg *PGBackend) SetURLSigKeys(keys []URLSigKey) error {
	pg.url.fromGeneric(keys)
	return pg.url.encrypt(pg.cfg.AESKey)
}

type pgCommonRecord struct {
	DataEncrypted []byte
}

type pgDNSSecRecord struct {
	Key tc.DNSSECKeysTrafficVault
	CDN string
	pgCommonRecord
}
type pgDNSSecTable struct {
	Records []pgDNSSecRecord
}

func (tbl *pgDNSSecTable) gatherKeys(db *sql.DB) error {
	sz, err := getSize(db, "dnssec")
	if err != nil {
		log.Println("PGDNSSec gatherKeys: unable to determine size of dnssec table")
	}
	tbl.Records = make([]pgDNSSecRecord, sz)

	rows, err := db.Query("SELECT cdn, data from dnssec")
	if err != nil {
		return fmt.Errorf("PGDNSSec gatherKeys: unable to query: %w", err)
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		if i > len(tbl.Records)-1 {
			return fmt.Errorf("PGDNSSec gatherKeys got more results than expected %v", len(tbl.Records))
		}
		err := rows.Scan(&tbl.Records[i].CDN, &tbl.Records[i].DataEncrypted)
		if err != nil {
			return fmt.Errorf("PGDNSSec gatherKeys unable to scan row: %w", err)
		}
		i += 1
	}
	return nil
}
func (tbl *pgDNSSecTable) decrypt(aesKey []byte) error {
	for i, _ := range tbl.Records {
		err := decryptInto(aesKey, tbl.Records[i].DataEncrypted, &tbl.Records[i].Key)
		if err != nil {
			return fmt.Errorf("unable to decrypt into keys: %w", err)
		}
	}
	return nil
}
func (tbl *pgDNSSecTable) encrypt(aesKey []byte) error {
	for i, dns := range tbl.Records {
		data, err := json.Marshal(&dns.Key)
		if err != nil {
			return fmt.Errorf("encrypt issue marshalling keys: %w", err)
		}
		dat, err := encrypt(data, aesKey)
		if err != nil {
			return fmt.Errorf("encrypt error: %w", err)
		}
		tbl.Records[i].DataEncrypted = dat
	}
	return nil
}
func (tbl *pgDNSSecTable) toGeneric() []DNSSecKey {
	keys := make([]DNSSecKey, len(tbl.Records))

	for i, record := range tbl.Records {
		keys[i] = DNSSecKey{
			CDN:                    record.CDN,
			DNSSECKeysTrafficVault: record.Key,
		}
	}

	return keys
}
func (tbl *pgDNSSecTable) fromGeneric(keys []DNSSecKey) {
	tbl.Records = make([]pgDNSSecRecord, len(keys))

	for i, key := range keys {
		tbl.Records[i] = pgDNSSecRecord{
			Key: key.DNSSECKeysTrafficVault,
			CDN: key.CDN,
			pgCommonRecord: pgCommonRecord{
				DataEncrypted: nil,
			},
		}
	}
}
func (tbl *pgDNSSecTable) validate() []string {
	for i, record := range tbl.Records {
		if record.DataEncrypted == nil && len(record.Key) > 0 {
			return []string{fmt.Sprintf("DNSSEC Key %v: DataEncrypted is blank!", i)}
		}
	}
	return nil
}
func (tbl *pgDNSSecTable) insertKeys(db *sql.DB) error {
	queryFmt := "INSERT INTO dnssec (cdn, data) VALUES "
	stride := 2
	queryArgs := make([]interface{}, len(tbl.Records)*stride)
	for i, record := range tbl.Records {
		j := i * stride
		queryArgs[j] = record.CDN
		queryArgs[j+1] = record.DataEncrypted
	}
	return insertIntoTable(db, queryFmt, stride, queryArgs)
}

type pgSSLKeyRecord struct {
	Keys tc.DeliveryServiceSSLKeys
	pgCommonRecord

	// These records are stored on the table but are duplicated
	DeliveryService string
	CDN             string
}
type pgSSLKeyTable struct {
	Records []pgSSLKeyRecord
}

func (tbl *pgSSLKeyTable) insertKeys(db *sql.DB) error {
	queryFmt := "INSERT INTO sslkey (deliveryservice, data, cdn, version) VALUES "
	duplicateKeys := 2
	stride := 4
	queryArgs := make([]interface{}, len(tbl.Records)*stride*duplicateKeys)
	for i, record := range tbl.Records {
		j := i * duplicateKeys * stride

		queryArgs[j] = record.DeliveryService
		queryArgs[j+1] = record.DataEncrypted
		queryArgs[j+2] = record.CDN
		queryArgs[j+3] = record.Keys.Version.String()

		queryArgs[j+4] = record.DeliveryService
		queryArgs[j+5] = record.DataEncrypted
		queryArgs[j+6] = record.CDN
		queryArgs[j+7] = "latest"
	}
	return insertIntoTable(db, queryFmt, 4, queryArgs)
}
func (tbl *pgSSLKeyTable) gatherKeys(db *sql.DB) error {
	sz, err := getSize(db, "sslkey WHERE version='latest'")
	if err != nil {
		return fmt.Errorf("PGSSLKey gatherKeys unable to determine size of sslkey table: %w", err)
	}
	tbl.Records = make([]pgSSLKeyRecord, sz)

	rows, err := db.Query("SELECT data, deliveryservice, cdn from sslkey WHERE version = 'latest'")
	if err != nil {
		return fmt.Errorf("PGSSLKey gatherKeys unable to query: %w", err)
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		if i > len(tbl.Records)-1 {
			return errors.New("PGSSLKey gatherKeys: got more results than expected")
		}
		err := rows.Scan(&tbl.Records[i].DataEncrypted, &tbl.Records[i].DeliveryService, &tbl.Records[i].CDN)
		if err != nil {
			return fmt.Errorf("PGSSLKey gatherKeys unable to scan row: %w", err)
		}
		i += 1
	}
	return nil
}
func (tbl *pgSSLKeyTable) decrypt(aesKey []byte) error {
	for i, dns := range tbl.Records {
		err := decryptInto(aesKey, dns.DataEncrypted, &tbl.Records[i].Keys)
		if err != nil {
			return fmt.Errorf("unable to decrypt into keys: %w", err)
		}
	}
	return nil
}
func (tbl *pgSSLKeyTable) encrypt(aesKey []byte) error {
	for i, dns := range tbl.Records {
		data, err := json.Marshal(dns.Keys)
		if err != nil {
			return fmt.Errorf("encrypt issue marshalling keys: %w", err)
		}
		dat, err := encrypt(data, aesKey)
		if err != nil {
			return fmt.Errorf("encrypt error: %w", err)
		}
		tbl.Records[i].DataEncrypted = dat
	}
	return nil
}
func (tbl *pgSSLKeyTable) toGeneric() []SSLKey {
	keys := make([]SSLKey, len(tbl.Records))

	for i, record := range tbl.Records {
		keys[i] = SSLKey{
			DeliveryServiceSSLKeys: record.Keys,
		}
	}
	return keys
}
func (tbl *pgSSLKeyTable) fromGeneric(keys []SSLKey) {
	tbl.Records = make([]pgSSLKeyRecord, len(keys))

	for i, key := range keys {
		tbl.Records[i] = pgSSLKeyRecord{
			Keys: key.DeliveryServiceSSLKeys,
			pgCommonRecord: pgCommonRecord{
				DataEncrypted: nil,
			},
			DeliveryService: key.DeliveryService,
			CDN:             key.CDN,
		}
	}
}
func (tbl *pgSSLKeyTable) validate() []string {
	defaultKey := tc.DeliveryServiceSSLKeys{}
	var errors []string
	fmtStr := "SSL Key %v: %v"
	for i, record := range tbl.Records {
		if record.Keys == defaultKey {
			errors = append(errors, fmt.Sprintf(fmtStr, i, "DS SSL Keys are default!"))
		} else if record.Keys.Key == "" {
			errors = append(errors, fmt.Sprintf(fmtStr, i, "Key is blank!"))
		} else if record.Keys.CDN == "" {
			errors = append(errors, fmt.Sprintf(fmtStr, i, "CDN is blank!"))
		} else if record.Keys.DeliveryService == "" {
			errors = append(errors, fmt.Sprintf(fmtStr, i, "DS is blank!"))
		} else if record.DataEncrypted == nil {
			errors = append(errors, fmt.Sprintf(fmtStr, i, "DataEncrypted is blank!"))
		} else if record.Keys.Version.String() == "" {
			errors = append(errors, fmt.Sprintf(fmtStr, i, "Version is blank!"))
		}
	}
	return errors
}

type pgURLSigKeyRecord struct {
	Keys            tc.URLSigKeys
	DeliveryService string
	pgCommonRecord
}
type pgURLSigKeyTable struct {
	Records []pgURLSigKeyRecord
}

func (tbl *pgURLSigKeyTable) insertKeys(db *sql.DB) error {
	queryBase := "INSERT INTO url_sig_key (deliveryservice, data) VALUES "
	stride := 2
	queryArgs := make([]interface{}, len(tbl.Records)*stride)
	for i, record := range tbl.Records {
		j := i * stride
		queryArgs[j] = record.DeliveryService
		queryArgs[j+1] = record.DataEncrypted
	}
	return insertIntoTable(db, queryBase, stride, queryArgs)
}
func (tbl *pgURLSigKeyTable) gatherKeys(db *sql.DB) error {
	sz, err := getSize(db, "url_sig_key")
	if err != nil {
		log.Println("PGURLSigKey gatherKeys: unable to determine url_sig_key table size")
	}
	tbl.Records = make([]pgURLSigKeyRecord, sz)

	rows, err := db.Query("SELECT deliveryservice, data from url_sig_key")
	if err != nil {
		return fmt.Errorf("PGURLSigKey gatherKeys error creating query: %w", err)
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		if i > len(tbl.Records)-1 {
			return fmt.Errorf("PGURLSigKey gatherKeys: got more results than expected %v", len(tbl.Records))
		}
		err := rows.Scan(&tbl.Records[i].DeliveryService, &tbl.Records[i].DataEncrypted)
		if err != nil {
			return fmt.Errorf("PGURLSigKey gatherKeys: unable to scan row: %w", err)
		}
		i += 1
	}
	return nil
}
func (tbl *pgURLSigKeyTable) decrypt(aesKey []byte) error {
	for i, sig := range tbl.Records {
		err := decryptInto(aesKey, sig.DataEncrypted, &tbl.Records[i].Keys)
		if err != nil {
			return fmt.Errorf("unable to decrypt into keys: %w", err)
		}
	}
	return nil
}
func (tbl *pgURLSigKeyTable) encrypt(aesKey []byte) error {
	for i, sig := range tbl.Records {
		data, err := json.Marshal(&sig.Keys)
		if err != nil {
			return fmt.Errorf("encrypt issue marshalling keys: %w", err)
		}

		dat, err := encrypt(data, aesKey)
		if err != nil {
			return fmt.Errorf("encrypt error: %w", err)
		}
		tbl.Records[i].DataEncrypted = dat
	}
	return nil
}
func (tbl *pgURLSigKeyTable) toGeneric() []URLSigKey {
	keys := make([]URLSigKey, len(tbl.Records))

	for i, record := range tbl.Records {
		keys[i] = URLSigKey{
			DeliveryService: record.DeliveryService,
			URLSigKeys:      record.Keys,
		}
	}
	return keys
}
func (tbl *pgURLSigKeyTable) fromGeneric(keys []URLSigKey) {
	tbl.Records = make([]pgURLSigKeyRecord, len(keys))

	for i, key := range keys {
		tbl.Records[i] = pgURLSigKeyRecord{
			Keys:            key.URLSigKeys,
			DeliveryService: key.DeliveryService,
			pgCommonRecord: pgCommonRecord{
				DataEncrypted: nil,
			},
		}
	}
}
func (tbl *pgURLSigKeyTable) validate() []string {
	for i, record := range tbl.Records {
		if record.DataEncrypted == nil && len(record.Keys) > 0 {
			return []string{fmt.Sprintf("URl Sig Key %v: DataEncrypted is blank!", i)}
		}
	}
	return nil
}

type pgURISignKeyRecord struct {
	Keys            map[string]tc.URISignerKeyset
	DeliveryService string
	pgCommonRecord
}
type pgURISignKeyTable struct {
	Records []pgURISignKeyRecord
}

func (tbl *pgURISignKeyTable) insertKeys(db *sql.DB) error {
	queryFmt := "INSERT INTO uri_signing_key (deliveryservice, data) VALUES "
	stride := 2
	queryArgs := make([]interface{}, len(tbl.Records)*stride)
	for i, record := range tbl.Records {
		j := i * stride
		queryArgs[j] = record.DeliveryService
		queryArgs[j+1] = record.DataEncrypted
	}
	return insertIntoTable(db, queryFmt, stride, queryArgs)
}
func (tbl *pgURISignKeyTable) gatherKeys(db *sql.DB) error {
	sz, err := getSize(db, "uri_signing_key")
	if err != nil {
		log.Println("PGURISignKey gatherKeys: unable to determine size of uri_signing_key table")
	}
	tbl.Records = make([]pgURISignKeyRecord, sz)

	rows, err := db.Query("SELECT deliveryservice, data from uri_signing_key")
	if err != nil {
		return fmt.Errorf("PGURISignKey gatherKeys error while query: %w", err)
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		if i > len(tbl.Records)-1 {
			return fmt.Errorf("PGURISignKey gatherKeys: got more results than expected %v", len(tbl.Records))
		}
		err := rows.Scan(&tbl.Records[i].DeliveryService, &tbl.Records[i].DataEncrypted)
		if err != nil {
			return fmt.Errorf("PGURISignKey gatherKeys: unable to scan row: %w", err)
		}
		i += 1
	}
	return nil
}
func (tbl *pgURISignKeyTable) decrypt(aesKey []byte) error {
	for i, sign := range tbl.Records {
		err := decryptInto(aesKey, sign.DataEncrypted, &tbl.Records[i].Keys)
		if err != nil {
			return fmt.Errorf("unable to decrypt into keys: %w", err)
		}
	}
	return nil
}
func (tbl *pgURISignKeyTable) encrypt(aesKey []byte) error {
	for i, sign := range tbl.Records {
		data, err := json.Marshal(sign.Keys)
		if err != nil {
			return fmt.Errorf("encrypt issue marshalling keys: %w", err)
		}

		dat, err := encrypt(data, aesKey)
		if err != nil {
			return fmt.Errorf("encrypt error: %w", err)
		}
		tbl.Records[i].DataEncrypted = dat
	}
	return nil
}
func (tbl *pgURISignKeyTable) toGeneric() []URISignKey {
	keys := make([]URISignKey, len(tbl.Records))

	for i, record := range tbl.Records {
		keys[i] = URISignKey{
			DeliveryService: record.DeliveryService,
			Keys:            record.Keys,
		}
	}

	return keys
}
func (tbl *pgURISignKeyTable) fromGeneric(keys []URISignKey) {
	tbl.Records = make([]pgURISignKeyRecord, len(keys))

	for i, key := range keys {
		tbl.Records[i] = pgURISignKeyRecord{
			Keys:            key.Keys,
			DeliveryService: key.DeliveryService,
			pgCommonRecord: pgCommonRecord{
				DataEncrypted: nil,
			},
		}
	}
}
func (tbl *pgURISignKeyTable) validate() []string {
	for i, record := range tbl.Records {
		if record.DataEncrypted == nil && len(record.Keys) > 0 {
			return []string{fmt.Sprintf("URI Sign Key %v: DataEncrypted is blank!", i)}
		}
	}
	return nil
}

func getSize(db *sql.DB, table string) (int64, error) {
	rows, err := db.Query("SELECT COUNT(*) FROM " + table)
	if err != nil {
		return 0, err
	}
	var numRows int64
	if !rows.Next() {
		return 0, errors.New("no results returned for: " + table)
	}
	err = rows.Scan(&numRows)
	if err != nil {
		return 0, fmt.Errorf("error reading number of results for %v: %w", table, err)
	}
	return numRows, nil
}
func decrypt(record []byte, aesKey []byte) ([]byte, error) {
	unencrypted, err := util.AESDecrypt(record, aesKey)
	if err != nil {
		return nil, fmt.Errorf("unable to decrypt: %w", err)
	}
	return unencrypted, nil
}
func encrypt(record []byte, aesKey []byte) ([]byte, error) {
	encrypted, err := util.AESEncrypt(record, aesKey)
	if err != nil {
		return nil, err
	}
	return encrypted, nil
}
func decryptInto(aesKey []byte, encData []byte, value interface{}) error {
	data, err := decrypt(encData, aesKey)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &value)
	if err != nil {
		return err
	}
	return nil
}
func insertIntoTable(db *sql.DB, queryFmt string, stride int, queryArgs []interface{}) error {
	rows := len(queryArgs) / stride
	workStr := ""
	queryValueStr := make([]string, rows)
	for i, _ := range queryArgs {
		rowIndex := i % stride
		rowGroup := i / stride
		if rowIndex == 0 && i > 0 {
			queryValueStr[rowGroup-1] = "(" + workStr + ")"
			workStr = ""
		}
		if rowIndex == 0 {
			workStr += "$"
		} else {
			workStr += ",$"
		}
		workStr += strconv.Itoa(i + 1)
	}
	queryValueStr[len(queryValueStr)-1] = "(" + workStr + ")"
	query := queryFmt + strings.Join(queryValueStr, ",")

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("unable to open db transaction: %w", err)
	}
	res, err := tx.Exec(query, queryArgs...)
	if err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return fmt.Errorf("encountered error rolling back %v while handling error %v", err2, err)
		}
		return fmt.Errorf("error executing query '%v': %w", query, err)
	}
	if rows, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	} else if rows != int64(len(queryValueStr)) {
		return fmt.Errorf("wanted to insert %v rows, but inserted %v\n", len(queryValueStr), rows)
	}
	return tx.Commit()
}
