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
import { ProfileType, ResponseParameter, ResponseProfile } from "trafficops-types";

import { ProfileService } from "./profile.service";

import { APITestingModule } from ".";

describe("TestingProfileService", () => {
	let service: ProfileService;
	let profile: ResponseProfile;
	let param: ResponseParameter;

	beforeEach(async () => {
		TestBed.configureTestingModule({
			imports: [APITestingModule],
			providers: [
				ProfileService,
			]
		});
		service = TestBed.inject(ProfileService);
		expect((await service.getProfiles()).length).toBeGreaterThanOrEqual(2);
		expect((await service.getParameters()).length).toBeGreaterThanOrEqual(2);
		profile = service.profiles[1];
		param = service.parameters[1];
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	it("gets a single Profile by ID", async () => {
		await expectAsync(service.getProfiles(profile.id)).toBeResolvedTo(profile);
	});

	it("gets a single Profile by name", async () => {
		await expectAsync(service.getProfiles(profile.name)).toBeResolvedTo(profile);
	});

	it("throws an error when the requested Profile doesn't exist", async () => {
		await expectAsync(service.getProfiles(-1)).toBeRejected();
	});

	it("creates new Profiles", async () => {
		const newProfile = {
			cdn: 1,
			description: "",
			name: `test-${(new Date()).valueOf()}`,
			routingDisabled: false,
			type: ProfileType.ATS_PROFILE
		};
		const initialLength = service.profiles.length;
		const created = await service.createProfile(newProfile);
		expect(service.profiles.length).toEqual(initialLength+1);
		expect(service.profiles).toContain(created);
	});

	it("updates existing Profiles", async () => {
		const current = {...service.profiles[1]};
		current.name += String((new Date()).valueOf());
		const {name} = current;
		const updated = await service.updateProfile(current);
		expect(service.profiles[1]).toBe(updated);
		expect(updated.name).toBe(name);
	});

	it("throws an error when asked to update a non-existent Profile", async () => {
		await expectAsync(service.updateProfile({...service.profiles[1], id: -1})).toBeRejected();
	});

	it("deletes existing Profiles", async () => {
		const initialLength = service.profiles.length;
		const deleted = await service.deleteProfile(profile);
		expect(deleted).toEqual(profile);
		expect(service.profiles.length).toEqual(initialLength-1);
		expect(service.profiles).not.toContain(profile);
	});

	it("deletes an existing Profile by ID", async () => {
		const initialLength = service.profiles.length;
		const deleted = await service.deleteProfile(profile.id);
		expect(deleted).toEqual(profile);
		expect(service.profiles.length).toEqual(initialLength-1);
		expect(service.profiles).not.toContain(profile);
	});

	it("throws an error when asked to delete a non-existent Profile", async () => {
		await expectAsync(service.deleteProfile(-1)).toBeRejected();
	});

	it("pretends to import Profiles", async () => {
		const importable = {
			cdn: profile.cdnName ?? "",
			description: profile.description,
			name: profile.name,
			type: profile.type
		};
		const imported = await service.importProfile({parameters: [], profile: importable});
		expect(imported).toEqual({...importable, id: imported.id});
	});

	it("gets Parameters", async () => {
		await expectAsync(service.getParameters()).toBeResolvedTo(service.parameters);
	});
	it("gets a specific Parameter by ID", async () => {
		await expectAsync(service.getParameters(param.id)).toBeResolvedTo(param);
	});
	it("throws an error when a non-existent Parameter is requested", async () => {
		await expectAsync(service.getParameters(-1)).toBeRejected();
	});

	it("deletes existing Parameters", async () => {
		await expectAsync(service.deleteParameter(param)).toBeResolved();
		expect(service.parameters).not.toContain(param);
	});
	it("deletes existing Parameters by ID", async () => {
		await expectAsync(service.deleteParameter(param.id)).toBeResolved();
		expect(service.parameters).not.toContain(param);
	});
	it("throws an error when asked to delete a non-existent Parameter", async () => {
		await expectAsync(service.deleteParameter(-1)).toBeRejected();
	});

	it("creates new Parameters", async () => {
		const initialLength = service.parameters.length;
		const created = await service.createParameter({
			configFile: "conf.ig",
			name: "creation test",
			secure: false
		});
		expect(service.parameters.length).toEqual(initialLength+1);
		expect(service.parameters).toContain(created);
	});

	it("updates existing Parameters", async () => {
		const current = {
			...param,
			value: `${param.value} update test`
		};

		const updated = await service.updateParameter(current);
		expect(updated).toEqual(current);
		expect(service.parameters).not.toContain(param);
	});
	it("throws an error if asked to update a non-existent Parameter", async () => {
		await expectAsync(service.updateParameter({...param, id: -1})).toBeRejected();
	});

	it("gets Profiles assigned to a Parameter", async () => {
		const prof = service.profiles.find(p => p.params && p.params.length > 0);
		if (!prof) {
			return fail("No Profiles exist that have any Parameters - cannot test");
		}
		// This is safe because the above `find` check ensures that params is
		// neither null nor empty.
		const p = {...prof.params?.[0], lastUpdated: new Date()} as ResponseParameter;
		const profiles = await service.getProfilesByParam(p);
		expect(profiles.length).toBeGreaterThanOrEqual(1);
		expect(profiles).toContain(prof);
	});

	it("throws an error if asked to retrieve Profiles assigned to a non-existent Parameter", async () => {
		await expectAsync(service.getProfilesByParam(-1)).toBeRejected();
	});
});
