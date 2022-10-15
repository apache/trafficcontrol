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

import { SnapshotContentRouter, SnapshotContentServer } from "trafficops-types";

/**
 * A value is "diffable" iff it's a primitive type and not a complex object or
 * collection.
 */
export type Diffable = string | number | boolean | null;

/**
 * A single difference in value.
 */
export interface DiffVal<T extends string | symbol | number | boolean | null | bigint = Diffable> {
	newValue: T | undefined;
	oldValue: T | undefined;
};

/**
 * A collection of differences between two objects.
 */
export interface Diff {
	/**
	 * All of the fields that exist in either the new or pending configuration.
	 */
	readonly fields: Readonly<Record<string, Readonly<DiffVal>>>;
	/**
	 * Cached number of differences; will be equal to the number of entries in
	 * `fields`.
	 */
	readonly num: number;
}

/**
 * Builds a set of differences between two "diffable" objects.
 *
 * @param current The current object; its values are turned into DiffVal
 * `oldValue` property values.
 * @param pending The new object; its values are turned into DiffVal `newValue`
 * property values.
 * @returns A summary of the differences between `current` and `pending`.
 */
export function buildDiff(current: Record<string, Diffable>, pending: Record<string, Diffable>): Diff {
	const diff = {
		fields: {} as Record<string, DiffVal>,
		num: 0
	};
	const currentHas = new Set<string>();
	for (const [name, oldValue] of Object.entries(current)) {
		currentHas.add(name);
		const newValue = pending[name];
		if (oldValue !== newValue) {
			++diff.num;
		}
		diff.fields[name] = {
			newValue,
			oldValue
		};
	}

	for (const [name, newValue] of Object.entries(pending)) {
		if (!currentHas.has(name)) {
			++diff.num;
			diff.fields[name] = {
				newValue,
				oldValue: undefined
			};
		}
	}

	return diff;
}

/**
 * Expresses the differences between two unordered collections (that were
 * represented as an Array due to limitations of JSON encoding).
 */
interface ArrayDiff<T extends Diffable> {
	readonly deleted: Readonly<Set<T>>;
	readonly new: Readonly<Set<T>>;
	readonly unchanged: Readonly<Set<T>>;
}

/**
 * Expresses the differences between two ordered collections
 */
interface OrderedArrayDiff<T extends Diffable> {
	readonly changes: readonly DiffVal<T>[];
	readonly changed: boolean;
}

/**
 * Finds the differences between two homogeneously typed collections, where
 * changes to the ordering of the elements of said collections are considered
 * differences.
 *
 * @param current The current ordered collection.
 * @param pending The pending ordered collection.
 * @returns A summary of the differences between the two collections.
 */
function orderedArrayDiff<T extends Diffable>(current: readonly T[], pending: readonly T[]): OrderedArrayDiff<T> {
	let changes = new Array<DiffVal<T>>();
	let changed = false;
	let i;
	for (i = 0; i < current.length; ++i) {
		const cVal = current[i];
		const pVal = pending[i];
		changes[i] = {
			newValue: pVal,
			oldValue: cVal
		};
		if (cVal !== pVal) {
			changed = true;
		}
	}

	if (i < pending.length) {
		changes = changes.concat(pending.slice(i).map(v=>({newValue: v, oldValue: undefined})));
		changed = true;
	}

	return {changed, changes};
}

/**
 * Finds the difference between two homogenous arrays, where order doesn't
 * matter.
 *
 * @param current The first/current/old collection.
 * @param pending The second/pending/new collection.
 * @param ordered If given and `true`, the differences will consider order
 * changes to be differences.
 * @returns A summary of the differences between the two collections.
 */
export function arrayDiff<T extends Diffable>(current: readonly T[], pending: readonly T[], ordered?: false): ArrayDiff<T>;
/**
 * Finds the difference between two homogenous arrays, where order does matter.
 *
 * @param current The first/current/old collection.
 * @param pending The second/pending/new collection.
 * @param ordered If given and `true`, the differences will consider order
 * changes to be differences.
 * @returns A summary of the differences between the two collections.
 */
export function arrayDiff<T extends Diffable>(current: readonly T[], pending: readonly T[], ordered: true): OrderedArrayDiff<T>;
/**
 * Finds the difference between two homogenous arrays, where order does or
 * doesn't matter according to the value of `ordered`.
 *
 * @param current The first/current/old collection.
 * @param pending The second/pending/new collection.
 * @param ordered If given and `true`, the differences will consider order
 * changes to be differences.
 * @returns A summary of the differences between the two collections.
 */
export function arrayDiff<T extends Diffable>(
	current: readonly T[],
	pending: readonly T[],
	ordered = false
): ArrayDiff<T> | OrderedArrayDiff<T> {
	if (ordered) {
		return orderedArrayDiff(current, pending);
	}
	const deleted = new Set<T>();
	const unchanged = new Set<T>();
	const newVals = new Set(pending);
	for (const c of current) {
		if (newVals.has(c)) {
			unchanged.add(c);
			newVals.delete(c);
		} else {
			deleted.add(c);
		}
	}
	return {
		deleted,
		new: newVals,
		unchanged
	};
}

/**
 * Contains the actual differences in value between two Traffic Routers.
 */
export type RouterDifferences = {
	[k in keyof SnapshotContentRouter]: DiffVal;
};

/**
 * A summary of the differences between two Traffic Routers.
 */
interface RouterDiff {
	/** Whether there's any difference between the two. */
	changed: boolean;
	diff: RouterDifferences;
}

/**
 * Finds the differences between two separate definitions of the same Traffic
 * Router.
 *
 * @param current The old definition of the Traffic Router.
 * @param pending The new definition of the Traffic Router.
 * @returns A summary of the differences between `current` and `pending`.
 */
function diffRouters(current: SnapshotContentRouter, pending: SnapshotContentRouter): RouterDiff  {
	let changed = false;
	const diff = {} as {[k in keyof SnapshotContentRouter]: DiffVal};
	if (current["api.port"] !== pending["api.port"]) {
		changed = true;
	}
	diff["api.port"] = {
		newValue: pending["api.port"],
		oldValue: current["api.port"]
	};

	if (!changed && current["secure.api.port"] !== pending["secure.api.port"]) {
		changed = true;
	}
	diff["secure.api.port"] = {
		newValue: pending["secure.api.port"],
		oldValue: current["secure.api.port"]
	};
	if (!changed && current.fqdn !== pending.fqdn) {
		changed = true;
	}
	diff.fqdn = {
		newValue: pending.fqdn,
		oldValue: current.fqdn
	};
	if (!changed && current.httpsPort !== pending.httpsPort) {
		changed = true;
	}
	diff.httpsPort = {
		newValue: pending.httpsPort,
		oldValue: current.httpsPort
	};
	if (!changed && current.ip !== pending.ip) {
		changed = true;
	}
	diff.ip = {
		newValue: pending.ip,
		oldValue: current.ip
	};
	if (!changed && current.ip6 !== pending.ip6) {
		changed = true;
	}
	diff.ip6 = {
		newValue: pending.ip6,
		oldValue: current.ip6
	};
	if (!changed && current.location !== pending.location) {
		changed = true;
	}
	diff.location = {
		newValue: pending.location,
		oldValue: current.location
	};
	if (!changed && current.port !== pending.port) {
		changed = true;
	}
	diff.port = {
		newValue: pending.port,
		oldValue: current.port
	};
	if (!changed && current.profile !== pending.profile) {
		changed = true;
	}
	diff.profile = {
		newValue: pending.profile,
		oldValue: current.profile
	};
	if (!changed && current.status !== pending.status) {
		changed = true;
	}
	diff.status = {
		newValue: pending.status,
		oldValue: current.status
	};

	return {changed, diff};
}

/**
 * Summarizes the differences between two sets of CDN Snapshot Traffic Routers.
 */
export interface RoutersDiff {
	readonly changed: readonly Readonly<RouterDifferences>[];
	readonly changes: number;
	readonly deleted: readonly Readonly<SnapshotContentRouter>[];
	readonly new: readonly Readonly<SnapshotContentRouter>[];
	readonly unchanged: readonly Readonly<SnapshotContentRouter>[];
}

/**
 * Finds the differences between two sets of Snapshot Traffic Routers.
 *
 * @param current The current set of Traffic Routers.
 * @param pending The new set of Traffic Routers.
 * @returns A summary of the differences between the router sets.
 */
export function routerDifferences(
	current: Record<string, SnapshotContentRouter>,
	pending: Record<string, SnapshotContentRouter>
): RoutersDiff {
	const currentHas = new Set<string>();
	const diffs = {
		changed: new Array<RouterDifferences>(),
		changes: 0,
		deleted: new Array<SnapshotContentRouter>(),
		new: new Array<SnapshotContentRouter>(),
		unchanged: new Array<SnapshotContentRouter>()
	};

	for (const [k, v] of Object.entries(current)) {
		currentHas.add(k);
		const p = pending[k];
		if (!p) {
			diffs.deleted.push(v);
			++diffs.changes;
		} else {
			const diff = diffRouters(v, p);
			if (diff.changed) {
				++diffs.changes;
				diffs.changed.push(diff.diff);
			} else {
				diffs.unchanged.push(v);
			}
		}
	}

	diffs.new = Object.entries(pending).filter(r=>!currentHas.has(r[0])).map(r=>r[1]);
	diffs.changes += diffs.new.length;
	return diffs;
}

/**
 * Expresses the differences between two CDN Snapshot cache servers' respective
 * `deliveryServices` properties.
 */
interface ServerDSesDiff {
	changed: {
		[xmlID: string]: readonly DiffVal<string>[];
	};
	deleted: {
		[xmlID: string]: readonly string[];
	};
	new: {
		[xmlID: string]: readonly string[];
	};
	unchanged: {
		[xmlID: string]: readonly string[];
	};
}

/**
 * Contains the actual differences in value between two cache servers.
 */
export type ServerDifferences = {
	[k in Exclude<keyof SnapshotContentServer, "capabilities" | "deliveryServices">]: DiffVal;
} & {
	capabilities: {
		deleted: Set<string>;
		new: Set<string>;
		unchanged: Set<string>;
	};
	deliveryServices?: ServerDSesDiff;
};

/**
 * A summary of the differences between two Traffic Routers.
 */
interface ServerDiff {
	/** Whether there's any difference between the two. */
	changed: boolean;
	diff: ServerDifferences;
}

/**
 * Summarizes the differences between two sets of CDN Snapshot cache servers.
 */
export interface ServersDiff {
	readonly changed: readonly Readonly<ServerDifferences>[];
	readonly changes: number;
	readonly deleted: readonly Readonly<SnapshotContentServer>[];
	readonly new: readonly Readonly<SnapshotContentServer>[];
	readonly unchanged: readonly Readonly<SnapshotContentServer>[];
}

/**
 * Finds the differences between the (defined) `deliveryServices` properties of
 * two CDN Snapshot `"EDGE"`-type cache servers.
 *
 * @param current The current set of Delivery Services.
 * @param pending The pending set of Delivery Services.
 * @returns A summary of the differences between the two sets of Delivery
 * Services.
 */
function diffDSes(current: Record<string, readonly string[]>, pending: Record<string, readonly string[]>): ServerDSesDiff {
	const diff: ServerDSesDiff = {
		changed: {},
		deleted: {},
		new: {},
		unchanged: {}
	};

	const currentHas = new Set();
	for (const [xmlID, ds] of Object.entries(current)) {
		currentHas.add(xmlID);
		const pDS = pending[xmlID];
		if (!pDS) {
			diff.deleted[xmlID] = ds;
		} else {
			const arrDiff = orderedArrayDiff(ds, pDS);
			if (arrDiff.changed) {
				diff.changed[xmlID] = arrDiff.changes;
			} else {
				diff.unchanged[xmlID] = ds;
			}
		}
	}

	diff.new = Object.fromEntries(Object.entries(pending).filter(([xmlID])=>!currentHas.has(xmlID)));

	return diff;
}

/**
 * Small helper function that just checks if the given object has any keys at
 * all.
 *
 * @param x The object to check.
 * @returns `true` if `x` has any "own" properties, `false` otherwise.
 */
function hasKeys(x: Record<PropertyKey, unknown>): boolean {
	return Object.keys(x).length > 0;
}

/**
 * Finds the differences between two separate definitions of the same cache
 * server.
 *
 * @param current The old definition of the cache server.
 * @param pending The new definition of the cache server.
 * @returns A summary of the differences between `current` and `pending`.
 */
function diffServers(current: SnapshotContentServer, pending: SnapshotContentServer): ServerDiff  {
	let changed = false;
	const diff = {} as ServerDifferences;
	if (current.cacheGroup !== pending.cacheGroup) {
		changed = true;
	}
	diff.cacheGroup = {
		newValue: pending.cacheGroup,
		oldValue: current.cacheGroup
	};

	if (!changed && current.hashCount !== pending.hashCount) {
		changed = true;
	}
	diff.hashCount = {
		newValue: pending.hashCount,
		oldValue: current.hashCount
	};
	if (!changed && current.fqdn !== pending.fqdn) {
		changed = true;
	}
	diff.fqdn = {
		newValue: pending.fqdn,
		oldValue: current.fqdn
	};
	if (!changed && current.httpsPort !== pending.httpsPort) {
		changed = true;
	}
	diff.httpsPort = {
		newValue: pending.httpsPort,
		oldValue: current.httpsPort
	};
	if (!changed && current.ip !== pending.ip) {
		changed = true;
	}
	diff.interfaceName = {
		newValue: pending.interfaceName,
		oldValue: current.interfaceName
	};
	if (!changed && current.interfaceName !== pending.interfaceName) {
		changed = true;
	}
	diff.ip = {
		newValue: pending.ip,
		oldValue: current.ip
	};
	if (!changed && current.ip6 !== pending.ip6) {
		changed = true;
	}
	diff.ip6 = {
		newValue: pending.ip6,
		oldValue: current.ip6
	};
	if (!changed && current.hashId !== pending.hashId) {
		changed = true;
	}
	diff.hashId = {
		newValue: pending.hashId,
		oldValue: current.hashId
	};
	if (!changed && current.port !== pending.port) {
		changed = true;
	}
	diff.port = {
		newValue: pending.port,
		oldValue: current.port
	};
	if (!changed && current.profile !== pending.profile) {
		changed = true;
	}
	diff.profile = {
		newValue: pending.profile,
		oldValue: current.profile
	};
	if (!changed && current.status !== pending.status) {
		changed = true;
	}
	diff.status = {
		newValue: pending.status,
		oldValue: current.status
	};
	if (!changed && current.locationId !== pending.locationId) {
		changed = true;
	}
	diff.locationId = {
		newValue: pending.locationId,
		oldValue: current.locationId
	};
	if (!changed && current.routingDisabled !== pending.routingDisabled) {
		changed = true;
	}
	diff.routingDisabled = {
		newValue: !!pending.routingDisabled,
		oldValue: !!current.routingDisabled
	};

	diff.capabilities = arrayDiff(current.capabilities ?? [], pending.capabilities ?? []);
	if (!changed && (diff.capabilities.deleted.size !== 0 || diff.capabilities.new.size !== 0)) {
		changed = true;
	}

	if (current.type === "EDGE") {
		if (pending.type === "MID") {
			changed = true;
			diff.type = {
				newValue: "MID",
				oldValue: "EDGE"
			};
			diff.deliveryServices = {
				changed: {},
				deleted: current.deliveryServices,
				new: {},
				unchanged: {}
			};
		} else {
			diff.type = {
				newValue: "EDGE",
				oldValue: "EDGE"
			};
			const dsDiff = diffDSes(current.deliveryServices, pending.deliveryServices);
			diff.deliveryServices = dsDiff;
			if (!changed && (hasKeys(dsDiff.changed) || hasKeys(dsDiff.deleted) || hasKeys(dsDiff.new))) {
				changed = true;
			}
		}
	} else if (pending.type === "EDGE") {
		changed = true;
		diff.type = {
			newValue: "EDGE",
			oldValue: "MID"
		};
		diff.deliveryServices = {
			changed: {},
			deleted: {},
			new: pending.deliveryServices,
			unchanged: {}
		};
	} else {
		diff.type = {
			newValue: "MID",
			oldValue: "MID"
		};
	}

	return {changed, diff};
}

/**
 * Finds the differences between a current Snapshot's `contentServers` and a
 * pending Snapshot's `contentServers`.
 *
 * @param current The current Snapshot's server set.
 * @param pending The pending Snapshot's server set.
 * @returns A summary of the differences between the two sets of servers.
 */
export function serverDifferences(
	current: Record<string, SnapshotContentServer>,
	pending: Record<string, SnapshotContentServer>
): ServersDiff {
	const currentHas = new Set<string>();
	const diffs = {
		changed: new Array<ServerDifferences>(),
		changes: 0,
		deleted: new Array<SnapshotContentServer>(),
		new: new Array<SnapshotContentServer>(),
		unchanged: new Array<SnapshotContentServer>()
	};

	for (const [k, v] of Object.entries(current)) {
		currentHas.add(k);
		const p = pending[k];
		if (!p) {
			diffs.deleted.push(v);
			++diffs.changes;
		} else {
			const diff = diffServers(v, p);
			if (diff.changed) {
				++diffs.changes;
				diffs.changed.push(diff.diff);
			} else {
				diffs.unchanged.push(v);
			}
		}
	}

	diffs.new = Object.entries(pending).filter(r=>!currentHas.has(r[0])).map(r=>r[1]);
	diffs.changes += diffs.new.length;
	return diffs;
}
