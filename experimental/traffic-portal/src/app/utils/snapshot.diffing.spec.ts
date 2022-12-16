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

/* eslint-disable @typescript-eslint/naming-convention */
import type { SnapshotContentRouter, SnapshotContentServer } from "trafficops-types";

import { arrayDiff, buildDiff, routerDifferences, serverDifferences, type Diffable } from "./snapshot.diffing";

describe("CDN Snapshot diffing utilities", () => {
	it("can find differences between generic record objects", () => {
		const a: Record<PropertyKey, Diffable> = {
			foo: "bar",
			gee: "whiz",
			test: "quest"
		};

		const b: Record<PropertyKey, Diffable> = {
			fizz: "buzz",
			gee: "willikers",
			test: "quest"
		};

		const diff = buildDiff(a, b);
		expect(diff.num).toBe(3);
		const {foo, fizz, gee, test} = diff.fields;
		if (foo) {
			expect(foo.newValue).toBeUndefined();
			expect(foo.oldValue).toBe(a.foo);
		} else {
			fail("diff shows no 'foo' property");
		}

		if (fizz) {
			expect(fizz.newValue).toBe(b.fizz);
			expect(fizz.oldValue).toBeUndefined();
		} else {
			fail("diff shows no 'fizz' property");
		}

		if (gee) {
			expect(gee.newValue).toBe(b.gee);
			expect(gee.oldValue).toBe(a.gee);
		} else {
			fail("diff shows no 'gee' property");
		}

		if (test) {
			expect(test.newValue).toBe(a.test);
			expect(test.newValue).toBe(b.test);
			expect(test.newValue).toBe(test.oldValue);
		} else {
			fail("diff shows no 'test' property");
		}
	});

	it("diffs arrays without considering order", () => {
		const a = ["foo", "bar", "test", "quest"];
		const b = ["fizz", "test", "quest", "buzz"];

		const diff = arrayDiff(a, b);
		expect(diff.deleted).toHaveSize(2);
		expect(diff.deleted).toContain("foo");
		expect(diff.deleted).toContain("bar");

		expect(diff.new).toHaveSize(2);
		expect(diff.new).toContain("fizz");
		expect(diff.new).toContain("buzz");

		expect(diff.unchanged).toHaveSize(2);
		expect(diff.unchanged).toContain("test");
		expect(diff.unchanged).toContain("quest");
	});

	it("diffs arrays while considering order", () => {
		const a = ["foo", "bar", "test", "quest"];
		const b = ["fizz", "quest", "test", "buzz", "entirely new value"];

		let diff = arrayDiff(a, b, true);
		expect(diff.changed).toBeTrue();

		if (diff.changes.length === 5) {
			expect(diff.changes[0].newValue).toBe("fizz");
			expect(diff.changes[0].oldValue).toBe("foo");
			expect(diff.changes[1].newValue).toBe("quest");
			expect(diff.changes[1].oldValue).toBe("bar");
			expect(diff.changes[2].newValue).toBe("test");
			expect(diff.changes[2].oldValue).toBe("test");
			expect(diff.changes[3].newValue).toBe("buzz");
			expect(diff.changes[3].oldValue).toBe("quest");
			expect(diff.changes[4].newValue).toBe("entirely new value");
			expect(diff.changes[4].oldValue).toBeUndefined();
		} else {
			fail(`array diff length should have been 5, got: ${diff.changes.length}`);
		}

		diff = arrayDiff(a, a, true);
		expect(diff.changed).toBeFalse();
		if (diff.changes.length === a.length) {
			for (let i = 0; i < a.length; ++i) {
				expect(diff.changes[i].newValue).toBe(a[i]);
				expect(diff.changes[i].oldValue).toBe(a[i]);
			}
		} else {
			fail(`diff with no changes should have the same length as its argument (${a.length}), got: ${diff.changes.length}`);
		}
	});

	it("correctly calculates the differences between contentRouter entries", () => {
		const current: Record<string, SnapshotContentRouter> = {
			"fizz.buzz": {
				"api.port": "8",
				fqdn: "test.quest.cdn.test",
				httpsPort: 9,
				ip: "0.0.0.3",
				ip6: "::3",
				location: "test.quest location",
				port: 10,
				profile: "test.quest profile",
				"secure.api.port": "11",
				status: "ONLINE"
			},
			"foo.bar": {
				"api.port": "4",
				fqdn: "foo.bar.cdn.test",
				httpsPort: 5,
				ip: "0.0.0.2",
				ip6: "::2",
				location: "foo.bar location",
				port: 6,
				profile: "foo.bar profile",
				"secure.api.port": "7",
				status: "REPORTED"
			},
			"test.quest": {
				"api.port": "1",
				fqdn: "test.quest.cdn.test",
				httpsPort: 2,
				ip: "0.0.0.1",
				ip6: "::1",
				location: "test.quest location",
				port: 3,
				profile: "test.quest profile",
				"secure.api.port": "4",
				status: "some other third status"
			},
		};
		const pending: Record<string, SnapshotContentRouter> = {
			"fizz.buzz": {
				"api.port": "8",
				fqdn: "test.quest.cdn.test",
				httpsPort: 9,
				ip: "0.0.0.3",
				ip6: "::3",
				location: "test.quest location",
				port: 10,
				profile: "test.quest profile",
				"secure.api.port": "11",
				status: "ONLINE"
			},
			"foo.bar": {
				"api.port": "12",
				fqdn: "foo.bar.changed.cdn.test",
				httpsPort: 13,
				ip: "0.0.0.4",
				ip6: "::4",
				location: "foo.bar changed location",
				port: 14,
				profile: "foo.bar changed profile",
				"secure.api.port": "15",
				status: "ADMIN_UP"
			},
			"gee.whiz": {
				"api.port": "1",
				fqdn: "test.quest.cdn.test",
				httpsPort: 2,
				ip: "0.0.0.1",
				ip6: "::1",
				location: "test.quest location",
				port: 3,
				profile: "test.quest profile",
				"secure.api.port": "4",
				status: "some other third status"
			}
		};
		const diff = routerDifferences(current, pending);
		expect(diff.changes).toBe(diff.changed.length + diff.deleted.length + diff.new.length);
		expect(diff.changes).toBe(3);
		if (diff.changed.length === 1) {
			const changed = diff.changed[0];
			expect(changed).toEqual({
				"api.port": {
					newValue: "12",
					oldValue: "4",
				},
				fqdn: {
					newValue: "foo.bar.changed.cdn.test",
					oldValue: "foo.bar.cdn.test"
				},
				httpsPort: {
					newValue: 13,
					oldValue: 5
				},
				ip: {
					newValue: "0.0.0.4",
					oldValue: "0.0.0.2"
				},
				ip6: {
					newValue: "::4",
					oldValue: "::2"
				},
				location: {
					newValue: "foo.bar changed location",
					oldValue: "foo.bar location"
				},
				port: {
					newValue: 14,
					oldValue: 6
				},
				profile: {
					newValue: "foo.bar changed profile",
					oldValue: "foo.bar profile"
				},
				"secure.api.port": {
					newValue: "15",
					oldValue: "7"
				},
				status: {
					newValue: "ADMIN_UP",
					oldValue: "REPORTED"
				}
			});
		} else {
			fail(`should have shown one changed router, got: ${diff.changed.length}`);
		}

		if (diff.new.length !== 1) {
			return fail(`should have show one new router, got: ${diff.new.length}`);
		}
		expect(diff.new[0]).toEqual(pending["gee.whiz"]);
		if (diff.deleted.length !== 1) {
			return fail(`should have show one deleted router, got: ${diff.deleted.length}`);
		}
		expect(diff.deleted[0]).toEqual(current["test.quest"]);

		if (diff.unchanged.length !== 1) {
			return fail(`should have show one unchanged router, got: ${diff.unchanged.length}`);
		}
		expect(diff.unchanged[0]).toEqual(current["fizz.buzz"]);
		expect(diff.unchanged[0]).toEqual(pending["fizz.buzz"]);
	});

	it("detects changes to routers in any field", () => {
		const routerA = {
			"api.port": "8",
			fqdn: "test.quest.cdn.test",
			httpsPort: 9,
			ip: "0.0.0.3",
			ip6: "::3",
			location: "test.quest location",
			port: 10,
			profile: "test.quest profile",
			"secure.api.port": "11",
			status: "ONLINE"
		};
		const routerB = {
			...routerA
		};
		const current = {a: routerA};
		const pending = {a: routerB};
		let diff = routerDifferences(current, pending);
		expect(diff.changes).toBe(0);
		routerA["api.port"] += "7";
		diff = routerDifferences(current, pending);
		expect(diff.changes).toBe(1);
		routerA["api.port"] = routerB["api.port"];
		routerA.fqdn += ".foo";
		diff = routerDifferences(current, pending);
		expect(diff.changes).toBe(1);
		routerA.fqdn = routerB.fqdn;
		++routerA.httpsPort;
		diff = routerDifferences(current, pending);
		expect(diff.changes).toBe(1);
		routerA.httpsPort = routerB.httpsPort;
		routerA.ip = "1.2.3.4";
		diff = routerDifferences(current, pending);
		expect(diff.changes).toBe(1);
		routerA.ip = routerB.ip;
		routerA.ip6 = "1::";
		diff = routerDifferences(current, pending);
		expect(diff.changes).toBe(1);
		routerA.ip6 = routerB.ip6;
		routerA.location += "foo";
		diff = routerDifferences(current, pending);
		expect(diff.changes).toBe(1);
		routerA.location = routerB.location;
		++routerA.port;
		diff = routerDifferences(current, pending);
		expect(diff.changes).toBe(1);
		routerA.port = routerB.port;
		routerA.profile += "foo";
		diff = routerDifferences(current, pending);
		expect(diff.changes).toBe(1);
		routerA.profile = routerB.profile;
		routerA["secure.api.port"] += "7";
		diff = routerDifferences(current, pending);
		expect(diff.changes).toBe(1);
		routerA["secure.api.port"] = routerB["secure.api.port"];
		routerA.status += " but not really";
		diff = routerDifferences(current, pending);
		expect(diff.changes).toBe(1);
	});

	it("correctly calculates the differences between capability-less MID servers", () => {
		const current: Record<string, SnapshotContentServer> = {
			"fizz.buzz": {
				cacheGroup: "fizz.buzz cache group",
				capabilities: [],
				fqdn: "fizz.buzz.cdn.test",
				hashCount: 1,
				hashId: "hash ID 1",
				httpsPort: 1,
				interfaceName: "eth0",
				ip: "0.0.0.1",
				ip6: "::1",
				locationId: "fizz.buzz location",
				port: 2,
				profile: "fizz.buzz profile",
				routingDisabled: 0,
				status: "ONLINE",
				type: "MID"
			},
			"foo.bar": {
				cacheGroup: "foo.bar cache group",
				capabilities: [],
				fqdn: "foo.bar.cdn.test",
				hashCount: 2,
				hashId: "hash ID 2",
				httpsPort: 3,
				interfaceName: "eth1",
				ip: "0.0.0.2",
				ip6: "::2",
				locationId: "foo.bar location",
				port: 4,
				profile: "foo.bar profile",
				routingDisabled: 0,
				status: "REPORTED",
				type: "MID"
			},
			"test.quest": {
				cacheGroup: "test.quest cache group",
				capabilities: [],
				fqdn: "test.quest.cdn.test",
				hashCount: 4,
				hashId: "hash ID 4",
				httpsPort: 7,
				interfaceName: "eth2",
				ip: "0.0.0.3",
				ip6: ":3",
				locationId: "test.quest location",
				port: 8,
				profile: "test.quest profile",
				routingDisabled: 0,
				status: "test.quest status",
				type: "MID"
			},
		};
		const pending: Record<string, SnapshotContentServer> = {
			"fizz.buzz": {
				cacheGroup: "fizz.buzz cache group",
				capabilities: [],
				fqdn: "fizz.buzz.cdn.test",
				hashCount: 1,
				hashId: "hash ID 1",
				httpsPort: 1,
				interfaceName: "eth0",
				ip: "0.0.0.1",
				ip6: "::1",
				locationId: "fizz.buzz location",
				port: 2,
				profile: "fizz.buzz profile",
				routingDisabled: 0,
				status: "ONLINE",
				type: "MID"
			},
			"foo.bar": {
				cacheGroup: "foo.bar changed cache group",
				capabilities: [],
				deliveryServices: {},
				fqdn: "foo.bar.changed.cdn.test",
				hashCount: 3,
				hashId: "hash ID 3",
				httpsPort: 5,
				interfaceName: "eth2",
				ip: "0.0.0.3",
				ip6: "::3",
				locationId: "foo.bar changed location",
				port: 6,
				profile: "foo.bar changed profile",
				routingDisabled: 1,
				status: "NOT REPORTED",
				type: "EDGE"
			},
			"gee.whiz": {
				cacheGroup: "gee.whiz cache group",
				capabilities: [],
				fqdn: "gee.whiz.cdn.test",
				hashCount: 5,
				hashId: "hash ID 5",
				httpsPort: 9,
				interfaceName: "eth3",
				ip: "0.0.0.4",
				ip6: "::4",
				locationId: "gee.whiz location",
				port: 10,
				profile: "gee.whiz profile",
				routingDisabled: 0,
				status: "gee.whiz status",
				type: "MID"
			}
		};
		const diff = serverDifferences(current, pending);
		expect(diff.changes).toBe(diff.changed.length + diff.deleted.length + diff.new.length);
		expect(diff.changes).toBe(3);
		if (diff.changed.length === 1) {
			const changed = diff.changed[0];
			expect(changed).toEqual({
				cacheGroup: {
					newValue: "foo.bar changed cache group",
					oldValue: "foo.bar cache group"
				},
				capabilities: {
					deleted: new Set(),
					new: new Set(),
					unchanged: new Set()
				},
				deliveryServices: {
					changed: {},
					deleted: {},
					new: {},
					unchanged: {}
				},
				fqdn: {
					newValue: "foo.bar.changed.cdn.test",
					oldValue: "foo.bar.cdn.test"
				},
				hashCount: {
					newValue: 3,
					oldValue: 2
				},
				hashId: {
					newValue: "hash ID 3",
					oldValue: "hash ID 2"
				},
				httpsPort: {
					newValue: 5,
					oldValue: 3
				},
				interfaceName: {
					newValue: "eth2",
					oldValue: "eth1"
				},
				ip: {
					newValue: "0.0.0.3",
					oldValue: "0.0.0.2"
				},
				ip6: {
					newValue: "::3",
					oldValue: "::2"
				},
				locationId: {
					newValue: "foo.bar changed location",
					oldValue: "foo.bar location"
				},
				port: {
					newValue: 6,
					oldValue: 4
				},
				profile: {
					newValue: "foo.bar changed profile",
					oldValue: "foo.bar profile"
				},
				// These are coerced to boolean values
				routingDisabled: {
					newValue: true,
					oldValue: false
				},
				status: {
					newValue: "NOT REPORTED",
					oldValue: "REPORTED"
				},
				type: {
					newValue: "EDGE",
					oldValue: "MID"
				},
			});
		} else {
			fail(`should have shown one changed router, got: ${diff.changed.length}`);
		}

		if (diff.new.length !== 1) {
			return fail(`should have show one new router, got: ${diff.new.length}`);
		}
		expect(diff.new[0]).toEqual(pending["gee.whiz"]);
		if (diff.deleted.length !== 1) {
			return fail(`should have show one deleted router, got: ${diff.deleted.length}`);
		}
		expect(diff.deleted[0]).toEqual(current["test.quest"]);

		if (diff.unchanged.length !== 1) {
			return fail(`should have show one unchanged router, got: ${diff.unchanged.length}`);
		}
		expect(diff.unchanged[0]).toEqual(current["fizz.buzz"]);
		expect(diff.unchanged[0]).toEqual(pending["fizz.buzz"]);
	});

	it("correctly diffs servers with only capability differences", () => {
		const current: Record<string, SnapshotContentServer> = {
			"test.quest": {
				cacheGroup: "test.quest cache group",
				capabilities: [
					"a",
					"b",
					"c"
				],
				fqdn: "test.quest.cdn.test",
				hashCount: 4,
				hashId: "hash ID 4",
				httpsPort: 7,
				interfaceName: "eth2",
				ip: "0.0.0.3",
				ip6: ":3",
				locationId: "test.quest location",
				port: 8,
				profile: "test.quest profile",
				routingDisabled: 0,
				status: "test.quest status",
				type: "MID"
			},
		};
		const pending: Record<string, SnapshotContentServer> = {
			"test.quest": {
				cacheGroup: "test.quest cache group",
				capabilities: [
					"a",
					"c",
					"d"
				],
				fqdn: "test.quest.cdn.test",
				hashCount: 4,
				hashId: "hash ID 4",
				httpsPort: 7,
				interfaceName: "eth2",
				ip: "0.0.0.3",
				ip6: ":3",
				locationId: "test.quest location",
				port: 8,
				profile: "test.quest profile",
				routingDisabled: 0,
				status: "test.quest status",
				type: "MID"
			},
		};
		const diff = serverDifferences(current, pending);
		expect(diff.changed).toHaveSize(1);
		expect(diff.deleted).toHaveSize(0);
		expect(diff.new).toHaveSize(0);
		expect(diff.unchanged).toHaveSize(0);
		if (diff.changed.length !== 1) {
			return fail(`should have been exactly one changed server, got: ${diff.changed.length}`);
		}
		const serverDiff = diff.changed[0];
		expect(serverDiff.cacheGroup.newValue).toBe(serverDiff.cacheGroup.oldValue);
		expect(serverDiff.fqdn.newValue).toBe(serverDiff.fqdn.oldValue);
		expect(serverDiff.hashCount.newValue).toBe(serverDiff.hashCount.oldValue);
		expect(serverDiff.hashId.newValue).toBe(serverDiff.hashId.oldValue);
		expect(serverDiff.httpsPort.newValue).toBe(serverDiff.httpsPort.oldValue);
		expect(serverDiff.interfaceName.newValue).toBe(serverDiff.interfaceName.oldValue);
		expect(serverDiff.ip.newValue).toBe(serverDiff.ip.oldValue);
		expect(serverDiff.ip6.newValue).toBe(serverDiff.ip6.oldValue);
		expect(serverDiff.locationId.newValue).toBe(serverDiff.locationId.oldValue);
		expect(serverDiff.port.newValue).toBe(serverDiff.port.oldValue);
		expect(serverDiff.profile.newValue).toBe(serverDiff.profile.oldValue);
		expect(serverDiff.routingDisabled.newValue).toBe(serverDiff.routingDisabled.oldValue);
		expect(serverDiff.status.newValue).toBe(serverDiff.status.oldValue);
		expect(serverDiff.type.newValue).toBe(serverDiff.type.oldValue);

		expect(serverDiff.deliveryServices).toBeUndefined();
		const caps = serverDiff.capabilities;
		expect(caps.deleted).toHaveSize(1);
		expect(caps.deleted).toContain("b");
		expect(caps.new).toHaveSize(1);
		expect(caps.new).toContain("d");
		expect(caps.unchanged).toHaveSize(2);
		expect(caps.unchanged).toContain("a");
		expect(caps.unchanged).toContain("c");
	});

	it("correctly finds the ordered differences between a server's `deliveryServices`", () => {
		const current: Record<string, SnapshotContentServer> = {
			"fizz.buzz": {
				cacheGroup: "test.quest cache group",
				capabilities: [],
				deliveryServices: {
					a: ["1", "2", "3"],
					b: ["4", "5", "6"],
					c: ["7", "8", "9"]
				},
				fqdn: "fizz.buzz",
				hashCount: 4,
				hashId: "hash ID 4",
				httpsPort: 7,
				interfaceName: "eth2",
				ip: "0.0.0.3",
				ip6: ":3",
				locationId: "test.quest location",
				port: 8,
				profile: "test.quest profile",
				routingDisabled: 0,
				status: "test.quest status",
				type: "EDGE"
			},
			"foo.bar": {
				cacheGroup: "test.quest cache group",
				capabilities: [],
				deliveryServices: {
					a: ["1", "2", "3"]
				},
				fqdn: "foo.bar",
				hashCount: 4,
				hashId: "hash ID 4",
				httpsPort: 7,
				interfaceName: "eth2",
				ip: "0.0.0.3",
				ip6: ":3",
				locationId: "test.quest location",
				port: 8,
				profile: "test.quest profile",
				routingDisabled: 0,
				status: "test.quest status",
				type: "EDGE"
			},
			"gee.whiz": {
				cacheGroup: "test.quest cache group",
				capabilities: [],
				fqdn: "gee.whiz",
				hashCount: 4,
				hashId: "hash ID 4",
				httpsPort: 7,
				interfaceName: "eth2",
				ip: "0.0.0.3",
				ip6: ":3",
				locationId: "test.quest location",
				port: 8,
				profile: "test.quest profile",
				routingDisabled: 0,
				status: "test.quest status",
				type: "MID"
			},
			"test.quest": {
				cacheGroup: "test.quest cache group",
				capabilities: [],
				fqdn: "test.quest",
				hashCount: 4,
				hashId: "hash ID 4",
				httpsPort: 7,
				interfaceName: "eth2",
				ip: "0.0.0.3",
				ip6: ":3",
				locationId: "test.quest location",
				port: 8,
				profile: "test.quest profile",
				routingDisabled: 0,
				status: "test.quest status",
				type: "MID"
			},
		};
		const pending: Record<string, SnapshotContentServer> = {
			"fizz.buzz": {
				cacheGroup: "test.quest cache group",
				capabilities: [],
				deliveryServices: {
					a: ["1", "2", "3"],
					b: ["5", "6", "7"],
					d: ["8", "9", "10"]
				},
				fqdn: "fizz.buzz",
				hashCount: 4,
				hashId: "hash ID 4",
				httpsPort: 7,
				interfaceName: "eth2",
				ip: "0.0.0.3",
				ip6: ":3",
				locationId: "test.quest location",
				port: 8,
				profile: "test.quest profile",
				routingDisabled: 0,
				status: "test.quest status",
				type: "EDGE"
			},
			"foo.bar": {
				cacheGroup: "test.quest cache group",
				capabilities: [],
				fqdn: "foo.bar",
				hashCount: 4,
				hashId: "hash ID 4",
				httpsPort: 7,
				interfaceName: "eth2",
				ip: "0.0.0.3",
				ip6: ":3",
				locationId: "test.quest location",
				port: 8,
				profile: "test.quest profile",
				routingDisabled: 0,
				status: "test.quest status",
				type: "MID"
			},
			"gee.whiz": {
				cacheGroup: "test.quest cache group",
				capabilities: [],
				deliveryServices: {
					a: ["1", "2", "3"]
				},
				fqdn: "gee.whiz",
				hashCount: 4,
				hashId: "hash ID 4",
				httpsPort: 7,
				interfaceName: "eth2",
				ip: "0.0.0.3",
				ip6: ":3",
				locationId: "test.quest location",
				port: 8,
				profile: "test.quest profile",
				routingDisabled: 0,
				status: "test.quest status",
				type: "EDGE"
			},
			"test.quest": {
				cacheGroup: "test.quest cache group",
				capabilities: [],
				fqdn: "test.quest",
				hashCount: 4,
				hashId: "hash ID 4",
				httpsPort: 7,
				interfaceName: "eth2",
				ip: "0.0.0.3",
				ip6: ":3",
				locationId: "test.quest location",
				port: 8,
				profile: "test.quest profile",
				routingDisabled: 0,
				status: "test.quest status",
				type: "MID"
			},
		};
		const diff = serverDifferences(current, pending);
		expect(diff.changes).toBe(3);
		if (diff.changed.length !== 3) {
			return fail(`on basis of type and DS assignment, should have found 3 changed servers; got: ${diff.changed.length}`);
		}
		let diffedSrv = diff.changed[0];
		expect(diffedSrv.cacheGroup.newValue).toBe(diffedSrv.cacheGroup.oldValue);
		expect(diffedSrv.fqdn.newValue).toBe(diffedSrv.fqdn.oldValue);
		expect(diffedSrv.fqdn.newValue).toBe("fizz.buzz");
		expect(diffedSrv.hashCount.newValue).toBe(diffedSrv.hashCount.oldValue);
		expect(diffedSrv.hashId.newValue).toBe(diffedSrv.hashId.oldValue);
		expect(diffedSrv.httpsPort.newValue).toBe(diffedSrv.httpsPort.oldValue);
		expect(diffedSrv.interfaceName.newValue).toBe(diffedSrv.interfaceName.oldValue);
		expect(diffedSrv.ip.newValue).toBe(diffedSrv.ip.oldValue);
		expect(diffedSrv.ip6.newValue).toBe(diffedSrv.ip6.oldValue);
		expect(diffedSrv.locationId.newValue).toBe(diffedSrv.locationId.oldValue);
		expect(diffedSrv.port.newValue).toBe(diffedSrv.port.oldValue);
		expect(diffedSrv.profile.newValue).toBe(diffedSrv.profile.oldValue);
		expect(diffedSrv.routingDisabled.newValue).toBe(diffedSrv.routingDisabled.oldValue);
		expect(diffedSrv.status.newValue).toBe(diffedSrv.status.oldValue);
		expect(diffedSrv.type.newValue).toBe(diffedSrv.type.oldValue);
		expect(diffedSrv.capabilities.deleted).toHaveSize(0);
		expect(diffedSrv.capabilities.new).toHaveSize(0);
		expect(diffedSrv.capabilities.unchanged).toHaveSize(0);
		let dsDiff = diffedSrv.deliveryServices;
		if (!dsDiff) {
			return fail("dsDiff for the first contentServer should've been defined");
		}
		expect(dsDiff.unchanged).toHaveSize(1);
		let d = dsDiff.unchanged.a;
		if (!d) {
			return fail("the first contentServer should've had 'a' marked as unchanged.");
		}
		if (d.length !== 3) {
			return fail(`the first contentServer's unchanged DS length should've been 3, got: ${d.length}`);
		}
		expect(d[0]).toBe("1");
		expect(d[1]).toBe("2");
		expect(d[2]).toBe("3");

		expect(dsDiff.deleted).toHaveSize(1);
		d = dsDiff.deleted.c;
		if (!d) {
			return fail("the first contentServer should've had 'c' marked as deleted");
		}
		if (d.length !== 3) {
			return fail(`the first contentServer's deleted DS length should been 3, got: ${d.length}`);
		}
		expect(d[0]).toBe("7");
		expect(d[1]).toBe("8");
		expect(d[2]).toBe("9");

		expect(dsDiff.new).toHaveSize(1);
		d = dsDiff.new.d;
		if (!d) {
			return fail("the first contentServer should've had 'c' marked as deleted");
		}
		if (d.length !== 3) {
			return fail(`the first contentServer's deleted DS length should been 3, got: ${d.length}`);
		}
		expect(d[0]).toBe("8");
		expect(d[1]).toBe("9");
		expect(d[2]).toBe("10");

		expect(dsDiff.changed).toHaveSize(1);
		const changedD = dsDiff.changed.b;
		if (!changedD) {
			return fail("the first contentServer should've had 'c' marked as deleted");
		}
		if (changedD.length !== 3) {
			return fail(`the first contentServer's deleted DS length should been 3, got: ${d.length}`);
		}
		expect(changedD[0]).toEqual({newValue: "5", oldValue: "4"});
		expect(changedD[1]).toEqual({newValue: "6", oldValue: "5"});
		expect(changedD[2]).toEqual({newValue: "7", oldValue: "6"});

		diffedSrv = diff.changed[1];
		expect(diffedSrv.cacheGroup.newValue).toBe(diffedSrv.cacheGroup.oldValue);
		expect(diffedSrv.fqdn.newValue).toBe(diffedSrv.fqdn.oldValue);
		expect(diffedSrv.fqdn.newValue).toBe("foo.bar");
		expect(diffedSrv.hashCount.newValue).toBe(diffedSrv.hashCount.oldValue);
		expect(diffedSrv.hashId.newValue).toBe(diffedSrv.hashId.oldValue);
		expect(diffedSrv.httpsPort.newValue).toBe(diffedSrv.httpsPort.oldValue);
		expect(diffedSrv.interfaceName.newValue).toBe(diffedSrv.interfaceName.oldValue);
		expect(diffedSrv.ip.newValue).toBe(diffedSrv.ip.oldValue);
		expect(diffedSrv.ip6.newValue).toBe(diffedSrv.ip6.oldValue);
		expect(diffedSrv.locationId.newValue).toBe(diffedSrv.locationId.oldValue);
		expect(diffedSrv.port.newValue).toBe(diffedSrv.port.oldValue);
		expect(diffedSrv.profile.newValue).toBe(diffedSrv.profile.oldValue);
		expect(diffedSrv.routingDisabled.newValue).toBe(diffedSrv.routingDisabled.oldValue);
		expect(diffedSrv.status.newValue).toBe(diffedSrv.status.oldValue);
		expect(diffedSrv.type.newValue).not.toBe(diffedSrv.type.oldValue);
		expect(diffedSrv.type.oldValue).toBe("EDGE");
		expect(diffedSrv.type.newValue).toBe("MID");
		expect(diffedSrv.capabilities.deleted).toHaveSize(0);
		expect(diffedSrv.capabilities.new).toHaveSize(0);
		expect(diffedSrv.capabilities.unchanged).toHaveSize(0);
		expect(diffedSrv.deliveryServices).toBeDefined();
		dsDiff = diffedSrv.deliveryServices;
		if (!dsDiff) {
			return fail("dsDiff for the second contentServer should've been defined");
		}
		expect(dsDiff.changed).toHaveSize(0);
		expect(dsDiff.new).toHaveSize(0);
		expect(dsDiff.unchanged).toHaveSize(0);
		expect(dsDiff.deleted).toHaveSize(1);
		d = dsDiff.deleted.a;
		expect(d).toEqual(["1", "2", "3"]);

		diffedSrv = diff.changed[2];
		expect(diffedSrv.cacheGroup.newValue).toBe(diffedSrv.cacheGroup.oldValue);
		expect(diffedSrv.fqdn.newValue).toBe(diffedSrv.fqdn.oldValue);
		expect(diffedSrv.fqdn.newValue).toBe("gee.whiz");
		expect(diffedSrv.hashCount.newValue).toBe(diffedSrv.hashCount.oldValue);
		expect(diffedSrv.hashId.newValue).toBe(diffedSrv.hashId.oldValue);
		expect(diffedSrv.httpsPort.newValue).toBe(diffedSrv.httpsPort.oldValue);
		expect(diffedSrv.interfaceName.newValue).toBe(diffedSrv.interfaceName.oldValue);
		expect(diffedSrv.ip.newValue).toBe(diffedSrv.ip.oldValue);
		expect(diffedSrv.ip6.newValue).toBe(diffedSrv.ip6.oldValue);
		expect(diffedSrv.locationId.newValue).toBe(diffedSrv.locationId.oldValue);
		expect(diffedSrv.port.newValue).toBe(diffedSrv.port.oldValue);
		expect(diffedSrv.profile.newValue).toBe(diffedSrv.profile.oldValue);
		expect(diffedSrv.routingDisabled.newValue).toBe(diffedSrv.routingDisabled.oldValue);
		expect(diffedSrv.status.newValue).toBe(diffedSrv.status.oldValue);
		expect(diffedSrv.type.newValue).not.toBe(diffedSrv.type.oldValue);
		expect(diffedSrv.type.oldValue).toBe("MID");
		expect(diffedSrv.type.newValue).toBe("EDGE");
		expect(diffedSrv.capabilities.deleted).toHaveSize(0);
		expect(diffedSrv.capabilities.new).toHaveSize(0);
		expect(diffedSrv.capabilities.unchanged).toHaveSize(0);
		expect(diffedSrv.deliveryServices).toBeDefined();
		dsDiff = diffedSrv.deliveryServices;
		if (!dsDiff) {
			return fail("dsDiff for the third contentServer should've been defined");
		}
		expect(dsDiff.changed).toHaveSize(0);
		expect(dsDiff.deleted).toHaveSize(0);
		expect(dsDiff.unchanged).toHaveSize(0);
		expect(dsDiff.new).toHaveSize(1);
		d = dsDiff.new.a;
		expect(d).toEqual(["1", "2", "3"]);

		expect(diff.unchanged).toHaveSize(1);
		expect(diff.deleted).toHaveSize(0);
		expect(diff.new).toHaveSize(0);

		const unchangedSrv = diff.unchanged[0];
		const testQuest = current["test.quest"];
		expect(unchangedSrv).toEqual(testQuest);
		expect(pending["test.quest"]).toEqual(testQuest);
		expect(unchangedSrv.fqdn).toBe("test.quest");
		expect(unchangedSrv.type).toBe("MID");
		expect(unchangedSrv.capabilities).toHaveSize(0);
		expect((unchangedSrv as Record<PropertyKey, unknown>).deliveryServices).toBeUndefined();
	});

	it("detects changes to servers in any field", () => {
		let serverA: SnapshotContentServer = {
			cacheGroup: "test cg",
			capabilities: [],
			fqdn: "test.quest.cdn.test",
			hashCount: 1,
			hashId: "test",
			httpsPort: 9,
			interfaceName: "eth0",
			ip: "0.0.0.3",
			ip6: "::3",
			locationId: "test.quest location",
			port: 10,
			profile: "test.quest profile",
			routingDisabled: 0,
			status: "ONLINE",
			type: "MID"
		};
		const serverB = {
			...serverA
		};
		const current: Record<PropertyKey, SnapshotContentServer> = {a: serverA};
		const pending = {a: serverB};
		let diff = serverDifferences(current, pending);
		expect(diff.changes).toBe(0);
		serverA.cacheGroup += "7";
		diff = serverDifferences(current, pending);
		expect(diff.changes).toBe(1);
		serverA.cacheGroup = serverB.cacheGroup;
		serverA.capabilities = ["test"];
		diff = serverDifferences(current, pending);
		expect(diff.changes).toBe(1);
		serverA.capabilities = [];
		serverA.hashCount += 5;
		diff = serverDifferences(current, pending);
		expect(diff.changes).toBe(1);
		serverA.hashCount = serverB.hashCount;
		serverA.hashId += ".foo";
		diff = serverDifferences(current, pending);
		expect(diff.changes).toBe(1);
		serverA.hashId = serverB.hashId;
		serverA.fqdn += ".foo";
		diff = serverDifferences(current, pending);
		expect(diff.changes).toBe(1);
		serverA.fqdn = serverB.fqdn;
		serverA.httpsPort = 9001;
		diff = serverDifferences(current, pending);
		expect(diff.changes).toBe(1);
		serverA.httpsPort = serverB.httpsPort;
		serverA.interfaceName = " but now it's different";
		diff = serverDifferences(current, pending);
		expect(diff.changes).toBe(1);
		serverA.interfaceName = serverB.interfaceName;
		serverA.ip = "1.2.3.4";
		diff = serverDifferences(current, pending);
		expect(diff.changes).toBe(1);
		serverA.ip = serverB.ip;
		serverA.ip6 = "1::";
		diff = serverDifferences(current, pending);
		expect(diff.changes).toBe(1);
		serverA.ip6 = serverB.ip6;
		serverA.locationId += "foo";
		diff = serverDifferences(current, pending);
		expect(diff.changes).toBe(1);
		serverA.locationId = serverB.locationId;
		++serverA.port;
		diff = serverDifferences(current, pending);
		expect(diff.changes).toBe(1);
		serverA.port = serverB.port;
		serverA.profile += "foo";
		diff = serverDifferences(current, pending);
		expect(diff.changes).toBe(1);
		serverA.profile = serverB.profile;
		serverA.routingDisabled = serverA.routingDisabled ? 0 : 1;
		diff = serverDifferences(current, pending);
		expect(diff.changes).toBe(1);
		serverA.routingDisabled = serverB.routingDisabled;
		serverA.status += " but not really";
		diff = serverDifferences(current, pending);
		expect(diff.changes).toBe(1);
		serverA.status = serverB.status;
		serverA = {...serverA, deliveryServices: {}, type: "EDGE"};
		current.a = serverA;
		diff = serverDifferences(current, pending);
		expect(diff.changes).toBe(1);
	});
});
