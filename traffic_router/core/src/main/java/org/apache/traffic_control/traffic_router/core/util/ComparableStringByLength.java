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

package org.apache.traffic_control.traffic_router.core.util;

public class ComparableStringByLength implements Comparable<ComparableStringByLength> {
	final private String string;

	public ComparableStringByLength(final String string) {
		if (string == null || string.length() == 0) {
			throw new IllegalArgumentException("String parameter must be non-null and non-empty");
		}

		this.string = string;
	}

	@Override
	public int compareTo(final ComparableStringByLength other) {
		if (string.length() == other.string.length()) {
			return string.compareTo(other.string);
		}

		return (string.length() > other.string.length()) ? -1 : 1;
	}

	@Override
	public String toString() {
		return string;
	}

	@Override
	public boolean equals(final Object other) {
		if (this == other) {
			return true;
		}

		if (other == null) {
			return false;
		}

		if (getClass() != other.getClass() && String.class != other.getClass()) {
			return false;
		}

		if (String.class == other.getClass()) {
			return string.equals(other);
		}

		return string.equals(((ComparableStringByLength) other).string);
	}

	@Override
	public int hashCode() {
		return string.hashCode();
	}
}
