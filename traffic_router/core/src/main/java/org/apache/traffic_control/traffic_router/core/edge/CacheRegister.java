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

package org.apache.traffic_control.traffic_router.core.edge;

import java.util.*;
import java.util.stream.Collectors;

import com.fasterxml.jackson.databind.JsonNode;

import org.apache.traffic_control.traffic_router.core.ds.DeliveryService;
import org.apache.traffic_control.traffic_router.core.ds.DeliveryServiceMatcher;
import org.apache.traffic_control.traffic_router.core.request.Request;

@SuppressWarnings("PMD.LooseCoupling")
public class CacheRegister {
	private final Map<String, CacheLocation> configuredLocations;
	private final Map<String, TrafficRouterLocation> edgeTrafficRouterLocations;
	private JsonNode trafficRouters;
	private Map<String,Cache> allCaches;
	private TreeSet<DeliveryServiceMatcher> deliveryServiceMatchers;
	private Map<String, DeliveryService> dsMap;
	private Map<String, DeliveryService> fqdnToDeliveryServiceMap;
	private JsonNode config;
	private JsonNode stats;
	private int edgeTrafficRouterCount;

	public CacheRegister() {
		configuredLocations = new HashMap<String, CacheLocation>();
		edgeTrafficRouterLocations = new HashMap<String, TrafficRouterLocation>();
	}

	public CacheLocation getCacheLocation(final String id) {
		return configuredLocations.get(id);
	}

	public Set<CacheLocation> getCacheLocations() {
		final Set<CacheLocation> result = new HashSet<CacheLocation>(configuredLocations.size());
		result.addAll(configuredLocations.values());
		return result;
	}

	public CacheLocation getCacheLocationById(final String id) {
		for (final CacheLocation location : configuredLocations.values()) {
			if (id.equals(location.getId())) {
				return location;
			}
		}

		return null;
	}

	public TrafficRouterLocation getEdgeTrafficRouterLocation(final String id) {
		return edgeTrafficRouterLocations.get(id);
	}

	public List<TrafficRouterLocation> getEdgeTrafficRouterLocations() {
		return new ArrayList<TrafficRouterLocation>(edgeTrafficRouterLocations.values());
	}

	private void setEdgeTrafficRouterCount(final int count) {
		this.edgeTrafficRouterCount = count;
	}

	public int getEdgeTrafficRouterCount() {
		return edgeTrafficRouterCount;
	}

	public List<Node> getAllEdgeTrafficRouters() {
		final List<Node> edgeTrafficRouters = new ArrayList<>();

		for (final TrafficRouterLocation location : getEdgeTrafficRouterLocations()) {
			edgeTrafficRouters.addAll(location.getTrafficRouters());
		}

		return edgeTrafficRouters;
	}

	/**
	 * Sets the configured locations.
	 * 
	 * @param locations
	 *            the new configured locations
	 */
	public void setConfiguredLocations(final Set<CacheLocation> locations) {
		configuredLocations.clear();
		for (final CacheLocation newLoc : locations) {
			configuredLocations.put(newLoc.getId(), newLoc);
		}
	}

	public void setEdgeTrafficRouterLocations(final Collection<TrafficRouterLocation> locations) {
		int count = 0;

		edgeTrafficRouterLocations.clear();

		for (final TrafficRouterLocation newLoc : locations) {
			edgeTrafficRouterLocations.put(newLoc.getId(), newLoc);

			final List<Node> trafficRouters = newLoc.getTrafficRouters();

			if (trafficRouters != null) {
				count += trafficRouters.size();
			}
		}

		setEdgeTrafficRouterCount(count);
	}

	public boolean hasEdgeTrafficRouters() {
		return !edgeTrafficRouterLocations.isEmpty();
	}

	public void setCacheMap(final Map<String,Cache> map) {
		allCaches = map;
	}

	public Map<String,Cache> getCacheMap() {
		return allCaches;
	}

	public Set<DeliveryServiceMatcher> getDeliveryServiceMatchers(final DeliveryService deliveryService) {
	    return this.deliveryServiceMatchers.stream()
				.filter(deliveryServiceMatcher -> deliveryServiceMatcher.getDeliveryService().getId().equals(deliveryService.getId()))
				.collect(Collectors.toCollection(TreeSet::new));
	}

	public void setDeliveryServiceMatchers(final TreeSet<DeliveryServiceMatcher> matchers) {
		this.deliveryServiceMatchers = matchers;
	}

	/**
	 * Gets the first {@link DeliveryService} that matches the {@link Request}.
	 * 
	 * @param request
	 *            the request to match
	 * @return the DeliveryService that matches the request
	 */
	public DeliveryService getDeliveryService(final Request request) {
		final String requestName = request.getHostname();
		final Map<String, DeliveryService> map = getFQDNToDeliveryServiceMap();
		if (map != null) {
			final DeliveryService ds = map.get(requestName);
			if (ds != null) {
				return ds;
			}
		}
		if (deliveryServiceMatchers == null) {
			return null;
		}

		for (final DeliveryServiceMatcher m : deliveryServiceMatchers) {
			if (m.matches(request)) {
				return m.getDeliveryService();
			}
		}

		return null;
	}

	public DeliveryService getDeliveryService(final String deliveryServiceId) {
		return dsMap.get(deliveryServiceId);
	}

	public List<CacheLocation> filterAvailableCacheLocations(final String deliveryServiceId) {
		final DeliveryService deliveryService = dsMap.get(deliveryServiceId);

		if (deliveryService == null) {
			return null;
		}

		return deliveryService.filterAvailableLocations(getCacheLocations());
	}

	public void setDeliveryServiceMap(final Map<String, DeliveryService> dsMap) {
		this.dsMap = dsMap;
	}

	public Map<String, DeliveryService> getFQDNToDeliveryServiceMap() {
		return fqdnToDeliveryServiceMap;
	}

	public void setFQDNToDeliveryServiceMap(final Map<String, DeliveryService> fqdnToDeliveryServiceMap) {
		this.fqdnToDeliveryServiceMap = fqdnToDeliveryServiceMap;
	}

	public JsonNode getTrafficRouters() {
		return trafficRouters;
	}
	public void setTrafficRouters(final JsonNode o) {
		trafficRouters = o;
	}

	public void setConfig(final JsonNode o) {
		config = o;
	}
	public JsonNode getConfig() {
		return config;
	}

	public Map<String, DeliveryService> getDeliveryServices() {
		return this.dsMap;
	}

	public JsonNode getStats() {
		return stats;
	}

	public void setStats(final JsonNode stats) {
		this.stats = stats;
	}

}
