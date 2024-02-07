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
import type {
	ProfileImport,
	ProfileImportResponse,
	RequestParameter,
	RequestProfile,
	ResponseParameter,
	ResponseProfile,
} from "trafficops-types";

import { APIService } from "./base-api.service";

/**
 * ProfileService exposes API functionality related to Profiles.
 */
@Injectable()
export class ProfileService extends APIService {

	/**
	 * Injects the Angular HTTP client service into the parent constructor.
	 *
	 * @param http The Angular HTTP client service.
	 */
	constructor(http: HttpClient) {
		super(http);
	}

	/**
	 * Retrieves a single Profile from the API.
	 *
	 * @param idOrName Either the integral, unique identifier (number) or name
	 * (string) of the specific Profile to retrieve.
	 * @returns The requested Profile.
	 */
	public async getProfiles(idOrName: number | string): Promise<ResponseProfile>;
	/**
	 * Retrieves Profiles from the API.
	 *
	 * @returns The requested Profiles.
	 */
	public async getProfiles(): Promise<Array<ResponseProfile>>;
	/**
	 * Retrieves one or more Profiles from the API.
	 *
	 * @param idOrName Specify either the integral, unique identifier (number)
	 * of a specific Profile to retrieve, or its name (string).
	 * @returns The requested Profile(s).
	 */
	public async getProfiles(idOrName?: number | string): Promise<Array<ResponseProfile> | ResponseProfile> {
		const path = "profiles";
		if (idOrName !== undefined) {
			let params;
			switch (typeof idOrName) {
				case "number":
					params = {id: idOrName};
					break;
				case "string":
					params = {name: idOrName};
			}
			const r = await this.get<[ResponseProfile]>(path, undefined, params).toPromise();
			return r[0];
		}
		return this.get<Array<ResponseProfile>>(path).toPromise();
	}

	/**
	 * Retrieves Profiles associated with a Parameter from the API.
	 *
	 * @param p Either a {@link ResponseParameter} or an integral, unique identifier of a Parameter, for which the
	 * Profiles are to be retrieved.
	 * @returns The requested Profile(s).
	 */
	public async getProfilesByParam(p: number| ResponseParameter): Promise<Array<ResponseProfile>> {
		let id: number;
		if (typeof p === "number") {
			id = p;
		} else {
			id = p.id;
		}

		const path = "profiles";
		const params = {param: id};
		const r = await this.get<Array<ResponseProfile>>(path, undefined, params).toPromise();
		return r;
	}

	/**
	 * Creates a new profile.
	 *
	 * @param profile The profile to create.
	 * @returns The created profile.
	 */
	public async createProfile(profile: RequestProfile): Promise<ResponseProfile> {
		return this.post<ResponseProfile>("profiles", profile).toPromise();
	}

	/**
	 * Replaces the current definition of a profile with the one given.
	 *
	 * @param profile The new profile.
	 * @returns The updated profile.
	 */
	public async updateProfile(profile: ResponseProfile): Promise<ResponseProfile> {
		const path = `profiles/${profile.id}`;
		return this.put<ResponseProfile>(path, profile).toPromise();
	}

	/**
	 * Deletes an existing Profile.
	 *
	 * @param profile The Profile to delete, or just its ID.
	 * @returns The deleted Profile.
	 */
	public async deleteProfile(profile: number | ResponseProfile): Promise<ResponseProfile> {
		const id = typeof (profile) === "number" ? profile : profile.id;
		return this.delete<ResponseProfile>(`profiles/${id}`).toPromise();
	}

	/**
	 * Imports a Profile along with all its associated Parameters.
	 *
	 * @param importJSON The specification of the Profile to be imported/created.
	 * @returns The created Profile.
	 */
	public async importProfile(importJSON: ProfileImport): Promise<ProfileImportResponse>{
		return this.post<ProfileImportResponse>("profiles/import", importJSON).toPromise();
	}

	/**
	 * Retrieves all Parameters from Traffic Ops.
	 *
	 * @returns The requested Parameters.
	 */
	public async getParameters(): Promise<Array<ResponseParameter>>;
	/**
	 * Retrieves a single Parameter from Traffic Ops.
	 *
	 * @param id The integral, unique identifier of the specific Parameter to
	 * retrieve.
	 * @returns The requested Parameter(s).
	 */
	public async getParameters(id: number): Promise<ResponseParameter>;
	/**
	 * Retrieves a Parameter or Parameters from the API.
	 *
	 * @param id If given, only the Parameter with this integral, unique
	 * identifier will be returned.
	 * @returns The requested Parameter(s).
	 */
	public async getParameters(id?: number): Promise<Array<ResponseParameter> | ResponseParameter> {
		const path = "parameters";
		if (id !== undefined) {
			const params = {id};
			const r = await this.get<[ResponseParameter]>(path, undefined, params).toPromise();
			if (r.length !== 1) {
				throw new Error(`Traffic Ops responded with ${r.length} Parameters by identifier ${id}`);
			}
			return r[0];
		}
		return this.get<Array<ResponseParameter>>(path).toPromise();
	}

	/**
	 * Deletes an existing Parameter.
	 *
	 * @param typeOrId The ID of the Parameter to delete.
	 * @returns The deleted Parameter.
	 */
	public async deleteParameter(typeOrId: number | ResponseParameter): Promise<void> {
		const id = typeof(typeOrId) === "number" ? typeOrId : typeOrId.id;
		return this.delete(`parameters/${id}`).toPromise();
	}

	/**
	 * Creates a new Parameter.
	 *
	 * @param parameter The Parameter to create.
	 * @returns The created Parameter.
	 */
	public async createParameter(parameter: RequestParameter): Promise<ResponseParameter> {
		return this.post<ResponseParameter>("parameters", parameter).toPromise();
	}

	/**
	 * Replaces the current definition of a Parameter with the one given.
	 *
	 * @param parameter The new Parameter.
	 * @returns The updated Parameter.
	 */
	public async updateParameter(parameter: ResponseParameter): Promise<ResponseParameter> {
		const path = `parameters/${parameter.id}`;
		return this.put<ResponseParameter>(path, parameter).toPromise();
	}
}
