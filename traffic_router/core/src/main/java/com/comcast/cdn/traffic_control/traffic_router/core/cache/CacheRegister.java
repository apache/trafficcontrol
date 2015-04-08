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

package com.comcast.cdn.traffic_control.traffic_router.core.cache;

import java.util.ArrayList;
import java.util.Collection;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;

import org.apache.log4j.Logger;
import org.json.JSONObject;

import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryServiceMatcher;
import com.comcast.cdn.traffic_control.traffic_router.core.request.Request;

public class CacheRegister implements CacheLocationManager {
	private final Map<String, CacheLocation> configuredLocations;
	private JSONObject trafficRouters;
	private Map<String,Cache> allCaches;
	private List<DeliveryServiceMatcher> dnsServiceMatchers;
	private List<DeliveryServiceMatcher> httpServiceMatchers;
	private Map<String, DeliveryService> dsMap;
	private JSONObject config;

	public CacheRegister() {
		configuredLocations = new HashMap<String, CacheLocation>();
	}

	/*
	 * (non-Javadoc)
	 * 
	 * @see com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheLocationManager#getCacheLocation(java.lang.String)
	 */
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

	@Override
	@SuppressWarnings("PMD")
	public Collection<CacheLocation> getCacheLocations(String zoneId) {
		if(zoneId == null) { zoneId = ""; }
		final List<CacheLocation> result = new ArrayList<CacheLocation>(configuredLocations.size());
		for (final CacheLocation location : configuredLocations.values()) {
			if (strsEqual(location.getZoneId(),zoneId)) {
				result.add(location);
			}
		}
		return result;
	}
	@SuppressWarnings("PMD")
	private boolean strsEqual(String a, String b) {
		if(a == null) { a = ""; }
		if(b == null) { b = ""; }
		return a.equals(b);
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

	public void print(final Logger log) {
		final Collection<CacheLocation> locs = configuredLocations.values();
		for(CacheLocation loc : locs) {
			final String id = loc.getId();
			final List<Cache> caches = loc.getCaches();
			log.warn("CacheLocation: "+id);
			for(Cache c : caches) {
				log.warn("\tCache: "+c.getId() + " : " + c.toString());
			}
		}
	}

	public void setCacheMap(final Map<String,Cache> map) {
		allCaches = map;
	}
	public void mergeState(final CacheRegister cacheRegister) {
		if(allCaches != null && cacheRegister.allCaches != null) {
			for(String key : allCaches.keySet()) {
				final Cache orig = cacheRegister.allCaches.get(key);
				final Cache n = allCaches.get(key);
				if(orig == null) { continue; }
				if(orig.hasAuthority) {
					n.setIsAvailable(orig.isAvailable);
				}
			}
		}
	}
	public Map<String,Cache> getCacheMap() {
		return allCaches;
	}
	
	public void setDnsDeliveryServiceMatchers(final List<DeliveryServiceMatcher> dnsServices) {
		this.dnsServiceMatchers = dnsServices;
	}
	public void setHttpDeliveryServiceMatchers(final List<DeliveryServiceMatcher> httpServices) {
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
		DeliveryService result = null;

		List<DeliveryServiceMatcher> matchers = null;
		if (isHttp) {
			matchers = httpServiceMatchers;
		} else {
			matchers = dnsServiceMatchers;
		}

		if (matchers != null) {
			for (final DeliveryServiceMatcher m : matchers) {
				if (m.matches(request)) {
					result = m.getDeliveryService();
					break;
				}
			}
		}

		return result;
	}

	public DeliveryService getDeliveryService(final String deliveryServiceId) {
		return dsMap.get(deliveryServiceId);
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

}
