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

package com.comcast.cdn.traffic_control.traffic_router.core.hash;

import com.comcast.cdn.traffic_control.traffic_router.core.ds.Dispersion;
//import org.apache.log4j.Logger;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.SortedMap;
import java.util.TreeMap;

public class ConsistentHasher {
//	private static final Logger LOGGER = Logger.getLogger(ConsistentHasher.class);

	final private MD5HashFunction hashFunction = new MD5HashFunction();

	public <T extends Hashable> T selectHashable(final List<T> hashables, final Dispersion dispersion, final String s) {
		return selectHashables(hashables, dispersion, s).get(0);
	}

	public <T extends Hashable> List<T> selectHashables(final List<T> hashables, final Dispersion dispersion, final String s) {

		final SortedMap<Double, T> sortedHashables = sortHashables(hashables, s);
		final List<T> selectedHashables = new ArrayList<T>();

		for (final T hashable : sortedHashables.values()) {
			if (selectedHashables.size() >= dispersion.getLimit()) {
				break;
			}

			selectedHashables.add(hashable);
		}
		if (dispersion.isShuffled()) {
			Collections.shuffle(selectedHashables);
		}

		return (dispersion.getLimit() <= selectedHashables.size()) ? selectedHashables.subList(0, dispersion.getLimit()) : selectedHashables;
	}

	private <T extends Hashable> SortedMap<Double, T> sortHashables(final List<T> hashables, final String s) {
		final double hash = hashFunction.hash(s);
		final SortedMap<Double, T> hashableMap = new TreeMap<Double, T>();

		for (final T hashable : hashables) {
			final double closestHash = hashable.getClosestHash(hash);
			double hashDelta = Math.abs(hash - closestHash);

			if (hashableMap.containsKey(hashDelta)) {
				long bits = Double.doubleToLongBits(hashDelta);
				do {
					bits++;
					hashDelta = Double.longBitsToDouble(bits);
				} while (hashableMap.containsKey(hashDelta));
			}

			hashableMap.put(hashDelta, hashable);
		}
		return hashableMap;
	}
}
