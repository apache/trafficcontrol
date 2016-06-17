package com.comcast.cdn.traffic_control.traffic_router.core.hash;

import java.util.Arrays;
import java.util.List;
import java.util.TreeSet;

public class DefaultHashable implements Hashable {
	private Double[] hashes;

	@Override
	public double getClosestHash(final double hash) {
		return hashes[NumberSearcher.findClosest(hashes, hash)];
	}

	@Override
	public DefaultHashable generateHashes(final String hashId, final int hashCount) {
		final TreeSet<Double> hashSet = new TreeSet<Double>();
		final MD5HashFunction hashFunction = new MD5HashFunction();

		for (int i = 0; i < hashCount; i++) {
			hashSet.add(hashFunction.hash(hashId + "--" + i));
		}

		hashes = new Double[hashSet.size()];
		System.arraycopy(hashSet.toArray(),0,hashes,0,hashSet.size());
		return this;
	}

	@Override
	public List<Double> getHashValues() {
		return Arrays.asList(hashes);
	}
}
