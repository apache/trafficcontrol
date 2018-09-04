package tc

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
// DeliveryServiceRegexResponse ...
type DeliveryServiceRegexResponse struct {
	Response []DeliveryServiceRegexes `json:"response"`
}

// DeliveryServiceRegexes ...
type DeliveryServiceRegexes struct {
	Regexes []DeliveryServiceRegex `json:"regexes"`
	DSName  string                 `json:"dsName"`
}

// DeliveryServiceRegex ...
type DeliveryServiceRegex struct {
	Type      string `json:"type"`
	SetNumber int    `json:"setNumber"`
	Pattern   string `json:"pattern"`
}

type DeliveryServiceIDRegexResponse struct {
	Response []DeliveryServiceIDRegex `json:"response"`
}

type DeliveryServiceIDRegex struct {
	ID        int    `json:"id"`
	Type      int    `json:"type"`
	TypeName  string `json:"typeName"`
	SetNumber int    `json:"setNumber"`
	Pattern   string `json:"pattern"`
}

type DeliveryServiceRegexPost struct {
	Type      int    `json:"type"`
	SetNumber int    `json:"setNumber"`
	Pattern   string `json:"pattern"`
}
