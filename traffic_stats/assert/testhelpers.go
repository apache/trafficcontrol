/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package assert

import (
	"reflect"
	"testing"
)

// NotNil takes a value, and checks whether that value is not nil, failing the test
// if not
func NotNil(t *testing.T, a interface{}) {
	if a == nil {
		t.Errorf("expected %v to be non-nil, but it is", a)
	}
}

// Nil takes a value, and checks whether that value is nil, failing the test
// if not
func Nil(t *testing.T, a interface{}) {
	if a != nil {
		t.Errorf("expected %v to be nil, but it's not", a)
	}
}

// Equal takes two values, checks whether they are equal (using a deep equal
// check) and fails the test if not
func Equal(t *testing.T, a, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("expected %v to equal %v, but it does not", a, b)
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
