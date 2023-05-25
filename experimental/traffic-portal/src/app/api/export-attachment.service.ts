/*
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
import { Injectable } from "@angular/core";
import { ProfileExport, ResponseProfile } from "trafficops-types";

import { environment } from "src/environments/environment";

/**
 * Defines & handles api endpoints related to export attachments/json/text
 * Here we are not using the base-api-service to tap the response.
 */
@Injectable()
export class ExportAttachmentService {

	/**
	 * The API version used by the service(s) - this will be overridden by the
	 * environment if a different API version is therein found.
	 */
	public apiVersion = "4.0";

	constructor(private readonly http: HttpClient) {
		if (environment.apiVersion) {
			this.apiVersion = environment.apiVersion;
		}
	}

	/**
	 * Exports profile
	 *
	 * @param profileId Id of the profile to export.
	 * @returns profile export object.
	 */
	public async exportProfile(profileId: number | ResponseProfile): Promise<ProfileExport>{
		const id = typeof (profileId) === "number" ? profileId : profileId.id;
		return this.http.get<ProfileExport>(`/api/${this.apiVersion}/profiles/${id}/export`).toPromise();
	}
}
