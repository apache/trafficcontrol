package enum

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

import "strings"

type SteeringType string

const (
	SteeringTypeOrder     SteeringType = "STEERING_ORDER"
	SteeringTypeWeight    SteeringType = "STEERING_WEIGHT"
	SteeringTypeGeoOrder  SteeringType = "STEERING_GEO_ORDER"
	SteeringTypeGeoWeight SteeringType = "STEERING_GEO_WEIGHT"
	SteeringTypeInvalid   SteeringType = ""
)

func SteeringTypeFromString(s string) SteeringType {
	s = strings.ToLower(strings.Replace(s, "_", "", -1))
	switch s {
	case "steeringorder":
		return SteeringTypeOrder
	case "steeringweight":
		return SteeringTypeWeight
	case "steeringgeoorder":
		return SteeringTypeGeoOrder
	case "steeringgeoweight":
		return SteeringTypeGeoWeight
	default:
		return SteeringTypeInvalid
	}
}

// String returns a string representation of this steering type.
func (t SteeringType) String() string {
	switch t {
	case SteeringTypeOrder:
		fallthrough
	case SteeringTypeWeight:
		fallthrough
	case SteeringTypeGeoOrder:
		fallthrough
	case SteeringTypeGeoWeight:
		return string(t)
	default:
		return "INVALID"
	}
}

