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

import { Type } from "../../models";
import { APIService } from "./apiservice";

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
@Injectable({providedIn: "root"})
export class TypeService extends APIService {
	public getTypes(idOrName: number | string): Observable<Type>;
	public getTypes(): Observable<Array<Type>>;
	/**
	 * Gets one or all Types from Traffic Ops
	 *
	 * @param idOrName Either the integral, unique identifier (number) or name (string) of a single Type to be returned.
	 * @returns An Observable that will emit the requested Type(s).
	 */
	public getTypes(idOrName?: number | string): Observable<Type | Array<Type>> {
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
			return this.get<[Type]>(path, undefined, params).pipe(map(
				r => r[0]
			));
		}
		return this.get<Array<Type>>(path);
	}

	/**
	 * Gets all Types used by specific database table.
	 *
	 * @param useInTable The database table for which to retrieve Types.
	 * @returns An Observable that emits the requested Types.
	 */
	public getTypesInTable(useInTable: UseInTable): Observable<Array<Type>> {
		return this.get<Array<Type>>("types", undefined, {useInTable});
	}

	/**
	 * Gets all Server Types.
	 *
	 * @returns All Types that have 'server' as their 'useInTable'.
	 */
	public getServerTypes(): Observable<Array<Type>> {
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
