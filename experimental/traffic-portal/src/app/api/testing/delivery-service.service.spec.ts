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
import {
	ProfileType,
	type ResponseDeliveryServiceSSLKey,
	type RequestAnyMapDeliveryService,
	type RequestSteeringDeliveryService,
	type ResponseDeliveryService
} from "trafficops-types";

import { ProfileService } from "../profile.service";

import { DeliveryServiceService } from "./delivery-service.service";

import { APITestingModule } from ".";

const now = new Date();
const twoSecondsAgo = new Date(now.getTime() - 2000);

describe("DeliveryServiceService", () => {
	let service: DeliveryServiceService;
	let testDS: ResponseDeliveryService;

	beforeEach(async () => {
		TestBed.configureTestingModule({
			imports: [APITestingModule],
			providers: [
				DeliveryServiceService,
			]
		});
		service = TestBed.inject(DeliveryServiceService);
		expect(service.deliveryServiceTypes.length).toBeGreaterThanOrEqual(1);
		const requestDS: RequestAnyMapDeliveryService = {
			active: false,
			cacheurl: null,
			cdnId: 1,
			displayName: "Test Quest",
			dscp: 2,
			ecsEnabled: false,
			geoLimit: 0,
			geoProvider: 1,
			httpBypassFqdn: null,
			infoUrl: null,
			initialDispersion: 2,
			logsEnabled: false,
			orgServerFqdn: "",
			regionalGeoBlocking: false,
			remapText: null,
			routingName: "cdn",
			serviceCategory: null,
			tenantId: 1,
			typeId: service.deliveryServiceTypes[0].id,
			xmlId: "xml",
		};
		testDS = await service.createDeliveryService(requestDS);
		expect(testDS).toBeTruthy();
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	describe("Stats methods", () => {
		const interval = "60s";

		it("generates KBPS stats", async () => {
			const stats = await service.getDSKBPS(testDS, twoSecondsAgo, now, interval, false);
			expect(stats.series).toBeDefined();
			expect(stats.series?.values.length).toBeGreaterThanOrEqual(1);
		});

		it("generates KBPS stats, returning only the data series", async () => {
			const stats = await service.getDSKBPS(testDS, twoSecondsAgo, now, interval, true, true);
			expect(Array.isArray(stats)).toBeTrue();
			expect(stats.length).toBeGreaterThanOrEqual(1);
		});

		it("throws an error when data generation is requested for a non-existent Delivery Service", async () => {
			await expectAsync(service.getDSKBPS("", twoSecondsAgo, now, interval, false)).toBeRejected();
		});

		it("generates TPS stats", async () => {
			const stats = await service.getDSTPS(testDS, twoSecondsAgo, now, interval);
			expect(stats.series).toBeDefined();
			expect(stats.series?.values.length).toBeGreaterThanOrEqual(1);
		});
		it("throws an error when TPS stats are requested for a non-existent DS", async () => {
			await expectAsync(service.getDSTPS(`${testDS.xmlId}-${(new Date()).valueOf()}`, twoSecondsAgo, now, interval)).toBeRejected();
		});

		it("generates all TPS data for a DS", async () => {
			const stats = await service.getAllDSTPSData(testDS, twoSecondsAgo, now, interval);
			expect(stats.informational).toBeUndefined();
			expect(stats.serverError.dataSet.data.length).toBeGreaterThanOrEqual(1);
			expect(stats.clientError.dataSet.data.length).toBeGreaterThanOrEqual(1);
			expect(stats.redirection.dataSet.data.length).toBeGreaterThanOrEqual(1);
			expect(stats.success.dataSet.data.length).toBeGreaterThanOrEqual(1);
			expect(stats.total.dataSet.data.length).toBeGreaterThanOrEqual(1);
		});

		it("throws an error when stats are requested for a non-existent Delivery Service", async () => {
			await expectAsync(service.getAllDSTPSData("", twoSecondsAgo, now, interval)).toBeRejected();
		});

		it("gets DS health", async () => {
			const health = await service.getDSHealth(testDS);
			expect(health.cacheGroups).not.toBeNull();
			expect(health.totalOffline).toBeGreaterThan(0);
			expect(health.totalOnline).toBeGreaterThan(0);
		});
		it("gets DS health by ID", async () => {
			const health = await service.getDSHealth(testDS.id);
			expect(health.cacheGroups).not.toBeNull();
			expect(health.totalOffline+health.totalOnline).toBeCloseTo(100);
		});
		it("throws an error when health is requested for a non-existent DS", async () => {
			await expectAsync(service.getDSHealth(-1)).toBeRejected();
		});

		it("gets DS capacity", async () => {
			const cap = await service.getDSCapacity(testDS);
			expect(cap.availablePercent).toBeGreaterThan(0);
			expect(cap.maintenancePercent).toBeGreaterThan(0);
			expect(cap.unavailablePercent).toBeGreaterThan(0);
			expect(cap.utilizedPercent).toBeGreaterThan(0);
		});
		it("gets DS capacity by DS ID", async () => {
			const cap = await service.getDSCapacity(testDS.id);
			expect(cap.availablePercent).toBeGreaterThan(0);
			expect(cap.maintenancePercent).toBeGreaterThan(0);
			expect(cap.unavailablePercent).toBeGreaterThan(0);
			expect(cap.utilizedPercent).toBeGreaterThan(0);
		});
		it("throws an error when capacity is requested for a non-existent DS", async () => {
			await expectAsync(service.getDSCapacity(-1)).toBeRejected();
		});
	});

	describe("basic DS operations methods", () => {
		it("gets multiple Delivery Services", async () => {
			await expectAsync(service.getDeliveryServices()).toBeResolvedTo(service.deliveryServices);
		});

		it("gets a single Delivery Service by ID", async () => {
			await expectAsync(service.getDeliveryServices(testDS.id)).toBeResolvedTo(testDS);
		});

		it("gets a single Delivery Service by xml_id", async () => {
			await expectAsync(service.getDeliveryServices(testDS.xmlId)).toBeResolvedTo(testDS);
		});

		it("throws an error when a single, non-existent DS is requested", async () => {
			await expectAsync(service.getDeliveryServices(-1)).toBeRejected();
		});

		it("creates Delivery Services with Profiles", async () => {
			const steeringType = service.deliveryServiceTypes.find(t => t.name === "STEERING");
			if (!steeringType) {
				return fail("no steering-type for DSes");
			}

			const profile = await TestBed.inject(ProfileService).createProfile({
				cdn: 1,
				description: "",
				name: `DS Test Profile-${(new Date()).valueOf()}`,
				routingDisabled: false,
				type: ProfileType.DS_PROFILE
			});

			const xmlId = `DS-with-Profile Creation Test-${(new Date()).valueOf()}`;
			const ds: RequestSteeringDeliveryService = {
				active: true,
				cacheurl: null,
				cdnId: 1,
				consistentHashQueryParams: ["format"],
				displayName: xmlId,
				dscp: 2,
				geoLimit: 0,
				geoProvider: 0,
				httpBypassFqdn: null,
				infoUrl: null,
				logsEnabled: true,
				profileId: profile.id,
				regionalGeoBlocking: false,
				remapText: null,
				tenantId: 1,
				tlsVersions: ["5.27"],
				typeId: steeringType.id,
				xmlId,
			};

			const newDS = await service.createDeliveryService(ds);
			expect(newDS.profileDescription).toBe(profile.description);
			expect(newDS.profileName).toBe(profile.name);
			expect(newDS.profileId).toBe(profile.id);
		});
	});

	describe("SSL key methods", () => {
		let keys: ResponseDeliveryServiceSSLKey[];
		beforeEach(async () => {
			keys = service.dsSSLKeys;
			expect(keys.length).toBeGreaterThanOrEqual(1);
		});

		it("gets a DS's SSL Keys", async () => {
			const key = keys[0];
			const ds = await service.getDeliveryServices(key.deliveryservice);
			if (!ds) {
				return fail("found an SSL key with no corresponding DS");
			}
			await expectAsync(service.getSSLKeys(ds)).toBeResolvedTo(key);
		});

		it("throws an error when asked to retrieve keys for a non-existent DS", async () => {
			await expectAsync(service.getSSLKeys("")).toBeRejected();
		});
	});

	it("exposes its DS types", async () => {
		await expectAsync(service.getDSTypes()).toBeResolvedTo(service.deliveryServiceTypes);
	});

	it("gets (empty) steering configurations", async () => {
		await expectAsync(service.getSteering()).toBeResolvedTo([]);
	});
});
