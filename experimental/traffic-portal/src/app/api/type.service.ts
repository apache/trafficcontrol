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
import type { RequestType, TypeFromResponse } from "trafficops-types";

import { APIService } from "./base-api.service";

/** The allowed values for the 'useInTables' query parameter of GET requests to /types. */
type UseInTable = "cachegroup" |
"server" |
"deliveryservice" |
"to_extension" |
"federation_resolver" |
"regex" |
"staticdnsentry" |
"steering_target";

/**
 * TypeService exposes API functionality relating to Types.
 */
@Injectable()
export class TypeService extends APIService {
	/**
	 * Gets all Types from Traffic Ops.
	 *
	 * @returns The requested Types.
	 */
	public async getTypes(): Promise<Array<TypeFromResponse>>;
	/**
	 * Gets a specific Type from Traffic Ops.
	 *
	 * @param idOrName Either the integral, unique identifier (number) or name
	 * (string) of the Type to be returned.
	 * @returns The requested Type.
	 */
	public async getTypes(idOrName: number | string): Promise<TypeFromResponse>;
	/**
	 * Gets one or all Types from Traffic Ops.
	 *
	 * @param idOrName Optionally the integral, unique identifier (number) or
	 * name (string) of a single Type to be returned.
	 * @returns The requested Type(s).
	 */
	public async getTypes(idOrName?: number | string): Promise<TypeFromResponse | Array<TypeFromResponse>> {
		const path = "types";
		if (idOrName !== undefined) {
			let params;
			switch (typeof idOrName) {
				case "string":
					params = {name: idOrName};
					break;
				case "number":
					params = {id: idOrName};
			}
			const r = await this.get<[TypeFromResponse]>(path, undefined, params).toPromise();
			if (r.length !== 1) {
				throw new Error(`Traffic Ops responded with ${r.length} Types by identifier ${idOrName}`);
			}
			return r[0];
		}
		return this.get<Array<TypeFromResponse>>(path).toPromise();
	}

	/**
	 * Gets all Types used by specific database table.
	 *
	 * @param useInTable The database table for which to retrieve Types.
	 * @returns The requested Types.
	 */
	public async getTypesInTable(useInTable: UseInTable): Promise<Array<TypeFromResponse>> {
		return this.get<Array<TypeFromResponse>>("types", undefined, {useInTable}).toPromise();
	}

	/**
	 * Gets all Server Types.
	 *
	 * @returns All Types that have 'server' as their 'useInTable'.
	 */
	public async getServerTypes(): Promise<Array<TypeFromResponse>> {
		return this.getTypesInTable("server");
	}

	/**
	 * Deletes an existing type.
	 *
	 * @param typeOrId Id of the type to delete.
	 * @returns The deleted type.
	 */
	public async deleteType(typeOrId: number | TypeFromResponse): Promise<TypeFromResponse> {
		const id = typeof(typeOrId) === "number" ? typeOrId : typeOrId.id;
		return this.delete<TypeFromResponse>(`types/${id}`).toPromise();
	}

	/**
	 * Creates a new type.
	 *
	 * @param type The type to create.
	 * @returns The created type.
	 */
	public async createType(type: RequestType): Promise<TypeFromResponse> {
		return this.post<TypeFromResponse>("types", type).toPromise();
	}

	/**
	 * Replaces the current definition of a type with the one given.
	 *
	 * @param type The new type.
	 * @returns The updated type.
	 */
	public async updateType(type: TypeFromResponse): Promise<TypeFromResponse> {
		const path = `types/${type.id}`;
		return this.put<TypeFromResponse>(path, type).toPromise();
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
