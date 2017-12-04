package tc

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

// StatsSummaryResponse ...
type StatsSummaryResponse struct {
	Response []StatsSummary `json:"response"`
}

// StatsSummary ...
type StatsSummary struct {
	CDNName         string `json:"cdnName"`
	DeliveryService string `json:"deliveryServiceName"`
	StatName        string `json:"statName"`
	StatValue       string `json:"statValue"`
	SummaryTime     string `json:"summaryTime"`
	StatDate        string `json:"statDate"`
}

// LastUpdated ...
type LastUpdated struct {
	Version  string `json:"version"`
	Response struct {
		SummaryTime string `json:"summaryTime"`
	} `json:"response"`
}
