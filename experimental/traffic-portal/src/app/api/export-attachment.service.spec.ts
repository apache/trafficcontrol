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

import { ExportAttachmentService } from "./export-attachment.service";

describe("ExportAttachmentService", () => {
	let service: ExportAttachmentService;
	let httpTestingController: HttpTestingController;
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
				ExportAttachmentService,
			]
		});
		service = TestBed.inject(ExportAttachmentService);
		httpTestingController = TestBed.inject(HttpTestingController);
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
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
});
