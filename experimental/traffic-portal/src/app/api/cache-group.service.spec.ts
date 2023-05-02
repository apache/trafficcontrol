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
import { ResponseCacheGroup } from "trafficops-types";

import { CacheGroupService } from "./cache-group.service";

describe("CacheGroupService", () => {
	let service: CacheGroupService;
	let httpTestingController: HttpTestingController;

	const cg = {
		fallbackToClosest: true,
		fallbacks: [],
		id: 1,
		lastUpdated: new Date(),
		latitude: 0,
		localizationMethods: [],
		longitude: 0,
		name: "test",
		parentCachegroupId: null,
		parentCachegroupName: null,
		secondaryParentCachegroupId: null,
		secondaryParentCachegroupName: null,
		shortName: "test",
		typeId: 1,
		typeName: "EDGE_LOC"
	};

	beforeEach(() => {
		TestBed.configureTestingModule({
			imports: [HttpClientTestingModule],
			providers: [
				CacheGroupService,
			]
		});
		service = TestBed.inject(CacheGroupService);
		httpTestingController = TestBed.inject(HttpTestingController);
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	describe("Cache Group methods", () => {

		it("gets multiple Cache Groups", async () => {
			const responseP = service.getCacheGroups();
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/cachegroups`);
			expect(req.request.method).toBe("GET");
			expect(req.request.body).toBeNull();
			expect(req.request.params.keys().length).toBe(0);
			const data = {
				response: [
					cg
				]
			};
			req.flush(data);
			await expectAsync(responseP).toBeResolvedTo(data.response);
		});

		it("gets a single Cache Group by ID", async () => {
			const responseP = service.getCacheGroups(cg.id);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/cachegroups`);
			expect(req.request.method).toBe("GET");
			expect(req.request.body).toBeNull();
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("id")).toBe(String(cg.id));
			const data = {
				response: [
					cg
				]
			};
			req.flush(data);
			await expectAsync(responseP).toBeResolvedTo(cg);
		});

		it("gets a single Cache Group by name", async () => {
			const responseP = service.getCacheGroups(cg.name);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/cachegroups`);
			expect(req.request.method).toBe("GET");
			expect(req.request.body).toBeNull();
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("name")).toBe(cg.name);
			const data = {
				response: [
					cg
				]
			};
			req.flush(data);
			await expectAsync(responseP).toBeResolvedTo(cg);
		});

		it("throws an error when multiple Cache Groups share an ID", async () => {
			const responseP = service.getCacheGroups(cg.id);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/cachegroups`);
			expect(req.request.method).toBe("GET");
			expect(req.request.body).toBeNull();
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("id")).toBe(String(cg.id));
			const data = {
				response: [
					cg,
					{
						fallbackToClosest: true,
						fallbacks: [],
						id: cg.id,
						lastUpdated: new Date(),
						latitude: 0,
						localizationMethods: [],
						longitude: 0,
						name: `${cg.name}quest`,
						parentCachegroupId: null,
						parentCachegroupName: null,
						secondaryParentCachegroupId: null,
						secondaryParentCachegroupName: null,
						shortName: `${cg.name}quest`,
						typeId: 1,
						typeName: "EDGE_LOC"
					}
				]
			};
			req.flush(data);
			await expectAsync(responseP).toBeRejected();
		});

		it("deletes a Cache Group", async () => {
			const responseP = service.deleteCacheGroup(cg);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/cachegroups/${cg.id}`);
			expect(req.request.method).toBe("DELETE");
			expect(req.request.body).toBeNull();
			expect(req.request.params.keys().length).toBe(0);
			req.flush({alerts: [{level: "success", text: "deleted the Cache Group"}]});
			await expectAsync(responseP).toBeResolved();
		});

		it("deletes a Cache Group by ID", async () => {
			const responseP = service.deleteCacheGroup(cg.id);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/cachegroups/${cg.id}`);
			expect(req.request.method).toBe("DELETE");
			expect(req.request.body).toBeNull();
			expect(req.request.params.keys().length).toBe(0);
			req.flush({alerts: [{level: "success", text: "deleted the Cache Group"}]});
			await expectAsync(responseP).toBeResolved();
		});

		it("creates a Cache Group", async () => {
			const responseP = service.createCacheGroup(cg);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/cachegroups`);
			expect(req.request.method).toBe("POST");
			expect(req.request.body).toEqual(cg);
			expect(req.request.params.keys().length).toBe(0);
			req.flush({response: cg});
			await expectAsync(responseP).toBeResolvedTo(cg);
		});

		it("updates a Cache Group", async () => {
			const responseP = service.updateCacheGroup(cg);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/cachegroups/${cg.id}`);
			expect(req.request.method).toBe("PUT");
			expect(req.request.body).toEqual(cg);
			expect(req.request.params.keys().length).toBe(0);
			req.flush({response: cg});
			await expectAsync(responseP).toBeResolvedTo(cg);
		});

		it("updates a Cache Group by ID", async () => {
			const responseP = service.updateCacheGroup(cg.id, cg);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/cachegroups/${cg.id}`);
			expect(req.request.method).toBe("PUT");
			expect(req.request.body).toEqual(cg);
			expect(req.request.params.keys().length).toBe(0);
			req.flush({response: cg});
			await expectAsync(responseP).toBeResolvedTo(cg);
		});

		it("throws an error for invalid call signatures to updateCacheGroup", async () => {
			const responseP = (service as unknown as {updateCacheGroup: (id: number) => Promise<ResponseCacheGroup>}).updateCacheGroup(
				cg.id
			);
			httpTestingController.expectNone({method: "PUT"});
			await expectAsync(responseP).toBeRejected();
		});

		describe("queueing and de-queueing updates", () => {
			it("queues updates on a Cache Group", async () => {
				let responseP = service.queueCacheGroupUpdates(cg, 1, "queue");
				let req = httpTestingController.expectOne(`/api/${service.apiVersion}/cachegroups/${cg.id}/queue_update`);
				expect(req.request.method).toBe("POST");
				expect(req.request.body).toEqual({
					action: "queue",
					cdnId: 1
				});
				expect(req.request.params.keys().length).toBe(0);
				req.flush({});
				await expectAsync(responseP).toBeResolved();

				responseP = service.queueCacheGroupUpdates(cg, "testquest");
				req = httpTestingController.expectOne(`/api/${service.apiVersion}/cachegroups/${cg.id}/queue_update`);
				expect(req.request.method).toBe("POST");
				expect(req.request.body).toEqual({
					action: "queue",
					cdn: "testquest"
				});
				expect(req.request.params.keys().length).toBe(0);
				req.flush({});
				await expectAsync(responseP).toBeResolved();

				const cdn = {
					dnssecEnabled: false,
					domainName: "test.quest",
					id: 1,
					lastUpdated: new Date(),
					name: "TestQuest"
				};
				responseP = service.queueCacheGroupUpdates(cg.id, cdn);
				req = httpTestingController.expectOne(`/api/${service.apiVersion}/cachegroups/${cg.id}/queue_update`);
				expect(req.request.method).toBe("POST");
				expect(req.request.body).toEqual({
					action: "queue",
					cdn: cdn.name
				});
				expect(req.request.params.keys().length).toBe(0);
				req.flush({});
				await expectAsync(responseP).toBeResolved();
			});

			it("queues updates on a Cache Group using a literal request", async () => {
				const queueRequest = {action: "queue" as const, cdn: "testquest"};
				const responseP = service.queueCacheGroupUpdates(cg, queueRequest);
				const req = httpTestingController.expectOne(`/api/${service.apiVersion}/cachegroups/${cg.id}/queue_update`);
				expect(req.request.method).toBe("POST");
				expect(req.request.body).toEqual(queueRequest);
				expect(req.request.params.keys().length).toBe(0);
				req.flush({});

				await expectAsync(responseP).toBeResolved();
			});
		});
	});

	describe("Divisions methods", () => {
		const div = {
			id: 1,
			lastUpdated: new Date(),
			name: "testquest"
		};

		it("creates a division", async () => {
			const responseP = service.createDivision(div);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/divisions`);
			expect(req.request.method).toBe("POST");
			expect(req.request.body).toEqual(div);
			expect(req.request.params.keys().length).toBe(0);
			req.flush({response: div});

			await expectAsync(responseP).toBeResolvedTo(div);
		});

		it("gets divisions", async () => {
			const responseP = service.getDivisions();
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/divisions`);
			expect(req.request.method).toBe("GET");
			expect(req.request.body).toBeNull();
			expect(req.request.params.keys().length).toBe(0);
			req.flush({response: [div]});

			await expectAsync(responseP).toBeResolvedTo([div]);
		});

		it("gets a single division by ID", async () => {
			const responseP = service.getDivisions(1);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/divisions`);
			expect(req.request.method).toBe("GET");
			expect(req.request.body).toBeNull();
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("id")).toBe("1");
			req.flush({response: [div]});

			await expectAsync(responseP).toBeResolvedTo(div);
		});

		it("gets a single division by name", async () => {
			const responseP = service.getDivisions("testquest");
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/divisions`);
			expect(req.request.method).toBe("GET");
			expect(req.request.body).toBeNull();
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("name")).toBe("testquest");
			req.flush({response: [div]});

			await expectAsync(responseP).toBeResolvedTo(div);
		});

		it("deletes a division", async () => {
			const responseP = service.deleteDivision(div);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/divisions/${div.id}`);
			expect(req.request.method).toBe("DELETE");
			expect(req.request.body).toBeNull();
			expect(req.request.params.keys().length).toBe(0);
			req.flush({response: div});

			await expectAsync(responseP).toBeResolvedTo(div);
		});

		it("deletes a division by ID", async () => {
			const responseP = service.deleteDivision(div.id);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/divisions/${div.id}`);
			expect(req.request.method).toBe("DELETE");
			expect(req.request.body).toBeNull();
			expect(req.request.params.keys().length).toBe(0);
			req.flush({response: div});

			await expectAsync(responseP).toBeResolvedTo(div);
		});

		it("updates divisions", async () => {
			const responseP = service.updateDivision(div);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/divisions/${div.id}`);
			expect(req.request.method).toBe("PUT");
			expect(req.request.body).toEqual(div);
			expect(req.request.params.keys().length).toBe(0);
			req.flush({response: div});

			await expectAsync(responseP).toBeResolvedTo(div);
		});
	});

	describe("Regions methods", () => {
		const reg = {
			division: 1,
			divisionName: "testing division",
			id: 1,
			lastUpdated: new Date(),
			name: "testquest",
		};

		it("gets regions", async () => {
			const responseP = service.getRegions();
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/regions`);
			expect(req.request.method).toBe("GET");
			expect(req.request.body).toBeNull();
			expect(req.request.params.keys().length).toBe(0);
			req.flush({response: [reg]});

			await expectAsync(responseP).toBeResolvedTo([reg]);
		});
		it("gets a single region by name", async () => {
			const responseP = service.getRegions(reg.name);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/regions`);
			expect(req.request.method).toBe("GET");
			expect(req.request.body).toBeNull();
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("name")).toBe(reg.name);
			req.flush({response: [reg]});

			await expectAsync(responseP).toBeResolvedTo(reg);
		});
		it("gets a single region by ID", async () => {
			const responseP = service.getRegions(reg.id);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/regions`);
			expect(req.request.method).toBe("GET");
			expect(req.request.body).toBeNull();
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("id")).toBe(String(reg.id));
			req.flush({response: [reg]});

			await expectAsync(responseP).toBeResolvedTo(reg);
		});
		it("creates a new region", async () => {
			const responseP = service.createRegion(reg);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/regions`);
			expect(req.request.method).toBe("POST");
			expect(req.request.body).toEqual(reg);
			expect(req.request.params.keys().length).toBe(0);
			req.flush({response: reg});

			await expectAsync(responseP).toBeResolvedTo(reg);
		});
		it("updates a region", async () => {
			const responseP = service.updateRegion(reg);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/regions/${reg.id}`);
			expect(req.request.method).toBe("PUT");
			expect(req.request.body).toEqual(reg);
			expect(req.request.params.keys().length).toBe(0);
			req.flush({response: reg});

			await expectAsync(responseP).toBeResolvedTo(reg);
		});
		it("deletes a region", async () => {
			const responseP = service.deleteRegion(reg);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/regions`);
			expect(req.request.method).toBe("DELETE");
			expect(req.request.body).toBeNull();
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("id")).toBe(String(reg.id));
			req.flush({alerts: []});

			await expectAsync(responseP).toBeResolved();
		});
		it("deletes a region by ID", async () => {
			const responseP = service.deleteRegion(reg.id);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/regions`);
			expect(req.request.method).toBe("DELETE");
			expect(req.request.body).toBeNull();
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("id")).toBe(String(reg.id));
			req.flush({alerts: []});

			await expectAsync(responseP).toBeResolved();
		});
	});

	describe("ASN methods", () => {
		const asn = {
			asn: 100,
			cachegroup: cg.name,
			cachegroupId: cg.id,
			id: 1,
			lastUpdated: new Date()
		};
		it("gets multiple ASNs", async () => {
			const responseP = service.getASNs();
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/asns`);
			expect(req.request.method).toBe("GET");
			expect(req.request.body).toBeNull();
			req.flush({response: [asn]});

			await expectAsync(responseP).toBeResolvedTo([asn]);
		});
		it("gets a single ASNs by ID", async () => {
			const responseP = service.getASNs(asn.id);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/asns`);
			expect(req.request.method).toBe("GET");
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("id")).toBe(String(asn.id));
			expect(req.request.body).toBeNull();
			req.flush({response: [asn]});

			await expectAsync(responseP).toBeResolvedTo(asn);
		});
		it("deletes an ASN", async () => {
			const responseP = service.deleteASN(asn);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/asns/${asn.id}`);
			expect(req.request.method).toBe("DELETE");
			expect(req.request.body).toBeNull();
			expect(req.request.params.keys().length).toBe(0);
			req.flush({alerts: []});

			await expectAsync(responseP).toBeResolved();
		});
		it("deletes an ASN by ID", async () => {
			const responseP = service.deleteASN(asn.id);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/asns/${asn.id}`);
			expect(req.request.method).toBe("DELETE");
			expect(req.request.body).toBeNull();
			expect(req.request.params.keys().length).toBe(0);
			req.flush({alerts: []});

			await expectAsync(responseP).toBeResolved();
		});
		it("creates a new ASN", async () => {
			const responseP = service.createASN(asn);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/asns`);
			expect(req.request.method).toBe("POST");
			expect(req.request.body).toEqual(asn);
			req.flush({response: asn});

			await expectAsync(responseP).toBeResolved();
		});
		it("updates an existing ASN", async () => {
			const responseP = service.updateASN(asn);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/asns/${asn.id}`);
			expect(req.request.method).toBe("PUT");
			expect(req.request.body).toEqual(asn);
			req.flush({response: asn});

			await expectAsync(responseP).toBeResolved();
		});
	});

	describe("Coordinate methods", () => {
		const coord = {
			id: 1,
			lastUpdated: new Date(),
			latitude: 1.0,
			longitude: -1.0,
			name: "test"
		};
		it("gets multiple Coordinates", async () => {
			const responseP = service.getCoordinates();
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/coordinates`);
			expect(req.request.method).toBe("GET");
			expect(req.request.body).toBeNull();
			req.flush({response: [coord]});

			await expectAsync(responseP).toBeResolvedTo([coord]);
		});
		it("gets a single Coordinate by ID", async () => {
			const responseP = service.getCoordinates(coord.id);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/coordinates`);
			expect(req.request.method).toBe("GET");
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("id")).toBe(String(coord.id));
			expect(req.request.body).toBeNull();
			req.flush({response: [coord]});

			await expectAsync(responseP).toBeResolvedTo(coord);
		});
		it("gets a single Coordinate by name", async () => {
			const responseP = service.getCoordinates(coord.name);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/coordinates`);
			expect(req.request.method).toBe("GET");
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("name")).toBe(coord.name);
			expect(req.request.body).toBeNull();
			req.flush({response: [coord]});

			await expectAsync(responseP).toBeResolvedTo(coord);
		});
		it("deletes a Coordinate", async () => {
			const responseP = service.deleteCoordinate(coord);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/coordinates`);
			expect(req.request.method).toBe("DELETE");
			expect(req.request.body).toBeNull();
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("id")).toBe(String(coord.id));
			req.flush({alerts: []});

			await expectAsync(responseP).toBeResolved();
		});
		it("deletes a Coordinate by ID", async () => {
			const responseP = service.deleteCoordinate(coord.id);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/coordinates`);
			expect(req.request.method).toBe("DELETE");
			expect(req.request.body).toBeNull();
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("id")).toBe(String(coord.id));
			req.flush({alerts: []});

			await expectAsync(responseP).toBeResolved();
		});
		it("creates a new Coordinate", async () => {
			const responseP = service.createCoordinate(coord);
			const req = httpTestingController.expectOne(`/api/${service.apiVersion}/coordinates`);
			expect(req.request.method).toBe("POST");
			expect(req.request.body).toEqual(coord);
			req.flush({response: coord});

			await expectAsync(responseP).toBeResolved();
		});
		it("updates an existing ASN", async () => {
			const responseP = service.updateCoordinate(coord);
			const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/coordinates`);
			expect(req.request.method).toBe("PUT");
			expect(req.request.body).toEqual(coord);
			expect(req.request.params.keys().length).toBe(1);
			expect(req.request.params.get("id")).toBe(String(coord.id));
			req.flush({response: coord});

			await expectAsync(responseP).toBeResolved();
		});
	});

	afterEach(() => {
		httpTestingController.verify();
	});
});
