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

import { EventEmitter } from "@angular/core";
import { of, type Subscription } from "rxjs";

import { Differences } from "./differences.class";

/**
 * A class that's used simply to test the abstract class it extends.
 */
class TestingDifferencesSubclass extends Differences<null, null> {
	public changesPending = new EventEmitter<number>();
	public snapshots = of({current: {}, pending: {}});

	public subscription: Subscription;

	constructor() {
		super();
		this.subscription = this.snapshots.subscribe();
	}

	/**
	 * Angular lifecycle hook.
	 */
	public ngOnInit(): void {
		// do nothing.
	}
}

describe("'Differences' abstract class", () => {
	let instance: TestingDifferencesSubclass;

	beforeEach(() => {
		instance = new TestingDifferencesSubclass();
	});

	it("gets the right string for changes pending", () => {
		expect(instance.pendingChangesStr(0)).toBe("0 changes pending");
		expect(instance.pendingChangesStr(1)).toBe("1 change pending");
		expect(instance.pendingChangesStr(2)).toBe("2 changes pending");
		expect(instance.pendingChangesStr(9001)).toBe("9001 changes pending");

		const arr = [5];
		expect(instance.pendingChangesStr(arr)).toBe("1 change pending");
		arr.push(5);
		expect(instance.pendingChangesStr(arr)).toBe("2 changes pending");
		arr.splice(0, 2);
		expect(instance.pendingChangesStr(arr)).toBe("0 changes pending");

		const obj: Record<PropertyKey, unknown> = {};
		expect(instance.pendingChangesStr(obj)).toBe("0 changes pending");
		obj.test = 5;
		expect(instance.pendingChangesStr(obj)).toBe("1 change pending");
		obj.quest = 5;
		expect(instance.pendingChangesStr(obj)).toBe("2 changes pending");
	});

	it("unsubscribes on component destruction", () => {
		instance.ngOnDestroy();
		expect(instance.subscription.closed).toBeTrue();
	});
});
