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
		toDate = new Date();
		toDate.setTime(1609199999000);
		// 03/15/2015 02:15:03
		fromDate = new Date();
		fromDate.setTime(1426385703000);
	});
	it("Years", () => {
		const str = relativeTimeString(toDate.getTime() - fromDate.getTime());
		expect(str.substring(4)).toBe(" years ago");
		expect(str.substring(0, 4)).toBe("5.80");
	});
	it("Months", () => {
		fromDate.setFullYear(toDate.getFullYear());
		const str = relativeTimeString(toDate.getTime() - fromDate.getTime());
		expect(str.substring(4)).toBe(" months ago");
		expect(str.substring(0, 4)).toBe("9.50");
	});
	it("Weeks", () => {
		fromDate.setFullYear(toDate.getFullYear(), toDate.getMonth());
		const str = relativeTimeString(toDate.getTime() - fromDate.getTime());
		expect(str.substring(4)).toBe(" weeks ago");
		expect(str.substring(0, 4)).toBe("1.82");
	});
	it("Days", () => {
		fromDate.setFullYear(toDate.getFullYear(), toDate.getMonth(), toDate.getDate()-3);
		const str = relativeTimeString(toDate.getTime() - fromDate.getTime());
		expect(str.substring(4)).toBe(" days ago");
		expect(str.substring(0, 4)).toBe("2.86");
	});
	it("Hours", () => {
		fromDate.setFullYear(toDate.getFullYear(), toDate.getMonth(), toDate.getDate());
		fromDate.setHours(toDate.getHours()-3);
		const str = relativeTimeString(toDate.getTime() - fromDate.getTime());
		expect(str.substring(4)).toBe(" hours ago");
		expect(str.substring(0, 4)).toBe("3.75");
	});
	it("Minutes", () => {
		fromDate.setFullYear(toDate.getFullYear(), toDate.getMonth(), toDate.getDate());
		fromDate.setHours(toDate.getHours());
		const str = relativeTimeString(toDate.getTime() - fromDate.getTime());
		expect(str.substring(5)).toBe(" minutes ago");
		expect(str.substring(0, 5)).toBe("44.93");
	});
	it("Seconds", () => {
		fromDate.setFullYear(toDate.getFullYear(), toDate.getMonth(), toDate.getDate());
		fromDate.setHours(toDate.getHours(), toDate.getMinutes());
		const str = relativeTimeString(toDate.getTime() - fromDate.getTime());
		expect(str.substring(2)).toBe(" seconds ago");
		expect(str.substring(0, 2)).toBe("56");
	});
});
