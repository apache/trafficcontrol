/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import * as https from "https";

import axios from "axios";
import {AugmentedBrowser} from "nightwatch/globals/globals";
import {
	CDN,
	GeoLimit, GeoProvider, LoginRequest,
	Protocol,
	RequestDeliveryService,
	ResponseCDN,
	ResponseDeliveryService
} from "trafficops-types";

/**
 * Creates data necessary for the e2e tests, currently just a CDN, and a DS
 *
 * @param augBrowser The browser object
 * @returns Created Delivery Service
 */
async function createData(augBrowser: AugmentedBrowser): Promise<ResponseDeliveryService> {
	const apiUrl = `${augBrowser.globals.trafficOpsURL}/api/${augBrowser.globals.apiVersion}`;
	const client = axios.create({
		httpsAgent: new https.Agent({
			rejectUnauthorized: false
		})
	});
	let accessToken = "";
	const loginReq: LoginRequest = {
		p: augBrowser.globals.adminPass,
		u: augBrowser.globals.adminUser
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
		return Promise.reject(e);
	}
	if(accessToken === "") {
		return Promise.reject("Access token is not set");
	}
	client.defaults.headers.common = { Cookie: accessToken };

	const cdn: CDN = {
		dnssecEnabled: false, domainName: "tests.com", name: "testCDN"
	};
	let respCDN: ResponseCDN;
	try {
		let resp = await client.post(`${apiUrl}/cdns`, JSON.stringify(cdn));
		respCDN = resp.data.response;

		const ds: RequestDeliveryService = {
			active: false,
			cacheurl: null,
			cdnId: respCDN.id,
			displayName: "test DS",
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
			xmlId: "testDS"
		};
		resp = await client.post(`${apiUrl}/deliveryservices`, JSON.stringify(ds));
		const respDS: ResponseDeliveryService = resp.data.response[0];

		return Promise.resolve(respDS);
	} catch(e) {
		return Promise.reject(e);
	}
}

describe("Bootstrap Spec", () => {
	it("Create Data", async (): Promise<void> => {
		const augBrowser = browser as AugmentedBrowser;
		await createData(augBrowser).then(ds => {
			console.log(`Successfully created DS '${ds.displayName}'`);
		}).catch(err => {
			augBrowser.assert.fail(err, null, "Exception occurred while creating data");
		}).finally(() => {
			augBrowser.end();
		});
	});
});
