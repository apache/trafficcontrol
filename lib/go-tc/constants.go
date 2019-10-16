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

type ErrorConstant string

func (e ErrorConstant) Error() string { return string(e) }

const DBError = ErrorConstant("database access error")
const NilTenantError = ErrorConstant("tenancy is enabled but request tenantID is nil")
const TenantUserNotAuthError = ErrorConstant("user not authorized for requested tenant")
const TenantDSUserNotAuthError = ErrorConstant("user not authorized for requested delivery service tenant")

const ApplicationJson = "application/json"
const Gzip = "gzip"
const ContentType = "Content-Type"
const ContentEncoding = "Content-Encoding"
const ContentTypeTextPlain = "text/plain"

type AlertLevel int

const (
	SuccessLevel AlertLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

var alertLevels = [4]string{"success", "info", "warning", "error"}

func (a AlertLevel) String() string {
	return alertLevels[a]
}

const CachegroupCoordinateNamePrefix = "from_cachegroup_"
