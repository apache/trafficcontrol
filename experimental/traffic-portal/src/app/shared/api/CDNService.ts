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

import { CDN } from "../../models";
import { APIService } from "./APIService";


/**
 * CDNService expose API functionality relating to CDNs.
 */
@Injectable()
export class CDNService extends APIService {
	public async getCDNs(id: number): Promise<CDN>;
	public async getCDNs(): Promise<Map<string, CDN>>;
	/**
	 * Gets one or all CDNs from Traffic Ops
	 *
	 * @param id The integral, unique identifier of a single CDN to be returned
	 * @returns Either a Map of CDN names to full CDN objects, or a single CDN, depending on whether `id` was
	 * 	passed.
	 * (In the event that `id` is passed but does not match any CDN, `null` will be emitted)
	 */
	public async getCDNs(id?: number): Promise<Map<string, CDN> | CDN> {
		const path = "cdns";
		if (id) {
			return this.get<[CDN]>(path, undefined, {id: String(id)}).toPromise().then(
				r => r[0]
			).catch(
				e => {
					console.error(`Failed to get CDN #${id}`, e);
					return {
						dnssecEnabled: false,
						domainName: "",
						id: -1,
						name: "",
					};
				}
			);
		}
		return this.get<Array<CDN>>(path).toPromise().then(
			r => new Map<string, CDN>(r.map(c=>[c.name, c]))
		).catch(
			e => {
				console.error("Failed to get CDNs:", e);
				return new Map();
			}
		);
	}

	/**
	 * Injects the Angular HTTP client service into the parent constructor.
	 *
	 * @param http The Angular HTTP client service.
	 */
	constructor(http: HttpClient) {
		super(http);
	}
}
