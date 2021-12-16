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

package org.apache.traffic_control.traffic_router.core.hash;

import java.util.Arrays;

@SuppressWarnings("PMD.ClassNamingConventions")
public class NumberSearcher {
	public static int findClosest(final Double[] numbers, final double target) {
		final int index = Arrays.binarySearch(numbers, target);
		if (index >= 0) {
			return index;
		}

		final int biggerThanIndex = -(index + 1);
		if (biggerThanIndex == numbers.length) {
			return numbers.length - 1;
		}

		if (biggerThanIndex == 0) {
			return 0;
		}

		final int smallerThanIndex = biggerThanIndex - 1;

		final double biggerThanDelta = Math.abs(numbers[biggerThanIndex] - target);
		final double smallerThanDelta = Math.abs(numbers[smallerThanIndex] - target);

		return (biggerThanDelta < smallerThanDelta) ? biggerThanIndex : smallerThanIndex;
	}
}
