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

import type { Type } from "src/app/models";

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
	public async getTypes(idOrName: number | string): Promise<Type>;
	public async getTypes(): Promise<Array<Type>>;
	/**
	 * Gets one or all Types from Traffic Ops
	 *
	 * @param idOrName Either the integral, unique identifier (number) or name (string) of a single Type to be returned.
	 * @returns The requested Type(s).
	 */
	public async getTypes(idOrName?: number | string): Promise<Type | Array<Type>> {
		const path = "types";
		if (idOrName !== undefined) {
			let params;
			switch (typeof idOrName) {
				case "string":
					params = {name: idOrName};
					break;
				case "number":
					params = {id: String(idOrName)};
			}
			return this.get<[Type]>(path, undefined, params).toPromise().then(
				r => r[0]
			).catch(
				e => {
					console.error("Failed to get Type:", e);
					return {
						id: -1,
						name: ""
					};
				}
			);
		}
		return this.get<Array<Type>>(path).toPromise().catch(
			e => {
				console.error("Failed to get Types:", e);
				return [];
			}
		);
	}

	/**
	 * Gets all Types used by specific database table.
	 *
	 * @param useInTable The database table for which to retrieve Types.
	 * @returns The requested Types.
	 */
	public async getTypesInTable(useInTable: UseInTable): Promise<Array<Type>> {
		return this.get<Array<Type>>("types", undefined, {useInTable}).toPromise().catch(
			(e) => {
				console.error("Failed to get Types:", e);
				return [];
			}
		);
	}

	/**
	 * Gets all Server Types.
	 *
	 * @returns All Types that have 'server' as their 'useInTable'.
	 */
	public async getServerTypes(): Promise<Array<Type>> {
		return this.getTypesInTable("server");
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
