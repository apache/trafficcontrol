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

// Steering holds information about a steering delivery service.
type Steering struct {
	DeliveryService DeliveryServiceName      `json:"deliveryService"`
	ClientSteering  bool                     `json:"clientSteering"`
	Targets         []SteeringSteeringTarget `json:"targets"`
	Filters         []SteeringFilter         `json:"filters"`
}

// SteeringResponse is the type of a response from Traffic Ops to a request
// to its /steering endpoint.
type SteeringResponse struct {
	Response []Steering `json:"response"`
	Alerts
}

// SteeringFilter is a filter for a target delivery service.
type SteeringFilter struct {
	DeliveryService DeliveryServiceName `json:"deliveryService"`
	Pattern         string              `json:"pattern"`
}

// SteeringSteeringTarget is a target delivery service of a steering delivery
// service.
type SteeringSteeringTarget struct {
	Order           int32               `json:"order"`
	Weight          int32               `json:"weight"`
	DeliveryService DeliveryServiceName `json:"deliveryService"`
	GeoOrder        *int                `json:"geoOrder,omitempty"`
	Longitude       *float64            `json:"longitude,omitempty"`
	Latitude        *float64            `json:"latitude,omitempty"`
}
