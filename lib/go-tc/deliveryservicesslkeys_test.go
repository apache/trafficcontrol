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
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-util"
)

// - Perl: 2018-08-21+14:26:06
// - Portal: 2018-08-21 14:14:42
// - RFC3339

func TestCDNDNSSECGenerateReqDateUnmarshalJSON(t *testing.T) {
	type T struct {
		D CDNDNSSECGenerateReqDate `json:"d"`
	}
	obj := T{}

	// epoch
	s := `{"d": 1534869508}`
	if err := json.Unmarshal([]byte(s), &obj); err != nil {
		t.Fatalf("unmarshalling epoch CDNDNSSECGenerateReqDate: error expected nil, actual:  %+v", err)
	}
	if st := int64(1534869508); int64(obj.D) != st {
		t.Fatalf("unmarshalling epoch CDNDNSSECGenerateReqDate: error expected %+v, actual:  %+v", st, int64(obj.D))
	}

	// RFC3339
	s = `{"d": "2018-08-21T10:58:56-06:00"}`
	if err := json.Unmarshal([]byte(s), &obj); err != nil {
		t.Fatalf("unmarshalling RFC3339 CDNDNSSECGenerateReqDate: error expected nil, actual:  %+v", err)
	}
	if st := int64(1534870736); int64(obj.D) != st {
		t.Fatalf("unmarshalling RFC3339 CDNDNSSECGenerateReqDate: error expected %+v, actual:  %+v", st, int64(obj.D))
	}

	// old Portal date format
	s = `{"d": "2018-10-31 1:12:06"}`
	if err := json.Unmarshal([]byte(s), &obj); err != nil {
		t.Fatalf("unmarshalling old Portal CDNDNSSECGenerateReqDate: error expected nil, actual:  %+v", err)
	}
	if st := int64(1540948326); int64(obj.D) != st {
		t.Fatalf("unmarshalling old Portal CDNDNSSECGenerateReqDate: error expected %+v, actual:  %+v", st, int64(obj.D))
	}

	// old Perl date format
	s = `{"d": "2018-04-19+02:04:05"}`
	if err := json.Unmarshal([]byte(s), &obj); err != nil {
		t.Fatalf("unmarshalling old Perl CDNDNSSECGenerateReqDate: error expected nil, actual:  %+v", err)
	}
	if st := int64(1524103445); int64(obj.D) != st {
		t.Fatalf("unmarshalling old Perl CDNDNSSECGenerateReqDate: error expected %+v, actual:  %+v", st, int64(obj.D))
	}

	// invalid date format
	s = `{"d": "Mon Jan 12 15:04:05 -0700 MST 2018"}`
	if err := json.Unmarshal([]byte(s), &obj); err == nil {
		t.Fatalf("unmarshalling invalid CDNDNSSECGenerateReqDate: error expected, actual nil")
	}
	s = `{"d": "01-25-2018"}`
	if err := json.Unmarshal([]byte(s), &obj); err == nil {
		t.Fatalf("unmarshalling invalid format CDNDNSSECGenerateReqDate: error expected, actual nil")
	}
	s = `{"d": "foo"}`
	if err := json.Unmarshal([]byte(s), &obj); err == nil {
		t.Fatalf("unmarshalling invalid string CDNDNSSECGenerateReqDate: error expected, actual nil")
	}
	s = `{"d": "42.0"}`
	if err := json.Unmarshal([]byte(s), &obj); err == nil {
		t.Fatalf("unmarshalling invalid float string CDNDNSSECGenerateReqDate: error expected, actual nil")
	}
	s = `{"d": 42.0}`
	if err := json.Unmarshal([]byte(s), &obj); err == nil {
		t.Fatalf("unmarshalling invalid float CDNDNSSECGenerateReqDate: error expected, actual nil")
	}
}

func TestSSLKeysReqValidate(t *testing.T) {
	req := DeliveryServiceAddSSLKeysReq{}
	req.CDN = util.StrPtr("foo")
	req.DeliveryService = util.StrPtr("bar")
	req.Key = util.StrPtr("bar")
	ver := util.JSONIntStr(1)
	req.Version = &ver
	cert := DeliveryServiceSSLKeysCertificate{}
	cert.Crt = ""
	cert.CSR = ""
	cert.Key = ""
	req.Certificate = &cert
	if err := req.Validate(nil); err == nil {
		t.Error("expected validation to return an error")
	}
}
