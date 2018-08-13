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
import java.net.MalformedURLException;
import java.net.URL;
import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collection;
import java.util.Date;
import java.util.HashMap;
import java.util.HashSet;
import java.util.Iterator;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.Set;
import java.util.TreeSet;
import java.util.function.Predicate;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.stream.Collectors;

import com.comcast.cdn.traffic_control.traffic_router.core.ds.SteeringWatcher;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.FederationsWatcher;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.GeolocationDatabaseUpdater;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.NetworkNode;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.NetworkUpdater;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.DeepNetworkUpdater;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.RegionalGeoUpdater;

import com.comcast.cdn.traffic_control.traffic_router.core.secure.CertificatesPoller;
import com.comcast.cdn.traffic_control.traffic_router.core.secure.CertificatesPublisher;
import com.comcast.cdn.traffic_control.traffic_router.core.util.JsonUtils;
import com.comcast.cdn.traffic_control.traffic_router.core.util.JsonUtilsException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.log4j.Logger;

import com.comcast.cdn.traffic_control.traffic_router.core.edge.Cache;
import com.comcast.cdn.traffic_control.traffic_router.core.edge.CacheLocation;
import com.comcast.cdn.traffic_control.traffic_router.core.edge.CacheRegister;
import com.comcast.cdn.traffic_control.traffic_router.core.edge.Cache.DeliveryServiceReference;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryServiceMatcher;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryServiceMatcher.Type;
import com.comcast.cdn.traffic_control.traffic_router.core.monitor.TrafficMonitorWatcher;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouterManager;
import com.comcast.cdn.traffic_control.traffic_router.core.util.TrafficOpsUtils;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker;
import com.comcast.cdn.traffic_control.traffic_router.geolocation.Geolocation;
import com.comcast.cdn.traffic_control.traffic_router.core.request.HTTPRequest;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.AnonymousIp;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.AnonymousIpConfigUpdater;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.AnonymousIpDatabaseUpdater;

@SuppressWarnings("PMD.TooManyFields")
public class ConfigHandler {
	private static final Logger LOGGER = Logger.getLogger(ConfigHandler.class);

	private static long lastSnapshotTimestamp = 0;
	private static final Object configSync = new Object();
	public final static String CONTENT_SERVERS_KEY = "contentServers";
	public final static String DELIVERY_SERVICES_KEY = "deliveryServices";
	public final static String CONFIG_KEY = "config";
	final static String DS_URL = "DS_URL";
	final static String NOT_DS_URL = "NOT_DS_URL";

	private TrafficRouterManager trafficRouterManager;
	private GeolocationDatabaseUpdater geolocationDatabaseUpdater;
	private StatTracker statTracker;
	private String configDir;
	private String trafficRouterId;
	private TrafficOpsUtils trafficOpsUtils;

	private NetworkUpdater networkUpdater;
	private DeepNetworkUpdater deepNetworkUpdater;
	private FederationsWatcher federationsWatcher;
	private RegionalGeoUpdater regionalGeoUpdater;
	private AnonymousIpConfigUpdater anonymousIpConfigUpdater;
	private AnonymousIpDatabaseUpdater anonymousIpDatabaseUpdater;
	private SteeringWatcher steeringWatcher;
	private CertificatesPoller certificatesPoller;
	private CertificatesPublisher certificatesPublisher;
	private BlockingQueue<Boolean> publishStatusQueue;
	private final AtomicBoolean cancelled = new AtomicBoolean(false);
	private final AtomicBoolean isProcessing = new AtomicBoolean(false);

	private final static String NEUSTAR_POLLING_URL = "neustar.polling.url";
	private final static String NEUSTAR_POLLING_INTERVAL = "neustar.polling.interval";
	private final static String LOCALIZATION_METHODS = "localizationMethods";

	public String getConfigDir() {
		return configDir;
	}

	public String getTrafficRouterId() {
		return trafficRouterId;
	}

	public GeolocationDatabaseUpdater getGeolocationDatabaseUpdater() {
		return geolocationDatabaseUpdater;
	}

	public NetworkUpdater getNetworkUpdater() {
		return networkUpdater;
	}

	public DeepNetworkUpdater getDeepNetworkUpdater() {
		return deepNetworkUpdater;
	}

	public RegionalGeoUpdater getRegionalGeoUpdater() {
		return regionalGeoUpdater;
	}

	public AnonymousIpConfigUpdater getAnonymousIpConfigUpdater() {
		return anonymousIpConfigUpdater;
	}

	public AnonymousIpDatabaseUpdater getAnonymousIpDatabaseUpdater() {
		return anonymousIpDatabaseUpdater;
	}

	@SuppressWarnings({"PMD.AvoidCatchingThrowable"})
	public boolean processConfig(final String snapJson, final String compJson) throws JsonUtilsException, IOException {
		isProcessing.set(true);
		LOGGER.debug("Entered processConfig");
		if (snapJson == null) {
			cancelled.set(false);
			isProcessing.set(false);
			publishStatusQueue.clear();
			LOGGER.info("Exiting processConfig: No json data to process because snapshot was NULL.");
			return false;
		}

		Date date;
		final ObjectMapper mapper = new ObjectMapper();
		final JsonNode jo = mapper.readTree(snapJson);
		final JsonNode stats = JsonUtils.getJsonNode(jo, "stats");
		final ObjectMapper compMapper = new ObjectMapper();
		JsonNode cjo = null;
		if (compJson != null) {
			cjo = compMapper.readTree(compJson);
		}
		synchronized (configSync){
			// Check to see if this is a new Snapshot
			final long sts = getSnapshotTimestamp(stats);
			date = new Date(sts * 1000L);
			if (sts <= getLastSnapshotTimestamp()){
				cancelled.set(false);
				isProcessing.set(false);
				publishStatusQueue.clear();
				LOGGER.info("Exiting processConfig: Incoming CrConfig snapshot timestamp (" + sts + ") is older than " +
						"the loaded timestamp (" + getLastSnapshotTimestamp() + "); unable to process");
				return false;
			}
			try{
				// Search for updates, adds and deletes to delivery services
				final SnapshotEventsProcessor snapshotEventsProcessor = SnapshotEventsProcessor
						.diffCrConfigs(jo, cjo);

				if (snapshotEventsProcessor.shouldReloadConfig()){
					if (initializeAll(jo, snapshotEventsProcessor)){
						ConfigHandler.setLastSnapshotTimestamp(sts);
						return true;
					}
				}else if (processChangeEvents(jo, snapshotEventsProcessor)){
					ConfigHandler.setLastSnapshotTimestamp(sts);
					return true;
				}
			} catch (ParseException e){
				LOGGER.error("Exiting processConfig: Failed to process config for snapshot from " + date, e);
				return false;
			} finally{
				isProcessing.set(false);
				cancelled.set(false);
				publishStatusQueue.clear();
			}
		}
		return false;
	}

	@SuppressWarnings({"PMD.AvoidCatchingThrowable"})
	private boolean processChangeEvents(final JsonNode crConfigJo,
			final SnapshotEventsProcessor snapshotEventsProcessor) throws ParseException, JsonUtilsException, IOException {
		LOGGER.debug("In processChangeEvents");
		final CacheRegister cacheRegister = new CacheRegister();
		if (trafficRouterManager.getTrafficRouter() != null) {
			final CacheRegister prevCr = trafficRouterManager.getTrafficRouter().getCacheRegister();
			cacheRegister.shallowCopy(prevCr);
		}
		final JsonNode config = JsonUtils.getJsonNode(crConfigJo, CONFIG_KEY);
		cacheRegister.setConfig(config);
		parseRegionalGeoConfig(config, snapshotEventsProcessor);
		parseAnonymousIpConfig(config, snapshotEventsProcessor);
		updateCertsPublisher(snapshotEventsProcessor);
		final List<DeliveryService> httpsDeliveryServices = snapshotEventsProcessor.getSSLEnabledChangeEvents();
		if (cancelledWhileWaitForCerts(httpsDeliveryServices)) {
			return false;
		}
		// updates, creates and removes the DeliveryServices in cacheRegister
		parseDeliveryServiceMatchSets(snapshotEventsProcessor, cacheRegister);
		final JsonNode contentServers = JsonUtils.getJsonNode(crConfigJo, ConfigHandler.CONTENT_SERVERS_KEY);
		parseCacheConfig(snapshotEventsProcessor, contentServers, cacheRegister);
		trafficRouterManager.setCacheRegister(cacheRegister, snapshotEventsProcessor);
		trafficRouterManager.getTrafficRouter().configurationChanged();
		NetworkNode.getInstance().clearCacheLocations();
		NetworkNode.getDeepInstance().clearCacheLocations(true);
		return true;
	}

	private boolean cancelledWhileWaitForCerts(final List<DeliveryService> httpsDeliveryServices) {
		if (httpsDeliveryServices != null && !httpsDeliveryServices.isEmpty()){
			httpsDeliveryServices.forEach(ds -> LOGGER.info("Checking for certificate for " + ds.getId()));
			try{
				publishStatusQueue.put(true);
			} catch (InterruptedException e){
				LOGGER.warn("Failed to notify certificates publisher we're waiting for certificates", e);
			}
			while (!cancelled.get() && !publishStatusQueue.isEmpty()){
				try{
					LOGGER.info("Waiting for https certificates to support new config " + String
							.format("%x", publishStatusQueue.hashCode()));
					Thread.sleep(1000L);
				} catch (Exception t){
					LOGGER.warn("Interrupted while waiting for status on publishing ssl certs", t);
				}
			}
		}
		if (cancelled.get()) {
			LOGGER.info("Exiting cancelledWhileWaitForCerts: processing of config was CANCELLED because a newer " +
					"one is ready.");
			return true;
		}
		return false;
	}

	public CertificatesPoller getCertificatesPoller() {
		return certificatesPoller;
	}

	private void updateCertsPublisher(final SnapshotEventsProcessor snapshotEventsProcessor) {
		Collection<DeliveryService> deliveryServices = null;

		if (snapshotEventsProcessor.getCreationEvents() != null && !snapshotEventsProcessor.getCreationEvents()
				.isEmpty()) {
			deliveryServices = snapshotEventsProcessor.getCreationEvents().values();
			getCertificatesPublisher().getDeliveryServices().addAll(deliveryServices);
		}

		if (snapshotEventsProcessor.getUpdateEvents() != null && !snapshotEventsProcessor.getUpdateEvents().isEmpty()) {
			getCertificatesPublisher().getDeliveryServices().replaceAll(deliveryService -> {
				return findFirst(
					snapshotEventsProcessor.getUpdateEvents().values(),
					updatedDs ->  updatedDs.getId().equals(deliveryService.getId())
				)
				.orElse(deliveryService);
			});
		}

		if (snapshotEventsProcessor.getDeleteEvents() != null && !snapshotEventsProcessor.getDeleteEvents().isEmpty()) {
			getCertificatesPublisher().getDeliveryServices().removeIf(ds -> {
				return findFirst(
					snapshotEventsProcessor.getDeleteEvents().values(),
					deletedDs -> deletedDs.getId().equals(ds.getId())
				)
				.isPresent();
			});
		}

		getCertificatesPoller().restart();
	}


	@SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.NPathComplexity", "PMD.AvoidCatchingThrowable"})
	private boolean initializeAll (final JsonNode jo, final SnapshotEventsProcessor snapshotEventsProcessor )
			throws JsonUtilsException, ParseException, IOException {
		final Map<String, DeliveryService> deliveryServiceMap = snapshotEventsProcessor.getCreationEvents();
		final JsonNode config = JsonUtils.getJsonNode(jo, CONFIG_KEY);
		final JsonNode stats = JsonUtils.getJsonNode(jo, "stats");
		final CacheRegister cacheRegister = new CacheRegister();
		parseGeolocationConfig(config);
		parseCoverageZoneNetworkConfig(config);
		parseDeepCoverageZoneNetworkConfig(config);
		parseRegionalGeoConfig(jo, snapshotEventsProcessor);
		parseAnonymousIpConfig(jo, snapshotEventsProcessor);
		cacheRegister.setTrafficRouters(JsonUtils.getJsonNode(jo, "contentRouters"));
		cacheRegister.setConfig(config);
		cacheRegister.setStats(stats);
		parseTrafficOpsConfig(config, stats);
		parseCertificatesConfig(config);
		final ArrayList<DeliveryService> deliveryServices = new ArrayList<>();

		if (deliveryServiceMap != null && !deliveryServiceMap.values().isEmpty()) {
			deliveryServices.addAll(deliveryServiceMap.values());
			getCertificatesPublisher().setDeliveryServices(deliveryServices);
		}

		getCertificatesPoller().restart();
		final List<DeliveryService> httpsDeliveryServices = deliveryServices.stream()
				.filter(ds -> !ds.isDns() && ds.isSslEnabled()).collect(Collectors.toList());

		if (cancelledWhileWaitForCerts(httpsDeliveryServices)) {
			return false;
		}

		parseDeliveryServiceMatchSets(deliveryServiceMap, cacheRegister);
		parseLocationConfig(JsonUtils.getJsonNode(jo, "edgeLocations"), cacheRegister);
		parseCacheConfig(JsonUtils.getJsonNode(jo, ConfigHandler.CONTENT_SERVERS_KEY), cacheRegister);
		parseMonitorConfig(JsonUtils.getJsonNode(jo, "monitors"));
		federationsWatcher.configure(config);
		steeringWatcher.configure(config);
		trafficRouterManager.setCacheRegister(cacheRegister);
		trafficRouterManager.getNameServer().setEcsEnable(JsonUtils.optBoolean(config, "ecsEnable", false));
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
		 *
		 * Update for DDC (Dynamic Deep Caching): NetworkNode now has a 2nd singleton (deepInstance) that is managed
		 * similarly to the non-deep instance. However, instead of clearing a NetworkNode's CacheLocation, only the
		 * Caches are cleared from the CacheLocation then lazily loaded at request time.
		 */
		NetworkNode.getInstance().clearCacheLocations();
		NetworkNode.getDeepInstance().clearCacheLocations(true);
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

	public void setDeepNetworkUpdater(final DeepNetworkUpdater dnu) {
		this.deepNetworkUpdater = dnu;
	}

	public void setRegionalGeoUpdater(final RegionalGeoUpdater regionalGeoUpdater) {
		this.regionalGeoUpdater = regionalGeoUpdater;
	}

	public void setAnonymousIpConfigUpdater(final AnonymousIpConfigUpdater anonymousIpConfigUpdater) {
		this.anonymousIpConfigUpdater = anonymousIpConfigUpdater;
	}

	public void setAnonymousIpDatabaseUpdater(final AnonymousIpDatabaseUpdater anonymousIpDatabaseUpdater) {
		this.anonymousIpDatabaseUpdater = anonymousIpDatabaseUpdater;
	}

	/**
	 * Parses the Traffic Ops config
	 *
	 * @param config the config section
	 * @param stats  the stats section
	 * @throws JsonUtilsException
	 */
	private void parseTrafficOpsConfig(final JsonNode config, final JsonNode stats) throws JsonUtilsException {
		if (stats.has("tm_host")) {
			trafficOpsUtils.setHostname(JsonUtils.getString(stats, "tm_host"));
		} else if (stats.has("to_host")) {
			trafficOpsUtils.setHostname(JsonUtils.getString(stats, "to_host"));
		} else {
			throw new JsonUtilsException("Unable to find to_host or tm_host in stats section of TrConfig; unable to build TrafficOps URLs");
		}

		trafficOpsUtils.setCdnName(JsonUtils.optString(stats, "CDN_name", null));
		trafficOpsUtils.setConfig(config);
	}

	/**
	 * Parses the cache information from the configuration and updates the {@link CacheRegister}.
	 * @param snap the {@link SnapshotEventsProcessor}
	 * @param contentServers
	 * @param cacheRegister the {@link CacheRegister}
	 */
	private void parseCacheConfig(final SnapshotEventsProcessor snap,
			final JsonNode contentServers, final CacheRegister cacheRegister) throws
			ParseException, JsonUtilsException {
		final Map<String, Cache> existingCacheMap = cacheRegister.getCacheMap();
		existingCacheMap.entrySet().removeIf(exEntry-> {
			if (snap.getDeleteCacheEvents().contains(exEntry.getKey())){
				return true;
			}
			return false;
		});
		final Set<CacheLocation> existingLocations = cacheRegister.getCacheLocations();
		existingLocations.forEach(loc->{
			loc.getCaches().forEach(cacheEntry-> {
				final String cacheId = cacheEntry.getId();
				if (snap.getDeleteCacheEvents().contains(cacheId)) {
					loc.removeCache(cacheId);
				}
			});
		});
		cacheRegister.setConfiguredLocations(existingLocations);
		// scan remaining caches for ds link changes
		for (final String cacheName : existingCacheMap.keySet()) {
			// If its in the modified list then replace the mappings in existing
			if (!snap.getMappingEvents().contains(cacheName)) {
				LOGGER.debug("No delivery service mapping changes for cache named: " + cacheName);
				continue;
			}
			LOGGER.info("Found delivery service MAPPING CHANGE FOR CACHE: " + cacheName);
			final Cache existingCache = existingCacheMap.get(cacheName);
			final JsonNode cacheJo = JsonUtils.getJsonNode(contentServers, cacheName);
			final String locationId = JsonUtils.getString(cacheJo, "locationId");
			final CacheLocation loc = cacheRegister.getCacheLocation(locationId);
			final Cache newCache = buildCache(cacheName, cacheJo, loc);
			final Collection<DeliveryServiceReference> existingDsRefs = existingCache.getDeliveryServices();
			for (final DeliveryServiceReference existingDsr : existingDsRefs) {
				final String dsId = existingDsr.getDeliveryServiceId();
				DeliveryService theDs = cacheRegister.getDeliveryService(dsId);
				if (theDs != null) {
					final DeliveryServiceReference ndsr = newCache.getDeliveryService(dsId);
					if (ndsr != null){
						addNewStats(ndsr, theDs);
					}
					continue;
				}
				// This means the delivery service has been deleted, so we can only find it in the deleted list
				theDs = snap.getDeleteEvents().get(dsId);
				if (theDs == null) {
					// this means there is a bug somewhere
					// allow NullPointerException
					throw new ParseException("parseCacheConfig: This DeliveryService has been deleted from the" +
							" Cache " +
							"Register, but it is not in the list of Delete events: " + dsId);
				}
				statTracker.removeTracks(theDs, existingDsr.getAliases());
			}
			loc.removeCache(cacheName);
			loc.addCache(newCache);
			existingCacheMap.put(cacheName, newCache);
			newCache.getDeliveryServices()
					.forEach(dsr -> LOGGER.info("cache = " + cacheName + " DSR = " + dsr.getDeliveryServiceId()+", " +
							"aliases: "+dsr.getAliases()));
		}
		cacheRegister.setCacheMap(existingCacheMap);
	}

	private void addNewStats(final DeliveryServiceReference dsr, final DeliveryService aDs) {
		List<String> newAliases = null;
		if (dsr != null) {
			newAliases = dsr.getAliases();
		}
		if (newAliases != null) {
			final List<StatTracker.Track> newTracks = StatTracker.createTracksForDs(aDs, newAliases);
			newTracks.forEach(newTrack -> statTracker.saveTrack(newTrack));
		}
	}

	/**
	 * Parses the cache information from the configuration and updates the {@link CacheRegister}.
	 *
	 * @param contentServers the JsonNode reference to the cache servers section
	 * @param cacheRegister
	 * @throws JsonUtilsException, ParseException
	 */
	@SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.AvoidDeeplyNestedIfStmts", "PMD.NPathComplexity"})
	private void parseCacheConfig(final JsonNode contentServers, final CacheRegister cacheRegister) throws
			JsonUtilsException, ParseException {
		final Map<String, List<String>> statMap = new HashMap<String, List<String>>();
		final Map<String, Cache> cacheMap = new HashMap<>();
		final Iterator<String> cacheIter = contentServers.fieldNames();
		while (cacheIter.hasNext()){
			final String cacheId = cacheIter.next();
			final JsonNode cacheJo = JsonUtils.getJsonNode(contentServers, cacheId);
			final CacheLocation loc = cacheRegister.getCacheLocation(JsonUtils.getString(cacheJo, "locationId"));
			if (loc != null){
				final Cache cache = buildCache(cacheId, cacheJo, loc);
				final Collection<DeliveryServiceReference> dsRs = cache.getDeliveryServices();
				final Iterator<DeliveryServiceReference> dsIter = dsRs.iterator();
				while (dsIter.hasNext()){
					/* technically this could be more than just a string or array,
					 * but, as we only have had those two types, let's not worry about the future
					 */
					final DeliveryServiceReference dsr = dsIter.next();
					final String ds = dsr.getDeliveryServiceId();
					List<String> dsNames = statMap.get(ds);

					if (dsNames == null){
						dsNames = new ArrayList<>();
					}

					final String tld = JsonUtils.optString(cacheRegister.getConfig(), "domain_name");
					for (final String name : dsr.getAliases()){
						if (name.endsWith(tld)){
							final String reName = name.replaceAll("^.*?\\.", "");
							if (!dsNames.contains(reName)){
								dsNames.add(reName);
							}
						}else if (!dsNames.contains(name)){
							dsNames.add(name);
						}
					}
					statMap.put(ds, dsNames);
				}

				loc.addCache(cache);
				cacheMap.put(cache.getId(), cache);
			}
		}
		cacheRegister.setCacheMap(cacheMap);
		statTracker.initialize(statMap, cacheRegister);

	}

	public Cache buildCache(final String cacheId, final JsonNode cacheJo, final CacheLocation loc)
			throws JsonUtilsException, ParseException{
		String hashId = cacheId;
		// not only must we check for the key, but also if it's null; problems with consistent hashing can arise if we use a null value as the hashId
		if (cacheJo.has("hashId") && cacheJo.get("hashId").textValue() != null){
			hashId = cacheJo.get("hashId").textValue();
		}

		Geolocation geolocation = null;
		if ( loc != null)
		{
			geolocation = loc.getGeolocation();
		}
		final Cache cache = new Cache(cacheId, hashId, JsonUtils.optInt(cacheJo, "hashCount"), geolocation);
		cache.setFqdn(JsonUtils.getString(cacheJo, "fqdn"));
		cache.setPort(JsonUtils.getInt(cacheJo, "port"));

		final String ip = JsonUtils.getString(cacheJo, "ip");
		final String ip6 = JsonUtils.optString(cacheJo, "ip6");

		try{
			cache.setIpAddress(ip, ip6, 0);
		} catch (UnknownHostException e){
			LOGGER.warn(e + " : " + ip);
		}

		if (cacheJo.has(DELIVERY_SERVICES_KEY)){
			final List<DeliveryServiceReference> references = parseDsRefs(cacheJo);
			cache.setDeliveryServices(references);
		}
		return cache;
	}

	private List<DeliveryServiceReference> parseDsRefs(final JsonNode cacheJson) throws JsonUtilsException,
			ParseException{
		final List<DeliveryServiceReference> dsRefs = new ArrayList<>();
		if (cacheJson.has(ConfigHandler.DELIVERY_SERVICES_KEY)) {
			final JsonNode dsRefsJson = JsonUtils.getJsonNode(cacheJson, ConfigHandler.DELIVERY_SERVICES_KEY);
			final Iterator<String> dsrfs = dsRefsJson.fieldNames();
			while (dsrfs.hasNext()) {
				final String dsid = dsrfs.next();
				final JsonNode dso = dsRefsJson.get(dsid);
				DeliveryServiceReference newDsref = null;
				if (dso.isArray() && dso.size() > 0) {
					final List<String> dsLinks = parseDsAliases(dso);
					newDsref = new DeliveryServiceReference(dsid, dsLinks.get(0), dsLinks);
					dsRefs.add(newDsref);
				} else {
					dsRefs.add(new DeliveryServiceReference(dsid, dso.toString()));
				}
			}
		}
		return dsRefs;
	}

	private List<String> parseDsAliases(final JsonNode dso) {
		final List<String> dsNames = new ArrayList<>();

		for (final JsonNode nameNode : dso){
			final String name = nameNode.textValue();

			if (!dsNames.contains(name)){
				dsNames.add(name);
			}
		}
		return dsNames;
	}

	@SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.NPathComplexity", "PMD.AvoidCatchingThrowable"})
	private void parseDeliveryServiceMatchSets(final SnapshotEventsProcessor snap,
			final CacheRegister cacheRegister) throws JsonUtilsException {
		final List<JsonUtilsException> jues = new ArrayList<>();
		if (cacheRegister.getDeliveryServiceMatchers() == null) {
			cacheRegister.setDeliveryServiceMatchers(new TreeSet<>());
		}
		if (cacheRegister.getDeliveryServices() == null) {
			cacheRegister.setDeliveryServiceMap(new HashMap<>());
		}
		final TreeSet<DeliveryServiceMatcher> matchers = cacheRegister.getDeliveryServiceMatchers();
		final Map<String, DeliveryService> masterDsMap = cacheRegister.getDeliveryServices();
		final Map<String, DeliveryService> deliveryServiceMap = snap.getChangeEvents();
		deliveryServiceMap.forEach((deliveryServiceId, deliveryService) -> {
			try {
				final JsonNode matchsets = deliveryService.getMatchsets();
				for (final JsonNode matchset : matchsets) {
					final DeliveryServiceMatcher deliveryServiceMatcher = new DeliveryServiceMatcher(deliveryService);
					matchers.removeIf(dsm -> dsm.getDeliveryService().getId().equals(deliveryServiceId));
					matchers.add(deliveryServiceMatcher);
					for (final JsonNode matcherJo : JsonUtils.getJsonNode(matchset, "matchlist")) {
						final Type type = Type.valueOf(JsonUtils.getString(matcherJo, "match-type"));
						final String target = JsonUtils.optString(matcherJo, "target");
						deliveryServiceMatcher.addMatch(type, JsonUtils.getString(matcherJo, "regex"), target);
					}
				}
				final String rurl = deliveryService.getGeoRedirectUrl();
				if (rurl != null) {
					final int idx = rurl.indexOf("://");
					if (idx < 0) {
						//this is a relative url, belongs to this ds
						deliveryService.setGeoRedirectUrlType(DS_URL);
					} else {
						//this is a url with protocol, must check further
						//first, parse the url, if url invalid it will throw Exception
						final URL url = new URL(rurl);
						//make a fake HTTPRequest for the redirect url
						final HTTPRequest req = new HTTPRequest(url);
						deliveryService.setGeoRedirectFile(url.getFile());
						//try select the ds by the redirect fake HTTPRequest
						final DeliveryService rds = cacheRegister.getDeliveryService(req, true);
						if (rds == null || rds.getId().equals(deliveryService.getId())) {
							//the redirect url not belongs to this ds
							deliveryService.setGeoRedirectUrlType(NOT_DS_URL);
						} else {
							deliveryService.setGeoRedirectUrlType(DS_URL);
						}
					}
				}
				masterDsMap.put(deliveryServiceId, deliveryService);
			} catch (JsonUtilsException e) {
				jues.add(e);
			} catch (MalformedURLException mue) {
				LOGGER.error("fatal error, failed to init NGB redirect with Exception: " + mue.getMessage());
			}
		});
		if (!jues.isEmpty()) {
			throw jues.get(0);
		}
		// remove delivery services that have been deleted in the snapshot
		snap.getDeleteEvents().forEach((deliveryServiceId, deliveryService) -> {
			matchers.removeIf(dsm -> dsm.getDeliveryService().getId().equals(deliveryServiceId));
			masterDsMap.remove(deliveryServiceId);
		});
		cacheRegister.setDeliveryServiceMap(masterDsMap);
		cacheRegister.setDeliveryServiceMatchers(matchers);
	}

	private void parseDeliveryServiceMatchSets(final Map<String, DeliveryService> deliveryServiceMap,
			final CacheRegister cacheRegister) throws JsonUtilsException {
		final TreeSet<DeliveryServiceMatcher> matchers = new TreeSet<>();
		final List<JsonUtilsException> jues = new ArrayList<>();

		deliveryServiceMap.forEach((deliveryServiceId, deliveryService) -> {
			try {
				final JsonNode matchsets = deliveryService.getMatchsets();

				for (final JsonNode matchset : matchsets) {
					final DeliveryServiceMatcher deliveryServiceMatcher = new DeliveryServiceMatcher(deliveryService);

					matchers.add(deliveryServiceMatcher);

					for (final JsonNode matcherJo : JsonUtils.getJsonNode(matchset, "matchlist")) {
						final Type type = Type.valueOf(JsonUtils.getString(matcherJo, "match-type"));
						final String target = JsonUtils.optString(matcherJo, "target");
						deliveryServiceMatcher.addMatch(type, JsonUtils.getString(matcherJo, "regex"), target);
					}
				}
			} catch (JsonUtilsException e) {
				jues.add(e);
			}
		});

		if (!jues.isEmpty()) {
			throw jues.get(0);
		}

		cacheRegister.setDeliveryServiceMap(deliveryServiceMap);
		cacheRegister.setDeliveryServiceMatchers(matchers);
		initGeoFailedRedirect(deliveryServiceMap, cacheRegister);
	}


	private void initGeoFailedRedirect(final Map<String, DeliveryService> dsMap, final CacheRegister cacheRegister) {
		final Iterator<String> itr = dsMap.keySet().iterator();
		while (itr.hasNext()) {
			final DeliveryService ds = dsMap.get(itr.next());
			//check if it's relative path or not
			final String rurl = ds.getGeoRedirectUrl();
			if (rurl == null) {
				continue;
			}

			try {
				final int idx = rurl.indexOf("://");

				if (idx < 0) {
					//this is a relative url, belongs to this ds
					ds.setGeoRedirectUrlType(DS_URL);
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
				if (rds == null || !rds.getId().equals(ds.getId())) {
					//the redirect url not belongs to this ds
					ds.setGeoRedirectUrlType(NOT_DS_URL);
					continue;
				}

				ds.setGeoRedirectUrlType(DS_URL);
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
	 * @throws JsonUtilsException
	 */
	private void parseGeolocationConfig(final JsonNode config) throws JsonUtilsException {
		String pollingUrlKey = "geolocation.polling.url";

		if (config.has("alt.geolocation.polling.url")) {
			pollingUrlKey = "alt.geolocation.polling.url";
		}

		getGeolocationDatabaseUpdater().setDataBaseURL(
			JsonUtils.getString(config, pollingUrlKey),
			JsonUtils.optLong(config, "geolocation.polling.interval")
		);

		if (config.has(NEUSTAR_POLLING_URL)) {
			System.setProperty(NEUSTAR_POLLING_URL, JsonUtils.getString(config, NEUSTAR_POLLING_URL));
		}

		if (config.has(NEUSTAR_POLLING_INTERVAL)) {
			System.setProperty(NEUSTAR_POLLING_INTERVAL, JsonUtils.getString(config, NEUSTAR_POLLING_INTERVAL));
		}
	}

	private void parseCertificatesConfig(final JsonNode config) {
		final String pollingInterval = "certificates.polling.interval";
		if (config.has(pollingInterval)) {
			try {
				System.setProperty(pollingInterval, JsonUtils.getString(config, pollingInterval));
			} catch (Exception e) {
				LOGGER.warn("Failed to set system property " + pollingInterval + " from configuration object: " + e
						.getMessage());
			}
		}
	}

	private void parseAnonymousIpConfig(final JsonNode config, final SnapshotEventsProcessor dsep) {
		final String anonymousPollingUrl = "anonymousip.polling.url";
		final String anonymousPollingInterval = "anonymousip.polling.interval";
		final String anonymousPolicyConfiguration = "anonymousip.policy.configuration";
		final long interval = JsonUtils.optLong(config, anonymousPollingInterval);
		final String configUrl = JsonUtils.optString(config, anonymousPolicyConfiguration, null);
		final String databaseUrl = JsonUtils.optString(config, anonymousPollingUrl, null);

		if (configUrl == null) {
			LOGGER.info(anonymousPolicyConfiguration + " policy not configured; stopping service updater and disabling feature");
			getAnonymousIpConfigUpdater().stopServiceUpdater();
			AnonymousIp.getCurrentConfig().enabled = false;
			return;
		}

		if (databaseUrl == null) {
			LOGGER.info(anonymousPollingUrl + " DB url not configured; stopping service updater and disabling feature");
			getAnonymousIpDatabaseUpdater().stopServiceUpdater();
			AnonymousIp.getCurrentConfig().enabled = false;
			return;
		}

		Collection<DeliveryService> deliveryServices = dsep.getCreationEvents().values();
		if (findFirst(deliveryServices, ds -> ds.isAnonymousIpEnabled()).isPresent()) {
			getAnonymousIpConfigUpdater().setDataBaseURL(configUrl, interval);
			getAnonymousIpDatabaseUpdater().setDataBaseURL(databaseUrl, interval);
			AnonymousIp.getCurrentConfig().enabled = true;
			LOGGER.debug("new Anonymous Blocking in use, scheduling service updaters and enabling feature");
			return;
		}

		deliveryServices = dsep.getUpdateEvents().values();
		if (findFirst(deliveryServices, ds -> ds.isAnonymousIpEnabled()).isPresent()) {
			getAnonymousIpConfigUpdater().setDataBaseURL(configUrl, interval);
			getAnonymousIpDatabaseUpdater().setDataBaseURL(databaseUrl, interval);
			AnonymousIp.getCurrentConfig().enabled = true;
			LOGGER.debug("added Anonymous Blocking in use, scheduling service updaters and enabling feature");
			return;
		}

		if (!AnonymousIp.getCurrentConfig().enabled) {
			LOGGER.debug("No DS using anonymous ip blocking - disabling feature");
			getAnonymousIpConfigUpdater().cancelServiceUpdater();
			getAnonymousIpDatabaseUpdater().cancelServiceUpdater();
		}
	}

	/**
	 * Parses the ConverageZoneNetwork database configuration and updates the database if the URL has
	 * changed.
	 *
	 * @param config
	 * @throws JsonUtilsException
	 */
	private void parseCoverageZoneNetworkConfig(final JsonNode config) throws JsonUtilsException {
		getNetworkUpdater().setDataBaseURL(
				JsonUtils.getString(config, "coveragezone.polling.url"),
				JsonUtils.optLong(config, "coveragezone.polling.interval")
		);
	}

	private void parseDeepCoverageZoneNetworkConfig(final JsonNode config) {
		getDeepNetworkUpdater().setDataBaseURL(
			JsonUtils.optString(config, "deepcoveragezone.polling.url", null),
			JsonUtils.optLong(config, "deepcoveragezone.polling.interval")
		);
	}

	private Optional<DeliveryService> findFirst(final Collection<DeliveryService> deliveryServices,
												final Predicate<DeliveryService> dsTester) {
		return deliveryServices.stream()
				.filter(ds -> dsTester.test(ds))
				.findFirst();
	}

	private void parseRegionalGeoConfig(final JsonNode config, final SnapshotEventsProcessor dsep) {
		final String url = JsonUtils.optString(config, "regional_geoblock.polling.url", null);
		final long interval = JsonUtils.optLong(config, "regional_geoblock.polling.interval");
		final RegionalGeoUpdater rgu = getRegionalGeoUpdater();

		if (url == null) {
			LOGGER.info("regional_geoblock.polling.url not configured; stopping service updater");
			if (rgu != null){
				rgu.stopServiceUpdater();
			}
			return;
		}

		Collection<DeliveryService> deliveryServices = dsep.getCreationEvents().values();
		if (findFirst(deliveryServices, ds -> ds.isRegionalGeoEnabled()).isPresent()) {
			rgu.setDataBaseURL(url, interval);
		} else {
			deliveryServices = dsep.getUpdateEvents().values();
			if (findFirst(deliveryServices, ds -> ds.isRegionalGeoEnabled()).isPresent()) {
				rgu.setDataBaseURL(url, interval);
			} else {
				rgu.cancelServiceUpdater();
			}
		}
	}

	/**
	 * Creates a {@link Map} of location IDs to {@link Geolocation}s for every Location
	 * included in the configuration that has both a latitude and a longitude specified.
	 *
	 * @param locationsJo the locations section of the config
	 * @return the {@link Map}, empty if there are no Locations that have both a latitude and
	 * longitude specified
	 * @throws JsonUtilsException
	 */
	private void parseLocationConfig(final JsonNode locationsJo, final CacheRegister cacheRegister) throws
			JsonUtilsException {
		final Set<CacheLocation> locations = new HashSet<CacheLocation>(locationsJo.size());

		final Iterator<String> locIter = locationsJo.fieldNames();
		while (locIter.hasNext()) {
			final String loc = locIter.next();
			final JsonNode jo = JsonUtils.getJsonNode(locationsJo, loc);
			List<String> backupCacheGroups = null;
			boolean useClosestOnBackupFailure = true;

			if (jo != null && jo.has("backupLocations")) {
				final JsonNode backupConfigJson = JsonUtils.getJsonNode(jo, "backupLocations");
				backupCacheGroups = new ArrayList<>();
				if (backupConfigJson.has("list")) {
					for (final JsonNode cacheGroup : JsonUtils.getJsonNode(backupConfigJson, "list")) {
						backupCacheGroups.add(cacheGroup.asText());
					}
					useClosestOnBackupFailure = JsonUtils.optBoolean(backupConfigJson, "fallbackToClosest", false);
				}

			}

			final Set<CacheLocation.LocalizationMethod> enabledLocalizationMethods = parseLocalizationMethods(loc, jo);

			try {
				locations.add(
						new CacheLocation(
								loc,
								new Geolocation(
										JsonUtils.getDouble(jo, "latitude"),
										JsonUtils.getDouble(jo, "longitude")),
								backupCacheGroups,
								useClosestOnBackupFailure,
								enabledLocalizationMethods));
			} catch (JsonUtilsException e) {
				LOGGER.warn(e, e);
			}
		}
		cacheRegister.setConfiguredLocations(locations);
	}

	private Set<CacheLocation.LocalizationMethod> parseLocalizationMethods(final String loc, final JsonNode jo) throws JsonUtilsException {
		final Set<CacheLocation.LocalizationMethod> enabledLocalizationMethods = new HashSet<>();
		if (jo != null && jo.hasNonNull(LOCALIZATION_METHODS) && JsonUtils.getJsonNode(jo, LOCALIZATION_METHODS).isArray()) {
			final JsonNode localizationMethodsJson = JsonUtils.getJsonNode(jo, LOCALIZATION_METHODS);
			for (final JsonNode methodJson : localizationMethodsJson) {
				if (methodJson.isNull() || !methodJson.isTextual()) {
					LOGGER.error("Location '" + loc + "' has a non-string localizationMethod, skipping");
					continue;
				}
				final String method = methodJson.asText();
				try {
					enabledLocalizationMethods.add(CacheLocation.LocalizationMethod.valueOf(method));
				} catch (IllegalArgumentException e) {
					LOGGER.error("Location '" + loc + "' has an unknown localizationMethod (" + method + "), skipping");
					continue;
				}
			}
		}
		// by default or if NO localization methods are explicitly enabled, enable ALL
		if (enabledLocalizationMethods.isEmpty()) {
			enabledLocalizationMethods.addAll(Arrays.asList(CacheLocation.LocalizationMethod.values()));
		}
		return enabledLocalizationMethods;
	}

	/**
	 * Creates a {@link Map} of Monitors used by {@link TrafficMonitorWatcher} to pull TrConfigs.
	 *
	 * @param monitors the monitors section of the TrafficRouter Configuration
	 * @throws JsonUtilsException, ParseException
	 */
	private void parseMonitorConfig(final JsonNode monitors) throws JsonUtilsException, ParseException {
		final List<String> monitorList = new ArrayList<String>();

		for (final JsonNode jo : monitors) {
			final String fqdn = JsonUtils.getString(jo, "fqdn");
			final int port = JsonUtils.optInt(jo, "port", 80);
			final String status = JsonUtils.getString(jo, "status");

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
	 * @param stats the stats section of the TrafficRouter Configuration
	 * @return long
	 * @throws JsonUtilsException
	 */
	private long getSnapshotTimestamp(final JsonNode stats) throws JsonUtilsException {
		return JsonUtils.getLong(stats, "date");
	}

	public StatTracker getStatTracker() {
		return statTracker;
	}

	public void setStatTracker(final StatTracker statTracker) {
		this.statTracker = statTracker;
	}

	static long getLastSnapshotTimestamp() {
		return lastSnapshotTimestamp;
	}

	static java.time.Instant getLastSnapshotInstant() {
		return java.time.Instant.ofEpochSecond(lastSnapshotTimestamp);
	}

	static void setLastSnapshotTimestamp(final long lastSnapshotTimestamp) {
		LOGGER.info("Last Snapshot Timestamp set: "+lastSnapshotTimestamp);
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
			} else {
				LOGGER.warn("Failed parsing request header from config");
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

