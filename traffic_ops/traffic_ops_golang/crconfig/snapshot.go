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
	"encoding/json"
	"errors"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-log"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/traffic_ops/traffic_ops_golang/monitoring"
)

// Snapshot takes the CRConfig JSON-serializable object (which may be generated via crconfig.Make), and writes it to the snapshot table.
// It also takes the monitoring config JSON and writes it to the snapshot table.
func Snapshot(tx *sql.Tx, crc *tc.CRConfig, monitoringJSON *monitoring.Monitoring) error {
	log.Debugln("calling Snapshot")
	bts, err := json.Marshal(crc)
	if err != nil {
		return errors.New("marshalling JSON: " + err.Error())
	}
	date := time.Now()
	if crc.Stats.DateUnixSeconds != nil {
		date = time.Unix(*crc.Stats.DateUnixSeconds, 0)
	}

	btstm, err := json.Marshal(monitoringJSON)
	if err != nil {
		return errors.New("marshalling JSON: " + err.Error())
	}

	log.Debugf("calling Snapshot, writing %+v\n", date)
	q := `insert into snapshot (cdn, crconfig, last_updated, monitoring) values ($1, $2, $3, $4) on conflict(cdn) do update set crconfig=$2, last_updated=$3, monitoring=$4`
	if _, err := tx.Exec(q, crc.Stats.CDNName, bts, date, btstm); err != nil {
		return errors.New("Error inserting the crconfig and monitoring snapshot into database: " + err.Error())
	}
	return nil
}

// GetSnapshot gets the snapshot for the given CDN.
// If the CDN does not exist, false is returned.
// If the CDN exists, but the snapshot does not, the string for an empty JSON object "{}" is returned.
// An error is only returned on database error, never if the CDN or snapshot does not exist.
func GetSnapshot(tx *sql.Tx, cdn string) (string, bool, error) {
	log.Debugln("calling GetSnapshot")

	snapshot := sql.NullString{}
	// cdn left join snapshot, so we get a row with null if the CDN exists but the snapshot doesn't, and no rows if the CDN doesn't exist.
	q := `
SELECT s.crconfig AS snapshot
FROM cdn AS c
LEFT JOIN snapshot AS s ON s.cdn = c.name
WHERE c.name = $1
`
	if err := tx.QueryRow(q, cdn).Scan(&snapshot); err != nil {
		if err == sql.ErrNoRows {
			// CDN doesn't exist
			return "", false, nil
		}
		return "", false, errors.New("Error querying crconfig snapshot: " + err.Error())
	}
	if !snapshot.Valid {
		// CDN exists, but snapshot doesn't
		return `{}`, true, nil
	}
	return snapshot.String, true, nil
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
	if !monitorSnapshot.Valid || monitorSnapshot.String == "{}" {
		log.Errorln("Monitoring Snapshot didn't exist! Generating on-the-fly! This will cause race conditions in Traffic Monitor until a Snapshot is created!")
		monitoringJSON, err := monitoring.GetMonitoringJSON(tx, cdn)
		if err != nil {
			return "", false, errors.New("creating monitor snapshot (none existed): " + err.Error())
		}
		bts, err := json.Marshal(monitoringJSON)
		if err != nil {
			return "", false, errors.New("marshalling monitor snapshot (none existed): " + err.Error())
		}
		return string(bts), true, nil
	}
	return monitorSnapshot.String, true, nil
}
