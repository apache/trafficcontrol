/*
*
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

import { dateReviver, parseHTTPDate } from "./date";

describe("date utilities", () => {
	it("revives legacy, custom TO timestamps into dates", ()=>{
		const data = '{"notADate": "testquest", "myDate": "2022-01-02 03:04:05+00"}';
		const parsed = JSON.parse(data, dateReviver);
		const keys = Object.keys(parsed);
		expect(keys.length).toBe(2);
		expect(keys).toContain("notADate");
		expect(keys).toContain("myDate");
		expect(typeof(parsed.notADate)).toBe("string");

		const parsedDate = parsed.myDate;
		if (!(parsedDate instanceof Date)) {
			return fail(`expected "${parsedDate}" to be a Date`);
		}
		expect(parsedDate.getUTCFullYear()).toBe(2022);
		expect(parsedDate.getUTCMonth()).toBe(0);
		expect(parsedDate.getUTCDate()).toBe(2);
		expect(parsedDate.getUTCHours()).toBe(3);
		expect(parsedDate.getUTCMinutes()).toBe(4);
		expect(parsedDate.getUTCSeconds()).toBe(5);
		expect(parsedDate.getUTCMilliseconds()).toBe(0);
	});

	it("revives RFC3339 timestamps into dates", () => {
		const data = '{"notADate": "testquest", "myDate": "2022-01-02T03:04:05Z"}';
		const parsed = JSON.parse(data, dateReviver);
		const keys = Object.keys(parsed);
		expect(keys.length).toBe(2);
		expect(keys).toContain("notADate");
		expect(keys).toContain("myDate");
		expect(typeof(parsed.notADate)).toBe("string");

		const parsedDate = parsed.myDate;
		if (!(parsedDate instanceof Date)) {
			return fail(`expected "${parsedDate}" to be a Date`);
		}
		expect(parsedDate.getUTCFullYear()).toBe(2022);
		expect(parsedDate.getUTCMonth()).toBe(0);
		expect(parsedDate.getUTCDate()).toBe(2);
		expect(parsedDate.getUTCHours()).toBe(3);
		expect(parsedDate.getUTCMinutes()).toBe(4);
		expect(parsedDate.getUTCSeconds()).toBe(5);
		expect(parsedDate.getUTCMilliseconds()).toBe(0);
	});

	it("revives RFC3339 timestamps with sub-second precision into dates", () => {
		const data = '{"notADate": "testquest", "myDate": "2022-01-02T03:04:05.6789Z"}';
		const parsed = JSON.parse(data, dateReviver);
		const keys = Object.keys(parsed);
		expect(keys.length).toBe(2);
		expect(keys).toContain("notADate");
		expect(keys).toContain("myDate");
		expect(typeof(parsed.notADate)).toBe("string");

		const parsedDate = parsed.myDate;
		if (!(parsedDate instanceof Date)) {
			return fail(`expected "${parsedDate}" to be a Date`);
		}
		expect(parsedDate.getUTCFullYear()).toBe(2022);
		expect(parsedDate.getUTCMonth()).toBe(0);
		expect(parsedDate.getUTCDate()).toBe(2);
		expect(parsedDate.getUTCHours()).toBe(3);
		expect(parsedDate.getUTCMinutes()).toBe(4);
		expect(parsedDate.getUTCSeconds()).toBe(5);
		expect(parsedDate.getUTCMilliseconds()).toBe(678);
	});

	it("leaves unparsable dates alone", () => {
		const data = {
			lastAuthenticated: "not a valid date",
			lastUpdated: "9999-99-99T99:99:99.99Z"
		};
		const parsed = JSON.parse(JSON.stringify(data), dateReviver);
		expect(parsed).toEqual(data);
	});

	it("parses HTTP header dates", () => {
		let date = "Sun, 02 Jan 2022 03:04:05 GMT";
		let parsed = parseHTTPDate(date);
		expect(parsed.getUTCFullYear()).toBe(2022);
		expect(parsed.getUTCMonth()).toBe(0);
		expect(parsed.getUTCDate()).toBe(2);
		expect(parsed.getUTCHours()).toBe(3);
		expect(parsed.getUTCMinutes()).toBe(4);
		expect(parsed.getUTCSeconds()).toBe(5);
		expect(parsed.getUTCMilliseconds()).toBe(0);

		date = "Mon, 07 Feb 2022 03:04:05 GMT";
		parsed = parseHTTPDate(date);
		expect(parsed.getUTCFullYear()).toBe(2022);
		expect(parsed.getUTCMonth()).toBe(1);
		expect(parsed.getUTCDate()).toBe(7);
		expect(parsed.getUTCHours()).toBe(3);
		expect(parsed.getUTCMinutes()).toBe(4);
		expect(parsed.getUTCSeconds()).toBe(5);
		expect(parsed.getUTCMilliseconds()).toBe(0);

		date = "Tue, 01 Mar 2022 03:04:05 GMT";
		parsed = parseHTTPDate(date);
		expect(parsed.getUTCFullYear()).toBe(2022);
		expect(parsed.getUTCMonth()).toBe(2);
		expect(parsed.getUTCDate()).toBe(1);
		expect(parsed.getUTCHours()).toBe(3);
		expect(parsed.getUTCMinutes()).toBe(4);
		expect(parsed.getUTCSeconds()).toBe(5);
		expect(parsed.getUTCMilliseconds()).toBe(0);

		date = "Wed, 06 Apr 2022 03:04:05 GMT";
		parsed = parseHTTPDate(date);
		expect(parsed.getUTCFullYear()).toBe(2022);
		expect(parsed.getUTCMonth()).toBe(3);
		expect(parsed.getUTCDate()).toBe(6);
		expect(parsed.getUTCHours()).toBe(3);
		expect(parsed.getUTCMinutes()).toBe(4);
		expect(parsed.getUTCSeconds()).toBe(5);
		expect(parsed.getUTCMilliseconds()).toBe(0);

		date = "Thu, 05 May 2022 03:04:05 GMT";
		parsed = parseHTTPDate(date);
		expect(parsed.getUTCFullYear()).toBe(2022);
		expect(parsed.getUTCMonth()).toBe(4);
		expect(parsed.getUTCDate()).toBe(5);
		expect(parsed.getUTCHours()).toBe(3);
		expect(parsed.getUTCMinutes()).toBe(4);
		expect(parsed.getUTCSeconds()).toBe(5);
		expect(parsed.getUTCMilliseconds()).toBe(0);

		date = "Fri, 03 Jun 2022 03:04:05 GMT";
		parsed = parseHTTPDate(date);
		expect(parsed.getUTCFullYear()).toBe(2022);
		expect(parsed.getUTCMonth()).toBe(5);
		expect(parsed.getUTCDate()).toBe(3);
		expect(parsed.getUTCHours()).toBe(3);
		expect(parsed.getUTCMinutes()).toBe(4);
		expect(parsed.getUTCSeconds()).toBe(5);
		expect(parsed.getUTCMilliseconds()).toBe(0);

		date = "Sat, 02 Jul 2022 03:04:05 GMT";
		parsed = parseHTTPDate(date);
		expect(parsed.getUTCFullYear()).toBe(2022);
		expect(parsed.getUTCMonth()).toBe(6);
		expect(parsed.getUTCDate()).toBe(2);
		expect(parsed.getUTCHours()).toBe(3);
		expect(parsed.getUTCMinutes()).toBe(4);
		expect(parsed.getUTCSeconds()).toBe(5);
		expect(parsed.getUTCMilliseconds()).toBe(0);

		date = "Mon, 01 Aug 2022 03:04:05 GMT";
		parsed = parseHTTPDate(date);
		expect(parsed.getUTCFullYear()).toBe(2022);
		expect(parsed.getUTCMonth()).toBe(7);
		expect(parsed.getUTCDate()).toBe(1);
		expect(parsed.getUTCHours()).toBe(3);
		expect(parsed.getUTCMinutes()).toBe(4);
		expect(parsed.getUTCSeconds()).toBe(5);
		expect(parsed.getUTCMilliseconds()).toBe(0);

		date = "Thu, 01 Sep 2022 03:04:05 GMT";
		parsed = parseHTTPDate(date);
		expect(parsed.getUTCFullYear()).toBe(2022);
		expect(parsed.getUTCMonth()).toBe(8);
		expect(parsed.getUTCDate()).toBe(1);
		expect(parsed.getUTCHours()).toBe(3);
		expect(parsed.getUTCMinutes()).toBe(4);
		expect(parsed.getUTCSeconds()).toBe(5);
		expect(parsed.getUTCMilliseconds()).toBe(0);

		date = "Sat, 01 Oct 2022 03:04:05 GMT";
		parsed = parseHTTPDate(date);
		expect(parsed.getUTCFullYear()).toBe(2022);
		expect(parsed.getUTCMonth()).toBe(9);
		expect(parsed.getUTCDate()).toBe(1);
		expect(parsed.getUTCHours()).toBe(3);
		expect(parsed.getUTCMinutes()).toBe(4);
		expect(parsed.getUTCSeconds()).toBe(5);
		expect(parsed.getUTCMilliseconds()).toBe(0);

		date = "Tue, 01 Nov 2022 03:04:05 GMT";
		parsed = parseHTTPDate(date);
		expect(parsed.getUTCFullYear()).toBe(2022);
		expect(parsed.getUTCMonth()).toBe(10);
		expect(parsed.getUTCDate()).toBe(1);
		expect(parsed.getUTCHours()).toBe(3);
		expect(parsed.getUTCMinutes()).toBe(4);
		expect(parsed.getUTCSeconds()).toBe(5);
		expect(parsed.getUTCMilliseconds()).toBe(0);

		date = "Thu, 01 Dec 2022 03:04:05 GMT";
		parsed = parseHTTPDate(date);
		expect(parsed.getUTCFullYear()).toBe(2022);
		expect(parsed.getUTCMonth()).toBe(11);
		expect(parsed.getUTCDate()).toBe(1);
		expect(parsed.getUTCHours()).toBe(3);
		expect(parsed.getUTCMinutes()).toBe(4);
		expect(parsed.getUTCSeconds()).toBe(5);
		expect(parsed.getUTCMilliseconds()).toBe(0);
	});

	it("fails to parse invalid dates", ()=>{
		expect(()=>parseHTTPDate("not a date")).toThrow();
		expect(()=>parseHTTPDate("Thu, 01 NaN 2022 03:04:05 GMT")).toThrow();
	});
});
