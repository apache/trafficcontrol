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
import type { ResponseCDN } from "trafficops-types";

import { APIService } from "./base-api.service";

/**
 * CDNService expose API functionality relating to CDNs.
 */
@Injectable()
export class CDNService extends APIService {

	constructor(http: HttpClient) {
		super(http);
	}

	public async getCDNs(id: number): Promise<ResponseCDN>;
	public async getCDNs(): Promise<Array<ResponseCDN>>;
	/**
	 * Gets one or all CDNs from Traffic Ops
	 *
	 * @param id The integral, unique identifier of a single CDN to be returned
	 * @returns Either a Map of CDN names to full CDN objects, or a single CDN, depending on whether `id` was
	 * 	passed.
	 * (In the event that `id` is passed but does not match any CDN, `null` will be emitted)
	 */
	public async getCDNs(id?: number): Promise<Array<ResponseCDN> | ResponseCDN> {
		const path = "cdns";
		if (id) {
			const cdn = await this.get<[ResponseCDN]>(path, undefined, {id: String(id)}).toPromise();
			if (cdn.length !== 1) {
				throw new Error(`${cdn.length} CDNs found by ID ${id}`);
			}
			return cdn;
		}
		return this.get<Array<ResponseCDN>>(path).toPromise();
	}
}
