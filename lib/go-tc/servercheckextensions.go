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
	"regexp"

	"github.com/apache/trafficcontrol/v8/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/v8/lib/go-util"
	validation "github.com/go-ozzo/ozzo-validation"
)

// ServerCheckExtensionNullable represents a server check extension used by Traffic Ops.
type ServerCheckExtensionNullable struct {
	ID                    *int       `json:"id" db:"id"`
	Name                  *string    `json:"name" db:"name"`
	Version               *string    `json:"version" db:"version"`
	InfoURL               *string    `json:"info_url" db:"info_url"`
	ScriptFile            *string    `json:"script_file" db:"script_file"`
	IsActive              *int       `json:"isactive" db:"isactive"`
	AdditionConfigJSON    *string    `json:"additional_config_json" db:"additional_config_json"`
	Description           *string    `json:"description" db:"description"`
	ServercheckShortName  *string    `json:"servercheck_short_name" db:"servercheck_short_name"`
	ServercheckColumnName *string    `json:"-" db:"servercheck_column_name"`
	Type                  *string    `json:"type" db:"type_name"`
	TypeID                *int       `json:"-" db:"type"`
	LastUpdated           *TimeNoMod `json:"-" db:"last_updated"`
}

// ServerCheckExtensionResponse represents the response from Traffic Ops when getting ServerCheckExtension.
type ServerCheckExtensionResponse struct {
	Response []ServerCheckExtensionNullable `json:"response"`
	Alerts
}

// ServerCheckExtensionPostResponse represents the response from Traffic Ops when creating ServerCheckExtension.
type ServerCheckExtensionPostResponse struct {
	Response ServerCheckExtensionID `json:"supplemental"`
	Alerts
}

// A ServerCheckExtensionID contains an identified for a Servercheck Extension.
type ServerCheckExtensionID struct {
	ID int `json:"id"`
}

// Validate ensures that the ServerCheckExtensionNullable request body is valid for creation.
func (e *ServerCheckExtensionNullable) Validate(tx *sql.Tx) error {
	checkRegexType := regexp.MustCompile(`^CHECK_EXTENSION_`)
	errs := tovalidate.ToErrors(validation.Errors{
		"name":        validation.Validate(e.Name, validation.NotNil),
		"version":     validation.Validate(e.Version, validation.NotNil),
		"info_url":    validation.Validate(e.InfoURL, validation.NotNil),
		"script_file": validation.Validate(e.ScriptFile, validation.NotNil),
		"type":        validation.Validate(e.Type, validation.NotNil, validation.Match(checkRegexType)),
		"isactive":    validation.Validate(e.IsActive, validation.NotNil),
	})
	if e.ID != nil {
		errs = append(errs, errors.New("ServerCheckExtension update not supported; delete and re-add."))
	}
	if e.IsActive != nil && !(*e.IsActive == 0 || *e.IsActive == 1) {
		errs = append(errs, errors.New("isactive can only be 0 or 1."))
	}
	if len(errs) > 0 {
		return util.JoinErrs(errs)
	}
	return nil
}
