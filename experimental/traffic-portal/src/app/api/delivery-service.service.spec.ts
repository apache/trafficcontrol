/**
 * @license Apache-2.0
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
import { HttpClientTestingModule, HttpTestingController } from "@angular/common/http/testing";
import { TestBed } from "@angular/core/testing";
import { DSStats, DSStatsMetricType, ResponseDeliveryServiceSSLKey } from "trafficops-types";

import { constructDataSetFromResponse, DeliveryServiceService } from "./delivery-service.service";

const now = new Date();
const twoSecondsAgo = new Date(now.getTime() - 2000);

/** A dummy DS for testing */
export const testDS = {
	active: false,
	anonymousBlockingEnabled: false,
	cacheurl: null,
	ccrDnsTtl: null,
	cdnId: 1,
	cdnName: "cdn",
	checkPath: null,
	consistentHashQueryParams: null,
	consistentHashRegex: null,
	deepCachingType: "NEVER" as const,
	displayName: "Test Quest",
	dnsBypassCname: null,
	dnsBypassIp: null,
	dnsBypassIp6: null,
	dnsBypassTtl: null,
	dscp: 2,
	ecsEnabled: false,
	edgeHeaderRewrite: null,
	exampleURLs: [],
	firstHeaderRewrite: null,
	fqPacingRate: 0,
	geoLimit: 0,
	geoLimitCountries: null,
	geoLimitRedirectURL: null,
	geoProvider: 1,
	globalMaxMbps: null,
	globalMaxTps: null,
	httpBypassFqdn: null,
	id: 1,
	infoUrl: null,
	initialDispersion: 2,
	innerHeaderRewrite: null,
	ipv6RoutingEnabled: true,
	lastHeaderRewrite: null,
	lastUpdated: new Date(),
	logsEnabled: false,
	longDesc: null,
	matchList: [],
	maxDnsAnswers: null,
	maxOriginConnections: 0,
	maxRequestHeaderBytes: 0,
	midHeaderRewrite: null,
	missLat: 0,
	missLong: 0,
	multiSiteOrigin: false,
	orgServerFqdn: "",
	originShield: null,
	profileDescription: null,
	profileId: null,
	profileName: null,
	protocol: null,
	qstringIgnore: null,
	rangeRequestHandling: null,
	rangeSliceBlockSize: null,
	regexRemap: null,
	regionalGeoBlocking: false,
	remapText: null,
	routingName: "cdn",
	serviceCategory: null,
	signed: false,
	signingAlgorithm: null,
	sslKeyVersion: null,
	tenant: "root",
	tenantId: 1,
	tlsVersions: null,
	topology: null,
	trRequestHeaders: null,
	trResponseHeaders: null,
	type: "HTTP",
	typeId: 1,
	xmlId: "testquest",
};

/** A dummy DS SSL Key for testing */
export const testDSSSLKeys: ResponseDeliveryServiceSSLKey = {
	cdn: testDS.cdnName,
	certificate: {crt: "", csr: "", key: ""},
	deliveryservice: testDS.xmlId,
	expiration: new Date(),
	version: ""
};

/**
 * Generates a basic set of DSStats data for use in tests.
 *
 * @param label The label for the data series i.e. the type of data it contains.
 * @returns A sample data set for testing purposes.
 */
function getDSStats(label: DSStatsMetricType): DSStats {
	return {
		series: {
			columns: ["time", "mean"],
			count: 0,
			name: `${label}.ds.1min`,
			values: [
				[
					twoSecondsAgo,
					null
				],
				[
					now,
					0
				]
			],
		},
		summary: {
			average: 1,
			count: 2,
			fifthPercentile: 3,
			max: 4,
			min: 5,
			ninetyEightPercentile: 6,
			ninetyFifthPercentile: 7
		}
	};
}

describe("Delivery Service API utilities", () => {
	it("throws an error when attempting to convert a data set with no series", () => {
		expect(() => constructDataSetFromResponse({})).toThrowError("invalid data set response");
	});

	it("constructs a data set from DSStats with no summary", () => {
		const output = constructDataSetFromResponse({
			series: {
				columns: ["time", "mean"],
				count: 2,
				name: "kbps.ds.1min",
				values: [
					[twoSecondsAgo, null],
					[now, 2]
				]
			}
		});

		expect(output).toEqual({
			dataSet: {
				data: [
					{
						t: now,
						y: (2).toFixed(3)
					}
				],
				label: "kbps"
			},
			fifthPercentile: -1,
			max: -1,
			mean: -1,
			min: -1,
			ninetyEighthPercentile: -1,
			ninetyFifthPercentile: -1
		});
	});

	it("constructs a data set from DSStats including a summary", () => {
		const output = constructDataSetFromResponse({
			series: {
				columns: ["time", "mean"],
				count: 2,
				name: "kbps.ds.1min",
				values: [
					[twoSecondsAgo, null],
					[now, 2]
				]
			},
			summary: {
				average: 1,
				count: 2,
				fifthPercentile: 3,
				max: 4,
				min: 5,
				ninetyEightPercentile: 6,
				ninetyFifthPercentile: 7
			}
		});

		expect(output).toEqual({
			dataSet: {
				data: [
					{
						t: now,
						y: (2).toFixed(3)
					}
				],
				label: "kbps"
			},
			fifthPercentile: 3,
			max: 4,
			mean: 1,
			min: 5,
			ninetyEighthPercentile: 6,
			ninetyFifthPercentile: 7
		});
	});
});

describe("DeliveryServiceService", () => {
	let service: DeliveryServiceService;
	let httpTestingController: HttpTestingController;

	beforeEach(() => {
		TestBed.configureTestingModule({
			imports: [HttpClientTestingModule],
			providers: [
				DeliveryServiceService,
			]
		});
		service = TestBed.inject(DeliveryServiceService);
		httpTestingController = TestBed.inject(HttpTestingController);
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	describe("Stats methods", () => {
		const response: DSStats = {
			series: {
				columns: ["time", "mean"],
				count: 0,
				name: "kbps.ds.1min",
				values: [
					[
						twoSecondsAgo,
						null
					],
					[
						now,
						0
					]
				],
			},
			summary: {
				average: 0,
				count: 0,
				fifthPercentile: 0,
				max: 0,
				min: 0,
				ninetyEightPercentile: 0,
				ninetyFifthPercentile: 0
			}
		};

		const health = {
			response: {
				cacheGroups: [{
					name: "name",
					offline: 1,
					online: 99,
				}],
				totalOffline: 1,
				totalOnline: 99
			}
		};

		const capacity = {
			response: {
				availablePercent: 80,
				maintenancePercent: 5,
				unavailablePercent: 5,
				utilizedPercent: 10,
			}
		};

		const interval = "60s";

		it("sends requests for KBPS stats", async () => {
			const responseP = service.getDSKBPS(testDS, twoSecondsAgo, now, interval, false);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/deliveryservice_stats`);
			expect(req.request.params.keys().length).toBe(6);
			expect(req.request.params.get("deliveryServiceName")).toBe(testDS.xmlId);
			expect(req.request.params.get("endDate")).toBe(now.toISOString());
			expect(req.request.params.get("startDate")).toBe(twoSecondsAgo.toISOString());
			expect(req.request.params.get("metricType")).toBe("kbps");
			expect(req.request.params.get("interval")).toBe(interval);
			expect(req.request.params.get("serverType")).toBe("edge");
			expect(req.request.method).toBe("GET");
			req.flush({response});
			await expectAsync(responseP).toBeResolvedTo(response);
		});

		it("sends requests for KBPS stats, returning only the data series", async () => {
			const responseP = service.getDSKBPS(testDS.xmlId, twoSecondsAgo, now, interval, true, true);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/deliveryservice_stats`);
			expect(req.request.params.keys().length).toBe(7);
			expect(req.request.params.get("deliveryServiceName")).toBe(testDS.xmlId);
			expect(req.request.params.get("endDate")).toBe(now.toISOString());
			expect(req.request.params.get("startDate")).toBe(twoSecondsAgo.toISOString());
			expect(req.request.params.get("metricType")).toBe("kbps");
			expect(req.request.params.get("exclude")).toBe("summary");
			expect(req.request.params.get("interval")).toBe(interval);
			expect(req.request.params.get("serverType")).toBe("mid");
			expect(req.request.method).toBe("GET");
			req.flush({response});
			await expectAsync(responseP).toBeResolvedTo([{t: now, y: (0).toFixed(3)}]);
		});

		it("throws an error when only data is requested, but TO omits it", async () => {
			const responseP = service.getDSKBPS(testDS.xmlId, twoSecondsAgo, now, interval, true, true);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/deliveryservice_stats`);
			req.flush({response: {}});
			await expectAsync(responseP).toBeRejectedWithError("no data series found");
		});

		it("sends requests for TPS stats", async () => {
			let responseP = service.getDSTPS(testDS, twoSecondsAgo, now, interval, false);
			let req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/deliveryservice_stats`);
			expect(req.request.params.keys().length).toBe(6);
			expect(req.request.params.get("deliveryServiceName")).toBe(testDS.xmlId);
			expect(req.request.params.get("endDate")).toBe(now.toISOString());
			expect(req.request.params.get("startDate")).toBe(twoSecondsAgo.toISOString());
			expect(req.request.params.get("metricType")).toBe("tps_total");
			expect(req.request.params.get("interval")).toBe(interval);
			expect(req.request.params.get("serverType")).toBe("edge");
			expect(req.request.method).toBe("GET");
			req.flush({response});
			await expectAsync(responseP).toBeResolvedTo(response);

			responseP = service.getDSTPS(testDS.xmlId, twoSecondsAgo, now, interval, true);
			req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/deliveryservice_stats`);
			expect(req.request.params.keys().length).toBe(6);
			expect(req.request.params.get("deliveryServiceName")).toBe(testDS.xmlId);
			expect(req.request.params.get("endDate")).toBe(now.toISOString());
			expect(req.request.params.get("startDate")).toBe(twoSecondsAgo.toISOString());
			expect(req.request.params.get("metricType")).toBe("tps_total");
			expect(req.request.params.get("interval")).toBe(interval);
			expect(req.request.params.get("serverType")).toBe("mid");
			expect(req.request.method).toBe("GET");
			req.flush({response});
			await expectAsync(responseP).toBeResolvedTo(response);
		});

		it("gets all TPS data for a DS", async () => {
			const totalResponse = getDSStats(DSStatsMetricType.TPS_TOTAL);
			const successResponse = getDSStats(DSStatsMetricType.TPS_2XX);
			const redirectionResponse = getDSStats(DSStatsMetricType.TPS_3XX);
			const clientErrorResponse = getDSStats(DSStatsMetricType.TPS_4XX);
			const serverErrorResponse = getDSStats(DSStatsMetricType.TPS_5XX);
			const responseP = service.getAllDSTPSData(testDS, twoSecondsAgo, now, interval);

			const requests = httpTestingController.match(r => r.url === `/api/${service.apiVersion}/deliveryservice_stats`);
			if (requests.length !== 5) {
				return fail(`expected getting all TPS stats to make 5 requests, found: ${requests.length}`);
			}

			for (const req of requests) {
				expect(req.request.params.keys().length).toBe(6);
				expect(req.request.params.get("serverType")).toBe("edge");
				expect(req.request.params.get("interval")).toBe(interval);
				expect(req.request.params.get("deliveryServiceName")).toBe(testDS.xmlId);
				expect(req.request.params.get("startDate")).toBe(twoSecondsAgo.toISOString());
				expect(req.request.params.get("endDate")).toBe(now.toISOString());
				expect(req.request.method).toBe("GET");

				const metricType = req.request.params.get("metricType");
				switch (metricType) {
					case "tps_total":
						req.flush({response: totalResponse});
						break;
					case "tps_2xx":
						req.flush({response: successResponse});
						break;
					case "tps_3xx":
						req.flush({response: redirectionResponse});
						break;
					case "tps_4xx":
						req.flush({response: clientErrorResponse});
						break;
					case "tps_5xx":
						req.flush({response: serverErrorResponse});
						break;
					default:
						fail(`unexpected metricType: ${metricType}`);
				}
			}

			await expectAsync(responseP).toBeResolvedTo({
				clientError: constructDataSetFromResponse(clientErrorResponse),
				redirection: constructDataSetFromResponse(redirectionResponse),
				serverError: constructDataSetFromResponse(serverErrorResponse),
				success: constructDataSetFromResponse(successResponse),
				total: constructDataSetFromResponse(totalResponse),
			});
		});

		it("throws an error when encountering an unknown data series type", async () => {
			const responseP = service.getAllDSTPSData(testDS.xmlId, twoSecondsAgo, now, interval, true);

			const requests = httpTestingController.match(r => r.url === `/api/${service.apiVersion}/deliveryservice_stats`);
			if (requests.length !== 5) {
				return fail(`expected getting all TPS stats to make 5 requests, found: ${requests.length}`);
			}

			for (const req of requests) {
				expect(req.request.params.keys().length).toBe(6);
				expect(req.request.params.get("serverType")).toBe("mid");
				expect(req.request.params.get("interval")).toBe(interval);
				expect(req.request.params.get("deliveryServiceName")).toBe(testDS.xmlId);
				expect(req.request.params.get("startDate")).toBe(twoSecondsAgo.toISOString());
				expect(req.request.params.get("endDate")).toBe(now.toISOString());
				expect(req.request.method).toBe("GET");

				req.flush({
					response: {
						series: {
							columns: ["time", "mean"],
							count: 0,
							name: "invalid",
							values: [
								[
									twoSecondsAgo,
									null
								],
								[
									now,
									0
								]
							],
						},
						summary: {
							average: 1,
							count: 2,
							fifthPercentile: 3,
							max: 4,
							min: 5,
							ninetyEightPercentile: 6,
							ninetyFifthPercentile: 7
						}
					}
				});
			}

			await expectAsync(responseP).toBeRejected();
		});

		it("gets DS health", async () => {
			const responseP = service.getDSHealth(testDS);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/deliveryservices/${testDS.id}/health`);
			expect(req.request.params.keys().length).toBe(0);
			expect(req.request.method).toBe("GET");
			req.flush(health);
			await expectAsync(responseP).toBeResolvedTo(health.response);
		});

		it("gets DS health by DS ID", async () => {
			const responseP = service.getDSHealth(testDS.id);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/deliveryservices/${testDS.id}/health`);
			expect(req.request.params.keys().length).toBe(0);
			expect(req.request.method).toBe("GET");
			req.flush(health);
			await expectAsync(responseP).toBeResolvedTo(health.response);
		});

		it("gets DS capacity", async () => {
			const responseP = service.getDSCapacity(testDS);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/deliveryservices/${testDS.id}/capacity`);
			expect(req.request.params.keys().length).toBe(0);
			expect(req.request.method).toBe("GET");
			req.flush(capacity);
			await expectAsync(responseP).toBeResolvedTo(capacity.response);
		});

		it("gets DS capacity by DS ID", async () => {
			const responseP = service.getDSCapacity(testDS.id);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/deliveryservices/${testDS.id}/capacity`);
			expect(req.request.params.keys().length).toBe(0);
			expect(req.request.method).toBe("GET");
			req.flush(capacity);
			await expectAsync(responseP).toBeResolvedTo(capacity.response);
		});
	});

	describe("basic DS operations methods", () => {
		it("gets multiple Delivery Services", async () => {
			const responseP = service.getDeliveryServices();
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/deliveryservices`);
			expect(req.request.params.keys().length).toBe(0);
			expect(req.request.method).toBe("GET");
			req.flush({response: [testDS]});
			await expectAsync(responseP).toBeResolvedTo([testDS]);
		});

		it("gets a single Delivery Service by ID", async () => {
			const responseP = service.getDeliveryServices(testDS.id);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/deliveryservices`);
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("id")).toBe(String(testDS.id));
			expect(req.request.method).toBe("GET");
			req.flush({response: [testDS]});
			await expectAsync(responseP).toBeResolvedTo(testDS);
		});

		it("throws an error when more than one Delivery Service exists by a given XMLID", async () => {
			const responseP = service.getDeliveryServices(testDS.xmlId);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/deliveryservices`);
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("xmlId")).toBe(testDS.xmlId);
			expect(req.request.method).toBe("GET");
			req.flush({response: [testDS, testDS]});
			await expectAsync(responseP).toBeRejected();
		});

		it("submits requests to create new Delivery Services", async () => {
			const responseP = service.createDeliveryService(testDS);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/deliveryservices`);
			expect(req.request.method).toBe("POST");
			expect(req.request.body).toEqual(testDS);
			req.flush({response: testDS});
			await expectAsync(responseP).toBeResolvedTo(testDS);
		});
	});

	it("only requests DS types when the cache is empty", async () => {
		const types = [
			{
				description: "HTTP desc",
				id: 1,
				lastUpdated: new Date(),
				name: "HTTP",
				useInTable: "deliveryService"
			},
			{
				description: "DNS desc",
				id: 2,
				lastUpdated: new Date(),
				name: "DNS",
				useInTable: "deliveryService"
			},
		];

		let responseP = service.getDSTypes();
		const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/types`);
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("useInTable")).toBe("deliveryservice");
		expect(req.request.method).toBe("GET");
		req.flush({response: types});
		await expectAsync(responseP).toBeResolvedTo(types);

		responseP = service.getDSTypes();
		httpTestingController.expectNone((): true => true);
		await expectAsync(responseP).toBeResolvedTo(types);
	});

	it("gets steering configurations", async () => {
		const response = [{
			clientSteering: false,
			deliveryService: "testquest",
			filters: [],
			targets: []
		}];
		const responseP = service.getSteering();
		const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/steering`);
		expect(req.request.params.keys().length).toBe(0);
		expect(req.request.method).toBe("GET");
		req.flush({response});
		await expectAsync(responseP).toBeResolvedTo(response);
	});

	it("gets DS ssl keys", async () => {
		let resp = service.getSSLKeys(testDS.xmlId);
		let req = httpTestingController.expectOne(r => r.url ===
			`/api/${service.apiVersion}/deliveryservices/xmlId/${testDS.xmlId}/sslkeys`);
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("decode")).toBe("true");
		expect(req.request.method).toBe("GET");
		req.flush({response: testDSSSLKeys});

		await expectAsync(resp).toBeResolvedTo(testDSSSLKeys);

		resp = service.getSSLKeys(testDS);
		req = httpTestingController.expectOne(r => r.url ===
			`/api/${service.apiVersion}/deliveryservices/xmlId/${testDS.xmlId}/sslkeys`);
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("decode")).toBe("true");
		expect(req.request.method).toBe("GET");
		req.flush({response: testDSSSLKeys});

		await expectAsync(resp).toBeResolvedTo(testDSSSLKeys);
	});

	afterEach(() => {
		httpTestingController.verify();
	});
});
