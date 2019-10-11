package v14

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

import "testing"

import "github.com/apache/trafficcontrol/lib/go-log"
import "github.com/apache/trafficcontrol/lib/go-tc"

// These capabilities are defined during the setup process in todb.go.
// ANY TIME THOSE ARE CHANGED THIS MUST BE UPDATED.
var staticCapabilities = []tc.Capability {
	tc.Capability{
		Name: "all-read",
		Description: "Full read access",
	},
	tc.Capability {
		Name: "all-write",
		Description: "Full write access",
	},
	tc.Capability {
		Name: "cdn-read",
		Description: "View CDN configuration",
	},
}

func TestCapabilities(t *testing.T) {
	WithObjs(t, []TCObj{Capabilities}, func() {
		ReplaceTestCapability(t)
		GetTestCapabilities(t)
	})
}

func CreateTestCapabilities(t *testing.T) {
	for _,c := range testData.Capabilities {
		resp, _, err := TOSession.CreateCapability(c)
		log.Debugln("Response: ", c.Name, " ", resp)
		if err != nil {
			t.Errorf("could not create capability: %v", err)
		}
	}
}

func GetTestCapabilities(t *testing.T) {
	testDataLen := len(testData.Capabilities) + len(staticCapabilities)
	capMap := make(map[string]string, testDataLen)

	for _,c := range testData.Capabilities {
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
	for _,c := range staticCapabilities {
		capMap[c.Name] = c.Description
	}


	caps, _, err := TOSession.GetCapabilities()
	if err != nil {
		t.Fatalf("could not get all capabilities: %v", err)
	}
	if len(caps) != testDataLen {
		t.Fatalf("response returned different number of capabilities than those that exist; got %d, want %d", len(caps), testDataLen)
	}

	for _,c := range caps {
		if desc, ok := capMap[c.Name]; !ok {
			t.Errorf("capability '%s' found in response, but not in test data!", c.Name)
		} else {
			if desc != c.Description {
				t.Errorf("capability '%s' has description '%s' in response, but had '%s' in the test data", c.Name, c.Description, desc)
			}
			delete(capMap, c.Name)
		}
	}

	for c,_ := range capMap {
		t.Errorf("Capability '%s' existed in the test data but didn't appear in the response!", c)
	}
}

func ReplaceTestCapability(t *testing.T) {
	if len(testData.Capabilities) < 1 {
		t.Fatalf("No test capabilities!")
	}

	c := testData.Capabilities[0]
	c.Name += "TEST"
	c.Description += "REPLACE TEST"

	alerts, _, err := TOSession.ReplaceCapabilityByName(testData.Capabilities[0].Name, c)
	log.Debugln("alerts:", alerts)
	if err != nil {
		t.Fatalf("Failed to replace capability %s: %v", testData.Capabilities[0].Name, err)
	}

	// The old one shouldn't exist anymore
	resp, _, err := TOSession.GetCapability(testData.Capabilities[0].Name)
	if err == nil {
		t.Errorf("Expected an error, but got response: %v", resp)
	}

	// ... but the new one should
	resp, _, err = TOSession.GetCapability(c.Name)
	if err != nil {
		t.Errorf("Expected to get new capability, error: %v", err)
	}

	// Now we need to replace it again, so the DELETE can catch it
	alerts, _, err = TOSession.ReplaceCapabilityByName(c.Name, testData.Capabilities[0])
	log.Debugln("alerts:", alerts)
	if err != nil {
		t.Errorf("Failed to re-replace capability: %s", err)
	}

	log.Debugln("response:", resp)
}

func DeleteTestCapabilities(t *testing.T) {
	for _,c := range testData.Capabilities {
		resp, _, err := TOSession.DeleteCapability(c.Name)
		log.Debugln("Response: ", c.Name, " ", resp)
		if err != nil {
			t.Errorf("could not delete capability: %v", err)
		}
	}
}
