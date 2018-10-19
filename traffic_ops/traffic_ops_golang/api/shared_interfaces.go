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
	// Update returns any user error, any system error, and the HTTP error code to be returned if there was an error.
	Update() (error, error, int)
}

type Identifier interface {
	GetKeys() (map[string]interface{}, bool)
	GetType() string
	GetAuditName() string
	GetKeyFieldsInfo() []KeyFieldInfo
}

type Creator interface {
	// Create returns any user error, any system error, and the HTTP error code to be returned if there was an error.
	Create() (error, error, int)
	SetKeys(map[string]interface{})
}

type Deleter interface {
	// Delete returns any user error, any system error, and the HTTP error code to be returned if there was an error.
	Delete() (error, error, int)
}

type Validator interface {
	Validate() error
}

type Tenantable interface {
	IsTenantAuthorized(user *auth.CurrentUser) (bool, error)
}

type Reader interface {
	// Read returns the object to write to the user, any user error, any system error, and the HTTP error code to be returned if there was an error.
	Read() ([]interface{}, error, error, int)
}
