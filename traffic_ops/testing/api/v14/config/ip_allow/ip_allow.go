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

	. "github.com/apache/trafficcontrol/traffic_ops/testing/api/v14/config"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/test"
)

func parseConfig(config string) test.Error {
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

		// The differences between the cases:
		// IP
		// IP/0
		// IP-IP
		// If I control how the base IP is understood then
		// it is very easy to split up based on finding "/" or "-"
		// then validate with my own method.
		// That being said, I don't know enough about CIDR to know
		// if the pure textual representation is semantically sufficient.
		// The only new thing about the CIDR notation I learned is that
		// the subnet prefix has a clear maximum value that is the maximum
		// number of digits used in an ip.
		//
		// Note: The CIDR formatted like IP/32 equals IP

		if err := ValidateIPRange(rhs); err != nil {
			return ErrorContext.AddErrorCode(InvalidIPRange, err)
		}

		/*
			// CASE 3: 0.0.0.0
			// A fully-specified IP range
			if ip := net.ParseIP(rhs); ip != nil {
				return nil
			}

			// CASE 2: 0.0.0.0/0
			// A fully-specified CIDR range
			if _, _, err := net.ParseCIDR(rhs); err == nil {
				return nil
			}

			// CASE 1: 0.0.0.0-255.255.255.255
			// A fully-specified range to another fully-specified range
			ips := strings.Split(rhs, "-")
			if len(ips) != 2 {
				// Note: Need a new error constant
				return ErrorContext.NewError(InvalidIP, "invalid range format")
			}
			if ip := net.ParseIP(ips[0]); ip == nil {
				return ErrorContext.NewError(InvalidIP, `"%s"`, rhs)
			}
			if ip := net.ParseIP(ips[1]); ip == nil {
				return ErrorContext.NewError(InvalidIP, `"%s"`, rhs)
			} */

	case "action":
		switch rhs {
		case "ip_allow":
		case "ip_deny":
		default:
			return ErrorContext.NewError(InvalidAction)
		}
	case "method":
		methods := strings.Split(rhs, "|")
		for _, method := range methods {
			switch method {
			// assuming all methods are valid
			// RFC 2616-9 for list of all methods
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
				return ErrorContext.NewError(InvalidMethod, `invalid method "%v"`, method)
			}
		}

	default:
		return ErrorContext.NewError(InvalidLabel)
	}

	return nil
}

// count:
// 1 src_ip OR 1 dest_ip (total=1)
// 1(exact) action
// 1(max) or 0 method
//
// Note:
// It might make sense to have a simple parseConfigRule, and give a list of functions to this
// method. I might want to consider this to avoid copying and pasting.. BUT I should be careful
// as well.
//
// I think this idea of keeping count is cool, but golang makes it difficult to implement
// 1) without interfaces or 2) without verbosity. I could treat strings as ints
// label -> label_grouping
// label_grouping -> count
// label -> count

func parseConfigRule(rule string) test.Error {

	//var destination bool
	//var action bool
	var match []string
	var err test.Error

	rule = strings.Trim(rule, "\t ")
	if rule == "" || strings.HasPrefix(rule, "#") {
		return nil
	}

	assignments := strings.Fields(rule)
	last := len(assignments) - 1
	if last < 1 {
		return ErrorContext.NewError(NotEnoughAssignments)
	}

	// neither the rhs or lhs can contain any whitespace
	assignment := regexp.MustCompile(`([a-z_\-\d]+)=(\S+)`)
	for _, elem := range assignments {
		match = assignment.FindStringSubmatch(strings.ToLower(elem))
		if match == nil {
			return ErrorContext.NewError(BadAssignmentMatch, `could not match assignment: "%v"`, elem)
		}

		err = parseLabelValue(match[1], match[2])
		if err == nil {
			//if count[match[1]]++; count[match[1]] == 2 {
			//return ErrorContext.NewError(ExcessLabel, `the label "%s" can only be used once per rule`, match[1])
			continue
		}

		if err.Code() == InvalidLabel {
			return err
		} else {
			return err.Prepend(`could not parse %s:`, match[1])
		}

	}

	//if !destination {
	//	return ErrorContext.NewError(MissingLabel, "missing primary destination label")
	//}

	return nil
}
