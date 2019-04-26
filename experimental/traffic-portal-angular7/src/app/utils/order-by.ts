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

function cmpr(a, b): number {
	if (a === null) {
		if (b === null) {
			return 0;
		}
		return 1;
	} else if (b === null) {
		return -1;
	}

	if (typeof(a) !== typeof(b)) {
		throw new Error();
	}

	if (a === b) {
		return 0;
	}
	if (a < b) {
		return -1;
	}
	return 1;
}

export function orderBy (value: Array<any>, property: string | Array<string>): Array<any> {
	return value.sort((a: any, b: any) => {
		let bail = false;
		if (!a.hasOwnProperty(property)) {
			console.error('object', a, "has no property '" + property + "'!");
			bail = true;
		}
		if (!b.hasOwnProperty(property)) {
			console.error('object', b, "has no property '" + property + "'!");
			bail = true;
		}

		if (bail) {
			return 0;
		}

		let props: Array<string>;
		if (typeof(property) === 'string') {
			props = [property];
		} else {
			props = property;
		}

		for (let p of props) {

			let aProp;
			let bProp;

			try {
				/* tslint:disable */
				aProp = a[p];
				bProp = b[p];
				/* tslint:enable */
			} catch (e) {
				console.error(e);
				return 0;
			}

			let result: number;
			try {
				result = cmpr(aProp, bProp);
			} catch (e) {
				console.error("property '" + p + "' is not the same type on objects", a, 'and', b, '! (' + e.toString() + ')');
				return 0;
			}

			if (result !== 0) {
				return result;
			}
		}

		return 0;

	});
}
