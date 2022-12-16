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
import type { CDN, ResponseCDN, Snapshot } from "trafficops-types";

import { APIService } from "./base-api.service";

/**
 * CDNService expose API functionality relating to CDNs.
 */
@Injectable()
export class CDNService extends APIService {

	constructor(http: HttpClient) {
		super(http);
	}

	public async getCDNs(id: number | string): Promise<ResponseCDN>;
	public async getCDNs(): Promise<Map<string, ResponseCDN>>;
	/**
	 * Gets one or all CDNs from Traffic Ops
	 *
	 * @param id The integral, unique identifier of a single CDN to be returned
	 * @returns Either a Map of CDN names to full CDN objects, or a single CDN, depending on whether `id` was
	 * 	passed.
	 * (In the event that `id` is passed but does not match any CDN, `null` will be emitted)
	 */
	public async getCDNs(id?: number | string): Promise<Map<string, ResponseCDN> | ResponseCDN> {
		const path = "cdns";
		if (typeof(id) === "number") {
			const r = await this.get<[ResponseCDN]>(path, undefined, {id: String(id)}).toPromise();
			if (r.length !== 1) {
				throw new Error(`got ${r.length} CDNs with ID ${id}`);
			}
			return r[0];
		}
		if (typeof(id) === "string") {
			const r = await this.get<[ResponseCDN]>(path, undefined, {name: id}).toPromise();
			if (r.length !== 1) {
				throw new Error(`got ${r.length} CDNs with Name ${id}`);
			}
			return r[0];
		}
		const cdns = await this.get<Array<ResponseCDN>>(path).toPromise();
		return new Map(cdns.map(c=>[c.name, c]));
	}

	/**
	 * Gets the *current* Snapshot for a given CDN.
	 *
	 * @param cdn The CDN for which to fetch a Snapshot.
	 * @returns The current Snapshot of the requested CDN.
	 */
	 public async getCurrentSnapshot(cdn: CDN | string): Promise<Snapshot> {
		const name = typeof(cdn) === "string" ? cdn : cdn.name;
		return this.get<Snapshot>(`cdns/${name}/snapshot`).toPromise();
	}

	/**
	 * Gets the *pending* Snapshot for a given CDN.
	 *
	 * @param cdn The CDN for which to fetch a Snapshot.
	 * @returns The current Snapshot of the requested CDN.
	 */
	 public async getPendingSnapshot(cdn: CDN | string): Promise<Snapshot> {
		const name = typeof(cdn) === "string" ? cdn : cdn.name;
		return this.get<Snapshot>(`cdns/${name}/snapshot/new`).toPromise();
	}
}
