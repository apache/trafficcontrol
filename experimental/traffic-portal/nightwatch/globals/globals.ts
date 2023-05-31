/*
*
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
import { AxiosError } from "axios";
import { NightwatchBrowser } from "nightwatch";
import type { AsnDetailPageObject } from "nightwatch/page_objects/cacheGroups/asnDetail";
import type { AsnsPageObject } from "nightwatch/page_objects/cacheGroups/asnsTable";
import type { CacheGroupDetailPageObject } from "nightwatch/page_objects/cacheGroups/cacheGroupDetails";
import type { CacheGroupsPageObject } from "nightwatch/page_objects/cacheGroups/cacheGroupsTable";
import type { CoordinateDetailPageObject } from "nightwatch/page_objects/cacheGroups/coordinateDetail";
import type { CoordinatesPageObject } from "nightwatch/page_objects/cacheGroups/coordinatesTable";
import type { DivisionDetailPageObject } from "nightwatch/page_objects/cacheGroups/divisionDetail";
import type { DivisionsPageObject } from "nightwatch/page_objects/cacheGroups/divisionsTable";
import type { RegionDetailPageObject } from "nightwatch/page_objects/cacheGroups/regionDetail";
import type { RegionsPageObject } from "nightwatch/page_objects/cacheGroups/regionsTable";
import type { CDNDetailPageObject } from "nightwatch/page_objects/cdns/cdnDetail";
import type { CommonPageObject } from "nightwatch/page_objects/common";
import type { DeliveryServiceCardPageObject } from "nightwatch/page_objects/deliveryServices/deliveryServiceCard";
import type { DeliveryServiceDetailPageObject } from "nightwatch/page_objects/deliveryServices/deliveryServiceDetail";
import type { DeliveryServiceInvalidPageObject } from "nightwatch/page_objects/deliveryServices/deliveryServiceInvalidationJobs";
import type { LoginPageObject } from "nightwatch/page_objects/login";
import type { ProfileDetailPageObject } from "nightwatch/page_objects/profiles/profileDetail";
import type { ProfilePageObject } from "nightwatch/page_objects/profiles/profilesTable";
import type { PhysLocDetailPageObject } from "nightwatch/page_objects/servers/physLocDetail";
import type { PhysLocTablePageObject } from "nightwatch/page_objects/servers/physLocTable";
import type { ServersDetailPageObject } from "nightwatch/page_objects/servers/serversDetail";
import type { ServersTablePageObject } from "nightwatch/page_objects/servers/serversTable";
import type { StatusDetailPageObject } from "nightwatch/page_objects/statuses/statusDetail";
import type { StatusesTablePageObject } from "nightwatch/page_objects/statuses/statusesTable";
import type { ChangeLogsPageObject } from "nightwatch/page_objects/users/changeLogs";
import type { RoleDetailPageObject } from "nightwatch/page_objects/users/roleDetail";
import type { RolesPageObject } from "nightwatch/page_objects/users/rolesTable";
import type { TenantDetailPageObject } from "nightwatch/page_objects/users/tenantDetail";
import type { TenantsPageObject } from "nightwatch/page_objects/users/tenants";
import type { UsersPageObject } from "nightwatch/page_objects/users/users";
import {
	ResponseCDN,
	ResponseDeliveryService,
	ResponseTenant,
	TypeFromResponse,
	ResponseASN,
	ResponseDivision,
	ResponseRegion,
	ResponseCacheGroup,
	ResponsePhysicalLocation,
	ResponseCoordinate,
	ResponseStatus,
	ResponseProfile,
	ResponseServer, ResponseServerCapability, ResponseRole,
} from "trafficops-types";

import * as config from "../config.json";
import { DataClient, generateUniqueString } from "../dataClient";
import type { CapabilitiesPageObject } from "../page_objects/servers/capabilities/capabilitiesTable";
import type { CapabilityDetailsPageObject } from "../page_objects/servers/capabilities/capabilityDetails";
import type { TypeDetailPageObject } from "../page_objects/types/typeDetail";
import type { TypesPageObject } from "../page_objects/types/typesTable";

declare module "nightwatch" {
	/**
	 * Defines the global nightwatch browser type with our types mixed in.
	 */
	export interface NightwatchCustomPageObjects {
		common: () => CommonPageObject;
		cacheGroups: {
			cacheGroupDetails: () => CacheGroupDetailPageObject;
			cacheGroupsTable: () => CacheGroupsPageObject;
			coordinateDetail: () => CoordinateDetailPageObject;
			coordinatesTable: () => CoordinatesPageObject;
			divisionDetail: () => DivisionDetailPageObject;
			divisionsTable: () => DivisionsPageObject;
			regionDetail: () => RegionDetailPageObject;
			regionsTable: () => RegionsPageObject;
			asnsTable: () => AsnsPageObject;
			asnDetail: () => AsnDetailPageObject;
		};
		cdns: {
			cdnDetail: () => CDNDetailPageObject;
		};
		deliveryServices: {
			deliveryServiceCard: () => DeliveryServiceCardPageObject;
			deliveryServiceDetail: () => DeliveryServiceDetailPageObject;
			deliveryServiceInvalidationJobs: () => DeliveryServiceInvalidPageObject;
		};
		login: () => LoginPageObject;
		profiles: {
			profileTable: () => ProfilePageObject;
			profileDetail: () => ProfileDetailPageObject;
		};
		servers: {
			capabilities: {
				capabilityDetails: () => CapabilityDetailsPageObject;
				capabilitiesTable: () => CapabilitiesPageObject;
			};
			physLocDetail: () => PhysLocDetailPageObject;
			physLocTable: () => PhysLocTablePageObject;
			serversTable: () => ServersTablePageObject;
			serversDetail: () => ServersDetailPageObject;
		};
		statuses: {
			statusesTable: () => StatusesTablePageObject;
			statusDetail: () => StatusDetailPageObject;
		};
		users: {
			changeLogs: () => ChangeLogsPageObject;
			roles: () => RolesPageObject;
			roleDetail: () => RoleDetailPageObject;
			tenants: () => TenantsPageObject;
			tenantDetail: () => TenantDetailPageObject;
			users: () => UsersPageObject;
		};
		types: {
			typesTable: () => TypesPageObject;
			typeDetail: () => TypeDetailPageObject;
		};
	}

	/**
	 * Defines the additional types needed for the test environment.
	 */
	export interface NightwatchGlobals {
		adminPass: string;
		adminUser: string;
		trafficOpsURL: string;
		apiVersion: string;
		uniqueString: string;
		testData: CreatedData;
	}
}

/**
 * Contains the data created by the client before the test suite runs.
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
	physLoc: ResponsePhysicalLocation;
	region: ResponseRegion;
	role: ResponseRole;
	steeringDS: ResponseDeliveryService;
	tenant: ResponseTenant;
	type: TypeFromResponse;
	statuses: ResponseStatus;
	profile: ResponseProfile;
}

let testData = {};
let client: DataClient;
let dataCreateFailed = false;

const globals = {
	adminPass: config.adminPass,
	adminUser: config.adminUser,
	after: async (done: () => void): Promise<void> => {
		if (dataCreateFailed){
			return done();
		} else if(client.loggedIn) {
			try {
				await client.createData(generateUniqueString());
			} catch(e) {
				console.error("Idempotency test failed, err:", e);
				throw e;
			}
			console.log("Data creation is idempotent");
		} else {
			console.log("Client not logged in, skipping idempotency test");
		}
		done();
	},
	afterEach: (browser: NightwatchBrowser, done: () => void): void => {
		browser.end(() => {
			done();
		});
	},
	apiVersion: "4.0",
	before: async (done: () => void): Promise<void> => {
		client = new DataClient(globals.trafficOpsURL, globals.apiVersion, globals.adminUser, globals.adminPass);
		try {
			testData = await client.createData(globals.uniqueString);
		} catch(e) {
			dataCreateFailed = true;
			console.error("Request for", globals.trafficOpsURL, "failed:", (e as AxiosError).message);
			throw e;
		}
		done();
	},
	beforeEach: (browser: NightwatchBrowser, done: () => void): void => {
		browser.globals.testData = testData as CreatedData;
		browser.page.login()
			.navigate().section.loginForm
			.loginAndWait(browser.globals.adminUser, browser.globals.adminPass);
		// This ensures that we call done after loginAndWait is finished
		browser.pause(1, () => {
			done();
		});
	},
	retryAssertionTimeout: config.retryAssertionTimeoutMS,
	testData,
	trafficOpsURL: config.to_url,
	uniqueString: generateUniqueString(),
	waitForConditionTimeout:config.waitForConditionTimeoutMS
};

module.exports = globals;
