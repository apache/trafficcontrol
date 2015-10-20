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

package com.comcast.cdn.traffic_control.traffic_router.core.router;

import java.io.IOException;
import java.net.InetAddress;
import java.net.MalformedURLException;
import java.net.URL;
import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.Collection;
import java.util.Collections;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Random;
import java.util.Set;
import java.util.SortedMap;
import java.util.TreeMap;

import org.apache.commons.pool.ObjectPool;
import org.apache.log4j.Logger;
import org.json.JSONException;
import org.json.JSONObject;
import org.xbill.DNS.Name;
import org.xbill.DNS.Zone;

import com.comcast.cdn.traffic_control.traffic_router.core.TrafficRouterException;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.Cache;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheLocation;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheRegister;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.InetRecord;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.Cache.DeliveryServiceReference;
import com.comcast.cdn.traffic_control.traffic_router.core.dns.ZoneManager;
import com.comcast.cdn.traffic_control.traffic_router.core.dns.DNSAccessRecord;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.Dispersion;
import com.comcast.cdn.traffic_control.traffic_router.core.hash.HashFunction;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.FederationRegistry;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.Geolocation;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.GeolocationException;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.GeolocationService;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.NetworkNode;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.NetworkNodeException;
import com.comcast.cdn.traffic_control.traffic_router.core.request.DNSRequest;
import com.comcast.cdn.traffic_control.traffic_router.core.request.HTTPRequest;
import com.comcast.cdn.traffic_control.traffic_router.core.request.Request;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track.ResultType;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track.RouteType;
import com.comcast.cdn.traffic_control.traffic_router.core.util.TrafficOpsUtils;
import com.comcast.cdn.traffic_control.traffic_router.core.util.CidrAddress;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track.ResultDetails;

public class TrafficRouter {
	public static final Logger LOGGER = Logger.getLogger(TrafficRouter.class);

	private final CacheRegister cacheRegister;
	private final ZoneManager zoneManager;
	private final GeolocationService geolocationService;
	private final GeolocationService geolocationService6;
	private final ObjectPool hashFunctionPool;
	private final FederationRegistry federationRegistry;

	private final Random random = new Random(System.nanoTime());
	private Set<String> requestHeaders = new HashSet<String>();

	public TrafficRouter(final CacheRegister cr, 
			final GeolocationService geolocationService, 
			final GeolocationService geolocationService6, 
			final ObjectPool hashFunctionPool,
			final StatTracker statTracker,
			final TrafficOpsUtils trafficOpsUtils,
			final FederationRegistry federationRegistry) throws IOException, JSONException, TrafficRouterException {
		this.cacheRegister = cr;
		this.geolocationService = geolocationService;
		this.geolocationService6 = geolocationService6;
		this.hashFunctionPool = hashFunctionPool;
		this.federationRegistry = federationRegistry;
		this.zoneManager = new ZoneManager(this, statTracker, trafficOpsUtils);
	}

	public ZoneManager getZoneManager() {
		return zoneManager;
	}

	/**
	 * Returns a {@link List} of all of the online {@link Cache}s that support the specified
	 * {@link DeliveryService}. If no online caches are found to support the specified
	 * DeliveryService an empty list is returned.
	 * 
	 * @param ds
	 *            the DeliveryService to check
	 * @return collection of supported caches
	 */
	protected List<Cache> getSupportingCaches(final List<Cache> caches, final DeliveryService ds) {
		for(int i = 0; i < caches.size(); i++) {
			final Cache cache = caches.get(i);
			boolean isAvailable = true;
			if(cache.hasAuthority()) {
				isAvailable = cache.isAvailable();
			}
			if (!isAvailable || !cacheSupportsDeliveryService(cache, ds)) {
				caches.remove(i);
				i--;
			}
		}
		return caches;
	}

	private boolean cacheSupportsDeliveryService(final Cache cache, final DeliveryService ds) {
		boolean result = false;
		for (final DeliveryServiceReference dsRef : cache.getDeliveryServices()) {
			if (dsRef.getDeliveryServiceId().equals(ds.getId())) {
				result = true;
				break;
			}
		}
		return result;
	}

	public CacheRegister getCacheRegister() {
		return cacheRegister;
	}
	protected DeliveryService selectDeliveryService(final Request request, final boolean isHttp) {
		if(cacheRegister==null) {
			LOGGER.warn("no caches yet");
			return null;
		}
		final DeliveryService ds = cacheRegister.getDeliveryService(request, isHttp);
		if (LOGGER.isDebugEnabled()) {
			LOGGER.debug("Selected DeliveryService: " + ds);
		}
		return ds;
	}

	boolean setState(final JSONObject states) throws UnknownHostException {
		setCacheStates(states.optJSONObject("caches"));
		setDsStates(states.optJSONObject("deliveryServices"));
		return true;
	}
	private boolean setDsStates(final JSONObject dsStates) {
		if(dsStates == null) {
			return false;
		}
		final Map<String, DeliveryService> dsMap = cacheRegister.getDeliveryServices();
		for (final String dsName : dsMap.keySet()) {
			dsMap.get(dsName).setState(dsStates.optJSONObject(dsName));
		}
		return true;
	}
	private boolean setCacheStates(final JSONObject cacheStates) {
		if(cacheStates == null) {
			return false;
		}
		final Map<String, Cache> cacheMap = cacheRegister.getCacheMap();
		if(cacheMap == null) { return false; }
		for (final String cacheName : cacheMap.keySet()) {
			final String monitorCacheName = cacheName.replaceFirst("@.*", "");
			final JSONObject state = cacheStates.optJSONObject(monitorCacheName);
			cacheMap.get(cacheName).setState(state);
		}
		return true;
	}

	protected static final String UNABLE_TO_ROUTE_REQUEST = "Unable to route request.";
	protected static final String URL_ERR_STR = "Unable to create URL.";

	public GeolocationService getGeolocationService() {
		return geolocationService;
	}
	public Geolocation getLocation(final String clientIP) throws GeolocationException {
		if(clientIP.contains(":")) {
			return geolocationService6.location(clientIP);
		}
		return geolocationService.location(clientIP);
	}

	/**
	 * Gets hashFunctionPool.
	 * 
	 * @return the hashFunctionPool
	 */
	public ObjectPool getHashFunctionPool() {
		return hashFunctionPool;
	}

	public List<Cache> getCachesByGeo(final Request request, final DeliveryService ds, final Geolocation clientLocation, final Map<String, Double> resultLocation) throws GeolocationException {
		final String zoneId = null; 
		// the specific use of the popularity zone
		// manager was not understood and not used
		// and was therefore was eliminated
		// final String zoneId = getZoneManager().getZone(request.getRequestedUrl());
		final int locationLimit = ds.getLocationLimit();
		int locationsTested = 0;
		final List<CacheLocation> cacheLocations = orderCacheLocations(request,
				getCacheRegister().getCacheLocations(zoneId), ds, clientLocation);
		for (final CacheLocation location : cacheLocations) {
			final List<Cache> caches = selectCache(location, ds);
			if (caches != null) {
				resultLocation.put("latitude", location.getGeolocation().getLatitude());
				resultLocation.put("longitude", location.getGeolocation().getLongitude());
				return caches;
			}
			locationsTested++;
			if(locationLimit != 0 && locationsTested >= locationLimit) {
				return null;
			}
		}
		return null;
	}
	protected List<Cache> selectCache(final Request request, final DeliveryService ds, final Track track) throws GeolocationException {
		final CacheLocation cacheLocation = getCoverageZoneCache(request.getClientIP());
		List<Cache> caches = selectCachesByCZ(ds, cacheLocation, track);

		if (caches != null) {
			return caches;
		}

		if (ds.isCoverageZoneOnly()) {
			LOGGER.warn(String.format("No Cache found in CZM (%s, ip=%s, path=%s), geo not supported", request.getType(), request.getClientIP(), request.getHostname()));
			track.setResult(ResultType.MISS);
			track.setResultDetails(ResultDetails.DS_CZ_ONLY);
		}
		else {
			LOGGER.warn(String.format("No Cache found by CZM (%s, ip=%s, path=%s)", request.getType(), request.getClientIP(), request.getHostname()));
			caches = selectCachesByGeo(request, ds, cacheLocation, track);
		}

		return caches;
	}

	public List<Cache> selectCachesByGeo(final Request request, final DeliveryService deliveryService, final CacheLocation cacheLocation, final Track track) throws GeolocationException {

		Geolocation clientLocation = null;

		try {
			clientLocation = getClientLocation(request, deliveryService, cacheLocation);
		} catch (GeolocationException e) {
			LOGGER.warn("Failed looking up Client GeoLocation: " + e.getMessage());
		}

		if (clientLocation == null) {
			track.setResultDetails(ResultDetails.DS_CLIENT_GEO_UNSUPPORTED);
			return null;
		}
		
		final Map<String, Double> resultLocation = new HashMap<String, Double>();

		final List<Cache> caches = getCachesByGeo(request, deliveryService, clientLocation, resultLocation);
		
		if (caches == null || caches.isEmpty()) {
			LOGGER.warn(String.format("No Cache found by Geo (%s, ip=%s, path=%s)", request.getType(), request.getClientIP(), request.getHostname()));
			track.setResultDetails(ResultDetails.GEO_NO_CACHE_FOUND);
		}

		track.setResult(ResultType.GEO);
		return caches;
	}

	public DNSRouteResult route(final DNSRequest request, final Track track) throws GeolocationException {
		track.setRouteType(RouteType.DNS, request.getHostname());

		final DeliveryService ds = selectDeliveryService(request, false);

		if (ds == null) {
			LOGGER.warn("[dns] No DeliveryService found for: " + request.getHostname());
			track.setResult(ResultType.STATIC_ROUTE);
			track.setResultDetails(ResultDetails.DS_NOT_FOUND);
			return null;
		}

		final DNSRouteResult result = new DNSRouteResult();

		if (!ds.isAvailable()) {
			LOGGER.warn("deliveryService not available: " + ds);
			result.setAddresses(ds.getFailureDnsResponse(request, track));
			return result;
		}

		final CacheLocation cacheLocation = getCoverageZoneCache(request.getClientIP());
		List<Cache> caches = selectCachesByCZ(ds, cacheLocation, track);

		if (caches != null) {
			track.setResult(ResultType.CZ);
			result.setAddresses(inetRecordsFromCaches(ds, caches));
			return result;
		}

		if (ds.isCoverageZoneOnly()) {
			LOGGER.info(String.format("No Cache found in CZM (%s, ip=%s, path=%s), geo not supported", request.getType(), request.getClientIP(), request.getHostname()));
			track.setResult(ResultType.MISS);
			track.setResultDetails(ResultDetails.DS_CZ_ONLY);
			result.setAddresses(ds.getFailureDnsResponse(request, track));
			return result;
		}

		try {
			final List<InetRecord> inetRecords = federationRegistry.findInetRecords(ds.getId(), CidrAddress.fromString(request.getClientIP()));

			if (inetRecords != null && !inetRecords.isEmpty()) {
				result.setAddresses(inetRecords);
				track.setResult(ResultType.FED);
				return result;
			}
		} catch (NetworkNodeException e) {
			LOGGER.error("Bad client address: '" + request.getClientIP() + "'");
		}

		LOGGER.info(String.format("No Cache found by CZM (%s, ip=%s, path=%s)", request.getType(), request.getClientIP(), request.getHostname()));
		caches = selectCachesByGeo(request, ds, cacheLocation, track);

		if (caches != null) {
			track.setResult(ResultType.GEO);
			result.setAddresses(inetRecordsFromCaches(ds, caches));
		}
		else {
			track.setResult(ResultType.MISS);
			result.setAddresses(ds.getFailureDnsResponse(request, track));
		}

		return result;
	}

	private List<InetRecord> inetRecordsFromCaches(final DeliveryService ds, final List<Cache> caches) {
		final List<InetRecord> addresses = new ArrayList<InetRecord>();
		final int maxDnsIps = ds.getMaxDnsIps();

		/*
		 * We also shuffle in NameServer when adding Records to the Message prior
		 * to sending it out, as the Records are sorted later when we fill the
		 * dynamic zone if DNSSEC is enabled. We shuffle here prior to pruning
		 * for maxDnsIps so that we ensure we are spreading load across all caches
		 * assigned to this delivery service.
		 */
		if (maxDnsIps > 0) {
			Collections.shuffle(caches, random);
		}

		int i = 0;

		for (final Cache cache : caches) {
			if (maxDnsIps!=0 && i >= maxDnsIps) {
				break;
			}

			i++;

			addresses.addAll(cache.getIpAddresses(ds.getTtls(), zoneManager, ds.isIp6RoutingEnabled()));
		}
		return addresses;
	}

	public Geolocation getClientLocation(final Request request, final DeliveryService ds, final CacheLocation cacheLocation) throws GeolocationException {
		Geolocation clientLocation;
		if (cacheLocation != null) {
			clientLocation = cacheLocation.getGeolocation();
		} else {
			clientLocation = getLocation(request.getClientIP());
			clientLocation = ds.supportLocation(clientLocation, request.getType());
		}
		return clientLocation;
	}

	private List<Cache> selectCachesByCZ(final DeliveryService ds, final CacheLocation cacheLocation, final Track track) {
		if (cacheLocation == null || !ds.isLocationAvailable(cacheLocation)) {
			return null;
		}

		final List<Cache> caches = selectCache(cacheLocation, ds);

		if (caches != null) {
			track.setResult(ResultType.CZ);
			track.setResultLocation(cacheLocation.getGeolocation());
		}

		return caches;
	}

	public HTTPRouteResult route(final HTTPRequest request, final Track track) throws MalformedURLException, GeolocationException {
		track.setRouteType(RouteType.HTTP, request.getHostname());

		final DeliveryService ds = selectDeliveryService(request, true);

		if (ds == null) {
			LOGGER.warn("No DeliveryService found for: " + request.getRequestedUrl());
			track.setResult(ResultType.DS_MISS);
			track.setResultDetails(ResultDetails.DS_NOT_FOUND);
			return null;
		}

		final HTTPRouteResult routeResult = new HTTPRouteResult();

		routeResult.setDeliveryService(ds);

		if (!ds.isAvailable()) {
			LOGGER.warn("deliveryService unavailable: " + ds);
			routeResult.setUrl(ds.getFailureHttpResponse(request, track));
			return routeResult;
		}

		final List<Cache> caches = selectCache(request, ds, track);

		if (caches == null) {
			routeResult.setUrl(ds.getFailureHttpResponse(request, track));
			return routeResult;
		}

		final Dispersion dispersion = ds.getDispersion();
		final Cache cache = dispersion.getCache(consistentHash(caches, request.getPath()));

		routeResult.setUrl(new URL(ds.createURIString(request, cache)));

		return routeResult;
	}

	protected CacheLocation getCoverageZoneCache(final String ip) {
		NetworkNode nn = null;
		try {
			nn = NetworkNode.getInstance().getNetwork(ip);
		} catch (NetworkNodeException e) {
			LOGGER.warn(e);
		}
		if (nn == null) {
			return null;
		}

		final String locId = nn.getLoc();
		final CacheLocation cl = nn.getCacheLocation();
		if(cl != null) {
			return cl;
		}
		if(locId == null) {
			return null;
		}

			// find CacheLocation
		final Collection<CacheLocation> caches = getCacheRegister()
				.getCacheLocations();
		for (final CacheLocation cl2 : caches) {
			if (cl2.getId().equals(locId)) {
				nn.setCacheLocation(cl2);
				return cl2;
			}
		}

		return null;
	}

	/**
	 * Utilizes the hashValues stored with each cache to select the cache that
	 * the specified hash should map to.
	 *
	 * @param caches
	 *            the list of caches to choose from
	 * @param hash
	 *            the hash value for the request
	 * @return a cache or null if no cache can be found to map to
	 */
	protected Cache consistentHashOld(final List<Cache> caches,
			final String request) {
		double hash = 0;
		HashFunction hashFunction = null;
		try {
			hashFunction = (HashFunction) hashFunctionPool.borrowObject();
			try {
				hash = hashFunction.hash(request);
			} catch (final Exception e) {
				LOGGER.debug(e.getMessage(), e);
			}
			hashFunctionPool.returnObject(hashFunction);
		} catch (final Exception e) {
			LOGGER.debug(e.getMessage(), e);
		}
		if (hash == 0) {
			LOGGER.warn("Problem with hashFunctionPool, request: " + request);
			return null;
		}

		return searchCacheOld(caches, hash);
	}

	private Cache searchCacheOld(final List<Cache> caches, final double hash) {
		Cache minCache = null;
		double minHash = Double.MAX_VALUE;
		Cache foundCache = null;
		double minDiff = Double.MAX_VALUE;

		for (final Cache cache : caches) {
			for (final double hashValue : cache.getHashValues()) {
				if (hashValue < minHash) {
					minCache = cache;
					minHash = hashValue;
				}
				final double diff = hashValue - hash;
				if ((diff >= 0) && (diff < minDiff)) {
					foundCache = cache;
					minDiff = diff;
				}
			}
		}

		final Cache result = (foundCache != null) ? foundCache : minCache;
		if (LOGGER.isDebugEnabled()) {
			LOGGER.debug("Selected cache: " + result);
		}
		return result;
	}

	/**
	 * Utilizes the hashValues stored with each cache to select the cache that
	 * the specified hash should map to.
	 * 
	 * @param caches
	 *            the list of caches to choose from
	 * @param request
	 *            the request string from which the hash will be generated
	 * @return a cache or null if no cache can be found to map to
	 */
	protected SortedMap<Double, Cache> consistentHash(final List<Cache> caches,
			final String request) {
		double hash = 0;
		HashFunction hashFunction = null;
		try {
			hashFunction = (HashFunction) hashFunctionPool.borrowObject();
			try {
				hash = hashFunction.hash(request);
			} catch (final Exception e) {
				LOGGER.debug(e.getMessage(), e);
			}
			hashFunctionPool.returnObject(hashFunction);
		} catch (final Exception e) {
			LOGGER.debug(e.getMessage(), e);
		}
		if (hash == 0) {
			LOGGER.warn("Problem with hashFunctionPool, request: " + request);
			return null;
		}

		final SortedMap<Double, Cache> cacheMap = new TreeMap<Double, Cache>();

		for (final Cache cache : caches) {
			final double r = cache.getClosestHash(hash);
			if (r == 0) {
				LOGGER.warn("Error: getClosestHash returned 0: " + cache);
				return null;
			}

			double diff = Math.abs(r - hash);

			if (cacheMap.containsKey(diff)) {
				LOGGER.warn("Error: cacheMap contains diff " + diff + "; incrementing to avoid collision");
				long bits = Double.doubleToLongBits(diff);

				while (cacheMap.containsKey(diff)) {
					bits++;
					diff = Double.longBitsToDouble(bits);
				}
			}

			cacheMap.put(diff, cache);
		}

		return cacheMap;
	}

	/**
	 * Returns a list {@link CacheLocation}s sorted by distance from the client.
	 * If the client's location could not be determined, then the list is
	 * unsorted.
	 * 
	 * @param request
	 *            the client's request
	 * @param cacheLocations
	 *            the collection of CacheLocations to order
	 * @param ds
	 * @return the ordered list of locations
	 */
	protected List<CacheLocation> orderCacheLocations(final Request request,
			final Collection<CacheLocation> cacheLocations,
			final DeliveryService ds,
			final Geolocation clientLocation) {
		final List<CacheLocation> locations = new ArrayList<CacheLocation>();
		for(final CacheLocation cl : cacheLocations) {
			if(ds.isLocationAvailable(cl)) {
				locations.add(cl);
			}
		}

		Collections.sort(locations, new CacheLocationComparator(
				clientLocation));

		return locations;
	}

	/**
	 * Selects a {@link Cache} from the {@link CacheLocation} provided.
	 * 
	 * @param location
	 *            the caches that will considered
	 * @param ds
	 *            the delivery service for the request
	 * @param request
	 *            the request to consider for cache selection
	 * @return the selected cache or null if none can be found
	 */
	private List<Cache> selectCache(final CacheLocation location,
			final DeliveryService ds) {
		if (LOGGER.isDebugEnabled()) {
			LOGGER.debug("Trying location: " + location.getId());
		}

		final List<Cache> caches = getSupportingCaches(location.getCaches(), ds);
		if (caches.isEmpty()) {
			if (LOGGER.isDebugEnabled()) {
				LOGGER.debug("No online, supporting caches were found at location: "
						+ location.getId());
			}
			return null;
		}

		return caches;//consistentHash(caches, request);List<Cache>
	}

	public Zone getZone(final Name qname, final int qtype, final InetAddress clientAddress, final boolean isDnssecRequest, final DNSAccessRecord.Builder builder) {
		return zoneManager.getZone(qname, qtype, clientAddress, isDnssecRequest, builder);
	}

	public void setRequestHeaders(final Set<String> requestHeaders) {
		this.requestHeaders = requestHeaders;
	}

	public Set<String> getRequestHeaders() {
		return requestHeaders;
	}
}
