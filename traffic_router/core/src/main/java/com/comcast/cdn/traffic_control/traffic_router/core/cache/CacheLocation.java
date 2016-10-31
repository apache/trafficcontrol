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

package com.comcast.cdn.traffic_control.traffic_router.core.cache;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;

import com.comcast.cdn.traffic_control.traffic_router.geolocation.Geolocation;

/**
 * A physical location that has caches.
 */
public class CacheLocation {

	private final String id;
	private final Geolocation geolocation;

	private final Map<String, Cache> caches;

	/**
	 * Creates a CacheLocation with the specified ID at the specified location.
	 * 
	 * @param id
	 *            the id of the location
	 * @param geolocation
	 *            the coordinates of this location
	 */
	public CacheLocation(final String id, final Geolocation geolocation) {
		this.id = id;
		this.geolocation = geolocation;
		caches = new HashMap<String, Cache>();
	}

	/**
	 * Adds the specified cache to this location.
	 * 
	 * @param cache
	 *            the cache to add
	 */
	public void addCache(final Cache cache) {
			caches.put(cache.getId(), cache);
	}

	@Override
	public boolean equals(final Object obj) {
		if (this == obj) {
			return true;
		} else if (obj instanceof CacheLocation) {
			final CacheLocation rhs = (CacheLocation) obj;
			return new EqualsBuilder()
			.append(getId(), rhs.getId())
			.isEquals();
		} else {
			return false;
		}
	}

	/**
	 * Retrieves the specified {@link Cache} from the location.
	 * 
	 * @param id
	 *            the ID for the desired <code>Cache</code>
	 * @return the cache or <code>null</code> if the cache doesn't exist
	 */
	public Cache getCache(final String id) {
		return caches.get(id);
	}

	/**
	 * Retrieves the {@link Set} of caches at this location.
	 * 
	 * @return the caches
	 */
	public List<Cache> getCaches() {
		return new ArrayList<Cache>(caches.values());
	}

	/**
	 * Gets geolocation.
	 * 
	 * @return the geolocation
	 */
	public Geolocation getGeolocation() {
		return geolocation;
	}

	/**
	 * Gets id.
	 * 
	 * @return the id
	 */
	public String getId() {
		return id;
	}


	/**
	 * Determines if the specified {@link Cache} exists at this location.
	 * 
	 * @param id
	 *            the <code>Cache</code> to check
	 * @return true if the <code>Cache</code> is at this location, false otherwise
	 */
	public boolean hasCache(final String id) {
		return caches.containsKey(id);
	}

	@Override
	public int hashCode() {
		return new HashCodeBuilder(1, 31)
		.append(getId())
		.toHashCode();
	}

	public Map<String,String> getProperties() {
		final Map<String,String> map = geolocation.getProperties();
		map.put("id", id);
		return map;
	}
}
