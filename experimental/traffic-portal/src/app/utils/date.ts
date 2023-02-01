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

/**
 * MalformedDateError is an Error that provides the raw date that caused it as a
 * readable property. Other than that, it's just like an Error.
 */
export class MalformedDateError extends Error {

	public readonly date: string;

	constructor(date: string) {
		super();
		this.message = `malformed date: ${date}`;
		this.date = date;
	}
}

/** Matches both the legacy, custom TO timestamp strings and RFC3339 with optional sub-second precision.  */
const datePattern = /^(\d{4})-(\d{2})-(\d{2})[T ](\d{2}):(\d{2}):(\d{2}(?:\.\d+)?)(?:[\+-]00|Z)$/;

/**
 * dateReviver is meant to be passed into a JSON.parse call as a "reviver"
 * callback. It causes strings that look like dates to be converted to Date
 * objects.
 *
 * This only supports UTC timestamps (which is all Traffic Ops is capable of
 * producing).
 *
 * If a timestamp has sub-milisecond precision, the trailing digits beyond the
 * thousandths place are truncated before parsing.
 *
 * Note that this will do this for **all** strings that look like dates! If, for
 * example, a Delivery Service's LongDescription contains only an RFC3339
 * datestamp, it will be improperly converted!
 *
 * @todo Find a way to specify object keys that should be left alone.
 *
 * @example
 *
 * const data = `{"notADate": "testquest", "myDate": "2022-01-01T00:00:00Z"}`;
 * const parsed = JSON.parse(data, dateReviver);
 * console.log(typeof parsed.notADate); // prints "string"
 * console.log(parsed.myDate instanceof Date); // prints true
 *
 * @param _ The name of the property being parsed - unused here.
 * @param v The value of the property being parsed.
 * @returns Either the parsed date, or just whatever the value is if it's not a
 * string that looks like a date.
 */
export function dateReviver(_: PropertyKey, v: unknown): Date | unknown {
	if (typeof v !== "string") {
		return v;
	}
	const matches = datePattern.exec(v.trim());
	if (!matches) {
		return v;
	}
	const [year, month, day, hour, minute] = matches.slice(1, 6).map(Number);
	let seconds;
	let ms = 0;

	if (matches[6].includes(".")) {
		const [secondsStr, msStr] = matches[6].split(".", 2);
		seconds = Number(secondsStr);
		ms = Number(msStr.slice(0, 3));
	} else {
		seconds = Number(matches[6]);
	}

	const date = new Date(0);
	date.setUTCFullYear(year, month-1, day);
	date.setUTCHours(hour, minute, seconds, ms);
	return date;
}

/** A MonthName is the abbreviated name of a month as it appears in HTTP header dates. */
type MonthName = "Jan"|"Feb"|"Mar"|"Apr"|"May"|"Jun"|"Jul"|"Aug"|"Sep"|"Oct"|"Nov"|"Dec";

/**
 * Checks if a string is a valid abbreviated month name.
 *
 * @param s The string to check.
 * @returns `true` if `s` is a valid abbreviated month name, `false` otherwise.
 */
function isMonthName(s: string): s is MonthName {
	switch(s) {
		case "Jan":
		case "Feb":
		case "Mar":
		case "Apr":
		case "May":
		case "Jun":
		case "Jul":
		case "Aug":
		case "Sep":
		case "Oct":
		case "Nov":
		case "Dec":
			return true;
	}
	return false;
}

/** Index with an abbreviated month name to obtain its number. */
const monthNumbers: Readonly<Record<MonthName, number>> = {
	// Month names are decided by RFC specs, can't be changed to match our conventions.
	// They should also be in the actual order of months in a year, not lexical order.
	/* eslint-disable @typescript-eslint/naming-convention */
	/* eslint-disable sort-keys */
	Jan: 0,
	Feb: 1,
	Mar: 2,
	Apr: 3,
	May: 4,
	Jun: 5,
	Jul: 6,
	Aug: 7,
	Sep: 8,
	Oct: 9,
	Nov: 10,
	Dec: 11
	/* eslint-enable @typescript-eslint/naming-convention */
	/* eslint-enable sort-keys */
};

/** Matches dates as formatted in HTTP headers e.g. "Date", "Last-Modified" etc. */
const httpDatePattern = /^(?:Mon|Tue|Wed|Thu|Fri|Sat|Sun),\s*(\d{2})\s+([A-Za-z]{3})\s+(\d{4})\s+(\d{2}):(\d{2}):(\d{2})\s+GMT$/;

/**
 * Parses a date as formatted in HTTP headers to a Javascript Date.
 *
 * @param raw The raw value of the header (or part of a header) being parsed.
 * @returns The Date represented by `raw`.
 * @throws {MalformedDateError} if `raw` fails to parse.
 */
export function parseHTTPDate(raw: string): Date {
	const matches = httpDatePattern.exec(raw.trim().replace(/\s\s+/g, " "));
	if (!matches || matches.length !== 7 || matches.some(x=>x===undefined)) {
		throw new MalformedDateError(raw);
	}
	const [, d, M, y, h, m, s] = matches;
	const [day, year, hour, minute, second] = [d, y, h, m, s].map(Number);
	if (!isMonthName(M)) {
		throw new MalformedDateError(raw);
	}
	const month = monthNumbers[M];

	const date = new Date(0);
	date.setUTCFullYear(year, month, day);
	date.setUTCHours(hour, minute, second);
	return date;
}
