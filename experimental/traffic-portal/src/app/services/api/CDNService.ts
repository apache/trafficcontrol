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

import { CDN } from "../../models";
import { APIService } from "./apiservice";


/**
 * CDNService expose API functionality relating to CDNs.
 */
@Injectable({providedIn: "root"})
export class CDNService extends APIService {
	public getCDNs(id: number): Observable<CDN>;
	public getCDNs(): Observable<Map<string, CDN>>;
	/**
	 * Gets one or all CDNs from Traffic Ops
	 *
	 * @param id The integral, unique identifier of a single CDN to be returned
	 * @returns An Observable that will emit either a Map of CDN names to full CDN objects, or a single CDN, depending on whether `id` was
	 * 	passed.
	 * (In the event that `id` is passed but does not match any CDN, `null` will be emitted)
	 */
	public getCDNs(id?: number): Observable<Map<string, CDN> | CDN> {
		const path = "cdns";
		if (id) {
			return this.get<[CDN]>(path, undefined, {id: String(id)}).pipe(map(
				r => r[0]
			));
		}
		return this.get<Array<CDN>>(path).pipe(map(
			r => new Map<string, CDN>(r.map(c=>[c.name, c]))
		));
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
