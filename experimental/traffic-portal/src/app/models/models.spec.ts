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
import {
	bypassable,
	defaultDeliveryService,
	QStringHandling,
	qStringHandlingToString,
	RangeRequestHandling,
	rangeRequestHandlingToString
} from "./deliveryservice";
import { checkMap, type Servercheck } from "./server";

describe("Delivery Service utilities", () => {
	it("converts Query String Handlings to human-readable strings", () => {
		let output = qStringHandlingToString(QStringHandling.USE);
		expect(output).toBe(
			"Use the query parameter string when deciding if a URL is cached, and pass it in upstream requests to the Mid-tier/origin"
		);
		output = qStringHandlingToString(QStringHandling.IGNORE);
		expect(output).toBe(
			"Do not use the query parameter string when deciding if a URL is cached, but do pass it in upstream requests to the " +
			"Mid-tier/origin"
		);
		output = qStringHandlingToString(QStringHandling.DROP);
		expect(output).toBe(
			"Immediately strip URLs of their query parameter strings before checking cached objects or making upstream requests"
		);
	});

	it("converts Range Request Handlings to human-readable strings", () => {
		let output = rangeRequestHandlingToString(RangeRequestHandling.NONE);
		expect(output).toBe("Do not cache Range requests");
		output = rangeRequestHandlingToString(RangeRequestHandling.BACKGROUND_FETCH);
		expect(output).toBe("Use the background_fetch plugin to serve Range requests while quietly caching the entire object");
		output = rangeRequestHandlingToString(RangeRequestHandling.CACHE_RANGE_REQUESTS);
		expect(output).toBe("Use the cache_range_requests plugin to directly cache object ranges");
	});

	it("can tell by a Delivery Service's Type whether or not it is eligible for bypassing", () => {
		const ds = {...defaultDeliveryService};
		// The defaultDeliveryService has undefined Type.
		expect(bypassable(ds)).toBeFalse();
		ds.type = "HTTP";
		expect(bypassable(ds)).toBeTrue();
		ds.type = "HTTP_LIVE";
		expect(bypassable(ds)).toBeTrue();
		ds.type = "HTTP_LIVE_NATNL";
		expect(bypassable(ds)).toBeTrue();
		ds.type = "DNS";
		expect(bypassable(ds)).toBeTrue();
		ds.type = "DNS_LIVE";
		expect(bypassable(ds)).toBeTrue();
		ds.type = "DNS_LIVE_NATNL";
		expect(bypassable(ds)).toBeTrue();
		ds.type = "ANY_MAP";
		expect(bypassable(ds)).toBeFalse();
		ds.type = "STEERING";
		expect(bypassable(ds)).toBeFalse();
		ds.type = "CLIENT_STEERING";
		expect(bypassable(ds)).toBeFalse();
	});
});

describe("Server utilities", () => {
	it("converts serverchecks into a Map", () => {
		const srv: Servercheck = {
			adminState: "ONLINE",
			cacheGroup: "",
			checks: undefined,
			hostName: "",
			id: 1,
			profile: "",
			revalPending: false,
			type: "",
			updPending: false
		};
		let map = checkMap(srv);
		expect(map.size).toBe(0);
		srv.checks = {};
		map = checkMap(srv);
		expect(map.size).toBe(0);
		srv.checks["10G"] = 1;
		srv.checks.ILO = 0;
		srv.checks.PING = 127;
		map = checkMap(srv);
		expect(map.size).toBe(3);
		expect(map.get("10G")).toBeTrue();
		expect(map.get("ILO")).toBeFalse();
		expect(map.get("PING")).toBe(127);
	});
});
