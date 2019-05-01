package tc

import (
	"database/sql"
)

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

type FoosResponse struct {
	Response []Foo `json:"response"`
}

type CreateFooResponse struct {
	Response []Foo `json:"response"`
	Alerts
}

type UpdateFooResponse struct {
	Response []Foo `json:"response"`
	Alerts
}

type FooResponse struct {
	Response Foo `json:"response"`
	Alerts
}

type DeleteFooResponse struct {
	Alerts
}

type FooV19 Foo // this type alias should always point to the latest minor version

type Foo struct {
	FooV18
	E *string `json:"E"`
}

type FooV18 struct {
	FooV17
	D *string `json:"D"`
}

type FooV17 struct {
	FooV16
	C *string `json:"C"`
}

type FooV16 struct { // FooV16
	FooV15
	B *string `json:"B"`
}

type FooV15 struct {
	ID   *int    `json:"id"`
	Name *string `json:"name"`
	A    *string `json:"A"`
}

func (foo *Foo) Sanitize() {
	// TODO: populate defaults for optional values
}

func (foo *Foo) Validate(tx *sql.Tx) error {
	foo.Sanitize()
	// TODO: add some validation
	return nil
}
