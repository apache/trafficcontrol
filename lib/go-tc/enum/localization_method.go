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

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// LocalizationMethod represents an enabled localization method for a cachegroup. The string values of this type should match the Traffic Ops values.
type LocalizationMethod string

const (
	LocalizationMethodCZ      = LocalizationMethod("CZ")
	LocalizationMethodDeepCZ  = LocalizationMethod("DEEP_CZ")
	LocalizationMethodGeo     = LocalizationMethod("GEO")
	LocalizationMethodInvalid = LocalizationMethod("INVALID")
)

// String returns a string representation of this localization method
func (m LocalizationMethod) String() string {
	switch m {
	case LocalizationMethodCZ:
		return string(m)
	case LocalizationMethodDeepCZ:
		return string(m)
	case LocalizationMethodGeo:
		return string(m)
	default:
		return "INVALID"
	}
}

func LocalizationMethodFromString(s string) LocalizationMethod {
	switch strings.ToLower(s) {
	case "cz":
		return LocalizationMethodCZ
	case "deep_cz":
		return LocalizationMethodDeepCZ
	case "geo":
		return LocalizationMethodGeo
	default:
		return LocalizationMethodInvalid
	}
}

func (m *LocalizationMethod) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return errors.New("LocalizationMethod cannot be null")
	}
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return errors.New(string(data) + " JSON not quoted")
	}
	*m = LocalizationMethodFromString(s)
	if *m == LocalizationMethodInvalid {
		return errors.New(s + " is not a LocalizationMethod")
	}
	return nil
}

func (m LocalizationMethod) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}

func (m *LocalizationMethod) Scan(value interface{}) error {
	if value == nil {
		return errors.New("LocalizationMethod cannot be null")
	}
	sv, err := driver.String.ConvertValue(value)
	if err != nil {
		return errors.New("failed to scan LocalizationMethod: " + err.Error())
	}

	switch v := sv.(type) {
	case []byte:
		*m = LocalizationMethodFromString(string(v))
		if *m == LocalizationMethodInvalid {
			return errors.New(string(v) + " is not a valid LocalizationMethod")
		}
		return nil
	default:
		return fmt.Errorf("failed to scan LocalizationMethod, unsupported input type: %T", value)
	}
}
