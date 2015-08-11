// Copyright 2015 Comcast Cable Communications Management, LLC

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

type Wrapper struct {
	Resp    interface{} `json:"response"`
	Error   string      `json:"error"`
	IsTable bool        `json:"isTable"`
	Version float64     `json:"version"`
}

//wraps the given interface r into a returned Wrapper
//prepped for encoding to stream
func MakeWrapper(r interface{}, err string, isTable bool) Wrapper {
	//version is hard coded to "1.1"
	//all of this is variable
	w := Wrapper{r, err, isTable, 1.1}
	return w
}
