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

import java.util.Arrays;
import java.util.List;
import java.util.TreeSet;

public class DefaultHashable implements Hashable {
	private static final int DEFAULT_HASH_COUNT = 1000;
	private Double[] hashes;

	@Override
	public double getClosestHash(final double hash) {
		return hashes[NumberSearcher.findClosest(hashes, hash)];
	}

	@Override
	public DefaultHashable generateHashes(final String hashId, final int hashCount) {
		final TreeSet<Double> hashSet = new TreeSet<Double>();
		final MD5HashFunction hashFunction = new MD5HashFunction();
		final int count = (hashCount > 0) ? hashCount : DEFAULT_HASH_COUNT;

		for (int i = 0; i < count; i++) {
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
