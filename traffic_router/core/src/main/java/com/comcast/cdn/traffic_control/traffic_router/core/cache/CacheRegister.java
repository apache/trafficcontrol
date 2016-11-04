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

import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.TreeSet;

import org.json.JSONObject;

import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryServiceMatcher;
import com.comcast.cdn.traffic_control.traffic_router.core.request.Request;

@SuppressWarnings("PMD.LooseCoupling")
public class CacheRegister implements CacheLocationManager {
	private final Map<String, CacheLocation> configuredLocations;
	private JSONObject trafficRouters;
	private Map<String,Cache> allCaches;
	private TreeSet<DeliveryServiceMatcher> dnsServiceMatchers;
	private TreeSet<DeliveryServiceMatcher> httpServiceMatchers;
	private Map<String, DeliveryService> dsMap;
	private JSONObject config;
	private JSONObject stats;

	public CacheRegister() {
		configuredLocations = new HashMap<String, CacheLocation>();
	}

	@Override
	public CacheLocation getCacheLocation(final String id) {
		return configuredLocations.get(id);
	}

	@Override
	public Set<CacheLocation> getCacheLocations() {
		final Set<CacheLocation> result = new HashSet<CacheLocation>(configuredLocations.size());
		result.addAll(configuredLocations.values());
		return result;
	}

	public CacheLocation getCacheLocationById(final String id) {
		for (final CacheLocation cacheLocation : configuredLocations.values()) {
			if (id.equals(cacheLocation.getId())) {
				return cacheLocation;
			}
		}

		return null;
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

	public void setCacheMap(final Map<String,Cache> map) {
		allCaches = map;
	}

	public Map<String,Cache> getCacheMap() {
		return allCaches;
	}
	
	public void setDnsDeliveryServiceMatchers(final TreeSet<DeliveryServiceMatcher> dnsServices) {
		this.dnsServiceMatchers = dnsServices;
	}

	public void setHttpDeliveryServiceMatchers(final TreeSet<DeliveryServiceMatcher> httpServices) {
		this.httpServiceMatchers = httpServices;
	}

	/**
	 * Gets the first {@link DeliveryService} that matches the {@link Request}.
	 * 
	 * @param request
	 *            the request to match
	 * @return the DeliveryService that matches the request
	 */
	public DeliveryService getDeliveryService(final Request request, final boolean isHttp) {

		TreeSet<DeliveryServiceMatcher> matchers = dnsServiceMatchers;
		if (isHttp) {
			matchers = httpServiceMatchers;
		}

		if (matchers == null) {
			return null;
		}

		for (final DeliveryServiceMatcher m : matchers) {
			if (m.matches(request)) {
				return m.getDeliveryService();
			}
		}

		return null;
	}

	public DeliveryService getDeliveryService(final String deliveryServiceId) {
		return dsMap.get(deliveryServiceId);
	}

	public List<CacheLocation> filterAvailableLocations(final String deliveryServiceId) {
		final DeliveryService deliveryService = dsMap.get(deliveryServiceId);

		if (deliveryService == null) {
			return null;
		}

		return deliveryService.filterAvailableLocations(getCacheLocations());
	}

	public void setDeliveryServiceMap(final Map<String, DeliveryService> dsMap) {
		this.dsMap = dsMap;
	}

	public JSONObject getTrafficRouters() {
		return trafficRouters;
	}
	public void setTrafficRouters(final JSONObject o) {
		trafficRouters = o;
	}

	public void setConfig(final JSONObject o) {
		config = o;
	}
	public JSONObject getConfig() {
		return config;
	}

	public Map<String, DeliveryService> getDeliveryServices() {
		return this.dsMap;
	}

	public JSONObject getStats() {
		return stats;
	}

	public void setStats(final JSONObject stats) {
		this.stats = stats;
	}

}
