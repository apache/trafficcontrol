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
	RequestASN,
	ResponseASN,
	RequestDivision,
	ResponseDivision,
	RequestRegion,
	ResponseRegion,
	ResponseCacheGroup,
	RequestCacheGroup,
	CDN,
	CacheGroupQueueResponse,
	CacheGroupQueueRequest,
	RequestCoordinate,
	ResponseCoordinate,
} from "trafficops-types";

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
					params = {id: idOrName};
			}
			const resp = await this.get<[ResponseCacheGroup]>(path, undefined, params).toPromise();
			if (resp.length !== 1) {
				throw new Error(`Traffic Ops returned wrong number of results for Cache Group identifier: ${params}`);
			}
			return resp[0];
		}
		return this.get<Array<ResponseCacheGroup>>(path).toPromise();
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

	/**
	 * Gets all Divisions from Traffic Ops.
	 *
	 * @returns The requested Divisions.
	 */
	public async getDivisions(): Promise<Array<ResponseDivision>>;
	/**
	 * Gets a single Division from Traffic Ops.
	 *
	 * @param nameOrID Either the name (string) or ID (number) of the single
	 * Division to be fetched.
	 * @returns The requested Division.
	 */
	public async getDivisions(nameOrID: string | number): Promise<ResponseDivision>;
	/**
	 * Gets a Division or Divisions from Traffic Ops.
	 *
	 * @param nameOrID If given, returns only the Division with the given name
	 * (string) or ID (number).
	 * @returns The requested Division or Divisions.
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
					params = {id: nameOrID};
			}
			const div = await this.get<[ResponseDivision]>(path, undefined, params).toPromise();
			return div[0];

		}
		return this.get<Array<ResponseDivision>>(path).toPromise();
	}

	/**
	 * Replaces the current definition of a Division with the one given.
	 *
	 * @param division The new Division.
	 * @returns The updated Division.
	 */
	public async updateDivision(division: ResponseDivision): Promise<ResponseDivision> {
		const path = `divisions/${division.id}`;
		return this.put<ResponseDivision>(path, division).toPromise();
	}

	/**
	 * Creates a new division.
	 *
	 * @param division The division to create.
	 * @returns The created division.
	 */
	public async createDivision(division: RequestDivision): Promise<ResponseDivision> {
		return this.post<ResponseDivision>("divisions", division).toPromise();
	}

	/**
	 * Deletes an existing Division.
	 *
	 * @param division The Division to be deleted, or just its ID.
	 * @returns The deleted Division.
	 */
	public async deleteDivision(division: number | ResponseDivision): Promise<ResponseDivision> {
		const id = typeof(division) === "number" ? division : division.id;
		return this.delete<ResponseDivision>(`divisions/${id}`).toPromise();
	}

	/**
	 * Gets all Regions from Traffic Ops.
	 *
	 * @returns The requested Regions.
	 */
	public async getRegions(): Promise<Array<ResponseRegion>>;
	/**
	 * Gets a single Region from Traffic Ops.
	 *
	 * @param nameOrID The name (string) or ID (number) of the single Region to
	 * be fetched.
	 * @returns The requested Region.
	 */
	public async getRegions(nameOrID: string | number): Promise<ResponseRegion>;
	/**
	 * Gets a Region or Regions from Traffic Ops.
	 *
	 * @param nameOrID If given, returns only the Region with the given name
	 * (string) or ID (number).
	 * @returns The requested Region or Regions.
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
					params = {id: nameOrID};
			}
			const r = await this.get<[ResponseRegion]>(path, undefined, params).toPromise();
			return r[0];
		}
		return this.get<Array<ResponseRegion>>(path).toPromise();
	}

	/**
	 * Replaces the current definition of a Region with the one given.
	 *
	 * @param region The new Region.
	 * @returns The updated Region.
	 */
	public async updateRegion(region: ResponseRegion): Promise<ResponseRegion> {
		const path = `regions/${region.id}`;
		return this.put<ResponseRegion>(path, region).toPromise();
	}

	/**
	 * Creates a new Region.
	 *
	 * @param region The region to create.
	 * @returns The created region.
	 */
	public async createRegion(region: RequestRegion): Promise<ResponseRegion> {
		return this.post<ResponseRegion>("regions", region).toPromise();
	}

	/**
	 * Deletes an existing Region.
	 *
	 * @param regionOrId The ID of the Region to delete.
	 * @returns The deleted Region.
	 */
	public async deleteRegion(regionOrId: number | ResponseRegion): Promise<void> {
		const id = typeof(regionOrId) === "number" ? regionOrId : regionOrId.id;
		await this.delete("regions", undefined, { id }).toPromise();
	}

	/**
	 * Gets all Coordinates from Traffic Ops.
	 *
	 * @returns The requested Coordinates.
	 */
	public async getCoordinates(): Promise<Array<ResponseCoordinate>>;
	/**
	 * Gets a single Coordinate from Traffic Ops.
	 *
	 * @param nameOrID The name (string) or ID (number) of the single Coordinate
	 * to be fetched.
	 * @returns The requested Coordinate.
	 */
	public async getCoordinates(nameOrID: string | number): Promise<ResponseCoordinate>;
	/**
	 * Gets a Coordinate or Coordinates from Traffic Ops.
	 *
	 * @param nameOrID If given, returns only the Coordinate with the given name
	 * (string) or ID (number).
	 * @returns The requested Coordinate or Coordinates.
	 */
	public async getCoordinates(nameOrID?: string | number): Promise<Array<ResponseCoordinate> | ResponseCoordinate> {
		const path = "coordinates";
		if(nameOrID) {
			let params;
			switch (typeof nameOrID) {
				case "string":
					params = {name: nameOrID};
					break;
				case "number":
					params = {id: nameOrID};
			}
			const r = await this.get<[ResponseCoordinate]>(path, undefined, params).toPromise();
			return r[0];
		}
		return this.get<Array<ResponseCoordinate>>(path).toPromise();
	}

	/**
	 * Replaces the current definition of a Coordinate with the one given.
	 *
	 * @param coordinate The new Coordinate.
	 * @returns The updated Coordinate.
	 */
	public async updateCoordinate(coordinate: ResponseCoordinate): Promise<ResponseCoordinate> {
		return this.put<ResponseCoordinate>("coordinates", coordinate, { id: coordinate.id }).toPromise();
	}

	/**
	 * Creates a new Coordinate.
	 *
	 * @param coordinate The Coordinate to create.
	 * @returns The created Coordinate.
	 */
	public async createCoordinate(coordinate: RequestCoordinate): Promise<ResponseCoordinate> {
		return this.post<ResponseCoordinate>("coordinates", coordinate).toPromise();
	}

	/**
	 * Deletes an existing Coordinate.
	 *
	 * @param coordinateOrId The ID of the Coordinate to delete.
	 * @returns The deleted Coordinate.
	 */
	public async deleteCoordinate(coordinateOrId: number | ResponseCoordinate): Promise<void> {
		const id = typeof(coordinateOrId) === "number" ? coordinateOrId : coordinateOrId.id;
		await this.delete("coordinates", undefined, { id }).toPromise();
	}

	/**
	 * Gets all ASNs from Traffic Ops.
	 *
	 * @returns All ASNs configured in Traffic Ops.
	 */
	public async getASNs(): Promise<Array<ResponseASN>>;
	/**
	 * Gets a single ASN from Traffic Ops.
	 *
	 * @param id The ID of the ASN to fetch.
	 * @returns The ASN with the given ID.
	 */
	public async getASNs(id: number): Promise<ResponseASN>;

	/**
	 * Gets ASNs from Traffic Ops.
	 *
	 * @param id If given, returns only the ASN with the given ID.
	 * @returns An Array of ASNs objects - or a single ASN object if `id`
	 * was given.
	 */
	public async getASNs(id?: number): Promise<Array<ResponseASN> | ResponseASN> {
		const path = "/asns";
		if(id) {
			const r = await this.get<[ResponseASN]>(path, undefined, { id }).toPromise();
			return r[0];

		}
		return this.get<Array<ResponseASN>>(path).toPromise();
	}

	/**
	 * Replaces the current definition of a ASN with the one given.
	 *
	 * @param asn The new ASN.
	 * @returns The updated ASN.
	 */
	public async updateASN(asn: ResponseASN): Promise<ResponseASN> {
		const path = `asns/${asn.id}`;
		return this.put<ResponseASN>(path, asn).toPromise();
	}

	/**
	 * Creates a new ASN.
	 *
	 * @param asn The ASN to create.
	 * @returns The created ASN.
	 */
	public async createASN(asn: RequestASN): Promise<ResponseASN> {
		return this.post<ResponseASN>("asns", asn).toPromise();
	}

	/**
	 * Deletes an existing ASN.
	 *
	 * @param asn The ASN to be deleted or ID of the ASN to delete.
	 * @returns The deleted ASN.
	 */
	public async deleteASN(asn: ResponseASN | number): Promise<void> {
		const id = typeof(asn) === "number" ? asn : asn.id;
		return this.delete(`asns/${id}`).toPromise();
	}

	constructor(http: HttpClient) {
		super(http);
	}
}
