package v3

/*

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

import (
	"testing"
)

func TestAPICapabilities(t *testing.T) {
	testCases := []struct {
		description string
		capability  string
		order       string
		first       string
		hasRecords  bool
	}{
		{
			description: "Successfully get all asns-write API Capabilities",
			capability:  "asns-write",
			hasRecords:  true,
		},
		{
			description: "Successfully get all asns-read API Capabilities",
			capability:  "asns-read",
			hasRecords:  true,
		},
		{
			description: "Successfully get all cache-groups-read API Capabilities",
			capability:  "cache-groups-read",
			hasRecords:  true,
		},
		{
			description: "Fail to get any API Capabilities with a bogus capability",
			capability:  "foo",
			hasRecords:  false,
		},
		{
			description: "Successfully get all API Capabilities in order of HTTP Method",
			order:       "httpMethod",
			first:       "GET",
		},
	}

	for _, c := range testCases {
		t.Run(c.description, func(t *testing.T) {
			caps, _, err := TOSession.GetAPICapabilities(c.capability, c.order)

			if err != nil {
				t.Fatalf("error retrieving API capabilities: %s", err.Error())
			}

			if len(caps.Response) == 0 && c.hasRecords {
				t.Fatalf("error: expected capability %s to have records, but found 0", c.capability)
			}

			if c.order != "" {
				if c.first != caps.Response[0].HTTPMethod {
					t.Fatalf("error: expected first element to be %s, got %s", c.first, caps.Response[0].HTTPMethod)
				}
			}
		})
	}

}
