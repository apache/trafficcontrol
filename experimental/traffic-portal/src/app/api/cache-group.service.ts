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
import type {
	RequestDivision,
	ResponseDivision,
	RequestRegion,
	ResponseRegion,
	RequestPhysicalLocation,
	ResponsePhysicalLocation,
	ResponseCacheGroup,
	RequestCacheGroup,
	CDN,
	CacheGroupQueueResponse,
	CacheGroupQueueRequest,
} from "trafficops-types";

import type { CacheGroup } from "src/app/models";

import { APIService } from "./base-api.service";

/**
 * Checks the type of an argument to
 * {@link CacheGroupService.queueCacheGroupUpdates}.
 *
 * @param x The object to check.
 * @returns Whether `x` is an {@link CacheGroupQueueRequest}.
 */
function isRequest(x: CacheGroupQueueRequest | CDN | string | number): x is CacheGroupQueueRequest {
	return Object.prototype.hasOwnProperty.call(x, "action");
}

/**
 * CDNService expose API functionality relating to CDNs.
 */
@Injectable()
export class CacheGroupService extends APIService {

	/**
	 * Gets a single Cache Group from Traffic Ops.
	 *
	 * @param idOrName Either the name or integral, unique identifier of the
	 * single Cache Group to be returned.
	 * @returns The Cache Group identified by `idOrName`.
	 * @throws {Error} When no matching Cache Group is found in Traffic Ops.
	 */
	public async getCacheGroups(idOrName: number | string): Promise<ResponseCacheGroup>;
	/**
	 * Gets Cache Groups from Traffic Ops.
	 *
	 * @returns All requested Cache Groups.
	 */
	public async getCacheGroups(): Promise<Array<ResponseCacheGroup>>;
	/**
	 * Gets one or all Cache Groups from Traffic Ops
	 *
	 * @param idOrName Optionally either the name or integral, unique identifier
	 * of a single Cache Group to be returned.
	 * @returns Either an Array of Cache Group objects, or a single Cache Group,
	 * depending on whether `idOrName` was 	passed.
	 * @throws {Error} In the event that `idOrName` is passed but does not match
	 * any Cache Group.
	 */
	public async getCacheGroups(idOrName?: number | string): Promise<Array<ResponseCacheGroup> | ResponseCacheGroup> {
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
			const resp = await this.get<[ResponseCacheGroup]>(path, undefined, params).toPromise();
			if (resp.length !== 1) {
				throw new Error(`Traffic Ops returned wrong number of results for Cache Group identifier: ${params}`);
			}
			const cg = resp[0];
			//  lastUpdated comes in as a string
			return {...cg, lastUpdated: new Date((cg.lastUpdated as unknown as string).replace("+00", "Z"))};
		}
		const r = await this.get<Array<ResponseCacheGroup>>(path).toPromise();
		return r.map(cg => ({...cg, lastUpdated: new Date((cg.lastUpdated as unknown as string).replace("+00", "Z"))}));
	}

	/**
	 * Deletes a Cache Group.
	 *
	 * @param cacheGroup The Cache Group to be deleted, or just its ID.
	 */
	public async deleteCacheGroup(cacheGroup: ResponseCacheGroup | number): Promise<void> {
		const id = typeof(cacheGroup) === "number" ? cacheGroup : cacheGroup.id;
		return this.delete(`cachegroups/${id}`).toPromise();
	}

	/**
	 * Creates a new Cache Group.
	 *
	 * @param cacheGroup The Cache Group to create.
	 */
	public async createCacheGroup(cacheGroup: RequestCacheGroup): Promise<ResponseCacheGroup> {
		return this.post<ResponseCacheGroup>("cachegroups", cacheGroup).toPromise();
	}

	/**
	 * Replaces an existing Cache Group with the provided new definition of a
	 * Cache Group.
	 *
	 * @param id The if of the Cache Group being updated.
	 * @param cacheGroup The new definition of the Cache Group.
	 */
	public async updateCacheGroup(id: number, cacheGroup: RequestCacheGroup): Promise<ResponseCacheGroup>;
	/**
	 * Replaces an existing Cache Group with the provided new definition of a
	 * Cache Group.
	 *
	 * @param cacheGroup The full new definition of the Cache Group being
	 * updated.
	 */
	public async updateCacheGroup(cacheGroup: ResponseCacheGroup): Promise<ResponseCacheGroup>;
	/**
	 * Replaces an existing Cache Group with the provided new definition of a
	 * Cache Group.
	 *
	 * @param cacheGroupOrID The full new definition of the Cache Group being
	 * updated, or just its ID.
	 * @param payload The new definition of the Cache Group. This is required if
	 * `cacheGroupOrID` is an ID, and ignored otherwise.
	 */
	public async updateCacheGroup(cacheGroupOrID: ResponseCacheGroup | number, payload?: RequestCacheGroup): Promise<ResponseCacheGroup> {
		let id;
		let body;
		if (typeof(cacheGroupOrID) === "number") {
			if (!payload) {
				throw new TypeError("invalid call signature - missing request payload");
			}
			body = payload;
			id = cacheGroupOrID;
		} else {
			body = cacheGroupOrID;
			({id} = cacheGroupOrID);
		}

		return this.put<ResponseCacheGroup>(`cachegroups/${id}`, body).toPromise();
	}

	/**
	 * Queues (or dequeues) updates on a Cache Group's servers.
	 *
	 * @param cacheGroupOrID The Cache Group on which updates will be queued, or
	 * just its ID.
	 * @param cdnOrIdentifier Either a CDN, its name, or its ID.
	 * @param action Used to determine the queue action to take. If not given,
	 * defaults to `queue`.
	 * @returns The API's response.
	 */
	public async queueCacheGroupUpdates(
		cacheGroupOrID: ResponseCacheGroup | number,
		cdnOrIdentifier: CDN | string | number,
		action?: "queue" | "dequeue"
	): Promise<CacheGroupQueueResponse>;
	/**
	 * Queues (or dequeues) updates on a Cache Group's servers.
	 *
	 * @param cacheGroupOrID The Cache Group on which updates will be queued, or
	 * just its ID.
	 * @param request The full (de/)queue request.
	 * @returns The API's response.
	 */
	public async queueCacheGroupUpdates(
		cacheGroupOrID: ResponseCacheGroup | number,
		request: CacheGroupQueueRequest
	): Promise<CacheGroupQueueResponse>;
	/**
	 * Queues (or dequeues) updates on a Cache Group's servers.
	 *
	 * @param cacheGroupOrID The Cache Group on which updates will be queued, or
	 * just its ID.
	 * @param cdnOrIdentifierOrRequest Either the full (de/)queue request or a
	 * CDN, its name, or its ID.
	 * @param action If `cdnOrIdentifierOrRequest` is not a full (de/)queue
	 * request, then this will be used to determine the queue action to take. If
	 * not given, defaults to `queue`.
	 * @returns The API's response.
	 */
	public async queueCacheGroupUpdates(
		cacheGroupOrID: ResponseCacheGroup | number,
		cdnOrIdentifierOrRequest: CacheGroupQueueRequest | CDN | string | number,
		action?: "queue" | "dequeue"
	): Promise<CacheGroupQueueResponse> {
		const cgID = typeof(cacheGroupOrID) === "number" ? cacheGroupOrID : cacheGroupOrID.id;
		const path = `cachegroups/${cgID}/queue_update`;
		let request: CacheGroupQueueRequest;
		if (isRequest(cdnOrIdentifierOrRequest)) {
			request = cdnOrIdentifierOrRequest;
		} else {
			action = action ?? "queue";
			switch (typeof(cdnOrIdentifierOrRequest)) {
				case "string":
					request = {
						action,
						cdn: cdnOrIdentifierOrRequest,
					};
					break;
				case "number":
					request = {
						action,
						cdnId: cdnOrIdentifierOrRequest,
					};
					break;
				default:
					request = {
						action,
						cdn: cdnOrIdentifierOrRequest.name
					};
			}
		}
		return this.post<CacheGroupQueueResponse>(path, request).toPromise();
	}

	public async getPhysicalLocations(): Promise<Array<ResponsePhysicalLocation>>;
	public async getPhysicalLocations(nameOrID: string | number): Promise<ResponsePhysicalLocation>;

	/**
	 * Gets an array of physicalLocations from Traffic Ops.
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
					params = {id: String(nameOrID)};
			}
			const r = await this.get<[ResponsePhysicalLocation]>(path, undefined, params).toPromise();
			return {...r[0], lastUpdated: new Date((r[0].lastUpdated as unknown as string).replace("+00", "Z"))};

		}
		const physicalLocations = await this.get<Array<ResponsePhysicalLocation>>(path).toPromise();
		return physicalLocations.map(
			d => ({...d, lastUpdated: new Date((d.lastUpdated as unknown as string).replace("+00", "Z"))})
		);
	}

	/**
	 * Replaces the current definition of a physicalLocation with the one given.
	 *
	 * @param physicalLocation The new physicalLocation.
	 * @returns The updated physicalLocation.
	 */
	public async updatePhysicalLocation(physicalLocation: ResponsePhysicalLocation): Promise<ResponsePhysicalLocation> {
		const path = `phys_locations/${physicalLocation.id}`;
		const response = await this.put<ResponsePhysicalLocation>(path, physicalLocation).toPromise();
		return {
			...response,
			lastUpdated: new Date((response.lastUpdated as unknown as string).replace(" ", "T").replace("+00", "Z"))
		};
	}

	/**
	 * Creates a new physicalLocation.
	 *
	 * @param physicalLocation The physicalLocation to create.
	 * @returns The created physicalLocation.
	 */
	public async createPhysicalLocation(physicalLocation: RequestPhysicalLocation): Promise<ResponsePhysicalLocation> {
		const response = await this.post<ResponsePhysicalLocation>("physicalLocations", physicalLocation).toPromise();
		return {
			...response,
			lastUpdated: new Date((response.lastUpdated as unknown as string).replace(" ", "T").replace("+00", "Z"))
		};
	}

	/**
	 * Deletes an existing physicalLocation.
	 *
	 * @param id Id of the physicalLocation to delete.
	 * @returns The deleted physicalLocation.
	 */
	public async deletePhysicalLocation(id: number): Promise<ResponsePhysicalLocation> {
		return this.delete<ResponsePhysicalLocation>(`phys_locations/${id}`).toPromise();
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
	 * @param regionOrId Id of the region to delete.
	 * @returns The deleted region.
	 */
	public async deleteRegion(regionOrId: number | ResponseRegion): Promise<void> {
		const id = typeof(regionOrId) === "number" ? regionOrId : regionOrId.id;
		await this.delete("regions/", undefined, { id : String(id) }).toPromise();
	}

	constructor(http: HttpClient) {
		super(http);
	}
}
