// Package llog provides logging utilities for library packages.
//
// This allows libraries to log if desired, while still allowing
// the library functions to have no side effects and not use
// global loggers, which may be different than the application's
// primary log library.
//
// This also allows users of the library to decide their logging
// level for this particular library. For example, an application
// may wish to generally log at the debug level, but not log
// debug messages for some particular library.
//
// Or, a user may wish to log warning messages from a library
// as errors. Setting two log levels to the same io.Writer
// is permissible.
//
// This is not itself a logging library. Rather, it allows
// applications using any log library to interface with libraries
// using llog, by constructing a Loggers object from their own
// logger.
package llog

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
)

// Log is an interface which library functions may accept
// in order to log without side effects, if the caller desires.
//
// Applications using a library which uses Log may use
// NewLog, or may themselves implement the interface.
//
// Library functions should immediately call DefaultIfNil to allow callers to
// pass a nil Log. Importantly, this allows applications using the library
// to avoid importing llog, if they don't want the library to log.
type Log interface {
	Errorf(format string, v ...interface{})
	Errorln(v ...interface{})
	Warnf(format string, v ...interface{})
	Warnln(v ...interface{})
	Infof(format string, v ...interface{})
	Infoln(v ...interface{})
	Debugf(format string, v ...interface{})
	Debugln(v ...interface{})
}

// LibInit initializes the Log for libraries.
//
// All public functions in Libraries using llog should immediately call LibInit.
// The return value should be assigned, like `log = llog.LibInit(log)`
//
// Applications creating a Log to pass to a library func should never call LibInit.
//
// This creates a Nop Log if the passed lg is nil, which allows applications
// to pass a nil Log.
//
// It may do other things in the future.
func LibInit(lg Log) Log {
	if lg == nil {
		return Nop()
	}
	return lg
}

// Nop returns a Log that never logs.
// Applications which don't want a library to log can pass Nop
// to libraries that take a Log.
func Nop() Log { return &loggers{} }

// New creates a new Log.
//
// Applications can use New to create a Log from their own internal log
// libraries and writers, to pass to libraries using liblog.
//
// Standard log example:
//
//	errLog := log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
//	mylib.MyFunc(liblog.New(errLog, nil, nil, nil), myArg)
//
// github.com/apache/trafficcontrol/v8/lib/go-log example:
//
//	import("github.com/apache/trafficcontrol/v8/lib/go-log")
//
//	log.Init(nil, os.Stderr, os.Stderr, nil, nil)
//
//	lLog := log.LLog() // lib/go-tc has a built-in llog helper
//
//	// alternatively, what the lib/go-log helper is doing internally:
//	ltow := func(lg *log.Logger) io.Writer {
//	  return llog.WriterFunc(func(p []byte) (n int, err error) {
//	  Logln(lg, string(p))
//	})
//	lLog = llog.New(ltow(Error), ltow(Warning), ltow(Info), ltow(Debug))
//
//	mylib.MyFunc(lLog, myArg)
//
// zap example:
//
//	import("go.uber.org/zap")
//
//	func main() {
//	  logger, _ := zap.NewProduction()
//	  sugar := logger.Sugar()
//	  zapErrLog := liblog.WriterFunc(func(p []byte) (n int, err error){
//	    logger.Sugar().Error(string(p))
//	  })
//	  mylib.MyFunc(liblog.New(zapErrLog, nil, nil, nil), myArg)
//	}
func New(err io.Writer, warn io.Writer, info io.Writer, debug io.Writer) Log {
	return &loggers{
		err:  err,
		warn: warn,
		inf:  info,
		dbg:  debug,
	}
}

// WriterFunc is an adapter to allow the use of ordinary functions as io.Writers.
// This behaves similar to http.HandlerFunc for http.Handler.
type WriterFunc func(p []byte) (n int, err error)

// Write implements io.Writer.
func (wf WriterFunc) Write(p []byte) (n int, err error) { return wf(p) }

type loggers struct {
	err  io.Writer
	warn io.Writer
	inf  io.Writer
	dbg  io.Writer
}

func (ls *loggers) Errorf(format string, v ...interface{}) { logf(ls.err, format, v...) }
func (ls *loggers) Errorln(v ...interface{})               { logln(ls.err, v...) }
func (ls *loggers) Warnf(format string, v ...interface{})  { logf(ls.warn, format, v...) }
func (ls *loggers) Warnln(v ...interface{})                { logln(ls.warn, v...) }
func (ls *loggers) Infof(format string, v ...interface{})  { logf(ls.inf, format, v...) }
func (ls *loggers) Infoln(v ...interface{})                { logln(ls.inf, v...) }
func (ls *loggers) Debugf(format string, v ...interface{}) { logf(ls.dbg, format, v...) }
func (ls *loggers) Debugln(v ...interface{})               { logln(ls.dbg, v...) }

func logf(wr io.Writer, format string, v ...interface{}) {
	if wr == nil {
		return
	}
	fmt.Fprintf(wr, format, v...)
}

func logln(wr io.Writer, v ...interface{}) {
	if wr == nil {
		return
	}
	fmt.Fprintln(wr, v...)
}
