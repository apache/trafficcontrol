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

import type { PhysicalLocation } from "src/app/models";

/**
 * PhysicalLocationService exposes API functionality relating to PhysicalLocations.
 */
@Injectable()
export class PhysicalLocationService {
	private readonly locs = [
		{
			address: "1600 Pennsylvania Avenue NW",
			city: "Washington",
			comments: "",
			email: "",
			id: 1,
			lastUpdated: new Date(),
			name: "test",
			phone: "",
			poc: "",
			region: "Washington, D.C",
			regionId: 1,
			shortName: "test",
			state: "DC",
			zip: "20500"
		}
	];

	public async getPhysicalLocations(idOrName: number | string): Promise<PhysicalLocation>;
	public async getPhysicalLocations(): Promise<Array<PhysicalLocation>>;
	/**
	 * Gets one or all PhysicalLocations from Traffic Ops
	 *
	 * @param idOrName Either the integral, unique identifier (number) or name (string) of a single PhysicalLocation to be returned.
	 * @returns The requested PhysicalLocation(s).
	 */
	public async getPhysicalLocations(idOrName?: number | string): Promise<PhysicalLocation | Array<PhysicalLocation>> {
		if (idOrName !== undefined) {
			let loc;
			switch (typeof idOrName) {
				case "string":
					loc = this.locs.find(l=>l.name === idOrName);
					break;
				case "number":
					loc = this.locs.find(l=>l.id === idOrName);
			}
			if (!loc) {
				throw new Error(`no such Physical Location: ${idOrName}`);
			}
			return loc;
		}
		return this.locs;
	}
}
