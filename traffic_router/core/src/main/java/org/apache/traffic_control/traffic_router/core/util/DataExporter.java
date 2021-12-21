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

package org.apache.traffic_control.traffic_router.core.util;

import java.io.IOException;
import java.io.InputStream;
import java.util.ArrayList;
import java.util.Collection;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Properties;

import com.fasterxml.jackson.databind.JsonNode;
import com.google.common.cache.CacheStats;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import org.apache.traffic_control.traffic_router.core.edge.Cache;
import org.apache.traffic_control.traffic_router.core.edge.CacheLocation;
import org.apache.traffic_control.traffic_router.core.edge.CacheRegister;
import org.apache.traffic_control.traffic_router.core.edge.InetRecord;
import org.apache.traffic_control.traffic_router.core.edge.Location;
import org.apache.traffic_control.traffic_router.core.edge.PropertiesAndCaches;
import org.apache.traffic_control.traffic_router.core.loc.NetworkNode;
import org.apache.traffic_control.traffic_router.core.loc.NetworkNodeException;
import org.apache.traffic_control.traffic_router.core.router.StatTracker;
import org.apache.traffic_control.traffic_router.core.router.TrafficRouter;
import org.apache.traffic_control.traffic_router.core.router.TrafficRouterManager;
import org.apache.traffic_control.traffic_router.core.status.model.CacheModel;
import org.apache.traffic_control.traffic_router.geolocation.Geolocation;
import org.apache.traffic_control.traffic_router.geolocation.GeolocationException;

public class DataExporter {
	private static final Logger LOGGER = LogManager.getLogger(DataExporter.class);
	private static final String NOT_FOUND_MESSAGE = "not found";

	private TrafficRouterManager trafficRouterManager;

	private StatTracker statTracker;

	private FederationExporter federationExporter;

	public void setTrafficRouterManager(final TrafficRouterManager trafficRouterManager) {
		this.trafficRouterManager = trafficRouterManager;
	}

	public void setStatTracker(final StatTracker statTracker) {
		this.statTracker = statTracker;
	}

	public StatTracker getStatTracker() {
		return this.statTracker;
	}

	public Map<String, String> getAppInfo() {
		final Map<String, String> globals = new HashMap<String, String>();
		System.getProperties().keys();

		final Properties props = new Properties();

		try (InputStream stream = getClass().getResourceAsStream("/version.prop")){
			props.load(stream);
		} catch (IOException e) {
			LOGGER.warn(e,e);
		}

		for (final Object key : props.keySet()) {
			globals.put((String) key, props.getProperty((String) key));
		}

		return globals;
	}

	public Map<String, Object> getCachesByIp(final String ip, final String geolocationProvider) {

		final Map<String, Object> map = new HashMap<String, Object>();
		map.put("requestIp", ip);

		final Location cl = getLocationFromCzm(ip);

		if (cl != null) {
			map.put("locationByCoverageZone", cl.getProperties());
		} else {
			map.put("locationByCoverageZone", NOT_FOUND_MESSAGE);
		}

		try {
			final Geolocation gl = trafficRouterManager.getTrafficRouter().getLocation(ip, geolocationProvider, "");

			if (gl != null) {
				map.put("locationByGeo", gl.getProperties());
			} else {
				map.put("locationByGeo", NOT_FOUND_MESSAGE);
			}
		} catch (GeolocationException e) {
			LOGGER.warn(e,e);
			map.put("locationByGeo", e.toString());
		}

		try {
			final CidrAddress cidrAddress = CidrAddress.fromString(ip);
			final List<Object> federationsList = federationExporter.getMatchingFederations(cidrAddress);

			if (federationsList.isEmpty()) {
				map.put("locationByFederation", NOT_FOUND_MESSAGE);
			} else {
				map.put("locationByFederation", federationsList);
			}
		} catch (NetworkNodeException e) {
			map.put("locationByFederation", NOT_FOUND_MESSAGE);
		}

		final CacheLocation clFromDCZ = trafficRouterManager.getTrafficRouter().getDeepCoverageZoneLocationByIP(ip);
		if (clFromDCZ != null) {
			map.put("locationByDeepCoverageZone", new PropertiesAndCaches(clFromDCZ));
		} else {
			map.put("locationByDeepCoverageZone", NOT_FOUND_MESSAGE);
		}

		return map;
	}

	private Location getLocationFromCzm(final String ip) {
		NetworkNode nn = null;

		try {
			nn = NetworkNode.getInstance().getNetwork(ip);
		} catch (NetworkNodeException e) {
			LOGGER.warn(e);
		}

		if (nn == null) { return null; }

		final String locId = nn.getLoc();
		final Location cl = nn.getLocation();

		if (cl != null) {
			return cl;
		}

		if (locId != null) {
			// find CacheLocation
			final TrafficRouter trafficRouter = trafficRouterManager.getTrafficRouter();
			final Collection<CacheLocation> caches = trafficRouter.getCacheRegister().getCacheLocations();

			for (final CacheLocation cl2 : caches) {
				if (cl2.getId().equals(locId)) {
					return cl2;
				}
			}
		}

		return null;
	}

	public List<String> getLocations() {
		final List<String> models = new ArrayList<String>();
		final TrafficRouter trafficRouter = trafficRouterManager.getTrafficRouter();

		for (final CacheLocation location : trafficRouter.getCacheRegister().getCacheLocations()) {
			models.add(location.getId());
		}

		Collections.sort(models);
		return models;
	}

	public List<CacheModel> getCaches(final String locationId) {
		final TrafficRouter trafficRouter = trafficRouterManager.getTrafficRouter();
		final CacheLocation location = trafficRouter.getCacheRegister().getCacheLocation(locationId);
		return getCaches(location);
	}

	public Map<String, Object> getCaches() {
		final Map<String, Object> models = new HashMap<String, Object>();
		final TrafficRouter trafficRouter = trafficRouterManager.getTrafficRouter();

		for (final CacheLocation location : trafficRouter.getCacheRegister().getCacheLocations()) {
			models.put(location.getId(), getCaches(location.getId()));
		}

		return models;
	}

	private List<CacheModel> getCaches(final CacheLocation location) {
		final List<CacheModel> models = new ArrayList<CacheModel>();

		for (final Cache cache : location.getCaches()) {
			final CacheModel model = new CacheModel();
			final List<String> ipAddresses = new ArrayList<String>();
			final List<InetRecord> ips = cache.getIpAddresses(null);

			if (ips != null) {
				for (final InetRecord address : ips) {
					ipAddresses.add(address.getAddress().getHostAddress());
				}
			}

			model.setCacheId(cache.getId());
			model.setFqdn(cache.getFqdn());
			model.setIpAddresses(ipAddresses);

			if (cache.hasAuthority()) {
				model.setCacheOnline(cache.isAvailable());
			} else {
				model.setCacheOnline(false);
			}

			models.add(model);
		}

		return models;
	}

	public int getCacheControlMaxAge() {
		int maxAge = 0;

		if (trafficRouterManager != null) {
			final TrafficRouter trafficRouter = trafficRouterManager.getTrafficRouter();

			if (trafficRouter != null) {
				final CacheRegister cacheRegister = trafficRouter.getCacheRegister();
				final JsonNode config = cacheRegister.getConfig();

				if (config != null) {
					maxAge = JsonUtils.optInt(config, "api.cache-control.max-age");
				}
			}
		}

		return maxAge;
	}

	public Map<String, Object> getStaticZoneCacheStats() {
		return createCacheStatsMap(trafficRouterManager.getTrafficRouter().getZoneManager().getStaticCacheStats());
	}

	public Map<String, Object> getDynamicZoneCacheStats() {
		return createCacheStatsMap(trafficRouterManager.getTrafficRouter().getZoneManager().getDynamicCacheStats());
	}

	private Map<String, Object> createCacheStatsMap(final CacheStats cacheStats) {
		final Map<String, Object> cacheStatsMap = new HashMap<String, Object>();
		cacheStatsMap.put("requestCount", cacheStats.requestCount());
		cacheStatsMap.put("hitCount", cacheStats.hitCount());
		cacheStatsMap.put("missCount", cacheStats.missCount());
		cacheStatsMap.put("hitRate", cacheStats.hitRate());
		cacheStatsMap.put("missRate", cacheStats.missRate());
		cacheStatsMap.put("evictionCount", cacheStats.evictionCount());
		cacheStatsMap.put("loadCount", cacheStats.loadCount());
		cacheStatsMap.put("loadSuccessCount", cacheStats.loadSuccessCount());
		cacheStatsMap.put("loadExceptionCount", cacheStats.loadExceptionCount());
		cacheStatsMap.put("loadExceptionRate", cacheStats.loadExceptionRate());
		cacheStatsMap.put("totalLoadTime", cacheStats.totalLoadTime());
		cacheStatsMap.put("averageLoadPenalty", cacheStats.averageLoadPenalty());
		return cacheStatsMap;
	}

	public void setFederationExporter(final FederationExporter federationExporter) {
		this.federationExporter = federationExporter;
	}
}
