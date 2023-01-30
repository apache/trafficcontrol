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
import * as https from "https";

import axios, { AxiosError } from "axios";
import { NightwatchBrowser } from "nightwatch";
import type { CacheGroupDetailPageObject } from "nightwatch/page_objects/cacheGroups/cacheGroupDetails";
import type { CacheGroupsPageObject } from "nightwatch/page_objects/cacheGroups/cacheGroupsTable";
import type { DivisionDetailPageObject } from "nightwatch/page_objects/cacheGroups/divisionDetail";
import type { DivisionsPageObject } from "nightwatch/page_objects/cacheGroups/divisionsTable";
import type { RegionDetailPageObject } from "nightwatch/page_objects/cacheGroups/regionDetail";
import type { RegionsPageObject } from "nightwatch/page_objects/cacheGroups/regionsTable";
import type { CommonPageObject } from "nightwatch/page_objects/common";
import type { DeliveryServiceCardPageObject } from "nightwatch/page_objects/deliveryServices/deliveryServiceCard";
import type { DeliveryServiceDetailPageObject } from "nightwatch/page_objects/deliveryServices/deliveryServiceDetail";
import type { DeliveryServiceInvalidPageObject } from "nightwatch/page_objects/deliveryServices/deliveryServiceInvalidationJobs";
import type { LoginPageObject } from "nightwatch/page_objects/login";
import type { PhysLocDetailPageObject } from "nightwatch/page_objects/servers/physLocDetail";
import type { PhysLocTablePageObject } from "nightwatch/page_objects/servers/physLocTable";
import type { ServersPageObject } from "nightwatch/page_objects/servers/servers";
import type { ChangeLogsPageObject } from "nightwatch/page_objects/users/changeLogs";
import type { TenantDetailPageObject } from "nightwatch/page_objects/users/tenantDetail";
import type { TenantsPageObject } from "nightwatch/page_objects/users/tenants";
import type { UsersPageObject } from "nightwatch/page_objects/users/users";
import {
	CDN,
	GeoLimit,
	GeoProvider,
	LoginRequest,
	Protocol,
	RequestDeliveryService,
	ResponseCDN,
	ResponseDeliveryService,
	RequestTenant,
	ResponseTenant,
	TypeFromResponse,
	RequestSteeringTarget,
	ResponseDivision,
	RequestDivision,
	ResponseRegion,
	RequestRegion,
	RequestCacheGroup,
	ResponseCacheGroup,
	ResponsePhysicalLocation,
	RequestPhysicalLocation
} from "trafficops-types";

declare module "nightwatch" {
	/**
	 * Defines the global nightwatch browser type with our types mixed in.
	 */
	export interface NightwatchCustomPageObjects {
		common: () => CommonPageObject;
		cacheGroups: {
			cacheGroupDetails: () => CacheGroupDetailPageObject;
			cacheGroupsTable: () => CacheGroupsPageObject;
			divisionDetail: () => DivisionDetailPageObject;
			divisionsTable: () => DivisionsPageObject;
			regionDetail: () => RegionDetailPageObject;
			regionsTable: () => RegionsPageObject;
		};
		deliveryServices: {
			deliveryServiceCard: () => DeliveryServiceCardPageObject;
			deliveryServiceDetail: () => DeliveryServiceDetailPageObject;
			deliveryServiceInvalidationJobs: () => DeliveryServiceInvalidPageObject;
		};
		login: () => LoginPageObject;
		servers: {
			physLocDetail: () => PhysLocDetailPageObject;
			physLocTable: () => PhysLocTablePageObject;
			servers: () => ServersPageObject;
		};
		users: {
			changeLogs: () => ChangeLogsPageObject;
			tenants: () => TenantsPageObject;
			tenantDetail: () => TenantDetailPageObject;
			users: () => UsersPageObject;
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
	cacheGroup: ResponseCacheGroup;
	cdn: ResponseCDN;
	division: ResponseDivision;
	ds: ResponseDeliveryService;
	ds2: ResponseDeliveryService;
	physLoc: ResponsePhysicalLocation;
	region: ResponseRegion;
	steeringDS: ResponseDeliveryService;
	tenant: ResponseTenant;
}

const testData = {};

const globals = {
	adminPass: "twelve12",
	adminUser: "admin",
	afterEach: (browser: NightwatchBrowser, done: () => void): void => {
		browser.end(() => {
			done();
		});
	},
	apiVersion: "3.1",
	before: async (done: () => void): Promise<void> => {
		const apiUrl = `${globals.trafficOpsURL}/api/${globals.apiVersion}`;
		const client = axios.create({
			httpsAgent: new https.Agent({
				rejectUnauthorized: false
			})
		});
		let accessToken = "";
		const loginReq: LoginRequest = {
			p: globals.adminPass,
			u: globals.adminUser
		};
		try {
			const logResp = await client.post(`${apiUrl}/user/login`, JSON.stringify(loginReq));
			if(logResp.headers["set-cookie"]) {
				for (const cookie of logResp.headers["set-cookie"]) {
					if(cookie.indexOf("access_token") > -1) {
						accessToken = cookie;
						break;
					}
				}
			}
		} catch (e) {
			console.error((e as AxiosError).message);
			throw e;
		}
		if(accessToken === "") {
			const e = new Error("Access token is not set");
			console.error(e.message);
			throw e;
		}
		client.defaults.headers.common = { Cookie: accessToken };

		const cdn: CDN = {
			dnssecEnabled: false, domainName: `tests${globals.uniqueString}.com`, name: `testCDN${globals.uniqueString}`
		};
		let respCDN: ResponseCDN;

		let resp = await client.get(`${apiUrl}/types`);
		const types: Array<TypeFromResponse> = resp.data.response;
		const httpType = types.find(typ => typ.name === "HTTP" && typ.useInTable === "deliveryservice");
		if(httpType === undefined) {
			throw new Error("Unable to find `HTTP` type");
		}
		const steeringType = types.find(typ => typ.name === "STEERING" && typ.useInTable === "deliveryservice");
		if(steeringType === undefined) {
			throw new Error("Unable to find `STEERING` type");
		}
		const steeringWeightType = types.find(typ => typ.name === "STEERING_WEIGHT" && typ.useInTable === "steering_target");
		if(steeringWeightType === undefined) {
			throw new Error("Unable to find `STEERING_WEIGHT` type");
		}
		const cgType = types.find(typ => typ.useInTable === "cachegroup");
		if (!cgType) {
			throw new Error("Unable to find any Cache Group Types");
		}

		try {
			const data = testData as CreatedData;
			resp = await client.post(`${apiUrl}/cdns`, JSON.stringify(cdn));
			respCDN = resp.data.response;
			console.log(`Successfully created CDN ${respCDN.name}`);
			data.cdn = respCDN;

			const ds: RequestDeliveryService = {
				active: false,
				cacheurl: null,
				cdnId: respCDN.id,
				displayName: `test DS${globals.uniqueString}`,
				dscp: 0,
				ecsEnabled: false,
				edgeHeaderRewrite: null,
				fqPacingRate: null,
				geoLimit: GeoLimit.NONE,
				geoProvider: GeoProvider.MAX_MIND,
				httpBypassFqdn: null,
				infoUrl: null,
				initialDispersion: 1,
				ipv6RoutingEnabled: false,
				logsEnabled: false,
				maxOriginConnections: 0,
				maxRequestHeaderBytes: 0,
				midHeaderRewrite: null,
				missLat: 0,
				missLong: 0,
				multiSiteOrigin: false,
				orgServerFqdn: "http://test.com",
				profileId: 1,
				protocol: Protocol.HTTP,
				qstringIgnore: 0,
				rangeRequestHandling: 0,
				regionalGeoBlocking: false,
				remapText: null,
				routingName: "test",
				signed: false,
				tenantId: 1,
				typeId: httpType.id,
				xmlId: `testDS${globals.uniqueString}`
			};
			resp = await client.post(`${apiUrl}/deliveryservices`, JSON.stringify(ds));
			let respDS: ResponseDeliveryService = resp.data.response[0];
			console.log(`Successfully created DS '${respDS.displayName}'`);
			data.ds = respDS;

			ds.displayName = `test DS2${globals.uniqueString}`;
			ds.xmlId = `testDS2${globals.uniqueString}`;
			resp = await client.post(`${apiUrl}/deliveryservices`, JSON.stringify(ds));
			respDS = resp.data.response[0];
			console.log(`Successfully created DS '${respDS.displayName}'`);
			data.ds2 = respDS;

			ds.displayName = `test steering DS${globals.uniqueString}`;
			ds.xmlId = `testSDS${globals.uniqueString}`;
			ds.typeId = steeringType.id;
			resp = await client.post(`${apiUrl}/deliveryservices`, JSON.stringify(ds));
			respDS = resp.data.response[0];
			console.log(`Successfully created DS '${respDS.displayName}'`);
			data.steeringDS = respDS;

			const target: RequestSteeringTarget = {
				targetId: data.ds.id,
				typeId: steeringWeightType.id,
				value: 1
			};
			await client.post(`${apiUrl}/steering/${data.steeringDS.id}/targets`, JSON.stringify(target));
			target.targetId = data.ds2.id;
			await client.post(`${apiUrl}/steering/${data.steeringDS.id}/targets`, JSON.stringify(target));
			console.log(`Created steering targets for ${data.steeringDS.displayName}`);

			const tenant: RequestTenant = {
				active: true,
				name: `testT${globals.uniqueString}`,
				parentId: 1
			};
			resp = await client.post(`${apiUrl}/tenants`, JSON.stringify(tenant));
			const respTenant: ResponseTenant = resp.data.response;
			console.log(`Successfully created Tenant ${respTenant.name}`);
			data.tenant = respTenant;

			const division: RequestDivision = {
				name: `testD${globals.uniqueString}`
			};
			resp = await client.post(`${apiUrl}/divisions`, JSON.stringify(division));
			const respDivision: ResponseDivision = resp.data.response;
			console.log(`Successfully created Division ${respDivision.name}`);
			data.division = respDivision;

			const region: RequestRegion = {
				division: 1,
				name: `testR${globals.uniqueString}`
			};
			resp = await client.post(`${apiUrl}/regions`, JSON.stringify(region));
			const respRegion: ResponseRegion = resp.data.response;
			console.log(`Successfully created Region ${respRegion.name}`);
			data.region = respRegion;

			const cacheGroup: RequestCacheGroup = {
				name: `test${globals.uniqueString}`,
				shortName: `test${globals.uniqueString}`,
				typeId: cgType.id
			};
			resp = await client.post(`${apiUrl}/cachegroups`, JSON.stringify(cacheGroup));
			const responseCG: ResponseCacheGroup = resp.data.response;
			console.log("Successfully created Cache Group:", responseCG);
			data.cacheGroup = responseCG;

			const physLoc: RequestPhysicalLocation = {
				address: "street",
				city: "city",
				comments: "someone set us up the bomb",
				email: "email@test.com",
				name: `phys${globals.uniqueString}`,
				phone: "111-867-5309",
				poc: "me",
				regionId: respRegion.id,
				shortName: `short${globals.uniqueString}`,
				state: "CA",
				zip: "80000"
			};
			resp = await client.post(`${apiUrl}/phys_locations`, JSON.stringify(physLoc));
			const respPhysLoc: ResponsePhysicalLocation = resp.data.response;
			respPhysLoc.region = respRegion.name;
			console.log(`Successfully created Phys Loc ${respPhysLoc.name}`);
			data.physLoc = respPhysLoc;
		} catch(e) {
			console.error((e as AxiosError).message);
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
	testData,
	trafficOpsURL: "https://localhost:6443",
	uniqueString: new Date().getTime().toString()
};

module.exports = globals;
