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

/**
 * Takes the difference between two times (in milliseconds) and returns a formatted string with relative time
 *
 * @param delta time delta in milliseconds
 * @returns Formatted string in the form of 'N X Ago' where X is anything between Seconds and Years
 */
export function relativeTimeString(delta: number): string {
	const SEC = 1000;
	const MIN = SEC * 60;
	const HOUR = MIN * 60;
	const DAY = HOUR * 24;
	const YEAR = 365 * DAY;
	const MONTH = YEAR / 12;
	const WEEK = MONTH / 4;
	if (delta > YEAR) {
		return `${(delta / YEAR).toFixed(2)} years ago`;
	} else if (delta > MONTH) {
		return  `${(delta / MONTH).toFixed(2)} months ago`;
	} else if (delta > WEEK) {
		return `${(delta/ WEEK).toFixed(2)} weeks ago`;
	} else if (delta > DAY) {
		return `${(delta / DAY).toFixed(2)} days ago`;
	} else if (delta > HOUR) {
		return `${(delta / HOUR).toFixed(2)} hours ago`;
	} else if (delta > MIN) {
		return `${(delta / MIN).toFixed(2)} minutes ago`;
	}
	return `${(delta / SEC).toFixed(0)} seconds ago`;

}
