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
import { HttpClient } from "@angular/common/http";
import { HttpClientTestingModule, HttpTestingController } from "@angular/common/http/testing";
import { TestBed } from "@angular/core/testing";
import { throwError } from "rxjs";
import { type Alert, AlertLevel } from "trafficops-types";

import { AlertService } from "../shared/alert/alert.service";

import { MiscAPIsService } from "./misc-apis.service";

const body = {
	dhcp: "yes" as const,
	disk: "sda",
	domainName: "domain-name",
	hostName: "host-name",
	interfaceMtu: 0,
	osVersionDir: "centos7",
	rootPass: "",
};

describe("MiscAPIsService", () => {
	let service: MiscAPIsService;
	let httpTestingController: HttpTestingController;
	let alert: Alert | null;

	beforeEach(() => {
		TestBed.configureTestingModule({
			imports: [HttpClientTestingModule],
			providers: [
				MiscAPIsService,
				{provide: AlertService, useValue: {
					newAlert: (a: Alert): void => {
						alert = a;
					}
				}}
			]
		});
		service = TestBed.inject(MiscAPIsService);
		httpTestingController = TestBed.inject(HttpTestingController);
		alert = null;
	});

	it("should be created", () => {
		expect(service).toBeTruthy();
	});

	it("sends requests for OS Versions", async () => {
		const responseP = service.getISOOSVersions();
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/osversions`);
		expect(req.request.method).toBe("GET");
		const data = {
			response: {
				// eslint-disable-next-line @typescript-eslint/naming-convention
				"CentOS 7": "centos7",
				// eslint-disable-next-line @typescript-eslint/naming-convention
				"Rocky Linux 8": "rocky8"
			}
		};
		req.flush(data);
		await expectAsync(responseP).toBeResolvedTo(data.response);
	});

	it("sends requests for ISO generation blobs", async () => {
		const responseP = service.generateISO(body);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/isos`);
		expect(req.request.method).toBe("POST");
		expect(req.request.body).toEqual(body);
		const response = new Blob();
		req.flush(response);
		await expectAsync(responseP).toBeResolvedTo(response);
	});

	it("throws an error when TO gives back an empty ISO", async () => {
		const responseP = service.generateISO(body);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/isos`);
		expect(req.request.method).toBe("POST");
		expect(req.request.body).toEqual(body);
		req.flush(null);
		await expectAsync(responseP).toBeRejected();
	});

	it("parses JSON-encoded error alerts when TO responds with an error", async () => {
		expect(alert).toBeNull();
		const responseP = service.generateISO(body);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/isos`);
		expect(req.request.method).toBe("POST");
		expect(req.request.body).toEqual(body);
		const errAlert = {
			level: AlertLevel.ERROR,
			text: "something wicked happened"
		};
		req.flush(new Blob([JSON.stringify({alerts: [errAlert]})]), {status: 500, statusText: "Internal Server Error"});
		await expectAsync(responseP).toBeRejectedWithError("POST isos failed with status 500 Internal Server Error");
		expect(alert).toEqual(errAlert);
	});

	it("handles invalid JSON body error responses", async () => {
		expect(alert).toBeNull();
		const responseP = service.generateISO(body);
		const req = httpTestingController.expectOne(`/api/${service.apiVersion}/isos`);
		expect(req.request.method).toBe("POST");
		expect(req.request.body).toEqual(body);
		req.flush(new Blob(['{"this": "json is" invalid}']), {status: 500, statusText: "Internal Server Error"});
		await expectAsync(responseP).toBeRejectedWithError("POST isos failed with status 500 Internal Server Error");
		expect(alert).toBeNull();
	});

	it("handles non-HTTP errors", async () => {
		const httpService = TestBed.inject(HttpClient);
		const spy = spyOn(httpService, "request").and.returnValue(throwError(new Error("something wicked happened")));
		expect(alert).toBeNull();
		const responseP = service.generateISO(body);
		await expectAsync(responseP).toBeRejectedWithError(/^POST isos failed: unknown error occurred:/);
		expect(spy).toHaveBeenCalled();
		expect(alert).toBeNull();
	});

	afterEach(() => {
		httpTestingController.verify();
	});
});
