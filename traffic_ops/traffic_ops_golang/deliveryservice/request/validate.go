package request

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

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservice"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/tovalidate"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
)

type requestStatus int

const (
	statusDraft = requestStatus(iota)
	statusSubmitted
	statusRejected
	statusPending
	statusComplete
	statusInvalid = requestStatus(-1)
)

var statusNames = [...]string{
	"draft",
	"submitted",
	"rejected",
	"pending",
	"complete",
}

func statusFromString(s string) requestStatus {
	t := strings.ToLower(s)
	for i, st := range statusNames {
		if t == st {
			return requestStatus(i)
		}
	}
	return statusInvalid
}

func (s requestStatus) name() string {
	i := int(s)
	if i < 0 || i > len(statusNames) {
		return "INVALID"
	}
	return statusNames[i]
}

// validTransition returns nil if the transition is allowed for the workflow, or an error if not
func (s requestStatus) validTransition(to requestStatus) error {
	if s == to {
		// no change -- always allowed
		return nil
	}

	// indicate if valid transitioning to this requestStatus
	switch to {
	case statusDraft:
		// can go back to draft if submitted or rejected
		if s == statusSubmitted || s == statusRejected {
			return nil
		}
	case statusSubmitted:
		// can go be submitted if draft or rejected
		if s == statusDraft || s == statusRejected {
			return nil
		}
	case statusRejected:
		// only submitted can be rejected
		if s == statusSubmitted {
			return nil
		}
	case statusPending:
		// only submitted can move to pending
		if s == statusSubmitted {
			return nil
		}
	case statusComplete:
		// only pending can be completed.  Completed can never change.
		if s == statusPending {
			return nil
		}
	}
	return errors.New("invalid transition from " + s.name() + " to " + to.name())
}

// Validate ensures all required fields are present and in correct form.  Also checks request JSON is complete and valid
func (req *TODeliveryServiceRequest) Validate(db *sqlx.DB) []error {
	validation.NewStringRule(tovalidate.NoPeriods, "cannot contain periods")
	st := statusFromString(req.Status)

	errMap := validation.Errors{
		"changeType":      validation.Validate(req.ChangeType, validation.Required),
		"deliveryservice": validation.Validate(req.DeliveryService, validation.Required),
		"status": validation.Validate(req.Status, validation.Required,
			validation.By(func(s interface{}) error { return st.validTransition(s.(requestStatus)) })),
	}

	errs := tovalidate.ToErrors(errMap)

	var ds deliveryservice.TODeliveryService
	err := json.Unmarshal([]byte(req.DeliveryService), &ds)
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	// ensure the deliveryservice requested is valid
	e := ds.Validate(db)
	errs = append(errs, e...)

	return errs
}
