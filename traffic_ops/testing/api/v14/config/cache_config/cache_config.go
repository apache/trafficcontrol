package cache_config

import (
	"regexp"
	"strings"

	. "github.com/apache/trafficcontrol/traffic_ops/testing/api/v14/config"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/test"
	"github.com/go-ozzo/ozzo-validation/is"
)

// Kindof an alias. I should change it back later
// maybe this would be fine if I wasn't already using the dot import
var cxt = ErrorContext

func parseCacheConfig(config string) test.Error {
	lines := strings.Split(config, "\n")

	if len(lines) == 1 {
		return parseCacheConfigRule(lines[0])
	}

	for i, ln := range lines {
		err := parseCacheConfigRule(ln)
		if err != nil {
			return err.Prepend("error on line %d: %v", i+1)
		}
	}
	return nil
}

type parser func(string, string) test.Error

func parsePrimaryDestinations(lhs string, rhs string) test.Error {

	switch lhs {
	case "dest_domain":
		// dest_host is an alias for dest_domain
		fallthrough
	case "dest_host":
		if err := is.Host.Validate(rhs); err != nil {
			return cxt.NewError(InvalidHost, "\"%s\" %v", rhs, err)
		}
	case "dest_ip":
		if err := is.IP.Validate(rhs); err != nil {
			return cxt.NewError(InvalidIP, "\"%s\" %v", rhs, err)
		}
	case "host_regex":
		// only makes sure the regex compiles, not that the regex generates valid hosts
		if _, err := regexp.Compile(rhs); err != nil {
			return cxt.NewError(InvalidRegex, "%v", err)
		}
	case "url_regex":
		// only makes sure the regex compiles, not that the regex generates valid urls
		if _, err := regexp.Compile(rhs); err != nil {
			return cxt.NewError(InvalidRegex, "%v", err)
		}
	default:
		return cxt.NewError(InvalidLabel)
	}

	return nil
}

func parseSecondarySpecifiers(lhs string, rhs string) test.Error {

	switch lhs {
	case "port":
		if err := is.Port.Validate(rhs); err != nil {
			return cxt.AddErrorCode(InvalidPort, err)
		}
	case "scheme":
		if rhs != "http" && rhs != "https" {
			return cxt.NewError(InvalidHTTPScheme)
		}
	case "prefix":
		// idk what validation to do on this
		// does a path prefix contain '/' at the start of it?
		// ignore for now..
		//	Same cross platform problem as below?
	case "suffix":
		// examples: gif jpeg
		// pure syntax: xxx.1 is a valid file name
		//	I doubt there is anything in a validation package for this.
		//	Even if there was, it would be silly since different platforms
		//	have difference specifications for file suffixes.
	case "method":
		// assuming all methods are valid
		// RFC 2616-9 for list of all methods
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
		default:
			return cxt.NewError(InvalidMethod, "invalid method \"%v\"", rhs)
		}

	case "time":
		// Invalid24HrTimeRange
		// example: 08:00-14:00
		// time range using 24-hour clock
		// I think it is shown how to parse this in ATS? or was that for pin-in-cache?

	case "src_ip":
		if err := is.IP.Validate(rhs); err != nil {
			return cxt.AddErrorCode(InvalidIP, err)
		}
	case "internal":
		if rhs != "true" && rhs != "false" {
			return cxt.NewError(InvalidBool)
		}

	default:
		return cxt.NewError(InvalidLabel)
	}

	return nil
}

func parseActions(lhs string, rhs string) test.Error {

	switch lhs {
	// traffic_ctl gives us an error for THIS
	// it kinda makes sense since it is easy to verify..
	case "action":
		// is the RHS case sensitive?
		switch rhs {
		case "never-cache":
		case "ignore-no-cache":
		case "ignore-client-no-cache":
		case "ignore-server-no-cache":
		default:
			return cxt.NewError(InvalidAction)
		}

	case "cache-responses-to-cookies":
		digit := rhs[0]
		if digit < '0' || '4' > digit || len(rhs) > 1 {
			cxt.NewError(InvalidCacheResponseType)
		}

	// I assume the order matters. I can double check. Look at notes.
	// All of these are time formats (eg. 1d4h)
	case "pin-in-cache":
		// cxt.NewError(InvalidHDMSTimeFmt)
	case "revalidate":
	case "ttl-in-cache":
	default:
		return cxt.NewError(InvalidLabel)
	}

	return nil
}

func parseCacheConfigRule(rule string) test.Error {

	var primaryDestination bool
	var match []string
	var err test.Error

	// no individual secondary specifier label can be used twice
	var count = map[string]int{
		"port":     0,
		"scheme":   0,
		"prefix":   0,
		"suffix":   0,
		"method":   0,
		"time":     0,
		"src_ip":   0,
		"internal": 0,
	}

	// Does ATS also delimit with tabs?
	assignments := strings.Split(rule, " ")
	last := len(assignments) - 1
	if last < 1 {
		return cxt.NewError(NotEnoughAssignments)
	}

	// neither the rhs or lhs can contain any whitespace
	assignment := regexp.MustCompile(`([a-z_\d]+)=(\S+)`)
	for _, elem := range assignments {
		match = assignment.FindStringSubmatch(elem)
		if match == nil {
			return cxt.NewError(BadAssignmentMatch, "could not match assignment: \"%v\"", elem)
		}

		// convert lhs to lowercase
		err = parsePrimaryDestinations(match[1], match[2])
		if err != nil && err.Code() != InvalidLabel {
			return err.Prepend("coult not parse primary destination from (%s): ", match[0])
		} else if primaryDestination {
			return cxt.NewError(ExcessLabel, "too many primary destinations")
		} else {
			primaryDestination = true
		}

		err = parseSecondarySpecifiers(match[1], match[2])
		if err != nil && err.Code() != InvalidLabel {
			return err.Prepend("could not parse secondary specifiers from (%s): ", match[0])
		} else {
			count[match[1]]++
			if count[match[1]] == 2 {
				return cxt.NewError(ExcessLabel, "the label \"%s\" can only be used once per rule", match[1])
			}
		}

		err = parseActions(match[1], match[2])
		if err != nil && err.Code() != InvalidLabel {
			return err.Prepend("could not parse action from (%s): ", match[0])
		}

		if err != nil && err.Code() == InvalidLabel {
			return err.Prepend("unknown label: ")
		}
	}

	return nil
}
