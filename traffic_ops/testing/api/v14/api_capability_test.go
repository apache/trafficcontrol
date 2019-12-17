package v14

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
	"fmt"
	"strings"
	"testing"
)

func TestAPICapabilities(t *testing.T) {
	WithObjs(t, []TCObj{}, func() {
		GetTestAPICapabilities(t)
	})
}

func GetTestAPICapabilities(t *testing.T) {
	testCases := []struct {
		description string
		capability  string
		order       string
		count       int
		err         string
	}{
		{
			description: "Successfully get all API Capabilities",
			capability:  "",
			order:       "",
			count:       4,
			err:         "",
		},
	}

	for _, c := range testCases {
		caps, _, err := TOSession.GetAPICapabilities(c.capability, c.order)

		fmt.Printf("\nerror is: %v\n", err)
		if err != nil && !strings.Contains(err.Error(), c.err) {
			t.Fatalf("error: expected error result %v, want: %v", err, c.err)
		}

		if len(caps.Response) != c.count {
			t.Fatalf("error: expected len(caps) to be %d, got %d", c.count, len(caps.Response))
		}
	}

}
