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
import { RequestPhysicalLocation, ResponsePhysicalLocation } from "trafficops-types";

/**
 * PhysicalLocationService exposes API functionality relating to PhysicalLocations.
 */
@Injectable()
export class PhysicalLocationService {
	private lastID = 1;
	private readonly physicalLocations: Array<ResponsePhysicalLocation> = [{
		address: "street",
		city: "city",
		comments: null,
		email: null,
		id: 1,
		lastUpdated: new Date(),
		name: "phys",
		phone: null,
		poc: null,
		region: "Region",
		regionId: 1,
		shortName: "short",
		state: "st",
		zip: "0000"
	}];

	public async getPhysicalLocations(): Promise<Array<ResponsePhysicalLocation>>;
	public async getPhysicalLocations(nameOrID: string | number): Promise<ResponsePhysicalLocation>;

	/**
	 * Gets an array of physicalLocations from Traffic Ops.
	 *
	 * @param nameOrID If given, returns only the ResponsePhysicalLocation with the given name
	 * (string) or ID (number).
	 * @returns An Array of ResponsePhysicalLocation objects - or a single ResponsePhysicalLocation object if 'nameOrID'
	 * was given.
	 */
	public async getPhysicalLocations(nameOrID?: string | number): Promise<Array<ResponsePhysicalLocation> | ResponsePhysicalLocation> {
		if(nameOrID) {
			let physicalLocation;
			switch (typeof nameOrID) {
				case "string":
					physicalLocation = this.physicalLocations.find(d=>d.name === nameOrID);
					break;
				case "number":
					physicalLocation = this.physicalLocations.find(d=>d.id === nameOrID);
			}
			if (!physicalLocation) {
				throw new Error(`no such PhysicalLocation: ${nameOrID}`);
			}
			return physicalLocation;
		}
		return this.physicalLocations;
	}

	/**
	 * Replaces the current definition of a physicalLocation with the one given.
	 *
	 * @param physicalLocation The new physicalLocation.
	 * @returns The updated physicalLocation.
	 */
	public async updatePhysicalLocation(physicalLocation: ResponsePhysicalLocation): Promise<ResponsePhysicalLocation> {
		const id = this.physicalLocations.findIndex(d => d.id === physicalLocation.id);
		if (id === -1) {
			throw new Error(`no such PhysicalLocation: ${physicalLocation.id}`);
		}
		this.physicalLocations[id] = physicalLocation;
		return physicalLocation;
	}

	/**
	 * Creates a new physicalLocation.
	 *
	 * @param physicalLocation The physicalLocation to create.
	 * @returns The created physicalLocation.
	 */
	public async createPhysicalLocation(physicalLocation: RequestPhysicalLocation): Promise<ResponsePhysicalLocation> {
		return {
			...physicalLocation,
			comments: physicalLocation.comments ?? null,
			email: physicalLocation.email ?? null,
			id: ++this.lastID,
			lastUpdated: new Date(),
			phone: physicalLocation.phone ?? null,
			poc: physicalLocation.poc ?? null,
			region: ""
		};
	}

	/**
	 * Deletes an existing physicalLocation.
	 *
	 * @param id Id of the physicalLocation to delete.
	 * @returns The deleted physicalLocation.
	 */
	public async deletePhysicalLocation(id: number): Promise<void> {
		const index = this.physicalLocations.findIndex(d => d.id === id);
		if (index === -1) {
			throw new Error(`no such PhysicalLocation: ${id}`);
		}
		return;
	}
}
