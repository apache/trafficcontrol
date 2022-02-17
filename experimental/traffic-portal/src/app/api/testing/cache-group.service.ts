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

import type { CacheGroup } from "src/app/models";

/**
 * CDNService expose API functionality relating to CDNs.
 */
@Injectable()
export class CacheGroupService {

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
}
