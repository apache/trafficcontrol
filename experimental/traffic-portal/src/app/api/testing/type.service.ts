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
import {RequestType, TypeFromResponse} from "trafficops-types";

/** The allowed values for the 'useInTables' query parameter of GET requests to /types. */
type UseInTable = "cachegroup" |
"server" |
"deliveryservice" |
"to_extension" |
"federation_resolver" |
"regex" |
"staticdnsentry" |
"steering_target";

/**
 * TypeService exposes API functionality relating to Types.
 */
@Injectable()
export class TypeService {
	private lastID = 20;
	private readonly types = [
		{
			description: "Mid Logical Location",
			id: 1,
			lastUpdated: new Date(),
			name: "MID_LOC",
			useInTable: "cachegroup"
		},
		{
			description: "Edge Logical Location",
			id: 2,
			lastUpdated: new Date(),
			name: "EDGE_LOC",
			useInTable: "cachegroup"
		},
		{
			description: "Origin Logical Site",
			id: 3,
			lastUpdated: new Date(),
			name: "ORG_LOC",
			useInTable: "cachegroup"
		},
		{
			description: "Traffic Control Component Location",
			id: 4,
			lastUpdated: new Date(),
			name: "TC_LOC",
			useInTable: "cachegroup"
		},
		{
			description: "Traffic Router Logical Location",
			id: 15,
			lastUpdated: new Date(),
			name: "TR_LOC",
			useInTable: "cachegroup"
		},
		{
			description: "No Content Routing - arbitrary remap at the edge, no Traffic Router config",
			id: 5,
			lastUpdated: new Date(),
			name: "ANY_MAP",
			useInTable: "deliveryservice"
		},
		{
			description: "Client-Controlled Steering Delivery Service",
			id: 6,
			lastUpdated: new Date(),
			name: "CLIENT_STEERING",
			useInTable: "deliveryservice"
		},
		{
			description: "DNS Content Routing",
			id: 7,
			lastUpdated: new Date(),
			name: "DNS",
			useInTable: "deliveryservice"
		},
		{
			description: "DNS Content routing, RAM cache, Local",
			id: 8,
			lastUpdated: new Date(),
			name: "DNS_LIVE",
			useInTable: "deliveryservice"
		},
		{
			description: "DNS Content routing, RAM cache, National",
			id: 9,
			lastUpdated: new Date(),
			name: "DNS_LIVE_NATNL",
			useInTable: "deliveryservice"
		},
		{
			description: "HTTP Content Routing",
			id: 10,
			lastUpdated: new Date(),
			name: "HTTP",
			useInTable: "deliveryservice"
		},
		{
			description: "HTTP Content routing cache in RAM",
			id: 11,
			lastUpdated: new Date(),
			name: "HTTP_LIVE",
			useInTable: "deliveryservice"
		},
		{
			description: "HTTP Content routing, RAM cache, National",
			id: 12,
			lastUpdated: new Date(),
			name: "HTTP_LIVE_NATNL",
			useInTable: "deliveryservice"
		},
		{
			description: "HTTP Content Routing, no caching",
			id: 13,
			lastUpdated: new Date(),
			name: "HTTP_NO_CACHE",
			useInTable: "deliveryservice"
		},
		{
			description: "Steering Delivery Service",
			id: 14,
			lastUpdated: new Date(),
			name: "STEERING",
			useInTable: "deliveryservice"
		},
		{
			description: "Edge Cache",
			id: 15,
			lastUpdated: new Date(),
			name: "EDGE",
			useInTable: "server"
		}
	];

	/**
	 * Gets all Types.
	 *
	 * @returns The requested Types.
	 */
	public async getTypes(): Promise<Array<TypeFromResponse>>;
	/**
	 * Gets a specific Type.
	 *
	 * @param idOrName Either the integral, unique identifier (number) or name
	 * (string) of the Type to be returned.
	 * @returns The requested Type.
	 */
	public async getTypes(idOrName: number | string): Promise<TypeFromResponse>;
	/**
	 * Gets one or all Types.
	 *
	 * @param idOrName Optionally the integral, unique identifier (number) or
	 * name (string) of a single Type to be returned.
	 * @returns The requested Type(s).
	 */
	public async getTypes(idOrName?: number | string): Promise<TypeFromResponse | Array<TypeFromResponse>> {
		if (idOrName !== undefined) {
			let type;
			switch (typeof idOrName) {
				case "string":
					type = this.types.filter(t=>t.name === idOrName)[0];
					break;
				case "number":
					type = this.types.filter(t=>t.id === idOrName)[0];
			}
			if (!type) {
				throw new Error(`no such Type: ${idOrName}`);
			}
			return type;
		}
		return this.types;
	}

	/**
	 * Gets all Types used by specific database table.
	 *
	 * @param useInTable The database table for which to retrieve Types.
	 * @returns The requested Types.
	 */
	public async getTypesInTable(useInTable: UseInTable): Promise<Array<TypeFromResponse>> {
		return this.types.filter(t=>t.useInTable === useInTable);
	}

	/**
	 * Gets all Server Types.
	 *
	 * @returns All Types that have 'server' as their 'useInTable'.
	 */
	public async getServerTypes(): Promise<Array<TypeFromResponse>> {
		return this.getTypesInTable("server");
	}

	/**
	 * Deletes an existing type.
	 *
	 * @param id Id of the type to delete.
	 * @returns The deleted type.
	 */
	public async deleteType(id: number): Promise<TypeFromResponse> {
		const index = this.types.findIndex(t => t.id === id);
		if (index === -1) {
			throw new Error(`no such Type: ${id}`);
		}
		return this.types.splice(index, 1)[0];
	}

	/**
	 * Creates a new type.
	 *
	 * @param type The type to create.
	 * @returns The created type.
	 */
	public async createType(type: RequestType): Promise<TypeFromResponse> {
		const t = {
			...type,
			id: ++this.lastID,
			lastUpdated: new Date()
		};
		this.types.push(t);
		return t;
	}
}
