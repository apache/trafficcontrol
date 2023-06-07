package tc

import (
	"encoding/json"
	"errors"
	"github.com/apache/trafficcontrol/lib/go-util"
	"time"
)

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

type DeliveryServiceRequestCommentV5 DeliveryServiceRequestCommentV50

type DeliveryServiceRequestCommentV50 struct {
	AuthorID                 IDNoMod   `json:"authorId" db:"author_id"`
	Author                   string    `json:"author"`
	DeliveryServiceRequestID int       `json:"deliveryServiceRequestId" db:"deliveryservice_request_id"`
	ID                       int       `json:"id" db:"id"`
	LastUpdated              time.Time `json:"lastUpdated" db:"last_updated"`
	Value                    string    `json:"value" db:"value"`
	XMLID                    string    `json:"xmlId" db:"xml_id"`
}

type DeliveryServiceRequestCommentsResponseV5 DeliveryServiceRequestCommentsResponseV50

type DeliveryServiceRequestCommentsResponseV50 struct {
	Response []DeliveryServiceRequestCommentV5 `json:"response"`
	Alerts
}

func (d DeliveryServiceRequestCommentsResponse) Upgrade() (DeliveryServiceRequestCommentsResponseV5, error) {
	deliveryServiceRequestCommentsResponse := DeliveryServiceRequestCommentsResponseV5{}
	deliveryServiceRequestCommentsResponse.Response = make([]DeliveryServiceRequestCommentV5, 0)
	for _, dsrc := range d.Response {
		dsrcV5, err := dsrc.Upgrade()
		if err != nil {
			return DeliveryServiceRequestCommentsResponseV5{}, err
		}
		deliveryServiceRequestCommentsResponse.Response = append(deliveryServiceRequestCommentsResponse.Response, dsrcV5)
	}
	deliveryServiceRequestCommentsResponse.Alerts = d.Alerts
	return deliveryServiceRequestCommentsResponse, nil
}

func (d DeliveryServiceRequestCommentsResponseV5) Downgrade() DeliveryServiceRequestCommentsResponse {
	deliveryServiceRequestCommentsResponse := DeliveryServiceRequestCommentsResponse{}
	deliveryServiceRequestCommentsResponse.Response = make([]DeliveryServiceRequestComment, 0)
	for _, dsrc := range d.Response {
		dsrcLegacy := dsrc.Downgrade()
		deliveryServiceRequestCommentsResponse.Response = append(deliveryServiceRequestCommentsResponse.Response, dsrcLegacy)
	}
	deliveryServiceRequestCommentsResponse.Alerts = d.Alerts
	return deliveryServiceRequestCommentsResponse
}

func (d DeliveryServiceRequestCommentV5) DowngradeToNullable() DeliveryServiceRequestCommentNullable {
	t := TimeNoModFromTime(d.LastUpdated)
	deliveryServiceRequestComment := DeliveryServiceRequestCommentNullable{
		AuthorID:                 &d.AuthorID,
		Author:                   &d.Author,
		DeliveryServiceRequestID: &d.DeliveryServiceRequestID,
		ID:                       &d.ID,
		LastUpdated:              t,
		Value:                    &d.Value,
		XMLID:                    &d.XMLID,
	}
	return deliveryServiceRequestComment
}

func (d DeliveryServiceRequestCommentNullable) UpgradeFromNullable() (DeliveryServiceRequestCommentV5, error) {
	deliveryServiceRequestCommentV5 := DeliveryServiceRequestCommentV5{}
	if d.AuthorID != nil {
		deliveryServiceRequestCommentV5.AuthorID = *d.AuthorID
	}
	if d.Author != nil {
		deliveryServiceRequestCommentV5.Author = *d.Author
	}
	if d.DeliveryServiceRequestID != nil {
		deliveryServiceRequestCommentV5.DeliveryServiceRequestID = *d.DeliveryServiceRequestID
	}
	if d.ID != nil {
		deliveryServiceRequestCommentV5.ID = *d.ID
	}
	if d.LastUpdated != nil {
		t := d.LastUpdated.Time
		updatedTime, err := util.ConvertTimeFormat(t, time.RFC3339)
		if err != nil {
			return DeliveryServiceRequestCommentV5{}, err
		}
		deliveryServiceRequestCommentV5.LastUpdated = *updatedTime
	}
	if d.Value != nil {
		deliveryServiceRequestCommentV5.Value = *d.Value
	}
	if d.XMLID != nil {
		deliveryServiceRequestCommentV5.XMLID = *d.XMLID
	}
	return deliveryServiceRequestCommentV5, nil
}

func (d DeliveryServiceRequestComment) Upgrade() (DeliveryServiceRequestCommentV5, error) {
	t := d.LastUpdated.Time
	updatedTime, err := util.ConvertTimeFormat(t, time.RFC3339)
	if err != nil {
		return DeliveryServiceRequestCommentV5{}, err
	}
	deliveryServiceRequestCommentV5 := DeliveryServiceRequestCommentV5{
		Author:                   d.Author,
		AuthorID:                 d.AuthorID,
		DeliveryServiceRequestID: d.DeliveryServiceRequestID,
		ID:                       d.ID,
		LastUpdated:              *updatedTime,
		Value:                    d.Value,
		XMLID:                    d.XMLID,
	}
	return deliveryServiceRequestCommentV5, nil
}

func (d DeliveryServiceRequestCommentV5) Downgrade() DeliveryServiceRequestComment {
	t := TimeNoModFromTime(d.LastUpdated)
	return DeliveryServiceRequestComment{
		AuthorID:                 d.AuthorID,
		Author:                   d.Author,
		DeliveryServiceRequestID: d.DeliveryServiceRequestID,
		ID:                       d.ID,
		LastUpdated:              *t,
		Value:                    d.Value,
		XMLID:                    d.XMLID,
	}
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface with a
// customized decoding to force the date format on LastUpdated.
func (d *DeliveryServiceRequestCommentV5) UnmarshalJSON(data []byte) error {
	type Alias DeliveryServiceRequestCommentV5
	resp := struct {
		LastUpdated string `json:"lastUpdated"`
		*Alias
	}{
		Alias: (*Alias)(d),
	}
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return err
	}

	d.LastUpdated, err = parseTimeV5(resp.LastUpdated)
	if err != nil {
		return errors.New("invalid timestamp given for lastUpdated: " + err.Error())
	}
	return nil
}

func parseTimeV5(ts string) (time.Time, error) {
	rt, err := time.Parse(time.RFC3339, ts)
	if err == nil {
		return rt, err
	}
	return time.Parse(dateFormat, ts)
}

// MarshalJSON implements the encoding/json.Marshaler interface with a
// customized encoding to force the date format on LastUpdated.
func (d DeliveryServiceRequestCommentV5) MarshalJSON() ([]byte, error) {
	type Alias DeliveryServiceRequestCommentV5
	resp := struct {
		LastUpdated string `json:"lastUpdated"`
		Alias
	}{
		LastUpdated: d.LastUpdated.Format(time.RFC3339),
		Alias:       (Alias)(d),
	}
	return json.Marshal(&resp)
}

// DeliveryServiceRequestCommentsResponse is a list of
// DeliveryServiceRequestComments as a response.
type DeliveryServiceRequestCommentsResponse struct {
	Response []DeliveryServiceRequestComment `json:"response"`
	Alerts
}

// DeliveryServiceRequestComment is a struct containing the fields for a delivery
// service request comment.
type DeliveryServiceRequestComment struct {
	AuthorID                 IDNoMod   `json:"authorId" db:"author_id"`
	Author                   string    `json:"author"`
	DeliveryServiceRequestID int       `json:"deliveryServiceRequestId" db:"deliveryservice_request_id"`
	ID                       int       `json:"id" db:"id"`
	LastUpdated              TimeNoMod `json:"lastUpdated" db:"last_updated"`
	Value                    string    `json:"value" db:"value"`
	XMLID                    string    `json:"xmlId" db:"xml_id"`
}

// DeliveryServiceRequestCommentNullable is a nullable struct containing the
// fields for a delivery service request comment.
type DeliveryServiceRequestCommentNullable struct {
	AuthorID                 *IDNoMod   `json:"authorId" db:"author_id"`
	Author                   *string    `json:"author"`
	DeliveryServiceRequestID *int       `json:"deliveryServiceRequestId" db:"deliveryservice_request_id"`
	ID                       *int       `json:"id" db:"id"`
	LastUpdated              *TimeNoMod `json:"lastUpdated" db:"last_updated"`
	Value                    *string    `json:"value" db:"value"`
	XMLID                    *string    `json:"xmlId" db:"xml_id"`
}
