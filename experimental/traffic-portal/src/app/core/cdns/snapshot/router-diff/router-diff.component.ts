import { Component, EventEmitter, Input, type OnInit, Output } from "@angular/core";
import type { Subscribable, Unsubscribable } from "rxjs";
import type { Snapshot, SnapshotContentRouter } from "trafficops-types";

import { type RouterDifferences, routerDifferences } from "src/app/utils/snapshot.diffing";

import { Differences } from "../differences.class";

/**
 * Controller for a component that displays the set of differences between a
 * CDN's current Snapshot's `contentRouters` and its pending Snapshot's
 * `contentRouters`.
 */
@Component({
	selector: "tp-router-diff[snapshots]",
	styleUrls: ["./router-diff.component.scss"],
	templateUrl: "./router-diff.component.html",
})
export class RouterDiffComponent extends Differences<SnapshotContentRouter, RouterDifferences> implements OnInit {

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

				const diff = routerDifferences(current.contentRouters ?? {}, pending.contentRouters ?? {});
				this.numChanges = diff.changes;
				this.new = diff.new;
				this.deleted = diff.deleted;
				this.unchanged = diff.unchanged;
				this.changed = diff.changed;
			}
		);
	}
}
