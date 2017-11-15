package forest

import (
	"os"
)

// Logger can be used for the testing.T parameter for forest functions
// when you need more control over what to log and how to handle fatals.
// The variable TestingT is a Logger with all enabled.
type Logger struct {
	InfoEnabled  bool
	ErrorEnabled bool
	ExitOnFatal  bool
}

// Logf formats its arguments according to the format, analogous to Printf, and records the text in the error log.
// The text will be printed only if the test fails or the -test.v flag is set.
func (l Logger) Logf(format string, args ...interface{}) {
	if l.InfoEnabled {
		LoggingPrintf("\tinfo : "+tabify(format)+"\n", args...)
	}
}

// Error is equivalent to Log followed by Fail.
func (l Logger) Error(args ...interface{}) {
	if l.ErrorEnabled {
		LoggingPrintf("\terror: "+tabify("%s")+"\n", args)
	}
}

// Fatal is equivalent to Log followed by FailNow.
func (l Logger) Fatal(args ...interface{}) {
	if l.ErrorEnabled {
		LoggingPrintf("\tfatal: "+tabify("%s")+"\n", args...)
	}
	if l.ExitOnFatal {
		os.Exit(1)
	}
}

// FailNow marks the function as having failed and stops its execution.
func (l Logger) FailNow() {
	if l.ExitOnFatal {
		os.Exit(1)
	}
}

// Fail marks the function as having failed but continues execution.
func (l Logger) Fail() {}
