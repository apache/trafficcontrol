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
import { ProfileExport, ProfileType } from "trafficops-types";

import { ProfileService } from "./profile.service";

describe("ProfileService", () => {
	let service: ProfileService;
	let httpTestingController: HttpTestingController;
	const profile = {
		cdn: 1,
		cdnName: "CDN",
		description: "",
		id: 1,
		lastUpdated: new Date(),
		name: "TestQuest",
		routingDisabled: false,
		type: ProfileType.ATS_PROFILE
	};
	const importProfile = {
		parameters:[],
		profile: {
			cdn: "CDN",
			description: "",
			id: 1,
			name: "TestQuest",
			type: ProfileType.ATS_PROFILE,
		}
	};
	const exportProfile: ProfileExport = {
		alerts: null,
		parameters:[],
		profile: {
			cdn: "ALL",
			description: "test",
			name: "TRAFFIC_ANALYTICS",
			type: ProfileType.TS_PROFILE
		}
	};

	beforeEach(() => {
		TestBed.configureTestingModule({
			imports: [HttpClientTestingModule],
			providers: [
				ProfileService,
			]
		});
		service = TestBed.inject(ProfileService);
		httpTestingController = TestBed.inject(HttpTestingController);
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	it("sends requests multiple Profiles", async () => {
		const responseP = service.getProfiles();
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/profiles`);
		expect(req.request.method).toBe("GET");
		expect(req.request.params.keys().length).toBe(0);
		req.flush({response: [profile]});
		await expectAsync(responseP).toBeResolvedTo([profile]);
	});

	it("sends requests for a single Profile by ID", async () => {
		const responseP = service.getProfiles(profile.id);
		const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/profiles`);
		expect(req.request.method).toBe("GET");
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("id")).toBe(String(profile.id));
		req.flush({response: [profile]});
		await expectAsync(responseP).toBeResolvedTo(profile);
	});

	it("sends requests for a single Profile by name", async () => {
		const responseP = service.getProfiles(profile.name);
		const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/profiles`);
		expect(req.request.method).toBe("GET");
		expect(req.request.params.keys().length).toBe(1);
		expect(req.request.params.get("name")).toBe(profile.name);
		req.flush({response: [profile]});
		await expectAsync(responseP).toBeResolvedTo(profile);
	});

	it("creates new Profiles", async () => {
		const responseP = service.createProfile(profile);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/profiles`);
		expect(req.request.method).toBe("POST");
		expect(req.request.params.keys().length).toBe(0);
		expect(req.request.body).toBe(profile);
		req.flush({response: profile});
		await expectAsync(responseP).toBeResolvedTo(profile);
	});

	it("deletes existing Profiles", async () => {
		const responseP = service.deleteProfile(profile);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/profiles/${profile.id}`);
		expect(req.request.method).toBe("DELETE");
		expect(req.request.params.keys().length).toBe(0);
		expect(req.request.body).toBeNull();
		req.flush({response: profile});
		await expectAsync(responseP).toBeResolvedTo(profile);
	});

	it("deletes an existing Profile by ID", async () => {
		const responseP = service.deleteProfile(profile.id);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/profiles/${profile.id}`);
		expect(req.request.method).toBe("DELETE");
		expect(req.request.params.keys().length).toBe(0);
		expect(req.request.body).toBeNull();
		req.flush({response: profile});
		await expectAsync(responseP).toBeResolvedTo(profile);
	});

	it("sends request for Export object by Profile ID", async () => {
		const id = 1;
		const response = service.exportProfile(id);
		const req = httpTestingController.expectOne(r => r.url === `/api/${service.apiVersion}/profiles/${id}/export`);
		expect(req.request.method).toBe("GET");
		expect(req.request.params.keys().length).toBe(0);
		req.flush(exportProfile);
		await expectAsync(response).toBeResolvedTo(exportProfile);
	});

	it("send request for import profile", async () => {
		const responseP = service.importProfile(importProfile);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/profiles/import`);
		expect(req.request.method).toBe("POST");
		expect(req.request.params.keys().length).toBe(0);
		expect(req.request.body).toBe(importProfile);
		req.flush({response: importProfile.profile});
		await expectAsync(responseP).toBeResolvedTo(importProfile.profile);
	});

	afterEach(() => {
		httpTestingController.verify();
	});
});
