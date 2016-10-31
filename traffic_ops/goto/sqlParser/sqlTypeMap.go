
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sqlParser

/*****************************************************
 * DIRECTORY:
 * 1. StringToType(b []byte, t string) interface{} {}
      Given a value of []byte type, returns that value in
	  type specified in t string.
   2. TypeToString(data interface{}) string {}
      Given a value of generic type, returns it in string type.
	  Useful for constructing string queries.
 ****************************************************/
import (
	"errors"
	"strconv"
)

//given a []byte b and type t, return b in t form
func StringToType(b []byte, t string) (interface{}, error) {
	//all unregistered types (datetime for now, etc) are type string
	s := string(b)

	if t == "bigint" || t == "int" || t == "integer" || t == "tinyint" {
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		return i, nil
	} else if t == "double" {
		float, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		return float, nil
	} else if t == "varchar" {
		return s, nil
	} else {
		return string(b), nil
	}

}

//given data of generic type, returns data of string type
func TypeToString(data interface{}) (string, error) {
	if bigint, ok := data.(int64); ok {
		return strconv.Itoa(int(bigint)), nil
	} else if intv, ok := data.(int32); ok {
		return strconv.Itoa(int(intv)), nil
	} else if tinyint, ok := data.(uint8); ok {
		return strconv.Itoa(int(tinyint)), nil
	} else if double, ok := data.(float64); ok {
		return strconv.FormatFloat(double, 'f', 2, 32), nil
	} else if str, ok := data.(string); ok {
		return str, nil
	} else {
		err := errors.New("SQLPARSER: Whoa, what is this type?")
		return "", err
	}

	return "", nil
}
