package log

/*
* Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
*/

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"
)

type mockTester struct {
	logbuf       string
	errbuf       string
	failed       bool
	helperCalled bool
	*testing.T
}

/*** testing.TB method overrides ***/
func (m *mockTester) Helper() {
	m.helperCalled = true
}
func (m *mockTester) Failed() bool {
	return m.failed
}
func (m *mockTester) Error(args ...any) {
	m.failed = true
	m.errbuf = fmt.Sprint(args...)
}
func (m *mockTester) Log(args ...any) {
	m.logbuf = fmt.Sprint(args...)
}

/*** None of these should ever be called ***/
func (m *mockTester) Fail() {
	m.failed = true
	m.T.Error("'Fail' called")
}
func (m *mockTester) FailNow() {
	m.failed = true
	m.T.Error("'FailNow' called")
}
func (m *mockTester) Fatal(args ...any) {
	m.failed = true
	m.T.Error("'Fatal' called")
}
func (m *mockTester) Errorf(fmtstr string, args ...any) {
	m.failed = true
	m.T.Error("'Errorf' called")
}
func (m *mockTester) Fatalf(fmtstr string, args ...any) {
	m.failed = true
	m.T.Error("'Fatalf' called")
}
func (m *mockTester) Logf(fmtstr string, args ...any) {
	m.failed = true
	m.T.Error("'Logf' called")
}
func (m *mockTester) Name() string {
	m.failed = true
	m.T.Error("'Name' called")
	return ""
}
func (m *mockTester) Skip(args ...any) {
	m.failed = true
	m.T.Error("'Skip' called")
}
func (m *mockTester) Skipf(fmtstr string, args ...any) {
	m.failed = true
	m.T.Error("'Skipf' called")
}
func (m *mockTester) SkipNow() {
	m.failed = true
	m.T.Error("'SkipNow' called")
}
func (m *mockTester) Skipped() bool {
	m.failed = true
	m.T.Error("'Skipped' called")
	return false
}

const datePattern = `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?Z`

func mkPattern(prefix, msg string) string {
	msg = strings.ReplaceAll(msg, `\`, `\\`)
	msg = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(msg, "[", `\[`), "]", `\]`), "(", `\(`), ")", `\)`)
	msg = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(msg, "{", `\{`), "}", `\}`), ".", `\.`), "^", `\^`)
	msg = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(msg, "$", `\$`), "+", `\+`), "*", `\*`), "?", `\?`)

	return "^" + prefix + "test\\.logger_test\\.go:\\d+: " + datePattern + ": " + msg + "\n$"
}

func matches(pattern, str string, t *testing.T) bool {
	match, err := regexp.Match(pattern, []byte(str))
	if err != nil {
		t.Errorf("failed to compile pattern '%s': %v", pattern, err)
		return false
	}
	return match
}

func TestInitTestLogging(t *testing.T) {
	m := mockTester{T: t}
	InitTestingLogging(&m, false)
	msg := "testing message"
	Debugln(msg)
	pattern := mkPattern(DebugPrefix, msg)
	if !matches(pattern, m.logbuf, t) {
		t.Errorf("expected last log message to match '%s', got: '%s'", pattern, m.logbuf)
	}
	Debugf("formatted: %d, %+v, %t", 1, []byte{1}, true)
	msg = "formatted: 1, [1], true"
	pattern = mkPattern(DebugPrefix, msg)
	if !matches(pattern, m.logbuf, t) {
		t.Errorf("expected last log message to match '%s', got: '%s'", pattern, m.logbuf)
	}

	msg = "testing message"
	Infoln(msg)
	pattern = mkPattern(InfoPrefix, msg)
	if !matches(pattern, m.logbuf, t) {
		t.Errorf("expected last log message to match '%s', got: '%s'", pattern, m.logbuf)
	}
	Infof("formatted: %d, %+v, %t", 1, []byte{1}, true)
	msg = "formatted: 1, [1], true"
	pattern = mkPattern(InfoPrefix, msg)
	if !matches(pattern, m.logbuf, t) {
		t.Errorf("expected last log message to match '%s', got: '%s'", pattern, m.logbuf)
	}

	msg = "testing message"
	Warnln(msg)
	pattern = mkPattern(WarnPrefix, msg)
	if !matches(pattern, m.logbuf, t) {
		t.Errorf("expected last log message to match '%s', got: '%s'", pattern, m.logbuf)
	}
	Warnf("formatted: %d, %+v, %t", 1, []byte{1}, true)
	msg = "formatted: 1, [1], true"
	pattern = mkPattern(WarnPrefix, msg)
	if !matches(pattern, m.logbuf, t) {
		t.Errorf("expected last log message to match '%s', got: '%s'", pattern, m.logbuf)
	}

	msg = "testing message"
	EventRaw(msg)
	if m.logbuf != msg+"\n" {
		t.Errorf("expected last log message to be exactly '%s', got: '%s'", msg, m.logbuf)
	}
	Eventf(time.Time{}, "formatted: %d, %+v, %t", 1, []byte{1}, true)
	msg = "formatted: 1, [1], true"
	pattern = fmt.Sprintf(eventFormat, eventTime(time.Time{}), "formatted: 1, [1], true\n")
	if pattern != m.logbuf {
		t.Errorf("expected last log message to be exactly '%s', got: '%s'", pattern, m.logbuf)
	}

	if m.Failed() {
		t.Error("none of the non-error logging functions should have caused a failure")
	}
	if m.errbuf != "" {
		t.Error("none of the non-error logging functions should have populated the error buffer:", m.errbuf)
	}

	m.logbuf = ""
	msg = "testing message"
	Errorln(msg)
	if m.logbuf != "" {
		t.Error("error should've been logged to the error stream, not the logging stream")
		m.logbuf = ""
	}
	if !m.Failed() {
		t.Error("error logging should have caused a failure")
	} else {
		m.failed = false
	}
	pattern = mkPattern(ErrPrefix, msg)
	if !matches(pattern, m.errbuf, t) {
		t.Errorf("expected last error message to match '%s', got: '%s'", pattern, m.errbuf)
	}
	Errorf("formatted: %d, %+v, %t", 1, []byte{1}, true)
	if m.logbuf != "" {
		t.Error("error should've been logged to the error stream, not the logging stream")
	}
	if !m.Failed() {
		t.Error("error logging should have caused a failure")
	}
	msg = "formatted: 1, [1], true"
	pattern = mkPattern(ErrPrefix, msg)
	if !matches(pattern, m.errbuf, t) {
		t.Errorf("expected last error message to match '%s', got: '%s'", pattern, m.errbuf)
	}
}

func TestLoggingWarningsAsErrors(t *testing.T) {
	m := mockTester{T: t}
	InitTestingLogging(&m, true)
	msg := "testing message"
	Warnln(msg)
	if m.logbuf != "" {
		t.Error("warning should've been logged to the error stream, not the logging stream")
		m.logbuf = ""
	}
	if !m.Failed() {
		t.Error("warning logging should have caused a failure")
	} else {
		m.failed = false
	}
	pattern := mkPattern(WarnPrefix, msg)
	if !matches(pattern, m.errbuf, t) {
		t.Errorf("expected last error message to match '%s', got: '%s'", pattern, m.errbuf)
	}
	Warnf("formatted: %d, %+v, %t", 1, []byte{1}, true)
	if m.logbuf != "" {
		t.Error("warning should've been logged to the error stream, not the logging stream")
	}
	if !m.Failed() {
		t.Error("warning logging should have caused a failure")
	}
	msg = "formatted: 1, [1], true"
	pattern = mkPattern(WarnPrefix, msg)
	if !matches(pattern, m.errbuf, t) {
		t.Errorf("expected last error message to match '%s', got: '%s'", pattern, m.errbuf)
	}
}
