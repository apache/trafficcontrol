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
import type { RequestOrigin, RequestOriginResponse } from "trafficops-types";

import { APIService } from "./base-api.service";

/**
 * OriginService exposes API functionality relating to Origins.
 */
@Injectable()
export class OriginService extends APIService {
	/**
	 * Gets a specific Origin from Traffic Ops.
	 *
	 * @param idOrName Either the integral, unique identifier (number) or name
	 * (string) of the Origin to be returned.
	 * @returns The requested Origin.
	 */
	public async getOrigins(
		idOrName: number | string
	): Promise<RequestOriginResponse>;
	/**
	 * Gets Origins from Traffic Ops.
	 *
	 * @returns An Array of all Origins from Traffic Ops.
	 */
	public async getOrigins(): Promise<Array<RequestOriginResponse>>;
	/**
	 * Gets one or all Origins from Traffic Ops.
	 *
	 * @param idOrName Optionally the integral, unique identifier (number) or
	 * name (string) of a single Origin to be returned.
	 * @returns The requested Origin(s).
	 */
	public async getOrigins(
		idOrName?: number | string
	): Promise<RequestOriginResponse | Array<RequestOriginResponse>> {
		const path = "origins";
		if (idOrName !== undefined) {
			let params;
			switch (typeof idOrName) {
				case "string":
					params = { name: idOrName };
					break;
				case "number":
					params = { id: idOrName };
			}
			const r = await this.get<[RequestOriginResponse]>(
				path,
				undefined,
				params
			).toPromise();
			if (r.length !== 1) {
				throw new Error(
					`Traffic Ops responded with ${r.length} Origins by identifier ${idOrName}`
				);
			}
			return r[0];
		}
		return this.get<Array<RequestOriginResponse>>(path).toPromise();
	}

	/**
	 * Deletes an existing Origin.
	 *
	 * @param originOrId The ID of the Origin to delete.
	 * @returns The deleted Origin.
	 */
	public async deleteOrigin(originOrId: number | RequestOriginResponse): Promise<RequestOriginResponse> {
		const id = typeof originOrId === "number" ? originOrId : originOrId.id;
		return this.delete<RequestOriginResponse>(`origins?id=${id}`).toPromise();
	}

	/**
	 * Creates a new Origin.
	 *
	 * @param origin The Origin to create.
	 * @returns The created Origin.
	 */
	public async createOrigin(origin: RequestOrigin): Promise<RequestOriginResponse> {
		return this.post<RequestOriginResponse>("origins", origin).toPromise();
	}

	/**
	 * Replaces the current definition of an Origin with the one given.
	 *
	 * @param origin The new Origin.
	 * @returns The updated Origin.
	 */
	public async updateOrigin(
		origin: RequestOriginResponse
	): Promise<RequestOriginResponse> {
		const path = `origins?id=${origin.id}`;
		return this.put<RequestOriginResponse>(path, origin).toPromise();
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
