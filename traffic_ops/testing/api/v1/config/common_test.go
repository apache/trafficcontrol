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
		{"", false},
		{" ", false},
		{"-", false},
		{"--", false},
		{" - ", false},
		{"asdf", false},
		{"08:00-asdf", false},
		{"asdf-08:00", false},
		{"-08:00", false},
		{"08:00-", false},
		{"08-09:00", false},
		{"09:00-10", false},
		{"09:00-10:0", false},
		{"9:00-10:0", false},
		{"08:00-32:00", false},
		{"32:00-33:00", false},
		{"08:00--16:00", false},
		{"08:00-16:00-", false},
		{"08:00-16:00-17:00", false},
		{"08:00-09:00 16:00-17:00", false},
		{"foo 16:00-17:00", false},
		{"16:00-17:00 foo", false},
	}

	for _, test := range tests {
		if err := Validate24HrTimeRange(test.time); (err == nil) != test.ok {
			if test.ok {
				t.Errorf(`
  test should have passed
  time: %v
  err: %v`, test.time, err)
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
  time: %v
  err: %v`, test.time, err)
			} else {
				t.Errorf(`
  test should not have passed
  time: %v`, test.time)
			}
		}
	}

}

// TestIPRange
// \__> parseIpP (only two simple branches)
// \__> expandIP (simple enough to test all 5 cases)
func TestIPRange(t *testing.T) {

	var tests = []struct {
		ip string
		ok bool
	}{
		{"", false},
		{"x", false},
		{"0", true},
		{"0.0", true},
		{"0.0.0", true},
		{"0.0.0.0", true},
		{"256.0.0.0", false},
		{"0.0.0.0/0", true},
		{"0.0.0.0/x", false},
		{"0.0.0/0", true},
		{"0.0.0.0-255.255.255.255", true},
		{"255.255.255.255-0.0.0.0", false},
		{"0.0.0-255.255.255.255", true},
		{"0.0.0.0-255.255.255", true},
		{"0.0.0.x-255.255.255.255", false},
		{"::/0", true},
		{"::-ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", true},
		{"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff-::", false},
		{"::-0", false},
		{"0.0.0.0-ffff::ffff", false},
		{"::-ffff::ffff", true},
	}

	for _, test := range tests {
		if err := ValidateIPRange(test.ip); (err == nil) != test.ok {
			if test.ok {
				t.Errorf(`
  test should have passed
  ip: %v
  error: %v`, test.ip, err)
			} else {
				t.Errorf(`
  test should not have passed
  ip: %v`, test.ip)
			}
		}
	}
}
