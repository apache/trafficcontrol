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
	"net/mail"
	"strconv"
)

// EmailAddress is an alias of net/mail.Address that implements JSON encoding and decoding, as well
// as scanning from database driver values.
type EmailAddress struct{ mail.Address }

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
func (a *EmailAddress) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		a = nil
		return nil
	}
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return fmt.Errorf("JSON %s not quoted: %v", data, err)
	}
	addr, err := mail.ParseAddress(s)
	if err != nil {
		return fmt.Errorf("Couldn't parse '%s' as an email address: %v", s, err)
	}
	*a = EmailAddress{*addr}
	return nil
}

// MarshalJSON implements the encoding/json.Marshaler interface.
func (a EmailAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

// Scan implements the database/sql.Scanner interface.
func (a *EmailAddress) Scan(src interface{}) error {
	if src == nil {
		a = nil
		return nil
	}

	switch src.(type) {
	case string:
		addr, err := mail.ParseAddress(src.(string))
		if err == nil {
			*a = EmailAddress{*addr}
		}
		return err
	case []byte:
		addr, err := mail.ParseAddress(string(src.([]byte)))
		if err == nil {
			*a = EmailAddress{*addr}
		}
		return err
	default:
		return fmt.Errorf("Type %T cannot represent an EmailAddress!", src)
	}
}
