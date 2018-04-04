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

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
)

// Snapshot takes the CRConfig JSON-serializable object (which may be generated via crconfig.Make), and writes it to the snapshot table.
func Snapshot(db *sql.DB, crc *tc.CRConfig) error {
	log.Errorln("DEBUG calling Snapshot")
	bts, err := json.Marshal(crc)
	if err != nil {
		return errors.New("marshalling JSON: " + err.Error())
	}
	date := time.Now()
	if crc.Stats.DateUnixSeconds != nil {
		date = time.Unix(*crc.Stats.DateUnixSeconds, 0)
	}
	log.Errorf("DEBUG calling Snapshot, writing %+v\n", date)
	q := `insert into snapshot (cdn, content, last_updated) values ($1, $2, $3) on conflict(cdn) do update set content=$2, last_updated=$3`
	if _, err := db.Exec(q, crc.Stats.CDNName, bts, date); err != nil {
		return errors.New("Error inserting the snapshot into database: " + err.Error())
	}
	return nil
}

// GetSnapshot gets the snapshot for the given CDN.
// If the CDN does not exist, false is returned.
// If the CDN exists, but the snapshot does not, the string for an empty JSON object "{}" is returned.
// An error is only returned on database error, never if the CDN or snapshot does not exist.
func GetSnapshot(db *sql.DB, cdn string) (string, bool, error) {
	log.Errorln("DEBUG calling GetSnapshot")

	snapshot := sql.NullString{}
	// cdn left join snapshot, so we get a row with null if the CDN exists but the snapshot doesn't, and no rows if the CDN doesn't exist.
	q := `
select s.content as snapshot
from cdn as c
left join snapshot as s on s.cdn = c.name
where c.name = $1
`
	if err := db.QueryRow(q, cdn).Scan(&snapshot); err != nil {
		if err == sql.ErrNoRows {
			// CDN doesn't exist
			return "", false, nil
		}
		return "", false, errors.New("Error querying snapshot: " + err.Error())
	}
	if !snapshot.Valid {
		// CDN exists, but snapshot doesn't
		return `{}`, true, nil
	}
	return snapshot.String, true, nil
}
