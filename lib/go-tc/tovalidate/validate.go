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
	"fmt"
)

// ToErrors converts a map of strings to errors into an array of errors.
//
// This is accomplished using `fmt.Errorf("'%v' %v", key, value)` where 'key'
// is the map key and 'value' is the error value to which it points - this
// means that error identity is NOT preserved. For example:
//
//     errMap := map[string]error{
//         "sql.ErrNoRows": sql.ErrNoRows,
//     }
//     errs := ToErrors(errMap)
//     if errors.Is(errs[0], sql.ErrNoRows) {
//         fmt.Println("true")
//     } else {
//         fmt.Println("false")
//     }
//
// ... will output 'false'.
func ToErrors(err map[string]error) []error {
	vErrors := []error{}
	for key, value := range err {
		if value != nil {
			errMsg := fmt.Errorf("'%v' %v", key, value)
			vErrors = append(vErrors, errMsg)
		}
	}
	return vErrors
}
