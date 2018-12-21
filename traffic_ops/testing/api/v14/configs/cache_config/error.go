package cache_config

import "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/test"

// This would even be common to other configs
// You could get more, but right now these are organized by the function they are called in
const (
	COMMON_BASE = iota + 10
	NOT_ENOUGH_ASSIGNMENTS
	BAD_ASSIGNMENT_MATCH
)

// Primary Destination
const (
	PD_BASE = iota + 20
	INVALID_DESTINATION_LABEL
	INVALID_HOST
	INVALID_IP
	INVALID_HOST_REGEX
	INVALID_URL_REGEX
)

// Secondary Specifier
const (
	SS_BASE = iota + 30
	INVALID_SPECIFIER_LABEL
	INVALID_PORT
	INVALID_SCHEME
	INVALID_PREFIX
	INVALID_SUFFIX
	INVALID_METHOD
	INVALID_TIME
	INVALID_SRC_IP
	INVALID_INTERNAL
)

// Action
const (
	A_BASE = iota + 40
	INVALID_ACTION_LABEL
	INVALID_ACTION
	INVALID_RESP_TO_COOKIES
	INVALID_PIN_IN_CACHE
	INVALID_REVALIDATE
	INVALID_TTL_IN_CACHE
)

// scoped to the package name
var ErrorContext *test.ErrorContext

func init() {
	iterableErrorCodes := []uint{
		INVALID_DESTINATION_LABEL,
		INVALID_HOST,
		INVALID_IP,
		INVALID_HOST_REGEX,
		INVALID_URL_REGEX,
		INVALID_SPECIFIER_LABEL,
		INVALID_PORT,
		INVALID_SCHEME,
		INVALID_PREFIX,
		INVALID_SUFFIX,
		INVALID_METHOD,
		INVALID_TIME,
		INVALID_SRC_IP,
		INVALID_INTERNAL,
		INVALID_ACTION_LABEL,
		INVALID_ACTION,
		INVALID_RESP_TO_COOKIES,
		INVALID_PIN_IN_CACHE,
		INVALID_REVALIDATE,
		INVALID_TTL_IN_CACHE,
		NOT_ENOUGH_ASSIGNMENTS,
		BAD_ASSIGNMENT_MATCH,
	}
	// Not many mappings are made since most error add context
	ErrorContext = test.NewErrorContext("cache config", iterableErrorCodes)
	ErrorContext.SetDefaultMessageForCode(NOT_ENOUGH_ASSIGNMENTS, "not enough assignments in rule")
	ErrorContext.TurnPanicOn()
}
