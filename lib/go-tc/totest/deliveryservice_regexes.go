package totest

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
	"fmt"
	"testing"

	"github.com/apache/trafficcontrol/v8/lib/go-tc"
	"github.com/apache/trafficcontrol/v8/lib/go-util/assert"
	toclient "github.com/apache/trafficcontrol/v8/traffic_ops/v5-client"
)

func CreateTestDeliveryServicesRegexes(t *testing.T, cl *toclient.Session, td TrafficControl) {
	for _, dsRegex := range td.DeliveryServicesRegexes {
		dsID := GetDeliveryServiceId(t, cl, dsRegex.DSName)()
		typeId := GetTypeId(t, cl, dsRegex.TypeName)
		dsRegexPost := tc.DeliveryServiceRegexPost{
			Type:      typeId,
			SetNumber: dsRegex.SetNumber,
			Pattern:   dsRegex.Pattern,
		}
		alerts, _, err := cl.PostDeliveryServiceRegexesByDSID(dsID, dsRegexPost, toclient.RequestOptions{})
		assert.NoError(t, err, "Could not create Delivery Service Regex: %v - alerts: %+v", err, alerts)
	}
}

func DeleteTestDeliveryServicesRegexes(t *testing.T, cl *toclient.Session, td TrafficControl, db *sql.DB) {
	for _, regex := range td.DeliveryServicesRegexes {
		err := execSQL(db, fmt.Sprintf("DELETE FROM deliveryservice_regex WHERE deliveryservice = '%v' and regex ='%v';", regex.DSID, regex.ID))
		assert.RequireNoError(t, err, "Unable to delete deliveryservice_regex by regex %v and ds %v: %v", regex.ID, regex.DSID, err)

		err = execSQL(db, fmt.Sprintf("DELETE FROM regex WHERE Id = '%v';", regex.ID))
		assert.RequireNoError(t, err, "Unable to delete regex %v: %v", regex.ID, err)
	}
}
