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

/**
 * Returns a fuzzy search "score" for the passed 'text' against the passed
 * 'pattern'. The higher a score is, the later in a fuzzy-sorted/searched list
 * the passed text (or the object it represents) ought to appear.
 *
 * @param text The text being scored
 * @param pattern The pattern being matched against. This should just be a list
 * of characters, not a regular expression or glob or anything.
 * @param threshold An optional threshold which, if given causes text which
 * scores above it to not be returned.
 * @returns The score of the passed text. If not all of the pattern's characters
 * could be found within it (or if it exceeds the optional threshold), this will
 * be Infinity.
 */
export function fuzzyScore(text: string, pattern: string, threshold: number = 0): number {
	if (pattern.length < 1) {
		return 0;
	}

	const p = Array.from(pattern);
	let char = p.shift();
	let score = 0;
	for (const l of text) {
		if (l === char) {
			char = p.shift();
			if (char === undefined) {
				break;
			}
		} else {
			++score;
		}
	}

	if (p.length > 0 || char !== undefined) {
		return Infinity;
	}

	if (threshold > 0 && score > threshold) {
		return Infinity;
	}
	return score;
}
