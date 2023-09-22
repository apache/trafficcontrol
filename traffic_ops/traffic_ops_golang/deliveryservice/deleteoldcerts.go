package deliveryservice

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
	"sync"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/config"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/dbhelpers"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/trafficvault"
)

// DeleteOldCerts asynchronously deletes HTTPS certificates in Riak which have no corresponding delivery service in the database.
//
// Note the delivery service may still be in the CRConfig! Therefore, this should only be called immediately after a CRConfig Snapshot.
//
// This creates a goroutine, and immediately returns. It returns an error if there was an error preparing the delete routine, such as an error creating a db transaction.
//
// Note because it is asynchronous, this may return a nil error, but the asynchronous goroutine may error when fetching or deleting the certificates. If such an error occurs, it will be logged to the error log.
//
// If certificate deletion is already being processed by a goroutine, another delete will be queued, and this immediately returns nil. Only one delete will ever be queued.
func DeleteOldCerts(db *sql.DB, tx *sql.Tx, cfg *config.Config, cdn tc.CDNName, tv trafficvault.TrafficVault) error {
	if !cfg.TrafficVaultEnabled {
		log.Infoln("deleting old delivery service certificates: Traffic Vault is not enabled, returning without cleaning up old certificates.")
		return nil
	}
	if cfg.DisableAutoCertDeletion {
		log.Infoln("automatic certificate deletion is disabled, returning without cleaning up old certificates")
		return nil
	}
	if db == nil {
		return errors.New("nil db")
	}
	if cfg == nil {
		return errors.New("nil config")
	}
	startOldCertDeleter(db, tx, time.Duration(cfg.DBQueryTimeoutSeconds)*time.Second, cdn, tv)
	cleanupOldCertDeleters(tx)
	return nil
}

// deleteOldDSCerts deletes the HTTPS certificates in Riak of delivery services which have been deleted in Traffic Ops.
func deleteOldDSCerts(tx *sql.Tx, cdn tc.CDNName, tv trafficvault.TrafficVault) error {
	dses, err := dbhelpers.GetCDNDSes(tx, cdn)
	if err != nil {
		return errors.New("getting ds names: " + err.Error())
	}

	if err := tv.DeleteOldDeliveryServiceSSLKeys(dses, string(cdn), tx, context.Background()); err != nil {
		return errors.New("getting ds keys from Traffic Vault: " + err.Error())
	}
	return nil
}

// deleteOldDSCertsDB takes a db, and creates a transaction to pass to deleteOldDSCerts.
func deleteOldDSCertsDB(db *sql.DB, dbTimeout time.Duration, cdn tc.CDNName, tv trafficvault.TrafficVault) {
	dbCtx, cancelTx := context.WithTimeout(context.Background(), dbTimeout)
	defer cancelTx()
	tx, err := db.BeginTx(dbCtx, nil)
	if err != nil {
		log.Errorln("Old Cert Deleter Job: beginning tx: " + err.Error())
		return
	}
	txCommit := false
	defer dbhelpers.CommitIf(tx, &txCommit)
	if err := deleteOldDSCerts(tx, cdn, tv); err != nil {
		log.Errorln("deleting old DS certificates: " + err.Error())
		return
	}
	txCommit = true
}

var oldCertDeleters = OldCertDeleters{D: map[tc.CDNName]*OldCertDeleter{}}

type OldCertDeleters struct {
	D map[tc.CDNName]*OldCertDeleter
	M sync.Mutex
}

// startOldCertDeleter tells the old cert deleter goroutine to start another delete job, creating the goroutine if it doesn't exist.
func startOldCertDeleter(db *sql.DB, tx *sql.Tx, dbTimeout time.Duration, cdn tc.CDNName, tv trafficvault.TrafficVault) {
	oldCertDeleter := getOrCreateOldCertDeleter(cdn)
	oldCertDeleter.Once.Do(func() {
		go doOldCertDeleter(oldCertDeleter.Start, oldCertDeleter.Die, db, dbTimeout, cdn, tv)
	})

	select {
	case oldCertDeleter.Start <- struct{}{}:
	default:
	}
}

func getOrCreateOldCertDeleter(cdn tc.CDNName) *OldCertDeleter {
	oldCertDeleters.M.Lock()
	defer oldCertDeleters.M.Unlock()
	oldCertDeleter, ok := oldCertDeleters.D[cdn]
	if !ok {
		oldCertDeleter = newOldCertDeleter()
		oldCertDeleters.D[cdn] = oldCertDeleter
	}
	return oldCertDeleter
}

// cleanupOldCertDeleters stops all cert deleter goroutines for CDNs which have been deleted in the database.
// Any error is logged, but not returned.
// This is designed to be called when starting a new cert deleter job, to clean up any old cert deleters from deleted CDNs.
// This should only be called immediately when snapshotting, and immediately after startOldCertDeleter, because otherwise a cert deleter may be removed before it can delete old certs for a given CDN.
func cleanupOldCertDeleters(tx *sql.Tx) {
	cdns, err := dbhelpers.GetCDNs(tx) // (map[tc.CDNName]struct{}, error) {
	if err != nil {
		log.Errorln("cleanupOldCertDeleters: getting cdns: " + err.Error())
		return
	}

	oldCertDeleters.M.Lock()
	defer oldCertDeleters.M.Unlock()

	for cdn, oldCertDeleter := range oldCertDeleters.D {
		if _, ok := cdns[cdn]; ok {
			continue
		}
		delete(oldCertDeleters.D, cdn)
		select {
		case oldCertDeleter.Die <- struct{}{}:
		default:
		}
	}
}

type OldCertDeleter struct {
	Start chan struct{}
	Die   chan struct{}
	Once  sync.Once
}

func newOldCertDeleter() *OldCertDeleter {
	return &OldCertDeleter{
		Start: make(chan struct{}, 1),
		Die:   make(chan struct{}, 1),
	}
}

func doOldCertDeleter(do chan struct{}, die chan struct{}, db *sql.DB, dbTimeout time.Duration, cdn tc.CDNName, tv trafficvault.TrafficVault) {
	for {
		select {
		case <-do:
			deleteOldDSCertsDB(db, dbTimeout, cdn, tv)
		case <-die:
			// Go selects aren't ordered, so double-check the do chan in case a race happened and a job came in at the same time as the die.
			select {
			case <-do:
				deleteOldDSCertsDB(db, dbTimeout, cdn, tv)
			default:
			}
			return
		}
	}
}
