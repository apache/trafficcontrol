package crconfig

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
	"errors"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

// Snapshot creates a new snapshot of the CDN. Returns whether the CDN exists, and any error creating the snapshot.
// Note this does not include the "deliveryservices" key and "contentServers/{server}/deliveryservices" keys. Delivery Services must be snapshotted separately.
func Snapshot(tx *sql.Tx, cdn tc.CDNName) error {
	if err := UpdateSnapshotTables(tx); err != nil {
		return errors.New("updating snapshot tables: " + err.Error())
	}

	// TODO remove crconfig column, or at least not-null.
	qry := `
INSERT INTO snapshot (cdn, crconfig, time) VALUES ($1, '', now())
ON CONFLICT(cdn) DO UPDATE SET time=now()
`
	if _, err := tx.Exec(qry, cdn); err != nil {
		return errors.New("inserting the cdn snapshot into database: " + err.Error())
	}
	return nil
}

// GetSnapshotMonitoring gets the monitor snapshot for the given CDN.
// If the CDN does not exist, false is returned.
// If the CDN exists, but the snapshot does not, the string for an empty JSON object "{}" is returned.
// An error is only returned on database error, never if the CDN or snapshot does not exist.
// Because all snapshotting is handled by the crconfig endpoints we have to also do the monitoring one
// here as well
func GetSnapshotMonitoring(tx *sql.Tx, cdn string) (string, bool, error) {
	log.Debugln("calling GetSnapshotMonitoring")

	monitorSnapshot := sql.NullString{}
	// cdn left join snapshot, so we get a row with null if the CDN exists but the snapshot doesn't, and no rows if the CDN doesn't exist.
	q := `
SELECT s.monitoring AS snapshot
FROM cdn AS c
LEFT JOIN snapshot AS s ON s.cdn = c.name
WHERE c.name = $1
`
	if err := tx.QueryRow(q, cdn).Scan(&monitorSnapshot); err != nil {
		if err == sql.ErrNoRows {
			// CDN doesn't exist
			return "", false, nil
		}
		return "", false, errors.New("Error querying monitor snapshot: " + err.Error())
	}
	if !monitorSnapshot.Valid {
		// CDN exists, but snapshot doesn't
		return `{}`, true, nil
	}
	return monitorSnapshot.String, true, nil
}

func SnapshotDS(tx *sql.Tx, ds tc.DeliveryServiceName) error {
	if err := UpdateSnapshotTablesForDS(tx, ds); err != nil {
		return errors.New("updating snapshot tables for ds '" + string(ds) + "': " + err.Error())
	}

	qry := `
INSERT INTO deliveryservice_snapshots (deliveryservice, time) VALUES ($1, now())
ON CONFLICT(deliveryservice) DO UPDATE SET time=now()
`
	if _, err := tx.Exec(qry, ds); err != nil {
		return errors.New("inserting the deliveryservice snapshot into database: " + err.Error())
	}
	return nil
}
