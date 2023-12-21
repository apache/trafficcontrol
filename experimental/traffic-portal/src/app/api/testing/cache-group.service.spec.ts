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
import { TestBed } from "@angular/core/testing";
import type { ResponseDivision } from "trafficops-types";

import { ProfileService } from "../profile.service";
import { ServerService } from "../server.service";

import { CacheGroupService as TestingCacheGroupService } from "./cache-group.service";

import { APITestingModule } from ".";

describe("TestingCacheGroupService", () => {
	let service: TestingCacheGroupService;

	beforeEach(() => {
		TestBed.configureTestingModule({
			imports: [APITestingModule]
		});
		service = TestBed.inject(TestingCacheGroupService);
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	describe("Cache Group methods", () => {
		beforeEach(() => {
			if (service.cacheGroups.length < 2) {
				return fail("need at least 2 Cache Groups in default testing data set");
			}
		});

		it("returns cache groups", async () => {
			expect(await service.getCacheGroups()).toEqual(service.cacheGroups);
		});

		it("gets a single Cache Group by name", async () => {
			const cg = service.cacheGroups[1];
			expect(await service.getCacheGroups(cg.name)).toEqual(cg);
		});

		it("gets a single Cache Group by ID", async () => {
			const cg = service.cacheGroups[1];
			expect(await service.getCacheGroups(cg.id)).toEqual(cg);
		});

		it("throws an error when a requested Cache Group doesn't exist", async () => {
			await expectAsync(service.getCacheGroups(-1)).toBeRejected();
		});

		it("deletes Cache Groups", async () => {
			const cg = service.cacheGroups[1];
			const initialLength = service.cacheGroups.length;

			await service.deleteCacheGroup(cg);
			expect(service.cacheGroups.length).toEqual(initialLength - 1);
			expect(service.cacheGroups).not.toContain(cg);
		});

		it("deletes Cache Groups by ID", async () => {
			const cg = service.cacheGroups[1];
			const initialLength = service.cacheGroups.length;

			await service.deleteCacheGroup(cg.id);
			expect(service.cacheGroups.length).toEqual(initialLength - 1);
			expect(service.cacheGroups).not.toContain(cg);
		});

		it("throws an error when asked to delete a non-existent Cache Group", async () => {
			await expectAsync(service.deleteCacheGroup(-1)).toBeRejected();
		});

		it("creates new Cache Groups", async () => {
			const initialLength = service.cacheGroups.length;
			const name = `test-${(new Date()).valueOf()}`;

			const cg = await service.createCacheGroup({
				name,
				shortName: name,
				typeId: 1
			});
			expect(service.cacheGroups.length).toEqual(initialLength+1);
			expect(service.cacheGroups).toContain(cg);
		});

		it("updates Cache Groups", async () => {
			const cg = {...service.cacheGroups[1]};
			cg.name += "quest";
			const {name} = cg;
			const updated = await service.updateCacheGroup(cg);
			expect(updated.name).toBe(name);
			expect(service.cacheGroups[1]).toBe(updated);
		});

		it("updates Cache Groups by ID", async () => {
			const cg = {...service.cacheGroups[1]};
			cg.name += "quest";
			const {name} = cg;
			const updated = await service.updateCacheGroup(cg.id, cg);
			expect(updated.name).toBe(name);
			expect(service.cacheGroups[1]).toBe(updated);
		});

		it("throws an error when asked to update a non-existent Cache Group", async () => {
			await expectAsync(service.updateCacheGroup(-1, {name: "", shortName: "", typeId: 1})).toBeRejected();
		});

		it("throws an error when called using an improper call signature", async () => {
			const responseP = (service as unknown as {updateCacheGroup: (id: number) => Promise<unknown>}).updateCacheGroup(1);
			await expectAsync(responseP).toBeRejected();
		});

		it("finds information for parentage", () => {
			const cg = service.cacheGroups[0];
			const parentage = service.getParents(cg.id, null);
			expect(parentage.parentCachegroupId).toBe(cg.id);
			expect(parentage.parentCachegroupName).toBe(cg.name);
			expect(parentage.secondaryParentCachegroupId).toBeNull();
			expect(parentage.secondaryParentCachegroupName).toBeNull();
		});
		it("finds information for secondary parentage", () => {
			const cg = service.cacheGroups[0];
			const parentage = service.getParents(null, cg.id);
			expect(parentage.parentCachegroupId).toBeNull();
			expect(parentage.parentCachegroupName).toBeNull();
			expect(parentage.secondaryParentCachegroupId).toBe(cg.id);
			expect(parentage.secondaryParentCachegroupName).toBe(cg.name);
		});
		it("throws errors when asked to find non-existent parents/secondary parents", () => {
			expect(()=>service.getParents(-1, null)).toThrow();
			expect(()=>service.getParents(null, -1)).toThrow();
		});

		describe("queuing/dequeing updates", () => {
			let serverService: ServerService;

			beforeEach(async () => {
				serverService = TestBed.inject(ServerService);
				const profileService = TestBed.inject(ProfileService);
				const cg = service.cacheGroups[0];
				const server = await serverService.createServer({
					cachegroupId: cg.id,
					cdnId: 1,
					domainName: "quest",
					hostName: "test",
					interfaces: [],
					physLocationId: 1,
					profileNames: [(await profileService.getProfiles())[0].name],
					statusId: (await serverService.getStatuses())[0].id,
					typeId: 1,
				});
				expect(server.updPending).toBeFalse();
			});

			it("queues cache group updates", async () => {
				const cg = service.cacheGroups[0];
				const cgServers = (await serverService.getServers()).filter(s => s.cachegroupId === cg.id);
				if(cgServers.length < 1) {
					return fail(`no servers in Cache Group #${cg.id} ('${cg.name}')`);
				}

				const filteredServers = cgServers.filter(s => s.cdnId === cgServers[0].cdnId);

				const response = await service.queueCacheGroupUpdates(cg, filteredServers[0].cdnId);
				expect(response.action).toBe("queue");
				expect(response.cachegroupID).toBe(cg.id);
				expect(response.cachegroupName).toBe(cg.name);
				expect(response.serverNames.sort()).toEqual(filteredServers.map(s => s.hostName).sort());

				// This is how the testing service implements this without importing the CDN Service
				expect(response.cdn).toBe(`${cgServers[0].cdnId}`);

				for (const server of await serverService.getServers()) {
					if (server.cachegroupId === cg.id && server.cdnId === filteredServers[0].cdnId) {
						expect(server.updPending).toBeTrue();
					}
				}
			});

			it("dequeues cache group updates", async () => {
				const cg = service.cacheGroups[0];
				const cgServers = (await serverService.getServers()).filter(s => s.cachegroupId === cg.id);
				if(cgServers.length < 1) {
					return fail(`no servers in Cache Group #${cg.id} ('${cg.name}')`);
				}

				const filteredServers = cgServers.filter(s => s.cdnId === cgServers[0].cdnId);

				const response = await service.queueCacheGroupUpdates(cg.id, filteredServers[0].cdnName, "dequeue");
				expect(response.action).toBe("dequeue");
				expect(response.cachegroupID).toBe(cg.id);
				expect(response.cachegroupName).toBe(cg.name);
				expect(response.serverNames.sort()).toEqual(filteredServers.map(s => s.hostName).sort());
				expect(response.cdn).toBe(cgServers[0].cdnName);

				for (const server of await serverService.getServers()) {
					if (server.cachegroupId === cg.id && server.cdnId === filteredServers[0].cdnId) {
						expect(server.updPending).toBeFalse();
					}
				}
			});

			it("queues cache group updates using a raw request body", async () => {
				const cg = service.cacheGroups[0];
				const cgServers = (await serverService.getServers()).filter(s => s.cachegroupId === cg.id);
				if(cgServers.length < 1) {
					return fail(`no servers in Cache Group #${cg.id} ('${cg.name}')`);
				}

				const filteredServers = cgServers.filter(s => s.cdnId === cgServers[0].cdnId);

				const response = await service.queueCacheGroupUpdates(cg, {action: "queue", cdnId: filteredServers[0].cdnId});
				expect(response.action).toBe("queue");
				expect(response.cachegroupID).toBe(cg.id);
				expect(response.cachegroupName).toBe(cg.name);
				expect(response.serverNames.sort()).toEqual(filteredServers.map(s => s.hostName).sort());

				// This is how the testing service implements this without importing the CDN Service
				expect(response.cdn).toBe(`${cgServers[0].cdnId}`);

				for (const server of await serverService.getServers()) {
					if (server.cachegroupId === cg.id && server.cdnId === filteredServers[0].cdnId) {
						expect(server.updPending).toBeTrue();
					}
				}
			});

			it("queues cache group updates using a raw request body containing CDN name", async () => {
				const cg = service.cacheGroups[0];
				const cgServers = (await serverService.getServers()).filter(s => s.cachegroupId === cg.id);
				if(cgServers.length < 1) {
					return fail(`no servers in Cache Group #${cg.id} ('${cg.name}')`);
				}

				const filteredServers = cgServers.filter(s => s.cdnId === cgServers[0].cdnId);

				const response = await service.queueCacheGroupUpdates(cg, {action: "queue", cdn: filteredServers[0].cdnName});
				expect(response.action).toBe("queue");
				expect(response.cachegroupID).toBe(cg.id);
				expect(response.cachegroupName).toBe(cg.name);
				expect(response.serverNames.sort()).toEqual(filteredServers.map(s => s.hostName).sort());
				expect(response.cdn).toBe(cgServers[0].cdnName);

				for (const server of await serverService.getServers()) {
					if (server.cachegroupId === cg.id && server.cdnId === filteredServers[0].cdnId) {
						expect(server.updPending).toBeTrue();
					}
				}
			});

			it("queues cache group updates using a CDN object", async () => {
				const cg = service.cacheGroups[0];
				const cgServers = (await serverService.getServers()).filter(s => s.cachegroupId === cg.id);
				if(cgServers.length < 1) {
					return fail(`no servers in Cache Group #${cg.id} ('${cg.name}')`);
				}

				const filteredServers = cgServers.filter(s => s.cdnId === cgServers[0].cdnId);

				const cdn = {
					dnssecEnabled: false,
					domainName: "-",
					id: 1,
					name: "ALL",
				};

				const response = await service.queueCacheGroupUpdates(cg, cdn);
				expect(response.action).toBe("queue");
				expect(response.cachegroupID).toBe(cg.id);
				expect(response.cachegroupName).toBe(cg.name);
				expect(response.serverNames.sort()).toEqual(filteredServers.map(s => s.hostName).sort());
				expect(response.cdn).toBe(cdn.name);

				for (const server of await serverService.getServers()) {
					if (server.cachegroupId === cg.id && server.cdnId === filteredServers[0].cdnId) {
						expect(server.updPending).toBeTrue();
					}
				}
			});

			it("throws an error when asked to queue updates on a non-existent Cache Group", async () => {
				await expectAsync(service.queueCacheGroupUpdates(-1, 1)).toBeRejected();
			});
		});
	});

	describe("Divisions methods", () => {
		beforeEach(() => {
			expect(service.divisions.length).toBeGreaterThanOrEqual(1);
		});

		it("gets Divisions", async () => {
			await expectAsync(service.getDivisions()).toBeResolvedTo(service.divisions);
		});

		it("gets a single Division by ID", async () => {
			const div = service.divisions[0];
			await expectAsync(service.getDivisions(div.id)).toBeResolvedTo(div);
		});

		it("gets a single Division by name", async () => {
			const div = service.divisions[0];
			await expectAsync(service.getDivisions(div.name)).toBeResolvedTo(div);
		});

		it("throws an error when asked to get a non-existent Division", async () => {
			await expectAsync(service.getDivisions(-1)).toBeRejected();
		});

		it("updates Divisions", async () => {
			const div = {...service.divisions[0]};
			div.name += "quest";
			const {name} = div;
			const updated = await service.updateDivision(div);
			expect(service.divisions[0].name).toEqual(name);
			expect(service.divisions[0]).toEqual(updated);
		});

		it("throws an error when asked to update a non-existent Division", async () => {
			await expectAsync(service.updateDivision({id: -1, lastUpdated: new Date(), name: ""})).toBeRejected();
		});

		it("creates Divisions", async () => {
			const initialLength = service.divisions.length;
			const div = await service.createDivision({name: `test-${(new Date()).valueOf()}`});
			expect(service.divisions).toContain(div);
			expect(service.divisions.length).toEqual(initialLength+1);
		});

		it("deletes Divisions", async () => {
			const div = service.divisions[0];
			const initialLength = service.divisions.length;
			await expectAsync(service.deleteDivision(div)).toBeResolvedTo(div);
			expect(service.divisions).not.toContain(div);
			expect(service.divisions.length).toEqual(initialLength-1);
		});

		it("throws an error when asked to delete a non-existent Division", async () => {
			await expectAsync(service.deleteDivision(-1)).toBeRejected();
		});
	});

	describe("Regions methods", () => {
		let div: ResponseDivision;

		beforeEach(() => {
			expect(service.regions.length).toBeGreaterThanOrEqual(1);
			expect(service.divisions.length).toBeGreaterThanOrEqual(1);
			div = service.divisions[0];
		});

		it("gets Regions", async () => {
			await expectAsync(service.getRegions()).toBeResolvedTo(service.regions);
		});

		it("gets a single Region by ID", async () => {
			const reg = service.regions[0];
			await expectAsync(service.getRegions(reg.id)).toBeResolvedTo(reg);
		});

		it("gets a single Region by name", async () => {
			const reg = service.regions[0];
			await expectAsync(service.getRegions(reg.name)).toBeResolvedTo(reg);
		});

		it("throws an error when asked to get a non-existent Region", async () => {
			await expectAsync(service.getRegions(-1)).toBeRejected();
		});

		it("updates Regions", async () => {
			const reg = {...service.regions[0]};
			reg.name += "quest";
			const {name} = reg;
			const updated = await service.updateRegion(reg);
			expect(service.regions[0].name).toEqual(name);
			expect(service.regions[0]).toEqual(updated);
		});

		it("throws an error when asked to update a non-existent Region", async () => {
			const reg = {
				division: div.id,
				divisionName: div.name,
				id: -1,
				lastUpdated: new Date(),
				name: ""
			};
			await expectAsync(service.updateRegion(reg)).toBeRejected();
		});

		it("creates Regions", async () => {
			const initialLength = service.regions.length;
			const reg = await service.createRegion({division: div.id, name: `test-${(new Date()).valueOf()}`});
			expect(service.regions).toContain(reg);
			expect(service.regions.length).toEqual(initialLength+1);
		});

		it("throws an error when attempting to create a Region in a non-existent Division", async () => {
			await expectAsync(service.createRegion({division: -1, name: ""})).toBeRejected();
		});

		it("deletes Regions", async () => {
			const reg = service.regions[0];
			const initialLength = service.regions.length;
			await expectAsync(service.deleteRegion(reg)).toBeResolved();
			expect(service.regions).not.toContain(reg);
			expect(service.regions.length).toEqual(initialLength-1);
		});

		it("throws an error when asked to delete a non-existent Region", async () => {
			await expectAsync(service.deleteRegion(-1)).toBeRejected();
		});
	});

	describe("Coordinates methods", () => {
		beforeEach(() => {
			expect(service.coordinates.length).toBeGreaterThanOrEqual(1);
		});

		it("gets Coordinates", async () => {
			await expectAsync(service.getCoordinates()).toBeResolvedTo(service.coordinates);
		});

		it("gets a single Coordinate by ID", async () => {
			const coord = service.coordinates[0];
			await expectAsync(service.getCoordinates(coord.id)).toBeResolvedTo(coord);
		});

		it("gets a single Coordinate by name", async () => {
			const coord = service.coordinates[0];
			await expectAsync(service.getCoordinates(coord.name)).toBeResolvedTo(coord);
		});

		it("throws an error when asked to get a non-existent Coordinate", async () => {
			await expectAsync(service.getCoordinates(-1)).toBeRejected();
		});

		it("updates Coordinates", async () => {
			const coord = {...service.coordinates[0]};
			coord.name += "quest";
			const {name} = coord;
			const updated = await service.updateCoordinate(coord);
			expect(service.coordinates[0].name).toEqual(name);
			expect(service.coordinates[0]).toEqual(updated);
		});

		it("throws an error when asked to update a non-existent Coordinate", async () => {
			const coord = {
				id: -1,
				lastUpdated: new Date(),
				latitude: 0,
				longitude: 0,
				name: ""
			};
			await expectAsync(service.updateCoordinate(coord)).toBeRejected();
		});

		it("creates Coordinates", async () => {
			const initialLength = service.coordinates.length;
			const coord = await service.createCoordinate({latitude: 0, longitude: 0, name: `test-${(new Date()).valueOf()}`});
			expect(service.coordinates).toContain(coord);
			expect(service.coordinates.length).toEqual(initialLength+1);
		});

		it("deletes Coordinates", async () => {
			const coord = service.coordinates[0];
			const initialLength = service.coordinates.length;
			await expectAsync(service.deleteCoordinate(coord)).toBeResolvedTo(undefined);
			expect(service.coordinates).not.toContain(coord);
			expect(service.coordinates.length).toEqual(initialLength-1);
		});

		it("throws an error when asked to delete a non-existent Coordinate", async () => {
			await expectAsync(service.deleteCoordinate(-1)).toBeRejected();
		});
	});

	describe("ASNs methods", () => {
		beforeEach(() => {
			expect(service.asns.length).toBeGreaterThanOrEqual(1);
			expect(service.cacheGroups.length).toBeGreaterThanOrEqual(1);
		});

		it("gets ASNs", async () => {
			await expectAsync(service.getASNs()).toBeResolvedTo(service.asns);
		});

		it("gets a single ASN by ID", async () => {
			const asn = service.asns[0];
			await expectAsync(service.getASNs(asn.id)).toBeResolvedTo(asn);
		});

		it("throws an error when asked to get a non-existent ASN", async () => {
			await expectAsync(service.getASNs(-1)).toBeRejected();
		});

		it("updates ASNs", async () => {
			const asnObj = {...service.asns[0]};
			asnObj.asn += 7;
			const {asn} = asnObj;
			const updated = await service.updateASN(asnObj);
			expect(service.asns[0].asn).toEqual(asn);
			expect(service.asns[0]).toEqual(updated);
		});

		it("throws an error when asked to update a non-existent ASN", async () => {
			const asn = {
				asn: 1,
				cachegroup: service.cacheGroups[0].name,
				cachegroupId: service.cacheGroups[0].id,
				id: -1,
				lastUpdated: new Date(),
			};
			await expectAsync(service.updateASN(asn)).toBeRejected();
		});

		it("creates ASNs", async () => {
			const initialLength = service.asns.length;
			const asn = await service.createASN({asn: 1, cachegroupId: service.cacheGroups[0].id});
			expect(service.asns).toContain(asn);
			expect(service.asns.length).toEqual(initialLength+1);
		});

		it("throws an error attempting to create an ASN in a non-existent Cache Group", async () => {
			await expectAsync(service.createASN({asn: 1, cachegroupId: -1})).toBeRejected();
		});

		it("deletes ASNs", async () => {
			const asn = service.asns[0];
			const initialLength = service.asns.length;
			await expectAsync(service.deleteASN(asn)).toBeResolvedTo(undefined);
			expect(service.asns).not.toContain(asn);
			expect(service.asns.length).toEqual(initialLength-1);
		});

		it("throws an error when asked to delete a non-existent ASN", async () => {
			await expectAsync(service.deleteASN(-1)).toBeRejected();
		});
	});
});
