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
import type { ResponseCoordinate } from "trafficops-types";

/**
 * CoordinateService exposes API functionality relating to Coordinates.
 */
@Injectable()
export class CoordinateService {
	private readonly coordinates: Array<ResponseCoordinate> = [
		{
			id: 1,
			lastUpdated: new Date(),
			latitude: 1.0,
			longitude: -1.0,
			name: "test_coordinate",
		},
	];

	public async getCoordinates(): Promise<Array<ResponseCoordinate>>;
	public async getCoordinates(
		nameOrID: string | number
	): Promise<ResponseCoordinate>;

	/**
	 * Gets one or all Coordinates from Traffic Ops.
	 *
	 * @param nameOrID If given, returns only the ResponseCoordinate with the given name
	 * (string) or ID (number).
	 * @returns An Array of ResponseCoordinate objects - or a single ResponseCoordinate object if 'nameOrID'
	 * was given.
	 */
	public async getCoordinates(
		nameOrID?: string | number
	): Promise<Array<ResponseCoordinate> | ResponseCoordinate> {
		if (nameOrID) {
			let coordinate;
			switch (typeof nameOrID) {
				case "string":
					coordinate = this.coordinates.find(
						(d) => d.name === nameOrID
					);
					break;
				case "number":
					coordinate = this.coordinates.find(
						(d) => d.id === nameOrID
					);
			}
			if (!coordinate) {
				throw new Error(`no such Coordinate: ${nameOrID}`);
			}
			return coordinate;
		}
		return this.coordinates;
	}
}
