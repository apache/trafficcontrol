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

import (
	"fmt"
	"testing"
)

var codes = []uint{
	0, 1, 2,
}

func basicMatch(t *testing.T, err error, expected string) {
	if err.Error() != expected {
		t.Errorf("\nexpected: %v\nactual: %v\n\n", expected, err)
	}
}

func TestErrorCode(t *testing.T) {
	// test basic formatting
	basicMatch(t, NewError(10, "no args"), "no args")
	basicMatch(t, NewError(11, "arg is %v", "awesome"), "arg is awesome")

	// assert error code works
	expected := 12
	err := NewError(expected, "this works")
	if err.Code() != expected {
		t.Errorf("\nexpected: %d\nactual: %d\n\n", expected, err.Code())
	}

	// test prepending errors and retaining the original error
	err = err.Prepend("note:")
	basicMatch(t, err, "note: this works")
	basicMatch(t, err.Cause(), "this works")
}

func TestErrorContextPanicMode(t *testing.T) {
	var cxt *ErrorContext

	cxt = NewErrorContext("panic test", codes)
	cxt.TurnPanicOn()

	defer func() {
		if err := recover(); err == nil {
			t.Errorf("panic was not triggered")
		}
	}()
	cxt.NewError(3, "invalid error code")
}

func TestErrorContextMapping(t *testing.T) {
	var err Error
	var cxt *ErrorContext

	cxt = NewErrorContext("mapping test", codes)

	mapErr := cxt.SetDefaultMessageForCode(0, "test error")
	if mapErr != nil {
		t.Errorf("mapping error code 0 incorrectly returned an error")
	}
	err = cxt.NewError(0)
	basicMatch(t, err, "test error")

	mapErr = cxt.SetDefaultMessageForCode(3, "invalid error code")
	if mapErr == nil {
		t.Errorf("mapping error code 3 should have created an error")
	}
}

// goes over internal errors
func TestErrorContextMisc(t *testing.T) {

	var err Error
	var cxt *ErrorContext

	cxt = NewErrorContext("bad input", codes)

	err = cxt.NewError(0)
	if err.Code() == 0 { // BadMsgLookup
		t.Errorf("default mapping for error code 0 should not exist yet")
	}

	err = cxt.NewError(0, fmt.Errorf("not a fmt string"))
	if err.Code() == 0 { // BadFmtString
		t.Errorf("non-string type should not have been interpreted as a string")
	}

	mapErr := cxt.SetDefaultMessageForCode(0, "give default")
	if mapErr == nil { // BadInitOrder
		t.Errorf("should have gotten error for attempting to modify error context after using initial configuration")
	}

	panicErr := cxt.TurnPanicOn()
	if panicErr == nil { // BadInitOrder
		t.Errorf("should not be able to turn panic on after attempting to create errors with error context")
	}

	cxt = NewErrorContext("dup mapping", codes)
	cxtMap := map[uint]string{
		0: "zero default",
		1: "one default",
	}
	cxt.AddDefaultErrorMessages(cxtMap)

	mapErr = cxt.SetDefaultMessageForCode(0, "second zero default")
	if mapErr == nil { // BadDupMessage
		t.Errorf("should have gotten error for duplicate mapping")
	}

	mapErr = cxt.SetDefaultMessageForCode(3, "bad err code")
	if mapErr == nil { // BadErrorCode
		t.Errorf("should have gotten error for bad error code in map function")
	}
}

func TestErrorContext(t *testing.T) {
	var err Error
	var cxt *ErrorContext

	cxt = NewErrorContext("basic test", codes)

	stats := cxt.GetErrorStats()
	if stats != nil {
		t.Errorf("no errors have been created, error stats should be nil")
	}

	err = cxt.NewError(0, "no args")
	basicMatch(t, err, "no args")

	err = cxt.NewError(1, "arg is %v", "awesome")
	basicMatch(t, err, "arg is awesome")

	err = cxt.NewError(3, "invalid error code")
	if err.Code() != int(BadErrorCode) {
		t.Errorf("sucessfully made an error with a code that the context shouldn't know")
	}

	cxt.NewError(0, "second no args")

	stats = cxt.GetErrorStats()
	if stats[0] != 2 || stats[1] != 1 || stats[2] != 0 {
		t.Errorf("one of the stats in %v didn't match [2, 1, 0]", stats)
	}
}
