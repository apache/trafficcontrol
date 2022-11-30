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

import { Component, EventEmitter, type OnInit } from "@angular/core";
import { ActivatedRoute } from "@angular/router";
import type { ResponseCDN, Snapshot, SnapshotStatsSection } from "trafficops-types";

import { CDNService } from "src/app/api";
import { buildDiff, type DiffVal, type Diff } from "src/app/utils/snapshot.diffing";

/**
 * Produces a new blank stats section of a CDN Snapshot.
 *
 * @returns A CDN Snapshot 'stats' section with all of the properties
 * initialized to 'blank' or 'zero' values.
 */
function newStatsSection(): SnapshotStatsSection {
	return {
		// eslint-disable-next-line @typescript-eslint/naming-convention
		CDN_name: "",
		date: new Date(0),
		// eslint-disable-next-line @typescript-eslint/naming-convention
		tm_host: "",
		// eslint-disable-next-line @typescript-eslint/naming-convention
		tm_path: "",
		// eslint-disable-next-line @typescript-eslint/naming-convention
		tm_user: "",
		// eslint-disable-next-line @typescript-eslint/naming-convention
		tm_version: ""
	};
}

/**
 * The controller for the page that shows Snapshot diffs
 */
@Component({
	selector: "tp-snapshot",
	styleUrls: ["./snapshot.component.scss"],
	templateUrl: "./snapshot.component.html",
})
export class SnapshotComponent implements OnInit {

	public cdn: ResponseCDN = {
		dnssecEnabled: false,
		domainName: "",
		id: -1,
		lastUpdated: new Date(0),
		name: "",
	};
	public currentSnapshot: Snapshot = {
		config: {},
		contentRouters: {},
		contentServers: {},
		deliveryServices: {},
		edgeLocations: {},
		monitors: {},
		stats: newStatsSection(),
		trafficRouterLocations: {}
	};
	public pendingSnapshot: Snapshot = {
		config: {},
		contentRouters: {},
		contentServers: {},
		deliveryServices: {},
		edgeLocations: {},
		monitors: {},
		stats: newStatsSection(),
		trafficRouterLocations: {}
	};

	public snaps = new EventEmitter<{current: Snapshot; pending: Snapshot}>();

	public configDiff: Diff = {
		fields: {},
		num: 0
	};

	private routerChangesPending = 0;
	private serverChangesPending = 0;

	/**
	 * The total number of changes to the Snapshot as a whole.
	 */
	public get totalChangesPending(): number {
		return this.routerChangesPending + this.serverChangesPending + this.configDiff.num;
	}

	/**
	 * An easily iterable collection of tuples containing the name of the field,
	 * and the differences between the old and new values.
	 */
	public get configDiffFields(): Array<[string, DiffVal]> {
		return Object.entries(this.configDiff.fields);
	}

	constructor(private readonly route: ActivatedRoute, private readonly api: CDNService) { }

	/**
	 * Angular lifecycle hook.
	 */
	public async ngOnInit(): Promise<void> {
		const name = this.route.snapshot.paramMap.get("name");
		if (name === null) {
			console.error("missing required route parameter 'name'");
			return;
		}

		const p = [
			this.api.getCDNs(name),
			this.api.getCurrentSnapshot(name),
			this.api.getPendingSnapshot(name)
		] as const;
		[this.cdn, this.currentSnapshot, this.pendingSnapshot] = await Promise.all(p);
		this.snaps.emit({
			current: this.currentSnapshot,
			pending: this.pendingSnapshot
		});
		const [currentConf, pendingConf] = [
			this.currentSnapshot.config ?? {},
			this.pendingSnapshot.config ?? {}
		];

		this.configDiff = buildDiff(currentConf, pendingConf);
	}

	/**
	 * Sets the amount of changes to a particular section of the Snapshot.
	 *
	 * @param category The category of changes which will be set to the given
	 * amount. This is any of the categories of CDN Snapshot sections that are
	 * implemented by sub-components.
	 * @param amt The amount to set for the given category.
	 */
	public setChangeAmount(category: "router" | "server", amt: number): void {
		switch (category) {
			case "router":
				this.routerChangesPending = amt;
				break;
			case "server":
				this.serverChangesPending = amt;
				break;
		}
	}

}
