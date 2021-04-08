package riaksvc

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
	"crypto/tls"
	"testing"
)

const goodRiakConfig = `
	   {
	       "user": "riakuser",
	       "password": "password",
	       "port": 8087,
	       "MaxTLSVersion": "1.1",
	       "tlsConfig": {
	           "insecureSkipVerify": true
	       }
	   }
`

func TestUnmarshalGoodRiakConfig(t *testing.T) {
	if cfg, err := unmarshalRiakConfig([]byte(goodRiakConfig)); err != nil {
		t.Errorf("unmarshalling good Riak config - expected: nil error, actual: %v", err)
	} else {
		if cfg.User != "riakuser" {
			t.Errorf("unmarshalling good Riak config - expected user: riakuser, actual: %s", cfg.User)
		}
		if cfg.Password != "password" {
			t.Errorf("unmarshalling good Riak config - expected password: password, actual: %s", cfg.Password)
		}
		if cfg.Port != 8087 {
			t.Errorf("unmarshalling good Riak config - expected port: 8087, actual: %d", cfg.Port)
		}
		if cfg.TlsConfig == nil {
			t.Fatal("unmarshalling good Riak config - expected TlsConfig: non-nil, actual: nil")
		}
		if cfg.TlsConfig.InsecureSkipVerify != true {
			t.Errorf("unmarshalling good Riak config - expected TlsConfig.InsecureSkipVerify: true, actual: %t", cfg.TlsConfig.InsecureSkipVerify)
		}
		if cfg.TlsConfig.MaxVersion != tls.VersionTLS11 {
			t.Errorf("unmarshalling good Riak config - expected TlsConfig.MaxVersion: %d, actual: %d", tls.VersionTLS11, cfg.TlsConfig.MaxVersion)
		}
	}
}

func TestUnmarshalBadRiakConfig(t *testing.T) {
	type TestCase struct {
		jsonStr string
		reason  string
	}
	testCases := []TestCase{
		{
			jsonStr: `{"user": "foo"}`,
			reason:  "missing password",
		},
		{
			jsonStr: `{"password": "password"}`,
			reason:  "missing user",
		},
		{
			jsonStr: `{"user": "user", password": "password", "MaxTLSVersion": "1234"}`,
			reason:  "invalid MaxTLSVersion",
		},
		{
			jsonStr: `asdf`,
			reason:  "invalid JSON",
		},
	}
	for _, testCase := range testCases {
		_, err := riakConfigLoad([]byte(testCase.jsonStr))
		if err == nil {
			t.Errorf("unmarshalling bad Riak config - expected error because %s, actual: no error", testCase.reason)
		}
	}
}
