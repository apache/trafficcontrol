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

	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/traffic_ops/traffic_ops_golang/deliveryservice"
	"github.com/jmoiron/sqlx"
)

// Validate ...
func (request *TODeliveryServiceRequest) Validate(db *sqlx.DB) []error {
	log.Debugf("Got request with %++v\n", request)
	var errs []error
	if len(request.ChangeType) == 0 {
		errs = append(errs, errors.New(`'changeType' is required`))
	}
	if len(request.Status) == 0 {
		errs = append(errs, errors.New(`'status' is required`))
	}
	if len(request.Request) == 0 {
		// TODO: validate request json has required deliveryservice fields
		errs = append(errs, errors.New(`'request' is required`))
	}

	var ds deliveryservice.TODeliveryService
	err := json.Unmarshal([]byte(request.Request), &ds)
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	// ensure the deliveryservice requested is valid
	e := ds.Validate(db)
	errs = append(errs, e...)

	return errs
}
