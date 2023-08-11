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
import { HttpClient, HttpErrorResponse } from "@angular/common/http";
import { Injectable } from "@angular/core";
import type { ISORequest, OSVersions } from "trafficops-types";

import { AlertService } from "../shared/alert/alert.service";
import { LoggingService } from "../shared/logging.service";

import { APIService, hasAlerts } from "./base-api.service";

/**
 * This service implements APIs that aren't specific to any given ATC object.
 * They can provide things like system information, access to the wacky external
 * tooling provided by TO, basically anything that doesn't fit in a different
 * API service.
 */
@Injectable()
export class MiscAPIsService extends APIService{

	constructor(http: HttpClient, private readonly alertsService: AlertService, private readonly log: LoggingService) {
		super(http);
	}

	/**
	 * Retrieves the operating system versions that can be used to generate
	 * system images through the Traffic Ops API.
	 *
	 * @returns A mapping of human-friendly operating system names to
	 * machine-readable OS IDs that can be used in subsequent requests to
	 * {@link MiscAPIsService.generateISO}.
	 */
	public async getISOOSVersions(): Promise<OSVersions> {
		return this.get<OSVersions>("osversions").toPromise();
	}

	/**
	 * Generates a system image.
	 *
	 * @param spec The specifications used to define what kind of image Traffic
	 * Ops will generate.
	 * @returns The generated system image.
	 */
	public async generateISO(spec: ISORequest): Promise<Blob> {
		const options = {
			body: spec,
			...this.defaultOptions,
			responseType: "blob" as const
		};
		let response;
		try {
			response = await this.http.request("post", `/api/${this.apiVersion}/isos`, options).toPromise();
		} catch (e) {
			if (e instanceof HttpErrorResponse) {
				try {
					const body = JSON.parse(await e.error.text());
					if (hasAlerts(body)) {
						body.alerts.forEach(a => this.alertsService.newAlert(a));
					}
				} catch (innerError) {
					this.log.error("during handling request failure, encountered an error trying to parse error-level alerts:", innerError);
				}
				throw new Error(`POST isos failed with status ${e.status} ${e.statusText}`);
			}
			throw new Error(`POST isos failed: unknown error occurred: ${e}`);
		}
		if (!response.body) {
			throw new Error(`POST isos returned no response body - ${response.status} ${response.statusText}`);
		}
		if (response.body.type !== "application/octet-stream") {
			this.log.warn("data returned by TO for ISO generation is of unrecognized MIME type", response.body.type);
		}
		return response.body;
	}
}
