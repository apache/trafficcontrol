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

package ip_allow

import (
	"regexp"
	"strings"

	"github.com/apache/trafficcontrol/traffic_ops/testing/api/v14/config"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/test"
)

// TODO Godoc
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

func parseLabelValue(lhs string, rhs string) test.Error {

	switch lhs {
	case "src_ip":
		fallthrough
	case "dest_ip":
		if err := config.ValidateIPRange(rhs); err != nil {
			return config.ErrorContext.AddErrorCode(config.InvalidIPRange, err)
		}
	case "action":
		switch rhs {
		case "ip_allow":
		case "ip_deny":
		default:
			return config.ErrorContext.NewError(config.InvalidAction)
		}
	case "method":
		methods := strings.Split(rhs, "|")
		for _, method := range methods {
			switch method {
			// assuming all methods are valid
			// see RFC 2616-9 for list of all methods
			// PUSH and PURGE are specific to ATS
			case "all":
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
				return config.ErrorContext.NewError(config.UnknownMethod, `unknown method "%v"`, method)
			}
		}

	default:
		return config.ErrorContext.NewError(config.InvalidLabel)
	}

	return nil
}

func parseConfigRule(rule string) test.Error {

	var ip bool
	var action bool

	var match []string
	var err test.Error

	rule = strings.Trim(rule, "\t ")
	if rule == "" || strings.HasPrefix(rule, "#") {
		return nil
	}

	assignments := strings.Fields(rule)
	last := len(assignments) - 1
	if last < 1 {
		return config.ErrorContext.NewError(config.NotEnoughAssignments)
	}

	// neither the rhs or lhs can contain any whitespace
	assignment := regexp.MustCompile(`([a-z_\-\d]+)=(\S+)`)
	for _, elem := range assignments {
		match = assignment.FindStringSubmatch(strings.ToLower(elem))
		if match == nil {
			return config.ErrorContext.NewError(config.BadAssignmentMatch, `could not match assignment: "%v"`, elem)
		}

		err = parseLabelValue(match[1], match[2])
		if err == nil {
			switch match[1] {
			case "action":
				if action {
					return config.ErrorContext.NewError(config.ExcessLabel, "only one action is allowed per rule")
				}
				action = true
			case "src_ip":
				fallthrough
			case "dest_ip":
				if ip {
					return config.ErrorContext.NewError(config.ExcessLabel, "only one of src_ip or dest_ip is allowed per rule")
				}
				ip = true
			}
			continue
		}

		if err.Code() == config.InvalidLabel {
			return err
		} else {
			return err.Prepend(`could not parse %s:`, match[1])
		}
	}

	if !ip {
		return config.ErrorContext.NewError(config.MissingLabel, "missing either src_ip or dest_ip label")
	}
	if !action {
		return config.ErrorContext.NewError(config.MissingLabel, "missing action label")
	}

	return nil
}
