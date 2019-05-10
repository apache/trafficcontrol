package tc

import "github.com/lib/pq"

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

type CacheAssignmentGroupResponse struct {
	Response []CacheAssignmentGroup      `json:"response"`
	Alerts
	Size     int                         `json:"size"`
	Limit    int                         `json:"limit"`
}

type CacheAssignmentGroup struct {
	CDNID       int         `json:"cdnId" db:"cdn_id"`
	Description string   	`json:"description" db:"description"`
	ID          int         `json:"id" db:"id"`
	LastUpdated TimeNoMod	`json:"lastUpdated" db:"last_updated"`
	Name        string      `json:"name" db:"name"`
	Servers	    []int       `json:"servers"`
}

type CacheAssignmentGroupNullable struct {
	CDNID       *int        `json:"cdnId" db:"cdn_id"`
	Description *string 	`json:"description" db:"description"`
	ID          *int    	`json:"id" db:"id"`
	LastUpdated *TimeNoMod	`json:"lastUpdated" db:"last_updated"`
	Name        *string 	`json:"name" db:"name"`
	Servers	    pq.Int64Array      `json:"servers" db:"servers"`
}

