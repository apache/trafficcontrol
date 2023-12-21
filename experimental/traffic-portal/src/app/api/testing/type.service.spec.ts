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

import { TypeService } from "./type.service";

import { APITestingModule } from ".";

describe("TestingTypeService", () => {
	let service: TypeService;

	beforeEach(() => {
		TestBed.configureTestingModule({
			imports: [APITestingModule],
		});
		service = TestBed.inject(TypeService);
		expect(service.types.length).toBeGreaterThanOrEqual(2);
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	it("gets multiple Types", async () => {
		await expectAsync(service.getTypes()).toBeResolvedTo(service.types);
	});

	it("gets a single Type by ID", async () => {
		const type = service.types[1];
		await expectAsync(service.getTypes(type.id)).toBeResolvedTo(type);
	});
	it("gets a single Type by name", async () => {
		const type = service.types[1];
		await expectAsync(service.getTypes(type.name)).toBeResolvedTo(type);
	});
	it("throws an error when asked to get a non-existent Type", async () => {
		await expectAsync(service.getTypes(-1)).toBeRejected();
	});
	it("gets Types for a specific table", async () => {
		const serverTypes = service.types.filter(t => t.useInTable === "server");
		expect(serverTypes.length).toBeGreaterThanOrEqual(1);
		expect(serverTypes.length).toBeLessThan(service.types.length);
		await expectAsync(service.getTypesInTable("server")).toBeResolvedTo(serverTypes);
	});
	it("gets server Types", async () => {
		const serverTypes = service.types.filter(t => t.useInTable === "server");
		expect(serverTypes.length).toBeGreaterThanOrEqual(1);
		expect(serverTypes.length).toBeLessThan(service.types.length);
		await expectAsync(service.getServerTypes()).toBeResolvedTo(serverTypes);
	});
	it("updates existing Types", async () => {
		const type = {...service.types[1]};
		type.name += String((new Date()).valueOf());
		const {name} = type;
		const updated = await service.updateType(type);
		expect(updated.name).toBe(name);
		expect(service.types[1]).toBe(updated);
	});
	it("throws an error when asked to update a non-existent Type", async () => {
		await expectAsync(service.updateType({...service.types[1], id: -1})).toBeRejected();
	});
	it("creates Types", async () => {
		const type = {
			description: "",
			name: String((new Date()).valueOf()),
			useInTable: "server"
		};
		const initialLength = service.types.length;
		const created = await service.createType(type);
		expect(service.types).toContain(created);
		expect(service.types.length).toEqual(initialLength+1);
	});
	it("deletes Types", async () => {
		const initialLength = service.types.length;
		const type = service.types[initialLength - 1];
		const deleted = await service.deleteType(type);
		expect(deleted).toEqual(type);
		expect(service.types).not.toContain(type);
		expect(service.types.length).toEqual(initialLength - 1);
	});
	it("deletes Types by ID", async () => {
		const initialLength = service.types.length;
		const type = service.types[initialLength - 1];
		const deleted = await service.deleteType(type.id);
		expect(deleted).toEqual(type);
		expect(service.types).not.toContain(type);
		expect(service.types.length).toEqual(initialLength - 1);
	});
	it("throws an error when asked to delete a Type that doesn't exist", async () => {
		const initialLength = service.types.length;
		await expectAsync(service.deleteType(-1)).toBeRejected();
		expect(service.types.length).toBe(initialLength);
	});
});
