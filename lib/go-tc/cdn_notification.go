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
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/v8/lib/go-util"

	"github.com/go-ozzo/ozzo-validation"
)

// CDNNotificationsResponse is a list of CDN notifications as a response.
type CDNNotificationsResponse struct {
	Response []CDNNotification `json:"response"`
	Alerts
}

// CDNNotificationRequest encodes the request data for the POST
// cdn_notifications endpoint.
type CDNNotificationRequest struct {
	CDN          string `json:"cdn"`
	Notification string `json:"notification"`
}

// CDNNotification is a notification created for a specific CDN.
type CDNNotification struct {
	ID           int       `json:"id" db:"id"`
	CDN          string    `json:"cdn" db:"cdn"`
	LastUpdated  time.Time `json:"lastUpdated" db:"last_updated"`
	Notification string    `json:"notification" db:"notification"`
	User         string    `json:"user" db:"user"`
}

// Validate validates the CDNNotificationRequest request is valid for creation.
func (n *CDNNotificationRequest) Validate(tx *sql.Tx) error {
	errs := validation.Errors{
		"cdn":          validation.Validate(n.CDN, validation.Required),
		"notification": validation.Validate(n.Notification, validation.Required),
	}
	return util.JoinErrs(tovalidate.ToErrors(errs))
}
