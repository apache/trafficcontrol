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

type errorCode struct {
	error
	cause error
	code  int
}

type ErrorCoder interface {
	Error() string
	ErrorCode() int
	Cause() error
	Prepend(string, ...interface{}) ErrorCoder
}

func NewError(code int, fmtStr string, fmtArgs ...interface{}) ErrorCoder {
	err := fmt.Errorf(fmtStr, fmtArgs...)
	return &errorCode{err, err, code}
}

func (e errorCode) ErrorCode() int {
	return e.code
}

// not a pointer receiver, does not modify itself
// in the new error the original cause is maintained
func (e errorCode) Prepend(fmtStr string, fmtArgs ...interface{}) ErrorCoder {
	err := fmt.Errorf(fmtStr, fmtArgs...)
	e.error = fmt.Errorf("%v %v", err, e.error)
	return e
}

func (e errorCode) Cause() error {
	return e.cause
}

func AddErrorCode(code int, err error) ErrorCoder {
	return NewError(code, "%v", err)
}

// ErrorContext regulates which error codes can be made and keeps a count
// of all the errors created through the context
//
// ErrorContext.NewError can be used like test.NewError
//
// Implementation details:
// contains a list of all valid error codes
//		- allows user to make sure they are creating the correct errors
//		- actually a map
//			lookup can be done without linear search
//			we can use the map to keep count of which errors are made
// contains mapping from error code to name (either for testing metainfo or used in case no args are given)
//		- not required for all error codes, or for any
//
// Recommended setup:
//
//	package some_regex_checker
//
//	const (
//		COMMON_BASE            = 10 + iota
//		NOT_ENOUGH_ASSIGNMENTS
//		BAD_ASSIGNMENT_MATCH
//		...
//	)
//
//	// scoped to the package name
//	var ErrorContext *test.ErrorContext
//
//	func init() {
//		errorCodes := []uint{
//			NOT_ENOUGH_ASSIGNMENTS,
//			BAD_ASSIGNMENT_MATCH,
//		}
//		ErrorContext = test.NewErrorContext("cache config", errorCodes)
//		err := ErrorContext.AddMapping(NOT_ENOUGH_ASSIGNMENTS, "not enough assignments in rule")
//		// check err
//		ErrorContext.TurnPanicOn()
//		// in non-init function create errors with the context
//	}
//
type ErrorContext struct {
	calledNewError  bool
	createdNewError bool
	doPanic         bool
	name            string
	codes           map[uint]uint
	description     map[uint]string
}

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

// although golang panics are highly discouraged, panic mode is made as an option
// made this decision partially because type assertions and map membership have similar options
// if a user doesn't have panic mode on, they should still terminate after running into an error
// panic is off by default, and must be turned on explicitly so that the user must make an active decision
// panic must be turned on before errors are created
// once turned on the panic mode can't be turned off
func (ec *ErrorContext) TurnPanicOn() error {
	if ec.calledNewError {
		return ec.internalError(BAD_INIT_TIMING, nil)
	}
	ec.doPanic = true
	return nil
}

// AddMapping gives a default message for a given error code.
// Mappings should be configured before errors are called for consistency.
//
// usage: ErrorContext.AddMapping(404, "not found")
//        err := ErrorContext.NewError(404)
//        err := ErrorContext.NewError(404, "not found: %v", prev_err")
// the err automatically has a default message
// the default message can be overridden with something else to add context to the error
//
// parameters:
//  `code` should exist in the error context, only one mapping can exist per error context
//  `msg` should be a plain string without special formatting
func (ec *ErrorContext) AddMapping(code uint, msg string) error {
	if ec.calledNewError {
		return ec.internalError(BAD_INIT_TIMING, nil)
	}
	if !ec.whitelisted(code) {
		return ec.internalError(BAD_MAP_CREATE, code)
	}
	if _, exists := ec.description[code]; exists {
		return ec.internalError(BAD_DUP_MAPPING, code)
	}
	ec.description[code] = msg
	return nil
}

// AddMap applies the `AddMapping` function for every element in the given map
// the function argument will not override the current contents of the map, everything is additive
// an error is returned if a duplicate code is added
func (ec *ErrorContext) AddMap(mapping map[uint]string) error {
	for code, desc := range mapping {
		if err := ec.AddMapping(code, desc); err != nil {
			return err
		}
	}
	return nil
}

// all internal errors
const (
	BAD_ERROR_CODE  = -1 - iota // when bad error code is given in creation of new error
	BAD_MAP_LOOKUP              // when creating an error with no error message, default message wasn't found
	BAD_MAP_CREATE              // when bad error code is given in creation of map
	BAD_DUP_MAPPING             // when a mapping is already made
	BAD_FMT_STRING              // when the fmt string isn't a string
	BAD_INIT_TIMING             // when the error context is modifed after having created an error
)

func (ec ErrorContext) internalError(code int, offender interface{}) ErrorCoder {

	var err error
	switch code {
	case BAD_ERROR_CODE:
		err = fmt.Errorf("error code %v not found in whitelist for \"%s\" error context", offender, ec.name)
	case BAD_MAP_CREATE:
		err = fmt.Errorf("when creating default error mapping, code %v wasn't found in the code whitelist for the \"%s\" error context", offender, ec.name)
	case BAD_DUP_MAPPING:
		err = fmt.Errorf("when creating default error mapping, code %v was already found as a mapping for the \"%s\" error context", offender, ec.name)
	case BAD_MAP_LOOKUP:
		err = fmt.Errorf("bad default error lookup for code %v in the \"%s\" error context", offender, ec.name)
	case BAD_FMT_STRING:
		err = fmt.Errorf("the leading argument \"%v\" could not be interpreted as a format string", offender)
	case BAD_INIT_TIMING:
		err = fmt.Errorf("tried to modify error context after creating error")
	}

	if ec.doPanic {
		panic(err)
	}
	return AddErrorCode(code, err)
}

// whitelist is defined by code map membership
func (ec ErrorContext) whitelisted(code uint) bool {
	_, ok := ec.codes[code]
	return ok
}

// makes sure every created error gets counted
func (ec *ErrorContext) newError(code uint, fmtStr string, fmtArgs ...interface{}) ErrorCoder {
	ec.codes[code]++
	ec.createdNewError = true
	return NewError(int(code), fmtStr, fmtArgs...)
}

// ErrorContext.NewError creates an error similar to test.NewError
// any error created in this manner must have a code that belongs to the error context
// the args is interpreted as a format string with args, but I'm using a ...interface{} because
// if no args are specified a lookup will be made to find the default configured string (see `AddMapping`)
//
// usage: NewError(404, "not found: %v", prev_err)
//        NewError(404) // (mapping must exist otherwise this is an error)
func (ec *ErrorContext) NewError(code uint, args ...interface{}) ErrorCoder {
	ec.calledNewError = true

	if ec.whitelisted(code) {

		// if no args given, try to find in mapping
		if len(args) == 0 {
			if errDesc, ok := ec.description[code]; ok {
				return ec.newError(code, errDesc)
			}
			return ec.internalError(BAD_MAP_LOOKUP, code)
		}

		// if args given, interpret as (fmtStr, fmtArgs)
		if fmtStr, ok := args[0].(string); ok {
			return ec.newError(code, fmtStr, args[1:]...)
		}
		return ec.internalError(BAD_FMT_STRING, args[0])
	}
	return ec.internalError(BAD_ERROR_CODE, code)
}

// AddErrorCode takes a regular error and gives it a code
// since this method is scoped to an error context, the code must exist in the whitelist
func (ec ErrorContext) AddErrorCode(code uint, err error) ErrorCoder {
	return ec.NewError(code, "%v", err)
}

// GetErrorStats returns the map of error codes
// ec.codes[code] is the number of times an error with the indexed code has been run
func (ec ErrorContext) GetErrorStats() map[uint]uint {
	if ec.createdNewError {
		return ec.codes
	}
	return nil
}
