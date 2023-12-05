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
import {RequestOrigin, ResponseOrigin} from "trafficops-types/dist/origin";

import { APIService } from "./base-api.service";

/** The allowed values for the 'useInTables' query parameter of GET requests to /origins. */
type UseInTable = "cachegroup" |
"server" |
"deliveryservice" |
"to_extension" |
"federation_resolver" |
"regex" |
"staticdnsentry" |
"steering_target";

/**
 * OriginService exposes API functionality relating to Origins.
 */
@Injectable()
export class OriginService extends APIService {
	/**
	 * Gets a specific Origin from Traffic Ops.
	 *
	 * @param idOrName Either the integral, unique identifier (number) or name
	 * (string) of the Origin to be returned.
	 * @returns The requested Origin.
	 */
	public async getOrigins(idOrName: number | string): Promise<ResponseOrigin>;
	/**
	 * Gets Origins from Traffic Ops.
	 *
	 * @param idOrName Either the integral, unique identifier (number) or name
	 * (string) of a single Origin to be returned.
	 * @returns The requested Origin(s).
	 */
	public async getOrigins(): Promise<Array<ResponseOrigin>>;
	/**
	 * Gets one or all Origins from Traffic Ops.
	 *
	 * @param idOrName Optionally the integral, unique identifier (number) or
	 * name (string) of a single Origin to be returned.
	 * @returns The requested Origin(s).
	 */
	public async getOrigins(idOrName?: number | string): Promise<ResponseOrigin | Array<ResponseOrigin>> {
		const path = "origins";
		if (idOrName !== undefined) {
			let params;
			switch (typeof idOrName) {
				case "string":
					params = {name: idOrName};
					break;
				case "number":
					params = {id: idOrName};
			}
			const r = await this.get<[ResponseOrigin]>(path, undefined, params).toPromise();
			if (r.length !== 1) {
				throw new Error(`Traffic Ops responded with ${r.length} Origins by identifier ${idOrName}`);
			}
			return r[0];
		}
		return this.get<Array<ResponseOrigin>>(path).toPromise();
	}

	/**
	 * Gets all Origins used by specific database table.
	 *
	 * @param useInTable The database table for which to retrieve Origins.
	 * @returns The requested Origins.
	 */
	public async getOriginsInTable(useInTable: UseInTable): Promise<Array<ResponseOrigin>> {
		return this.get<Array<ResponseOrigin>>("types", undefined, {useInTable}).toPromise();
	}

	/**
	 * Deletes an existing origin.
	 *
	 * @param typeOrId Id of the origin to delete.
	 * @returns The deleted origin.
	 */
	public async deleteOrigin(typeOrId: number | ResponseOrigin): Promise<ResponseOrigin> {
		const id = typeof(typeOrId) === "number" ? typeOrId : typeOrId.id;
		return this.delete<ResponseOrigin>(`origins/${id}`).toPromise();
	}

	/**
	 * Creates a new origin.
	 *
	 * @param origin The origin to create.
	 * @returns The created origin.
	 */
	public async createOrigin(origin: RequestOrigin): Promise<ResponseOrigin> {
		return this.post<ResponseOrigin>("origins", origin).toPromise();
	}

	/**
	 * Replaces the current definition of a origin with the one given.
	 *
	 * @param origin The new origin.
	 * @returns The updated origin.
	 */
	public async updateOrigin(origin: ResponseOrigin): Promise<ResponseOrigin> {
		const path = `origins/${origin.id}`;
		return this.put<ResponseOrigin>(path, origin).toPromise();
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
