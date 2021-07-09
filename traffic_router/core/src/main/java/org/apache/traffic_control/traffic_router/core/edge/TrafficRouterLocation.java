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

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;

import org.apache.traffic_control.traffic_router.geolocation.Geolocation;

/**
 * A physical location that has caches.
 */
public class TrafficRouterLocation extends Location {
	private final Map<String, Node> trafficRouters;

	/**
	 * Creates a TrafficRouteRLocation with the specified ID at the specified location.
	 * 
	 * @param id
	 *            the id of the location
	 * @param geolocation
	 *            the coordinates of this location
	 */
	public TrafficRouterLocation(final String id, final Geolocation geolocation) {
		super(id, geolocation);
		trafficRouters = new HashMap<String, Node>();
	}

	/**
	 * Adds the specified cache to this location.
	 * 
	 * @param name
	 *            the name of the Traffic Router to add
	 * @param trafficRouter
	 *            the Node representing a Traffic Router
	 */
	public void addTrafficRouter(final String name, final Node trafficRouter) {
			trafficRouters.put(name, trafficRouter);
	}

	@Override
	public boolean equals(final Object obj) {
		if (this == obj) {
			return true;
		} else if (obj instanceof TrafficRouterLocation) {
			final TrafficRouterLocation rhs = (TrafficRouterLocation) obj;
			return new EqualsBuilder()
			.append(getId(), rhs.getId())
			.isEquals();
		} else {
			return false;
		}
	}

	@Override
	public int hashCode() {
		return new HashCodeBuilder(1, 31)
		.append(getId())
		.toHashCode();
	}

	/**
	 * Retrieves the {@link Set} of Traffic Routers at this location.
	 * 
	 * @return the caches
	 */
	public List<Node> getTrafficRouters() {
		return new ArrayList<Node>(trafficRouters.values());
	}
}
