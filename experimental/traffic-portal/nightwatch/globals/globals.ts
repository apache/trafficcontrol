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

import axios, {AxiosError} from "axios";
import {NightwatchBrowser} from "nightwatch";
import type {CommonPageObject} from "nightwatch/page_objects/common";
import type {DeliveryServiceCardPageObject} from "nightwatch/page_objects/deliveryServiceCard";
import type {DeliveryServiceDetailPageObject} from "nightwatch/page_objects/deliveryServiceDetail";
import type {DeliveryServiceInvalidPageObject} from "nightwatch/page_objects/deliveryServiceInvalidationJobs";
import type {LoginPageObject} from "nightwatch/page_objects/login";
import type {ServersPageObject} from "nightwatch/page_objects/servers";
import type {UsersPageObject} from "nightwatch/page_objects/users";
import {
	CDN,
	GeoLimit, GeoProvider, LoginRequest,
	Protocol,
	RequestDeliveryService,
	ResponseCDN,
	ResponseDeliveryService
} from "trafficops-types";
import { ChangeLogsPageObject } from "nightwatch/page_objects/changeLogs";

declare module "nightwatch" {
	/**
	 * Defines the global nightwatch browser type with our types mixed in.
	 */
	export interface NightwatchCustomPageObjects {
		common: () => CommonPageObject;
		changeLogs: () => ChangeLogsPageObject;
		deliveryServiceCard: () => DeliveryServiceCardPageObject;
		deliveryServiceDetail: () => DeliveryServiceDetailPageObject;
		deliveryServiceInvalidationJobs: () => DeliveryServiceInvalidPageObject;
		login: () => LoginPageObject;
		servers: () => ServersPageObject;
		users: () => UsersPageObject;
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
	}
}

const globals = {
	adminPass: "twelve12",
	adminUser: "admin",
	afterEach: (browser: NightwatchBrowser, done: () => void): void => {
		browser.end(() => {
			done();
		});
	},
	apiVersion: "4.0",
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
			const resp = await client.post(`${apiUrl}/user/login`, JSON.stringify(loginReq));
			if(resp.headers["set-cookie"]) {
				for (const cookie of resp.headers["set-cookie"]) {
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
			console.error("Access token is not set");
			return Promise.reject();
		}
		client.defaults.headers.common = { Cookie: accessToken };

		const cdn: CDN = {
			dnssecEnabled: false, domainName: `tests${globals.uniqueString}.com`, name: `testCDN${globals.uniqueString}`
		};
		let respCDN: ResponseCDN;
		try {
			let resp = await client.post(`${apiUrl}/cdns`, JSON.stringify(cdn));
			respCDN = resp.data.response;

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
				typeId: 1,
				xmlId: `testDS${globals.uniqueString}`
			};
			resp = await client.post(`${apiUrl}/deliveryservices`, JSON.stringify(ds));
			const respDS: ResponseDeliveryService = resp.data.response[0];
			console.log(`Successfully created DS '${respDS.displayName}'`);
		} catch(e) {
			console.error((e as AxiosError).message);
			throw e;
		}
		done();
	},
	beforeEach: (browser: NightwatchBrowser, done: () => void): void => {
		browser.page.login()
			.navigate().section.loginForm
			.loginAndWait(browser.globals.adminUser, browser.globals.adminPass);
		// This ensures that we call done after loginAndWait is finished
		browser.pause(1, () => {
			done();
		});
	},
	trafficOpsURL: "https://localhost:6443",
	uniqueString: new Date().getTime().toString()
};

module.exports = globals;
