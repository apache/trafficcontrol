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
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	log "github.com/apache/trafficcontrol/lib/go-log"
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
	DeliveryService DeliveryService `json:"deliveryService"` // TODO version DeliveryServiceRequest
	Status          RequestStatus   `json:"status"`
	XMLID           string          `json:"-" db:"xml_id"`
}

// DeliveryServiceRequestNullable is used as part of the workflow to create,
// modify, or delete a delivery service.
type DeliveryServiceRequestNullable struct {
	AssigneeID      *int                     `json:"assigneeId,omitempty" db:"assignee_id"`
	Assignee        *string                  `json:"assignee,omitempty"`
	AuthorID        *IDNoMod                 `json:"authorId" db:"author_id"`
	Author          *string                  `json:"author"`
	ChangeType      *string                  `json:"changeType" db:"change_type"`
	CreatedAt       *TimeNoMod               `json:"createdAt" db:"created_at"`
	ID              *int                     `json:"id" db:"id"`
	LastEditedBy    *string                  `json:"lastEditedBy"`
	LastEditedByID  *IDNoMod                 `json:"lastEditedById" db:"last_edited_by_id"`
	LastUpdated     *TimeNoMod               `json:"lastUpdated" db:"last_updated"`
	DeliveryService *DeliveryServiceNullable `json:"deliveryService" db:"deliveryservice"` // TODO version DeliveryServiceRequest
	Status          *RequestStatus           `json:"status" db:"status"`
	XMLID           *string                  `json:"-" db:"xml_id"`
}

// UnmarshalJSON implements the json.Unmarshaller interface to suppress unmarshalling for IDNoMod
func (a *IDNoMod) UnmarshalJSON([]byte) error {
	return nil
}

// RequestStatus captures where in the workflow this request is
type RequestStatus string

const (
	// RequestStatusInvalid -- invalid state
	RequestStatusInvalid = RequestStatus("invalid")
	// RequestStatusDraft -- newly created; not ready to be reviewed
	RequestStatusDraft = RequestStatus("draft")
	// RequestStatusSubmitted -- newly created; ready to be reviewed
	RequestStatusSubmitted = RequestStatus("submitted")
	// RequestStatusRejected -- reviewed, but problems found
	RequestStatusRejected = RequestStatus("rejected")
	// RequestStatusPending -- reviewed and locked; ready to be implemented
	RequestStatusPending = RequestStatus("pending")
	// RequestStatusComplete -- implemented and locked
	RequestStatusComplete = RequestStatus("complete")
)

// RequestStatuses -- user-visible string associated with each of the above
var RequestStatuses = []RequestStatus{
	// "invalid" -- don't list here..
	"draft",
	"submitted",
	"rejected",
	"pending",
	"complete",
}

// UnmarshalJSON implements json.Unmarshaller
func (r *RequestStatus) UnmarshalJSON(b []byte) error {
	u, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}

	// just check to see if the string represents a valid requeststatus
	_, err = RequestStatusFromString(u)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, (*string)(r))
}

// MarshalJSON implements json.Marshaller
func (r RequestStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(r))
}

// Value implements driver.Valuer
func (r *RequestStatus) Value() (driver.Value, error) {
	v, err := json.Marshal(r)
	log.Debugf("value is %v; err is %v", v, err)
	v = []byte(strings.Trim(string(v), `"`))
	return v, err
}

// Scan implements sql.Scanner
func (r *RequestStatus) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("expected requeststatus in byte array form; got %T", src)
	}
	b = []byte(`"` + string(b) + `"`)
	return json.Unmarshal(b, r)
}

// RequestStatusFromString gets the status enumeration from a string
func RequestStatusFromString(rs string) (RequestStatus, error) {
	if rs == "" {
		return RequestStatusDraft, nil
	}
	for _, s := range RequestStatuses {
		if string(s) == rs {
			return s, nil
		}
	}
	return RequestStatusInvalid, errors.New(rs + " is not a valid RequestStatus name")
}

// ValidTransition returns nil if the transition is allowed for the workflow, an error if not
func (r RequestStatus) ValidTransition(to RequestStatus) error {
	if r == RequestStatusRejected || r == RequestStatusComplete {
		// once rejected or completed,  no changes allowed
		return errors.New(string(r) + " request cannot be changed")
	}

	if r == to {
		// no change -- always allowed
		return nil
	}

	// indicate if valid transitioning to this RequestStatus
	switch to {
	case RequestStatusDraft:
		// can go back to draft if submitted or rejected
		if r == RequestStatusSubmitted {
			return nil
		}
	case RequestStatusSubmitted:
		// can go be submitted if draft or rejected
		if r == RequestStatusDraft {
			return nil
		}
	case RequestStatusRejected:
		// only submitted can be rejected
		if r == RequestStatusSubmitted {
			return nil
		}
	case RequestStatusPending:
		// only submitted can move to pending
		if r == RequestStatusSubmitted {
			return nil
		}
	case RequestStatusComplete:
		// only submitted or pending requests can be completed
		if r == RequestStatusSubmitted || r == RequestStatusPending {
			return nil
		}
	}
	return errors.New("invalid transition from " + string(r) + " to " + string(to))
}
