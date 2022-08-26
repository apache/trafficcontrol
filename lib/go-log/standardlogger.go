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
	"log"
	"strings"
)

// StandardLogger returns a new Logger which will write the appropriate prefix for standard log Printf calls.
// This allows the given logger to prefix correctly when passed to third-party or standard library functions which only know about the standard Logger interface.
// It does this by wrapping logger's error print as a writer, and sets that as the new Logger's io.Writer.
//
// prefix is a prefix to add, which will be added immediately before messages, but after any existing prefix on logger and the timestamp.
func StandardLogger(logger *log.Logger, prefix string) *log.Logger {
	return log.New(&standardLoggerWriter{realLogger: logger, prefix: prefix}, "", 0)
}

type standardLoggerWriter struct {
	realLogger *log.Logger
	prefix     string
}

// Write writes to writer's underlying log, in the standard log format (note this is not the Event log format).
// The writer.realLogger may be nil, in which case this does nothing.
// This always returns len(p) and nil, claiming it successfully wrote even if it didn't.
func (writer *standardLoggerWriter) Write(p []byte) (n int, err error) {
	Logln(writer.realLogger, writer.prefix+strings.TrimSpace(string(p)))
	return len(p), nil
}
