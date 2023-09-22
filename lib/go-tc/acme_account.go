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

	"github.com/apache/trafficcontrol/v8/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/v8/lib/go-util"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// AcmeAccount is the information needed to access an account with an ACME provider.
type AcmeAccount struct {
	Email      *string `json:"email" db:"email"`
	PrivateKey *string `json:"privateKey" db:"private_key"`
	Uri        *string `json:"uri" db:"uri"`
	Provider   *string `json:"provider" db:"provider"`
}

// Validate validates the AcmeAccount request is valid for creation or update.
func (aa *AcmeAccount) Validate(tx *sql.Tx) error {
	errs := validation.Errors{
		"email":      validation.Validate(aa.Email, validation.Required, is.Email),
		"privateKey": validation.Validate(aa.PrivateKey, validation.Required),
		"uri":        validation.Validate(aa.Uri, validation.Required, is.URL),
		"provider":   validation.Validate(aa.Provider, validation.Required),
	}

	return util.JoinErrs(tovalidate.ToErrors(errs))
}
