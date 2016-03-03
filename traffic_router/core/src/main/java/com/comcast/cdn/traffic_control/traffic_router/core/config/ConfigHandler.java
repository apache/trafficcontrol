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

package com.comcast.cdn.traffic_control.traffic_router.core.config;

import java.io.IOException;
import java.net.UnknownHostException;
import java.net.URL;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.TreeSet;
import java.util.Iterator;
import java.io.PrintWriter;
import java.io.StringWriter;

import com.comcast.cdn.traffic_control.traffic_router.core.loc.FederationsWatcher;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.GeolocationDatabaseUpdater;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.NetworkNode;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.NetworkUpdater;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.RegionalGeoUpdater;

import org.apache.log4j.Logger;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import com.comcast.cdn.traffic_control.traffic_router.core.TrafficRouterException;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.Cache;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheLocation;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheRegister;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.Cache.DeliveryServiceReference;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryServiceMatcher;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryServiceMatcher.Type;
import com.comcast.cdn.traffic_control.traffic_router.core.monitor.TrafficMonitorWatcher;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouterManager;
import com.comcast.cdn.traffic_control.traffic_router.core.util.TrafficOpsUtils;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker;
import com.comcast.cdn.traffic_control.traffic_router.geolocation.Geolocation;
import com.comcast.cdn.traffic_control.traffic_router.core.request.HTTPRequest;


public class ConfigHandler {
	private static final Logger LOGGER = Logger.getLogger(ConfigHandler.class);

	private static long lastSnapshotTimestamp = 0;
	private static Object configSync = new Object();

	private TrafficRouterManager trafficRouterManager;
	private GeolocationDatabaseUpdater geolocationDatabaseUpdater;
	private StatTracker statTracker;
	private String configDir;
	private String trafficRouterId;
	private TrafficOpsUtils trafficOpsUtils;

	private NetworkUpdater networkUpdater;
	private FederationsWatcher federationsWatcher;
	private RegionalGeoUpdater regionalGeoUpdater;

	public String getConfigDir() {
		return configDir;
	}

	public String getTrafficRouterId() {
		return trafficRouterId;
	}

	public GeolocationDatabaseUpdater getGeolocationDatabaseUpdater() {
		return geolocationDatabaseUpdater;
	}
	public NetworkUpdater getNetworkUpdater () {
		return networkUpdater;
	}

	public RegionalGeoUpdater getRegionalGeoUpdater() {
		return regionalGeoUpdater;
	}

	public boolean processConfig(final String jsonStr) throws JSONException, IOException, TrafficRouterException  {
		if (jsonStr == null) {
			trafficRouterManager.setCacheRegister(null);
			return false;
		}

		synchronized(configSync) {
			final JSONObject jo = new JSONObject(jsonStr);
			LOGGER.info("Enter: processConfig");
			final JSONObject config = jo.getJSONObject("config");
			final JSONObject stats = jo.getJSONObject("stats");

			final long sts = getSnapshotTimestamp(stats);

			if (sts <= getLastSnapshotTimestamp()) {
				LOGGER.warn("Incoming TrConfig snapshot timestamp (" + sts + ") is older or equal to the loaded timestamp (" + getLastSnapshotTimestamp() + "); unable to process");
				return false;
			}

			try {
				parseGeolocationConfig(config);
				parseCoverageZoneNetworkConfig(config);
				parseRegionalGeoConfig(config);

				final CacheRegister cacheRegister = new CacheRegister();
				cacheRegister.setTrafficRouters(jo.getJSONObject("contentRouters"));
				cacheRegister.setConfig(config);
				cacheRegister.setStats(stats);
				parseTrafficOpsConfig(config, stats);
				parseDeliveryServiceConfig(jo.getJSONObject("deliveryServices"), cacheRegister);
				parseLocationConfig(jo.getJSONObject("edgeLocations"), cacheRegister);
				parseCacheConfig(jo.getJSONObject("contentServers"), cacheRegister);
				parseMonitorConfig(jo.getJSONObject("monitors"));
				NetworkNode.getInstance().clearCacheLocations();
				federationsWatcher.configure(config);

				trafficRouterManager.setCacheRegister(cacheRegister);
				trafficRouterManager.getTrafficRouter().setRequestHeaders(parseRequestHeaders(config.optJSONArray("requestHeaders")));
				setLastSnapshotTimestamp(sts);
			} catch (ParseException e) {
				LOGGER.error(e, e);
				return false;
			}
		}

		LOGGER.info("Exit: processConfig");

		return true;
	}

	public void setTrafficRouterManager(final TrafficRouterManager trafficRouterManager) {
		this.trafficRouterManager = trafficRouterManager;
	}

	public void setConfigDir(final String configDir) {
		this.configDir = configDir;
	}

	public void setTrafficRouterId(final String traffictRouterId) {
		this.trafficRouterId = traffictRouterId;
	}

	public void setGeolocationDatabaseUpdater(final GeolocationDatabaseUpdater geolocationDatabaseUpdater) {
		this.geolocationDatabaseUpdater = geolocationDatabaseUpdater;
	}
	public void setNetworkUpdater(final NetworkUpdater nu) {
		this.networkUpdater = nu;
	}

	public void setRegionalGeoUpdater(final RegionalGeoUpdater regionalGeoUpdater) {
		this.regionalGeoUpdater = regionalGeoUpdater;
	}

	/**
	 * Parses the Traffic Ops config
	 * @param config
	 *            the {@link TrafficRouterConfiguration} config section
	 * @param stats
	 *            the {@link TrafficRouterConfiguration} stats section
	 *
	 * @throws JSONException 
	 */
	private void parseTrafficOpsConfig(final JSONObject config, final JSONObject stats) throws JSONException {
		if (stats.has("tm_host")) {
			trafficOpsUtils.setHostname(stats.getString("tm_host"));
		} else if (stats.has("to_host")) {
			trafficOpsUtils.setHostname(stats.getString("to_host"));
		} else {
			throw new JSONException("Unable to find to_host or tm_host in stats section of TrConfig; unable to build TrafficOps URLs");
		}

		trafficOpsUtils.setCdnName(stats.getString("CDN_name"));
		trafficOpsUtils.setConfig(config);
	}

	/**
	 * Parses the cache information from the configuration and updates the {@link CacheRegister}.
	 *
	 * @param trConfig
	 *            the {@link TrafficRouterConfiguration}
	 * @throws JSONException 
	 */
	@SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.AvoidDeeplyNestedIfStmts"})
	private void parseCacheConfig(final JSONObject contentServers, final CacheRegister cacheRegister) throws JSONException {
		final Map<String,Cache> map = new HashMap<String,Cache>();
		final Map<String, List<String>> statMap = new HashMap<String, List<String>>();
		for (final String node : JSONObject.getNames(contentServers)) {
			final JSONObject jo = contentServers.getJSONObject(node);
			final CacheLocation loc = cacheRegister.getCacheLocation(jo.getString("locationId"));
			if (loc != null) {
				String hashId = node;
				if(jo.has("hashId")) {
					hashId = jo.optString("hashId");
				}
				final Cache cache = new Cache(node, hashId, jo.optInt("hashCount"));
				cache.setFqdn(jo.getString("fqdn"));
				//				generateCacheIPList(cache);
				cache.setPort(jo.getInt("port"));
//				cache.setAdminStatus(AdminStatus.valueOf(jo.getString("status")));
				final String ip = jo.getString("ip");
				final String ip6 = jo.optString("ip6");
				try {
					cache.setIpAddress(ip, ip6, 0);
				} catch (UnknownHostException e) {
					LOGGER.warn(e+" : "+ip);
				}

				if(jo.has("deliveryServices")) {
					final List<DeliveryServiceReference> references = new ArrayList<Cache.DeliveryServiceReference>();
					final JSONObject dsJos = jo.optJSONObject("deliveryServices");
					for(String ds : JSONObject.getNames(dsJos)) {
						/* technically this could be more than just a string or array,
						 * but, as we only have had those two types, let's not worry about the future
						 */
						final Object dso = dsJos.get(ds);

						List<String> dsNames = statMap.get(ds);

						if (dsNames == null) {
							dsNames = new ArrayList<String>();
						}

						if (dso instanceof JSONArray) {
							final JSONArray fqdnList = (JSONArray) dso;

							if (fqdnList != null && fqdnList.length() > 0) {
								for (int i = 0; i < fqdnList.length(); i++) {
									final String name = fqdnList.getString(i).toLowerCase();

									if (i == 0) {
										references.add(new DeliveryServiceReference(ds, name));
									}

									final String tld = cacheRegister.getConfig().optString("domain_name").toLowerCase();

									if (name.endsWith(tld)) {
										final String reName = name.replaceAll("^.*?\\.", "");

										if (!dsNames.contains(reName)) {
											dsNames.add(reName);
										}
									} else {
										if (!dsNames.contains(name)) {
											dsNames.add(name);
										}
									}
								}
							}
						} else {
							references.add(new DeliveryServiceReference(ds, dso.toString()));

							if (!dsNames.contains(dso.toString())) {
								dsNames.add(dso.toString());
							}
						}

						statMap.put(ds, dsNames);
					}
					cache.setDeliveryServices(references);
				}
				loc.addCache(cache);
				map.put(cache.getId(), cache);
			}
		}
		cacheRegister.setCacheMap(map);
		statTracker.initialize(statMap, cacheRegister);
	}

	/**
	 * Parses the {@link DeliveryService} information from the configuration and updates the
	 * {@link DeliveryServiceManager}.
	 * @param cacheRegister
	 *
	 * @param trConfig
	 *            the {@link TrafficRouterConfiguration}
	 * @throws JSONException 
	 */
	private void parseDeliveryServiceConfig(final JSONObject deliveryServices, final CacheRegister cacheRegister) throws JSONException {
		final TreeSet<DeliveryServiceMatcher> dnsServiceMatchers = new TreeSet<DeliveryServiceMatcher>();
		final TreeSet<DeliveryServiceMatcher> httpServiceMatchers = new TreeSet<DeliveryServiceMatcher>();
		final Map<String,DeliveryService> dsMap = new HashMap<String,DeliveryService>();

		for (String dsId : JSONObject.getNames(deliveryServices)) {
			final JSONObject dsJo = deliveryServices.getJSONObject(dsId);
			final JSONArray matchsets = dsJo.getJSONArray("matchsets");
			final DeliveryService ds = new DeliveryService(dsId, dsJo);
			boolean isDns = false;
			dsMap.put(dsId, ds);

			for (int i = 0; i < matchsets.length(); i++) {
				final JSONObject matchset = matchsets.getJSONObject(i);
				final String protocol = matchset.getString("protocol");

				final DeliveryServiceMatcher deliveryServiceMatcher = new DeliveryServiceMatcher(ds);

				if ("HTTP".equals(protocol)) {
					httpServiceMatchers.add(deliveryServiceMatcher);
				} else if ("DNS".equals(protocol)) {
					dnsServiceMatchers.add(deliveryServiceMatcher);
					isDns = true;
				}

				final JSONArray list = matchset.getJSONArray("matchlist");
				for (int j = 0; j < list.length(); j++) {
					final JSONObject matcherJo = list.getJSONObject(j);
					final Type type = Type.valueOf(matcherJo.getString("match-type"));
					final String target = matcherJo.optString("target");
					deliveryServiceMatcher.addMatch(type, matcherJo.getString("regex"), target);
				}
			}

			ds.setDns(isDns);
		}

		cacheRegister.setDeliveryServiceMap(dsMap);
		cacheRegister.setDnsDeliveryServiceMatchers(dnsServiceMatchers);
		cacheRegister.setHttpDeliveryServiceMatchers(httpServiceMatchers);
		initGeoFailedRedirect(dsMap, cacheRegister);
	}

	private void initGeoFailedRedirect(final Map<String, DeliveryService> dsMap, final CacheRegister cacheRegister) {
		final Iterator<String> itr = dsMap.keySet().iterator();
		while (itr.hasNext()) {
			final DeliveryService ds = dsMap.get(itr.next());
			String type = "INVALID_URL";
			//check if it's relative path or not
			final String rurl = ds.getGeoRedirectUrl();
			if (rurl == null) { continue; }

			try {
				final int idx = rurl.indexOf("://");

				if (idx < 0) {
					//this is a relative url, belongs to this ds
					type = "DS_URL";
				} else {
					//this is a url with protocol, must check further
					//first, parse the url, if url invalid it will throw Exception
					final URL url = new URL(rurl);

					//make a fake HTTPRequest for the redirect url
					final HTTPRequest req = new HTTPRequest();

					req.setPath(url.getPath());
					req.setQueryString(url.getQuery());
					req.setHostname(url.getHost());
					req.setRequestedUrl(rurl);

					ds.setGeoRedirectFile(url.getFile());
					//try select the ds by the redirect fake HTTPRequest
					final DeliveryService rds = cacheRegister.getDeliveryService(req, true);
					if (rds == null) {
						LOGGER.debug("No DeliveryService found for: "
								+ rurl);
						//the redirect url not belongs to any ds
						type = "NOT_DS_URL";
					} else {
						//check if it's the same ds
						if (rds.getId() == ds.getId()) { type = "DS_URL"; }
						else { type = "NOT_DS_URL"; }
					}
				}

				ds.setGeoRedirectUrlType(type);
			} catch (Exception e) {
				LOGGER.error("fatal error, failed to init NGB redirect with Exception: " + e);
				final StringWriter sw = new StringWriter();
				final PrintWriter pw = new PrintWriter(sw);
				e.printStackTrace(pw);
				LOGGER.error(sw.toString());
			}
		}
	}

	/**
	 * Parses the geolocation database configuration and updates the database if the URL has
	 * changed.
	 * 
	 * @param config
	 *            the {@link TrafficRouterConfiguration}
	 * @throws JSONException 
	 */
	private void parseGeolocationConfig(final JSONObject config) throws JSONException {
		String pollingUrlKey = "geolocation.polling.url";

		if (config.has("alt.geolocation.polling.url")) {
			pollingUrlKey = "alt.geolocation.polling.url";
		}

		getGeolocationDatabaseUpdater().setDataBaseURL(
			config.getString(pollingUrlKey),
			config.optLong("geolocation.polling.interval")
		);
	}
	/**
	 * Parses the ConverageZoneNetwork database configuration and updates the database if the URL has
	 * changed.
	 *
	 * @param trConfig
	 *            the {@link TrafficRouterConfiguration}
	 * @throws JSONException 
	 */
	private void parseCoverageZoneNetworkConfig(final JSONObject config) throws JSONException {
		getNetworkUpdater().setDataBaseURL(
				config.getString("coveragezone.polling.url"),
				config.optLong("coveragezone.polling.interval")
			);
	}

	private void parseRegionalGeoConfig(final JSONObject config) {
		final String url = config.optString("regional_geoblock.polling.url", null);
		if (url == null) {
			LOGGER.info("regional_geoblock.polling.url not configured");
			return;
		}

		final long interval = config.optLong("regional_geoblock.polling.interval");
		getRegionalGeoUpdater().setDataBaseURL(url, interval);
	}

	/**
	 * Creates a {@link Map} of location IDs to {@link Geolocation}s for every {@link Location}
	 * included in the configuration that has both a latitude and a longitude specified.
	 *
	 * @param trConfig
	 *            the TrafficRouterConfiguration
	 * @return the {@link Map}, empty if there are no Locations that have both a latitude and
	 *         longitude specified
	 * @throws JSONException 
	 */
	private void parseLocationConfig(final JSONObject locationsJo, final CacheRegister cacheRegister) throws JSONException {
		final Set<CacheLocation> locations = new HashSet<CacheLocation>(locationsJo.length());
		for (final String loc : JSONObject.getNames(locationsJo)) {
			final JSONObject jo = locationsJo.getJSONObject(loc);
			try {
				locations.add(new CacheLocation(loc, jo.optString("zoneId"), 
						new Geolocation(jo.getDouble("latitude"), jo.getDouble("longitude"))));
			} catch (JSONException e) {
				LOGGER.warn(e,e);
			}
		}
		cacheRegister.setConfiguredLocations(locations);
	}

	/**
	 * Creates a {@link Map} of Monitors used by {@link TrafficMonitorWatcher} to pull TrConfigs.
	 *
	 * @param trconfig.monitors
	 *            the monitors section of the TrafficRouter Configuration
	 * @return void
	 * @throws JSONException
	 */
	private void parseMonitorConfig(final JSONObject monitors) throws JSONException, ParseException {
		final List<String> monitorList = new ArrayList<String>();

		for (final String monitorKey : JSONObject.getNames(monitors)) {
			final JSONObject jo = monitors.getJSONObject(monitorKey);
			final String fqdn = jo.getString("fqdn");
			final int port = jo.optInt("port", 80);
			final String status = jo.getString("status");

			if ("ONLINE".equals(status)) {
				monitorList.add(fqdn + ":" + port);
			}
		}

		if (monitorList.isEmpty()) {
			throw new ParseException("Unable to locate any ONLINE monitors in the TrConfig: " + monitors);
		}

		TrafficMonitorWatcher.setOnlineMonitors(monitorList);
	}

	/**
	 * Returns the time stamp (seconds since the epoch) of the TrConfig snapshot.
	 *
	 * @param trconfig.stats
	 *            the stats section of the TrafficRouter Configuration
	 * @return long
	 * @throws JSONException
	 */
	private long getSnapshotTimestamp(final JSONObject stats) throws JSONException {
		return stats.getLong("date");
	}

	public StatTracker getStatTracker() {
		return statTracker;
	}

	public void setStatTracker(final StatTracker statTracker) {
		this.statTracker = statTracker;
	}

	private static long getLastSnapshotTimestamp() {
		return lastSnapshotTimestamp;
	}

	private static void setLastSnapshotTimestamp(final long lastSnapshotTimestamp) {
		ConfigHandler.lastSnapshotTimestamp = lastSnapshotTimestamp;
	}

	public void setFederationsWatcher(final FederationsWatcher federationsWatcher) {
		this.federationsWatcher = federationsWatcher;
	}

	public void setTrafficOpsUtils(final TrafficOpsUtils trafficOpsUtils) {
		this.trafficOpsUtils = trafficOpsUtils;
	}

	private Set<String> parseRequestHeaders(final JSONArray requestHeaders) {
		final Set<String> headers = new HashSet<String>();

		if (requestHeaders == null) {
			return headers;
		}

		for (int i = 0; i < requestHeaders.length(); i++) {
			try {
				headers.add(requestHeaders.getString(i));
			}
			catch (JSONException e) {
				LOGGER.warn("Failed parsing request header from config at position " + i, e);
			}
		}

		return headers;
	}
}
