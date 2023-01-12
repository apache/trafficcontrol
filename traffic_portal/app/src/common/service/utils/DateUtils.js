/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/**
 * @global
 * @typedef {import("moment")} moment
 */

const masks = {
	default: "ddd mmm dd yyyy HH:MM:ss",
	shortDate: "m/d/yy",
	mediumDate: "mmm d, yyyy",
	longDate: "mmmm d, yyyy",
	fullDate: "dddd, mmmm d, yyyy",
	shortTime: "h:MM TT",
	mediumTime: "h:MM:ss TT",
	longTime: "h:MM:ss TT Z",
	isoDate: "yyyy-mm-dd",
	isoTime: "HH:MM:ss",
	isoDateTime: "yyyy-mm-dd'T'HH:MM:ss",
	isoUtcDateTime: "UTC:yyyy-mm-dd'T'HH:MM:ss'Z'"
};

const i18n = {
	dayNames: [
		{
			short: "Sun",
			long: "Sunday"
		},
		{
			short: "Mon",
			long: "Monday"
		},
		{
			short: "Tue",
			long: "Tuesday"
		},
		{
			short: "Wed",
			long: "Wednesday"
		},
		{
			short: "Thu",
			long: "Thursday"
		},
		{
			short: "Fri",
			long: "Friday"
		},
		{
			short: "Sat",
			long: "Saturday"
		},
	],
	monthNames: [
		{
			long: "January",
			short: "Jan"
		},
		{
			long: "February",
			short: "Feb"
		},
		{
			long: "March",
			short: "Mar"
		},
		{
			long: "April",
			short: "Apr"
		},
		{
			long: "May",
			short: "May"
		},
		{
			long: "June",
			short: "Jun"
		},
		{
			long: "July",
			short: "Jul"
		},
		{
			long: "August",
			short: "Aug"
		},
		{
			long: "September",
			short: "Sep"
		},
		{
			long: "October",
			short: "Oct"
		},
		{
			long: "November",
			short: "Nov"
		},
		{
			long: "December",
			short: "Dec"
		},
	]
}

// source: http://blog.stevenlevithan.com/archives/date-time-format
const token = /d{1,4}|m{1,4}|yy(?:yy)?|([HhMsTt])\1?|[LloSZ]|"[^"]*"|'[^']*'/g;
const timezone = /\b(?:[PMCEA][SDP]T|(?:Pacific|Mountain|Central|Eastern|Atlantic) (?:Standard|Daylight|Prevailing) Time|(?:GMT|UTC)(?:[-+]\d{4})?)\b/g;
const timezoneClip = /[^-+\dA-Z]/g

/**
 * DateUtils provides utilities for dealing with dates, either as strings or as
 * actual `Date`s.
 */
class DateUtils {

	/**
	 * Formats a date.
	 *
	 * @param {Date | string} date The date to format - or a string that will be
	 * parsed into a `Date`.
	 * @param {string} [mask] A format to use when formatting `date`. It can
	 * be a literal string, or one of the named masks:
	 * - default: `ddd mmm dd yyyy HH:MM:ss`
	 * - shortDate: `m/d/yy`
	 * - mediumDate: `mmm d, yyyy`
	 * - longDate: `mmmm d, yyyy`
	 * - fullDate: `dddd, mmmm d, yyyy`
	 * - shortTime: `h:MM TT`
	 * - mediumTime: `h:MM:ss TT`
	 * - longTime: `h:MM:ss TT Z`
	 * - isoDate: `yyyy-mm-dd`
	 * - isoTime: `HH:MM:ss`
	 * - isoDateTime: `yyyy-mm-dd'T'HH:MM:ss`
	 * - isoUtcDateTime: `UTC:yyyy-mm-dd'T'HH:MM:ss'Z'`
	 * @param {boolean} [utc]
	 */
	dateFormat(date, mask, utc) {
		if (!date) {
			return '';
		}

		/** @type {string} */
		let maskVal;
		// You can't provide utc if you skip other args (use the "UTC:" mask prefix)
		if (arguments.length === 1 && typeof(date) === "string" && !/\d/.test(date)) {
			maskVal = masks[date];
			date = new Date();
		} else {
			maskVal = masks[mask] ?? mask ?? masks["default"];
		}

		// Passing date through Date applies Date.parse, if necessary
		date = date ? new Date(date) : new Date();
		if (isNaN(date.valueOf()))
			throw new SyntaxError("invalid date");

		// Allow setting the utc argument via the mask
		if (maskVal.slice(0, 4) === "UTC:") {
			maskVal = maskVal.slice(4);
			utc = true;
		}

		let d, D, m, y, H, M, s, L, o
		if (utc) {
			d = date.getUTCDate();
			D = date.getUTCDay();
			m = date.getUTCMonth();
			y = date.getUTCFullYear();
			H = date.getUTCHours();
			M = date.getUTCMinutes();
			s = date.getUTCSeconds();
			L = date.getUTCMilliseconds();
			o = 0;
		} else {
			d = date.getDate();
			D = date.getDay();
			m = date.getMonth();
			y = date.getFullYear();
			H = date.getHours();
			M = date.getMinutes();
			s = date.getSeconds();
			L = date.getMilliseconds();
			o = date.getTimezoneOffset();
		}
		const flags = {
			d,
			dd: String(d).padStart(2, "0"),
			ddd: i18n.dayNames[D].short,
			dddd: i18n.dayNames[D].long,
			m: m + 1,
			mm: String(m + 1).padStart(2, "0"),
			mmm: i18n.monthNames[m].short,
			mmmm: i18n.monthNames[m].long,
			yy: String(y).slice(2),
			yyyy: y,
			h: H % 12 || 12,
			hh: String(H % 12 || 12).padStart(2, "0"),
			H,
			HH: String(H).padStart(2, "0"),
			M,
			MM: String(M).padStart(2, "0"),
			s,
			ss: String(s).padStart(2, "0"),
			l: String(L).padStart(3, "0"),
			L: String(L > 99 ? Math.round(L / 10) : L).padStart(2, "0"),
			t: H < 12 ? "a" : "p",
			tt: H < 12 ? "am" : "pm",
			T: H < 12 ? "A" : "P",
			TT: H < 12 ? "AM" : "PM",
			Z: utc ? "UTC" : (String(date).match(timezone)?.pop() ?? "").replace(timezoneClip, ""),
			o: (o > 0 ? "-" : "+") + String(Math.floor(Math.abs(o) / 60) * 100 + Math.abs(o) % 60).padStart(4, "0"),
			S: ["th", "st", "nd", "rd"][d % 10 > 3 ? 0 : ((d % 100 - d % 10) !== 10 ? 1 : 0) * d % 10]
		};

		return maskVal.replace(token, $0 => flags[$0] ?? $0.slice(1, $0.length - 1));
	}

	/**
	 * Converts a date into a string that tells how much time is between the
	 * current time and the given date.
	 *
	 * @example
	 * // returns "1 hour ago"
	 * getRelativeTime(new Date(Date.now() - 60*60*1000));
	 *
	 * // returns "1 hour from now"
	 * getRelativeTime(new Date(Date.now() + 60*60*1000));
	 *
	 * @param {Date | string} date Either a Date object or a string that can be
	 * parsed by momentjs.
	 * @returns {string} A human readable description of how much time is
	 * between now and `date`.
	 */
	getRelativeTime(date) {
		return moment(date).fromNow();
	};

	/**
	 * Converts a date into a string that tells how much time is between the
	 * current time and the given date.
	 *
	 * This is meant for describing the time at which a user last logged in, so
	 * when the date isn't given it reports "Never logged in".
	 *
	 * @example
	 * // returns "1 hour ago"
	 * relativeLoginTime(new Date(Date.now() - 60*60*1000));
	 *
	 * // returns "Never logged in"
	 * relativeLoginTime(null);
	 *
	 * @param {Date | string | null | undefined} [date] Either a Date object or
	 * a string that can be
	 * parsed by momentjs.
	 * @returns {string} A human readable description of how much time is
	 * between now and `date`.
	 */
	relativeLoginTime(date) {
		if (date) {
			return this.getRelativeTime(date);
		}
		return "Never logged in";
	}
}

DateUtils.$inject = [];
module.exports = DateUtils;
