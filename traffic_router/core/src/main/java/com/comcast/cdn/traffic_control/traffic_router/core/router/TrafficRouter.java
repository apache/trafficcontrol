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
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Random;
import java.util.Set;
import java.util.SortedMap;
import java.util.TreeMap;
import java.io.PrintWriter;
import java.io.StringWriter;

import org.apache.commons.pool.ObjectPool;
import org.apache.log4j.Logger;
import org.json.JSONException;
import org.json.JSONObject;
import org.springframework.beans.BeansException;
import org.springframework.context.ApplicationContext;
import org.xbill.DNS.Name;
import org.xbill.DNS.Zone;

import com.comcast.cdn.traffic_control.traffic_router.core.TrafficRouterException;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.Cache;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheLocation;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheRegister;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.InetRecord;
import com.comcast.cdn.traffic_control.traffic_router.core.dns.ZoneManager;
import com.comcast.cdn.traffic_control.traffic_router.core.dns.DNSAccessRecord;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.Dispersion;
import com.comcast.cdn.traffic_control.traffic_router.core.hash.HashFunction;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.FederationRegistry;
import com.comcast.cdn.traffic_control.traffic_router.geolocation.Geolocation;
import com.comcast.cdn.traffic_control.traffic_router.geolocation.GeolocationException;
import com.comcast.cdn.traffic_control.traffic_router.geolocation.GeolocationService;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.NetworkNode;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.NetworkNodeException;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.RegionalGeo;
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
	private final boolean consistentDNSRouting;

	private final Random random = new Random(System.nanoTime());
	private Set<String> requestHeaders = new HashSet<String>();
	private static final Geolocation GEO_ZERO_ZERO = new Geolocation(0,0);
	private ApplicationContext applicationContext;

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
		this.consistentDNSRouting = cr.getConfig().optBoolean("consistent.dns.routing", false); // previous/default behavior
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
	public List<Cache> getSupportingCaches(final List<Cache> caches, final DeliveryService ds) {
		for(int i = 0; i < caches.size(); i++) {
			final Cache cache = caches.get(i);
			boolean isAvailable = true;
			if(cache.hasAuthority()) {
				isAvailable = cache.isAvailable();
			}
			if (!isAvailable || !cache.hasDeliveryService(ds.getId())) {
				caches.remove(i);
				i--;
			}
		}
		return caches;
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
		return clientIP.contains(":") ? geolocationService6.location(clientIP) : geolocationService.location(clientIP);
	}

	public Geolocation getLocation(final String clientIP, final DeliveryService deliveryService) throws GeolocationException {
		final String geolocationProvider = deliveryService.getGeolocationProvider();

		if (applicationContext == null) {
			LOGGER.error("ApplicationContext not set unable to use custom geolocation service providers");
		}

		if (geolocationProvider != null && !geolocationProvider.isEmpty() && applicationContext != null) {
			try {
				return ((GeolocationService) applicationContext.getBean(geolocationProvider)).location(clientIP);
			} catch (BeansException e) {
				LOGGER.error("Failed getting providing class '" + geolocationProvider + "' for geolocation for delivery service " + deliveryService.getId() + " falling back to maxmind");
			}
		}

		return getLocation(clientIP);
	}

	/**
	 * Gets hashFunctionPool.
	 * 
	 * @return the hashFunctionPool
	 */
	public ObjectPool getHashFunctionPool() {
		return hashFunctionPool;
	}

	public List<Cache> getCachesByGeo(final Request request, final DeliveryService ds, final Geolocation clientLocation, final Track track) throws GeolocationException {
		int locationsTested = 0;

		final int locationLimit = ds.getLocationLimit();
		final List<CacheLocation> cacheLocations = orderCacheLocations(request, getCacheRegister().getCacheLocations(null), ds, clientLocation);

		for (final CacheLocation location : cacheLocations) {
			final List<Cache> caches = selectCache(location, ds);
			if (caches != null) {
				track.setResultLocation(location.getGeolocation());
				if (track.getResultLocation().equals(GEO_ZERO_ZERO)) {
					LOGGER.error("Location " + location.getId() + " has Geolocation " + location.getGeolocation());
				}
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
			if (ds.getGeoRedirectUrl() != null) {
				//use the NGB redirect
				caches = enforceGeoRedirect(track, ds, request, null);
			} else {
				track.setResult(ResultType.MISS);
				track.setResultDetails(ResultDetails.DS_CZ_ONLY);
			}
		} else {
			caches = selectCachesByGeo(request, ds, cacheLocation, track);
		}

		return caches;
	}

	public List<Cache> selectCachesByGeo(final Request request, final DeliveryService deliveryService, final CacheLocation cacheLocation, final Track track) throws GeolocationException {

		Geolocation clientLocation = null;

		try {
			clientLocation = getClientLocation(request, deliveryService, cacheLocation, track);
		} catch (GeolocationException e) {
			LOGGER.warn("Failed looking up Client GeoLocation: " + e.getMessage());
		}

		if (clientLocation == null) {
			if (deliveryService.getGeoRedirectUrl() != null) {
				//will use the NGB redirect
				LOGGER.debug(String
						.format("client is blocked by geolimit, use the NGB redirect url: %s",
							deliveryService.getGeoRedirectUrl()));
				return enforceGeoRedirect(track, deliveryService, request, track.getClientGeolocation());
			} else {
				track.setResultDetails(ResultDetails.DS_CLIENT_GEO_UNSUPPORTED);
				return null;
			}
		}

		final List<Cache> caches = getCachesByGeo(request, deliveryService, clientLocation, track);
		
		if (caches == null || caches.isEmpty()) {
			track.setResultDetails(ResultDetails.GEO_NO_CACHE_FOUND);
		}

		track.setResult(ResultType.GEO);
		return caches;
	}

	public DNSRouteResult route(final DNSRequest request, final Track track) throws GeolocationException {
		track.setRouteType(RouteType.DNS, request.getHostname());

		final DeliveryService ds = selectDeliveryService(request, false);

		if (ds == null) {
			track.setResult(ResultType.STATIC_ROUTE);
			track.setResultDetails(ResultDetails.DS_NOT_FOUND);
			return null;
		}

		final DNSRouteResult result = new DNSRouteResult();

		if (!ds.isAvailable()) {
			result.setAddresses(ds.getFailureDnsResponse(request, track));
			return result;
		}

		final CacheLocation cacheLocation = getCoverageZoneCache(request.getClientIP());
		List<Cache> caches = selectCachesByCZ(ds, cacheLocation, track);

		if (caches != null) {
			track.setResult(ResultType.CZ);
			result.setAddresses(inetRecordsFromCaches(ds, caches, request));
			return result;
		}

		if (ds.isCoverageZoneOnly()) {
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

		caches = selectCachesByGeo(request, ds, cacheLocation, track);

		if (caches != null) {
			track.setResult(ResultType.GEO);
			result.setAddresses(inetRecordsFromCaches(ds, caches, request));
		} else {
			track.setResult(ResultType.MISS);
			result.setAddresses(ds.getFailureDnsResponse(request, track));
		}

		return result;
	}

	public List<InetRecord> inetRecordsFromCaches(final DeliveryService ds, final List<Cache> caches, final Request request) {
		final List<InetRecord> addresses = new ArrayList<InetRecord>();
		final int maxDnsIps = ds.getMaxDnsIps();
		List<Cache> selectedCaches;

		if (maxDnsIps > 0 && isConsistentDNSRouting()) { // only consistent hash if we must
			final SortedMap<Double, Cache> cacheMap = consistentHash(caches, request.getHostname());
			final Dispersion dispersion = ds.getDispersion();
			selectedCaches = dispersion.getCacheList(cacheMap);
		} else if (maxDnsIps > 0) {
			/*
			 * We also shuffle in NameServer when adding Records to the Message prior
			 * to sending it out, as the Records are sorted later when we fill the
			 * dynamic zone if DNSSEC is enabled. We shuffle here prior to pruning
			 * for maxDnsIps so that we ensure we are spreading load across all caches
			 * assigned to this delivery service.
			*/
			Collections.shuffle(caches, random);

			selectedCaches = new ArrayList<Cache>();

			for (final Cache cache : caches) {
				selectedCaches.add(cache);

				if (selectedCaches.size() >= maxDnsIps) {
					break;
				}
			}
		} else {
			selectedCaches = caches;
		}

		for (final Cache cache : selectedCaches) {
			addresses.addAll(cache.getIpAddresses(ds.getTtls(), zoneManager, ds.isIp6RoutingEnabled()));
		}

		return addresses;
	}

	public Geolocation getClientGeolocation(final Request request, final Track track, final DeliveryService deliveryService) throws GeolocationException {
		if (track.isClientGeolocationQueried()) {
			return track.getClientGeolocation();
		}

		final Geolocation clientGeolocation = getLocation(request.getClientIP(), deliveryService);
		track.setClientGeolocation(clientGeolocation);
		track.setClientGeolocationQueried(true);

		return clientGeolocation;
	}

	public Geolocation getClientLocation(final Request request, final DeliveryService ds, final CacheLocation cacheLocation, final Track track) throws GeolocationException {
		if (cacheLocation != null) {
			return cacheLocation.getGeolocation();
		}

		final Geolocation clientGeolocation = getClientGeolocation(request, track, ds);
		return ds.supportLocation(clientGeolocation, request.getType());
	}

	public List<Cache> selectCachesByCZ(final DeliveryService ds, final CacheLocation cacheLocation) {
		return selectCachesByCZ(ds, cacheLocation, null);
	}

	private List<Cache> selectCachesByCZ(final DeliveryService ds, final CacheLocation cacheLocation, final Track track) {
		if (cacheLocation == null || !ds.isLocationAvailable(cacheLocation)) {
			return null;
		}

		final List<Cache> caches = selectCache(cacheLocation, ds);

		if (caches != null && track != null) {
			track.setResult(ResultType.CZ);
			track.setResultLocation(cacheLocation.getGeolocation());
		}

		return caches;
	}

	public HTTPRouteResult route(final HTTPRequest request, final Track track) throws MalformedURLException, GeolocationException {
		track.setRouteType(RouteType.HTTP, request.getHostname());

		final DeliveryService ds = selectDeliveryService(request, true);

		if (ds == null) {
			track.setResult(ResultType.DS_MISS);
			track.setResultDetails(ResultDetails.DS_NOT_FOUND);
			return null;
		}

		final HTTPRouteResult routeResult = new HTTPRouteResult();

		routeResult.setDeliveryService(ds);

		if (!ds.isAvailable()) {
			routeResult.setUrl(ds.getFailureHttpResponse(request, track));
			return routeResult;
		}

		final List<Cache> caches = selectCache(request, ds, track);

		if (caches == null) {
			if (track.getResult() == ResultType.GEO_REDIRECT) {
				routeResult.setUrl(new URL(ds.getGeoRedirectUrl()));
				LOGGER.debug(String.format("NGB redirect to url: %s for request: %s", ds.getGeoRedirectUrl()
						, request.getRequestedUrl()));
				return routeResult;
			}

			routeResult.setUrl(ds.getFailureHttpResponse(request, track));
			return routeResult;
		}

		final Dispersion dispersion = ds.getDispersion();
		final Cache cache = dispersion.getCache(consistentHash(caches, request.getPath()));

		if (ds.isRegionalGeoEnabled()) {
			RegionalGeo.enforce(this, request, ds, cache, routeResult, track);
			return routeResult;
		}

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

		return (foundCache != null) ? foundCache : minCache;
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
				LOGGER.error(e.getMessage(), e);
			}
			hashFunctionPool.returnObject(hashFunction);
		} catch (final Exception e) {
			LOGGER.error(e.getMessage(), e);
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
	public List<CacheLocation> orderCacheLocations(final Request request, final Collection<CacheLocation> cacheLocations, final DeliveryService ds, final Geolocation clientLocation) {
		final List<CacheLocation> locations = new ArrayList<CacheLocation>();

		for(final CacheLocation cl : cacheLocations) {
			if(ds.isLocationAvailable(cl)) {
				locations.add(cl);
			}
		}

		Collections.sort(locations, new CacheLocationComparator(clientLocation));

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

	public boolean isConsistentDNSRouting() {
		return consistentDNSRouting;
	}

	private List<Cache> enforceGeoRedirect(final Track track, final DeliveryService ds,
			final Request request, final Geolocation queriedClientLocation) {

		final String urlType = ds.getGeoRedirectUrlType();
		track.setResult(ResultType.GEO_REDIRECT);
		try {
			if ("NOT_DS_URL".equals(urlType)) {
				//redirect url not belongs to this DS, just redirect it
				LOGGER.debug("geo redirect url not belongs to ds: " + ds.getGeoRedirectUrl());
				return null;
			} else if ("DS_URL".equals(urlType)) {
				Geolocation clientLocation = queriedClientLocation;

				//redirect url belongs to this DS, will try return the caches
				if (clientLocation == null) {
					LOGGER.debug("clientLocation null, try to query it");
					clientLocation = getLocation(request.getClientIP(), ds);

					if (clientLocation == null) { clientLocation = ds.getMissLocation(); }

					if (clientLocation == null) {
						LOGGER.error("cannot find a geo location for the client: " + request.getClientIP());
						// particular error was logged in ds.supportLocation
						track.setResult(ResultType.MISS);
						track.setResultDetails(ResultDetails.DS_CLIENT_GEO_UNSUPPORTED);
						return null;
					}
				}

				final List<Cache> caches = getCachesByGeo(request, ds, clientLocation, track);
				if (caches == null) {
					LOGGER.warn(String.format(
								"No Cache found by Geo in NGB redirect"));
					track.setResult(ResultType.MISS);
					track.setResultDetails(ResultDetails.GEO_NO_CACHE_FOUND);
					return null;
				}
				return caches;
			} else {
				LOGGER.error("invalid geo redirect url type");
				track.setResult(ResultType.MISS);
				track.setResultDetails(ResultDetails.GEO_NO_CACHE_FOUND);
				return null;
			}
		} catch (Exception e) {
			LOGGER.error("caught Exception when enforceGeoRedirect: " + e);
			final StringWriter sw = new StringWriter();
			final PrintWriter pw = new PrintWriter(sw);
			e.printStackTrace(pw);
			LOGGER.error(sw.toString());
			track.setResult(ResultType.ERROR);
			return null;
		}
	}

	public void setApplicationContext(final ApplicationContext applicationContext) throws BeansException {
		this.applicationContext = applicationContext;
	}
}
