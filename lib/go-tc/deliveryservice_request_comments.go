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

// DeliveryServiceRequestCommentV5 is a Delivery Service Request Comment as it appears in version 5 of the
// Traffic Ops API - it always points to the highest minor version in APIv5.
type DeliveryServiceRequestCommentV5 DeliveryServiceRequestCommentV50

// DeliveryServiceRequestCommentV50 is a struct containing the fields for a delivery
// service request comment, for API version 5.0.
type DeliveryServiceRequestCommentV50 struct {
	AuthorID                 IDNoMod   `json:"authorId" db:"author_id"`
	Author                   string    `json:"author"`
	DeliveryServiceRequestID int       `json:"deliveryServiceRequestId" db:"deliveryservice_request_id"`
	ID                       int       `json:"id" db:"id"`
	LastUpdated              time.Time `json:"lastUpdated" db:"last_updated"`
	Value                    string    `json:"value" db:"value"`
	XMLID                    string    `json:"xmlId" db:"xml_id"`
}

// DeliveryServiceRequestCommentsResponseV5 is a Delivery Service Request Comment Response as it appears in version 5
// of the Traffic Ops API - it always points to the highest minor version in APIv5.
type DeliveryServiceRequestCommentsResponseV5 DeliveryServiceRequestCommentsResponseV50

// DeliveryServiceRequestCommentsResponseV50 is a list of
// DeliveryServiceRequestCommentsV5 as a response, for API version 5.0.
type DeliveryServiceRequestCommentsResponseV50 struct {
	Response []DeliveryServiceRequestCommentV5 `json:"response"`
	Alerts
}

// DeliveryServiceRequestCommentsResponse is a list of
// DeliveryServiceRequestComments as a response.
type DeliveryServiceRequestCommentsResponse struct {
	Response []DeliveryServiceRequestComment `json:"response"`
	Alerts
}

// DeliveryServiceRequestComment is a struct containing the fields for a delivery
// service request comment.
type DeliveryServiceRequestComment struct {
	AuthorID                 IDNoMod   `json:"authorId" db:"author_id"`
	Author                   string    `json:"author"`
	DeliveryServiceRequestID int       `json:"deliveryServiceRequestId" db:"deliveryservice_request_id"`
	ID                       int       `json:"id" db:"id"`
	LastUpdated              TimeNoMod `json:"lastUpdated" db:"last_updated"`
	Value                    string    `json:"value" db:"value"`
	XMLID                    string    `json:"xmlId" db:"xml_id"`
}

// DeliveryServiceRequestCommentNullable is a nullable struct containing the
// fields for a delivery service request comment.
type DeliveryServiceRequestCommentNullable struct {
	AuthorID                 *IDNoMod   `json:"authorId" db:"author_id"`
	Author                   *string    `json:"author"`
	DeliveryServiceRequestID *int       `json:"deliveryServiceRequestId" db:"deliveryservice_request_id"`
	ID                       *int       `json:"id" db:"id"`
	LastUpdated              *TimeNoMod `json:"lastUpdated" db:"last_updated"`
	Value                    *string    `json:"value" db:"value"`
	XMLID                    *string    `json:"xmlId" db:"xml_id"`
}
