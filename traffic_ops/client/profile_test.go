/*
     Copyright 2015 Comcast Cable Communications Management, LLC

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

package client

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestProfile(t *testing.T) {
	fmt.Println("Running Profile Tests")
	text, err := ioutil.ReadFile("testdata/profiles.json")
	if err != nil {
		t.Skip("Skipping parameters test, no profiles.json found.")
	}

	profileList, err := profileUnmarshall(text)
	if err != nil {
		t.Fatal(err)
	}
	for _, profile := range profileList.Response {
		name := profile.Name
		if len(name) == 0 {
			t.Fatal("Profile result does not contain 'Name'")
		}
		if len(profile.Description) == 0 {
			t.Error("Description is null for profile: " + name)
		}
		if len(profile.LastUpdated) == 0 {
			t.Error("LastUpdated is null for profile: " + name)
		}
	}
}
