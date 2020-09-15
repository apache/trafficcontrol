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

package config

import "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/test"

// Error codes:
const (
	BadAssignmentMatch = iota + 1
	NotEnoughAssignments
	ExcessLabel
	InvalidLabel
	MissingLabel
	InvalidAction
	InvalidBool
	InvalidCacheCookieResponse
	InvalidHTTPScheme
	InvalidHost
	InvalidIP
	UnknownMethod
	InvalidIPRange
	InvalidPort
	InvalidRegex
	InvalidTimeFormatDHMS
	InvalidTimeRange24Hr
)

// ErrorContext contains the error codes mentioned above.
// Any error made must have one of those error codes.
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
		UnknownMethod,
		InvalidIPRange,
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
