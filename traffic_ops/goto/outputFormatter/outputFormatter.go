
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package outputFormatter

/******************************************************************
outputFormatter contains:
* Wrapper struct, which is written to the stream in server.go
  * * Resp interface{}, which is the response to the user query
    * Version, which is the version of the API
* MakeWrapper(r interface{}), which wraps r into a struct to encode
	*****************************************************************/
//import "fmt"

type ApiWrapper struct {
	Resp        interface{}       `json:"response"`
	Cols        map[string]Column `json:"columns"`
	ColWrappers []ColumnWrapper   `json:"colWrappers"`
	Error       string            `json:"error"`
	IsTable     bool              `json:"isTable"`
	Version     float64           `json:"version"`
}

//wraps the given interface r into a returned Wrapper
//prepped for encoding to stream
func MakeApiWrapper(r interface{}, c []string, ca []string, cm map[string]map[string]interface{}, err string, isTable bool) ApiWrapper {
	//version is hard coded to "1.1"
	//all of this is variable
	w := ApiWrapper{r, MakeColumns(c, ca, cm), MakeColumnWrappers(ca), err, isTable, 1.1}
	return w
}

type ColumnWrapper struct {
	Field        string `json:"field"`
	DisplayName  string `json:"displayName"`
	ColumnFilter bool   `json:"columnFilter"`
}

type Column struct {
	Alias            string                 `json:"colAlias"`
	ForeignKey       bool                   `json:"isForeignKey"`
	ForeignKeyValues map[string]interface{} `json:"foreignKeyValues"`
}

func MakeColumns(columns []string, aliases []string, fkMap map[string]map[string]interface{}) map[string]Column {
	c := make(map[string]Column)

	for idx, column := range columns {
		var w Column
		if fkMapVals, ok := fkMap[column]; ok {
			w = Column{aliases[idx], true, fkMapVals}
		} else {
			w = Column{column, false, nil}
		}
		c[column] = w
	}
	return c
}

func MakeColumnWrappers(columns []string) []ColumnWrapper {
	cw := make([]ColumnWrapper, 0)
	for _, column := range columns {
		w := ColumnWrapper{column, column, true}
		cw = append(cw, w)
	}

	return cw
}
