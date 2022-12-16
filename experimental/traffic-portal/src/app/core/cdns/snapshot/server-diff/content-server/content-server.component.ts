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

import { Component, Input } from "@angular/core";
import { SnapshotContentServer } from "trafficops-types";

/**
 * Controller for displaying a simple Snapshot content server.
 */
@Component({
	selector: "tp-content-server[server]",
	styleUrls: ["./content-server.component.scss"],
	templateUrl: "./content-server.component.html",
})
export class ContentServerComponent {

	@Input() public server!: SnapshotContentServer;
	@Input() public kind: "unchanged" | "new" | "deleted" = "unchanged";

	/**
	 * Yields iterable tuples for the described server's `deliveryServices`. If
	 * the server isn't an edge (or hasn't been set yet by Angular), this simply
	 * returns an empty collection.
	 *
	 * @returns Tuples where the first element is the DS's XMLID and the second
	 * is the DS value.
	 */
	public deliveryServices(): Iterable<[string, Array<string>]> {
		if (!this.server || this.server.type === "MID") {
			return [];
		}
		return Object.entries(this.server.deliveryServices);
	}

}
