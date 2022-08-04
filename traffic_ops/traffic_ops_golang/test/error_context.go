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

package test

import "fmt"

type errorImpl struct {
	error
	cause error
	code  int
}

// Error describes an interface that supports error codes.
// It also allows for the first (causal) error to be easily referenced without needing to parse it.
// Adding error context is supported by prepending a string, and keeping the current error code.
//
// usage:
//
//	func func1() Error {
//		return test.NewError(404, "not found")
//	}
//
//	func func2() Error {
//		err := func1()
//		if err != nil {
//			return err.Prepend("while in func2: ")
//		}
//	}
//
//	func main() {
//		err := func1()
//		err.Code()  // 404
//		err.Cause() // not found
//	}
type Error interface {
	Error() string
	Code() int
	Cause() error
	Prepend(string, ...interface{}) Error
}

// NewError constructs an error with a code.
func NewError(code int, fmtStr string, fmtArgs ...interface{}) Error {
	err := fmt.Errorf(fmtStr, fmtArgs...)
	return &errorImpl{err, err, code}
}

// Code returns the integer error code that the error message is associated with.
func (e errorImpl) Code() int {
	return e.code
}

// Prepend prepends the string built from fmtStr and fmtArgs to the error and returns it.
// Note is is not a pointer receiver, and does not modify itself. In the new error the
// original cause is maintained.
func (e errorImpl) Prepend(fmtStr string, fmtArgs ...interface{}) Error {
	err := fmt.Errorf(fmtStr, fmtArgs...)
	e.error = fmt.Errorf("%v %v", err, e.error)
	return e
}

// Cause returns the original error made, without extra context. It implements
// the causer interface defined in the errors package.
// see:
//
//	https://pkg.go.dev/github.com/pkg/errors#Cause
func (e errorImpl) Cause() error {
	return e.cause
}

// AddErrorCode takes an error and returns an instance satisfying the Error interface.
func AddErrorCode(code int, err error) Error {
	return NewError(code, "%v", err)
}

// ErrorContext regulates which error codes can be made and keeps a count
// of all the errors created through the context.
//
// `ErrorContext.NewError` can be used like `test.NewError` (see Error above)
// The primary difference is that the context prevents non-whitelisted error codes from being made.
//
// Implementation details:
//
//	 contains a list of all valid error codes
//		- allows user to make sure they are creating the correct errors
//		- actually a map
//			lookup can be done without linear search
//			we can use the map to keep count of which errors are made
//	 contains mapping from error code to name/default message
//		- not required for all error codes, or for any
//
// An example setup:
//
//	package some_regex_checker
//
//	const (
//		CommonBase            = 10 + iota
//		NotEnoughAssignments
//		BadAssignmentMatch
//		...
//	)
//
//	// scoped to the package name
//	var ErrorContext *test.ErrorContext
//
//	func init() {
//		errorCodes := []uint{
//			NotEnoughAssignments,
//			BadAssignmentMatch,
//		}
//
//		ErrorContext = test.NewErrorContext("cache config", errorCodes)
//		err := ErrorContext.SetDefaultMessageForCode(NotEnoughAssignments, "not enough assignments in rule")
//		// check err
//
//		ErrorContext.TurnPanicOn()
//	}
//
//	func main() {
//		// create a new user error with the context like this
//		err := ErrorContext.NewError(BadAssignmentMatch, "bad assignment match")
//
//		// there is no error code with 0, so this panics
//		err = ErrorContext.NewError(0, "some error msg")
//	}
type ErrorContext struct {
	calledNewError  bool
	createdNewError bool
	doPanic         bool
	name            string

	// codes is both a whitelist of codes and a count of which error codes have been called.
	codes map[uint]uint

	// description is a map from an error code to a default error message.
	description map[uint]string
}

// NewErrorContext constructs an error context with list of valid codes.
func NewErrorContext(contextName string, errCodes []uint) *ErrorContext {

	codeMap := make(map[uint]uint)
	for _, code := range errCodes {
		codeMap[code] = 0
	}

	descMap := make(map[uint]string)

	return &ErrorContext{
		calledNewError:  false,
		createdNewError: false,
		doPanic:         false,
		name:            contextName,
		codes:           codeMap,
		description:     descMap,
	}
}

// TurnPanicOn enables panic mode.
//
// When panic mode is on, ErrorContext methods that return errors will panic.
// Panic mode does not affect user-created errors. Panic mode can be used to
// assert the error context is set up correctly.
//
// Although golang panics are highly discouraged, panic mode is made as an option.
// This decision was partially made because type assertions and map membership have
// similar options. If a user doesn't have panic mode on, they should still terminate
// after running into an error. Panic is off by default, and must be turned on explicitly
// so that the user must make an active decision. Panic must be turned on before errors
// are created.
//
// Once turned on, panic mode can't be turned off.
func (ec *ErrorContext) TurnPanicOn() error {
	if ec.calledNewError {
		return ec.internalError(BadInitOrder, nil)
	}
	ec.doPanic = true
	return nil
}

// SetDefaultMessageForCode gives a default message for a given error code.
// Default messages must be configured before errors are created.
//
// parameters:
//
//	`code` should exist in the error context, only one default message mapping can exist per error context
//	`msg` should be a plain string without special formatting
//
// usage:
//
//	ErrorContext.SetDefaultMessageForCode(404, "not found")
//
//	// err has a default message
//	err := ErrorContext.NewError(404)
//
//	// the default message is overridden to add context to the error
//	err := ErrorContext.NewError(404, "not found: %v", prev_err")
func (ec *ErrorContext) SetDefaultMessageForCode(code uint, msg string) error {
	if ec.calledNewError {
		return ec.internalError(BadInitOrder, nil)
	}
	if !ec.whitelisted(code) {
		return ec.internalError(BadErrorCode, code)
	}
	if _, exists := ec.description[code]; exists {
		return ec.internalError(BadDupMessage, code)
	}
	ec.description[code] = msg
	return nil
}

// AddDefaultMessages applies the SetDefaultMessageForCode function for every element in the given map.
// The function does not override the current contents of the map, everything is additive. An error is
// returned if a duplicate code is added.
//
// parameter:
//
//	`mapping` is a map of error codes to their default error messages
func (ec *ErrorContext) AddDefaultErrorMessages(mapping map[uint]string) error {
	for code, desc := range mapping {
		if err := ec.SetDefaultMessageForCode(code, desc); err != nil {
			return err
		}
	}
	return nil
}

type ErrorContextInternalErrorCode int

// All internal errors for ErrorContext
const (
	iotaRef       ErrorContextInternalErrorCode = -iota
	BadErrorCode                                // when a bad error code is given in creation of new error
	BadDupMessage                               // when a default message already exists
	BadMsgLookup                                // when creating an error with no error message, default message wasn't found
	BadFmtString                                // when the fmt string isn't a string
	BadInitOrder                                // when the error context is modifed after having created an error
)

// internalError returns an error with a message depending on the error code given.
// The offender is the single item that the error is blamed on. Most offenders are an incorrect code.
func (ec ErrorContext) internalError(code ErrorContextInternalErrorCode, offender interface{}) Error {

	var err error
	switch code {
	case BadErrorCode:
		err = fmt.Errorf(`error code %v not found in whitelist for "%s" error context`, offender, ec.name)
	case BadDupMessage:
		err = fmt.Errorf(`code %v already has a default message for the "%s" error context`, offender, ec.name)
	case BadMsgLookup:
		err = fmt.Errorf(`bad default error lookup for code %v in the "%s" error context`, offender, ec.name)
	case BadFmtString:
		err = fmt.Errorf(`the leading argument "%v" could not be interpreted as a format string`, offender)
	case BadInitOrder:
		err = fmt.Errorf(`tried to modify error context after creating error`)
	}

	if ec.doPanic {
		panic(err)
	}
	return AddErrorCode(int(code), err)
}

// whitelisted is defined by code map membership
func (ec ErrorContext) whitelisted(code uint) bool {
	_, ok := ec.codes[code]
	return ok
}

// newError makes sure every created error gets counted
func (ec *ErrorContext) newError(code uint, fmtStr string, fmtArgs ...interface{}) Error {
	ec.codes[code]++
	ec.createdNewError = true
	return NewError(int(code), fmtStr, fmtArgs...)
}

// NewError for an ErrorContext creates an error similar to test.NewError
// Any error created in this manner must have a code that belongs to the error context.
// The args is interpreted as a format string with args, but `...interface{}` is used because
// if no args are specified a lookup will be made to find the default configured string (see SetDefaultMessageForCode)
//
// usage:
//
//	cxt.NewError(404, "not found: %v", prev_err)
//	cxt.NewError(404) // (default message must exist otherwise this is an error)
func (ec *ErrorContext) NewError(code uint, args ...interface{}) Error {
	ec.calledNewError = true

	if ec.whitelisted(code) {

		// if no args given, try to find a default
		if len(args) == 0 {
			if errDesc, ok := ec.description[code]; ok {
				return ec.newError(code, errDesc)
			}
			return ec.internalError(BadMsgLookup, code)
		}

		// if args given, interpret as (fmtStr, fmtArgs)
		if fmtStr, ok := args[0].(string); ok {
			return ec.newError(code, fmtStr, args[1:]...)
		}
		return ec.internalError(BadFmtString, args[0])
	}
	return ec.internalError(BadErrorCode, code)
}

// AddErrorCode takes a regular error and gives it a code.
// Since this method is scoped to an error context, the code must exist in the whitelist.
func (ec ErrorContext) AddErrorCode(code uint, err error) Error {
	return ec.NewError(code, "%v", err)
}

// GetErrorStats returns the map of error codes.
// usage:
//
//	stats := cxt.GetErrorStats()
//	stats[code] // represents the number of times an error with the code has been created
func (ec ErrorContext) GetErrorStats() map[uint]uint {
	if ec.createdNewError {
		return ec.codes
	}
	return nil
}
