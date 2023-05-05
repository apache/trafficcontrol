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
import type { RequestPhysicalLocation, ResponsePhysicalLocation } from "trafficops-types";

import { APIService } from "./base-api.service";

/**
 * PhysicalLocationService exposes API functionality relating to PhysicalLocations.
 */
@Injectable()
export class PhysicalLocationService extends APIService {
	/**
	 * Gets Physical Locations from Traffic Ops.
	 *
	 * @returns An Array of Physical Locations.
	 */
	public async getPhysicalLocations(): Promise<Array<ResponsePhysicalLocation>>;
	/**
	 * Gets a single Physical Location from Traffic Ops.
	 *
	 * @param nameOrID Specifies the name or ID of the Physical Location to
	 * fetch.
	 * @returns The requested Physical Location.
	 */
	public async getPhysicalLocations(nameOrID: string | number): Promise<ResponsePhysicalLocation>;
	/**
	 * Gets Physical Locations from Traffic Ops.
	 *
	 * @param nameOrID If given, returns only the PhysicalLocation with the given name
	 * (string) or ID (number).
	 * @returns An Array of PhysicalLocation objects - or a single PhysicalLocation object if 'nameOrID'
	 * was given.
	 */
	public async getPhysicalLocations(nameOrID?: string | number): Promise<Array<ResponsePhysicalLocation> | ResponsePhysicalLocation> {
		const path = "phys_locations";
		if(nameOrID) {
			let params;
			switch (typeof nameOrID) {
				case "string":
					params = {name: nameOrID};
					break;
				case "number":
					params = {id: nameOrID};
			}
			const r = await this.get<[ResponsePhysicalLocation]>(path, undefined, params).toPromise();
			return r[0];

		}
		return this.get<Array<ResponsePhysicalLocation>>(path).toPromise();
	}

	/**
	 * Replaces the current definition of a Physical Location with the one given.
	 *
	 * @param physicalLocation The new Physical Location.
	 * @returns The updated Physical Location.
	 */
	public async updatePhysicalLocation(physicalLocation: ResponsePhysicalLocation): Promise<ResponsePhysicalLocation> {
		const path = `phys_locations/${physicalLocation.id}`;
		return this.put<ResponsePhysicalLocation>(path, physicalLocation).toPromise();
	}

	/**
	 * Creates a new Physical Location.
	 *
	 * @param physicalLocation The Physical Location to create.
	 * @returns The created Physical Location.
	 */
	public async createPhysicalLocation(physicalLocation: RequestPhysicalLocation): Promise<ResponsePhysicalLocation> {
		return this.post<ResponsePhysicalLocation>("phys_locations", physicalLocation).toPromise();
	}

	/**
	 * Deletes an existing Physical Location.
	 *
	 * @param physLoc The Physical Location to be deleted (or its ID)
	 */
	public async deletePhysicalLocation(physLoc: ResponsePhysicalLocation | number): Promise<void> {
		const id = typeof(physLoc) === "number" ? physLoc : physLoc.id;
		return this.delete(`phys_locations/${id}`).toPromise();
	}

	constructor(http: HttpClient) {
		super(http);
	}
}
