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
import type { ResponseCDN } from "trafficops-types";

import { CDNService } from "./cdn.service";

import { APITestingModule } from ".";

describe("CDNService", () => {
	let service: CDNService;
	let cdn: ResponseCDN;

	beforeEach(() => {
		TestBed.configureTestingModule({
			imports: [APITestingModule],
		});
		service = TestBed.inject(CDNService);

		expect(service.cdns.length).toBeGreaterThanOrEqual(2);
		cdn = service.cdns[0];
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	it("gets multiple CDNs", async () => {
		await expectAsync(service.getCDNs()).toBeResolvedTo(service.cdns);
	});

	it("gets a single CDN by ID", async () => {
		await expectAsync(service.getCDNs(cdn.id)).toBeResolvedTo(cdn);
	});

	it("throws an error when requesting a specific, but non-existent, CDN", async () => {
		await expectAsync(service.getCDNs(-1)).toBeRejected();
	});

	it("creates a new CDN", async () => {
		const newCDN = {
			dnssecEnabled: false,
			domainName: "quest",
			name: `test-${(new Date()).valueOf()}`
		};
		const initialLength = service.cdns.length;
		const created = await service.createCDN(newCDN);
		expect(service.cdns).toContain(created);
		expect(service.cdns.length).toEqual(initialLength+1);
	});

	it("updates an existing CDN", async () => {
		const current = {...service.cdns[1]};
		current.name += "quest";
		const {name} = current;
		const updated = await service.updateCDN(current);
		expect(updated.name).toBe(name);
		expect(service.cdns[1]).toBe(updated);
	});

	it("updates an existing CDN by ID", async () => {
		const current = {...service.cdns[1]};
		current.name += "quest";
		const {name} = current;
		const updated = await service.updateCDN(current.id, current);
		expect(updated.name).toBe(name);
		expect(service.cdns[1]).toBe(updated);
	});

	it("throws an error when asked to update a non-existent CDN", async () => {
		await expectAsync(service.updateCDN(-1, {dnssecEnabled: false, domainName: "", name: ""})).toBeRejected();
	});

	it("throws an error for invalid call signatures to updateCDN", async () => {
		const responseP = (service as unknown as {updateCDN: (id: number) => Promise<unknown>}).updateCDN(service.cdns[0].id);
		await expectAsync(responseP).toBeRejected();
	});

	it("deletes an existing CDN", async () => {
		const initialLength = service.cdns.length;
		await service.deleteCDN(cdn);
		expect(service.cdns.length).toEqual(initialLength - 1);
		expect(service.cdns).not.toContain(cdn);
	});

	it("deletes an existing CDN by ID", async () => {
		const initialLength = service.cdns.length;
		await service.deleteCDN(cdn.id);
		expect(service.cdns.length).toEqual(initialLength - 1);
		expect(service.cdns).not.toContain(cdn);
	});

	it("throws an error when asked to delete a non-existent CDN", async () => {
		await expectAsync(service.deleteCDN(-1)).toBeRejected();
	});

	it("pretends to queue updates on a CDN", async () => {
		await expectAsync(service.queueServerUpdates(cdn)).toBeResolvedTo({action: "queue", cdnId: cdn.id});
	});
	it("throws an error when attempting to queue a non-existent CDN", async () => {
		await expectAsync(service.queueServerUpdates(-1)).toBeRejected();
	});
	it("pretends to dequeue updates on a CDN", async () => {
		await expectAsync(service.dequeueServerUpdates(cdn)).toBeResolvedTo({action: "dequeue", cdnId: cdn.id});
	});
	it("throws an error when attempting to de-queue a non-existent CDN", async () => {
		await expectAsync(service.dequeueServerUpdates(-1)).toBeRejected();
	});
});
