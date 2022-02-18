package assert

/*

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

import (
	"reflect"
	"testing"
)

// Equal asserts that two objects are equal.
func Equal(t *testing.T, a, b interface{}, msg string) bool {
	t.Helper()
	if a != b {
		t.Error(msg)
		return false
	}
	return true
}

// RequireEqual asserts that two objects are equal.
// It marks the test as failed and stops execution.
func RequireEqual(t *testing.T, a, b interface{}, msg string) {
	t.Helper()
	if !Equal(t, a, b, msg) {
		t.FailNow()
	}
}

// Error asserts that a function returned an error (i.e. not `nil`).
func Error(t *testing.T, err error, msg string) bool {
	t.Helper()
	if err == nil {
		t.Error(msg)
		return false
	}
	return true
}

// Exactly asserts that two objects are equal in value and type.
func Exactly(t *testing.T, a, b interface{}, msg string) bool {
	t.Helper()
	if !reflect.DeepEqual(a, b) {
		t.Error(msg)
		return false
	}
	return true
}

// GreaterOrEqual asserts that the first element is greater than or equal to the second
func GreaterOrEqual(t *testing.T, a, b int, msg string) bool {
	t.Helper()
	if a > b {
		return true
	}
	return Equal(t, a, b, msg)
}

// RequireGreaterOrEqual asserts that the first element is greater than or equal to the second
// It marks the test as failed and stops execution.
func RequireGreaterOrEqual(t *testing.T, a, b int, msg string) {
	t.Helper()
	if !GreaterOrEqual(t, a, b, msg) {
		t.FailNow()
	}
}

// NoError asserts that a function returned no error (i.e. `nil`).
func NoError(t *testing.T, err error, msg string) bool { t.Helper(); return Equal(t, err, nil, msg) }

// RequireNoError asserts that a function returned no error (i.e. `nil`).
// It marks the test as failed and stops execution.
func RequireNoError(t *testing.T, err error, msg string) {
	t.Helper()
	if !NoError(t, err, msg) {
		t.FailNow()
	}
}

// NotNil asserts that the specified object is not nil.
func NotNil(t *testing.T, a interface{}, msg string) bool {
	t.Helper()
	if a == nil {
		t.Error(msg)
		return false
	}
	return true
}

// RequireNotNil asserts that the specified object is not nil.
// It marks the test as failed and stops execution.
func RequireNotNil(t *testing.T, a interface{}, msg string) {
	t.Helper()
	if !NotNil(t, a, msg) {
		t.FailNow()
	}
}
