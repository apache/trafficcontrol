package cache_config

import (
	"regexp"
	"strings"

	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/test"
	"github.com/go-ozzo/ozzo-validation/is"
)

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

func parsePrimaryDestinations(lhs string, rhs string) test.Error {

	switch lhs {
	case "dest_domain":
		// dest_host is an alias for dest_domain
		fallthrough
	case "dest_host":
		if err := is.Host.Validate(rhs); err != nil {
			return ErrorContext.NewError(InvalidHost, "\"%s\" %v", rhs, err)
		}
	case "dest_ip":
		if err := is.IP.Validate(rhs); err != nil {
			return ErrorContext.NewError(InvalidIP, "\"%s\" %v", rhs, err)
		}
	case "host_regex":
		// only makes sure the regex compiles, not that the regex generates valid hosts
		if _, err := regexp.Compile(rhs); err != nil {
			return ErrorContext.NewError(InvalidHostRegex, "%v", err)
		}
	case "url_regex":
		// only makes sure the regex compiles, not that the regex generates valid urls
		if _, err := regexp.Compile(rhs); err != nil {
			return ErrorContext.NewError(InvalidURLRegex, "%v", err)
		}
	default:
		return ErrorContext.NewError(InvalidDestinationLabel)
	}

	return nil
}

func parseSecondarySpecifiers(lhs string, rhs string) test.Error {

	switch lhs {
	case "port":
		if err := is.Port.Validate(rhs); err != nil {
			return ErrorContext.AddErrorCode(InvalidPort, err)
		}
	case "scheme":
		if rhs != "http" && rhs != "https" {
			return ErrorContext.NewError(InvalidScheme)
		}
	case "prefix":
		// idk what validation to do on this
		// does a path prefix contain '/' at the start of it?
	case "suffix":
		// examples: gif jpeg
		// should I do a lookup of particular suffixes, or should I allow everything?
		// I think xxx.1 is a valid file name
		// what restrictions are there?
		// should I use an official validation package?
	case "method":
		// assuming all methods are valid
		// RFC link?
		// assuming both all uppercase and all lowercase should be valid
		// TODO convert to lowercase
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
			return ErrorContext.NewError(InvalidMethod, "invalid method \"%v\"", rhs)
		}

	case "time":
		// example: 08:00-14:00
		// time range using 24-hour clock
		// I think it is shown how to parse this in ATS? or was that for pin-in-cache?

	case "src_ip":
		if err := is.IP.Validate(rhs); err != nil {
			return ErrorContext.AddErrorCode(InvalidSrcIP, err)
		}
	case "internal":
		if rhs != "true" && rhs != "false" {
			return ErrorContext.NewError(0, "internal label must have a value of 'true' or 'false'")
		}

	default:
		return ErrorContext.NewError(InvalidSpecifierLabel)
	}

	return nil
}

func parseActions(lhs string, rhs string) test.Error {

	switch lhs {
	// traffic_ctl gives us an error for THIS
	// it kinda makes sense since it is easy to verify..
	case "action":
		// can only be one of 'never-cache', 'ignore-no-cache', 'ignore-client-no-cache', 'ignore-server-no-cache'
	case "cache-responses-to-cookies":
		// can only be one of 0, 1, 2, 3, or 4
	case "pin-in-cache": //1d4h (time format)
	case "revalidate": //time format
	case "ttl-in-cache": // I ASSUME the order matters. I can double check in the source code. I found it once before. Look at note sheet.
	default:
		return ErrorContext.NewError(InvalidAction)
	}

	return nil
}

func parseCacheConfigRule(rule string) test.Error {

	var match []string
	var err test.Error

	// Basically need to ensure there is at least 1 of the first 5 and not 2 of the second blob
	var label = map[string]int{

		// consider dropping
		// consider using for better error message
		"dest_domain": 0,
		"dest_host":   0,
		"dest_ip":     0,
		"host_regex":  0,
		"url_regex":   0,

		// count for individual items cannot be 2
		"port":     0,
		"scheme":   0,
		"prefix":   0,
		"suffix":   0,
		"method":   0,
		"time":     0,
		"src_ip":   0,
		"internal": 0,
	}
	// map parseXxx to "internal" and other keys..

	assignments := strings.Split(rule, " ")
	last := len(assignments) - 1
	if last < 1 {
		return ErrorContext.NewError(NotEnoughAssignments)
	}

	assignment := regexp.MustCompile(`([a-z_\d]+)=(\S+)`)

	for _, elem := range assignments {
		match = assignment.FindStringSubmatch(elem)
		if match == nil {
			return ErrorContext.NewError(BadAssignmentMatch, "could not match assignment: \"%v\"", elem)
		}

		// not sure if this is safe
		label[match[1]]++

		//if err = parsePrimaryDestinations(match[1], match[2]); err != nil {
		//if err = parseActions(match[1], match[2]); err != nil {
		// convert lhs to lowercase
		if err = parseSecondarySpecifiers(match[1], match[2]); err != nil {
			return err.Prepend("could not parse secondary specifiers from (%s): ", match[0])
		}
	}

	return nil
}
