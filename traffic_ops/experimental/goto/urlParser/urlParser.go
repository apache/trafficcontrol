
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package urlParser

/******************************************************************
* urlParser contains:
* * type Request struct,  which stores request information
  * * table/view from which information is queried
  * * parameters (in []string) that narrows table/view query field
* * ParseURL(urlString), which takes in a URL and parses it into a Request
*****************************************************************/

import (
	"strings"
)

type Request struct {
	Type string
	//can be for a table or a view
	TableName string
	//ex. "cachegroup < 50"
	//ex. "cachegroup >= 50"
	Parameters []string
}

//makes a new request given a string url
// XXX: why are we not using r.Url.Query?
func ParseURL(url string) Request {
	r := Request{"", "", make([]string, 0)}

	url = strings.ToLower(url)

	//replace less than/greater than symbols in url encode
	url = strings.Replace(url, "%3c", "<", -1)
	url = strings.Replace(url, "%3e", ">", -1)

	urlSections := strings.Split(url, "/")

	r.Type = urlSections[0]

	//title exists
	if len(urlSections) > 1 {
		titleParamStr := urlSections[1]

		// splits table name and parameters by "?"
		qMarkSplit := strings.Split(titleParamStr, "?")
		r.TableName = qMarkSplit[0]

		// if parameters exist, separate by "&"
		if len(qMarkSplit) > 1 {
			paramSplit := strings.Split(qMarkSplit[1], "&")
			for _, param := range paramSplit {
				if strings.HasPrefix(param, "format=") || strings.HasPrefix(param, "join=") {
					continue
				}
				r.Parameters = append(r.Parameters, param)
			}
		}
	}

	//second potential urlSection (after tableName & parameters) is specified id
	//by nature of SQLParser, this is considered as a parameter
	if len(urlSections) > 2 && urlSections[2] != "" {
		r.Parameters = append(r.Parameters, "id="+urlSections[2])
	}

	return r
}
