package v3

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
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-rfc"
	"github.com/apache/trafficcontrol/v8/lib/go-tc"
)

// These capabilities are defined during the setup process in todb.go.
// ANY TIME THOSE ARE CHANGED THIS MUST BE UPDATED.
var staticCapabilities = []tc.Capability{
	tc.Capability{
		Name:        "all-read",
		Description: "Full read access",
	},
	tc.Capability{
		Name:        "all-write",
		Description: "Full write access",
	},
	tc.Capability{
		Name:        "cdn-read",
		Description: "View CDN configuration",
	},
	tc.Capability{
		Name:        "asns-read",
		Description: "Read ASNs",
	},
	tc.Capability{
		Name:        "asns-write",
		Description: "Write ASNs",
	},
	tc.Capability{
		Name:        "cache-groups-read",
		Description: "Read CGs",
	},
}

func TestCapabilities(t *testing.T) {
	CreateTestCapabilities(t)
	GetTestCapabilitiesIMS(t)
	GetTestCapabilities(t)
}

func CreateTestCapabilities(t *testing.T) {
	db, err := OpenConnection()
	if err != nil {
		t.Fatal("cannot open db")
	}
	defer db.Close()
	dbInsertTemplate := `INSERT INTO capability (name, description) VALUES ('%v', '%v');`

	for _, c := range testData.Capabilities {
		err = execSQL(db, fmt.Sprintf(dbInsertTemplate, c.Name, c.Description))
		if err != nil {
			t.Errorf("could not create capability: %v", err)
		}
	}
}

func GetTestCapabilitiesIMS(t *testing.T) {
	var header http.Header
	header = make(map[string][]string)
	futureTime := time.Now().AddDate(0, 0, 1)
	time := futureTime.Format(time.RFC1123)
	header.Set(rfc.IfModifiedSince, time)
	testDataLen := len(testData.Capabilities) + len(staticCapabilities)
	capMap := make(map[string]string, testDataLen)

	for _, c := range testData.Capabilities {
		capMap[c.Name] = c.Description
		_, reqInf, err := TOSession.GetCapabilityWithHdr(c.Name, header)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err.Error())
		}
		if reqInf.StatusCode != http.StatusNotModified {
			t.Fatalf("Expected 304 status code, got %v", reqInf.StatusCode)
		}
	}
}

func GetTestCapabilities(t *testing.T) {
	testDataLen := len(testData.Capabilities) + len(staticCapabilities)
	capMap := make(map[string]string, testDataLen)

	for _, c := range testData.Capabilities {
		capMap[c.Name] = c.Description
		cap, _, err := TOSession.GetCapability(c.Name)
		if err != nil {
			t.Errorf("could not get capability '%s': %v", c.Name, err)
			continue
		}

		if cap.Name != c.Name {
			t.Errorf("requested capacity '%s' but got a capacity with the name '%s'", c.Name, cap.Name)
		}
		if cap.Description != c.Description {
			t.Errorf("capacity '%s' has the wrong description, want '%s' but got '%s'", c.Name, c.Description, cap.Description)
		}
	}

	// Hopefully this won't need to be done for much longer
	for _, c := range staticCapabilities {
		capMap[c.Name] = c.Description
	}

	caps, _, err := TOSession.GetCapabilities()
	if err != nil {
		t.Fatalf("could not get all capabilities: %v", err)
	}
	if len(caps) != testDataLen {
		t.Fatalf("response returned different number of capabilities than those that exist; got %d, want %d", len(caps), testDataLen)
	}

	for _, c := range caps {
		if desc, ok := capMap[c.Name]; !ok {
			t.Errorf("capability '%s' found in response, but not in test data!", c.Name)
		} else {
			if desc != c.Description {
				t.Errorf("capability '%s' has description '%s' in response, but had '%s' in the test data", c.Name, c.Description, desc)
			}
			delete(capMap, c.Name)
		}
	}

	for c, _ := range capMap {
		t.Errorf("Capability '%s' existed in the test data but didn't appear in the response!", c)
	}
}
