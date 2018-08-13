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

package com.comcast.cdn.traffic_control.traffic_router.core.dns;

import java.io.File;
import java.io.FileWriter;
import java.io.IOException;
import java.net.Inet6Address;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.security.GeneralSecurityException;
import java.security.NoSuchAlgorithmException;
import java.time.Duration;
import java.time.Instant;
import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.HashSet;
import java.util.Iterator;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.Callable;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;
import java.util.concurrent.LinkedBlockingQueue;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;
import java.util.stream.Collectors;

import com.comcast.cdn.traffic_control.traffic_router.core.config.SnapshotEventsProcessor;
import com.comcast.cdn.traffic_control.traffic_router.core.util.JsonUtils;
import com.comcast.cdn.traffic_control.traffic_router.core.util.JsonUtilsException;
import com.fasterxml.jackson.databind.JsonNode;
import org.apache.commons.io.IOUtils;
import org.apache.log4j.Logger;
import org.xbill.DNS.AAAARecord;
import org.xbill.DNS.ARecord;
import org.xbill.DNS.CNAMERecord;
import org.xbill.DNS.DClass;
import org.xbill.DNS.NSRecord;
import org.xbill.DNS.Name;
import org.xbill.DNS.RRset;
import org.xbill.DNS.Record;
import org.xbill.DNS.SOARecord;
import org.xbill.DNS.SetResponse;
import org.xbill.DNS.TextParseException;
import org.xbill.DNS.TXTRecord;
import org.xbill.DNS.Type;
import org.xbill.DNS.Zone;

import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.edge.Cache;
import com.comcast.cdn.traffic_control.traffic_router.core.edge.Cache.DeliveryServiceReference;
import com.comcast.cdn.traffic_control.traffic_router.core.edge.CacheLocation;
import com.comcast.cdn.traffic_control.traffic_router.core.edge.CacheRegister;
import com.comcast.cdn.traffic_control.traffic_router.core.edge.InetRecord;
import com.comcast.cdn.traffic_control.traffic_router.core.edge.Resolver;
import com.comcast.cdn.traffic_control.traffic_router.core.request.DNSRequest;
import com.comcast.cdn.traffic_control.traffic_router.core.router.DNSRouteResult;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouter;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouterManager;
import com.comcast.cdn.traffic_control.traffic_router.core.util.TrafficOpsUtils;
import com.google.common.cache.CacheBuilder;
import com.google.common.cache.CacheBuilderSpec;
import com.google.common.cache.CacheLoader;
import com.google.common.cache.CacheStats;
import com.google.common.cache.LoadingCache;
import com.google.common.cache.RemovalListener;
import com.google.common.cache.RemovalNotification;
import com.google.common.util.concurrent.ListenableFuture;
import com.google.common.util.concurrent.ListenableFutureTask;

public class ZoneManager extends Resolver {
	private static final Logger LOGGER = Logger.getLogger(ZoneManager.class);

	private final TrafficRouter trafficRouter;
	private static LoadingCache<ZoneKey, Zone> dynamicZoneCache = null;
	private static LoadingCache<ZoneKey, Zone> zoneCache = null;
	private static ScheduledExecutorService zoneMaintenanceExecutor = null;
	private static ExecutorService zoneExecutor = null;
	private static final int DEFAULT_PRIMER_LIMIT = 500;
	private final StatTracker statTracker;
	private static final String IP = "ip";
	private static final String IP6 = "ip6";

	private static File zoneDirectory;
	private static SignatureManager signatureManager;

	private static Name topLevelDomain;
	private static Set<String> dnsRoutingNames;
	private static final String AAAA = "AAAA";

	protected static enum ZoneCacheType {
		DYNAMIC, STATIC
	}

	public ZoneManager(final TrafficRouter tr, final StatTracker statTracker, final TrafficOpsUtils trafficOpsUtils, final TrafficRouterManager trafficRouterManager) throws IOException {
		this.trafficRouter = tr;
		this.statTracker = statTracker;
		initTopLevelDomain(tr.getCacheRegister());
		initDnsRoutingNames(tr.getCacheRegister());
		initSignatureManager(tr.getCacheRegister(), trafficOpsUtils, trafficRouterManager);
		initZoneCache(tr);
	}

	ZoneManager( final TrafficRouter trafficRouter,  final StatTracker statTracker) throws IOException {
		this.trafficRouter = trafficRouter;
		this.statTracker = statTracker;
		initTopLevelDomain(trafficRouter.getCacheRegister());
		initDnsRoutingNames(trafficRouter.getCacheRegister());
	}

	public static void destroy() {
		zoneMaintenanceExecutor.shutdownNow();
		zoneExecutor.shutdownNow();
		getSignatureManager().destroy();
	}

	/**
	 * Use this factory method when Traffic Router is being initialized or in cases when the zone
	 * caches need to be completely rebuilt.
	 * @param tr The TrafficRouter instance containing a prepared CacheRegister
	 * @param statTracker A StatsTracker instance for tracking statistics
	 * @param trafficOpsUtils
	 * @param trafficRouterManager The TrafficRouterManager which will eventually be managing the
	 *                             TrafficRouter instance 'tr'.
	 * @return a new instance of a ZoneManager with zone caches fully populated
	 * @throws IOException
	 */
	public static ZoneManager initialInstance(final TrafficRouter tr, final StatTracker statTracker,
			final TrafficOpsUtils trafficOpsUtils, final TrafficRouterManager trafficRouterManager)
			throws IOException {
		return new ZoneManager( tr, statTracker, trafficOpsUtils, trafficRouterManager);
	}

	/**
	 * Use this factory only to make changes to the zone caches and retrieve an associated ZoneManager. The
	 * initialInstance method should have already been called at least once previously.
	 * @param tr The TrafficRouter instance containing a CacheRegister prepared with the changes cooresponding
	 *           to the changes that will be made in the zone caches.
	 * @param statTracker A StatsTracker instance for tracking statistics
	 * @param sep The SnapshotEventsProcessor containing the events representing the changes to be
	 *            made to the zone caches
	 * @return
	 * @throws IOException
	 */
	public static ZoneManager snapshotInstance( final TrafficRouter tr, final StatTracker statTracker,
		final SnapshotEventsProcessor sep) throws IOException {
		final ZoneManager zmCopy = new ZoneManager( tr, statTracker);
		if (ZoneManager.dynamicZoneCache == null || ZoneManager.zoneCache == null)
		{
			LOGGER.error("Called snapshotInstance without having ever called initialInstance. Attempting to correct.");
			return null;
		}
		zmCopy.processChangeEvents(sep);
		return zmCopy;
	}

	public void processChangeEvents(final SnapshotEventsProcessor snapshotEvents) {
		getSignatureManager().refreshKeyMap();
		updateZoneCache(snapshotEvents.getChangeEvents());
	}

	protected void updateZoneCache(final List<String> zoneKeys) {
		final CacheRegister cacheRegister = getTrafficRouter().getCacheRegister();
		final Map<String, DeliveryService> deliveryServices = cacheRegister.getDeliveryServices();
		final Map<String, DeliveryService> targetDServices = new HashMap<>();
		deliveryServices.forEach((dsid, ds) -> {
			if (zoneKeys.contains(ds.getDomain()+'.')) {
				targetDServices.putIfAbsent(dsid, ds);
			}
		});

		updateZoneCache(targetDServices);
	}

	private static void initDnsRoutingNames(final CacheRegister cacheRegister) {
		final Set<String> dnsRoutingNames = new HashSet<>();
		for (final DeliveryService ds : cacheRegister.getDeliveryServices().values()) {
			if (ds.isDns()) {
				dnsRoutingNames.add(ds.getRoutingName());
			}
		}
		setDnsRoutingNames(dnsRoutingNames);
	}

	@SuppressWarnings("PMD.UseStringBufferForStringAppends")
	private static void initTopLevelDomain(final CacheRegister data) throws TextParseException {
		String tld = JsonUtils.optString(data.getConfig(), "domain_name");

		if (!tld.endsWith(".")) {
			tld = tld + ".";
		}

		setTopLevelDomain(new Name(tld));
	}

	private void initSignatureManager(final CacheRegister cacheRegister, final TrafficOpsUtils trafficOpsUtils, final TrafficRouterManager trafficRouterManager) {
		final SignatureManager sm = new SignatureManager(this, cacheRegister, trafficOpsUtils, trafficRouterManager);
		ZoneManager.signatureManager = sm;
	}

	private static long upznCnt = 0;
	@SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.NPathComplexity"})
	protected void updateZoneCache(final Map<String, DeliveryService> changeEvents) {
		synchronized (ZoneManager.class) {
			upznCnt++;
			final TrafficRouter trafficRouter = getTrafficRouter();
			final CacheRegister cacheRegister = trafficRouter.getCacheRegister();
			final JsonNode config = cacheRegister.getConfig();
			final String dspec = "expireAfterAccess=" + (JsonUtils
					.optString(config, "zonemanager.dynamic.response.expiration", "300s")); // default to 5 minutes
			final LoadingCache<ZoneKey, Zone> dzc = createZoneCache(ZoneCacheType.DYNAMIC, CacheBuilderSpec
					.parse(dspec));
			final LoadingCache<ZoneKey, Zone> zc = createZoneCache(ZoneCacheType.STATIC);
			dzc.putAll(ZoneManager.dynamicZoneCache.asMap());
			zc.putAll(ZoneManager.zoneCache.asMap());
			final ExecutorService initExecutor = Executors.newFixedThreadPool(calcThreadPoolSize(config));
			final int initTimeout = JsonUtils.optInt(config, "zonemanager.init.timeout", 10);
			final List<Runnable> generationTasks = new ArrayList<>();
			final BlockingQueue<Runnable> primingTasks = new LinkedBlockingQueue<>();

			try {
				if (LOGGER.isDebugEnabled()){
					LOGGER.debug("Refreshing zone data for these CHANGE EVENTS: " + changeEvents.keySet().toString());
				}
				final Map<String, List<Record>> zoneMap = new HashMap<String, List<Record>>();
				final Map<String, DeliveryService> dsMap = new HashMap<String, DeliveryService>();

				for (final DeliveryService ds : changeEvents.values()) {
					final String domain = ds.getDomain();

					if (domain == null) {
						continue;
					}

					dsMap.put(domain, ds);
				}
				loadZoneCacheFromDsMap(zoneMap, dsMap, cacheRegister, trafficRouter, zc, dzc, generationTasks, primingTasks);
				initExecutor.invokeAll(generationTasks.stream().map(Executors::callable).collect(Collectors.toList()));
				LOGGER.info("Refreshed the Zone cache for the "+upznCnt+" time since the last restart.") ;
				final Instant primingStart = Instant.now();
				final List<Future<Object>> futures = initExecutor.invokeAll(primingTasks.stream().map(Executors::callable).collect(Collectors.toList()), initTimeout, TimeUnit.MINUTES);
				final Instant primingEnd = Instant.now();
				if (futures.stream().anyMatch(Future::isCancelled)) {
					LOGGER.warn(String.format("Priming zone cache exceeded time limit of %d minute(s); continuing", initTimeout));
				} else {
					LOGGER.info(String.format("Priming zone cache completed in %s", Duration.between(primingStart, primingEnd).toString()));
				}

			} catch (final InterruptedException ex) {
				LOGGER.warn("Initialization of zone data exceeded time limit of 5 minutes; continuing", ex);
			} catch (IOException ex) {
				LOGGER.fatal("Caught fatal exception while generating zone data!", ex);
			}

			ZoneManager.dynamicZoneCache = dzc;
			ZoneManager.zoneCache = zc;
		}
	}

	public static SignatureManager getSignatureManager() {
		return signatureManager;
	}

	private static int calcThreadPoolSize(final JsonNode config) {
		int poolSize = 1;
		final double scale = JsonUtils.optDouble(config, "zonemanager.threadpool.scale", 0.75);
		final int cores = Runtime.getRuntime().availableProcessors();

		if (cores > 2) {
			final Double s = Math.floor((double) cores * scale);

			if (s.intValue() > 1) {
				poolSize = s.intValue();
			}
		}

		return poolSize;
	}

	private static long lznCnt = 0l;
	protected synchronized static void initZoneCache(final TrafficRouter tr){
		final CacheRegister cacheRegister = tr.getCacheRegister();
		final JsonNode config = cacheRegister.getConfig();

		final int poolSize = calcThreadPoolSize(config);
		final ExecutorService initExecutor = Executors.newFixedThreadPool(poolSize);
		final List<Runnable> generationTasks = new ArrayList<>();
		final BlockingQueue<Runnable> primingTasks = new LinkedBlockingQueue<>();

		final ExecutorService ze = Executors.newFixedThreadPool(poolSize);
		final ScheduledExecutorService me = Executors
				.newScheduledThreadPool(2); // 2 threads, one for static, one for dynamic, threads to refresh zones
		final int maintenanceInterval = JsonUtils
				.optInt(config, "zonemanager.cache.maintenance.interval", 300); // default 5 minutes
		final int initTimeout = JsonUtils.optInt(config, "zonemanager.init.timeout", 10);
		final LoadingCache<ZoneKey, Zone> dzc = createZoneCache(ZoneCacheType.DYNAMIC, getDynamicZoneCacheSpec(config, poolSize));
		final LoadingCache<ZoneKey, Zone> zc = createZoneCache(ZoneCacheType.STATIC);

		initZoneDirectory();

		try{
			LOGGER.info("Generating zone data");
			loadAllZones(tr, zc, dzc, generationTasks, primingTasks);
			lznCnt++;
			initExecutor.invokeAll(generationTasks.stream().map(Executors::callable).collect(Collectors.toList()));
			LOGGER.info("Zone generation complete - Reloaded all Zones "+lznCnt+" times since last restart.");
			final Instant primingStart = Instant.now();
			final List<Future<Object>> futures = initExecutor.invokeAll(primingTasks.stream().map(Executors::callable).collect(Collectors.toList()), initTimeout, TimeUnit.MINUTES);
			final Instant primingEnd = Instant.now();
			if (futures.stream().anyMatch(Future::isCancelled)) {
				LOGGER.warn(String.format("Priming zone cache exceeded time limit of %d minute(s); continuing", initTimeout));
			} else {
				LOGGER.info(String.format("Priming zone cache completed in %s", Duration.between(primingStart, primingEnd).toString()));
			}

			me.scheduleWithFixedDelay(getMaintenanceRunnable(dzc, ZoneCacheType.DYNAMIC, maintenanceInterval), 0, maintenanceInterval, TimeUnit.SECONDS);
			me.scheduleWithFixedDelay(getMaintenanceRunnable(zc, ZoneCacheType.STATIC, maintenanceInterval), 0, maintenanceInterval, TimeUnit.SECONDS);

			final ExecutorService tze = ZoneManager.zoneExecutor;
			final ScheduledExecutorService tme = ZoneManager.zoneMaintenanceExecutor;
			final LoadingCache<ZoneKey, Zone> tzc = ZoneManager.zoneCache;
			final LoadingCache<ZoneKey, Zone> tdzc = ZoneManager.dynamicZoneCache;

			ZoneManager.zoneExecutor = ze;
			ZoneManager.zoneMaintenanceExecutor = me;
			ZoneManager.dynamicZoneCache = dzc;
			ZoneManager.zoneCache = zc;

			if (tze != null) {
				tze.shutdownNow();
			}

			if (tme != null) {
				tme.shutdownNow();
			}

			if (tzc != null) {
				tzc.invalidateAll();
			}

			if (tdzc != null) {
				tdzc.invalidateAll();
			}
			LOGGER.info("Initialization of zone data completed");
		} catch (final InterruptedException ex) {
			LOGGER.warn(String.format("Initialization of zone data was interrupted, timeout of %d minute(s); continuing", initTimeout), ex);
		} catch (IOException ex){
			LOGGER.fatal("Caught fatal exception while generating zone data!", ex);
		}

	}

	private static CacheBuilderSpec getDynamicZoneCacheSpec(final JsonNode config, final int poolSize) {
		final List<String> cacheSpec = new ArrayList<>();
		cacheSpec.add("expireAfterAccess=" + JsonUtils.optString(config, "zonemanager.dynamic.response.expiration", "3600s")); // default to one hour
		cacheSpec.add("concurrencyLevel=" + JsonUtils.optString(config, "zonemanager.dynamic.concurrencylevel", String.valueOf(poolSize))); // default to pool size, 4 is the actual default
		cacheSpec.add("initialCapacity=" + JsonUtils.optInt(config, "zonemanager.dynamic.initialcapacity", 10000)); // set the initial capacity to avoid expensive resizing

		return CacheBuilderSpec.parse(cacheSpec.stream().collect(Collectors.joining(",")));
	}

	private static Runnable getMaintenanceRunnable(final LoadingCache<ZoneKey, Zone> cache, final ZoneCacheType type, final int refreshInterval) {
		return new Runnable() {
			public void run() {
				synchronized (cache){
					cache.cleanUp();

					for (final ZoneKey zoneKey : cache.asMap().keySet()){
						try{
							if (getSignatureManager().needsRefresh(type, zoneKey, refreshInterval)){
								cache.refresh(zoneKey);
							}
						} catch (RuntimeException ex){
							LOGGER.fatal("RuntimeException caught on " + zoneKey.getClass()
									.getSimpleName() + " for " + zoneKey.getName(), ex);
						}
					}
				}
			}
		};
	}

	private static void initZoneDirectory() {
		synchronized(LOGGER) {
			if (zoneDirectory.exists()) {
				for (final String entry : zoneDirectory.list()) {
					final File zone = new File(zoneDirectory.getPath(), entry);
					zone.delete();
				}

				final boolean deleted = zoneDirectory.delete();

				if (!deleted) {
					LOGGER.warn("Unable to delete " + zoneDirectory);
				}
			}

			zoneDirectory.mkdir();
		}
	}

	private static void writeZone(final Zone zone) throws IOException {
		synchronized(LOGGER) {
			if (!zoneDirectory.exists() && !zoneDirectory.mkdirs()) {
				LOGGER.error(zoneDirectory.getAbsolutePath() + " directory does not exist and cannot be created!");
			}

			final File zoneFile = new File(getZoneDirectory(), zone.getOrigin().toString());
			final FileWriter w = new FileWriter(zoneFile);
			LOGGER.info("writing: " + zoneFile.getAbsolutePath());
			IOUtils.write(zone.toMasterFile(), w);
			w.flush();
			w.close();
		}
	}

	public StatTracker getStatTracker() {
		return statTracker;
	}

	static LoadingCache<ZoneKey, Zone> createZoneCache(final ZoneCacheType cacheType) {
		return createZoneCache(cacheType, CacheBuilderSpec.parse(""));
	}

	private static LoadingCache<ZoneKey, Zone> createZoneCache(final ZoneCacheType cacheType, final CacheBuilderSpec spec) {
		final RemovalListener<ZoneKey, Zone> removalListener = new RemovalListener<ZoneKey, Zone>() {
			public void onRemoval(final RemovalNotification<ZoneKey, Zone> removal){
				if (LOGGER.isDebugEnabled()){
					LOGGER.debug(cacheType + " " + removal.getKey().getClass().getSimpleName() + " " + removal.getKey()
							.getName() + " evicted from cache: " + removal.getCause());
				}
			}
		};

		return CacheBuilder.from(spec).recordStats().removalListener(removalListener).build(
				new CacheLoader<ZoneKey, Zone>() {
					final boolean writeZone = (cacheType == ZoneCacheType.STATIC) ? true : false;

					public Zone load(final ZoneKey zoneKey) throws IOException, GeneralSecurityException {
						if (LOGGER.isDebugEnabled()){
							LOGGER.debug("loading " + cacheType + " " + zoneKey.getClass().getSimpleName() +
									" " + zoneKey.getName());
						}
						return loadZone(zoneKey, writeZone);
					}

				public ListenableFuture<Zone> reload(final ZoneKey zoneKey, final Zone prevZone) throws IOException, GeneralSecurityException {
					final ListenableFutureTask<Zone> zoneTask = ListenableFutureTask.create(new Callable<Zone>() {
						public Zone call() throws IOException, GeneralSecurityException {
							return loadZone(zoneKey, writeZone);
						}
					});

						zoneExecutor.execute(zoneTask);

						return zoneTask;
					}
				}
		);
	}

	public static Zone loadZone(final ZoneKey zoneKey, final boolean writeZone) throws IOException,
			GeneralSecurityException {
		final Name name = zoneKey.getName();
		List<Record> records = zoneKey.getRecords();
		zoneKey.updateTimestamp();

		if (zoneKey instanceof SignedZoneKey) {
			records = getSignatureManager().signZone(name, records, (SignedZoneKey) zoneKey);
		}

		final Zone zone = new Zone(name, records.toArray(new Record[records.size()]));

		if (writeZone) {
			writeZone(zone);
		}

		return zone;
	}

	private static void loadAllZones(final TrafficRouter tr, final LoadingCache<ZoneKey, Zone> zc, final LoadingCache<ZoneKey, Zone> dzc, final List<Runnable> generationTasks, final BlockingQueue<Runnable> primingTasks) throws IOException {
		final CacheRegister data = tr.getCacheRegister();
		final Map<String, List<Record>> zoneMap = new HashMap<String, List<Record>>();
		final Map<String, DeliveryService> dsMap = new HashMap<String, DeliveryService>();
		final String tld = getTopLevelDomain().toString(true); // Name.toString(true) - omit the trailing dot

		for (final DeliveryService ds : data.getDeliveryServices().values()) {
			final String domain = ds.getDomain();

			if (domain == null) {
				continue;
			}

			dsMap.put(domain, ds);
		}

		loadZoneCacheFromDsMap(zoneMap, dsMap, data, tr, zc, dzc, generationTasks, primingTasks);
	}

	private static void loadZoneCacheFromDsMap(final Map<String, List<Record>> zoneMap, final Map<String, DeliveryService> dsMap,
			final CacheRegister data, final TrafficRouter tr, final LoadingCache<ZoneKey, Zone> zc, final LoadingCache<ZoneKey, Zone> dzc,
			final List<Runnable> generationTasks, BlockingQueue<Runnable> primingTasks) throws java.io.IOException {
		final Map<String, List<Record>> superDomains = populateZoneMap(zoneMap, dsMap, data);
		final List<Record> superRecords = fillZones(zoneMap, dsMap, tr, zc, dzc, generationTasks, primingTasks);
		final List<Record> upstreamRecords = fillZones(superDomains, dsMap, tr, superRecords, zc, dzc, generationTasks, primingTasks);

		for (final Record record : upstreamRecords) {
			if (record.getType() == Type.DS) {
				LOGGER.info("Publish this DS record in the parent zone: " + record);
			}
		}
	}

	private static List<Record> fillZones(final Map<String, List<Record>> zoneMap, final Map<String, DeliveryService> dsMap, final TrafficRouter tr, final LoadingCache<ZoneKey, Zone> zc, final LoadingCache<ZoneKey, Zone> dzc, final List<Runnable> generationTasks, final BlockingQueue<Runnable> primingTasks)
			throws IOException {
		return fillZones(zoneMap, dsMap, tr, null, zc, dzc, generationTasks, primingTasks);
	}

	private static List<Record> fillZones(final Map<String, List<Record>> zoneMap, final Map<String, DeliveryService> dsMap, final TrafficRouter tr, final List<Record> superRecords, final LoadingCache<ZoneKey, Zone> zc, final LoadingCache<ZoneKey, Zone> dzc, final List<Runnable> generationTasks, final BlockingQueue<Runnable> primingTasks)
			throws IOException {
		final String hostname = InetAddress.getLocalHost().getHostName().replaceAll("\\..*", "");
		final List<Record> records = new ArrayList<Record>();

		for (final String domain : zoneMap.keySet()) {
			if (superRecords != null && !superRecords.isEmpty()) {
				final List<Record> recholder = zoneMap.get(domain);
				final List<Record> matchHolder = new ArrayList<>();
				// replace the DS records in the exiting record list for the domain with any newer one
				// found in the superRecords
				recholder.replaceAll( record -> {
					for ( final Record superRec : superRecords) {
						if (record.getName().equals(superRec.getName()) && record.getType() == superRec.getType()) {
							matchHolder.add(superRec);
							return superRec;
						}
					}
					return record;
				});

				// Add any new DS records that weren't in the existing record list for the domain
				superRecords.forEach(record -> {
					if (!matchHolder.contains(record)){
						recholder.add(record);
					}
				});
				zoneMap.replace(domain, recholder);
			}

			records.addAll(createZone(domain, zoneMap, dsMap, tr, zc, dzc, generationTasks, primingTasks, hostname));
		}

		return records;
	}

	@SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.NPathComplexity"})
	private static List<Record> createZone(final String domain, final Map<String, List<Record>> zoneMap, final Map<String, DeliveryService> dsMap, 
			final TrafficRouter tr, final LoadingCache<ZoneKey, Zone> zc, final LoadingCache<ZoneKey, Zone> dzc, final List<Runnable> generationTasks, final BlockingQueue<Runnable> primingTasks, final String hostname) throws IOException {
		final DeliveryService ds = dsMap.get(domain);
		final CacheRegister data = tr.getCacheRegister();
		final JsonNode trafficRouters = data.getTrafficRouters();
		final JsonNode config = data.getConfig();

		JsonNode ttl = null;
		JsonNode soa = null;

		if (ds != null) {
			ttl = ds.getTtls();
			soa = ds.getSoa();
		} else {
			ttl = config.get("ttls");
			soa = config.get("soa");
		}

		final Name name = newName(domain);
		final List<Record> list = zoneMap.get(domain);
		final Name admin = newName(ZoneUtils.getAdminString(soa, "admin", "traffic_ops", domain));
		list.add(new SOARecord(name, DClass.IN,
				ZoneUtils.getLong(ttl, "SOA", 86400), getGlueName(ds, trafficRouters
				.get(hostname), name, hostname), admin,
				ZoneUtils.getLong(soa, "serial", ZoneUtils.getSerial(data.getStats())),
				ZoneUtils.getLong(soa, "refresh", 28800),
				ZoneUtils.getLong(soa, "retry", 7200),
				ZoneUtils.getLong(soa, "expire", 604800),
				ZoneUtils.getLong(soa, "minimum", 60)));
		addTrafficRouters(list, trafficRouters, name, ttl, domain, ds);
		addStaticDnsEntries(list, ds, domain);

		final List<Record> records = new ArrayList<Record>();

		try {
			final long maxTTL = ZoneUtils.getMaximumTTL(list);
			records.addAll(getSignatureManager().generateDSRecords(name, maxTTL));
			list.addAll(getSignatureManager().generateDNSKEYRecords(name, maxTTL));
		} catch (NoSuchAlgorithmException ex) {
			LOGGER.fatal("Unable to create zone: " + ex.getMessage(), ex);
		}

		primeZoneCache(domain, name, list, records, tr, zc, dzc, generationTasks, primingTasks, ds);

		return records;
	}

	@SuppressWarnings("PMD.CyclomaticComplexity")
	private static void primeZoneCache(final String domain, final Name name, final List<Record> list, final List<Record> dsRecords, final TrafficRouter tr,
			final LoadingCache<ZoneKey, Zone> zc, final LoadingCache<ZoneKey, Zone> dzc, final List<Runnable> generationTasks,
			final BlockingQueue<Runnable> primingTasks, final DeliveryService ds) {
		generationTasks.add(() -> {
			try {
				final Zone zone = zc.get(getSignatureManager().generateZoneKey(name, list)); // cause the zone to be loaded into the new cache
				final CacheRegister data = tr.getCacheRegister();
				final JsonNode config = data.getConfig();
				final boolean primeDynCache = JsonUtils.optBoolean(config, "dynamic.cache.primer.enabled", true);

				// prime the dynamic zone cache
				if (primeDynCache && ds != null && ds.isDns()) {
					primingTasks.add(() -> {
						try {
							primeDNSDeliveryServices(domain, dsRecords, tr, dzc, zone, ds, data);
						} catch (TextParseException ex) {
							LOGGER.fatal("Unable to prime dynamic zone " + domain, ex);
						}
					});
				}
			} catch (ExecutionException ex) {
				LOGGER.fatal("Unable to load zone into cache: " + ex.getMessage(), ex);
			}
		});
	}

	@SuppressWarnings("PMD.CyclomaticComplexity")
	private static void primeDNSDeliveryServices(final String domain, final List<Record> dsRecords, final TrafficRouter tr, final LoadingCache<ZoneKey, Zone> dzc,
			final Zone zone, final DeliveryService ds, final CacheRegister data) throws TextParseException {
		final Name domainName = newName(domain);
		final Name edgeName = newName(ds.getRoutingName(), domain);
		final JsonNode config = data.getConfig();
		final int primerLimit = JsonUtils.optInt(config, "dynamic.cache.primer.limit", DEFAULT_PRIMER_LIMIT);

		if (LOGGER.isDebugEnabled()){
			LOGGER.debug("Priming " + edgeName);
		}

		final DNSRequest request = new DNSRequest(zone, domainName,  Type.A);
		request.setDnssec(getSignatureManager().isDnssecEnabled());
		request.setHostname(edgeName.toString(true)); // Name.toString(true) - omit the trailing dot

		for (final CacheLocation cacheLocation : data.getCacheLocations()) {
			final List<Cache> caches = tr.selectCachesByCZ(ds, cacheLocation);

			if (caches == null) {
				continue;
			}

			// calculate number of permutations if maxDnsIpsForLocation > 0 and we're not using consistent DNS routing
			int p = 1;

			if (ds.getMaxDnsIps() > 0 && !tr.isConsistentDNSRouting() && caches.size() > ds.getMaxDnsIps()) {
				for (int c = caches.size(); c > (caches.size() - ds.getMaxDnsIps()); c--) {
					p *= c;
				}
			}

			final Set<List<InetRecord>> pset = new HashSet<>();

			for (int i = 0; i < primerLimit; i++) {
				final List<InetRecord> records = tr.inetRecordsFromCaches(ds, caches, request);
				final DNSRouteResult result = new DNSRouteResult();
				result.setAddresses(records);

				if (!pset.contains(dsRecords)) {
						fillDynamicZone(dzc, zone, request, result);
					pset.add(records);
				}

				if (LOGGER.isDebugEnabled()){
					LOGGER.debug("Primed " + ds.getId() + " @ " + cacheLocation.getId() + "; permutation " + pset
							.size() + "/" + p);
				}

				if (pset.size() == p) {
					break;
				}
			}
		}
	}

	@SuppressWarnings("PMD.CyclomaticComplexity")
	private static void addStaticDnsEntries(final List<Record> list, final DeliveryService ds, final String domain)
			throws TextParseException, UnknownHostException {
		if (ds != null && ds.getStaticDnsEntries() != null) {

			final JsonNode entryList = ds.getStaticDnsEntries();

			for (final JsonNode staticEntry : entryList) {
				try {
					final String type = JsonUtils.getString(staticEntry, "type").toUpperCase();
					final String jsName = JsonUtils.getString(staticEntry, "name");
					final String value = JsonUtils.getString(staticEntry, "value");
					final Name name = newName(jsName, domain);
					long ttl = JsonUtils.optInt(staticEntry, "ttl");

					if (ttl == 0) {
						ttl = ZoneUtils.getLong(ds.getTtls(), type, 60);
					}
					switch (type) {
						case "A":
							list.add(new ARecord(name, DClass.IN, ttl, InetAddress.getByName(value)));
							break;
						case "AAAA":
							list.add(new AAAARecord(name, DClass.IN, ttl, InetAddress.getByName(value)));
							break;
						case "CNAME":
							list.add(new CNAMERecord(name, DClass.IN, ttl, new Name(value)));
							break;
						case "TXT":
							list.add(new TXTRecord(name, DClass.IN, ttl, new String(value)));
							break;
					}
				} catch (JsonUtilsException ex) {
					LOGGER.error(ex);
				}
			}
		}
	}

	@SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.NPathComplexity"})
	private static void addTrafficRouters(final List<Record> list, final JsonNode trafficRouters, final Name name,
			final JsonNode ttl, final String domain, final DeliveryService ds) throws TextParseException, UnknownHostException {

		final boolean ip6RoutingEnabled = (ds == null || (ds != null && ds.isIp6RoutingEnabled())) ? true : false;

		final Iterator<String> keyIter = trafficRouters.fieldNames();
		while (keyIter.hasNext()) {
			final String key = keyIter.next();
			final JsonNode trJo = trafficRouters.get(key);

			//TODO: Note when it doesn't have a status its ONLINE
			if (trJo.has("status") && "OFFLINE".equals(trJo.get("status").asText())) {
				// if "status": "OFFLINE"
				continue;
			}

			final Name trName = newName(key, domain);

			list.add(new NSRecord(name, DClass.IN, ZoneUtils.getLong(ttl, "NS", 60), getGlueName(ds, trJo, name, key)));
			list.add(new ARecord(trName,
					DClass.IN, ZoneUtils.getLong(ttl, "A", 60),
					InetAddress.getByName(JsonUtils.optString(trJo, IP))));

			String ip6 = JsonUtils.optString(trJo, IP6);

			if (ip6 != null && !ip6.isEmpty() && ip6RoutingEnabled) {
				ip6 = ip6.replaceAll("/.*", "");
				list.add(new AAAARecord(trName,
						DClass.IN,
						ZoneUtils.getLong(ttl, AAAA, 60),
						Inet6Address.getByName(ip6)));
			}

			if (ds != null && !ds.isDns()) {
				addHttpRoutingRecords(list, ds.getRoutingName(), domain, trJo, ttl, ip6RoutingEnabled);
			}
		}
	}

	private static void addHttpRoutingRecords(final List<Record> list, final String routingName, final String domain,
			final JsonNode trJo, final JsonNode ttl, final boolean addTrafficRoutersAAAA) throws TextParseException, UnknownHostException {
		final Name trName = newName(routingName, domain);
		list.add(new ARecord(trName,
				DClass.IN,
				ZoneUtils.getLong(ttl, "A", 60),
				InetAddress.getByName(JsonUtils.optString(trJo, IP))));
		String ip6 = JsonUtils.optString(trJo, IP6);
		if (addTrafficRoutersAAAA && ip6 != null && !ip6.isEmpty()) {
			ip6 = ip6.replaceAll("/.*", "");
			list.add(new AAAARecord(trName,
					DClass.IN,
					ZoneUtils.getLong(ttl, AAAA, 60),
					Inet6Address.getByName(ip6)));
		}
	}

	private static Name newName(final String hostname, final String domain) throws TextParseException {
		return newName(hostname + "." + domain);
	}

	private static Name newName(final String fqdn) throws TextParseException {
		if (fqdn.endsWith(".")) {
			return new Name(fqdn);
		} else {
			return new Name(fqdn + ".");
		}
	}

	private static Name getGlueName(final DeliveryService ds, final JsonNode trJo, final Name name, final String trName) throws TextParseException {
		if (ds == null && trJo != null && trJo.has("fqdn") && trJo.get("fqdn").textValue() != null) {
			return newName(trJo.get("fqdn").textValue());
		} else {
			final Name superDomain = new Name(new Name(name.toString(true)), 1);
			return newName(trName, superDomain.toString());
		}
	}

	@SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.NPathComplexity"})
	private static Map<String, List<Record>> populateZoneMap(final Map<String, List<Record>> zoneMap,
			final Map<String, DeliveryService> dsMap, final CacheRegister data) throws IOException {
		final Map<String, List<Record>> superDomains = new HashMap<String, List<Record>>();
		for (final String domain : dsMap.keySet()) {
			zoneMap.put(domain, new ArrayList<>());
		}
		for (final Cache c : data.getCacheMap().values()) {
			for (final DeliveryServiceReference dsr : c.getDeliveryServices()) {
				final DeliveryService ds = data.getDeliveryService(dsr.getDeliveryServiceId());
				if (ds == null) {
					LOGGER.error(dsr
							.getDeliveryServiceId() + " delivery service not found in CacheRegister. This means" +
							" the delivery service has been removed or its name changed and the link " +
							"in cache: " + c.getFqdn() + " has not been corrected. The DNS record will not be created" +
							".");
					continue;
				}
				final String fqdn = dsr.getFqdn();
				final String host = dsr.getHost();
				final String domain = dsr.getDomain();

				if (dsMap.containsKey(domain)) {
					final List<Record> zholder = zoneMap.get(domain);
					final String superdomain = domain.split("\\.", 2)[1];
					if (!superDomains.containsKey(superdomain)) {
						superDomains.put(superdomain, getExistingDSRecords(superdomain, ZoneManager.zoneCache));
					}
					if (ds.isDns() && host.equalsIgnoreCase(ds.getRoutingName())) {
						continue;
					}
					try {
						final Name name = newName(fqdn);
						final JsonNode ttl = ds.getTtls();
						try {
							zholder.add(new ARecord(name, DClass.IN, ZoneUtils.getLong(ttl, "A", 60), c.getIp4()));
						} catch (IllegalArgumentException e) {
							LOGGER.warn(e + " : " + c.getIp4());
						}
						final InetAddress ip6 = c.getIp6();
						if (ip6 != null && ds != null && ds.isIp6RoutingEnabled()) {
							zholder.add(new AAAARecord(name, DClass.IN, ZoneUtils.getLong(ttl, AAAA, 60), ip6));
						}
						if (LOGGER.isDebugEnabled()){
							LOGGER.debug("Added records for: "+ds.getId()+", - "+zholder.toString() );
						}
					} catch (org.xbill.DNS.TextParseException e) {
						LOGGER.error("Caught fatal exception while generating zone data for " + fqdn + "!", e);
					}
				}
			}
		}
		return superDomains;
	}

	private static List<Record> getExistingDSRecords(final String superdomain, final LoadingCache<ZoneKey, Zone> zc){
		final List<Record> superRecords = new ArrayList<>();
		if (zc == null) {
			return superRecords;
		}
		final Map<ZoneKey, Zone> existingZones = zc.asMap();
		for (final ZoneKey zoneKey : existingZones.keySet()) {
			if (zoneKey.getName().toString(true).equals(superdomain)) {
				final Iterator zoneIter = existingZones.get(zoneKey).iterator();
				while (zoneIter.hasNext()){
					final RRset rRset = (RRset) zoneIter.next();
					final Iterator rrs = rRset.rrs();
					while (rrs.hasNext()){
						final Record record = (Record) rrs.next();
						if (record.getType() == Type.DS){
							superRecords.add(record);
						}
					}
				}
				break;
			}
		}
		return superRecords;
	}

	/**
	 * Gets trafficRouter.
	 *
	 * @return the trafficRouter
	 */
	public TrafficRouter getTrafficRouter() {
		return trafficRouter;
	}

	/**
	 * Attempts to find a {@link Zone} that would contain the specified {@link Name}.
	 *
	 * @param name the Name to use to attempt to find the Zone
	 * @return the Zone to use to resolve the specified Name
	 */
	public Zone getZone(final Name name) {
		return getZone(name, 0);
	}

	/**
	 * Attempts to find a {@link Zone} that would contain the specified {@link Name}.
	 *
	 * @param name  the Name to use to attempt to find the Zone
	 * @param qtype the Type to use to control Zone ordering
	 * @return the Zone to use to resolve the specified Name
	 */
	public Zone getZone(final Name name, final int qtype) {
		final Map<ZoneKey, Zone> zoneMap = zoneCache.asMap();
		final List<ZoneKey> sorted = new ArrayList<ZoneKey>(zoneMap.keySet());

		Zone result = null;
		final Name target = name;

		Collections.sort(sorted);

		// put the superDomains at the beginning of the list so we look there first for DS records
		if (qtype == Type.DS) {
			Collections.reverse(sorted);
		}

		for (final ZoneKey key : sorted) {
			final Zone zone = zoneMap.get(key);
			final Name origin = zone.getOrigin();

			if (target.subdomain(origin)) {
				result = zone;
				break;
			}
		}

		return result;
	}

	/**
	 * Creates a dynamic zone that serves a set of A and AAAA records for the specified {@link Name}
	 * .
	 *
	 * @param staticZone
	 *            The Zone that would normally serve this request
	 * @param builder
	 *            DNSAccessRecord.Builder access logging
	 * @param request
	 * 	          DNSRequest representing the query
	 * @return the new Zone to serve the request or null if the static Zone should be used
	 */
	private Zone createDynamicZone(final Zone staticZone, final DNSAccessRecord.Builder builder, final DNSRequest request) {
		if (request.getClientIP()==null) {
			return staticZone;
		}

		final Track track = StatTracker.getTrack();

		try {
			final DNSRouteResult result = trafficRouter.route(request, track);

			if (result != null) {
				return fillDynamicZone(dynamicZoneCache, staticZone, request, result);
			} else {
				return null;
			}
		} catch (final Exception e) {
			LOGGER.error(e.getMessage(), e);
		} finally {
			builder.resultType(track.getResult());
			builder.resultDetails(track.getResultDetails());
			builder.resultLocation(track.getResultLocation());
			statTracker.saveTrack(track);
		}

		return null;
	}

	@SuppressWarnings({"PMD.CyclomaticComplexity", "PMD.NPathComplexity"})
	private static Zone fillDynamicZone(final LoadingCache<ZoneKey, Zone> dzc, final Zone staticZone, final DNSRequest request, final DNSRouteResult result) {
		if (result == null || result.getAddresses() == null) {
			return null;
		}

		try {
			int recordsAdded = 0;
			final List<Record> records = createZoneRecords(staticZone);

			for (final InetRecord address : result.getAddresses()) {
				final DeliveryService ds = result.getDeliveryService();
				Name name = request.getName();

				if (address.getType() == Type.NS) {
					name = staticZone.getOrigin();
				} else if (ds != null && (address.getType() == Type.A || address.getType() == Type.AAAA)) {
					final String routingName = ds.getRoutingName();
					name = new Name(routingName, staticZone.getOrigin()); // routingname.ds.cdn.tld
				}

				final Record record = createRecord(name, address);

				if (record != null) {
					records.add(record);
					recordsAdded++;
				}
			}

			if (recordsAdded > 0) {
				try {
					final ZoneKey zoneKey = getSignatureManager()
							.generateDynamicZoneKey(staticZone.getOrigin(), records, request.isDnssec());
					final Zone zone = dzc.get(zoneKey);
					return zone;
				} catch (ExecutionException e) {
					LOGGER.error(e, e);
				}

				return new Zone(staticZone.getOrigin(), records.toArray(new Record[records.size()]));
			}
		} catch (final IOException e) {
			LOGGER.error(e.getMessage(), e);
		}

		return null;
	}

	private static Record createRecord(final Name name, final InetRecord address) throws TextParseException {
		Record record = null;

		if (address.isAlias()) {
			record = new CNAMERecord(name, DClass.IN, address.getTTL(), newName(address.getAlias()));
		} else if (address.getType() == Type.NS) {
			final Name tld = getTopLevelDomain();
			String target = address.getTarget();

			// fix up target to be TR host name plus top level domain
			if (name.subdomain(tld) && !name.equals(tld)) {
				target = String.format("%s.%s", target.split("\\.", 2)[0], tld.toString());
			}

			record = new NSRecord(name, DClass.IN, address.getTTL(), newName(target));
		} else if (address.isInet4()) { // address instanceof Inet4Address
			record = new ARecord(name, DClass.IN, address.getTTL(), address.getAddress());
		} else if (address.isInet6()) {
			record = new AAAARecord(name, DClass.IN, address.getTTL(), address.getAddress());
		}

		return record;
	}

	private static List<Record> createZoneRecords(final Zone staticZone) throws IOException {
		final List<Record> records = new ArrayList<Record>();
		records.add(staticZone.getSOA());

		@SuppressWarnings("unchecked")
		final Iterator<Record> ns = staticZone.getNS().rrs();

		while (ns.hasNext()) {
			records.add(ns.next());
		}

		return records;
	}

	private List<InetRecord> lookup(final Name qname, final Zone zone, final int type) {
		final List<InetRecord> ipAddresses = new ArrayList<InetRecord>();
		final SetResponse sr = zone.findRecords(qname, type);

		if (sr.isSuccessful()) {
			final RRset[] answers = sr.answers();

			for (final RRset answer : answers) {
				@SuppressWarnings("unchecked") final Iterator<Record> it = answer.rrs();

				while (it.hasNext()) {
					final Record r = it.next();

					if (r instanceof ARecord) {
						final ARecord ar = (ARecord)r;
						ipAddresses.add(new InetRecord(ar.getAddress(), ar.getTTL()));
					} else if (r instanceof AAAARecord) {
						final AAAARecord ar = (AAAARecord)r;
						ipAddresses.add(new InetRecord(ar.getAddress(), ar.getTTL()));
					}
				}
			}

			return ipAddresses;
		}

		return null;
	}

	public List<InetRecord> resolve(final String fqdn) {
		try {
			final Name name = new Name(fqdn);
			final Zone zone = getZone(name);
			if (zone == null) {
				LOGGER.error("No zone - Defaulting to system resolver: " + fqdn);
				return super.resolve(fqdn);
			}
			return lookup(name, zone, Type.A);
		} catch (TextParseException e) {
			LOGGER.warn("TextParseException from: " + fqdn, e);
		}

		return null;
	}

	public List<InetRecord> resolve(final String fqdn, final String address, final DNSAccessRecord.Builder builder) throws UnknownHostException {
		try {
			final Name name = new Name(fqdn);
			Zone zone = getZone(name);
			final InetAddress addr = InetAddress.getByName(address);
			final int qtype = (addr instanceof Inet6Address) ? Type.AAAA : Type.A;
			final DNSRequest request = new DNSRequest(zone, name, qtype);
			request.setClientIP(addr.getHostAddress());
			request.setHostname(name.relativize(Name.root).toString());
			request.setDnssec(true);

			final Zone dynamicZone = createDynamicZone(zone, builder, request);

			if (dynamicZone != null) {
				zone = dynamicZone;
			}

			if (zone == null) {
				LOGGER.error("No such zone - Defaulting to system resolver: " + fqdn);
				return super.resolve(fqdn);
			}

			return lookup(name, zone, Type.A);
		} catch (TextParseException e) {
			LOGGER.error("TextParseException from: " + fqdn);
		}

		return null;
	}

	public Zone getZone(final Name qname, final int qtype, final InetAddress clientAddress, final boolean isDnssecRequest, final DNSAccessRecord.Builder builder) {
		final Zone zone = getZone(qname, qtype);

		if (zone == null) {
			return null;
		}

		final SetResponse sr = zone.findRecords(qname, qtype);

		if (sr.isSuccessful()) {
			return zone;
		} else if (isDnsRoutingName(qname.toString().toLowerCase().split("\\.")[0])) {

			final DNSRequest request = new DNSRequest(zone, qname, qtype);
			request.setClientIP(clientAddress.getHostAddress());
			request.setHostname(qname.relativize(Name.root).toString());
			request.setDnssec(isDnssecRequest);

			final Zone dynamicZone = createDynamicZone(zone, builder, request);

			if (dynamicZone != null) {
				return dynamicZone;
			}
		}

		return zone;
	}

	private static boolean isDnsRoutingName(final String routingName) {
		return getDnsRoutingNames().contains(routingName);
	}

	private static Set<String> getDnsRoutingNames() {
		return dnsRoutingNames;
	}

	private static void setDnsRoutingNames(final Set<String> dnsRoutingNames) {
		ZoneManager.dnsRoutingNames = dnsRoutingNames;
	}

	public static File getZoneDirectory() {
		return zoneDirectory;
	}

	public static void setZoneDirectory(final File zoneDirectory) {
		ZoneManager.zoneDirectory = zoneDirectory;
	}

	protected static Name getTopLevelDomain() {
		return topLevelDomain;
	}

	private static void setTopLevelDomain(final Name topLevelDomain) {
		ZoneManager.topLevelDomain = topLevelDomain;
	}

	public CacheStats getStaticCacheStats() {
		return zoneCache.stats();
	}

	public CacheStats getDynamicCacheStats() {
		return dynamicZoneCache.stats();
	}
}
