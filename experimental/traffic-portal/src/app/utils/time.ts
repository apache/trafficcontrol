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

const SEC_MS = 1000;
const MIN_MS = SEC_MS * 60;
const HOUR_MS = MIN_MS * 60;
const DAY_MS = HOUR_MS * 24;
const YEAR_MS = 365 * DAY_MS;
const MONTH_MS = YEAR_MS / 12;
const WEEK_MS = MONTH_MS / 4;
/**
 * Takes the difference between two times (in milliseconds) and returns a formatted string with relative time
 *
 * @param delta time delta in milliseconds
 * @returns Formatted string in the form of 'N X Ago' where X is anything between Seconds and Years
 */
export function relativeTimeString(delta: number): string {
	if (delta > YEAR_MS) {
		return `${(delta / YEAR_MS).toFixed(2)} years ago`;
	} else if (delta > MONTH_MS) {
		return  `${(delta / MONTH_MS).toFixed(2)} months ago`;
	} else if (delta > WEEK_MS) {
		return `${(delta/ WEEK_MS).toFixed(2)} weeks ago`;
	} else if (delta > DAY_MS) {
		return `${(delta / DAY_MS).toFixed(2)} days ago`;
	} else if (delta > HOUR_MS) {
		return `${(delta / HOUR_MS).toFixed(2)} hours ago`;
	} else if (delta > MIN_MS) {
		return `${(delta / MIN_MS).toFixed(2)} minutes ago`;
	}
	return `${(delta / SEC_MS).toFixed(0)} seconds ago`;

}
