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

import { environment } from "src/environments/environment";

import { LogLevel, Logger } from "./logging";

/**
 * Implements a single comparison between two values
 *
 * @param a A value to compare.
 * @param b A value to compare.
 * @returns 0 if ``a===b`` or if both a and b are ``null``, -1 if ``a<b`` or b is ``null`` and a is not, otherwise `1`
 * @throws whenever an attempt is made to compare values of different types.
 * This is calculated using ``typeof``, and so only primitive type is considered
 */
function cmpr(a: unknown, b: unknown): number {
	if (a === null) {
		if (b === null) {
			return 0;
		}
		return 1;
	}

	if (b === null) {
		return -1;
	}

	if (typeof(a) !== typeof(b)) {
		throw new Error();
	}

	if (a === b) {
		return 0;
	}
	// These can truly be any type, since all types in JS are comparable with <
	/* eslint-disable @typescript-eslint/no-explicit-any */
	if ((a as any) < (b as any)) {
		return -1;
	}
	/* eslint-enable @typescript-eslint/no-explicit-any */
	return 1;
}

// This can truly be anything, and there's no good way to type it to avoid 'any', so just this once I'm doing it.
/* eslint-disable @typescript-eslint/no-explicit-any */
/**
 * Returns the passed array sorted by the properties of each element as given by
 * the caller.
 *
 * Array elements which are ``undefined`` are unaffected by the sort (uses
 * ``Array.prototype.sort``). Elements are sorted by each element of the
 * property array sequentiall, e.g.
 * > orderBy([{foo: 1, bar: 2}, {foo: 1, bar: 1}], ['foo', 'bar'])
 * [{foo: 1, bar: 1}, {foo: 1, bar: 2}]
 * ``null`` properties are sorted to later positions than not-``null`` properties.
 * Array elements of different types will compare as "equal", but an error will be printed to the console.
 * If a property in the @param{property} array is encountered that has a different type on each object,
 * the objects are immediately considered "equal" without checking any remaining properties - but again
 * an error is printed to the console. Note that type checks are done using ``typeof``, so if properties
 * are not primitive types, they will be considered to be the same type.
 *
 * @param value The array to be sorted
 * @param property Either a single property name or an array of property names to sort by - in descending order of importance.
 * @returns The sorted array
 */
export function orderBy<T extends any>(value: Array<T>, property: string | Array<string>): Array<T> {
	const logger = new Logger(console, environment.production ? LogLevel.INFO : LogLevel.DEBUG, "orderBy call", false);
	return value.sort((a: any, b: any) => {
		/* eslint-enable @typescript-eslint/no-explicit-any */

		let props: Array<string>;
		if (typeof(property) === "string") {
			props = [property];
		} else {
			props = property;
		}

		for (const p of props) {

			let bail = false;
			if (!Object.prototype.hasOwnProperty.call(a, p)) {
				logger.debug("object", a, `has no property "${p}"!`);
				bail = true;
			}
			if (!Object.prototype.hasOwnProperty.call(b, p)) {
				logger.debug("object", b, `has no property "${p}"!`);
				bail = true;
			}

			if (bail) {
				return 0;
			}

			const aProp = a[p];
			const bProp = b[p];

			let result: number;
			try {
				result = cmpr(aProp, bProp);
			} catch (e) {
				logger.debug("property", p, "is not the same type on objects", a, "and", b, `! (${e})`);
				return 0;
			}

			if (result !== 0) {
				return result;
			}
		}

		return 0;

	});
}
