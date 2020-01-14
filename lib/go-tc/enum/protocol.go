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
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const DSProtocolHTTP = 0
const DSProtocolHTTPS = 1
const DSProtocolHTTPAndHTTPS = 2
const DSProtocolHTTPToHTTPS = 3

// Protocol represents an ATC-supported content delivery protocol.
type Protocol string

const (
	// ProtocolHTTP represents the HTTP/1.1 protocol as specified in RFC2616.
	ProtocolHTTP = Protocol("http")
	// ProtocolHTTPS represents the HTTP/1.1 protocol over a TCP connection secured by TLS
	ProtocolHTTPS = Protocol("https")
	// ProtocolHTTPtoHTTPS represents a redirection of unsecured HTTP requests to HTTPS
	ProtocolHTTPtoHTTPS = Protocol("http to https")
	// ProtocolHTTPandHTTPS represents the use of both HTTP and HTTPS
	ProtocolHTTPandHTTPS = Protocol("http and https")
	// ProtocolInvalid represents an invalid Protocol
	ProtocolInvalid = Protocol("")
)

// String implements the "Stringer" interface.
func (p Protocol) String() string {
	switch p {
	case ProtocolHTTP:
		fallthrough
	case ProtocolHTTPS:
		fallthrough
	case ProtocolHTTPtoHTTPS:
		fallthrough
	case ProtocolHTTPandHTTPS:
		return string(p)
	default:
		return "INVALIDPROTOCOL"
	}
}

// ProtocolFromString parses a string and returns the corresponding Protocol.
func ProtocolFromString(s string) Protocol {
	switch strings.Replace(strings.ToLower(s), "_", " ", -1) {
	case "http":
		return ProtocolHTTP
	case "https":
		return ProtocolHTTPS
	case "http to https":
		return ProtocolHTTPtoHTTPS
	case "http and https":
		return ProtocolHTTPandHTTPS
	default:
		return ProtocolInvalid
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (p *Protocol) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return errors.New("Protocol cannot be null")
	}
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return fmt.Errorf("JSON %s not quoted: %v", data, err)
	}
	*p = ProtocolFromString(s)
	if *p == ProtocolInvalid {
		return fmt.Errorf("%s is not a (supported) Protocol", s)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (p Protocol) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}
