package rfc

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
	"fmt"
	"net/url"
	"strconv"
)

// URL is an alias of net/url.URL that implements JSON encoding and decoding, as well as scanning
// from database driver values.
type URL struct{ url.URL }

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
func (u *URL) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		u = nil
		return nil
	}
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return fmt.Errorf("JSON %s not quoted: %v", data, err)
	}
	addr, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("Couldn't parse '%s' as a URL: %v", s, err)
	}
	*u = URL{*addr}
	return nil
}

// MarshalJSON implements the encoding/json.Marshaler interface.
func (u URL) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

// Scan implements the database/sql.Scanner interface.
func (u *URL) Scan(src interface{}) error {
	if src == nil {
		u = nil
		return nil
	}

	switch src.(type) {
	case string:
		addr, err := url.Parse(src.(string))
		if err == nil {
			*u = URL{*addr}
		}
		return err
	case []byte:
		addr, err := url.Parse(string(src.([]byte)))
		if err == nil {
			*u = URL{*addr}
		}
		return err
	default:
		return fmt.Errorf("Type %T cannot represent a URL!", src)
	}
}
