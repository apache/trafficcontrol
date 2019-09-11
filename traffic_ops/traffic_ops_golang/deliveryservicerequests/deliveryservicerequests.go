package deliveryservicerequests

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

import "encoding/json"
import "errors"
import "fmt"
import "net/http"
import "net/mail"

import "github.com/apache/trafficcontrol/lib/go-rfc"
import "github.com/apache/trafficcontrol/lib/go-tc"

import "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"

const msg = "From: %s\r\nTo:%s\r\nContent-Type: text/html\r\nSubject: Delivery Service Request for %s\r\n\r\n%s"

func Request(w http.ResponseWriter, r *http.Request) {
	inf, userErr, sysErr, errCode := api.NewInfo(r, nil, nil)
	tx := inf.Tx.Tx
	if userErr != nil || sysErr != nil {
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}
	defer inf.Close()

	var dsr tc.DeliveryServiceRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&dsr); err != nil {
		userErr = fmt.Errorf("Error parsing request: %v", err)
		errCode = http.StatusBadRequest
		api.HandleErr(w, r, tx, errCode, userErr, nil)
		return
	}

	addr, err := mail.ParseAddress(dsr.EmailTo)
	if err != nil {
		userErr = fmt.Errorf("'%s' is not a valid RFC5322 email address!", dsr.EmailTo)
		sysErr = fmt.Errorf("Parsing submitted email address: %s", err)
		errCode = http.StatusBadRequest
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	body, err := dsr.Details.Format()
	if err != nil {
		sysErr = fmt.Errorf("Failed to format email body: %v", err)
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, tx, errCode, nil, sysErr)
		return
	}

	body = fmt.Sprintf(msg, inf.Config.ConfigTO.EmailFrom, addr, dsr.Details.Customer, body)
	if ok, err := inf.SendMail(rfc.EmailAddress{*addr}, []byte(body)); !ok {
		api.HandleErr(w, r, tx, http.StatusServiceUnavailable, err, nil)
		return
	} else if err != nil {
		sysErr = fmt.Errorf("Failed to send email: %v", err)
		errCode = http.StatusBadGateway
		api.HandleErr(w, r, tx, errCode, nil, sysErr)
		return
	}

	alert := tc.Alerts{
		Alerts: []tc.Alert{
			tc.Alert{
				Level: tc.SuccessLevel.String(),
				Text: fmt.Sprintf("Delivery Service request sent to %s", dsr.EmailTo),
			},
		},
	}

	resp, err := json.Marshal(alert)
	if err != nil {
		sysErr = fmt.Errorf("Marshaling response: %v", err)
		userErr = errors.New("Email was sent, but an error occurred afterward")
		errCode = http.StatusInternalServerError
		api.HandleErr(w, r, tx, errCode, userErr, sysErr)
		return
	}

	w.Header().Set(tc.ContentType, tc.ApplicationJson)
	w.WriteHeader(http.StatusOK)
	w.Write(append(resp, '\n'))
}
