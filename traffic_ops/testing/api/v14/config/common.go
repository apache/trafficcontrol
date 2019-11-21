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

// Package config provides testing helpers for Traffic Ops API tests.
package config

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// militaryTimeFmt defines the 24hr format
// see https://golang.org/pkg/time/#Parse
const militaryTimeFmt = "15:04"

// Validate24HrTimeRange determines whether the provided
// string fits in a format such as "08:00-16:00".
func Validate24HrTimeRange(rng string) error {
	rangeFormat := regexp.MustCompile(`^(\S+)-(\S+)$`)
	match := rangeFormat.FindStringSubmatch(rng)
	if match == nil {
		return fmt.Errorf("string %v is not a range", rng)
	}

	t1, err := time.Parse(militaryTimeFmt, match[1])
	if err != nil {
		return errors.New("time range must be a 24Hr format")
	}

	t2, err := time.Parse(militaryTimeFmt, match[2])
	if err != nil {
		return errors.New("second time range must be a 24Hr format")
	}

	if t1.After(t2) {
		return errors.New("first time should be smaller than the second")
	}

	return nil
}

// ValidateDHMSTimeFormat determines whether the provided
// string fits in a format such as "1d8h", where the valid
// units are days, hours, minutes, and seconds.
func ValidateDHMSTimeFormat(time string) error {

	if time == "" {
		return errors.New("time string cannot be empty")
	}

	dhms := regexp.MustCompile(`^(\d+)([dhms])(\S*)$`)
	match := dhms.FindStringSubmatch(time)

	if match == nil {
		return errors.New("invalid time format")
	}

	var count = map[string]int{
		"d": 0,
		"h": 0,
		"m": 0,
		"s": 0,
	}
	for match != nil {
		if _, err := strconv.Atoi(match[1]); err != nil {
			return err
		}
		if count[match[2]]++; count[match[2]] == 2 {
			return fmt.Errorf("%s unit specified multiple times", match[2])
		}
		match = dhms.FindStringSubmatch(match[3])
	}

	return nil
}

// expandIP matches a general IPv4 pattern and expands the captured octet
// depending on the number of matches. This function does not validate
// that the ip is correct or even that it matches the general pattern, but
// it instead puts the input into a form that go's net package can validate.
func expandIP(ip string) string {

	// d.d.d.d with optional groups of (.d)
	ipRegex := regexp.MustCompile(`^(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:\.(\d+))?$`)
	match := ipRegex.FindStringSubmatch(ip)
	if match == nil {
		return ""
	}

	// ping supports expanding IPv4 addresses
	// PING 1     (0.0.0.1)
	// PING 1.2   (1.0.0.2)
	// PING 1.2.3 (1.2.0.3)
	if match[2] == "" {
		return fmt.Sprintf("0.0.0.%v", match[1])
	}
	if match[3] == "" {
		return fmt.Sprintf("%v.0.0.%v", match[1], match[2])
	}
	if match[4] == "" {
		return fmt.Sprintf("%v.%v.0.%v", match[1], match[2], match[3])
	}
	return ip
}

// parseIP first uses go's net package to test for the common case of a standard
// ipv4 or ipv6 address. If that doesn't pass, it is possible that the address is
// ipv4 written in shorthand notation. The ip is expanded, then we try again.
func parseIP(ip string) net.IP {
	if goip := net.ParseIP(ip); goip != nil {
		return goip
	}
	if goip := net.ParseIP(expandIP(ip)); goip != nil {
		return goip
	}
	return nil
}

// ValidateIPRange validates one of the following forms:
// 1) IP (supports shorthand of IPv4 and IPv6 addresses)
// 2) IP/n (CIDR)
// 3) IP1-IP2 where IP1 < IP2 and type(IP1) == type(IP2)
func ValidateIPRange(ip string) error {

	var err error
	var splt []string

	splt = strings.Split(ip, "-")
	if len(splt) == 2 {
		ip1 := parseIP(splt[0])
		ip2 := parseIP(splt[1])

		// both must be valid
		if ip1 == nil || ip2 == nil {
			return fmt.Errorf("invalid IP range: %v", ip)
		}

		// must be of the same type
		if (ip1.To4() == nil) != (ip2.To4() == nil) {
			return fmt.Errorf("invalid IP range: %v", ip)
		}

		// ip2 must be less than ip1
		if bytes.Compare(ip2, ip1) < 0 {
			return fmt.Errorf("invalid IP range: %v", ip)
		}
		return nil
	}

	if goip := parseIP(ip); goip != nil {
		return nil
	}

	// first try for CIDR
	if _, _, err = net.ParseCIDR(ip); err == nil {
		return nil
	}

	// if it looks like we have a CIDR pattern we try to expand the ip and try again
	splt = strings.Split(ip, "/")
	if len(splt) == 2 {
		ip = fmt.Sprintf("%v/%v", expandIP(splt[0]), splt[1])
		if _, _, err = net.ParseCIDR(ip); err == nil {
			return nil
		}
		return fmt.Errorf("invalid CIDR address: %v", ip)
	}

	return fmt.Errorf("invalid IP range: %v", ip)
}
