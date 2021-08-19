// Package nagios is used by Traffic Monitor for unknown reasons.
//
// This may be moved internal to Traffic Monitor at some point, so its use is
// discouraged.
package nagios

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
	"os"
	"strings"
)

// A Status is an exit code that can be passed to Exit.
type Status int

// These are the different values allowed for a Status.
//
// Note that many things consider any non-zero exit code to be indicative of an
// error causing the program to quit, despite any names seen here.
const (
	Ok       Status = 0
	Warning  Status = 1
	Critical Status = 2
)

// Exit causes the current running program to exit by calling os.Exit with the
// given Status as an exit code.
//
// If the passed msg is not an empty string, it will be logged to stdout (NOT
// stderr, and directly using 'fmt.Printf', which may bypass logging
// configurations). Trailing newlines are ensured and not duplicated.
func Exit(status Status, msg string) {
	if msg != "" {
		msg = strings.TrimRight(msg, "\n")
		fmt.Printf("%s\n", msg)
	}
	os.Exit(int(status))
}
