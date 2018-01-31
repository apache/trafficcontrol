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
	"strconv"
	"strings"

	log "github.com/apache/incubator-trafficcontrol/lib/go-log"
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
	AssigneeID      *int                     `json:"assigneeId,omitempty" db:"assignee_id"`
	Assignee        *string                  `json:"assignee,omitempty"`
	AuthorID        IDNoMod                  `json:"authorId" db:"author_id"`
	Author          string                   `json:"author"`
	ChangeType      string                   `json:"changeType" db:"change_type"`
	CreatedAt       *TimeNoMod               `json:"createdAt" db:"created_at"`
	ID              int                      `json:"id" db:"id"`
	LastEditedBy    string                   `json:"lastEditedBy"`
	LastEditedByID  IDNoMod                  `json:"lastEditedById" db:"last_edited_by_id"`
	LastUpdated     *TimeNoMod               `json:"lastUpdated" db:"last_updated"`
	DeliveryService *DeliveryServiceNullable `json:"deliveryService" db:"deliveryservice"`
	Status          RequestStatus            `json:"status" db:"status"`
	XMLID           string                   `json:"-" db:"xml_id"`
}

// UnmarshalJSON implements the json.Unmarshaller interface to suppress unmarshalling for IDNoMod
func (a *IDNoMod) UnmarshalJSON([]byte) error {
	return nil
}

// RequestStatus captures where in the workflow this request is
type RequestStatus int

const (
	// RequestStatusDraft -- newly created; not ready to be reviewed
	RequestStatusDraft = RequestStatus(iota) // default
	// RequestStatusSubmitted -- newly created; ready to be reviewed
	RequestStatusSubmitted
	// RequestStatusRejected -- reviewed, but problems found
	RequestStatusRejected
	// RequestStatusPending -- reviewed and locked; ready to be implemented
	RequestStatusPending
	// RequestStatusComplete -- implemented and locked
	RequestStatusComplete
	// RequestStatusInvalid -- placeholder
	RequestStatusInvalid = RequestStatus(-1)
)

// RequestStatusNames -- user-visible string associated with each of the above
var RequestStatusNames = [...]string{
	"draft",
	"submitted",
	"rejected",
	"pending",
	"complete",
}

// UnmarshalJSON implements json.Unmarshaller
func (r *RequestStatus) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	x, err := RequestStatusFromString(s)
	if err != nil {
		return err
	}
	log.Debugf("from string '%s':  status %v", string(s), x)
	*r = x
	return nil
}

// UnmarshalJSON implements json.Marshaller
func (r RequestStatus) MarshalJSON() ([]byte, error) {
	i := int(r)
	if i > len(RequestStatusNames) || i < 0 {
		return nil, errors.New("RequestStatus " + strconv.Itoa(i) + " out of range")
	}
	return json.Marshal(RequestStatusNames[i])
}

// RequestStatusFromString gets the status enumeration from a string
func RequestStatusFromString(s string) (RequestStatus, error) {
	if s == "" {
		return RequestStatusDraft, nil
	}
	t := strings.ToLower(s)
	for i, st := range RequestStatusNames {
		if t == st {
			return RequestStatus(i), nil
		}
	}
	return RequestStatusInvalid, errors.New(s + " is not a valid RequestStatus name")
}

// Name returns user-friendly string from the enumeration
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
