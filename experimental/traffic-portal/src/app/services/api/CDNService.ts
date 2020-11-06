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

import { Observable } from "rxjs";
import { map } from "rxjs/operators";

import { APIService } from "./apiservice";

import { CDN } from "../../models";

/**
 * CDNService expose API functionality relating to CDNs.
 */
@Injectable({providedIn: "root"})
export class CDNService extends APIService {
	/**
	 * Gets one or all CDNs from Traffic Ops
	 * @param id The integral, unique identifier of a single CDN to be returned
	 * @returns An Observable that will emit either a Map of CDN names to full CDN objects, or a single CDN, depending on whether `id` was
	 * 	passed.
	 * (In the event that `id` is passed but does not match any CDN, `null` will be emitted)
	 */
	public getCDNs (id?: number): Observable<Map<string, CDN> | CDN | undefined> {
		const path = `/api/${this.API_VERSION}/cdns`;
		if (id) {
			return this.get(`${path}?id=${id}`).pipe(map(
				r => {
					for (const c of (r.body.response as Array<CDN>)) {
						if (c.id === id) {
							return c;
						}
					}
				}
			));
		}
		return this.get(path).pipe(map(
			r => {
				const ret = new Map<string, CDN>();
				for (const c of (r.body.response as Array<CDN>)) {
					ret.set(c.name, c);
				}
				return ret;
			}
		));
	}

	constructor(http: HttpClient) {
		super(http);
	}
}
