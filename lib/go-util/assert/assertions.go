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
		output = output + " Messages: " + message
	}
	return output
}

// Equal asserts that two objects are equal.
func Equal(t *testing.T, a, b interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	if a != b {
		msg := failureOutput(fmt.Sprintf("Not equal. Expected: %v Actual: %v", a, b), msgAndArgs...)
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
		msg := failureOutput("An error is expected but got nil.", msgAndArgs...)
		t.Error(msg)
		return false
	}
	return true
}

// Exactly asserts that two objects are equal in value and type.
func Exactly(t *testing.T, a, b interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	if !reflect.DeepEqual(a, b) {
		msg := failureOutput(fmt.Sprintf("Not exactly equal. Expected: %v (%T) Actual: %v (%T)", a, a, b, b), msgAndArgs...)
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
	msg := failureOutput(fmt.Sprintf("\"%v\" is not greater than or equal to \"%v\"", a, b), msgAndArgs...)
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
		msg := failureOutput(fmt.Sprintf("Received unexpected error: %+v", err), msgAndArgs...)
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
		msg := failureOutput("Expected value not to be nil.", msgAndArgs...)
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

// NotEqual asserts that two objects are NOT equal.
func NotEqual(t *testing.T, a, b interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	if a == b {
		msg := failureOutput(fmt.Sprintf("Should not be: %v", b), msgAndArgs...)
		t.Error(msg)
		return false
	}
	return true
}

// RequireNotEqual asserts that two objects are NOT equal.
// It marks the test as failed and stops execution.
func RequireNotEqual(t *testing.T, a, b interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if !NotEqual(t, a, b, msgAndArgs...) {
		t.FailNow()
	}
}

// Empty takes a value and checks whether it is empty for the given type
// supports slices, arrays, channels, strings, and maps
func Empty(t *testing.T, a interface{}) {
	val := reflect.ValueOf(a)
	switch val.Kind() {
	case reflect.Slice:
		if val.Len() != 0 {
			t.Errorf("expected %v to be empty, but it is not", a)
		}
	case reflect.Array:
		if val.Len() != 0 {
			t.Errorf("expected %v to be empty, but it is not", a)
		}
	case reflect.Chan:
		if val.Len() != 0 {
			t.Errorf("expected %v to be empty, but it is not", a)
		}
	case reflect.String:
		if val.Len() != 0 {
			t.Errorf("expected %v to be empty, but it is not", a)
		}
	case reflect.Map:
		if val.Len() != 0 {
			t.Errorf("expected %v to be empty, but it is not", a)
		}
	default:
		t.Errorf("can't check that %v of type %T is empty", a, a)
	}
}

// NotEmpty takes a value and checks whether it isvnot empty for the given type
// supports slices, arrays, channels, strings, and maps
func NotEmpty(t *testing.T, a interface{}) {
	val := reflect.ValueOf(a)
	switch val.Kind() {
	case reflect.Slice:
		if val.Len() == 0 {
			t.Errorf("expected %v to not be empty, but it is", a)
		}
	case reflect.Array:
		if val.Len() == 0 {
			t.Errorf("expected %v to not be empty, but it is", a)
		}
	case reflect.Chan:
		if val.Len() == 0 {
			t.Errorf("expected %v to not be empty, but it is", a)
		}
	case reflect.String:
		if val.Len() == 0 {
			t.Errorf("expected %v to not be empty, but it is", a)
		}
	case reflect.Map:
		if val.Len() == 0 {
			t.Errorf("expected %v to not be empty, but it is", a)
		}
	default:
		t.Errorf("can't check that %v of type %T is not empty", a, a)
	}
}
