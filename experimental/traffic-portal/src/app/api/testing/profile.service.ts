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
import {
	type ProfileImport,
	type ProfileImportResponse,
	ProfileType,
	type RequestProfile,
	type ResponseProfile,
	type RequestParameter,
	type ResponseParameter,
} from "trafficops-types";

/**
 * ProfileService exposes API functionality related to Profiles.
 */
@Injectable()
export class ProfileService {
	private lastID = 10;
	private readonly profiles: ResponseProfile[] = [
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

	/**
	 * Retrieves all Profiles.
	 *
	 * @returns The requested Profiles.
	 */
	public async getProfiles(): Promise<Array<ResponseProfile>>;
	/**
	 * Retrieves a specific Profile from the API.
	 *
	 * @param idOrName Specify either the integral, unique identifier (number)
	 * of a specific Profile to retrieve, or its name (string).
	 * @returns The requested Profile.
	 */
	public async getProfiles(idOrName: number | string): Promise<ResponseProfile>;
	/**
	 * Retrieves one or allProfiles from the API.
	 *
	 * @param idOrName If given, only the Profile with this integral, unique
	 * identifier (number) or name (string) will be returned.
	 * @returns The requested Profile(s).
	 */
	public async getProfiles(idOrName?: number | string): Promise<Array<ResponseProfile> | ResponseProfile> {
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

	/**
	 * Creates a new Profile.
	 *
	 * @param profile The Profile to create.
	 * @returns The created Profile.
	 */
	public async createProfile(profile: RequestProfile): Promise<ResponseProfile> {
		const t = {
			...profile,
			cdnName: null,
			id: ++this.lastID,
			lastUpdated: new Date()
		};
		this.profiles.push(t);
		return t;
	}

	/**
	 * Updates an existing Profile.
	 *
	 * @param profile the Profile to update.
	 * @returns The success message.
	 */
	public async updateProfile(profile: ResponseProfile): Promise<ResponseProfile> {
		const id = this.profiles.findIndex(d => d.id === profile.id);
		if (id === -1) {
			throw new Error(`no such profile: ${profile.id}`);
		}
		this.profiles[id] = profile;
		return profile;
	}

	/**
	 * Deletes an existing Profile.
	 *
	 * @param profile The Profile to delete, or just its ID.
	 * @returns The deleted Profile.
	 */
	public async deleteProfile(profile: number | ResponseProfile): Promise<ResponseProfile> {
		const id = typeof(profile) === "number" ? profile : profile.id;
		const index = this.profiles.findIndex(t => t.id === id);
		if (index === -1) {
			throw new Error(`no such Profile #${id}`);
		}
		return this.profiles.splice(index, 1)[0];
	}

	/**
	 * Imports a Profile along with all its associated Parameters.
	 *
	 * @param profile The specification of the Profile to be imported/created.
	 * @returns The created Profile.
	 */
	public async importProfile(profile: ProfileImport): Promise<ProfileImportResponse> {
		const t = {
			...profile.profile,
			id: ++this.lastID,
		};
		return t;
	}

	private lastParamID = 20;
	private readonly parameters:  ResponseParameter[] = [
		{
			configFile: "cfg.txt",
			id: 1,
			lastUpdated: new Date(),
			name: "param1",
			profiles: [],
			secure: false,
			value: "10"
		}
	];

	/**
	 * Gets all Parameters.
	 *
	 * @returns The requested Parameters.
	 */
	public async getParameters(): Promise<Array<ResponseParameter>>;
	/**
	 * Gets a specific Parameter.
	 *
	 * @param id The integral, unique identifier of the specific Parameter to be
	 * returned.
	 * @returns The requested parameter.
	 */
	public async getParameters(id: number): Promise<ResponseParameter>;
	/**
	 * Gets one or all Parameters.
	 *
	 * @param id If given, only the Parameter with this unique, integral
	 * identifier will be returned.
	 * @returns The requested Parameter(s).
	 */
	public async getParameters(id?: number): Promise<ResponseParameter | Array<ResponseParameter>> {
		if (id !== undefined) {
			const parameter = this.parameters.filter(t=>t.id === id)[0];
			if (!parameter) {
				throw new Error(`no such Parameter: ${id}`);
			}
			return parameter;
		}
		return this.parameters;
	}

	/**
	 * Deletes a Parameter.
	 *
	 * @param typeOrId The Parameter to be deleted, or just its ID.
	 */
	public async deleteParameter(typeOrId: number | ResponseParameter): Promise<void> {
		const id = typeof typeOrId === "number" ? typeOrId : typeOrId.id;
		const idx = this.parameters.findIndex(p => p.id === id);
		if (idx < 0) {
			throw new Error(`no such Parameter: #${id}`);
		}
		this.parameters.splice(idx, 1);
	}

	/**
	 * Creates a new Parameter.
	 *
	 * @param parameter The Parameter to create.
	 * @returns The created Parameter.
	 */
	public async createParameter(parameter: RequestParameter): Promise<ResponseParameter> {
		const t = {
			...parameter,
			id: ++this.lastParamID,
			lastUpdated: new Date(),
			profiles: [],
			value: parameter.value ?? ""
		};
		this.parameters.push(t);
		return t;
	}

	/**
	 * Replaces an existing Parameter with the provided new definition of a
	 * Parameter.
	 *
	 * @param parameter The full new definition of the Parameter being
	 * updated.
	 * @returns The updated Parameter
	 */
	public async updateParameter(parameter: ResponseParameter): Promise<ResponseParameter> {
		const id = this.parameters.findIndex(d => d.id === parameter.id);
		if (id === -1) {
			throw new Error(`no such parameter: ${parameter.id}`);
		}
		this.parameters[id] = parameter;
		return parameter;
	}

	/**
	 * Retrieves Profiles associated with a Parameter from the API.
	 *
	 * @param parameter Either a {@link ResponseParameter} or an integral,
	 * unique identifier of a Parameter, for which the Profiles are to be
	 * retrieved.
	 * @returns The requested Profile(s).
	 */
	public async getProfilesByParam(parameter: number| ResponseParameter): Promise<Array<ResponseProfile>> {
		const id = typeof parameter === "number" ? parameter : parameter.id;
		if (id === -1) {
			throw new Error(`no such parameter: ${id}`);
		}
		const index = this.parameters.findIndex(d => d.id === id);
		const profiles = this.parameters[index].profiles;
		if (profiles === null) {
			return new Array<ResponseProfile>();
		}
		const returnedProfiles = new Array<ResponseProfile>();
		for (const val of profiles) {
			const p = this.getProfiles(val);
			returnedProfiles.push(await p);
		}
		return returnedProfiles;
	}
}
