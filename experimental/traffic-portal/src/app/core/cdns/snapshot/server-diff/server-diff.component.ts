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

import { Component, EventEmitter, type OnInit, type OnDestroy, Input, Output } from "@angular/core";
import type { Subscribable, Unsubscribable } from "rxjs";
import type { Snapshot, SnapshotContentServer } from "trafficops-types";

import { serverDifferences, type ServerDifferences } from "src/app/utils/snapshot.diffing";

import { Differences } from "../differences.class";

/**
 * Controller for a component that displays the differences between a CDN's
 * current Snapshot's `contentServers` and its pending Snapshot's
 * `contentServers`.
 */
@Component({
	selector: "tp-server-diff[snapshots]",
	styleUrls: ["./server-diff.component.scss"],
	templateUrl: "./server-diff.component.html",
})
export class ServerDiffComponent extends Differences<SnapshotContentServer, ServerDifferences> implements OnInit, OnDestroy {

	@Input() public snapshots!: Subscribable<{current: Snapshot; pending: Snapshot}>;
	@Output() public changesPending = new EventEmitter<number>();

	protected subscription!: Unsubscribable;

	constructor() {
		super();
		this.changesPending.emit(0);
	}

	/**
	 * Angular lifecycle hook.
	 */
	 public ngOnInit(): void {
		this.subscription = this.snapshots.subscribe(
			({current, pending}) => {
				const diff = serverDifferences(current.contentServers ?? {}, pending.contentServers ?? {});
				console.log("diff:", diff);
				this.numChanges = diff.changes;
				this.new = diff.new;
				this.deleted = diff.deleted;
				this.unchanged = diff.unchanged;
				this.changed = diff.changed;
			}
		);
	}

	/**
	 * Gets an iterable list of the Delivery Service Assignments of a given
	 * server.
	 *
	 * @param dss The Delivery Services to be expanded.
	 * @returns A list that can be iterated of the passed DS Assignments.
	 */
	public expandDSAssignments<T>(dss: Record<string, T>): Array<[string, T]> {
		return Object.entries(dss);
	}
}
