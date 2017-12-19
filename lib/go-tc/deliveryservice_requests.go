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

// GetDeliveryServiceRequestResponse ...
type GetDeliveryServiceRequestResponse struct {
	Response []DeliveryServiceRequest `json:"response"`
}

// CreateDeliveryServiceRequestResponse ...
type CreateDeliveryServiceRequestResponse struct {
	Response []DeliveryServiceRequest      `json:"response"`
	Alerts   []DeliveryServiceRequestAlert `json:"alerts"`
}

// UpdateDeliveryServiceRequestResponse ...
type UpdateDeliveryServiceRequestResponse struct {
	Response []DeliveryServiceRequest      `json:"response"`
	Alerts   []DeliveryServiceRequestAlert `json:"alerts"`
}

// DeliveryServiceRequestResponse ...
type DeliveryServiceRequestResponse struct {
	Response DeliveryServiceRequest        `json:"response"`
	Alerts   []DeliveryServiceRequestAlert `json:"alerts"`
}

// DeleteDeliveryServiceRequestResponse ...
type DeleteDeliveryServiceRequestResponse struct {
	Alerts []DeliveryServiceRequestAlert `json:"alerts"`
}

// DeliveryServiceRequestAlert ...
type DeliveryServiceRequestAlert struct {
	Level string `json:"level"`
	Text  string `json:"text"`
}

// DeliveryServiceRequest is used as part of the workflow to create, modify, or
// delete a delivery service.
type DeliveryServiceRequest struct {
	AssigneeID  int             `json:"assigneeId" db:"assignee_id"`
	AuthorID    int             `json:"authorId" db:"author_id"`
	ID          int             `json:"id" db:"id"`
	LastUpdated Time            `json:"lastUpdated" db:"last_updated"`
	Request     json.RawMessage `json:"request" db:"request"`
	Status      string          `json:"status" db:"status"`
}
