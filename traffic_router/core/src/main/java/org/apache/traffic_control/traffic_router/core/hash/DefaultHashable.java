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
import java.util.List;
import java.util.TreeSet;

public class DefaultHashable implements Hashable<DefaultHashable>, Comparable<DefaultHashable> {
	private Double[] hashes;
	private int order = 0;

	@Override
	public void setOrder(final int order) {
		this.order = order;
	}

	@Override
	public int getOrder() {
		return order;
	}

	@Override
	public boolean hasHashes() {
		return hashes.length > 0 ? true : false;
	}

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

	@Override
	public int compareTo(final DefaultHashable o) {
		if (this.getOrder() < 0 && o.getOrder() < 0) {
			return getOrder() < o.getOrder() ? 1 : getOrder() > o.getOrder() ? -1 : 0;
		} else {
			return getOrder() < o.getOrder() ? -1 : getOrder() > o.getOrder() ? 1 : 0;
		}
	}
}
