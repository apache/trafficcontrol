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
import { RequestDivision, ResponseDivision, RequestRegion, ResponseRegion } from "trafficops-types";

import type { CacheGroup } from "src/app/models";

/**
 * CDNService expose API functionality relating to CDNs.
 */
@Injectable()
export class CacheGroupService {
	private lastID = 10;

	private readonly divisions: Array<ResponseDivision> = [{
		id: 1,
		lastUpdated: new Date(),
		name: "Div1"
	}
	];
	private readonly regions: Array<ResponseRegion> = [{
		division: 1,
		divisionName: "div1",
		id: 1,
		lastUpdated: new Date(),
		name: "Reg1"
	}
	];
	private readonly cacheGroups = [
		{
			fallbackToClosest: true,
			fallbacks: [],
			id: 1,
			latitude: 0,
			localizationMethods: [],
			longitude: 0,
			name: "Mid",
			parentCacheGroupID: null,
			parentCacheGroupName: null,
			secondaryParentCacheGroupID: null,
			secondaryParentCacheGroupName: null,
			shortName: "Mid",
			typeId: 1,
			typeName: "MID_LOC"
		},
		{
			fallbackToClosest: true,
			fallbacks: [],
			id: 2,
			latitude: 0,
			localizationMethods: [],
			longitude: 0,
			name: "Edge",
			parentCacheGroupID: 1,
			parentCacheGroupName: "Mid",
			secondaryParentCacheGroupID: null,
			secondaryParentCacheGroupName: null,
			shortName: "Edge",
			typeId: 2,
			typeName: "EDGE_LOC"
		},
		{
			fallbackToClosest: true,
			fallbacks: [],
			id: 3,
			latitude: 0,
			localizationMethods: [],
			longitude: 0,
			name: "Origin",
			parentCacheGroupID: null,
			parentCacheGroupName: null,
			secondaryParentCacheGroupID: null,
			secondaryParentCacheGroupName: null,
			shortName: "Origin",
			typeId: 3,
			typeName: "ORG_LOC"
		},
		{
			fallbackToClosest: true,
			fallbacks: [],
			id: 4,
			latitude: 0,
			localizationMethods: [],
			longitude: 0,
			name: "Other",
			parentCacheGroupID: null,
			parentCacheGroupName: null,
			secondaryParentCacheGroupID: null,
			secondaryParentCacheGroupName: null,
			shortName: "Other",
			typeId: 4,
			typeName: "TC_LOC"
		}
	];

	public async getCacheGroups(idOrName: number | string): Promise<CacheGroup>;
	public async getCacheGroups(): Promise<Array<CacheGroup>>;
	/**
	 * Gets one or all CDNs from Traffic Ops
	 *
	 * @param idOrName Optionally either the name or integral, unique identifier of a single Cache Group to be returned.
	 * @returns Either an Array of CacheGroup objects, or a single CacheGroup, depending on whether
	 * `idOrName` was 	passed.
	 * @throws {Error} In the event that `idOrName` is passed but does not match any CacheGroup.
	 */
	public async getCacheGroups(idOrName?: number | string): Promise<Array<CacheGroup> | CacheGroup> {
		if (idOrName !== undefined) {
			let cacheGroup;
			switch (typeof(idOrName)) {
				case "string":
					cacheGroup = this.cacheGroups.filter(cg=>cg.name===idOrName)[0];
					break;
				case "number":
					cacheGroup = this.cacheGroups.filter(cg=>cg.id===idOrName)[0];
			}
			if (!cacheGroup) {
				throw new Error(`no such Cache Group: ${idOrName}`);
			}
			return cacheGroup;
		}
		return this.cacheGroups;
	}

	public async getDivisions(): Promise<Array<ResponseDivision>>;
	public async getDivisions(nameOrID: string | number): Promise<ResponseDivision>;

	/**
	 * Gets an array of divisions from Traffic Ops.
	 *
	 * @param nameOrID If given, returns only the ResponseDivision with the given name
	 * (string) or ID (number).
	 * @returns An Array of ResponseDivision objects - or a single ResponseDivision object if 'nameOrID'
	 * was given.
	 */
	public async getDivisions(nameOrID?: string | number): Promise<Array<ResponseDivision> | ResponseDivision> {
		if(nameOrID) {
			let division;
			switch (typeof nameOrID) {
				case "string":
					division = this.divisions.find(d=>d.name === nameOrID);
					break;
				case "number":
					division = this.divisions.find(d=>d.id === nameOrID);
			}
			if (!division) {
				throw new Error(`no such Division: ${nameOrID}`);
			}
			return division;
		}
		return this.divisions;
	}

	/**
	 * Replaces the current definition of a division with the one given.
	 *
	 * @param division The new division.
	 * @returns The updated division.
	 */
	public async updateDivision(division: ResponseDivision): Promise<ResponseDivision> {
		const id = this.divisions.findIndex(d => d.id === division.id);
		if (id === -1) {
			throw new Error(`no such Division: ${division.id}`);
		}
		this.divisions[id] = division;
		return division;
	}

	/**
	 * Creates a new division.
	 *
	 * @param division The division to create.
	 * @returns The created division.
	 */
	public async createDivision(division: RequestDivision): Promise<ResponseDivision> {
		return {
			...division,
			id: ++this.lastID,
			lastUpdated: new Date()
		};
	}

	/**
	 * Deletes an existing division.
	 *
	 * @param id Id of the division to delete.
	 * @returns The deleted division.
	 */
	public async deleteDivision(id: number): Promise<ResponseDivision> {
		const index = this.divisions.findIndex(d => d.id === id);
		if (index === -1) {
			throw new Error(`no such Division: ${id}`);
		}
		return this.divisions.splice(index, 1)[0];
	}

	public async getRegions(): Promise<Array<ResponseRegion>>;
	public async getRegions(nameOrID: string | number): Promise<ResponseRegion>;

	/**
	 * Gets an array of regions from Traffic Ops.
	 *
	 * @param nameOrID If given, returns only the ResponseRegion with the given name
	 * (string) or ID (number).
	 * @returns An Array of ResponseRegion objects - or a single ResponseRegion object if 'nameOrID'
	 * was given.
	 */
	public async getRegions(nameOrID?: string | number): Promise<Array<ResponseRegion> | ResponseRegion> {
		if(nameOrID) {
			let region;
			switch (typeof nameOrID) {
				case "string":
					region = this.regions.find(d=>d.name === nameOrID);
					break;
				case "number":
					region = this.regions.find(d=>d.id === nameOrID);
			}
			if (!region) {
				throw new Error(`no such Region: ${nameOrID}`);
			}
			return region;
		}
		return this.regions;
	}

	/**
	 * Replaces the current definition of a region with the one given.
	 *
	 * @param region The new region.
	 * @returns The updated region.
	 */
	public async updateRegion(region: ResponseRegion): Promise<ResponseRegion> {
		const id = this.regions.findIndex(d => d.id === region.id);
		if (id === -1) {
			throw new Error(`no such Region: ${region.id}`);
		}
		this.regions[id] = region;
		return region;
	}

	/**
	 * Creates a new region.
	 *
	 * @param region The region to create.
	 * @returns The created region.
	 */
	public async createRegion(region: RequestRegion): Promise<ResponseRegion> {
		return {
			divisionName: "Div1",
			...region,
			id: ++this.lastID,
			lastUpdated: new Date()
		};
	}

	/**
	 * Deletes an existing region.
	 *
	 * @param id Id of the region to delete.
	 * @returns The deleted region.
	 */
	public async deleteRegion(id: number): Promise<ResponseRegion> {
		const index = this.regions.findIndex(d => d.id === id);
		if (index === -1) {
			throw new Error(`no such Region: ${id}`);
		}
		return this.regions.splice(index, 1)[0];
	}
}
