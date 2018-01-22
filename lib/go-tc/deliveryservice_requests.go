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

import (
	"encoding/json"
)

// DeliveryServiceRequestNullable is used as part of the workflow to create,
// modify, or delete a delivery service.
type DeliveryServiceRequest struct {
	AssigneeID     int             `json:"assigneeId,omitempty" db:"assignee_id"`
	Assignee       string          `json:"assignee,omitempty"`
	AuthorID       int             `json:"authorId" db:"author_id"`
	Author         string          `json:"author"`
	ChangeType     string          `json:"changeType" db:"change_type"`
	CreatedAt      string          `json:"createdAt,omitempty" db:"created_at"`
	ID             int             `json:"id" db:"id"`
	LastEditedBy   string          `json:"lastEditedBy"`
	LastEditedByID int             `json:"lastEditedById" db:"last_edited_by_id"`
	LastUpdated    string          `json:"lastUpdated,omitempty" db:"last_updated"`
	Request        json.RawMessage `json:"request" db:"request"`
	Status         string          `json:"status" db:"status"`
}

// DeliveryServiceRequestNullable is used as part of the workflow to create,
// modify, or delete a delivery service.
type DeliveryServiceRequestNullable struct {
	AssigneeID     *int            `json:"assigneeId,omitempty" db:"assignee_id"`
	Assignee       *string         `json:"assignee,omitempty"`
	AuthorID       int             `json:"authorId" db:"author_id"`
	Author         string          `json:"author"`
	ChangeType     string          `json:"changeType" db:"change_type"`
	CreatedAt      Time            `json:"createdAt" db:"created_at"`
	ID             int             `json:"id" db:"id"`
	LastEditedBy   string          `json:"lastEditedBy"`
	LastEditedByID int             `json:"lastEditedById" db:"last_edited_by_id"`
	LastUpdated    Time            `json:"lastUpdated" db:"last_updated"`
	Request        json.RawMessage `json:"request" db:"request"`
	Status         string          `json:"status" db:"status"`
}
