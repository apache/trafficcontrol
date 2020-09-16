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
