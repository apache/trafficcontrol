/*
 * Copyright 2015 Comcast Cable Communications Management, LLC
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

package org.apache.traffic_control.traffic_router.core.edge;

import java.util.Map;

import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;

import org.apache.traffic_control.traffic_router.geolocation.Geolocation;

/**
 * A physical location that has caches.
 */
public class Location {

	private final String id;
	private final Geolocation geolocation;

	/**
	 * Creates a Location with the specified ID at the specified location.
	 * 
	 * @param id
	 *            the id of the location
	 * @param geolocation
	 *            the coordinates of this location
	 */
	public Location(final String id, final Geolocation geolocation) {
		this.id = id;
		this.geolocation = geolocation;
	}

	@Override
	public boolean equals(final Object obj) {
		if (this == obj) {
			return true;
		} else if (obj instanceof Location) {
			final Location rhs = (Location) obj;
			return new EqualsBuilder()
			.append(getId(), rhs.getId())
			.isEquals();
		} else {
			return false;
		}
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
