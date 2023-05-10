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

import axios, { AxiosError, AxiosInstance } from "axios";
import { CreatedData } from "nightwatch/globals/globals";
import {
	CDN,
	GeoLimit,
	GeoProvider,
	LoginRequest,
	ProfileType,
	Protocol,
	RequestASN,
	RequestCacheGroup,
	RequestCoordinate,
	RequestDeliveryService,
	RequestDivision,
	RequestPhysicalLocation,
	RequestProfile,
	RequestRegion,
	RequestServer,
	RequestStatus,
	RequestSteeringTarget,
	RequestTenant,
	RequestType,
	ResponseCacheGroup,
	ResponseDeliveryService,
	ResponseDivision,
	ResponsePhysicalLocation,
	ResponseProfile,
	ResponseRegion, ResponseStatus,
	TypeFromResponse
} from "trafficops-types";

/**
 * Generates a unique string used for tests, uses the current epoch time.
 */
export function generateUniqueString(): string {
	return new Date().getTime().toString();
}

/**
 * Defines the class used to create test data for the E2E environment
 */
export class DataClient {
	private readonly toURL: string;
	private readonly apiVersion: string;
	private readonly adminUser: string;
	private readonly adminPass: string;
	/** Tracks if the client has logged in */
	public loggedIn = false;
	/** Client used to talk to the TO API */
	private readonly client: AxiosInstance;

	public constructor(toURL: string, apiVersion: string, adminUser: string, adminPass: string) {
		this.toURL = toURL;
		this.apiVersion = apiVersion;
		this.adminUser = adminUser;
		this.adminPass = adminPass;

		this.client = axios.create({
			httpsAgent: new https.Agent({
				rejectUnauthorized: false
			})
		});
	}

	/**
	 * Creates data needed for the E2E tests
	 *
	 * @param id ID added to various fields to ensure that creation occurs regardless of environment
	 */
	public async createData(id: string): Promise<CreatedData> {
		const apiUrl = `${this.toURL}/api/${this.apiVersion}`;
		if (Object.keys(this.client.defaults.headers.common).indexOf("Cookie") === -1) {
			this.loggedIn = false;
			let accessToken = "";
			const loginReq: LoginRequest = {
				p: this.adminPass,
				u: this.adminUser
			};
			try {
				const logResp = await this.client.post(`${apiUrl}/user/login`, JSON.stringify(loginReq));
				if (logResp.headers["set-cookie"]) {
					for (const cookie of logResp.headers["set-cookie"]) {
						if (cookie.indexOf("access_token") > -1) {
							accessToken = cookie;
							break;
						}
					}
				}
			} catch (e) {
				console.error((e as AxiosError).message);
				throw e;
			}
			if (accessToken === "") {
				const e = new Error("Access token is not set");
				console.error(e.message);
				throw e;
			}
			this.loggedIn = true;
			this.client.defaults.headers.common = {Cookie: accessToken};
		}

		const cdn: CDN = {
			dnssecEnabled: false, domainName: `tests${id}.com`, name: `testCDN${id}`
		};

		let resp = await this.client.get(`${apiUrl}/types`);
		const types: Array<TypeFromResponse> = resp.data.response;
		const httpType = types.find(typ => typ.name === "HTTP" && typ.useInTable === "deliveryservice");
		if (httpType === undefined) {
			throw new Error("Unable to find `HTTP` type");
		}
		const steeringType = types.find(typ => typ.name === "STEERING" && typ.useInTable === "deliveryservice");
		if (steeringType === undefined) {
			throw new Error("Unable to find `STEERING` type");
		}
		const steeringWeightType = types.find(typ => typ.name === "STEERING_WEIGHT" && typ.useInTable === "steering_target");
		if (steeringWeightType === undefined) {
			throw new Error("Unable to find `STEERING_WEIGHT` type");
		}
		const cgType = types.find(typ => typ.useInTable === "cachegroup");
		if (!cgType) {
			throw new Error("Unable to find any Cache Group Types");
		}
		const edgeType = types.find(typ => typ.useInTable === "server" && typ.name === "EDGE");
		if (edgeType === undefined) {
			throw new Error("Unable to find `EDGE` type");
		}

		const data = {} as CreatedData;
		let url = `${apiUrl}/cdns`;
		try {
			resp = await this.client.post(url, JSON.stringify(cdn));
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
				xmlId: `testDS${id}`
			};
			url = `${apiUrl}/deliveryservices`;
			resp = await this.client.post(url, JSON.stringify(ds));
			let respDS: ResponseDeliveryService = resp.data.response[0];
			data.ds = respDS;

			ds.displayName = `test DS2${id}`;
			ds.xmlId = `testDS2${id}`;
			resp = await this.client.post(url, JSON.stringify(ds));
			respDS = resp.data.response[0];
			data.ds2 = respDS;

			ds.displayName = `test steering DS${id}`;
			ds.xmlId = `testSDS${id}`;
			ds.typeId = steeringType.id;
			resp = await this.client.post(url, JSON.stringify(ds));
			respDS = resp.data.response[0];
			data.steeringDS = respDS;

			const target: RequestSteeringTarget = {
				targetId: data.ds.id,
				typeId: steeringWeightType.id,
				value: 1
			};
			url = `${apiUrl}/steering/${data.steeringDS.id}/targets`;
			await this.client.post(url, JSON.stringify(target));
			target.targetId = data.ds2.id;
			await this.client.post(url, JSON.stringify(target));

			const tenant: RequestTenant = {
				active: true,
				name: `testT${id}`,
				parentId: 1
			};
			url = `${apiUrl}/tenants`;
			resp = await this.client.post(url, JSON.stringify(tenant));
			data.tenant = resp.data.response;

			const division: RequestDivision = {
				name: `testD${id}`
			};
			url = `${apiUrl}/divisions`;
			resp = await this.client.post(url, JSON.stringify(division));
			const respDivision: ResponseDivision = resp.data.response;
			data.division = respDivision;

			const region: RequestRegion = {
				division: respDivision.id,
				name: `testR${id}`
			};
			url = `${apiUrl}/regions`;
			resp = await this.client.post(url, JSON.stringify(region));
			const respRegion: ResponseRegion = resp.data.response;
			data.region = respRegion;

			const cacheGroup: RequestCacheGroup = {
				name: `test${id}`,
				shortName: `test${id}`,
				typeId: cgType.id
			};
			url = `${apiUrl}/cachegroups`;
			resp = await this.client.post(url, JSON.stringify(cacheGroup));
			const responseCG: ResponseCacheGroup = resp.data.response;
			data.cacheGroup = responseCG;

			const asn: RequestASN = {
				asn: +id,
				cachegroupId: responseCG.id
			};
			url = `${apiUrl}/asns`;
			resp = await this.client.post(url, JSON.stringify(asn));
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
			resp = await this.client.post(url, JSON.stringify(physLoc));
			const respPhysLoc: ResponsePhysicalLocation = resp.data.response;
			respPhysLoc.region = respRegion.name;
			data.physLoc = respPhysLoc;

			const coordinate: RequestCoordinate = {
				latitude: 0,
				longitude: 0,
				name: `coord${id}`
			};
			url = `${apiUrl}/coordinates`;
			resp = await this.client.post(url, JSON.stringify(coordinate));
			data.coordinate = resp.data.response;

			const type: RequestType = {
				description: "blah",
				name: `type${id}`,
				useInTable: "server"
			};
			url = `${apiUrl}/types`;
			resp = await this.client.post(url, JSON.stringify(type));

			data.type = resp.data.response;
			const status: RequestStatus = {
				description: "blah",
				name: `status${id}`,
			};
			url = `${apiUrl}/statuses`;
			resp = await this.client.post(url, JSON.stringify(status));
			const respStatus: ResponseStatus = resp.data.response;
			data.statuses = respStatus;

			const profile: RequestProfile = {
				cdn: respCDN.id,
				description: "blah",
				name: `profile${id}`,
				routingDisabled: false,
				type: ProfileType.ATS_PROFILE,
			};
			url = `${apiUrl}/profiles`;
			resp = await this.client.post(url, JSON.stringify(profile));
			const respProfile: ResponseProfile = resp.data.response;
			data.profile = respProfile;

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
			resp = await this.client.post(url, JSON.stringify(server));
			data.edgeServer = resp.data.response;
		} catch (e) {
			const ae = e as AxiosError;
			ae.message = `Request (${ae.config.method}) failed to ${url}`;
			ae.message += ae.response ? ` with response code ${ae.response.status}` : " with no response";
			throw ae;
		}


		return data;
	}
}
