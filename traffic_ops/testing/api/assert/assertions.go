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
	"fmt"
	"reflect"
	"testing"
)

// failureOutput checks if there is a message to be parsed and concatenates with a failure message.
func failureOutput(failureMessage string, msgAndArgs ...interface{}) string {
	output := failureMessage
	message := ""
	if len(msgAndArgs) == 1 {
		msg := msgAndArgs[0]
		if msgAsStr, ok := msg.(string); ok {
			message = msgAsStr
		}
		message = fmt.Sprintf("%+v", msg)
	}
	if len(msgAndArgs) > 1 {
		message = fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
	if len(message) > 0 {
		output = "\n" + output + "\nMessages: " + message + "\n"
	}
	return output
}

// Equal asserts that two objects are equal.
func Equal(t *testing.T, a, b interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	if a != b {
		msg := failureOutput(fmt.Sprintf("Error: Not equal: \n expected: %v\n actual  : %v", a, b), msgAndArgs...)
		t.Error(msg)
		return false
	}
	return true
}

// RequireEqual asserts that two objects are equal.
// It marks the test as failed and stops execution.
func RequireEqual(t *testing.T, a, b interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if !Equal(t, a, b, msgAndArgs...) {
		t.FailNow()
	}
}

// Error asserts that a function returned an error (i.e. not `nil`).
func Error(t *testing.T, err error, msgAndArgs ...interface{}) bool {
	t.Helper()
	if err == nil {
		msg := failureOutput("Error: An error is expected but got nil.", msgAndArgs...)
		t.Error(msg)
		return false
	}
	return true
}

// Exactly asserts that two objects are equal in value and type.
func Exactly(t *testing.T, a, b interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	if !reflect.DeepEqual(a, b) {
		msg := failureOutput(fmt.Sprintf("Error: Not equal: \n expected: %v\n actual  : %v", a, b), msgAndArgs...)
		t.Error(msg)
		return false
	}
	return true
}

// GreaterOrEqual asserts that the first element is greater than or equal to the second
func GreaterOrEqual(t *testing.T, a, b int, msgAndArgs ...interface{}) bool {
	t.Helper()
	if a >= b {
		return true
	}
	msg := failureOutput(fmt.Sprintf("Error: \"%v\" is not greater than or equal to \"%v\"", a, b), msgAndArgs...)
	t.Error(msg)
	return false
}

// RequireGreaterOrEqual asserts that the first element is greater than or equal to the second
// It marks the test as failed and stops execution.
func RequireGreaterOrEqual(t *testing.T, a, b int, msgAndArgs ...interface{}) {
	t.Helper()
	if !GreaterOrEqual(t, a, b, msgAndArgs...) {
		t.FailNow()
	}
}

// NoError asserts that a function returned no error (i.e. `nil`).
func NoError(t *testing.T, err error, msgAndArgs ...interface{}) bool {
	t.Helper()
	if err != nil {
		msg := failureOutput(fmt.Sprintf("Received unexpected error:\n%+v", err), msgAndArgs...)
		t.Error(msg)
		return false
	}
	return true
}

// RequireNoError asserts that a function returned no error (i.e. `nil`).
// It marks the test as failed and stops execution.
func RequireNoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	if !NoError(t, err, msgAndArgs...) {
		t.FailNow()
	}
}

// NotNil asserts that the specified object is not nil.
func NotNil(t *testing.T, a interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	if a == nil {
		msg := failureOutput("Error: Expected value not to be nil.", msgAndArgs...)
		t.Error(msg)
		return false
	}
	return true
}

// RequireNotNil asserts that the specified object is not nil.
// It marks the test as failed and stops execution.
func RequireNotNil(t *testing.T, a interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if !NotNil(t, a, msgAndArgs...) {
		t.FailNow()
	}
}
