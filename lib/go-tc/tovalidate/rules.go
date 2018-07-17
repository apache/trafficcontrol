// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tovalidate

import (
	"errors"
	"fmt"
	"strings"
)

// NoSpaces returns true if the string has no spaces
func NoSpaces(str string) bool {
	return !strings.ContainsAny(str, " ")
}

// NoPeriods returns true if the string has no periods
func NoPeriods(str string) bool {
	return !strings.ContainsAny(str, ".")
}

// IsOneOfString generates a validator function returning whether string is in the set of strings
func IsOneOfString(set ...string) func(string) bool {
	return func(s string) bool {
		for _, x := range set {
			if s == x {
				return true
			}
		}
		return false
	}
}

// IsOneOfStringICase is a case-insensitive version of IsOneOfString
func IsOneOfStringICase(set ...string) func(string) bool {
	var lowcased []string
	for _, s := range set {
		lowcased = append(lowcased, strings.ToLower(s))
	}
	return func(s string) bool {
		return IsOneOfString(lowcased...)(strings.ToLower(s))
	}
}

func IsGreaterThanZero(value interface{}) error {
	switch v := value.(type) {
	case *int:
		if v == nil || *v > 0 {
			return nil
		}
	case *float64:
		if v == nil || *v > 0 {
			return nil
		}
	default:
		return fmt.Errorf("IsGreaterThanZero validation failure: unknown type %T", value)
	}
	return errors.New("must be greater than zero")
}
