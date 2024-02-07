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
import type { RequestPhysicalLocation, ResponsePhysicalLocation } from "trafficops-types";

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

	/**
	 * Gets all Physical Locations.
	 *
	 * @returns All stored Physical Locations.
	 */
	public async getPhysicalLocations(): Promise<Array<ResponsePhysicalLocation>>;
	/**
	 * Gets a specific Physical Location.
	 *
	 * @param nameOrID The name (string) or ID (number) of the Physical Location
	 * to be returned.
	 * @returns The requested Physical Location.
	 */
	public async getPhysicalLocations(nameOrID: string | number): Promise<ResponsePhysicalLocation>;
	/**
	 * Gets one or all Physical Location(s).
	 *
	 * @param nameOrID If given, returns only the PhysicalLocation with the
	 * given name (string) or ID (number).
	 * @returns The requested Physical Location(s).
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
	 * Replaces the current definition of a Physical Location with the one given.
	 *
	 * @param physicalLocation The new Physical Location.
	 * @returns The updated Physical Location.
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
	 * Creates a new Physical Location.
	 *
	 * @param physicalLocation The Physical Location to create.
	 * @returns The created Physical Location.
	 */
	public async createPhysicalLocation(physicalLocation: RequestPhysicalLocation): Promise<ResponsePhysicalLocation> {
		const phys = {
			...physicalLocation,
			comments: physicalLocation.comments ?? null,
			email: physicalLocation.email ?? null,
			id: ++this.lastID,
			lastUpdated: new Date(),
			phone: physicalLocation.phone ?? null,
			poc: physicalLocation.poc ?? null,
			region: ""
		};
		this.physicalLocations.push(phys);
		return phys;
	}

	/**
	 * Deletes an existing Physical Location.
	 *
	 * @param physLoc The Physical Location to be deleted (or its ID)
	 */
	public async deletePhysicalLocation(physLoc: ResponsePhysicalLocation | number): Promise<void> {
		const id = typeof(physLoc) === "number" ? physLoc : physLoc.id;
		const index = this.physicalLocations.findIndex(d => d.id === id);
		if (index === -1) {
			throw new Error(`no such PhysicalLocation: ${id}`);
		}
		this.physicalLocations.splice(index, 1);
	}
}
