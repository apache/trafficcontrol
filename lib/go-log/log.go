// Package log contains utilities for sensible logging.
//
// ATC components written in Go should use this package rather than a
// third-party library. This package also provides log *levels*, which are not
// really supported by the standard log package, so it's generally better to
// use this than even that standard library.
//
// Inspired by https://www.goinggo.net/2013/11/using-log-package-in-go.html
//
// # Usage
//
// Normally, one will have some command-line program that takes in some kind of
// configuration either via command-line arguments or a configuration file, and
// want that configuration to define how logging happens. The easiest way to
// logging up and running is to put that configuration into a type that
// implements this package's Config interface, like:
//
//	type Configuration struct {
//		ErrorLogs     string `json:"errorLogs"`
//		WarningLogs   string `json:"warningLogs"`
//		InfoLogs      string `json:"infoLogs"`
//		DebugLogs     string `json:"debugLogs"`
//		EventLogs     string `json:"eventLogs"`
//		EnableFeature bool   `json:"enableFeature"`
//	}
//
//	func (c Configuration) ErrorLog() string { return c.ErrorLogs}
//	func (c Configuration) WarningLog() string { return c.WarningLogs }
//	func (c Configuration) InfoLog() string { return c.InfoLogs }
//	func (c Configuration) DebugLog() string { return c.DebugLogs }
//	func (c Configuration) EventLog() string { return c.EventLogs }
//
// Then the configuration can be unmarshaled from JSON and the resulting
// Configuration can be passed directly to InitCfg, and then everything is
// set-up and ready to go.
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
	"io"
	"log"
	"os"
	"time"

	"github.com/apache/trafficcontrol/v8/lib/go-llog"
)

// These are the loggers for the various log levels, which can be accessed
// directly, but in general the functions exported by this package should
// be used instead.
var (
	Debug   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	Event   *log.Logger
	// Access is NOT initialized by Init nor InitCfg in any way. To initialize
	// it, use InitAccess.
	Access *log.Logger
)

var (
	debugCloser  io.Closer
	infoCloser   io.Closer
	warnCloser   io.Closer
	errCloser    io.Closer
	eventCloser  io.Closer
	accessCloser io.Closer
)

func initLogger(logger **log.Logger, oldLogCloser *io.Closer, newLogWriter io.WriteCloser, logPrefix string, logFlags int) {
	if newLogWriter == nil {
		*logger = nil
		if *oldLogCloser != nil {
			(*oldLogCloser).Close()
			*oldLogCloser = nil
		}
		return
	}

	if *logger != nil {
		(*logger).SetOutput(newLogWriter)
	} else {
		*logger = log.New(newLogWriter, logPrefix, logFlags)
	}

	if *oldLogCloser != nil {
		(*oldLogCloser).Close()
	}
	*oldLogCloser = newLogWriter
}

// These are prefixes prepended to messages for the various log levels.
const (
	DebugPrefix = "DEBUG: "
	InfoPrefix  = "INFO: "
	WarnPrefix  = "WARNING: "
	ErrPrefix   = "ERROR: "
	EventPrefix = ""
)

// These constants are flags used to create the underlying "log.Logger",
// defining what's in the prefix for different types of log messages.
//
// Refer to the "log" package documentation for details.
const (
	DebugFlags = log.Lshortfile
	InfoFlags  = log.Lshortfile
	WarnFlags  = log.Lshortfile
	ErrFlags   = log.Lshortfile
	EventFlags = 0
)

// Init initializes the logs - except the Access log stream - with the given
// io.WriteClosers.
//
// If `Init` was previously called, existing loggers are Closed.
// If you have loggers which are not Closers or which must not be Closed, wrap
// them with `log.NopCloser`.
func Init(eventW, errW, warnW, infoW, debugW io.WriteCloser) {
	initLogger(&Debug, &debugCloser, debugW, DebugPrefix, DebugFlags)
	initLogger(&Info, &infoCloser, infoW, InfoPrefix, InfoFlags)
	initLogger(&Warning, &warnCloser, warnW, WarnPrefix, WarnFlags)
	initLogger(&Error, &errCloser, errW, ErrPrefix, ErrFlags)
	initLogger(&Event, &eventCloser, eventW, EventPrefix, EventFlags)
}

// InitAccess initializes the Access logging stream with the given
// io.WriteCloser.
func InitAccess(accessW io.WriteCloser) {
	initLogger(&Access, &accessCloser, accessW, EventPrefix, EventFlags)
}

// Logf should generally be avoided, use the Init or InitCfg and Errorf, Warnln,
// etc functions instead.
//
// It logs to the given logger, in the same format as the standard log
// functions. This should only be used in rare and unusual circumstances when
// the standard loggers and functions can't.
func Logf(logger *log.Logger, format string, v ...interface{}) {
	if logger == nil {
		return
	}
	logger.Output(stackFrame, time.Now().UTC().Format(timeFormat)+": "+fmt.Sprintf(format, v...))
}

// Logln should generally be avoided, use the built-in Init or InitCfg and Errorf, Warnln, etc functions instead.
// It logs to the given logger, in the same format as the standard log functions.
// This should only be used in rare and unusual circumstances when the standard loggers and functions can't.
func Logln(logger *log.Logger, v ...interface{}) {
	if logger == nil {
		return
	}
	logger.Output(stackFrame, time.Now().UTC().Format(timeFormat)+": "+fmt.Sprintln(v...))
}

const timeFormat = time.RFC3339Nano
const stackFrame = 3

// Errorf prints a formatted message to the Error logging stream.
//
// Formatting is handled identically to the fmt package, but the final message
// actually printed to the stream will be prefixed with the log level and
// timestamp just like Errorln. Also, if the final formatted message does not
// end in a newline character, one will be appended, so it's not necessary to
// do things like e.g.
//
//	Errorf("%v\n", 5)
//
// Instead, just:
//
//	Errorf("%v", 5)
//
// will suffice, and will print exactly the same message.
func Errorf(format string, v ...interface{}) { Logf(Error, format, v...) }

// Errorln prints a line to the Error stream.
//
// This will use default formatting for all of its passed parameters, as
// defined by the fmt package. The resulting line is first prefixed with the
// log level - in this case "ERROR: " - and the current timestamp, as well as
// the context in which it was logged. For example, if line 36 of foo.go
// invokes this like so:
//
//	log.Errorln("something wicked happened")
//
// ... the resulting log line will look like:
//
//	ERROR: foo.go:36: 2006-01-02T15:04:05Z: something wicked happened
//
// (assuming the log call occurred at the time package's formatting reference
// date/time). This allows multiple log streams to be directed to the same file
// for output and still be distinguishable by their log-level prefix.
func Errorln(v ...interface{}) { Logln(Error, v...) }

// Warnf prints a formatted message to the Warning logging stream.
//
// The message format is identical to Errorf, but using the log-level prefix
// "WARNING: ".
func Warnf(format string, v ...interface{}) { Logf(Warning, format, v...) }

// Warnln prints a line to the Warning stream.
//
// The message format is identical to Errorln, but using the log-level prefix
// "WARNING: ".
func Warnln(v ...interface{}) { Logln(Warning, v...) }

// Infof prints a formatted message to the Info logging stream.
//
// The message format is identical to Errorf, but using the log-level prefix
// "INFO: ".
func Infof(format string, v ...interface{}) { Logf(Info, format, v...) }

// Infoln prints a line to the Info stream.
//
// The message format is identical to Errorln, but using the log-level prefix
// "INFO: ".
func Infoln(v ...interface{}) { Logln(Info, v...) }

// Debugf prints a formatted message to the Debug logging stream.
//
// The message format is identical to Errorf, but using the log-level prefix
// "DEBUG: ".
func Debugf(format string, v ...interface{}) { Logf(Debug, format, v...) }

// Debugln prints a line to the Debug stream.
//
// The message format is identical to Errorln, but using the log-level prefix
// "DEBUG: ".
func Debugln(v ...interface{}) { Logln(Debug, v...) }

const eventFormat = "%.3f %s"

func eventTime(t time.Time) float64 {
	return float64(t.Unix()) + (float64(t.Nanosecond()) / float64(time.Second))
}

// Accessln prints an "access"-level log to the Access logging stream.
//
// This does NOT use the same message formatting as the other logging functions
// like Errorf. Instead, it prints only the information it's given, formatted
// according to the format string (if present) as defined by the fmt package.
// For example, if line 36 of foo.go invokes this like so:
//
//	log.Accessln("%s\n", "my message")
//
// ... the resulting log line will look like:
//
//	my message
func Accessln(v ...interface{}) {
	if Access != nil {
		Access.Println(v...)
	}
}

// Eventf prints an "event"-level log to the Event logging stream.
//
// This does NOT use the same message formatting as the other logging functions
// like Errorf. Instead, it prints the time at which the event occurred as the
// number of seconds since the Unix epoch to three decimal places of sub-second
// precision, followed by the trailing parameters formatted according to the
// format string as defined by the fmt package. For example, if line 36 of
// foo.go invokes this like so:
//
//	log.Eventf(time.Now(), "%s\n", "my message")
//
// ... the resulting log line will look like:
//
//	1136214245.000 my message
//
// Note that this WILL NOT add trailing newlines if the resulting formatted
// message string doesn't end with one.
func Eventf(t time.Time, format string, v ...interface{}) {
	if Event == nil {
		return
	}
	// 1484001185.287 ...
	Event.Printf(eventFormat, eventTime(t), fmt.Sprintf(format, v...))
}

// EventfRaw writes a formatted message to the Event log stream.
//
// The formatting is just like Eventf, but no timestamp prefix will be added.
func EventfRaw(format string, v ...interface{}) {
	if Event == nil {
		return
	}
	Event.Printf(format, v...)
}

// EventRaw writes the given string directly to the Event log stream with
// absolutely NO formatting - trailing newline, timestamp, log-level etc.
//
// Go's fmt.Printf/Sprintf etc. are very slow, so using this with string
// concatenation is by far the fastest way to log, and should be used for
// frequent logs.
func EventRaw(s string) {
	if Event == nil {
		return
	}
	Event.Output(stackFrame, s)
}

// Close calls `Close()` on the given Closer, and logs any error that results
// from that call.
//
// On error, the context is logged, followed by a colon and the error message.
// Specifically, it calls Errorf with the context and error using the format
// '%v: %v'.
//
// This is primarily designed to be used in `defer`, for example:
//
//	defer log.Close(resp.Body, "readData fetching /foo/bar")
//
// Note that some Go linting tools may not properly detect this as both Closing
// the Closer and checking the error value that Close returns, but it does both
// of those things, we swear.
func Close(c io.Closer, context string) {
	err := c.Close()
	if err != nil {
		Errorf("%v: %v", context, err)
	}
}

// Closef acts like Close, with a given format string and values, followed by a
// colon, the error message, and a newline.
//
// The given values are not coerced, concatenated, or printed unless an error
// occurs, so this is more efficient than Close in situations where the message
// passed to Close would have been generated with fmt.Sprintf.
//
// This actually results in two distinct lines being printed in the Error log
// stream. For example, if line 36 of foo.go invokes this like so:
//
//	log.Closef(resp.Body, "doing %s", "something")
//
// ... and resp.Body.Close() returns a non-nil error with the message "failed",
// the resulting log lines will look like:
//
//	ERROR: log.go:271: 2006-01-02T15:04:05Z: doing something
//	ERROR: log.go:272: 2006-01-02T15:04:05Z: : failed
//
// (please please please don't count on those line numbers).
func Closef(c io.Closer, contextFormat string, v ...interface{}) {
	err := c.Close()
	if err != nil {
		Errorf(contextFormat, v...)
		Errorf(": %v", err)
	}
}

// Write calls the Write method of the given Writer with the passed byte slice
// as an argument, and logs any error that arises from that call.
//
// On error, the context is logged, followed by a colon and the error message, and
// a trailing newline is guaranteed. For example, if line 36 of foo.go invokes
// this like so:
//
//	log.Write(someWriter, []byte("write me"), "trying to write")
//
// ... and someWriter.Write() returns a non-nil error with the message
// "failed", the resulting log line will look like:
//
//	ERROR: log.go:294: 2006-01-02T15:04:05Z: trying to write: failed
//
// (please please please don't count on that line number).
func Write(w io.Writer, b []byte, context string) {
	_, err := w.Write(b)
	if err != nil {
		Errorf("%v: %v", context, err)
	}
}

// Writef acts like Write, with a given format string and values, followed by a
// colon, the error message, and a newline.
//
// The given values are not coerced, concatenated, or printed unless an error
// occurs, so this is more efficient than Write in situations where the message
// passed to Write would have been generated with fmt.Sprintf.
//
// This actually results in two distinct lines being printed in the Error log
// stream. For example, if line 36 of foo.go invokes this like so:
//
//	log.Writef(someWriter, "doing %s", "something")
//
// ... and someWriter.Write() returns a non-nil error with the message
// "failed", the resulting log lines will look like:
//
//	ERROR: log.go:320: 2006-01-02T15:04:05Z: doing something
//	ERROR: log.go:321: 2006-01-02T15:04:05Z: : failed
//
// (please please please don't count on those line numbers).
func Writef(w io.Writer, b []byte, contextFormat string, v ...interface{}) {
	_, err := w.Write(b)
	if err != nil {
		Errorf(contextFormat, v...)
		Errorf(": %v", err)
	}
}

type nopCloser struct {
	io.Writer
}

// Close implements the io.Closer interface. This does nothing.
func (nopCloser) Close() error { return nil }

// NopCloser returns an io.WriteCloser which wraps the passed io.Writer with a
// Close method that just does nothing but return a nil error. This allows
// non-closable streams to be used for logging.
func NopCloser(w io.Writer) io.WriteCloser {
	return nopCloser{w}
}

// LogLocation is a location to log to. This may be stdout, stderr, null (/dev/null), or a valid file path.
type LogLocation string

const (
	// LogLocationStdout indicates the stdout IO stream.
	LogLocationStdout = "stdout"
	// LogLocationStderr indicates the stderr IO stream.
	LogLocationStderr = "stderr"
	// LogLocationNull indicates the null IO stream (/dev/null).
	LogLocationNull = "null"
)

// LogLocationFile specify where health client logs should go.
//
// Deprecated: It is unclear why this constant is in this package at all, and
// it will almost certainly be removed in the future.
const LogLocationFile = "/var/log/trafficcontrol/tc-health-client.log"

// StaticFileDir is the directory that contains static HTML and Javascript
// files for Traffic Monitor.
//
// Deprecated: It is unclear why this constant is in this package at all, and
// it will almost certainly be removed in the future.
const StaticFileDir = "/opt/traffic_monitor/static/"

// GetLogWriter creates a writable stream from the given LogLocation.
//
// If the requested log location is a file, it will be opened with the
// write-only, create, and append flags in permissions mode 0644. If that open
// operation causes an error, it is returned.
//
// As a special case, if the location is an empty string, the returned stream
// will be nil - this is the same behavior as passing LogLocationNull.
func GetLogWriter(location LogLocation) (io.WriteCloser, error) {
	switch location {
	case LogLocationStdout:
		return NopCloser(os.Stdout), nil
	case LogLocationStderr:
		return NopCloser(os.Stderr), nil
	case LogLocationNull:
		fallthrough
	case "":
		return nil, nil
	default:
		return os.OpenFile(string(location), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	}
}

// Config structures represent logging configurations, which can be accessed by
// this package for setting up the logging system via GetLogWriters.
type Config interface {
	ErrorLog() LogLocation
	WarningLog() LogLocation
	InfoLog() LogLocation
	DebugLog() LogLocation
	EventLog() LogLocation
}

// GetLogWriters uses the passed Config to set up and return log streams.
//
// This bails at the first encountered error creating a log stream, and returns
// the error without trying to create any further streams.
//
// The returned streams still need to be initialized, usually with Init. To
// create and initialize the logging streams at the same time, refer to
// InitCfg.
func GetLogWriters(cfg Config) (io.WriteCloser, io.WriteCloser, io.WriteCloser, io.WriteCloser, io.WriteCloser, error) {
	eventLoc := cfg.EventLog()
	errLoc := cfg.ErrorLog()
	warnLoc := cfg.WarningLog()
	infoLoc := cfg.InfoLog()
	debugLoc := cfg.DebugLog()

	eventW, err := GetLogWriter(eventLoc)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("getting log event writer %v: %v", eventLoc, err)
	}
	errW, err := GetLogWriter(errLoc)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("getting log error writer %v: %v", errLoc, err)
	}
	warnW, err := GetLogWriter(warnLoc)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("getting log warning writer %v: %v", warnLoc, err)
	}
	infoW, err := GetLogWriter(infoLoc)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("getting log info writer %v: %v", infoLoc, err)
	}
	debugW, err := GetLogWriter(debugLoc)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("getting log debug writer %v: %v", debugLoc, err)
	}
	return eventW, errW, warnW, infoW, debugW, nil
}

// InitCfg uses the passed configuration to both create and initialize all
// logging streams.
func InitCfg(cfg Config) error {
	eventW, errW, warnW, infoW, debugW, err := GetLogWriters(cfg)
	if err != nil {
		return err
	}
	Init(eventW, errW, warnW, infoW, debugW)
	return nil
}

// LLog returns an llog.Log, for passing to libraries using llog.
//
// Note the returned Log will have writers tied to loggers at its time of creation.
// Thus, it's safe to reuse LLog if an application never re-initializes the loggers,
// such as when reloading config on a HUP signal.
// If the application re-initializes loggers, LLog should be called again to get
// a new llog.Log associated with the new log locations.
func LLog() llog.Log {
	// ltow converts a log.Logger into an io.Writer
	// This is relatively inefficient. If performance is necessary, this package could be made to
	// keep track of the original io.Writer, to avoid an extra function call and string copy.
	ltow := func(lg *log.Logger) io.Writer {
		return llog.WriterFunc(func(p []byte) (n int, err error) {
			Logln(lg, string(p))
			return len(p), nil
		})
	}
	return llog.New(ltow(Error), ltow(Warning), ltow(Info), ltow(Debug))
}
