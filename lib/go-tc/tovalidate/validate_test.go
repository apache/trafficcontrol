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
)

func ExampleToError() {
	errs := map[string]error{
		"propA": errors.New("bad value"),
		"propB": errors.New("cannot be blank"),
	}
	fmt.Println(ToError(errs))
	// Output: 'propA' bad value, 'propB' cannot be blank
}
