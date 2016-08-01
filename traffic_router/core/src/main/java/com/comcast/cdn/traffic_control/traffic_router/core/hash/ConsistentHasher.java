package com.comcast.cdn.traffic_control.traffic_router.core.hash;

import org.apache.log4j.Logger;

import java.util.ArrayList;
import java.util.Collection;
import java.util.Collections;
import java.util.List;
import java.util.Random;
import java.util.SortedMap;
import java.util.TreeMap;

public class ConsistentHasher {
	private static final Logger LOGGER = Logger.getLogger(ConsistentHasher.class);

	final private MD5HashFunction hashFunction = new MD5HashFunction();

	public <T extends Hashable> T selectHashable(final List<T> hashables, final String s, final boolean shuffle) {
		if (hashables.isEmpty()) {
			LOGGER.warn("Cannot select a hashable from an empty list!");
			return null;
		}

		if (shuffle) {
			return hashables.get(new Random(System.currentTimeMillis()).nextInt(hashables.size()));
		}

		final Collection<T> values = sortHashables(hashables, s).values();

		if (values.isEmpty()) {
			LOGGER.warn("Failed to generate sorted hashables from given hashables list of size " + hashables.size());
			return null;
		}

		return values.iterator().next();
	}

	public <T extends Hashable> List<T> selectHashables(final List<T> hashables, final int limit, final String s, final boolean shuffle) {
		if (shuffle) {
			Collections.shuffle(hashables);
			return (limit <= hashables.size()) ? hashables.subList(0, limit) : hashables;
		}

		final SortedMap<Double, T> sortedHashables = sortHashables(hashables, s);
		final List<T> selectedHashables = new ArrayList<T>();

		for (final T hashable : sortedHashables.values()) {
			if (selectedHashables.size() >= limit) {
				break;
			}

			selectedHashables.add(hashable);
		}

		return selectedHashables;
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
