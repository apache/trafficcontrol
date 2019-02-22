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

package config

import "testing"

func Test24HrTimeRange(t *testing.T) {

	var tests = []struct {
		time string
		ok   bool
	}{
		{"15:04", false},
		{"16:00-08:00", false},
		{"16:00-8:00", false},
		{"xxx-8:00", false},
		{"8:00-xxx", false},
		{"8:00-16:00", true},
		{"08:00-16:00", true},
	}

	for _, test := range tests {
		if err := Validate24HrTimeRange(test.time); (err == nil) != test.ok {
			if test.ok {
				t.Errorf(`
  test should have passed
  time: %v`, test.time)
			} else {
				t.Errorf(`
  test should not have passed
  time: %v`, test.time)
			}
		}
	}

}

func TestDHMSTimeFormat(t *testing.T) {

	var tests = []struct {
		time string
		ok   bool
	}{
		{"1s", true},
		{"1m", true},
		{"1h", true},
		{"1d", true},
		{"1d2h", true},
		{"1m2s", true},
		{"1d2h3m4s", true},
		{"1s2h3m4d", true},
		{"10000000000000000000000s", false},
		{"1s2s", false},
		{"1x", false},
		{"1", false},
		{"x", false},
		{"", false},
	}

	for _, test := range tests {
		if err := ValidateDHMSTimeFormat(test.time); (err == nil) != test.ok {
			if test.ok {
				t.Errorf(`
  test should have passed
  time: %v`, test.time)
			} else {
				t.Errorf(`
  test should not have passed
  time: %v`, test.time)
			}
		}
	}

}
