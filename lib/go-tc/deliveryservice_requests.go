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
	"errors"
	"strings"
)

// IDNoMod type is used to suppress JSON unmarshalling
type IDNoMod int

// DeliveryServiceRequest is used as part of the workflow to create,
// modify, or delete a delivery service.
type DeliveryServiceRequest struct {
	AssigneeID      int             `json:"assigneeId,omitempty"`
	Assignee        string          `json:"assignee,omitempty"`
	AuthorID        IDNoMod         `json:"authorId"`
	Author          string          `json:"author"`
	ChangeType      string          `json:"changeType"`
	CreatedAt       *TimeNoMod      `json:"createdAt"`
	ID              int             `json:"id"`
	LastEditedBy    string          `json:"lastEditedBy,omitempty"`
	LastEditedByID  IDNoMod         `json:"lastEditedById,omitempty"`
	LastUpdated     *TimeNoMod      `json:"lastUpdated"`
	DeliveryService DeliveryService `json:"deliveryService"`
	Status          RequestStatus   `json:"status"`
	XMLID           string          `json:"-" db:"xml_id"`
}

// DeliveryServiceRequestNullable is used as part of the workflow to create,
// modify, or delete a delivery service.
type DeliveryServiceRequestNullable struct {
	AssigneeID      *int            `json:"assigneeId,omitempty" db:"assignee_id"`
	Assignee        *string         `json:"assignee,omitempty"`
	AuthorID        IDNoMod         `json:"authorId" db:"author_id"`
	Author          string          `json:"author"`
	ChangeType      string          `json:"changeType" db:"change_type"`
	CreatedAt       *TimeNoMod      `json:"createdAt" db:"created_at"`
	ID              int             `json:"id" db:"id"`
	LastEditedBy    string          `json:"lastEditedBy"`
	LastEditedByID  IDNoMod         `json:"lastEditedById" db:"last_edited_by_id"`
	LastUpdated     *TimeNoMod      `json:"lastUpdated" db:"last_updated"`
	DeliveryService json.RawMessage `json:"deliveryService" db:"deliveryservice"`
	Status          RequestStatus   `json:"status" db:"status"`
	XMLID           string          `json:"-" db:"xml_id"`
}

// UnmarshalJSON implements the json.Unmarshaller interface to suppress unmarshalling for IDNoMod
func (a *IDNoMod) UnmarshalJSON([]byte) error {
	return nil
}

//
type RequestStatus int

const (
	RequestStatusDraft = RequestStatus(iota) // default
	RequestStatusSubmitted
	RequestStatusRejected
	RequestStatusPending
	RequestStatusComplete
	RequestStatusInvalid = RequestStatus(-1)
)

var RequestStatusNames = [...]string{
	"draft",
	"submitted",
	"rejected",
	"pending",
	"complete",
}

func RequestStatusFromString(s string) RequestStatus {
	t := strings.ToLower(s)
	for i, st := range RequestStatusNames {
		if t == st {
			return RequestStatus(i)
		}
	}
	return RequestStatusInvalid
}

func (s RequestStatus) Name() string {
	i := int(s)
	if i < 0 || i > len(RequestStatusNames) {
		return "INVALID"
	}
	return RequestStatusNames[i]
}

// ValidTransition returns nil if the transition is allowed for the workflow, an error if not
func (s RequestStatus) ValidTransition(to RequestStatus) error {
	if s == to {
		// no change -- always allowed
		return nil
	}

	// indicate if valid transitioning to this RequestStatus
	switch to {
	case RequestStatusDraft:
		// can go back to draft if submitted or rejected
		if s == RequestStatusSubmitted || s == RequestStatusRejected {
			return nil
		}
	case RequestStatusSubmitted:
		// can go be submitted if draft or rejected
		if s == RequestStatusDraft || s == RequestStatusRejected {
			return nil
		}
	case RequestStatusRejected:
		// only submitted can be rejected
		if s == RequestStatusSubmitted {
			return nil
		}
	case RequestStatusPending:
		// only submitted can move to pending
		if s == RequestStatusSubmitted {
			return nil
		}
	case RequestStatusComplete:
		// only pending can be completed.  Completed can never change.
		if s == RequestStatusPending {
			return nil
		}
	}
	return errors.New("invalid transition from " + s.Name() + " to " + to.Name())
}
