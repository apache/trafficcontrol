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
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/auth"
)

type CRUDer interface {
	Creator
	Reader
	Updater
	Deleter
	Identifier
	Validator
}

type Updater interface {
	Update() (error, tc.ApiErrorType)
}

type Identifier interface {
	GetKeys() (map[string]interface{}, bool)
	GetType() string
	GetAuditName() string
	GetKeyFieldsInfo() []KeyFieldInfo
}

type Creator interface {
	Create() (error, tc.ApiErrorType)
	SetKeys(map[string]interface{})
}

type Deleter interface {
	Delete() (error, tc.ApiErrorType)
}

type Validator interface {
	Validate() error
}

type Tenantable interface {
	IsTenantAuthorized(user *auth.CurrentUser) (bool, error)
}

type Reader interface {
	Read(parameters map[string]string) ([]interface{}, []error, tc.ApiErrorType)
}
