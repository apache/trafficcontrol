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

import "testing"

type testErrorLogger struct {
	t testing.TB
}

func (tl testErrorLogger) Write(data []byte) (int, error) {
	tl.t.Error(string(data))
	return len(data), nil
}
func (tl testErrorLogger) Close() error {
	return nil
}

type testLogger struct {
	t testing.TB
}

func (tl testLogger) Write(data []byte) (int, error) {
	tl.t.Log(string(data))
	return len(data), nil
}
func (tl testLogger) Close() error {
	return nil
}

// InitTestingLogging initializes all logging functions to write their outputs
// to the logging output of the given testing context. If warningsAreErrors is
// true, warnings are logged using t.Error - marking the test as failed.
// Otherwise, errors are logged using t.Log. In either case, any errors logged
// cause the test to be marked as failed and are logged using t.Error, and the
// Info, Event, and Debug streams are always logged using t.Log.
func InitTestingLogging(tb testing.TB, warningsAreErrors bool) {
	tb.Helper()
	errWriter := testErrorLogger{t: tb}
	if writer := (testLogger{t: tb}); warningsAreErrors {
		Init(writer, errWriter, errWriter, writer, writer)
	} else {
		Init(writer, errWriter, writer, writer, writer)
	}
}
