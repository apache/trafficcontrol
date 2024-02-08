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
import type {
	CacheGroupQueueRequest,
	CacheGroupQueueResponse,
	CDN,
	RequestASN,
	ResponseASN,
	RequestCacheGroup,
	ResponseCacheGroup,
	RequestCoordinate,
	ResponseCoordinate,
	RequestDivision,
	RequestRegion,
	ResponseDivision,
	ResponseRegion,
} from "trafficops-types";

import { ServerService } from "./server.service";

/**
 * The names of properties of {@link ResponseCacheGroup}s that define its
 * primary parentage.
 */
type ParentKeys = "parentCachegroupId" | "parentCachegroupName";
/**
 * The names of properties of {@link ResponseCacheGroup}s that define its
 * secondary parentage.
 */
type SecondaryParentKeys = "secondaryParentCachegroupId" | "secondaryParentCachegroupName";
/**
 * The names of properties of {@link ResponseCacheGroup}s that define its
 * parentage; both primary and secondary.
 */
type AllParentageKeys = ParentKeys | SecondaryParentKeys;

/**
 * The parts of a Cache Group pertaining to its primary parentage.
 */
type Parentage = {
	parentCachegroupId: null;
	parentCachegroupName: null;
} | {
	parentCachegroupId: number;
	parentCachegroupName: string;
};

/**
 * The parts of a Cache Group pertaining to its secondary parentage.
 */
type SecondaryParentage = {
	secondaryParentCachegroupId: null;
	secondaryParentCachegroupName: null;
} | {
	secondaryParentCachegroupId: number;
	secondaryParentCachegroupName: string;
};

/**
 * Contains all information about a Cache Groups parents, primary and secondary.
 */
type AllParentage = Parentage & SecondaryParentage;

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
export class CacheGroupService {
	private lastID = 10;

	private readonly asns: Array<ResponseASN> = [{
		asn: 0,
		cachegroup: "Mid",
		cachegroupId: 1,
		id: 1,
		lastUpdated: new Date()
	}
	];
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
	private readonly cacheGroups: Array<ResponseCacheGroup> = [
		{
			fallbackToClosest: true,
			fallbacks: [],
			id: 1,
			lastUpdated: new Date(),
			latitude: 0,
			localizationMethods: [],
			longitude: 0,
			name: "Mid",
			parentCachegroupId: null,
			parentCachegroupName: null,
			secondaryParentCachegroupId: null,
			secondaryParentCachegroupName: null,
			shortName: "Mid",
			typeId: 1,
			typeName: "MID_LOC"
		},
		{
			fallbackToClosest: true,
			fallbacks: [],
			id: 2,
			lastUpdated: new Date(),
			latitude: 0,
			localizationMethods: [],
			longitude: 0,
			name: "Edge",
			parentCachegroupId: 1,
			parentCachegroupName: "Mid",
			secondaryParentCachegroupId: null,
			secondaryParentCachegroupName: null,
			shortName: "Edge",
			typeId: 2,
			typeName: "EDGE_LOC"
		},
		{
			fallbackToClosest: true,
			fallbacks: [],
			id: 3,
			lastUpdated: new Date(),
			latitude: 0,
			localizationMethods: [],
			longitude: 0,
			name: "Origin",
			parentCachegroupId: null,
			parentCachegroupName: null,
			secondaryParentCachegroupId: null,
			secondaryParentCachegroupName: null,
			shortName: "Origin",
			typeId: 3,
			typeName: "ORG_LOC"
		},
		{
			fallbackToClosest: true,
			fallbacks: [],
			id: 4,
			lastUpdated: new Date(),
			latitude: 0,
			localizationMethods: [],
			longitude: 0,
			name: "Other",
			parentCachegroupId: null,
			parentCachegroupName: null,
			secondaryParentCachegroupId: null,
			secondaryParentCachegroupName: null,
			shortName: "Other",
			typeId: 4,
			typeName: "TC_LOC"
		}
	];
	private readonly coordinates: Array<ResponseCoordinate> = [{
		id: 1,
		lastUpdated: new Date(),
		latitude: 0,
		longitude: 0,
		name: "Coord1"
	}
	];

	constructor(private readonly servers: ServerService) {}

	/**
	 * Gets all Cache Groups.
	 *
	 * @returns All stored Cache Groups.
	 */
	public async getCacheGroups(): Promise<Array<ResponseCacheGroup>>;
	/**
	 * Gets a single Cache Group.
	 *
	 * @param idOrName Either the name or integral, unique identifier of the
	 * single Cache Group to be returned.
	 * @returns The requested Cache Group.
	 * @throws {Error} In the event that `idOrName` is passed but does not match
	 * any CacheGroup.
	 */
	public async getCacheGroups(idOrName: number | string): Promise<ResponseCacheGroup>;
	/**
	 * Gets one or all Cache Groups.
	 *
	 * @param idOrName Optionally either the name or integral, unique identifier
	 * of a single Cache Group to be returned.
	 * @returns Either all stored Cache Groups, or a single Cache Group,
	 * depending on whether `idOrName` was 	passed.
	 * @throws {Error} In the event that `idOrName` is passed but does not match
	 * any CacheGroup.
	 */
	public async getCacheGroups(idOrName?: number | string): Promise<Array<ResponseCacheGroup> | ResponseCacheGroup> {
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

	/**
	 * Deletes a Cache Group.
	 *
	 * @param cacheGroup The Cache Group to be deleted, or just its ID.
	 */
	public async deleteCacheGroup(cacheGroup: ResponseCacheGroup | number): Promise<void> {
		const id = typeof(cacheGroup) === "number" ? cacheGroup : cacheGroup.id;
		const idx = this.cacheGroups.findIndex(cg => cg.id === id);
		if (idx < 0) {
			throw new Error(`no such Cache Group: #${id}`);
		}
		this.cacheGroups.splice(idx, 1);
	}

	/**
	 * Gets the names of parents for a Cache Group from their IDs.
	 *
	 * @param parentID The ID of a parent Cache Group (or not).
	 * @param secondaryParentID The ID of a "secondary" parent Cache Group (or
	 * not).
	 * @returns The parentage portion of a Cache Group.
	 */
	private getParents(parentID: number | null | undefined, secondaryParentID: number | null | undefined): AllParentage {
		let parent: Parentage = {
			parentCachegroupId: null,
			parentCachegroupName: null
		};
		if (typeof(parentID) === "number") {
			const p = this.cacheGroups.find(cg => cg.id === parentID);
			if (!p) {
				throw new Error(`no such parent Cache Group: #${parentID}`);
			}
			parent = {
				parentCachegroupId: p.id,
				parentCachegroupName: p.name
			};
		}

		let secondaryParent: SecondaryParentage = {
			secondaryParentCachegroupId: null,
			secondaryParentCachegroupName: null
		};
		if (typeof(secondaryParentID) === "number") {
			const p = this.cacheGroups.find(cg => cg.id === secondaryParentID);
			if (!p) {
				throw new Error(`no such secondary parent Cache Group: #${secondaryParentID}`);
			}
			secondaryParent = {
				secondaryParentCachegroupId: p.id,
				secondaryParentCachegroupName: p.name
			};
		}

		return {
			...parent,
			...secondaryParent
		};
	}

	/**
	 * Creates a new Cache Group.
	 *
	 * @param cacheGroup The Cache Group to create.
	 */
	public async createCacheGroup(cacheGroup: RequestCacheGroup): Promise<ResponseCacheGroup> {
		const cg = {
			...cacheGroup,
			...this.getParents(cacheGroup.parentCachegroupId, cacheGroup.secondaryParentCachegroupId),
			fallbackToClosest: cacheGroup.fallbackToClosest ?? false,
			fallbacks: cacheGroup.fallbacks ?? [],
			id: ++this.lastID,
			lastUpdated: new Date(),
			latitude: cacheGroup.latitude ?? 0,
			localizationMethods: cacheGroup.localizationMethods ?? [],
			longitude: cacheGroup.longitude ?? 0,
			typeName: "",
		};
		this.cacheGroups.push(cg);
		return cg;
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
		let idx;
		let cg: Omit<ResponseCacheGroup, AllParentageKeys>;
		let parentCachegroupId;
		let secondaryParentCachegroupId;
		if (typeof(cacheGroupOrID) === "number") {
			if (!payload) {
				throw new TypeError("invalid call signature - missing request payload");
			}
			idx = this.cacheGroups.findIndex(c => c.id === cacheGroupOrID);
			cg = {
				...payload,
				fallbackToClosest: payload.fallbackToClosest ?? false,
				fallbacks: payload.fallbacks ?? [],
				id: ++this.lastID,
				lastUpdated: new Date(),
				latitude: payload.latitude ?? 0,
				localizationMethods: payload.localizationMethods ?? [],
				longitude: payload.longitude ?? 0,
				typeName: "",
			};
			parentCachegroupId = payload.parentCachegroupId;
			secondaryParentCachegroupId = payload.secondaryParentCachegroupId;
		} else {
			idx = this.cacheGroups.findIndex(c => c.id === cacheGroupOrID.id);
			cg = {
				...cacheGroupOrID,
				lastUpdated: new Date()
			};
			parentCachegroupId = cacheGroupOrID.parentCachegroupId;
			secondaryParentCachegroupId = cacheGroupOrID.secondaryParentCachegroupId;
		}

		if (idx < 0) {
			throw new Error(`no such Cache Group: #${cacheGroupOrID}`);
		}

		const final = {
			...cg,
			...this.getParents(parentCachegroupId, secondaryParentCachegroupId)
		};

		this.cacheGroups[idx] = final;

		return final;
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
		const cachegroupID = typeof(cacheGroupOrID) === "number" ? cacheGroupOrID : cacheGroupOrID.id;
		const cg = this.cacheGroups.find(c => c.id === cachegroupID);
		if (!cg) {
			throw new Error(`no such Cache Group: #${cachegroupID}`);
		}

		let cdn;
		if (isRequest(cdnOrIdentifierOrRequest)) {
			action = cdnOrIdentifierOrRequest.action;
			cdn = cdnOrIdentifierOrRequest.cdn ?? cdnOrIdentifierOrRequest.cdnId;
		} else {
			action = action ?? "queue";
			switch (typeof(cdnOrIdentifierOrRequest)) {
				case "string":
				case "number":
					cdn = cdnOrIdentifierOrRequest;
					break;
				default:
					cdn = cdnOrIdentifierOrRequest.name;
			}
		}
		const updPendingValue = action === "queue";
		const serverNames = [];
		for (const server of this.servers.servers) {
			if (server.cachegroupId === cachegroupID && (server.cdnId === cdn || server.cdnName === cdn)) {
				server.updPending = updPendingValue;
				serverNames.push(server.hostName);
			}
		}
		return {
			action,
			cachegroupID,
			cachegroupName: cg.name,
			cdn: String(cdn),
			serverNames,
		};
	}

	/**
	 * Gets all Divisions.
	 *
	 * @returns The requested Divisions.
	 */
	public async getDivisions(): Promise<Array<ResponseDivision>>;
	/**
	 * Gets a single Division.
	 *
	 * @param nameOrID Either the name (string) or ID (number) of the single
	 * Division to be returned.
	 * @returns The requested Division.
	 */
	public async getDivisions(nameOrID: string | number): Promise<ResponseDivision>;
	/**
	 * Gets a Division or Divisions.
	 *
	 * @param nameOrID If given, returns only the Division with the given name
	 * (string) or ID (number).
	 * @returns The requested Division or Divisions.
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
	 * Replaces the current definition of a Division with the one given.
	 *
	 * @param division The new Division.
	 * @returns The updated Division.
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
	 * Creates a new Division.
	 *
	 * @param division The Division to create.
	 * @returns The created Division.
	 */
	public async createDivision(division: RequestDivision): Promise<ResponseDivision> {
		const div = {
			...division,
			id: ++this.lastID,
			lastUpdated: new Date()
		};
		this.divisions.push(div);
		return div;
	}

	/**
	 * Deletes an existing Division.
	 *
	 * @param id Id of the Division to delete.
	 * @returns The deleted Division.
	 */
	public async deleteDivision(id: number): Promise<ResponseDivision> {
		const index = this.divisions.findIndex(d => d.id === id);
		if (index === -1) {
			throw new Error(`no such Division: ${id}`);
		}
		return this.divisions.splice(index, 1)[0];
	}

	/**
	 * Gets all Regions.
	 *
	 * @returns The requested Regions.
	 */
	public async getRegions(): Promise<Array<ResponseRegion>>;
	/**
	 * Gets a single Region.
	 *
	 * @param nameOrID The name (string) or ID (number) of the single Region to
	 * be returned.
	 * @returns The requested Region.
	 */
	public async getRegions(nameOrID: string | number): Promise<ResponseRegion>;
	/**
	 * Gets a Region or Regions.
	 *
	 * @param nameOrID If given, returns only the Region with the given name
	 * (string) or ID (number).
	 * @returns The requested Region or Regions.
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
	 * Replaces the current definition of a Region with the one given.
	 *
	 * @param region The new Region.
	 * @returns The updated Region.
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
	 * Creates a new Region.
	 *
	 * @param region The Region to create.
	 * @returns The created Region.
	 */
	public async createRegion(region: RequestRegion): Promise<ResponseRegion> {
		const reg = {
			divisionName: this.divisions.find(d => d.id === region.division)?.name ?? "",
			...region,
			id: ++this.lastID,
			lastUpdated: new Date()
		};
		this.regions.push(reg);
		return reg;
	}

	/**
	 * Deletes an existing Region.
	 *
	 * @param id Id of the Region to delete.
	 * @returns The deleted Region.
	 */
	public async deleteRegion(id: number | ResponseRegion): Promise<ResponseRegion> {
		const index = this.regions.findIndex(d => d.id === id);
		if (index === -1) {
			throw new Error(`no such Region: ${id}`);
		}
		return this.regions.splice(index, 1)[0];
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
		if(nameOrID) {
			let coordinate;
			switch (typeof nameOrID) {
				case "string":
					coordinate = this.coordinates.find(c=>c.name === nameOrID);
					break;
				case "number":
					coordinate = this.coordinates.find(c=>c.id === nameOrID);
			}
			if (!coordinate) {
				throw new Error(`no such Coordinate: ${nameOrID}`);
			}
			return coordinate;
		}
		return this.coordinates;
	}

	/**
	 * Replaces the current definition of a Coordinate with the one given.
	 *
	 * @param coordinate The new Coordinate.
	 * @returns The updated Coordinate.
	 */
	public async updateCoordinate(coordinate: ResponseCoordinate): Promise<ResponseCoordinate> {
		const id = this.coordinates.findIndex(c => c.id === coordinate.id);
		if (id === -1) {
			throw new Error(`no such Coordinate: ${coordinate.id}`);
		}
		this.coordinates[id] = coordinate;
		return coordinate;
	}

	/**
	 * Creates a new Coordinate.
	 *
	 * @param coordinate The Coordinate to create.
	 * @returns The created Coordinate.
	 */
	public async createCoordinate(coordinate: RequestCoordinate): Promise<ResponseCoordinate> {
		const crd = {
			...coordinate,
			id: ++this.lastID,
			lastUpdated: new Date()
		};
		this.coordinates.push(crd);
		return crd;
	}

	/**
	 * Deletes an existing Coordinate.
	 *
	 * @param id Id of the Coordinate to delete.
	 * @returns The deleted Coordinate.
	 */
	public async deleteCoordinate(id: number): Promise<ResponseCoordinate> {
		const index = this.coordinates.findIndex(c => c.id === id);
		if (index === -1) {
			throw new Error(`no such Coordinate: ${id}`);
		}
		return this.coordinates.splice(index, 1)[0];
	}

	/**
	 * Gets all ASNs.
	 *
	 * @returns All stored ASNs.
	 */
	public async getASNs(): Promise<Array<ResponseASN>>;
	/**
	 * Gets a single ASN.
	 *
	 * @param id The ID of the ASN to fetch.
	 * @returns The ASN with the given ID.
	 */
	public async getASNs(id: number): Promise<ResponseASN>;

	/**
	 * Gets all ASNs.
	 *
	 * @param id If given, returns only the ASN with the given ID.
	 * @returns An Array of ASNs objects - or a single ASN object if `id`
	 * was given.
	 */
	public async getASNs(id?: number): Promise<Array<ResponseASN> | ResponseASN> {
		if(id) {
			const asn = this.asns.find(a=>a.id === id);
			if (!asn) {
				throw new Error(`no such asn with id: ${id}`);
			}
			return asn;
		}
		return this.asns;
	}

	/**
	 * Replaces the current definition of a ASN with the one given.
	 *
	 * @param asn The new ASN.
	 * @returns The updated ASN.
	 */
	public async updateASN(asn: ResponseASN): Promise<ResponseASN> {
		const id = this.asns.findIndex(a => a.id === asn.id);
		if (id === -1) {
			throw new Error(`no such ASN: ${asn.id}`);
		}
		this.asns[id] = asn;
		return asn;
	}

	/**
	 * Creates a new ASN.
	 *
	 * @param asn The ASN to create.
	 * @returns The created ASN.
	 */
	public async createASN(asn: RequestASN): Promise<ResponseASN> {
		const sn = {
			...asn,
			cachegroup: this.cacheGroups.find(cg => cg.id === asn.cachegroupId)?.name ?? "",
			cachegroupId: asn.cachegroupId,
			id: ++this.lastID,
			lastUpdated: new Date()
		};
		this.asns.push(sn);
		return sn;
	}

	/**
	 * Deletes an existing ASN.
	 *
	 * @param asn The ASN to be deleted or ID of the ASN to delete..
	 */
	public async deleteASN(asn: ResponseASN | number): Promise<void> {
		const index = this.asns.findIndex(a => a.asn === asn);
		if (index === -1) {
			throw new Error(`no such asn: ${asn}`);
		}
		this.asns.splice(index, 1);
	}
}
