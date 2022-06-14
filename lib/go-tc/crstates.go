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
	"encoding/json"
	"time"
)

// CRStates includes availability data for caches and delivery services, as gathered and aggregated by this Traffic Monitor. It is designed to be served at an API endpoint primarily for Traffic Routers (Content Router) to consume.
type CRStates struct {
	Caches          map[CacheName]IsAvailable                       `json:"caches"`
	DeliveryService map[DeliveryServiceName]CRStatesDeliveryService `json:"deliveryServices"`
}

// CRStatesDeliveryService contains data about the availability of a particular delivery service, and which caches in that delivery service have been marked as unavailable.
type CRStatesDeliveryService struct {
	DisabledLocations []CacheGroupName `json:"disabledLocations"`
	IsAvailable       bool             `json:"isAvailable"`
}

// IsAvailable contains whether the given cache or delivery service is available. It is designed for JSON serialization, namely in the Traffic Monitor 1.0 API.
type IsAvailable struct {
	IsAvailable    bool      `json:"isAvailable"`
	Ipv4Available  bool      `json:"ipv4Available"`
	Ipv6Available  bool      `json:"ipv6Available"`
	DirectlyPolled bool      `json:"-"`
	Status         string    `json:"status"`
	LastPoll       time.Time `json:"lastPoll"`
}

// NewCRStates creates a new CR states object, initializing pointer members.
func NewCRStates(cacheCap, dsCap int) CRStates {
	return CRStates{
		Caches:          make(map[CacheName]IsAvailable, cacheCap),
		DeliveryService: make(map[DeliveryServiceName]CRStatesDeliveryService, dsCap),
	}
}

// Copy creates a deep copy of this object. It does not mutate, and is thus safe for multiple goroutines.
func (a CRStates) Copy() CRStates {
	b := NewCRStates(len(a.Caches), len(a.DeliveryService))
	for k, v := range a.Caches {
		b.Caches[k] = v
	}
	for k, v := range a.DeliveryService {
		b.DeliveryService[k] = v
	}
	return b
}

// CopyDeliveryServices creates a deep copy of the delivery service availability data. It does not mutate, and is thus safe for multiple goroutines.
func (a CRStates) CopyDeliveryServices() map[DeliveryServiceName]CRStatesDeliveryService {
	b := make(map[DeliveryServiceName]CRStatesDeliveryService, len(a.DeliveryService))
	for k, v := range a.DeliveryService {
		b[k] = v
	}
	return b
}

// CopyCaches creates a deep copy of the cache availability data. It does not mutate, and is thus safe for multiple goroutines.
func (a CRStates) CopyCaches() map[CacheName]IsAvailable {
	b := make(map[CacheName]IsAvailable, len(a.Caches))
	for k, v := range a.Caches {
		b[k] = v
	}
	return b
}

// CRStatesMarshall serializes the given CRStates into bytes.
func CRStatesMarshall(states CRStates) ([]byte, error) {
	return json.Marshal(states)
}

// CRStatesUnMarshall takes bytes of a JSON string, and unmarshals them into a CRStates object.
func CRStatesUnMarshall(body []byte) (CRStates, error) {
	var crStates CRStates
	err := json.Unmarshal(body, &crStates)
	return crStates, err
}
