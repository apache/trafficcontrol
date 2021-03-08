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

import { Observable } from "rxjs";
import { map } from "rxjs/operators";

import { Parameter, Profile } from "../../models";

import { APIService } from "./apiservice";

/**
 * Shared mapping function for converting Parameter 'lastUpdated' fields to actual dates.
 *
 * @param p The Parameter being converted.
 * @returns the converted parameter.
 */
function paramMap(p: Parameter): Parameter {
	if (p.lastUpdated) {
		p.lastUpdated = new Date((p.lastUpdated as unknown as string).replace("+00", "Z"));
	}
	return p;
}

/**
 * ServerService exposes API functionality related to Servers.
 */
@Injectable({providedIn: "root"})
export class ProfileService extends APIService {

	/**
	 * Injects the Angular HTTP client service into the parent constructor.
	 *
	 * @param http The Angular HTTP client service.
	 */
	constructor(http: HttpClient) {
		super(http);
	}

	public getProfiles(idOrName: number | string): Observable<Profile>;
	public getProfiles(): Observable<Array<Profile>>;
	/**
	 * Retrieves Profiles from the API.
	 *
	 * @param idOrName Specify either the integral, unique identifier (number) of a specific Profile to retrieve, or its name (string).
	 * @returns An Observable that will emit the requested Profile(s).
	 */
	public getProfiles(idOrName?: number | string): Observable<Array<Profile> | Profile> {
		const path = "profiles";
		if (idOrName !== undefined) {
			let params;
			switch (typeof idOrName) {
				case "number":
					params = {id: String(idOrName)};
					break;
				case "string":
					params = {name: idOrName};
			}
			return this.get<[Profile]>(path, undefined, params).pipe(map(
				r => {
					const p = r[0];
					p.lastUpdated = new Date((p.lastUpdated as unknown as string).replace("+00", "Z"));
					if (p.params) {
						p.params.map(paramMap);
					}
					return p;
				}
			));
		}
		return this.get<Array<Profile>>(path).pipe(map(
			r => r.map(
				profile => {
					profile.lastUpdated = new Date(profile.lastUpdated as unknown as string);
					if (profile.params) {
						profile.params.map(paramMap);
					}
					return profile;
				}
			)
		));
	}
}
