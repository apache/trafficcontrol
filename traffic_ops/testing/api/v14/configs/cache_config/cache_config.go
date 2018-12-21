package cache_config

import (
	"regexp"
	"strings"

	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/test" // Error interface
	"github.com/go-ozzo/ozzo-validation/is"
)

func parseCacheConfig(config string) test.Error {
	lines := strings.Split(config, "\n")

	// if the config is only one line, don't add line error information
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
			return ErrorContext.NewError(INVALID_HOST, "\"%s\" %v", rhs, err)
		}

	case "dest_ip":
		//return fmt.Errorf("invalid ip: \"%s\"", rhs)

		// assuming both IP4 and IP6 are allowed and that there are no ranges or wildcards (unlike src_ip in ip_allow.config)
		// I'm not even sure what is.IP checks other than IP4 and IP6
		// Also, I'm not sure what the difference is between src_ip and dest_ip 100%, what does source and destination refer to in this case?
		if err := is.IP.Validate(rhs); err != nil {
			return ErrorContext.NewError(INVALID_IP, "\"%s\" %v", rhs, err)
		}

	case "host_regex":
		// right now I only make sure the regex compiles, not that the regex generates valid hosts
		if _, err := regexp.Compile(rhs); err != nil {
			return ErrorContext.NewError(INVALID_HOST_REGEX, "%v", err)
		}
	case "url_regex":
		// right now I only make sure the regex compiles, not that the regex generates valid urls
		if _, err := regexp.Compile(rhs); err != nil {
			return ErrorContext.NewError(INVALID_URL_REGEX, "%v", err)
		}

	default:
		// todo: if no message is given, use a map to get associated string?
		// todo: NewError is from a closure function
		// the closure must verify that all codes are valid and provide mappings for default messages
		// for simplicity of implementation, I can also consider using a struct instead
		// errorContext.New(INVALID_DESTINATION)
		return ErrorContext.NewError(INVALID_DESTINATION_LABEL, "invalid field given for primary destination")
	}

	return nil
}

func parseSecondarySpecifiers(lhs string, rhs string) test.Error {

	// note: `scheme://host:port/path_prefix`
	switch lhs {
	case "port":
		if err := is.Port.Validate(rhs); err != nil {
			return ErrorContext.AddErrorCode(22, err)
		}

	case "scheme":
		if rhs != "http" && rhs != "https" {
			return ErrorContext.NewError(23, "invalid scheme (must be either http or https)")
		}

	case "prefix":
		// this (from the remap.config docs) shows that .html can exist as a path_prefix
		//map http://example.com/dist_get_user http://127.0.0.1:8001/denied.html
		// not sure how to validate this, all I know is that it has no spaces
		// worst case, here is a reference: https://tools.ietf.org/html/rfc3986#section-1.1
		// golang's net/url package has: [scheme:][//[userinfo@]host][/]path[?query][#fragment]
		// seems to indicate I don't worry about [/] or query or fragment
		// thinking it makese sense to add prefix to `www.test.com/`
		// grep for prefix (cache.config examples specifically)
	case "suffix":
		// examples: gif jpeg
		// should I do a lookup of particular suffixes, or should I allow everything?
		// I think xxx.1 is a valid file name
		// what restrictions are there?
		// should I use an official validation package?
		// just say no special characters? [^/\?%*:|"<>. ]
		// would be safe to say no whitespace is allowed
		// [\d\w] not sure if this is what I need..
		// now that I think of it I should allow both upper and lowercase

	case "method":
		// assuming all methods are valid
		// not sure if any of these inherently don't make sense
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
			return ErrorContext.NewError(0, "invalid method \"%v\"", rhs)
		}

	case "time":
		// example: 08:00-14:00
		// time range using 24-hour clock
		// I think it is shown how to parse this in ATS? or was that for pin-in-cache?

	case "src_ip":
		// client ip address
		// not sure if this uses range
		if err := is.IP.Validate(rhs); err != nil {
			return ErrorContext.AddErrorCode(0, err)
		}

	case "internal":
		if rhs != "true" && rhs != "false" {
			return ErrorContext.NewError(0, "internal label must have a value of 'true' or 'false'")
		}

	default:
		return ErrorContext.NewError(21, "invalid field given for secondary specifier")
	}

	return nil
}

func parseActions(lhs string, rhs string) test.Error {

	switch lhs {
	case "action":
	case "cache-responses-to-cookies":
	case "pin-in-cache": //1h
	case "revalidate":
	case "ttl-in-cache": //1d2h
	default:
		return ErrorContext.NewError(31, "invalid field given for action")
	}

	return nil
}

func parseCacheConfigRule(rule string) test.Error {

	var match []string
	var err test.Error

	/*
		inside parseCacheConfigRule, we only need to count secondary specifiers to make sure none is made twice
		perhaps it would be helpful to count primary destinations as well to give a more specific error

		a separate function should be made to return label counts explicitly
		the label counts should not be returned by the parser so that it has a single responsibility
		the label counts would be used to track postive test coverage (should be exclusive to tests)
	*/

	var label = map[string]int{
		"dest_domain": 0,
		"dest_host":   0,
		"dest_ip":     0,
		"host_regex":  0,
		"url_regex":   0,
		"port":        0,
		"scheme":      0,
		"prefix":      0,
		"suffix":      0,
		"method":      0,
		"time":        0,
		"src_ip":      0,
		"internal":    0,
	}

	assignments := strings.Split(rule, " ")
	last := len(assignments) - 1
	if last < 1 {
		return ErrorContext.NewError(NOT_ENOUGH_ASSIGNMENTS)
	}

	// neither the rhs or the lhs should contain a space
	assignment := regexp.MustCompile(`([a-z_\d]+)=(\S+)`)

	// head is primary specifier, body is list of secondary specifiers, tail is the action
	// not actually sure order is supposed to matter in the list of rules, but not assuming an order might give me a better error message
	head, body, tail := assignments[0], assignments[1:last], assignments[last]

	match = assignment.FindStringSubmatch(head)
	if match == nil {
		// consider this: the error message does not specify 'primary desination' specifically
		// then ErrorContext.NewError can also be made to create (errcode, offender)
		// only not sure if adding quotes makes sense to do for every offender
		return ErrorContext.NewError(BAD_ASSIGNMENT_MATCH, "could not match primary destination assignment: \"%v\"", head)
	}

	label[match[1]]++
	if err = parsePrimaryDestinations(match[1], match[2]); err != nil {

		//if err.ErrorCode() == 21 {
		// if we have two primary specifiers or one action the user has the wrong order
		//}

		return err.Prepend("could not parse primary destination: ")
	}

	for _, elem := range body {
		match = assignment.FindStringSubmatch(elem)
		if match == nil {
			return ErrorContext.NewError(BAD_ASSIGNMENT_MATCH, "could not match secondary specifier assignment: \"%v\"", elem)
		}

		label[match[1]]++
		if err = parseSecondarySpecifiers(match[1], match[2]); err != nil {
			// hm.. not sure I like the name
			// just change to prepend? make sense to me
			return err.Prepend("could not parse secondary specifiers from (%s): ", match[0])
		}
	}

	match = assignment.FindStringSubmatch(tail)
	if match == nil {
		return ErrorContext.NewError(BAD_ASSIGNMENT_MATCH, "could not match action assignment: \"%v\"", tail)
	}

	label[match[1]]++
	if err = parseActions(match[1], match[2]); err != nil {
		return err.Prepend("could not parse action from (%s): ", match[0])
	}

	return nil
}
