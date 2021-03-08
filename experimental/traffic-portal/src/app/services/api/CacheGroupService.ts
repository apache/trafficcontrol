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

import { CacheGroup } from "src/app/models";
import { APIService } from "./apiservice";


/**
 * CDNService expose API functionality relating to CDNs.
 */
@Injectable({providedIn: "root"})
export class CacheGroupService extends APIService {
	public getCacheGroups(idOrName: number | string): Observable<CacheGroup>;
	public getCacheGroups(): Observable<Array<CacheGroup>>;
	/**
	 * Gets one or all CDNs from Traffic Ops
	 *
	 * @param idOrName Optionally either the name or integral, unique identifier of a single Cache Group to be returned.
	 * @returns An Observable that will emit either an Array of CacheGroup objects, or a single CacheGroup, depending on whether
	 * `idOrName` was 	passed.
	 * @throws {Error} In the event that `idOrName` is passed but does not match any CacheGroup.
	 */
	public getCacheGroups(idOrName?: number | string): Observable<Array<CacheGroup> | CacheGroup> {
		const path = "cachegroups";
		if (idOrName !== undefined) {
			let params;
			switch (typeof(idOrName)) {
				case "string":
					params = {name: idOrName};
					break;
				case "number":
					params = {id: String(idOrName)};
			}
			return this.get<[CacheGroup]>(path, undefined, params).pipe(map(
				r => {
					const cg = r[0];
					if (cg.id !== idOrName) {
						throw new Error(`Traffic Ops returned no match for ID ${idOrName}`);
					}
					//  lastUpdated comes in as a string
					cg.lastUpdated = cg.lastUpdated ? new Date((cg.lastUpdated as unknown as string).replace("+00", "Z")) : undefined;
					return cg;
				}
			));
		}
		return this.get<Array<CacheGroup>>(path).pipe(map(r => r.map(
			cg => {
				if (cg.lastUpdated) {
					cg.lastUpdated = new Date((cg.lastUpdated as unknown as string).replace("+00", "Z"));
				}
				return cg;
			}
		)));
	}

	/**
	 * Injects the Angular HTTP client service into the parent constructor.
	 *
	 * @param http The Angular HTTP client service.
	 */
	constructor(http: HttpClient) {
		super(http);
	}
}
