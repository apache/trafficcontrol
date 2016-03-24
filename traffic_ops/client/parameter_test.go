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

func TestParameter(t *testing.T) {
	fmt.Println("Running Parameter Tests")
	text, err := ioutil.ReadFile("testdata/parameters.json")
	if err != nil {
		t.Skip("Skipping parameters test, no parameters.json found.")
	}

	parameterList, err := paramUnmarshall(text)
	if err != nil {
		t.Fatal(err)
	}
	for _, parameter := range parameterList.Response {
		name := parameter.Name
		if len(name) == 0 {
			t.Fatal("param result does not contain 'Name'")
		}
		if len(parameter.ConfigFile) == 0 {
			t.Errorf("ConfigFile parameter is null for parameter: %s", name)
		}
		if len(parameter.LastUpdated) == 0 {
			t.Errorf("LastUpdate parameter is null for parameter: %s", name)
		}
		if len(parameter.Value) == 0 {
			t.Errorf("Value parameter is null for parameter: %s", name)
		}
	}
}
