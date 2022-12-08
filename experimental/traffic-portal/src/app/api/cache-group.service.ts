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
import type { RequestDivision, ResponseDivision, RequestRegion, ResponseRegion } from "trafficops-types";

import type { CacheGroup } from "src/app/models";

import { APIService } from "./base-api.service";

/**
 * CDNService expose API functionality relating to CDNs.
 */
@Injectable()
export class CacheGroupService extends APIService {
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
		const path = "cachegroups";
		if (idOrName !== undefined) {
			let params;
			switch (typeof(idOrName)) {
				case "string":
					params = {name: idOrName};
					break;
				case "number":
					params = {id: String(idOrName)};
			}
			return this.get<[CacheGroup]>(path, undefined, params).toPromise().then(
				r => {
					const cg = r[0];
					if (cg.id !== idOrName) {
						throw new Error(`Traffic Ops returned no match for ID ${idOrName}`);
					}
					//  lastUpdated comes in as a string
					cg.lastUpdated = cg.lastUpdated ? new Date((cg.lastUpdated as unknown as string).replace("+00", "Z")) : undefined;
					return cg;
				}
			).catch(
				e => {
					console.error("Failed to get Cache Group with identifier", idOrName, ":", e);
					return {
						fallbackToClosest: false,
						fallbacks: [],
						latitude: 0,
						localizationMethods: [],
						longitude: 0,
						name: "",
						parentCacheGroupID: -1,
						parentCacheGroupName: "",
						secondaryParentCacheGroupID: -1,
						secondaryParentCacheGroupName: "",
						shortName: "",
						typeId: -1,
						typeName: ""
					};
				}
			);
		}
		return this.get<Array<CacheGroup>>(path).toPromise().then(r => r.map(
			cg => {
				if (cg.lastUpdated) {
					cg.lastUpdated = new Date((cg.lastUpdated as unknown as string).replace("+00", "Z"));
				}
				return cg;
			}
		)).catch(
			e => {
				console.error("Failed to get Cache Groups:", e);
				return [];
			}
		);
	}

	public async getDivisions(): Promise<Array<ResponseDivision>>;
	public async getDivisions(nameOrID: string | number): Promise<ResponseDivision>;

	/**
	 * Gets an array of divisions from Traffic Ops.
	 *
	 * @param nameOrID If given, returns only the Division with the given name
	 * (string) or ID (number).
	 * @returns An Array of Division objects - or a single Division object if 'nameOrID'
	 * was given.
	 */
	public async getDivisions(nameOrID?: string | number): Promise<Array<ResponseDivision> | ResponseDivision> {
		const path = "divisions";
		if(nameOrID) {
			let params;
			switch (typeof nameOrID) {
				case "string":
					params = {name: nameOrID};
					break;
				case "number":
					params = {id: String(nameOrID)};
			}
			const r = await this.get<[ResponseDivision]>(path, undefined, params).toPromise();
			return {...r[0], lastUpdated: new Date((r[0].lastUpdated as unknown as string).replace("+00", "Z"))};

		}
		const divisions = await this.get<Array<ResponseDivision>>(path).toPromise();
		return divisions.map(
			d => ({...d, lastUpdated: new Date((d.lastUpdated as unknown as string).replace("+00", "Z"))})
		);
	}

	/**
	 * Replaces the current definition of a division with the one given.
	 *
	 * @param division The new division.
	 * @returns The updated division.
	 */
	public async updateDivision(division: ResponseDivision): Promise<ResponseDivision> {
		const path = `divisions/${division.id}`;
		const response = await this.put<ResponseDivision>(path, division).toPromise();
		return {
			...response,
			lastUpdated: new Date((response.lastUpdated as unknown as string).replace(" ", "T").replace("+00", "Z"))
		};
	}

	/**
	 * Creates a new division.
	 *
	 * @param division The division to create.
	 * @returns The created division.
	 */
	public async createDivision(division: RequestDivision): Promise<ResponseDivision> {
		const response = await this.post<ResponseDivision>("divisions", division).toPromise();
		return {
			...response,
			lastUpdated: new Date((response.lastUpdated as unknown as string).replace(" ", "T").replace("+00", "Z"))
		};
	}

	/**
	 * Deletes an existing division.
	 *
	 * @param id Id of the division to delete.
	 * @returns The deleted division.
	 */
	public async deleteDivision(id: number): Promise<ResponseDivision> {
		return this.delete<ResponseDivision>(`divisions/${id}`).toPromise();
	}

	public async getRegions(): Promise<Array<ResponseRegion>>;
	public async getRegions(nameOrID: string | number): Promise<ResponseRegion>;

	/**
	 * Gets an array of regions from Traffic Ops.
	 *
	 * @param nameOrID If given, returns only the Region with the given name
	 * (string) or ID (number).
	 * @returns An Array of Region objects - or a single Region object if 'nameOrID'
	 * was given.
	 */
	public async getRegions(nameOrID?: string | number): Promise<Array<ResponseRegion> | ResponseRegion> {
		const path = "regions";
		if(nameOrID) {
			let params;
			switch (typeof nameOrID) {
				case "string":
					params = {name: nameOrID};
					break;
				case "number":
					params = {id: String(nameOrID)};
			}
			const r = await this.get<[ResponseRegion]>(path, undefined, params).toPromise();
			return {...r[0], lastUpdated: new Date((r[0].lastUpdated as unknown as string).replace("+00", "Z"))};

		}
		const regions = await this.get<Array<ResponseRegion>>(path).toPromise();
		return regions.map(
			d => ({...d, lastUpdated: new Date((d.lastUpdated as unknown as string).replace("+00", "Z"))})
		);
	}

	/**
	 * Replaces the current definition of a region with the one given.
	 *
	 * @param region The new region.
	 * @returns The updated region.
	 */
	public async updateRegion(region: ResponseRegion): Promise<ResponseRegion> {
		const path = `regions/${region.id}`;
		const response = await this.put<ResponseRegion>(path, region).toPromise();
		return {
			...response,
			lastUpdated: new Date((response.lastUpdated as unknown as string).replace(" ", "T").replace("+00", "Z"))
		};
	}

	/**
	 * Creates a new region.
	 *
	 * @param region The region to create.
	 * @returns The created region.
	 */
	public async createRegion(region: RequestRegion): Promise<ResponseRegion> {
		const response = await this.post<ResponseRegion>("regions", region).toPromise();
		return {
			...response,
			lastUpdated: new Date((response.lastUpdated as unknown as string).replace(" ", "T").replace("+00", "Z"))
		};
	}

	/**
	 * Deletes an existing region.
	 *
	 * @param id Id of the region to delete.
	 * @returns The deleted region.
	 */
	public async deleteRegion(id: number): Promise<ResponseRegion> {
		return this.delete<ResponseRegion>(`regions/${id}`).toPromise();
	}

	constructor(http: HttpClient) {
		super(http);
	}
}
