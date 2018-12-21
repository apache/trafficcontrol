package test

import (
	"fmt"
	"testing"
)

var codes []uint

func init() {
	codes = []uint{
		0, 1, 2,
	}
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
	if err.ErrorCode() != expected {
		t.Errorf("\nexpected: %d\nactual: %d\n\n", expected, err.ErrorCode())
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
	var err ErrorCoder
	var cxt *ErrorContext

	cxt = NewErrorContext("mapping test", codes)

	mapErr := cxt.AddMapping(0, "test error")
	if mapErr != nil {
		t.Errorf("mapping error code 0 incorrectly returned an error")
	}
	err = cxt.NewError(0)
	basicMatch(t, err, "test error")

	mapErr = cxt.AddMapping(3, "invalid error code")
	if mapErr == nil {
		t.Errorf("mapping error code 3 should have created an error")
	}
}

// goes over internal errors
func TestErrorContextMisc(t *testing.T) {

	var err ErrorCoder
	var cxt *ErrorContext

	cxt = NewErrorContext("bad input", codes)

	err = cxt.NewError(0)
	if err.ErrorCode() == 0 { // BAD_MAP_LOOKUP
		t.Errorf("default mapping for error code 0 should not exist yet")
	}

	err = cxt.NewError(0, fmt.Errorf("not a fmt string"))
	if err.ErrorCode() == 0 { // BAD_FMT_STRING
		t.Errorf("non-string type should not have been interpreted as a string")
	}

	mapErr := cxt.AddMapping(0, "give default")
	if mapErr == nil { // BAD_INIT_TIMING
		t.Errorf("should have gotten error for attempting to modify error context after using initial configuration")
	}

	panicErr := cxt.TurnPanicOn()
	if panicErr == nil { // BAD_INIT_TIMING
		t.Errorf("should not be able to turn panic on after attempting to create errors with error context")
	}

	cxt = NewErrorContext("dup mapping", codes)
	cxtMap := map[uint]string{
		0: "zero default",
		1: "one default",
	}
	cxt.AddMap(cxtMap)

	mapErr = cxt.AddMapping(0, "second zero default")
	if mapErr == nil { // BAD_DUP_MAPPING
		t.Errorf("should have gotten error for duplicate mapping")
	}

	mapErr = cxt.AddMapping(3, "bad err code")
	if mapErr == nil { // BAD_MAP_CREATE
		t.Errorf("should have gotten error for bad error code in map function")
	}
}

func TestErrorContext(t *testing.T) {
	var err ErrorCoder
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
	if err.ErrorCode() != BAD_ERROR_CODE {
		t.Errorf("sucessfully made an error with a code that the context shouldn't know")
	}

	cxt.NewError(0, "second no args")

	stats = cxt.GetErrorStats()
	if stats[0] != 2 || stats[1] != 1 || stats[2] != 0 {
		t.Errorf("one of the stats in %v didn't match [2, 1, 0]", stats)
	}
}
