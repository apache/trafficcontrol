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

import { Profile, ProfileType } from "src/app/models";

/**
 * ProfileService exposes API functionality related to Profiles.
 */
@Injectable()
export class ProfileService {

	private readonly profiles = [
		{
			cdn: 1,
			cdnName: "ALL",
			description: "Global Traffic Ops profile, DO NOT DELETE",
			id: 1,
			lastUpdated: new Date(),
			name: "GLOBAL",
			params: [
				{
					configFile: "global",
					id: 1,
					lastUpdated: null,
					name: "tm.instance_name",
					profiles: null,
					secure: false,
					value: "Traffic Ops CDN"
				},
				{
					configFile: "global",
					id: 2,
					lastUpdated: null,
					name: "tm.toolname",
					profiles: null,
					secure: false,
					value: "Traffic Ops"
				},
				{
					configFile: "regex_revalidate.config",
					id: 3,
					lastUpdated: null,
					name: "maxRevalDurationDays",
					profiles: null,
					secure: false,
					value: "90"
				},
				{
					configFile: "global",
					id: 4,
					lastUpdated: null,
					name: "tm.url",
					profiles: null,
					secure: false,
					value: "https://trafficops.infra.ciab.test:443/"
				},
				{
					configFile: "global",
					id: 5,
					lastUpdated: null,
					name: "tm.instance_name",
					profiles: null,
					secure: false,
					value: "CDN-In-A-Box"
				},
				{
					configFile: "CRConfig.json",
					id: 6,
					lastUpdated: null,
					name: "geolocation.polling.url",
					profiles: null,
					secure: false,
					value: "https://static.infra.ciab.test/GeoLite2-City.mmdb.gz"
				},
				{
					configFile: "CRConfig.json",
					id: 7,
					lastUpdated: null,
					name: "geolocation6.polling.url",
					profiles: null,
					secure: false,
					value: "https://static.infra.ciab.test/GeoLite2-City.mmdb.gz"
				},
				{
					configFile: "global",
					id: 8,
					lastUpdated: null,
					name: "use_reval_pending",
					profiles: null,
					secure: false,
					value: "1"
				},
				{
					configFile: "global",
					id: 9,
					lastUpdated: null,
					name: "default_geo_miss_latitude",
					profiles: null,
					secure: false,
					value: "0"
				},
				{
					configFile: "global",
					id: 10,
					lastUpdated: null,
					name: "default_geo_miss_longitude",
					profiles: null,
					secure: false,
					value: "-1"
				}
			],
			routingDisabled: false,
			type: ProfileType.UNK_PROFILE
		},
		{
			cdn: 2,
			cdnName: "test",
			description: "Edge Cache - Apache Traffic Server",
			id: 2,
			lastUpdated: new Date(),
			name: "EDGE_TIER_ATS_CACHE",
			routingDisabled: false,
			type: ProfileType.ATS_PROFILE
		}
	];

	public async getProfiles(idOrName: number | string): Promise<Profile>;
	public async getProfiles(): Promise<Array<Profile>>;
	/**
	 * Retrieves Profiles from the API.
	 *
	 * @param idOrName Specify either the integral, unique identifier (number) of a specific Profile to retrieve, or its name (string).
	 * @returns The requested Profile(s).
	 */
	public async getProfiles(idOrName?: number | string): Promise<Array<Profile> | Profile> {
		if (idOrName !== undefined) {
			let profile;
			switch (typeof idOrName) {
				case "number":
					profile = this.profiles.filter(p=>p.id === idOrName)[0];
					break;
				case "string":
					profile = this.profiles.filter(p=>p.name === idOrName)[0];
			}
			if (!profile) {
				throw new Error(`no such Profile: ${idOrName}`);
			}
			return profile;
		}

		return this.profiles.map(
			p => ({
				cdn: p.cdn,
				cdnName: p.cdnName,
				description: p.description,
				id: p.id,
				lastUpdated: p.lastUpdated,
				name: p.name,
				routingDisabled: p.routingDisabled,
				type: p.type
			})
		);
	}
}
