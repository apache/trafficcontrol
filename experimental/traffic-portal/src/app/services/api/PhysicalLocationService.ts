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

import { PhysicalLocation } from "../../models";
import { APIService } from "./apiservice";

/**
 * PhysicalLocationService exposes API functionality relating to PhysicalLocations.
 */
@Injectable({providedIn: "root"})
export class PhysicalLocationService extends APIService {
	public async getPhysicalLocations(idOrName: number | string): Promise<PhysicalLocation>;
	public async getPhysicalLocations(): Promise<Array<PhysicalLocation>>;
	/**
	 * Gets one or all PhysicalLocations from Traffic Ops
	 *
	 * @param idOrName Either the integral, unique identifier (number) or name (string) of a single PhysicalLocation to be returned.
	 * @returns The requested PhysicalLocation(s).
	 */
	public async getPhysicalLocations(idOrName?: number | string): Promise<PhysicalLocation | Array<PhysicalLocation>> {
		const path = "phys_locations";
		let prom;
		if (idOrName !== undefined) {
			let params;
			switch (typeof idOrName) {
				case "string":
					params = {name: idOrName};
					break;
				case "number":
					params = {id: String(idOrName)};
			}
			prom = this.get<[PhysicalLocation]>(path, undefined, params).toPromise().then(
				r => r[0]
			);
		} else {
			prom = this.get<Array<PhysicalLocation>>(path).toPromise();
		}
		return prom;
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
