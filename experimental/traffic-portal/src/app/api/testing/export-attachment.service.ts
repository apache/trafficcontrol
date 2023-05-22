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

import { Injectable } from "@angular/core";
import { ProfileExport, ProfileType } from "trafficops-types";

/**
 * ExportAttachmentService exposes API functionality related to export of profile data as attachment.
 */
@Injectable()
export class ExportAttachmentService {
	private readonly profileExport: ProfileExport = {
		alerts: null,
		parameters:[],
		profile: {
			cdn: "ALL",
			description: "test",
			name: "TRAFFIC_ANALYTICS",
			type: ProfileType.TS_PROFILE
		},
	};

	/**
	 * Export Profile object from the API.
	 *
	 * @param id Specify unique identifier (number) of a specific Profile to retrieve the export object.
	 * @returns The requested Profile as attachment.
	 */
	public async exportProfile(id?: number): Promise<ProfileExport> {
		if( id !== undefined){
			const exportProfile = this.profileExport;
			return exportProfile;
		}
		return this.profileExport;
	}

}
