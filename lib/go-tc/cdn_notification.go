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

	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-util"

	"github.com/go-ozzo/ozzo-validation"
)

// CDNNotification is a notification created for a specific CDN
type CDNNotification struct {
	CDN          *string    `json:"cdn" db:"cdn"`
	LastUpdated  *TimeNoMod `json:"lastUpdated,omitempty" db:"last_updated"`
	Notification *string    `json:"notification" db:"notification"`
	User         *string    `json:"user" db:"user"`
}

// Validate validates the CDNNotification request is valid for creation.
func (n *CDNNotification) Validate(tx *sql.Tx) error {
	errs := validation.Errors{
		"cdn":       validation.Validate(n.CDN, validation.Required),
	}
	return util.JoinErrs(tovalidate.ToErrors(errs))
}
