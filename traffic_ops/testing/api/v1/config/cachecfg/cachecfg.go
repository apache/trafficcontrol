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

package cachecfg

import (
	"regexp"
	"strings"

	"github.com/apache/trafficcontrol/traffic_ops/testing/api/v1/config"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/test"
	"github.com/go-ozzo/ozzo-validation/is"
)

// Parse takes a string presumed to be an ATS cache.config and validates that it is
// syntatically correct.
//
// The general format of a cache config is three types of labels separated by spaces:
//
//  primary_destination=value secondary_specifier=value action=value
//
// For a full description of how to format a cache config, refer to the ATS documentation
// for the cache config:
// https://docs.trafficserver.apache.org/en/latest/admin-guide/files/cache.config.en.html
//
func Parse(config string) test.Error {
	lines := strings.Split(config, "\n")

	if len(lines) == 1 {
		return parseConfigRule(lines[0])
	}

	for i, ln := range lines {
		err := parseConfigRule(ln)
		if err != nil {
			return err.Prepend("error on line %d: ", i+1)
		}
	}

	return nil
}

func parsePrimaryDestinations(lhs string, rhs string) test.Error {

	switch lhs {
	case "dest_domain":
		// dest_host is an alias for dest_domain
		fallthrough
	case "dest_host":
		if err := is.Host.Validate(rhs); err != nil {
			return config.ErrorContext.NewError(config.InvalidHost, `"%s" %v`, rhs, err)
		}
	case "dest_ip":
		if err := is.IP.Validate(rhs); err != nil {
			return config.ErrorContext.NewError(config.InvalidIP, `"%s" %v`, rhs, err)
		}
	case "host_regex":
		fallthrough
	case "url_regex":
		// only makes sure the regex compiles, not that the regex generates anything valid
		if _, err := regexp.Compile(rhs); err != nil {
			return config.ErrorContext.NewError(config.InvalidRegex, "%v", err)
		}
	default:
		return config.ErrorContext.NewError(config.InvalidLabel)
	}

	return nil
}

func parseSecondarySpecifiers(lhs string, rhs string) test.Error {

	switch lhs {
	case "port":
		if err := is.Port.Validate(rhs); err != nil {
			return config.ErrorContext.AddErrorCode(config.InvalidPort, err)
		}
	case "scheme":
		if rhs != "http" && rhs != "https" {
			return config.ErrorContext.NewError(config.InvalidHTTPScheme)
		}
	case "prefix":
		// no clear validation to perform
	case "suffix":
		// examples: gif jpeg
		// no clear validation to perform
	case "method":
		// assuming all methods are valid
		// see RFC 2616-9 for list of all methods
		// PURGE and PUSH are specific to ATS
		switch rhs {
		case "get":
		case "put":
		case "post":
		case "delete":
		case "trace":
		case "options":
		case "head":
		case "connect":
		case "patch":
		case "purge":
		case "push":
		default:
			return config.ErrorContext.NewError(config.UnknownMethod, `unknown method "%v"`, rhs)
		}

	case "time":
		if err := config.Validate24HrTimeRange(rhs); err != nil {
			return config.ErrorContext.AddErrorCode(config.InvalidTimeRange24Hr, err)
		}
	case "src_ip":
		if err := is.IP.Validate(rhs); err != nil {
			return config.ErrorContext.AddErrorCode(config.InvalidIP, err)
		}
	case "internal":
		if rhs != "true" && rhs != "false" {
			return config.ErrorContext.NewError(config.InvalidBool)
		}

	default:
		return config.ErrorContext.NewError(config.InvalidLabel)
	}

	return nil
}

func parseActions(lhs string, rhs string) test.Error {

	switch lhs {
	case "action":
		switch rhs {
		case "never-cache":
		case "ignore-no-cache":
		case "ignore-client-no-cache":
		case "ignore-server-no-cache":
		default:
			return config.ErrorContext.NewError(config.InvalidAction)
		}

	case "cache-responses-to-cookies":
		digit := rhs[0]
		if digit < '0' || '4' > digit || len(rhs) > 1 {
			return config.ErrorContext.NewError(config.InvalidCacheCookieResponse)
		}

	// All of these are time formats
	case "pin-in-cache":
		fallthrough
	case "revalidate":
		fallthrough
	case "ttl-in-cache":
		err := config.ValidateDHMSTimeFormat(rhs)
		if err != nil {
			return config.ErrorContext.AddErrorCode(config.InvalidTimeFormatDHMS, err)
		}
	default:
		return config.ErrorContext.NewError(config.InvalidLabel)
	}

	return nil
}

func parseConfigRule(rule string) test.Error {

	var err test.Error

	rule = strings.TrimSpace(rule)
	if rule == "" || strings.HasPrefix(rule, "#") {
		return nil
	}

	assignments := strings.Fields(rule)
	last := len(assignments) - 1
	if last < 1 {
		return config.ErrorContext.NewError(config.NotEnoughAssignments)
	}

	// no individual secondary specifier label can be used twice
	count := map[string]int{
		"port":     0,
		"scheme":   0,
		"prefix":   0,
		"suffix":   0,
		"method":   0,
		"time":     0,
		"src_ip":   0,
		"internal": 0,
	}

	// neither the rhs or lhs can contain any whitespace
	assignment := regexp.MustCompile(`([a-z_\-\d]+)=(\S+)`)

	destination := false
	action := false

	for _, elem := range assignments {
		match := assignment.FindStringSubmatch(strings.ToLower(elem))
		if match == nil {
			return config.ErrorContext.NewError(config.BadAssignmentMatch, `could not match assignment: "%v"`, elem)
		}

		err = parsePrimaryDestinations(match[1], match[2])
		if err == nil {
			if destination {
				return config.ErrorContext.NewError(config.ExcessLabel, "too many primary destination labels")
			} else {
				destination = true
				continue
			}
		}
		if err.Code() != config.InvalidLabel {
			return err.Prepend(`coult not parse primary destination from "%s": `, match[0])
		}

		err = parseSecondarySpecifiers(match[1], match[2])
		if err == nil {
			if count[match[1]]++; count[match[1]] == 2 {
				return config.ErrorContext.NewError(config.ExcessLabel, `the label "%s" can only be used once per rule`, match[1])
			}
			continue
		}
		if err.Code() != config.InvalidLabel {
			return err.Prepend(`could not parse secondary specifier from "%s": `, match[0])
		}

		err = parseActions(match[1], match[2])
		if err == nil {
			action = true
			continue
		}

		if err.Code() == config.InvalidLabel {
			return err
		} else {
			return err.Prepend(`could not parse action from "%s": `, match[0])
		}

	}

	if !destination {
		return config.ErrorContext.NewError(config.MissingLabel, "missing primary destination label")
	}

	if !action {
		return config.ErrorContext.NewError(config.MissingLabel, "missing action lablel")
	}

	return nil
}
