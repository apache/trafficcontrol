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
	"time"
)

// CDNLock is a struct to store the details of a lock that a user wishes to acquire on a CDN.
type CDNLock struct {
	UserName        string    `json:"userName" db:"username"`
	CDN             string    `json:"cdn" db:"cdn"`
	Message         *string   `json:"message" db:"message"`
	Soft            *bool     `json:"soft" db:"soft"`
	SharedUserNames []string  `json:"sharedUserNames" db:"shared_usernames"`
	LastUpdated     time.Time `json:"lastUpdated" db:"last_updated"`
}

// CDNLockCreateResponse is a struct to store the response of a CREATE operation on a lock.
type CDNLockCreateResponse struct {
	Response CDNLock `json:"response"`
	Alerts
}

// CDNLocksGetResponse is a struct to store the response of a GET operation on locks.
type CDNLocksGetResponse struct {
	Response []CDNLock `json:"response"`
	Alerts
}

// CDNLockDeleteResponse is a struct to store the response of a DELETE operation on a lock.
type CDNLockDeleteResponse CDNLockCreateResponse
