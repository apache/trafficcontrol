package deliveryservice

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
	"encoding/base64"
	"fmt"
	"math"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

func TestMakeDSRecordText(t *testing.T) {
	ksk := tc.DNSSECKeyV11{
		DSRecord: &tc.DNSSECKeyDSRecordV11{
			Algorithm:  7,
			DigestType: 1,
		},
		Public: "",
	}
	_, err := MakeDSRecordText(ksk, 0)
	if err == nil {
		t.Error("Expected a blank 'Public' field to yield an error")
	} else {
		t.Logf("Got expected error from blank 'Public': %v", err)
	}

	ksk.Public = "not a base64 string"
	_, err = MakeDSRecordText(ksk, 0)
	if err == nil {
		t.Error("Expected an invalid (non-base64-encoded-string) 'Public' field to yield an error")
	} else {
		t.Logf("Got expected error from non-base64 'Public': %v", err)
	}

	ksk.Public = base64.RawStdEncoding.EncodeToString([]byte("x x x x 4 5 x invalid! x"))
	_, err = MakeDSRecordText(ksk, 0)
	if err == nil {
		t.Error("Expected an invalid (bad public key) 'Public' field to yield an error")
	} else {
		t.Logf("Got expected error from 'Public' with invalid public key: %v", err)
	}

	key := base64.StdEncoding.EncodeToString([]byte("This is a public key, I swear"))

	ksk.Public = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("x x x x %d 5 x %s x", uint64(math.Pow(2, 16)+1), key)))
	_, err = MakeDSRecordText(ksk, 0)
	if err == nil {
		t.Error("Expected a 'Public' field with a 'flags' too wide to fit in a uint16 to yield an error")
	} else {
		t.Logf("Got expected error from 'flags' too wide: %v", err)
	}

	ksk.Public = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("x x x x 4 %d x %s x", uint64(math.Pow(2, 8)+1), key)))
	_, err = MakeDSRecordText(ksk, 0)
	if err == nil {
		t.Error("Expected a 'Public' field with a 'protocol' too wide to fit in a uint8 to yield an error")
	} else {
		t.Logf("Got expected error from 'protocol' too wide: %v", err)
	}

	ksk.Public = base64.StdEncoding.EncodeToString([]byte("x x x x 4 5 x " + key + " x"))
	_, err = MakeDSRecordText(ksk, 0)
	if err != nil {
		t.Errorf("Unexpected error for a valid 'Public' field: %v", err)
	}
}
