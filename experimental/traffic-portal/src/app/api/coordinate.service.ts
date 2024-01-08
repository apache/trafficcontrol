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
import type { ResponseCoordinate } from "trafficops-types";

import { APIService } from "./base-api.service";

/**
 * CoordinateService exposes API functionality relating to Coordinates.
 */
@Injectable()
export class CoordinateService extends APIService {
	/**
	 * Gets a specific Coordinate from Traffic Ops.
	 *
	 * @param idOrName Either the integral, unique identifier (number) or name
	 * (string) of the Coordinate to be returned.
	 * @returns The requested Coordinate.
	 */
	public async getCoordinates(
		idOrName: number | string
	): Promise<ResponseCoordinate>;
	/**
	 * Gets Coordinates from Traffic Ops.
	 *
	 * @returns An Array of all Coordinates from Traffic Ops.
	 */
	public async getCoordinates(): Promise<Array<ResponseCoordinate>>;

	/**
	 * Gets one or all Coordinates from Traffic Ops.
	 *
	 * @param idOrName Optionally the integral, unique identifier (number) or
	 * name (string) of a single Coordinate to be returned.
	 * @returns The requested Coordinate(s).
	 */
	public async getCoordinates(
		idOrName?: number | string
	): Promise<ResponseCoordinate | Array<ResponseCoordinate>> {
		const path = "coordinates";
		if (idOrName !== undefined) {
			let params;
			switch (typeof idOrName) {
				case "string":
					params = { name: idOrName };
					break;
				case "number":
					params = { id: idOrName };
			}
			const r = await this.get<[ResponseCoordinate]>(
				path,
				undefined,
				params
			).toPromise();
			if (r.length !== 1) {
				throw new Error(
					`Traffic Ops responded with ${r.length} Coordinates by identifier ${idOrName}`
				);
			}
			return r[0];
		}
		return this.get<Array<ResponseCoordinate>>(path).toPromise();
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
