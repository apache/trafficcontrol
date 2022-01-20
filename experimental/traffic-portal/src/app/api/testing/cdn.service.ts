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

import { CDN } from "src/app/models";

/**
 * CDNService expose API functionality relating to CDNs.
 */
@Injectable()
export class CDNService {

	private readonly cdns = new Map([
		[
			"ALL",
			{
				dnssecEnabled: false,
				domainName: "-",
				id: 1,
				lastUpdated: new Date(),
				name: "ALL"
			}
		],
		[
			"test",
			{
				dnssecEnabled: false,
				domainName: "mycdn.test.test",
				id: 1,
				lastUpdated: new Date(),
				name: "test"
			}
		]
	]);

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
		if (id !== undefined) {
			const cdn = Array.from(this.cdns.values()).filter(c=>c.id===id)[0];
			if (!cdn) {
				throw new Error(`no such CDN #${id}`);
			}
			return cdn;
		}
		return this.cdns;
	}
}
