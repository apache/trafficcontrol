package riak

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestLog(t *testing.T) {
	EnableDebugLogging = true

	tests := []struct {
		setLoggerFunc func(*log.Logger)
		logFunc       func(string, string, ...interface{})
		prefix        string
	}{
		{
			SetErrorLogger,
			logError,
			"[ERROR]",
		},
		{
			SetLogger,
			logWarn,
			"[WARNING]",
		},
		{
			SetLogger,
			logDebug,
			"[DEBUG]",
		},
	}

	for _, tt := range tests {
		buf := &bytes.Buffer{}
		logger := log.New(buf, "", log.LstdFlags)
		tt.setLoggerFunc(logger)
		tt.logFunc("[test]", "Hello %s!", "World")

		actual := buf.String()
		suffix := fmt.Sprintf("%s %s", tt.prefix, "[test] Hello World!\n")

		if !strings.HasSuffix(actual, suffix) {
			t.Errorf("Expected %s to end with %s", actual, suffix)
		}
	}
}

func TestLogln(t *testing.T) {
	EnableDebugLogging = true

	tests := []struct {
		setLoggerFunc func(*log.Logger)
		logFunc       func(string, ...interface{})
		prefix        string
	}{
		{
			SetErrorLogger,
			logErrorln,
			"[ERROR]",
		},
		{
			SetLogger,
			logWarnln,
			"[WARNING]",
		},
		{
			SetLogger,
			logDebugln,
			"[DEBUG]",
		},
	}

	for _, tt := range tests {
		buf := &bytes.Buffer{}
		logger := log.New(buf, "", log.LstdFlags)
		tt.setLoggerFunc(logger)
		tt.logFunc("[test]", "Hello", "World!")

		actual := buf.String()
		suffix := fmt.Sprintf("%s %s", tt.prefix, "[test] [Hello World!]\n")

		if !strings.HasSuffix(actual, suffix) {
			t.Errorf("Expected %s to end with %s", actual, suffix)
		}
	}
}

func TestDebugDisabled(t *testing.T) {
	EnableDebugLogging = false

	buf := &bytes.Buffer{}
	logger := log.New(buf, "", log.LstdFlags)
	SetLogger(logger)

	logDebug("[test]", "Hello %s!", "World")
	logDebugln("[test]", "Hello", "World")

	actual := buf.String()

	if len(actual) != 0 {
		t.Errorf("Debug was disabled but got %s", actual)
	}
}
