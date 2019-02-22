package config

import "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/test"

const (
	BadAssignmentMatch = iota + 10
	NotEnoughAssignments
	ExcessLabel
	InvalidLabel
	MissingLabel
)

const (
	InvalidAction = iota + 20
	InvalidBool
	InvalidCacheCookieResponse
	InvalidHTTPScheme
	InvalidHost
	InvalidIP
	InvalidMethod
	InvalidPort
	InvalidRegex
	InvalidTimeFormatDHMS
	InvalidTimeRange24Hr
)

var ErrorContext *test.ErrorContext

func init() {
	iterableErrorCodes := []uint{
		BadAssignmentMatch,
		NotEnoughAssignments,
		ExcessLabel,
		InvalidLabel,
		MissingLabel,
		InvalidAction,
		InvalidBool,
		InvalidCacheCookieResponse,
		InvalidHTTPScheme,
		InvalidHost,
		InvalidIP,
		InvalidMethod,
		InvalidPort,
		InvalidRegex,
		InvalidTimeFormatDHMS,
		InvalidTimeRange24Hr,
	}

	ErrorContext = test.NewErrorContext("cache config", iterableErrorCodes)

	ErrorContext.SetDefaultMessageForCode(InvalidLabel,
		"invalid label")
	ErrorContext.SetDefaultMessageForCode(InvalidAction,
		"invalid action")
	ErrorContext.SetDefaultMessageForCode(NotEnoughAssignments,
		"not enough assignments in rule")
	ErrorContext.SetDefaultMessageForCode(InvalidHTTPScheme,
		"invalid scheme (must be either http or https)")
	ErrorContext.SetDefaultMessageForCode(InvalidBool,
		"label must have a value of 'true' or 'false'")
	ErrorContext.SetDefaultMessageForCode(InvalidCacheCookieResponse,
		"Value for cache-responses-to-cookies must be an integer in the range 0..4")

	ErrorContext.TurnPanicOn()
}
