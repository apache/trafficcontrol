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

// DeliveryServiceRegexResponse is the type of a Traffic Ops API response to
// the /deliveryservices_regexes endpoint - in all API verions (at the time of
// this writing).
type DeliveryServiceRegexResponse struct {
	Response []DeliveryServiceRegexes `json:"response"`
	Alerts
}

// DeliveryServiceRegexes structures associate a set of Delivery Service
// Regular Expressions (Delivery Service "Regexes") with a particular Delivery
// Service.
type DeliveryServiceRegexes struct {
	Regexes []DeliveryServiceRegex `json:"regexes"`
	// The XMLID of the Delivery Service to which the Regexes belong - NOT its
	// Display Name.
	DSName string `json:"dsName"`
}

// DeliveryServiceRegex is a regular expression used for routing to a Delivery
// Service.
type DeliveryServiceRegex struct {
	Type      string `json:"type"`
	SetNumber int    `json:"setNumber"`
	Pattern   string `json:"pattern"`
}

// DeliveryServiceIDRegexResponse is a list of DeliveryServiceIDRegexes.
type DeliveryServiceIDRegexResponse struct {
	Response []DeliveryServiceIDRegex `json:"response"`
	Alerts
}

// DeliveryServiceIDRegex holds information relating to a single routing regular
// expression of a delivery service, e.g., one of those listed at the
// deliveryservices/{{ID}}/regexes TO API route.
type DeliveryServiceIDRegex struct {
	ID        int    `json:"id"`
	Type      int    `json:"type"`
	TypeName  string `json:"typeName"`
	SetNumber int    `json:"setNumber"`
	Pattern   string `json:"pattern"`
}

// DeliveryServiceRegexesTest is used to represent the entire
// deliveryservice_regex for testing.
//
// This is ONLY meant to be used by testing code internal to ATC, do NOT use
// this to represent real CDN objects of any kind.
type DeliveryServiceRegexesTest struct {
	DSName string `json:"dsName"`
	DSID   int
	DeliveryServiceIDRegex
}

// DeliveryServiceRegexPost holds all of the information necessary to create a
// new routing regular expression for a delivery service.
type DeliveryServiceRegexPost struct {
	Type      int    `json:"type"`
	SetNumber int    `json:"setNumber"`
	Pattern   string `json:"pattern"`
}
