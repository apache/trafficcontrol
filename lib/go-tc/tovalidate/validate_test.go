package tovalidate

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"errors"
	"fmt"
	"testing"
)

func ExampleToError() {
	errs := map[string]error{
		"propA": errors.New("bad value"),
		"propB": errors.New("cannot be blank"),
	}
	err := ToError(errs).Error()
	// Iteration order of Go maps is random, so this is the best we can do.
	fmt.Println(
		err == "'propA' bad value, 'propB' cannot be blank" ||
			err == "'propB' cannot be blank, 'propA' bad value",
	)
	// Output: true
}

func TestToError(t *testing.T) {
	var errs map[string]error
	err := ToError(errs)
	if err != nil {
		t.Error("a nil error map should yield a nil error, got:", err)
	}
	errs = map[string]error{}
	err = ToError(errs)
	if err != nil {
		t.Error("an empty error map should yield a nil error, got:", err)
	}
	errs["something"] = nil
	err = ToError(errs)
	if err != nil {
		t.Error("an error map with no non-nil errors should yield a nil error, got:", err)
	}
}
