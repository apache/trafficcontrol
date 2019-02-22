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
			return err.Prepend("error on line %d: ", i+1)
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
		// ignore..
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
		if err := Validate24HrTimeRange(rhs); err != nil {
			return cxt.AddErrorCode(InvalidTimeRange24Hr, err)
		}
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
	case "action":
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
			return cxt.NewError(InvalidCacheCookieResponse)
		}

	// All of these are time formats
	case "pin-in-cache":
		fallthrough
	case "revalidate":
		fallthrough
	case "ttl-in-cache":
		err := ValidateDHMSTimeFormat(rhs)
		if err != nil {
			return cxt.AddErrorCode(InvalidTimeFormatDHMS, err)
		}
	default:
		return cxt.NewError(InvalidLabel)
	}

	return nil
}

func parseCacheConfigRule(rule string) test.Error {

	var destination bool
	var action bool
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

	if strings.Trim(rule, "\t ") == "" {
		return nil
	}

	assignments := strings.Split(rule, " ")
	last := len(assignments) - 1
	if last < 1 {
		return cxt.NewError(NotEnoughAssignments)
	}

	// neither the rhs or lhs can contain any whitespace
	assignment := regexp.MustCompile(`([a-z_\-\d]+)=(\S+)`)
	for _, elem := range assignments {
		match = assignment.FindStringSubmatch(strings.ToLower(elem))
		if match == nil {
			return cxt.NewError(BadAssignmentMatch, "could not match assignment: \"%v\"", elem)
		}

		err = parsePrimaryDestinations(match[1], match[2])
		if err == nil {
			if destination {
				return cxt.NewError(ExcessLabel, "too many primary destination labels")
			} else {
				destination = true
				continue
			}
		}
		if err.Code() != InvalidLabel {
			return err.Prepend("coult not parse primary destination from (%s): ", match[0])
		}

		err = parseSecondarySpecifiers(match[1], match[2])
		if err == nil {
			count[match[1]]++
			if count[match[1]] == 2 {
				return cxt.NewError(ExcessLabel, "the label \"%s\" can only be used once per rule", match[1])
			}
			continue
		}
		if err.Code() != InvalidLabel {
			return err.Prepend("could not parse secondary specifiers from (%s): ", match[0])
		}

		err = parseActions(match[1], match[2])
		if err == nil {
			action = true
			continue
		}

		if err.Code() == InvalidLabel {
			return err
		} else {
			return err.Prepend("could not parse action from (%s): ", match[0])
		}

	}

	if !destination {
		return cxt.NewError(MissingLabel, "missing primary destination label")
	}

	if !action {
		return cxt.NewError(MissingLabel, "missing action lablel")
	}

	return nil
}
