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

package org.apache.traffic_control.traffic_router.core.ds;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Random;
import java.util.SortedMap;

import org.apache.traffic_control.traffic_router.core.util.JsonUtils;
import com.fasterxml.jackson.databind.JsonNode;

import org.apache.traffic_control.traffic_router.core.edge.Cache;

public class Dispersion {
	public final static int DEFAULT_LIMIT = 1;
	public final static boolean DEFAULT_SHUFFLED = true;

	private int limit = DEFAULT_LIMIT;
	private boolean shuffled = DEFAULT_SHUFFLED;

	public Dispersion(final JsonNode dsJo) {
		final JsonNode jo = dsJo.get("dispersion");

		if (jo != null) {
			final int limit = JsonUtils.optInt(jo, "limit", DEFAULT_LIMIT);
			if (limit != 0) {
				this.setLimit(limit);
			}
			this.setShuffled(JsonUtils.optBoolean(jo, "shuffled", DEFAULT_SHUFFLED));
		} else if (dsJo.has("maxDnsIpsForLocation")) {
			// if no specific dispersion, use maxDnsIpsForLocation (should be DNS DSs only)
			this.setLimit(dsJo.get("maxDnsIpsForLocation").asInt(DEFAULT_LIMIT));
		}
	}

	public int getLimit() {
		return limit;
	}

	private void setLimit(final int limit) {
		this.limit = limit;
	}

	public boolean isShuffled() {
		return shuffled;
	}

	private void setShuffled(final boolean shuffled) {
		this.shuffled = shuffled;
	}

	// Used by Http Routing functions
	public Cache getCache(final SortedMap<Double, Cache> cacheMap) {
		if (cacheMap == null) {
			return null;
		}

		final List<Cache> cacheList = this.getCacheList(cacheMap);

		if (cacheList != null) {
			return cacheList.get(0);
		}

		return null;
	}

	// Used by DNS Routing functions
	public List<Cache> getCacheList(final SortedMap<Double, Cache> cacheMap) {
		if (cacheMap == null) {
			return null;
		}

		final List<Cache> cacheList = new ArrayList<Cache>();

		for (final Cache c : cacheMap.values()) {
			cacheList.add(c);

			if (this.getLimit() != 0 && cacheList.size() == this.getLimit()) {
				break;
			}
		}

		if (cacheList.size() > 1 && this.isShuffled()) {
			Collections.shuffle(cacheList, new Random());
		}

		return cacheList;
	}
}
