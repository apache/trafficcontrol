package tce

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
	"errors"
	"strconv"
	"strings"
)

// DeepCachingType represents a Delivery Service's deep caching type. The string values of this type should match the Traffic Ops values.
type DeepCachingType string

const (
	DeepCachingTypeNever   = DeepCachingType("") // default value
	DeepCachingTypeAlways  = DeepCachingType("ALWAYS")
	DeepCachingTypeInvalid = DeepCachingType("INVALID")
)

// String returns a string representation of this deep caching type
func (t DeepCachingType) String() string {
	switch t {
	case DeepCachingTypeAlways:
		return string(t)
	case DeepCachingTypeNever:
		return "NEVER"
	default:
		return "INVALID"
	}
}

// DeepCachingTypeFromString returns a DeepCachingType from its string representation, or DeepCachingTypeInvalid if the string is not a valid type.
func DeepCachingTypeFromString(s string) DeepCachingType {
	switch strings.ToLower(s) {
	case "always":
		return DeepCachingTypeAlways
	case "never":
		return DeepCachingTypeNever
	case "":
		// default when omitted
		return DeepCachingTypeNever
	default:
		return DeepCachingTypeInvalid
	}
}

// UnmarshalJSON unmarshals a JSON representation of a DeepCachingType (i.e. a string) or returns an error if the DeepCachingType is invalid
func (t *DeepCachingType) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*t = DeepCachingTypeNever
		return nil
	}
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return errors.New(string(data) + " JSON not quoted")
	}
	*t = DeepCachingTypeFromString(s)
	if *t == DeepCachingTypeInvalid {
		return errors.New(string(data) + " is not a DeepCachingType")
	}
	return nil
}

// MarshalJSON marshals into a JSON representation
func (t DeepCachingType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}
