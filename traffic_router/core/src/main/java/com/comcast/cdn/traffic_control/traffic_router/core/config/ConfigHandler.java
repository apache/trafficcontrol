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

package com.comcast.cdn.traffic_control.traffic_router.core.config;

import java.io.IOException;
import java.net.UnknownHostException;
import java.net.URL;
import java.util.ArrayList;
import java.util.Date;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.TreeSet;
import java.util.Iterator;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.stream.Collectors;

import com.comcast.cdn.traffic_control.traffic_router.core.ds.SteeringWatcher;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.FederationsWatcher;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.GeolocationDatabaseUpdater;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.NetworkNode;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.NetworkUpdater;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.RegionalGeoUpdater;

import com.comcast.cdn.traffic_control.traffic_router.core.secure.CertificatesPoller;
import com.comcast.cdn.traffic_control.traffic_router.core.secure.CertificatesPublisher;
import org.apache.log4j.Logger;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

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

@SuppressWarnings("PMD.TooManyFields")
public class ConfigHandler {
	private static final Logger LOGGER = Logger.getLogger(ConfigHandler.class);

	private static long lastSnapshotTimestamp = 0;
	private static Object configSync = new Object();
	private static String deliveryServicesKey = "deliveryServices";

	private TrafficRouterManager trafficRouterManager;
	private GeolocationDatabaseUpdater geolocationDatabaseUpdater;
	private StatTracker statTracker;
	private String configDir;
	private String trafficRouterId;
	private TrafficOpsUtils trafficOpsUtils;

	private NetworkUpdater networkUpdater;
	private FederationsWatcher federationsWatcher;
	private RegionalGeoUpdater regionalGeoUpdater;
	private SteeringWatcher steeringWatcher;
	private CertificatesPoller certificatesPoller;
	private CertificatesPublisher certificatesPublisher;
	private BlockingQueue<Boolean> publishStatusQueue;
	private final AtomicBoolean cancelled = new AtomicBoolean(false);
	private final AtomicBoolean isProcessing = new AtomicBoolean(false);

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

	@SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.NPathComplexity", "PMD.AvoidCatchingThrowable"})
	public boolean processConfig(final String jsonStr) throws JSONException, IOException  {
		isProcessing.set(true);
		LOGGER.info("Entered processConfig");
		if (jsonStr == null) {
			trafficRouterManager.setCacheRegister(null);
			cancelled.set(false);
			isProcessing.set(false);
			publishStatusQueue.clear();
			LOGGER.info("Exiting processConfig: No json data to process");
			return false;
		}

		Date date;
		synchronized(configSync) {
			final JSONObject jo = new JSONObject(jsonStr);
			final JSONObject config = jo.getJSONObject("config");
			final JSONObject stats = jo.getJSONObject("stats");

			final long sts = getSnapshotTimestamp(stats);
			date = new Date(sts * 1000L);

			if (sts <= getLastSnapshotTimestamp()) {
				cancelled.set(false);
				isProcessing.set(false);
				publishStatusQueue.clear();
				LOGGER.info("Exiting processConfig: Incoming TrConfig snapshot timestamp (" + sts + ") is older or equal to the loaded timestamp (" + getLastSnapshotTimestamp() + "); unable to process");
				return false;
			}

			try {
				parseGeolocationConfig(config);
				parseCoverageZoneNetworkConfig(config);
				parseRegionalGeoConfig(jo);

				final CacheRegister cacheRegister = new CacheRegister();
				final JSONObject deliveryServicesJson = jo.getJSONObject("deliveryServices");
				cacheRegister.setTrafficRouters(jo.getJSONObject("contentRouters"));
				cacheRegister.setConfig(config);
				cacheRegister.setStats(stats);
				parseTrafficOpsConfig(config, stats);

				final Map<String, DeliveryService> deliveryServiceMap = parseDeliveryServiceConfig(jo.getJSONObject(deliveryServicesKey));

				parseCertificatesConfig(config);
				certificatesPublisher.setDeliveryServicesJson(deliveryServicesJson);
				final ArrayList<DeliveryService> deliveryServices = new ArrayList<>();

				if (deliveryServiceMap != null && !deliveryServiceMap.values().isEmpty()) {
					deliveryServices.addAll(deliveryServiceMap.values());
				}

				if (deliveryServiceMap != null && !deliveryServiceMap.values().isEmpty()) {
					certificatesPublisher.setDeliveryServices(deliveryServices);
				}

				certificatesPoller.restart();

				final List<DeliveryService> httpsDeliveryServices = deliveryServices.stream().filter(ds -> !ds.isDns() && ds.isSslEnabled()).collect(Collectors.toList());
				httpsDeliveryServices.forEach(ds -> LOGGER.info("Checking for certificate for " + ds.getId()));

				if (!httpsDeliveryServices.isEmpty()) {
					try {
						publishStatusQueue.put(true);
					} catch (InterruptedException e) {
						LOGGER.warn("Failed to notify certificates publisher we're waiting for certificates", e);
					}
				}

				if (cancelled.get()) {
					cancelled.set(false);
					isProcessing.set(false);
					publishStatusQueue.clear();
					LOGGER.info("Exiting processConfig: processing of config with timestamp " + date + " was cancelled");
					return false;
				}

				parseDeliveryServiceMatchSets(deliveryServicesJson, deliveryServiceMap, cacheRegister);
				parseLocationConfig(jo.getJSONObject("edgeLocations"), cacheRegister);
				parseCacheConfig(jo.getJSONObject("contentServers"), cacheRegister);
				parseMonitorConfig(jo.getJSONObject("monitors"));
				NetworkNode.getInstance().clearCacheLocations();
				federationsWatcher.configure(config);
				steeringWatcher.configure(config);
				steeringWatcher.setCacheRegister(cacheRegister);
				trafficRouterManager.setCacheRegister(cacheRegister);
				trafficRouterManager.getTrafficRouter().setRequestHeaders(parseRequestHeaders(config.optJSONArray("requestHeaders")));
				trafficRouterManager.getTrafficRouter().configurationChanged();
				setLastSnapshotTimestamp(sts);
			} catch (ParseException e) {
				isProcessing.set(false);
				cancelled.set(false);
				publishStatusQueue.clear();
				LOGGER.error("Exiting processConfig: Failed to process config for snapshot from " + date, e);
				return false;
			}
		}

		LOGGER.info("Exit: processConfig, successfully applied snapshot from " + date);
		isProcessing.set(false);
		cancelled.set(false);
		publishStatusQueue.clear();
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
	private void parseCacheConfig(final JSONObject contentServers, final CacheRegister cacheRegister) throws JSONException, ParseException {
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

				if(jo.has(deliveryServicesKey)) {
					final List<DeliveryServiceReference> references = new ArrayList<Cache.DeliveryServiceReference>();
					final JSONObject dsJos = jo.optJSONObject(deliveryServicesKey);
					for (final String ds : JSONObject.getNames(dsJos)) {
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

	private Map<String, DeliveryService> parseDeliveryServiceConfig(final JSONObject allDeliveryServices) throws JSONException {
		final Map<String,DeliveryService> deliveryServiceMap = new HashMap<>();

		for (final String deliveryServiceId : JSONObject.getNames(allDeliveryServices)) {
			final JSONObject deliveryServiceJson = allDeliveryServices.getJSONObject(deliveryServiceId);
			final DeliveryService deliveryService = new DeliveryService(deliveryServiceId, deliveryServiceJson);
			boolean isDns = false;

			final JSONArray matchsets = deliveryServiceJson.getJSONArray("matchsets");

			for (int i = 0; i < matchsets.length(); i++) {
				final JSONObject matchset = matchsets.getJSONObject(i);
				final String protocol = matchset.getString("protocol");
				if ("DNS".equals(protocol)) {
					isDns = true;
				}
			}

			deliveryService.setDns(isDns);
			deliveryServiceMap.put(deliveryServiceId, deliveryService);
		}

		return deliveryServiceMap;
	}

	private void parseDeliveryServiceMatchSets(final JSONObject allDeliveryServices, final Map<String, DeliveryService> deliveryServiceMap, final CacheRegister cacheRegister) throws JSONException {
		final TreeSet<DeliveryServiceMatcher> dnsServiceMatchers = new TreeSet<>();
		final TreeSet<DeliveryServiceMatcher> httpServiceMatchers = new TreeSet<>();

		for (final String deliveryServiceId : JSONObject.getNames(allDeliveryServices)) {
			final JSONObject deliveryServiceJson = allDeliveryServices.getJSONObject(deliveryServiceId);
			final JSONArray matchsets = deliveryServiceJson.getJSONArray("matchsets");
			final DeliveryService deliveryService = deliveryServiceMap.get(deliveryServiceId);

			for (int i = 0; i < matchsets.length(); i++) {
				final JSONObject matchset = matchsets.getJSONObject(i);
				final String protocol = matchset.getString("protocol");

				final DeliveryServiceMatcher deliveryServiceMatcher = new DeliveryServiceMatcher(deliveryService);

				if ("HTTP".equals(protocol)) {
					httpServiceMatchers.add(deliveryServiceMatcher);
				} else if ("DNS".equals(protocol)) {
					dnsServiceMatchers.add(deliveryServiceMatcher);
				}

				final JSONArray list = matchset.getJSONArray("matchlist");
				for (int j = 0; j < list.length(); j++) {
					final JSONObject matcherJo = list.getJSONObject(j);
					final Type type = Type.valueOf(matcherJo.getString("match-type"));
					final String target = matcherJo.optString("target");
					deliveryServiceMatcher.addMatch(type, matcherJo.getString("regex"), target);
				}
			}
		}

		cacheRegister.setDeliveryServiceMap(deliveryServiceMap);
		cacheRegister.setDnsDeliveryServiceMatchers(dnsServiceMatchers);
		cacheRegister.setHttpDeliveryServiceMatchers(httpServiceMatchers);
		initGeoFailedRedirect(deliveryServiceMap, cacheRegister);
	}

	private void initGeoFailedRedirect(final Map<String, DeliveryService> dsMap, final CacheRegister cacheRegister) {
		final Iterator<String> itr = dsMap.keySet().iterator();
		while (itr.hasNext()) {
			final DeliveryService ds = dsMap.get(itr.next());
			//check if it's relative path or not
			final String rurl = ds.getGeoRedirectUrl();
			if (rurl == null) { continue; }

			try {
				final int idx = rurl.indexOf("://");

				if (idx < 0) {
					//this is a relative url, belongs to this ds
					ds.setGeoRedirectUrlType("DS_URL");
					continue;
				}
				//this is a url with protocol, must check further
				//first, parse the url, if url invalid it will throw Exception
				final URL url = new URL(rurl);

				//make a fake HTTPRequest for the redirect url
				final HTTPRequest req = new HTTPRequest(url);

				ds.setGeoRedirectFile(url.getFile());
				//try select the ds by the redirect fake HTTPRequest
				final DeliveryService rds = cacheRegister.getDeliveryService(req, true);
				if (rds == null || rds.getId() != ds.getId()) {
					//the redirect url not belongs to this ds
					ds.setGeoRedirectUrlType("NOT_DS_URL");
					continue;
				}

				ds.setGeoRedirectUrlType("DS_URL");
			} catch (Exception e) {
				LOGGER.error("fatal error, failed to init NGB redirect with Exception: " + e.getMessage());
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

		if (config.has("neustar.polling.url")) {
			System.setProperty("neustar.polling.url", config.getString("neustar.polling.url"));
		}

		if (config.has("neustar.polling.interval")) {
			System.setProperty("neustar.polling.interval", config.getString("neustar.polling.interval"));
		}
	}

	private void parseCertificatesConfig(final JSONObject config) {
		final String pollingInterval = "certificates.polling.interval";
		if (config.has(pollingInterval)) {
			try {
				System.setProperty(pollingInterval, config.getString(pollingInterval));
			} catch (Exception e) {
				LOGGER.warn("Failed to set system property " + pollingInterval + " from configuration object: " + e.getMessage());
			}
		}
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

	private void parseRegionalGeoConfig(final JSONObject jo) throws JSONException {
		final JSONObject config = jo.getJSONObject("config");
		final String url = config.optString("regional_geoblock.polling.url", null);

		if (url == null) {
			LOGGER.info("regional_geoblock.polling.url not configured; stopping service updater");
			getRegionalGeoUpdater().stopServiceUpdater();
			return;
		}

		if (jo.has(deliveryServicesKey)) {
			final JSONObject dss = jo.getJSONObject(deliveryServicesKey);
			for (final String ds : JSONObject.getNames(dss)) {
				if (dss.getJSONObject(ds).has("regionalGeoBlocking") &&
						dss.getJSONObject(ds).getString("regionalGeoBlocking").equals("true")) {
					final long interval = config.optLong("regional_geoblock.polling.interval");
					getRegionalGeoUpdater().setDataBaseURL(url, interval);
					return;
				}
			}
		}

		getRegionalGeoUpdater().cancelServiceUpdater();
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
				locations.add(new CacheLocation(loc, new Geolocation(jo.getDouble("latitude"), jo.getDouble("longitude"))));
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

	public void setSteeringWatcher(final SteeringWatcher steeringWatcher) {
		this.steeringWatcher = steeringWatcher;
	}

	public void setCertificatesPoller(final CertificatesPoller certificatesPoller) {
		this.certificatesPoller = certificatesPoller;
	}

	public CertificatesPublisher getCertificatesPublisher() {
		return certificatesPublisher;
	}

	public void setCertificatesPublisher(final CertificatesPublisher certificatesPublisher) {
		this.certificatesPublisher = certificatesPublisher;
	}

	public BlockingQueue<Boolean> getPublishStatusQueue() {
		return publishStatusQueue;
	}

	public void setPublishStatusQueue(final BlockingQueue<Boolean> publishStatusQueue) {
		this.publishStatusQueue = publishStatusQueue;
	}

	public void cancelProcessConfig() {
		if (isProcessing.get()) {
			cancelled.set(true);
		}
	}

	public boolean isProcessingConfig() {
		return isProcessing.get();
	}
}
