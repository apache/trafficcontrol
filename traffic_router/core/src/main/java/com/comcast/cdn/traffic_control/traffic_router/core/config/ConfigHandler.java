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
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.log4j.Logger;
import org.json.JSONException;

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

	private final static String NEUSTAR_POLLING_URL = "neustar.polling.url";
	private final static String NEUSTAR_POLLING_INTERVAL = "neustar.polling.interval";

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
			final ObjectMapper mapper = new ObjectMapper();
			final JsonNode jo = mapper.readTree(jsonStr);
			final JsonNode config = jo.get("config");
			final JsonNode stats = jo.get("stats");

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
				final JsonNode deliveryServicesJson = jo.get("deliveryServices");
				cacheRegister.setTrafficRouters(jo.get("contentRouters"));
				cacheRegister.setConfig(config);
				cacheRegister.setStats(stats);
				parseTrafficOpsConfig(config, stats);

				final Map<String, DeliveryService> deliveryServiceMap = parseDeliveryServiceConfig(jo.get(deliveryServicesKey));

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

					while (!cancelled.get() && !publishStatusQueue.isEmpty()) {
						try {
							LOGGER.info("Waiting for https certificates to support new config " + String.format("%x", publishStatusQueue.hashCode()));
							Thread.sleep(1000L);
						} catch (Throwable t) {
							LOGGER.warn("Interrupted while waiting for status on publishing ssl certs", t);
						}
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
				parseLocationConfig(jo.get("edgeLocations"), cacheRegister);
				parseCacheConfig(jo.get("contentServers"), cacheRegister);
				parseMonitorConfig(jo.get("monitors"));

				federationsWatcher.configure(config);
				steeringWatcher.configure(config);
				steeringWatcher.setCacheRegister(cacheRegister);
				trafficRouterManager.setCacheRegister(cacheRegister);
				trafficRouterManager.getTrafficRouter().setRequestHeaders(parseRequestHeaders(config.get("requestHeaders")));
				trafficRouterManager.getTrafficRouter().configurationChanged();

				/*
				 * NetworkNode uses lazy loading to associate CacheLocations with NetworkNodes at request time in TrafficRouter.
				 * Therefore this must be done last, as any thread that holds a reference to the CacheRegister might contain a reference
				 * to a Cache that no longer exists. In that case, the old CacheLocation and List<Cache> will be set on a
				 * given CacheLocation within a NetworkNode, leading to an OFFLINE cache to be served, or an ONLINE cache to
				 * never have traffic routed to it, as the old List<Cache> does not contain the Cache that was moved to ONLINE.
				 * NetworkNode is a singleton and is managed asynchronously. As long as we swap out the CacheRegister first,
				 * then clear cache locations, the lazy loading should work as designed. See issue TC-401 for details.
				 */
				NetworkNode.getInstance().clearCacheLocations();
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
	private void parseTrafficOpsConfig(final JsonNode config, final JsonNode stats) throws IOException {
		if (stats.has("tm_host")) {
			trafficOpsUtils.setHostname(stats.get("tm_host").asText());
		} else if (stats.has("to_host")) {
			trafficOpsUtils.setHostname(stats.get("to_host").asText());
		} else {
			throw new IOException("Unable to find to_host or tm_host in stats section of TrConfig; unable to build TrafficOps URLs");
		}

		trafficOpsUtils.setCdnName(stats.has("CDN_name") ? stats.get("CDN_name").textValue() : null);
		trafficOpsUtils.setConfig(config);
	}

	/**
	 * Parses the cache information from the configuration and updates the {@link CacheRegister}.
	 *
	 * @param trConfig
	 *            the {@link TrafficRouterConfiguration}
	 * @throws JSONException 
	 */
	@SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.AvoidDeeplyNestedIfStmts", "PMD.NPathComplexity"})
	private void parseCacheConfig(final JsonNode contentServers, final CacheRegister cacheRegister) throws ParseException {
		final Map<String,Cache> map = new HashMap<String,Cache>();
		final Map<String, List<String>> statMap = new HashMap<String, List<String>>();


		final Iterator<String> nodeIter = contentServers.fieldNames();
		while (nodeIter.hasNext()) {
			final String node = nodeIter.next();
			final JsonNode jo = contentServers.get(node);
			final CacheLocation loc = cacheRegister.getCacheLocation(jo.get("locationId").asText());

			if (loc != null) {
				String hashId = node;
				// not only must we check for the key, but also if it's null; problems with consistent hashing can arise if we use a null value as the hashId
				if (jo.has("hashId") && jo.get("hashId").textValue() != null) {
					hashId = jo.get("hashId").textValue();
				}

				final Cache cache = new Cache(node, hashId, jo.has("hashCount") ? jo.get("hashCount").asInt() : 0);
				cache.setFqdn(jo.get("fqdn").asText());
				cache.setPort(jo.get("port").asInt());

				final String ip = jo.get("ip").asText();
				final String ip6 = jo.has("ip6") ? jo.get("ip6").asText() : "";

				try {
					cache.setIpAddress(ip, ip6, 0);
				} catch (UnknownHostException e) {
					LOGGER.warn(e + " : " + ip);
				}

				if (jo.has(deliveryServicesKey)) {
					final List<DeliveryServiceReference> references = new ArrayList<Cache.DeliveryServiceReference>();
					final JsonNode dsJos = jo.get(deliveryServicesKey);

					final Iterator<String> dsIter = dsJos.fieldNames();
					while (dsIter.hasNext()) {
						final String ds = dsIter.next();
						final JsonNode dso = dsJos.get(ds);

						List<String> dsNames = statMap.get(ds);

						if (dsNames == null) {
							dsNames = new ArrayList<String>();
						}

						if (dso.isArray()) {
							if (dso != null && dso.size() > 0) {
								int i = 0;
								for (final JsonNode nameNode : dso) {
									final String name = nameNode.asText().toLowerCase();
									if (i == 0) {
										references.add(new DeliveryServiceReference(ds, name));
									}

									final String tld = cacheRegister.getConfig().has("domain_name") ? cacheRegister.getConfig().get("domain_name").asText().toLowerCase() : "";

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

									i++;
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

	private Map<String, DeliveryService> parseDeliveryServiceConfig(final JsonNode allDeliveryServices) {
		final Map<String,DeliveryService> deliveryServiceMap = new HashMap<>();

		final Iterator<String> deliveryServiceIter = allDeliveryServices.fieldNames();
		while (deliveryServiceIter.hasNext()) {
			final String deliveryServiceId = deliveryServiceIter.next();
			final JsonNode deliveryServiceJson = allDeliveryServices.get(deliveryServiceId);
			final DeliveryService deliveryService = new DeliveryService(deliveryServiceId, deliveryServiceJson);
			boolean isDns = false;

			final JsonNode matchsets = deliveryServiceJson.get("matchsets");

			for (final JsonNode matchset : matchsets) {
				final String protocol = matchset.has("protocol") ? matchset.get("protocol").asText(): null;
				if ("DNS".equals(protocol)) {
					isDns = true;
				}
			}

			deliveryService.setDns(isDns);
			deliveryServiceMap.put(deliveryServiceId, deliveryService);
		}

		return deliveryServiceMap;
	}

	private void parseDeliveryServiceMatchSets(final JsonNode allDeliveryServices, final Map<String, DeliveryService> deliveryServiceMap, final CacheRegister cacheRegister) {
		final TreeSet<DeliveryServiceMatcher> dnsServiceMatchers = new TreeSet<>();
		final TreeSet<DeliveryServiceMatcher> httpServiceMatchers = new TreeSet<>();

		final Iterator<String> deliveryServiceIds = allDeliveryServices.fieldNames();
		while (deliveryServiceIds.hasNext()) {
			final String deliveryServiceId = deliveryServiceIds.next();
			final JsonNode deliveryServiceJson = allDeliveryServices.get(deliveryServiceId);
			final JsonNode matchsets = deliveryServiceJson.get("matchsets");
			final DeliveryService deliveryService = deliveryServiceMap.get(deliveryServiceId);

			for (final JsonNode matchset : matchsets) {
				final String protocol = matchset.get("protocol").asText();

				final DeliveryServiceMatcher deliveryServiceMatcher = new DeliveryServiceMatcher(deliveryService);

				if ("HTTP".equals(protocol)) {
					httpServiceMatchers.add(deliveryServiceMatcher);
				} else if ("DNS".equals(protocol)) {
					dnsServiceMatchers.add(deliveryServiceMatcher);
				}

				for (final JsonNode matcherJo : matchset.get("matchlist")) {
					final Type type = Type.valueOf(matcherJo.get("match-type").asText());
					final String target = matcherJo.has("target") ? matcherJo.get("target").asText() : "";
					deliveryServiceMatcher.addMatch(type, matcherJo.get("regex").asText(), target);
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
	@SuppressWarnings("PMD.NPathComplexity")
	private void parseGeolocationConfig(final JsonNode config) {
		String pollingUrlKey = "geolocation.polling.url";

		if (config.has("alt.geolocation.polling.url")) {
			pollingUrlKey = "alt.geolocation.polling.url";
		}

		getGeolocationDatabaseUpdater().setDataBaseURL(
			config.has(pollingUrlKey) ? config.get(pollingUrlKey).textValue() : null,
			config.has("geolocation.polling.interval") ? config.get("geolocation.polling.interval").asLong() : 0
		);

		if (config.has(NEUSTAR_POLLING_URL)) {
			System.setProperty(NEUSTAR_POLLING_URL, config.has(NEUSTAR_POLLING_URL) ? config.get(NEUSTAR_POLLING_URL).textValue() : "");
		}

		if (config.has(NEUSTAR_POLLING_INTERVAL)) {
			System.setProperty(NEUSTAR_POLLING_INTERVAL, config.has(NEUSTAR_POLLING_INTERVAL) ? config.get(NEUSTAR_POLLING_INTERVAL).textValue() : "");
		}
	}

	private void parseCertificatesConfig(final JsonNode config) {
		final String pollingInterval = "certificates.polling.interval";
		if (config.has(pollingInterval)) {
			try {
				System.setProperty(pollingInterval, config.has(pollingInterval) ? config.get(pollingInterval).textValue() : "");
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
	private void parseCoverageZoneNetworkConfig(final JsonNode config) {
		getNetworkUpdater().setDataBaseURL(
				config.has("coveragezone.polling.url") ? config.get("coveragezone.polling.url").textValue() : null,
				config.has("coveragezone.polling.interval") ? config.get("coveragezone.polling.interval").asLong() : 5
			);
	}

	private void parseRegionalGeoConfig(final JsonNode jo) {
		final JsonNode config = jo.get("config");
		final String url = config.has("regional_geoblock.polling.url") ? config.get("regional_geoblock.polling.url").asText(null) : null;

		if (url == null) {
			LOGGER.info("regional_geoblock.polling.url not configured; stopping service updater");
			getRegionalGeoUpdater().stopServiceUpdater();
			return;
		}

		if (jo.has(deliveryServicesKey)) {
			final JsonNode dss = jo.get(deliveryServicesKey);
			for(final JsonNode ds : dss) {
				if (ds.has("regionalGeoBlocking") &&
						ds.get("regionalGeoBlocking").asText().equals("true")) {
					final long interval = config.has("regional_geoblock.polling.interval") ? config.get("regional_geoblock.polling.interval").asLong() : 0;
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
	private void parseLocationConfig(final JsonNode locationsJo, final CacheRegister cacheRegister) {
		final Set<CacheLocation> locations = new HashSet<CacheLocation>(locationsJo.size());

		final Iterator<String> locIter = locationsJo.fieldNames();
		while (locIter.hasNext()) {
			final String loc = locIter.next();
			final JsonNode jo = locationsJo.get(loc);
			final double latitude = jo.has("latitude") ? jo.get("latitude").asDouble() : 0.0;
			final double longitude = jo.has("longitude") ? jo.get("longitude").asDouble() : 0.0;
			locations.add(new CacheLocation(loc, new Geolocation(latitude, longitude)));
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
	private void parseMonitorConfig(final JsonNode monitors) throws ParseException {
		final List<String> monitorList = new ArrayList<String>();

		for (final JsonNode jo : monitors) {
			final String fqdn = jo.get("fqdn").asText();
			final int port = jo.has("port") ? jo.get("port").asInt(80) : 80;
			final String status = jo.get("status").asText();

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
	private long getSnapshotTimestamp(final JsonNode stats) {
		return stats.has("date") ? stats.get("date").longValue() : 0;
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

	private Set<String> parseRequestHeaders(final JsonNode requestHeaders) {
		final Set<String> headers = new HashSet<String>();

		if (requestHeaders == null) {
			return headers;
		}

		for (final JsonNode header : requestHeaders) {
			if (header != null) {
				headers.add(header.asText());
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
