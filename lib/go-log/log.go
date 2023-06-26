// Inspired by https://www.goinggo.net/2013/11/using-log-package-in-go.html
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

	"github.com/apache/trafficcontrol/lib/go-llog"
)

var (
	Debug        *log.Logger
	Info         *log.Logger
	Warning      *log.Logger
	Error        *log.Logger
	Event        *log.Logger
	Access       *log.Logger
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

const DebugPrefix = "DEBUG: "
const InfoPrefix = "INFO: "
const WarnPrefix = "WARNING: "
const ErrPrefix = "ERROR: "
const EventPrefix = ""

const DebugFlags = log.Lshortfile
const InfoFlags = log.Lshortfile
const WarnFlags = log.Lshortfile
const ErrFlags = log.Lshortfile
const EventFlags = 0

// Init initailizes the logs with the given io.WriteClosers. If `Init` was previously called, existing loggers are Closed. If you have loggers which are not Closers or which must not be Closed, wrap them with `log.NopCloser`.
func Init(eventW, errW, warnW, infoW, debugW io.WriteCloser) {
	initLogger(&Debug, &debugCloser, debugW, DebugPrefix, DebugFlags)
	initLogger(&Info, &infoCloser, infoW, InfoPrefix, InfoFlags)
	initLogger(&Warning, &warnCloser, warnW, WarnPrefix, WarnFlags)
	initLogger(&Error, &errCloser, errW, ErrPrefix, ErrFlags)
	initLogger(&Event, &eventCloser, eventW, EventPrefix, EventFlags)
}

func InitAccess(accessW io.WriteCloser) {
	initLogger(&Access, &accessCloser, accessW, EventPrefix, EventFlags)
}

// Logf should generally be avoided, use the built-in Init or InitCfg and Errorf, Warnln, etc functions instead.
// It logs to the given logger, in the same format as the standard log functions.
// This should only be used in rare and unusual circumstances when the standard loggers and functions can't.
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

func Errorf(format string, v ...interface{}) { Logf(Error, format, v...) }
func Errorln(v ...interface{})               { Logln(Error, v...) }
func Warnf(format string, v ...interface{})  { Logf(Warning, format, v...) }
func Warnln(v ...interface{})                { Logln(Warning, v...) }
func Infof(format string, v ...interface{})  { Logf(Info, format, v...) }
func Infoln(v ...interface{})                { Logln(Info, v...) }
func Debugf(format string, v ...interface{}) { Logf(Debug, format, v...) }
func Debugln(v ...interface{})               { Logln(Debug, v...) }

const eventFormat = "%.3f %s"

func eventTime(t time.Time) float64 {
	return float64(t.Unix()) + (float64(t.Nanosecond()) / 1e9)
}

func Accessln(v ...interface{}) {
	if Access != nil {
		Access.Println(v...)
	}
}

// event log entries (TM event.log, TR access.log, etc)
func Eventf(t time.Time, format string, v ...interface{}) {
	if Event == nil {
		return
	}
	// 1484001185.287 ...
	Event.Printf(eventFormat, eventTime(t), fmt.Sprintf(format, v...))
}

// EventfRaw writes to the event log with no prefix.
func EventfRaw(format string, v ...interface{}) {
	if Event == nil {
		return
	}
	Event.Printf(format, v...)
}

// EventRaw writes to the event log with no prefix, and no newline. Go's Printf is slow, using this with string concatenation is by far the fastest way to log, and should be used for frequent logs.
func EventRaw(s string) {
	if Event == nil {
		return
	}
	Event.Output(stackFrame, s)
}

// Close calls `Close()` on the given Closer, and logs any error. On error, the context is logged, followed by a colon, the error message, and a newline. This is primarily designed to be used in `defer`, for example, `defer log.Close(resp.Body, "readData fetching /foo/bar")`.
func Close(c io.Closer, context string) {
	err := c.Close()
	if err != nil {
		Errorf("%v: %v", context, err)
	}
}

// Closef acts like Close, with a given format string and values, followed by a colon, the error message, and a newline. The given values are not coerced, concatenated, or printed unless an error occurs, so this is more efficient than `Close()`.
func Closef(c io.Closer, contextFormat string, v ...interface{}) {
	err := c.Close()
	if err != nil {
		Errorf(contextFormat, v...)
		Errorf(": %v", err)
	}
}

// Write calls `Write()` on the given Writer, and logs any error. On error, the context is logged, followed by a colon, the error message, and a newline.
func Write(w io.Writer, b []byte, context string) {
	_, err := w.Write(b)
	if err != nil {
		Errorf("%v: %v", context, err)
	}
}

// Writef acts like Write, with a given format string and values, followed by a colon, the error message, and a newline. The given values are not coerced, concatenated, or printed unless an error occurs, so this is more efficient than `Write()`.
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

func (nopCloser) Close() error { return nil }

func NopCloser(w io.Writer) io.WriteCloser {
	return nopCloser{w}
}

// LogLocation is a location to log to. This may be stdout, stderr, null (/dev/null), or a valid file path.
type LogLocation string

const (
	// LogLocationStdout indicates the stdout IO stream
	LogLocationStdout = "stdout"
	// LogLocationStderr indicates the stderr IO stream
	LogLocationStderr = "stderr"
	// LogLocationNull indicates the null IO stream (/dev/null)
	LogLocationNull = "null"
	// LogLocationFile specify where health client logs should go.
	LogLocationFile = "/var/log/trafficcontrol/tc-health-client.log"
	// StaticFileDir is the directory that contains static html and js files.
	StaticFileDir = "/opt/traffic_monitor/static/"
)

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

type Config interface {
	ErrorLog() LogLocation
	WarningLog() LogLocation
	InfoLog() LogLocation
	DebugLog() LogLocation
	EventLog() LogLocation
}

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
