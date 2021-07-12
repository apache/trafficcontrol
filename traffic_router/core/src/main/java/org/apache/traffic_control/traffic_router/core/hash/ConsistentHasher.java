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

import org.apache.traffic_control.traffic_router.core.ds.Dispersion;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.NoSuchElementException;
import java.util.SortedMap;
import java.util.TreeMap;

public class ConsistentHasher {
	final private MD5HashFunction hashFunction = new MD5HashFunction();

	public <T extends Hashable> T selectHashable(final List<T> hashables, final Dispersion dispersion, final String s) {
		final List<T> selectedHashables = selectHashables(hashables, dispersion, s);
		return !selectedHashables.isEmpty() ? selectedHashables.get(0) : null;
	}

	public <T extends Hashable> List<T> selectHashables(final List<T> hashables, final String s) {
		return selectHashables(hashables, null, s);
	}

	public <T extends Hashable> List<T> selectHashables(final List<T> hashables, final Dispersion dispersion, final String s) {
		final SortedMap<Double, T> sortedHashables = sortHashables(hashables, s);
		final List<T> selectedHashables = new ArrayList<T>();

		for (final T hashable : sortedHashables.values()) {
			if (dispersion != null && selectedHashables.size() >= dispersion.getLimit()) {
				break;
			}

			selectedHashables.add(hashable);
		}

		if (dispersion != null && dispersion.isShuffled()) {
			Collections.shuffle(selectedHashables);
		}

		return selectedHashables;
	}

	@SuppressWarnings("PMD.EmptyCatchBlock")
	private <T extends Hashable> SortedMap<Double, T> sortHashables(final List<T> hashables, final String s) {
		final double hash = hashFunction.hash(s);
		final SortedMap<Double, T> hashableMap = new TreeMap<Double, T>();
		final List<T> zeroHashes = new ArrayList<T>();

		for (final T hashable : hashables) {
			if (!hashable.hasHashes()) {
				zeroHashes.add(hashable);
				continue;
			}

			final double closestHash = hashable.getClosestHash(hash);
			final double hashDelta = getSafePositiveHash(hashableMap, Math.abs(hash - closestHash));

			hashableMap.put(hashDelta, hashable);
		}

		return synthesizeZeroHashes(hashableMap, zeroHashes);
	}

	/*
	 * The following provides the ability to use zero weights/hashCounts, with or without ordering. The primary
	 * use case is for multi-location routing, but this could also apply to caches. See TC-261.
	 * Because this method returns a SortedMap, we need a means to find the "lowest" and "highest" values in the
	 * hashableMap, then decrement or increment that number within the bounds of Double such that we don't wrap.
	 * Wrapping is dangerous, as it could cause something intended for the tail of the list to appear at the head.
	 */
	@SuppressWarnings({"PMD.EmptyCatchBlock"})
	private <T extends Hashable> SortedMap<Double, T> synthesizeZeroHashes(final SortedMap<Double, T> hashableMap, final List<T> zeroHashes) {
		if (zeroHashes.isEmpty()) {
			return hashableMap;
		}

		double minHash = 0;
		double maxHash = 0;

		try {
			minHash = hashableMap.firstKey();
			maxHash = hashableMap.lastKey();
		} catch (NoSuchElementException ex) {
			// hashableMap is empty; ignore
		}

		Collections.sort(zeroHashes); // sort by order if specified, default is 0 if unspecified

		// add any hashables that don't have hashes to the head/tail of the SortedMap
		for (final T hashable : zeroHashes) {
			if (hashable.getOrder() >= 0) { // append
				final double syntheticHash = getSafePositiveHash(hashableMap, maxHash);
				hashableMap.put(syntheticHash, hashable);
				maxHash = syntheticHash;
			} else { // negative order specified, prepend
				final double syntheticHash = getSafeNegativeHash(hashableMap, minHash);
				hashableMap.put(syntheticHash, hashable);
				minHash = syntheticHash;
			}
		}

		return hashableMap;
	}

	private <T extends Hashable> double getSafePositiveHash(final SortedMap<Double, T> hashableMap, final double hash) {
		return getSafeHash(hashableMap, hash, true);
	}

	private <T extends Hashable> double getSafeNegativeHash(final SortedMap<Double, T> hashableMap, final double hash) {
		return getSafeHash(hashableMap, hash, false);
	}

	private <T extends Hashable> double getSafeHash(final SortedMap<Double, T> hashableMap, final double hash, final boolean add) {
		if (!hashableMap.containsKey(hash)) {
			return hash;
		}

		double syntheticHash = hash;
		long bits = Double.doubleToLongBits(syntheticHash);
		do {
			bits = (add) ? ++bits : --bits;
			syntheticHash = Double.longBitsToDouble(bits);
		} while (hashableMap.containsKey(syntheticHash));

		/*
		 * This shouldn't happen unless we wrap, return safest option if we do, replacing whatever key exists.
		 * If we return a wrapped value, we could incorrectly put the hashable at the head or tail of the SortedMap.
		 */
		if (add && syntheticHash < hash || !add && syntheticHash > hash) {
			return hash;
		}

		return syntheticHash;
	}
}
