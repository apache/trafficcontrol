/*
   Copyright 2016 Comcast Cable Communications Management, LLC

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

package client

const apiBase = "/api/1.2"
const dsPath = "/deliveryservices"

func deliveryServicesEp() string {
	return apiBase + dsPath + ".json"
}

func deliveryServiceBaseEp(id string) string {
	return apiBase + dsPath + "/" + id
}

func deliveryServiceEp(id string) string {
	return deliveryServiceBaseEp(id) + ".json"
}

func deliveryServiceStateEp(id string) string {
	return deliveryServiceBaseEp(id) + "/state.json"
}

func deliveryServiceHealthEp(id string) string {
	return deliveryServiceBaseEp(id) + "/health.json"
}

func deliveryServiceCapacityEp(id string) string {
	return deliveryServiceBaseEp(id) + "/capacity.json"
}

func deliveryServiceRoutingEp(id string) string {
	return deliveryServiceBaseEp(id) + "/routing.json"
}

func deliveryServiceServerEp(page, limit string) string {
	return apiBase + "/deliveryserviceserver.json?page=" + page + "&limit=" + limit
}
