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

import { PhysicalLocationService } from "./physical-location.service";

import { APITestingModule } from ".";

describe("PhysicalLocationService", () => {
	let service: PhysicalLocationService;

	beforeEach(async () => {
		TestBed.configureTestingModule({
			imports: [APITestingModule],
			providers: [
				PhysicalLocationService,
			]
		});
		service = TestBed.inject(PhysicalLocationService);
		expect(service.physicalLocations.length).toBeGreaterThanOrEqual(1);
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	it("gets multiple Physical Locations", async () => {
		await expectAsync(service.getPhysicalLocations()).toBeResolvedTo(service.physicalLocations);
	});

	it("gets a single Physical Location by ID", async () => {
		const physLoc = service.physicalLocations[0];
		await expectAsync(service.getPhysicalLocations(physLoc.id)).toBeResolvedTo(physLoc);
	});

	it("gets a single Physical Location by name", async () => {
		const physLoc = service.physicalLocations[0];
		await expectAsync(service.getPhysicalLocations(physLoc.name)).toBeResolvedTo(physLoc);
	});

	it("throws an error when the specifically requested Physical Location doesn't exist", async () => {
		await expectAsync(service.getPhysicalLocations(-1)).toBeRejected();
	});

	it("creates new Physical Locations", async () => {
		const name = `test-${(new Date()).valueOf()}`;
		const physLoc = {
			address: "",
			city: "",
			name,
			regionId: 1,
			shortName: name,
			state: "",
			zip: ""
		};
		const initialLength = service.physicalLocations.length;
		const created = await service.createPhysicalLocation(physLoc);
		expect(service.physicalLocations).toContain(created);
		expect(service.physicalLocations.length).toEqual(initialLength+1);
	});

	it("updates existing Physical Locations", async () => {
		const physLoc = {...service.physicalLocations[0]};
		physLoc.name += String((new Date()).valueOf());
		const {name} = physLoc;
		const updated = await service.updatePhysicalLocation(physLoc);
		expect(updated.name).toBe(name);
		expect(service.physicalLocations[0]).toBe(updated);
	});

	it("throws an error when asked to update a non-existent Physical Location", async () => {
		await expectAsync(service.updatePhysicalLocation({...service.physicalLocations[0], id: -1})).toBeRejected();
	});

	it("deletes Physical Locations", async () => {
		const physLoc = service.physicalLocations[0];
		const initialLength = service.physicalLocations.length;
		await service.deletePhysicalLocation(physLoc);
		expect(service.physicalLocations).not.toContain(physLoc);
		expect(service.physicalLocations.length).toEqual(initialLength-1);
	});

	it("deletes Physical Locations by ID", async () => {
		const physLoc = service.physicalLocations[0];
		const initialLength = service.physicalLocations.length;
		await service.deletePhysicalLocation(physLoc.id);
		expect(service.physicalLocations).not.toContain(physLoc);
		expect(service.physicalLocations.length).toEqual(initialLength-1);
	});

	it("throws an error when asked to delete a Physical Location that doesn't exist", async () => {
		await expectAsync(service.deletePhysicalLocation(-1)).toBeRejected();
	});
});
