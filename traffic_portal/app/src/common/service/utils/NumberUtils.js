/*


 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.

 */

const sizes = ["B", "Kb", "Mb", "Gb", "Tb", "Pb"];
const k = 1000;

/**
 * NumberUtils provides methods for interacting with various numeric quantities,
 * unitless or otherwise.
 */
class NumberUtils {

	/**
	 *
	 * @param {import("angular").IFilterService} $filter
	 */
	constructor($filter) {
		this.$filter = $filter;
	}

	/**
	 * Adds commas after every three digits (right-to-left) in the decimal
	 * representation of a number.
	 *
	 * @param {`${number}` | number} nStr The number or numeric string to
	 * format.
	 * @returns {string}
	 */
	addCommas(nStr) {
		const x = String(nStr).split('.');
		let x1 = x[0];
		const x2 = x.length > 1 ? "." + x[1] : "";
		const rgx = /(\d+)(\d{3})/;
		while (rgx.test(x1)) {
			x1 = x1.replace(rgx, "$1" + "," + "$2");
		}
		return x1 + x2;
	}

	/**
	 * This function takes big scary kilobit numbers and 'shrinks' them to a
	 * friendly version e.g. 10,000 kilobits is easier read as 10 megabits.
	 *
	 * @param {number} kilounits
	 * @returns {[number, string]} The quantity and its units.
	 */
	shrink(kilounits) {
		if (!kilounits)
			return [0, "Kb"];
		const units = kilounits * 1000;
		let i = Math.floor(Math.log(units) / Math.log(k));
		if (i < 1) { i = 1; } // kilobits is the lowest we will go
		if (i > 5) { i = 5; } // petabits is the highest we will go
		return [Math.round((units / Math.pow(k, i)) * 100) / 100, sizes[i]];
	}

	/**
	 * Converts a number in one set of units to the provided size.
	 *
	 * @param {number} kilounits The quantity to convert.
	 * @param {string} size The size to which to convert `kilounits`. Available
	 * sizes are:
	 * - B - **bits** *not* **bytes**
	 * - Kb - Kilobits
	 * - Mb - Megabits
	 * - Gb - Gigabits
	 * - Tb - Terabits
	 * - Pb - Petabits
	 *
	 * @returns The converted quantity, or zero if the unit/size is not
	 * recognized.
	 */
	convertTo(kilounits, size) {
		if (kilounits === 0)
			return 0;
		const units = kilounits * 1000;
		const i = sizes.indexOf(size);
		if (i === -1) {
			return 0;
		}
		return Math.round((units / Math.pow(k, i)) * 100) / 100;
	}

	/**
	 * Finds the arithmetic mean of a set of values.
	 *
	 * @param {number[]} arr
	 */
	average(arr) {
		return arr.reduce((memo, num) => memo + num, 0) / arr.length;
	}

	/**
	 * Converts a fraction to a unit ratio.
	 *
	 * @example
	 * // returns "N/A"
	 * ratio(5, 0);
	 *
	 * // returns "N/A"
	 * ratio(0, 2);
	 *
	 * // returns "10:1"
	 * ratio(5, 2);
	 *
	 * @param {number} numerator
	 * @param {number} denominator
	 * @returns A unit ratio e.g. "2:1", or "N/A" if a ratio cannot be computed
	 * from the input numbers.
	 */
	ratio(numerator, denominator) {
		if (numerator === 0 || denominator === 0) {
			return "N/A";
		}
		return `${this.$filter("number")(numerator / denominator, 2)}:1`;
	}
}

NumberUtils.$inject = ["$filter"];
module.exports = NumberUtils;
