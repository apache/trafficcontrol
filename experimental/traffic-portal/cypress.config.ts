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
import { promises as fs } from "fs";
import * as https from "https";

import axios, { AxiosError } from "axios";
import { defineConfig } from "cypress";
import type {
	CDN,
	LoginRequest,
	ProfileType,
	RequestASN,
	RequestCacheGroup,
	RequestCoordinate,
	RequestDeliveryService,
	RequestDivision,
	RequestParameter,
	RequestPhysicalLocation,
	RequestProfile,
	RequestRegion,
	RequestRole,
	RequestServer,
	RequestServerCapability,
	RequestStatus,
	RequestSteeringTarget,
	RequestTenant,
	RequestType,
	ResponseCacheGroup,
	ResponseDeliveryService,
	ResponseDivision,
	ResponseParameter,
	ResponsePhysicalLocation,
	ResponseProfile,
	ResponseRegion,
	ResponseStatus,
	TypeFromResponse
} from "trafficops-types";

import type { CreatedData } from "./cypress/support/testing.data";

import PluginEvents = Cypress.PluginEvents;

/**
 * Creates mock data needed for E2E testing.
 *
 * Ideally this functionality would go in a different file, but for some reason
 * Cypress gets very upset if you import anything but a type from any TypeScript
 * file in this config file.
 *
 * @param toURL The URL of a Traffic Ops instance e.g. 'https://traffic.ops/'.
 * @param apiVersion The version of the API to use e.g. '4.1'.
 * @param adminUser The username of the 'admin' Role user that will be used
 * to set up testing data.
 * @param adminPass The password of the 'admin' Role user that will be used to
 * set up testing data.
 * @returns The data that was created and the unique string that was appended to
 * all names that are required to be unique.
 */
async function createData(toURL: string, apiVersion: string, adminUser: string, adminPass: string): Promise<CreatedData> {
	const apiUrl = `${toURL}/api/${apiVersion}`;
	const client = axios.create({
		httpsAgent: new https.Agent({
			rejectUnauthorized: false
		})
	});

	if (!Object.keys(client.defaults.headers.common).includes("Cookie")) {
		let accessToken = "";
		const loginReq: LoginRequest = {
			p: adminPass,
			u: adminUser
		};
		try {
			const logResp = await client.post(`${apiUrl}/user/login`, JSON.stringify(loginReq));
			if (logResp.headers["set-cookie"]) {
				for (const cookie of logResp.headers["set-cookie"]) {
					if (cookie.includes("access_token")) {
						accessToken = cookie;
						break;
					}
				}
			}
		} catch (e) {
			// eslint-disable-next-line no-console
			console.error((e as AxiosError).message);
			throw e;
		}
		if (accessToken === "") {
			const e = new Error("Access token is not set");
			// eslint-disable-next-line no-console
			console.error(e.message);
			throw e;
		}
		// eslint-disable-next-line @typescript-eslint/naming-convention
		client.defaults.headers.common = {Cookie: accessToken};
	}

	let resp = await client.get(`${apiUrl}/types`);
	const types: Array<TypeFromResponse> = resp.data.response;
	const httpType = types.find(typ => typ.name === "HTTP" && typ.useInTable === "deliveryservice");
	if (!httpType) {
		throw new Error("Unable to find `HTTP` type");
	}
	const steeringType = types.find(typ => typ.name === "STEERING" && typ.useInTable === "deliveryservice");
	if (!steeringType) {
		throw new Error("Unable to find `STEERING` type");
	}
	const steeringWeightType = types.find(typ => typ.name === "STEERING_WEIGHT" && typ.useInTable === "steering_target");
	if (!steeringWeightType) {
		throw new Error("Unable to find `STEERING_WEIGHT` type");
	}
	const cgType = types.find(typ => typ.useInTable === "cachegroup");
	if (!cgType) {
		throw new Error("Unable to find any Cache Group Types");
	}
	const edgeType = types.find(typ => typ.useInTable === "server" && typ.name === "EDGE");
	if (!edgeType) {
		throw new Error("Unable to find `EDGE` type");
	}

	const id = (new Date()).getTime().toString();
	const data = {
		uniqueString: id
	} as CreatedData;

	let url = `${apiUrl}/cdns`;
	try {
		const cdn: CDN = {
			dnssecEnabled: false, domainName: `tests${id}.com`, name: `testCDN${id}`
		};

		resp = await client.post(url, JSON.stringify(cdn));
		const respCDN = resp.data.response;
		data.cdn = respCDN;

		const ds: RequestDeliveryService = {
			active: false,
			cacheurl: null,
			cdnId: respCDN.id,
			displayName: `test DS${id}`,
			dscp: 0,
			ecsEnabled: false,
			edgeHeaderRewrite: null,
			fqPacingRate: null,
			geoLimit: 0,
			geoProvider: 0,
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
			protocol: 0,
			qstringIgnore: 0,
			rangeRequestHandling: 0,
			regionalGeoBlocking: false,
			remapText: null,
			routingName: "test",
			signed: false,
			tenantId: 1,
			typeId: httpType.id,
			xmlId: `testDS${id}`
		};
		url = `${apiUrl}/deliveryservices`;
		resp = await client.post(url, JSON.stringify(ds));
		let respDS: ResponseDeliveryService = resp.data.response[0];
		data.ds = respDS;

		ds.displayName = `test DS2${id}`;
		ds.xmlId = `testDS2${id}`;
		resp = await client.post(url, JSON.stringify(ds));
		respDS = resp.data.response[0];
		data.ds2 = respDS;

		ds.displayName = `test steering DS${id}`;
		ds.xmlId = `testSDS${id}`;
		ds.typeId = steeringType.id;
		resp = await client.post(url, JSON.stringify(ds));
		respDS = resp.data.response[0];
		data.steeringDS = respDS;

		const target: RequestSteeringTarget = {
			targetId: data.ds.id,
			typeId: steeringWeightType.id,
			value: 1
		};
		url = `${apiUrl}/steering/${data.steeringDS.id}/targets`;
		await client.post(url, JSON.stringify(target));
		target.targetId = data.ds2.id;
		await client.post(url, JSON.stringify(target));

		const tenant: RequestTenant = {
			active: true,
			name: `testT${id}`,
			parentId: 1
		};
		url = `${apiUrl}/tenants`;
		resp = await client.post(url, JSON.stringify(tenant));
		data.tenant = resp.data.response;

		const division: RequestDivision = {
			name: `testD${id}`
		};
		url = `${apiUrl}/divisions`;
		resp = await client.post(url, JSON.stringify(division));
		const respDivision: ResponseDivision = resp.data.response;
		data.division = respDivision;

		const region: RequestRegion = {
			division: respDivision.id,
			name: `testR${id}`
		};
		url = `${apiUrl}/regions`;
		resp = await client.post(url, JSON.stringify(region));
		const respRegion: ResponseRegion = resp.data.response;
		data.region = respRegion;

		const cacheGroup: RequestCacheGroup = {
			name: `test${id}`,
			shortName: `test${id}`,
			typeId: cgType.id
		};
		url = `${apiUrl}/cachegroups`;
		resp = await client.post(url, JSON.stringify(cacheGroup));
		const responseCG: ResponseCacheGroup = resp.data.response;
		data.cacheGroup = responseCG;

		const asn: RequestASN = {
			asn: +id,
			cachegroupId: responseCG.id
		};
		url = `${apiUrl}/asns`;
		resp = await client.post(url, JSON.stringify(asn));
		data.asn = resp.data.response;

		const physLoc: RequestPhysicalLocation = {
			address: "street",
			city: "city",
			comments: "someone set us up the bomb",
			email: "email@test.com",
			name: `phys${id}`,
			phone: "111-867-5309",
			poc: "me",
			regionId: respRegion.id,
			shortName: `short${id}`,
			state: "CA",
			zip: "80000"
		};
		url = `${apiUrl}/phys_locations`;
		resp = await client.post(url, JSON.stringify(physLoc));
		const respPhysLoc: ResponsePhysicalLocation = resp.data.response;
		respPhysLoc.region = respRegion.name;
		data.physLoc = respPhysLoc;

		const coordinate: RequestCoordinate = {
			latitude: 0,
			longitude: 0,
			name: `coord${id}`
		};
		url = `${apiUrl}/coordinates`;
		resp = await client.post(url, JSON.stringify(coordinate));
		data.coordinate = resp.data.response;

		const type: RequestType = {
			description: "blah",
			name: `type${id}`,
			useInTable: "server"
		};
		url = `${apiUrl}/types`;
		resp = await client.post(url, JSON.stringify(type));

		data.type = resp.data.response;
		const status: RequestStatus = {
			description: "blah",
			name: `status${id}`,
		};
		url = `${apiUrl}/statuses`;
		resp = await client.post(url, JSON.stringify(status));
		const respStatus: ResponseStatus = resp.data.response;
		data.status = respStatus;

		const profile: RequestProfile = {
			cdn: respCDN.id,
			description: "blah",
			name: `profile${id}`,
			routingDisabled: false,
			type: "ATS_PROFILE" as ProfileType,
		};
		url = `${apiUrl}/profiles`;
		resp = await client.post(url, JSON.stringify(profile));
		const respProfile: ResponseProfile = resp.data.response;
		data.profile = respProfile;

		const parameter: RequestParameter = {
			configFile: "cfg.txt",
			name: `param${id}`,
			secure: false,
			value: "10",
		};
		url = `${apiUrl}/parameters`;
		resp = await client.post(url, JSON.stringify(parameter));
		const responseParameter: ResponseParameter = resp.data.response;
		data.parameter = responseParameter;

		const server: RequestServer = {
			cachegroupId: responseCG.id,
			cdnId: respCDN.id,
			domainName: "domain.com",
			hostName: id,
			interfaces: [{
				ipAddresses: [{
					address: "192.160.1.0",
					gateway: null,
					serviceAddress: true
				}],
				maxBandwidth: 0,
				monitor: true,
				mtu: 1500,
				name: "eth0"
			}],
			physLocationId: respPhysLoc.id,
			profileNames: [respProfile.name],
			statusId: respStatus.id,
			typeId: edgeType.id

		};
		url = `${apiUrl}/servers`;
		resp = await client.post(url, JSON.stringify(server));
		data.edgeServer = resp.data.response;

		const capability: RequestServerCapability = {
			name: `test${id}`
		};
		url = `${apiUrl}/server_capabilities`;
		resp = await client.post(url, JSON.stringify(capability));
		data.capability = resp.data.response;

		const role: RequestRole = {
			description: "Has access to everything - cannot be modified or deleted",
			name: `admin${id}`,
			permissions: [
				"ALL"
			]
		};
		url = `${apiUrl}/roles`;
		resp = await client.post(url, JSON.stringify(role));
		data.role = resp.data.response;
	} catch (e) {
		const ae = e as AxiosError;
		ae.message = `Request (${ae.config?.method}) failed to ${url}`;
		ae.message += ae.response ? ` with response code ${ae.response.status}` : " with no response";
		throw ae;
	}

	return data;
}

/** The schema for the Traffic Ops connection configuration file. */
interface TOConfig {
	adminPass: string;
	adminUser: string;
	apiVersion: string;
	toURL: string;
}

export default defineConfig({
	component: {
		devServer: {
			bundler: "webpack",
			framework: "angular",
		},
		specPattern: "**/*.cy.ts"
	},
	e2e: {
		baseUrl: "http://localhost:4200",
		setupNodeEvents(on: PluginEvents) {
			on("before:run", async () => {
				const toConfig: TOConfig = JSON.parse(await fs.readFile("cypress/fixtures/to.config.json", {encoding: "utf-8"}));
				const data = await createData(toConfig.toURL, toConfig.apiVersion, toConfig.adminUser, toConfig.adminPass);
				let formattedData = JSON.stringify(data, null, "\t");
				formattedData += "\n";
				return fs.writeFile("cypress/fixtures/test.data.json", formattedData);
			});
		},
	},
	experimentalInteractiveRunEvents: true,
});
