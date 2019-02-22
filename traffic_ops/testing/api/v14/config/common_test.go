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
