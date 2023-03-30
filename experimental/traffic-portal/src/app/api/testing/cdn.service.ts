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
import { CDNQueueResponse, RequestCDN, ResponseCDN } from "trafficops-types";

import { CDNService as ConcreteCDNService } from "../cdn.service";

import { APITestingService } from "./base-api.service";

/**
 * CDNService expose API functionality relating to CDNs.
 */
@Injectable()
export class CDNService extends APITestingService implements ConcreteCDNService{
	private lastID = 10;

	private readonly cdns = [
		{
			dnssecEnabled: false,
			domainName: "-",
			id: 1,
			lastUpdated: new Date(),
			name: "ALL",
		},
		{
			dnssecEnabled: false,
			domainName: "mycdn.test.test",
			id: 2,
			lastUpdated: new Date(),
			name: "test",
		}
	];

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
		if (id !== undefined) {
			const cdn = this.cdns.find(c => c.id === id);
			if (!cdn) {
				throw new Error(`no such CDN #${id}`);
			}
			return cdn;
		}
		return this.cdns;
	}

	/**
	 * Deletes a CDN.
	 *
	 * @param cdn The CDN to be deleted, or just its ID.
	 */
	public async deleteCDN(cdn: ResponseCDN | number): Promise<void> {
		const id = typeof cdn === "number" ? cdn : cdn.id;
		const idx = this.cdns.findIndex(c => c.id === id);
		if (idx < 0) {
			throw new Error(`no such CDN: #${id}`);
		}
		this.cdns.splice(idx, 1);
	}

	/**
	 * Creates a new CDN.
	 *
	 * @param cdn The CDN to create.
	 */
	public async createCDN(cdn: RequestCDN): Promise<ResponseCDN> {
		const c = {
			...cdn,
			id: ++this.lastID,
			lastUpdated: new Date(),
		};
		this.cdns.push(c);
		return c;
	}
	/**
	 * Queue updates to servers by a CDN
	 *
	 * @param cdn The CDN or id to queue server updates for
	 */
	public async queueServerUpdates(cdn: ResponseCDN | number): Promise<CDNQueueResponse> {
		const id = typeof cdn === "number" ? cdn : cdn.id;
		const realCDN = this.cdns.find(c => c.id === id);
		if (!realCDN) {
			throw new Error(`No CDN ${id}`);
		}
		return {
			action: "queue",
			cdnId: id
		};
	}

	/**
	 * Dequeue updates to servers by a CDN
	 *
	 * @param cdn The CDN or id to dequeue server updates for
	 */
	public async dequeueServerUpdates(cdn: ResponseCDN | number): Promise<CDNQueueResponse> {
		const id = typeof cdn === "number" ? cdn : cdn.id;
		const realCDN = this.cdns.find(c => c.id === id);
		if (!realCDN) {
			throw new Error(`No CDN ${id}`);
		}
		return {
			action: "dequeue",
			cdnId: id
		};
	}

	/**
	 * Replaces an existing CDN with the provided new definition of a
	 * CDN.
	 *
	 * @param id The if of the CDN being updated.
	 * @param cdn The new definition of the CDN.
	 */
	public async updateCDN(id: number, cdn: RequestCDN): Promise<ResponseCDN>;
	/**
	 * Replaces an existing CDN with the provided new definition of a
	 * CDN.
	 *
	 * @param cdn The full new definition of the CDN being
	 * updated.
	 */
	public async updateCDN(cdn: ResponseCDN): Promise<ResponseCDN>;
	/**
	 * Replaces an existing CDN with the provided new definition of a
	 * CDN.
	 *
	 * @param cdnOrID The full new definition of the CDN being
	 * updated, or just its ID.
	 * @param payload The new definition of the CDN. This is required if
	 * `cdnOrID` is an ID, and ignored otherwise.
	 */
	public async updateCDN(cdnOrID: ResponseCDN | number, payload?: RequestCDN): Promise<ResponseCDN> {
		let idx;
		let cdn;
		if (typeof cdnOrID === "number") {
			if (!payload) {
				throw new TypeError("invalid call signature - missing request payload");
			}
			idx = this.cdns.findIndex(c => c.id === cdnOrID);
			cdn = {
				...payload,
				id: ++this.lastID,
				lastUpdated: new Date(),
			};
		} else {
			idx = this.cdns.findIndex(c => c.id === cdnOrID.id);
			cdn = {
				...cdnOrID,
				lastUpdated: new Date()
			};
		}

		if (idx < 0) {
			throw new Error(`no such CDN: #${cdnOrID}`);
		}

		this.cdns[idx] = cdn;
		return cdn;
	}
}
