package api

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
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/auth"
	"github.com/jmoiron/sqlx"
)

type Updater interface {
	Update(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType)
	Identifier
	Validator
}

type Identifier interface {
	GetID() int
	GetType() string
	GetAuditName() string
}

type Inserter interface {
	Insert(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType)
	SetID(int)
	Identifier
	Validator
}

type Deleter interface {
	Delete(db *sqlx.DB, user auth.CurrentUser) (error, tc.ApiErrorType)
	SetID(int)
	Identifier
}

type Validator interface {
	Validate(db *sqlx.DB) []error
}

type Tenantable interface {
	IsTenantAuthorized(user auth.CurrentUser, db *sqlx.DB) (bool, error)
}
