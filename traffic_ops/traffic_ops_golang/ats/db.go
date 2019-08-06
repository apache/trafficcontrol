package ats

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

	"github.com/apache/trafficcontrol/lib/go-atscfg"
	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
)

func GetProfileParamData(tx *sql.Tx, profileID int, configFile string) (map[string]string, error) {
	// TODO add another func to return a slice, for things that don't need a map, for performance? Does it make a difference?
	qry := `
SELECT
  p.name,
  p.value
FROM
  parameter p
  join profile_parameter pp on p.id = pp.parameter
  JOIN profile pr on pr.id = pp.profile
WHERE
  pr.id = $1
  AND p.config_file = $2
`
	rows, err := tx.Query(qry, profileID, configFile)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	params := map[string]string{}
	for rows.Next() {
		name := ""
		val := ""
		if err := rows.Scan(&name, &val); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		if name == "location" {
			continue
		}

		if _, ok := params[name]; ok {
			log.Warnf("Profile %v has multiple parameters '%v' assigned! ATS config generation ignoring value '%v'!", profileID, name, params[name])
		}

		params[name] = val
	}
	return params, nil
}

type ProfileData struct {
	ID   int
	Name string
}

// GetProfileData returns the necessary info about the profile, whether it exists, and any error.
func GetProfileData(tx *sql.Tx, id int) (ProfileData, bool, error) {
	// TODO implement, determine what fields are necessary
	qry := `
SELECT
  p.name
FROM
  profile p
WHERE
  p.id = $1
`
	v := ProfileData{ID: id}
	if err := tx.QueryRow(qry, id).Scan(&v.Name); err != nil {
		if err == sql.ErrNoRows {
			return ProfileData{}, false, nil
		}
		return ProfileData{}, false, errors.New("querying: " + err.Error())
	}
	return v, true, nil
}

func GetProfileDS(tx *sql.Tx, profileID int) ([]atscfg.ProfileDS, error) {
	qry := `
SELECT
  dstype.name AS ds_type,
  (SELECT o.protocol::text || '://' || o.fqdn || rtrim(concat(':', o.port::text), ':')
    FROM origin o
    WHERE o.deliveryservice = ds.id
    AND o.is_primary) as org_server_fqdn
FROM
  deliveryservice ds
  JOIN type as dstype ON ds.type = dstype.id
WHERE
  ds.id IN (
    SELECT DISTINCT deliveryservice
    FROM deliveryservice_server
    WHERE server IN (SELECT id FROM server WHERE profile = $1)
  )p
`
	rows, err := tx.Query(qry, profileID)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	dses := []atscfg.ProfileDS{}
	for rows.Next() {
		d := atscfg.ProfileDS{}
		if err := rows.Scan(&d.Type, &d.OriginFQDN); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		d.Type = tc.DSTypeFromString(string(d.Type))
		dses = append(dses, d)
	}
	return dses, nil
}

// GetProfileParamValue gets the value of a parameter assigned to a profile, by name and config file.
// Returns the parameter, whether it existed, and any error.
func GetProfileParamValue(tx *sql.Tx, profileID int, configFile string, name string) (string, bool, error) {
	qry := `
SELECT
  p.value
FROM
  parameter p
  JOIN profile_parameter pp ON p.id = pp.parameter
WHERE
  pp.profile = $1
  AND p.config_file = $2
  AND p.name = $3
`
	val := ""
	if err := tx.QueryRow(qry, profileID, configFile, name).Scan(&val); err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, errors.New("querying: " + err.Error())
	}
	return val, true, nil
}

// GetProfileIDFromName returns the profile's ID, whether it exists, and any error.
func GetProfileIDFromName(tx *sql.Tx, profileName string) (int, bool, error) {
	qry := `SELECT id from profile where name = $1`
	id := 0
	if err := tx.QueryRow(qry, profileName).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, false, nil
		}
		return 0, false, errors.New("querying: " + err.Error())
	}
	return id, true, nil
}

type Parameter struct {
	Name       string
	ConfigFile string
	Value      string
}

func GetParamsByName(tx *sql.Tx, paramName string) ([]Parameter, error) {
	// TODO implement, determine what fields are necessary
	qry := `
SELECT
  p.value,
  p.config_file
FROM
  parameter p
WHERE
  p.name = $1
`
	rows, err := tx.Query(qry, paramName)
	if err != nil {
		return nil, errors.New("querying: " + err.Error())
	}
	defer rows.Close()

	params := []Parameter{}
	for rows.Next() {
		pa := Parameter{Name: paramName}
		if err := rows.Scan(&pa.Value, &pa.ConfigFile); err != nil {
			return nil, errors.New("scanning: " + err.Error())
		}
		params = append(params, pa)
	}
	return params, nil
}
