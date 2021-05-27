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

	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-util"
)

// CdnLock is a struct to store the details of a lock that a user wishes to acquire on a CDN.
type CdnLock struct {
	UserName    string    `json:"userName" db:"username"`
	Cdn         string    `json:"cdn" db:"cdn"`
	Message     *string   `json:"message" db:"message"`
	Soft        *bool     `json:"soft" db:"soft"`
	LastUpdated TimeNoMod `json:"lastUpdated" db:"last_updated"`
}

type CdnLockCreateResponse struct {
	Response CdnLock `json:"response"`
	Alerts
}

type CdnLocksGetResponse struct {
	Response []CdnLock `json:"response"`
	Alerts
}

type CdnLockDeleteResponse CdnLockCreateResponse

func (c CdnLock) Validate(tx *sql.Tx) error {
	errs := validation.Errors{
		"cdn": validation.Validate(c.Cdn, validation.Required),
	}
	return util.JoinErrs(tovalidate.ToErrors(errs))
}
