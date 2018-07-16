package v13

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

// Other endpoints define their own alert. Why not just use tc.Alert?
type FederationAlert struct {
	Level string `json:"level"`
	Text  string `json:"text"`
}

type CDNFederationResponse struct {
	Response []CDNFederation `json:"response"`
}

type CreateCDNFederationResponse struct {
	Response CDNFederation     `json:"response"`
	Alerts   []FederationAlert `json:"alerts"`
}

type UpdateCDNFederationResponse struct {
	Reponse CDNFederation     `json:"response"`
	Alerts  []FederationAlert `json:"alerts"`
}

type DeleteCDNFederationResponse struct {
	Alerts []FederationAlert `json:"alerts"`
}

type CDNFederation struct {
	ID          *int    `json:"id" db:"id"`
	CName       *string `json:"cname" db:"cname"`
	Ttl         *int    `json:"ttl" db:"ttl"`
	Description *string `json:"description" db:"description"`

	//omitempty only works with primitive types and pointers
	*DeliveryServiceIDs `json:"deliveryService,omitempty"`

	//Extra datapoint for ease in Read function (cannot Scan after StructScan)
	Name *string `json:"-" db:"cdn_name"`
}

type DeliveryServiceIDs struct {
	ID    *int    `json:"id,omitempty" db:"ds_id"`
	XmlId *string `json:"xmlId,omitempty" db:"xml_id"`
}
