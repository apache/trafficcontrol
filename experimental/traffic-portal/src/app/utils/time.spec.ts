/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { relativeTimeString } from "src/app/utils/time";

describe("RelativeTimeString", () => {
	let toDate: Date;
	let fromDate: Date;

	beforeEach(() => {
		// 12/28/2020 23:59:59
		toDate = new Date(1609199999000);
		// 03/15/2015 02:15:03
		fromDate = new Date(1426385703000);
	});
	it("Years", () => {
		const str = relativeTimeString(toDate.getTime() - fromDate.getTime());
		expect(str.substring(4)).toBe(" years ago");
		expect(str.substring(0, 4)).toBeCloseTo(5.8, 1);
	});
	it("Months", () => {
		fromDate.setUTCFullYear(toDate.getUTCFullYear());
		const str = relativeTimeString(toDate.getTime() - fromDate.getTime());
		expect(str.substring(4)).toBe(" months ago");
		expect(str.substring(0, 4)).toBeCloseTo(9.5, 1);
	});
	it("Weeks", () => {
		fromDate.setUTCFullYear(toDate.getUTCFullYear(), toDate.getUTCMonth());
		const str = relativeTimeString(toDate.getTime() - fromDate.getTime());
		expect(str.substring(4)).toBe(" weeks ago");
		expect(str.substring(0, 4)).toBeCloseTo(1.8, 1);
	});
	it("Days", () => {
		fromDate.setUTCFullYear(toDate.getUTCFullYear(), toDate.getUTCMonth(), toDate.getUTCDate()-3);
		const str = relativeTimeString(toDate.getTime() - fromDate.getTime());
		expect(str.substring(4)).toBe(" days ago");
		expect(str.substring(0, 4)).toBeCloseTo(3.9, 1);
	});
	it("Hours", () => {
		fromDate.setUTCFullYear(toDate.getUTCFullYear(), toDate.getUTCMonth(), toDate.getUTCDate());
		fromDate.setUTCHours(toDate.getUTCHours()-3);
		const str = relativeTimeString(toDate.getTime() - fromDate.getTime());
		expect(str.substring(4)).toBe(" hours ago");
		expect(str.substring(0, 4)).toBeCloseTo(3.7, 1);
	});
	it("Minutes", () => {
		fromDate.setUTCFullYear(toDate.getUTCFullYear(), toDate.getUTCMonth(), toDate.getUTCDate());
		fromDate.setUTCHours(toDate.getUTCHours());
		const str = relativeTimeString(toDate.getTime() - fromDate.getTime());
		expect(str.substring(5)).toBe(" minutes ago");
		expect(str.substring(0, 5)).toBeCloseTo(44.9, 1);
	});
	it("Seconds", () => {
		fromDate.setUTCFullYear(toDate.getUTCFullYear(), toDate.getUTCMonth(), toDate.getUTCDate());
		fromDate.setUTCHours(toDate.getUTCHours(), toDate.getUTCMinutes());
		const str = relativeTimeString(toDate.getTime() - fromDate.getTime());
		expect(str.substring(2)).toBe(" seconds ago");
		expect(str.substring(0, 2)).toBe("56");
	});
});
