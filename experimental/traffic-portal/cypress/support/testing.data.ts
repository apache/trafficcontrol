/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

import type {
	ResponseASN,
	ResponseCDN,
	ResponseCacheGroup,
	ResponseCoordinate,
	ResponseDeliveryService,
	ResponseDivision,
	ResponseParameter,
	ResponsePhysicalLocation,
	ResponseProfile,
	ResponseRegion,
	ResponseRole,
	ResponseServer,
	ResponseServerCapability,
	ResponseStatus,
	ResponseTenant,
	TypeFromResponse
} from "trafficops-types";

/**
 * The mock data created for use in E2E testing.
 */
export interface CreatedData {
	asn: ResponseASN;
	cacheGroup: ResponseCacheGroup;
	capability: ResponseServerCapability;
	cdn: ResponseCDN;
	coordinate: ResponseCoordinate;
	division: ResponseDivision;
	ds: ResponseDeliveryService;
	ds2: ResponseDeliveryService;
	edgeServer: ResponseServer;
	parameter: ResponseParameter;
	physLoc: ResponsePhysicalLocation;
	region: ResponseRegion;
	role: ResponseRole;
	steeringDS: ResponseDeliveryService;
	tenant: ResponseTenant;
	type: TypeFromResponse;
	status: ResponseStatus;
	profile: ResponseProfile;
	uniqueString: string;
}

/**
 * Contains data used by the E2E tests to authenticate with Traffic Ops.
 */
export interface LoginData {
	username: string;
	password: string;
}
