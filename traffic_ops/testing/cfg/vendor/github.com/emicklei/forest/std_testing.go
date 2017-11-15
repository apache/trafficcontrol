package forest

import (
	"fmt"
	"strings"
)

// T is the interface that this package is using from standard testing.T
type T interface {
	// Logf formats its arguments according to the format, analogous to Printf, and records the text in the error log.
	// The text will be printed only if the test fails or the -test.v flag is set.
	Logf(format string, args ...interface{})
	// Error is equivalent to Log followed by Fail.
	Error(args ...interface{})
	// Fatal is equivalent to Log followed by FailNow.
	Fatal(args ...interface{})
	// FailNow marks the function as having failed and stops its execution.
	FailNow()
	// Fail marks the function as having failed but continues execution.
	Fail()
}

// TestingT provides a sub-api of testing.T. Its purpose is to allow the use of this package in TestMain(m).
var TestingT = Logger{InfoEnabled: true, ErrorEnabled: true, ExitOnFatal: true}

// LoggingPrintf is the function used by TestingT to produce logging on Logf,Error and Fatal.
var LoggingPrintf = fmt.Printf

func tabify(format string) string {
	if strings.HasPrefix(format, "\n") {
		return strings.Replace(format, "\n", "\n\t\t", 1)
	}
	return format
}
