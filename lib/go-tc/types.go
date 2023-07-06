package tc

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
	"time"
)

// TypesResponse is the type of a response from Traffic Ops to a GET request
// made to its /types API endpoint.
type TypesResponse struct {
	Response []Type `json:"response"`
	Alerts
}

// Type contains information about a given Type in Traffic Ops.
type Type struct {
	ID          int       `json:"id"`
	LastUpdated TimeNoMod `json:"lastUpdated"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UseInTable  string    `json:"useInTable"`
}

// TypeNullable contains information about a given Type in Traffic Ops.
type TypeNullable struct {
	ID          *int       `json:"id" db:"id"`
	LastUpdated *TimeNoMod `json:"lastUpdated" db:"last_updated"`
	Name        *string    `json:"name" db:"name"`
	Description *string    `json:"description" db:"description"`
	UseInTable  *string    `json:"useInTable" db:"use_in_table"`
}

// TypesResponseV5 is type struct response used for the latest minor version associated with api major version 5.
type TypesResponseV5 = TypesResponseV50

// TypesResponseV50 is the type of response (for RFC3339) from Traffic Ops to a GET Request
// made to its /types API endpoint.
type TypesResponseV50 struct {
	Response []TypeV50 `json:"response"`
	Alerts
}

// TypeV5 contains information about a given Type in Traffic Ops used for the latest minor version associated with api major version 5.
type TypeV5 = TypeV50

// TypeV50 contains information about a given Type in Traffic Ops.
type TypeV50 struct {
	ID          int       `json:"id"`
	LastUpdated time.Time `json:"lastUpdated"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UseInTable  string    `json:"useInTable"`
}

// GetTypeData returns the type's name and use_in_table, true/false if the
// query returned data, and any error.
//
// TODO: Move this to the DB helpers package.
func GetTypeData(tx *sql.Tx, id int) (string, string, bool, error) {
	name := ""
	var useInTablePtr *string
	if err := tx.QueryRow(`SELECT name, use_in_table from type where id=$1`, id).Scan(&name, &useInTablePtr); err != nil {
		if err == sql.ErrNoRows {
			return "", "", false, nil
		}
		return "", "", false, errors.New("querying type data: " + err.Error())
	}
	useInTable := ""
	if useInTablePtr != nil {
		useInTable = *useInTablePtr
	}
	return name, useInTable, true, nil
}

// ValidateTypeID validates that the typeID references a type with the expected
// use_in_table string and returns "" and an error if the typeID is invalid. If
// valid, the type's name is returned.
//
// TODO: Move this to the DB helpers package.
func ValidateTypeID(tx *sql.Tx, typeID *int, expectedUseInTable string) (string, error) {
	if typeID == nil {
		return "", errors.New("missing property: 'typeId'")
	}

	typeName, useInTable, ok, err := GetTypeData(tx, *typeID)
	if err != nil {
		return "", errors.New("validating type: " + err.Error())
	}
	if !ok {
		return "", errors.New("type not found")
	}
	if useInTable != expectedUseInTable {
		return "", errors.New("type is not a valid " + expectedUseInTable + " type")
	}
	return typeName, nil
}
