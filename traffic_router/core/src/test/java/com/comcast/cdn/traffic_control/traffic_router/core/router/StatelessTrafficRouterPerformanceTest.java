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
import java.net.MalformedURLException;
import java.net.URL;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Random;
import java.util.concurrent.atomic.AtomicInteger;

import org.apache.commons.pool.ObjectPool;
import org.junit.Before;
import org.springframework.context.ApplicationContext;
import org.springframework.context.support.ClassPathXmlApplicationContext;

import com.comcast.cdn.traffic_control.traffic_router.core.TrafficRouterException;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.Cache;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheLocation;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheRegister;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.Dispersion;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.Geolocation;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.GeolocationException;
import com.comcast.cdn.traffic_control.traffic_router.core.loc.GeolocationService;
import com.comcast.cdn.traffic_control.traffic_router.core.request.HTTPRequest;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouter;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker;
import com.comcast.cdn.traffic_control.traffic_router.core.router.StatTracker.Track;

public class StatelessTrafficRouterPerformanceTest  extends TrafficRouter {


    public StatelessTrafficRouterPerformanceTest(CacheRegister cr,
			GeolocationService geolocationService, ObjectPool hashFunctionPool)
			throws IOException {
		super(cr, geolocationService, null, hashFunctionPool, null, "edge", "ccr");
	}

	@Before
    public void setUp() throws Exception {
		final ApplicationContext context = new ClassPathXmlApplicationContext("/applicationContext.xml");
		TrafficRouter trafficRouter = (TrafficRouter) context.getBean("TrafficRouter");
		HTTPRequest req = new HTTPRequest();
		req.setClientIP("10.0.0.15");
		req.setPath("/QualityLevels(96000)/Fragments(audio_eng=20720000000)");
		req.setQueryString("");
		req.setHostname("somehost.cdn.net");
		req.setRequestedUrl("http://somehost.cdn.net/QualityLevels(96000)/Fragments(audio_eng=20720000000)");
		Track track = StatTracker.getTrack();
		try {
			URL url = trafficRouter.route(req, track);
			LOGGER.warn("url: "+url);
		} catch (Exception e2) {
			e2.printStackTrace();
		}

		String[] rstrs = getRandomStrings(1000000);
		try {
			((StatelessTrafficRouterPerformanceTest) trafficRouter).routeTest(req, rstrs);
		} catch (TrafficRouterException e) {
			e.printStackTrace();
		}
		for(String str : rstrs) {
			req.setPath(str);
			req.setRequestedUrl("http://somehost.cdn.net/"+str);

		}

		System.exit(0);
    }
	static class RandomString {

		private static final char[] symbols = new char[36];

		static {
			for (int idx = 0; idx < 10; ++idx)
				symbols[idx] = (char) ('0' + idx);
			for (int idx = 10; idx < 36; ++idx)
				symbols[idx] = (char) ('a' + idx - 10);
		}

		private final Random random = new Random();
		private final char[] buf;

		public RandomString(int length) {
			if (length < 1)
				throw new IllegalArgumentException("length < 1: " + length);
			buf = new char[length];
		}

		public String nextString() {
			for (int idx = 0; idx < buf.length; ++idx) 
				buf[idx] = symbols[random.nextInt(symbols.length)];
			return new String(buf);
		}

	}
	private static String[] getRandomStrings(int cnt) {
		RandomString rs = new RandomString(10);
		String[] strs = new String[cnt];
		for(int i = 0; i < cnt; i++) {
			strs[i] = rs.nextString();
		}
		return strs;
	}

	public URL routeTest(final HTTPRequest request, String[] rstrs) throws TrafficRouterException, GeolocationException {
		if (LOGGER.isDebugEnabled()) {
			LOGGER.debug("Attempting to route HTTPRequest: " + request.getRequestedUrl());
		}

		final String ip = request.getClientIP();
		final DeliveryService ds = selectDeliveryService(request, true);
		if(ds == null) {
			return null;
		}
		final StatTracker.Track track = StatTracker.getTrack();
		List<Cache> caches = selectCache(request, ds, track, true);
		Dispersion dispersion = ds.getDispersion();
		Cache cache = dispersion.getCache(consistentHash(caches, request.getPath()));
		try {
			if(cache != null) { return new URL(ds.createURIString(request, cache)); }
		} catch (final MalformedURLException e) {
			LOGGER.error(e.getMessage(), e);
			throw new TrafficRouterException(URL_ERR_STR, e);
		}
		LOGGER.warn("No Cache found in CoverageZoneMap for HTTPRequest.getClientIP: "+ip);

		final String zoneId = null; // getZoneManager().getZone(request.getRequestedUrl());
		Geolocation clientLocation = getGeolocationService().location(request.getClientIP());
		final List<CacheLocation> cacheLocations = orderCacheLocations(request,
				getCacheRegister().getCacheLocations(zoneId), ds, clientLocation);
		for (final CacheLocation location : cacheLocations) {
			if (LOGGER.isDebugEnabled()) {
				LOGGER.debug("Trying location: " + location.getId());
			}

			caches = getSupportingCaches(location.getCaches(), ds);
			if (caches.isEmpty()) {
				if (LOGGER.isDebugEnabled()) {
					LOGGER.debug("No online, supporting caches were found at location: " + location.getId());
				}
				return null;
			}

			cache = dispersion.getCache(consistentHash(caches, request.getPath()));
			LOGGER.warn("cache selected: " + cache.getId());

			Map<String, AtomicInteger> m = new HashMap<String, AtomicInteger>();
			long time = System.currentTimeMillis();
			for(String str : rstrs) {
				cache = dispersion.getCache(consistentHash(caches, str));
				AtomicInteger i = m.get(cache.getId());
				if(i == null) {
					i = new AtomicInteger(0);
					m.put(cache.getId(), i);
				}
				i.incrementAndGet();
			}
			time = System.currentTimeMillis() - time;
			LOGGER.warn(String.format("time: %d", time));
			for(String id : m.keySet()) {
				LOGGER.warn(String.format("cache(%s): %d", id, m.get(id).get()));
			}
			
			m = new HashMap<String, AtomicInteger>();
			time = System.currentTimeMillis();
			for(String str : rstrs) {
				cache = consistentHashOld(caches, str);
				AtomicInteger i = m.get(cache.getId());
				if(i == null) {
					i = new AtomicInteger(0);
					m.put(cache.getId(), i);
				}
				i.incrementAndGet();
			}
			time = System.currentTimeMillis() - time;
			LOGGER.warn(String.format("time: %d", time));
			for(String id : m.keySet()) {
				LOGGER.warn(String.format("cache(%s): %d", id, m.get(id).get()));
			}
			
			
			if (cache != null) {
				try {
					return new URL(ds.createURIString(request, cache));

				} catch (final MalformedURLException e) {
					LOGGER.error(e.getMessage(), e);
					throw new TrafficRouterException(URL_ERR_STR, e);
				}
			}
		}

		LOGGER.info(UNABLE_TO_ROUTE_REQUEST);
		throw new TrafficRouterException(UNABLE_TO_ROUTE_REQUEST);
	}

}
