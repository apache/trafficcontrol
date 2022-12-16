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

import { Directive, EventEmitter, type OnDestroy } from "@angular/core";
import type { Subscribable, Unsubscribable } from "rxjs";
import type { Snapshot } from "trafficops-types";

/**
 * Differences classes describe the differences between a property of a CDN
 * Snapshot between the CDN's current Snapshot and pending Snapshot.
 */
@Directive()
export abstract class Differences<T, U> implements OnDestroy {
	public abstract snapshots: Subscribable<{current: Snapshot; pending: Snapshot}>;
	public abstract changesPending: EventEmitter<number>;

	public new: readonly T[] = [];
	public deleted: readonly T[] = [];
	public unchanged: readonly T[] = [];
	public changed: readonly U[] = [];

	protected abstract subscription: Unsubscribable;

	private changes = 0;

	/**
	 * The number of changes in the pending Snapshot from the current Snapshot.
	 */
	public get numChanges(): number{
		return this.changes;
	}
	public set numChanges(num: number) {
		this.changes = num;
		this.changesPending.emit(num);
	}

	public abstract ngOnInit(): void;

	/**
	 * Angular lifecycle hook.
	 */
	public ngOnDestroy(): void {
		this.subscription.unsubscribe();
	}

	/**
	 * Returns a string describing how many changes are pending in a given
	 * collection.
	 *
	 * @param collection The collection of things being described, or just the
	 * size of said collection.
	 * @returns A string describing how many changes are pending in a given
	 * collection.
	 */
	 public pendingChangesStr(collection: number | readonly unknown[] | Record<PropertyKey, unknown>): string {
		let len;
		if (Array.isArray(collection)) {
			len = collection.length;
		} else if (typeof(collection) === "number") {
			len = collection;
		} else {
			len = Object.keys(collection).length;
		}
		if (len === 1) {
			return "1 change pending";
		}
		return `${len} changes pending`;
	}
}
