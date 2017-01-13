/*
   Copyright 2015 Comcast Cable Communications Management, LLC

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package testHelper

import (
	"fmt"
	"testing"

	"github.com/cihub/seelog"
)

// Succeed is the Unicode codepoint for a check mark.
const Succeed = "\u2713"

// Failed is the Unicode codepoint for an X mark.
const Failed = "\u2717"

func init() {
	// An example seelog.xml is located at conf/seelog.xml
	// If you want to use this config - copy it to your project
	seelogConfig := "conf/seelog.xml.test"
	logger, err := seelog.LoggerFromConfigAsFile(seelogConfig)
	if err != nil {
		err := fmt.Errorf("FOO Error creating Logger from seelog file: %s", seelogConfig)
		seelog.Error(err)
	}
	defer seelog.Flush()
	seelog.ReplaceLogger(logger)
}

// Context is a summary of the test being run.
func Context(t *testing.T, msg string, args ...interface{}) {
	t.Log(fmt.Sprintf(msg, args...))
}

// Fatal contains details of a failed test and stops its execution.
func Fatal(t *testing.T, msg string, args ...interface{}) {
	m := fmt.Sprintf(msg, args...)
	t.Fatal(fmt.Sprintf("\t %-80s", m), Failed)
}

// Error contails details of a failed test and continues execution.
func Error(t *testing.T, msg string, args ...interface{}) {
	m := fmt.Sprintf(msg, args...)
	t.Error(fmt.Sprintf("\t %-80s", m), Failed)
}

// Success contains details of a successful test.
func Success(t *testing.T, msg string, args ...interface{}) {
	m := fmt.Sprintf(msg, args...)
	t.Log(fmt.Sprintf("\t %-80s", m), Succeed)
}
