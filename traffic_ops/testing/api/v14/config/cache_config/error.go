package cache_config

import "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/test"

const OK = 0

const (
	CommonBase = iota + 10
	NotEnoughAssignments
	BadAssignmentMatch
)

// Previously I made it so that I had three different sections..
// I think this is bad because there is obvious overlap between
// InvalidIP and InvalidSrcIP
//
// I think when I made it I assumed that just because something is
// valid in one config it doesn't necessarily mean it is valid in
// a different config.
//
// For instance, the scheme can only be http or https, but maybe
// somewhere else other schemes are allowed.

// ATS has such poor error communication itself I strongly think
// I should stop trying to emulate it. I should _simply_ allow
// things that are similar.. bring it up in the PR.
//
// If there is anything specific that sticks out I should be notified.
// I will follow the documentation for things like the scheme..
// That is easy and at least orthogonol to that documentation.

// A bad IP on a PD should be exactly like a bad IP on a SS,
// so there is really no need to test that. Thinking of permutations
// like that would be absolutely insane.

const (
// Things should go in here if they are common I guess...
)

const (
	PrimaryDestinationBase = iota + 20
	InvalidDestinationLabel
	InvalidHost
	InvalidIP
	InvalidHostRegex
	InvalidURLRegex
)

const (
	SecondarySpecifierBase = iota + 30
	InvalidSpecifierLabel
	InvalidPort
	InvalidScheme
	InvalidPrefix
	InvalidSuffix
	InvalidMethod
	InvalidTime
	InvalidSrcIP
	InvalidInternal
)

const (
	ActionBase = iota + 40
	InvalidActionLabel
	InvalidAction
	InvalidRespToCookies
	InvalidPinInCache
	InvalidRevalidate
	InvalidTTLInCache
)

var ErrorContext *test.ErrorContext

func init() {
	iterableErrorCodes := []uint{
		InvalidDestinationLabel,
		InvalidHost,
		InvalidIP,
		InvalidHostRegex,
		InvalidURLRegex,
		InvalidSpecifierLabel,
		InvalidPort,
		InvalidScheme,
		InvalidPrefix,
		InvalidSuffix,
		InvalidMethod,
		InvalidTime,
		InvalidSrcIP,
		InvalidInternal,
		InvalidActionLabel,
		InvalidAction,
		InvalidRespToCookies,
		InvalidPinInCache,
		InvalidRevalidate,
		InvalidTTLInCache,
		NotEnoughAssignments,
		BadAssignmentMatch,
	}
	ErrorContext = test.NewErrorContext("cache config", iterableErrorCodes)

	// This will need to be re-thought about slightly.
	// This error should now basically never occur.
	ErrorContext.SetDefaultMessageForCode(InvalidDestinationLabel, "invalid given for primary destination")
	ErrorContext.SetDefaultMessageForCode(InvalidSpeciferLabel, "invalid given for secondary specifier")
	ErrorContext.SetDefaultMessageForCode(InvalidAction, "invalid label for action")

	ErrorContext.SetDefaultMessageForCode(NotEnoughAssignments, "not enough assignments in rule")

	ErrorContext.SetDefaultMessageForCode(InvalidScheme, "invalid scheme (must be either http or https)")
	ErrorContext.TurnPanicOn()
}
