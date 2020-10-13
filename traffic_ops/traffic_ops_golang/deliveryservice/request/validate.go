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
	"database/sql"
	"errors"
	"fmt"

	"github.com/apache/trafficcontrol/lib/go-log"
	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/tovalidate"
	"github.com/apache/trafficcontrol/lib/go-util"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Validate ensures all required fields are present and in correct form.  Also checks request JSON is complete and valid
func validateLegacy(dsr tc.DeliveryServiceRequestV15, tx *sql.Tx) error {
	if tx == nil {
		log.Errorln("validating a legacy Delivery Service Request: nil transaction was passed")
	}

	fromStatus := tc.RequestStatusDraft
	if dsr.ID != nil && *dsr.ID > 0 {
		err := tx.QueryRow(`SELECT status FROM deliveryservice_request WHERE id=$1`, *dsr.ID).Scan(&fromStatus)

		if err != nil {
			log.Errorf("querying for dsr by ID %d: %v", *dsr.ID, err)
			return errors.New("unknown error")
		}
	}

	validTransition := func(s interface{}) error {
		if s == nil {
			return errors.New("cannot transition to nil status")
		}
		toStatus, ok := s.(*tc.RequestStatus)
		if !ok {
			return fmt.Errorf("Expected *tc.RequestStatus type,  got %T", s)
		}
		return fromStatus.ValidTransition(*toStatus)
	}

	errMap := validation.Errors{
		"changeType":      validation.Validate(dsr.ChangeType, validation.Required),
		"deliveryservice": validation.Validate(dsr.DeliveryService, validation.Required),
		"status":          validation.Validate(dsr.Status, validation.Required, validation.By(validTransition)),
	}
	errs := tovalidate.ToErrors(errMap)
	// ensure the deliveryservice requested is valid
	e := dsr.DeliveryService.Validate(tx)

	errs = append(errs, e)

	return util.JoinErrs(errs)
}
