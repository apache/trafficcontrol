
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mooseFixture

// loosely after the json package.

/******************************************************************
* mooseFixture contains:
*****************************************************************/

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var camelingRegex = regexp.MustCompile("[0-9A-Za-z]+")

func UpperCamelCase(src string) string {
	byteSrc := []byte(src)
	chunks := camelingRegex.FindAll(byteSrc, -1)
	for idx, val := range chunks {
		if idx > 0 {
			chunks[idx] = bytes.Title(val)
		}
	}
	camel := string(bytes.Join(chunks, nil))
	return strings.ToUpper(string(camel[0])) + camel[1:]
}

type Encoder struct {
	w   io.Writer
	err error
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (enc *Encoder) Encode(tableName string, v interface{}) error {
	if enc.err != nil {
		return enc.err
	}

	tableName = UpperCamelCase(tableName) // strings.ToUpper(string(tableName[0])) + tableName[1:]
	enc.w.Write([]byte("package Fixtures::Integration::" + tableName + ";\n\n"))
	enc.w.Write([]byte("# Do not edit! Generated code.\n"))
	enc.w.Write([]byte("# See https://github.com/Comcast/traffic_control/wiki/The%20Kabletown%20example\n\n"))
	enc.w.Write([]byte("use Moose;\n"))
	enc.w.Write([]byte("extends 'DBIx::Class::EasyFixture';\n"))
	enc.w.Write([]byte("use namespace::autoclean;\n\n"))
	enc.w.Write([]byte("my %definition_for = (\n"))

	m := v.([]map[string]interface{})
	for rowNum, rowMap := range m {
		enc.w.Write([]byte("'" + strconv.Itoa(rowNum) + "' => { new => '" + tableName + "', => using => { "))

		for key, val := range rowMap {
			var keyval string
			var ok bool
			if keyval, ok = val.(string); ok {
				enc.w.Write([]byte(key + " => '" + keyval + "', "))
			} else if val == nil {
				enc.w.Write([]byte(key + " => undef, "))
			} else {

				fmt.Println("Error on ", rowMap["id"], " key ", key, " - not a string!")
			}
		}
		enc.w.Write([]byte("}, }, \n"))
	}

	enc.w.Write([]byte("); \n\n"))
	enc.w.Write([]byte("sub name {\n		return \"" + tableName + "\";\n}\n\n"))
	enc.w.Write([]byte("sub get_definition { \n		my ( $self, $name ) = @_;\n		return $definition_for{$name};\n}\n\n"))
	enc.w.Write([]byte("sub all_fixture_names {\n		return keys %definition_for;\n}\n\n"))
	enc.w.Write([]byte("__PACKAGE__->meta->make_immutable;\n"))
	enc.w.Write([]byte("1;\n"))

	return enc.err
}
