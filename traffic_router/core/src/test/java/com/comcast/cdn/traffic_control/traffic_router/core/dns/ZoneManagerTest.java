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

package com.comcast.cdn.traffic_control.traffic_router.core.dns;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertTrue;

import java.io.File;
import java.io.FileReader;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.Collection;
import java.util.HashMap;
import java.util.HashSet;
import java.util.Map;
import java.util.Set;

import org.apache.log4j.Logger;
import org.json.JSONArray;
import org.json.JSONObject;
import org.json.JSONTokener;
import org.junit.Before;
import org.junit.BeforeClass;
import org.junit.Test;
import org.springframework.context.ApplicationContext;
import org.xbill.DNS.Name;
import org.xbill.DNS.TextParseException;
import org.xbill.DNS.Type;
import org.xbill.DNS.Zone;

import com.comcast.cdn.traffic_control.traffic_router.core.TestBase;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.Cache;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.Cache.DeliveryServiceReference;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheLocation;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheRegister;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouter;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouterManager;
import com.google.common.net.InetAddresses;

public class ZoneManagerTest {
	private static final Logger LOGGER = Logger.getLogger(ZoneManagerTest.class);
	private static ApplicationContext context;
	private TrafficRouterManager trafficRouterManager;
	private String defaultDnsRoutingName;
	private Map<String, InetAddress> netMap = new HashMap<String, InetAddress>();
	private DNSAccessRecord.Builder builder;

	@BeforeClass
	public static void setUpBeforeClass() throws Exception {
		try {
			context = TestBase.getContext();
		} catch(Exception e) {
			e.printStackTrace();
		}
	}

	@Before
	public void setUp() throws Exception {
		trafficRouterManager = (TrafficRouterManager) context.getBean("trafficRouterManager");
		defaultDnsRoutingName = (String) context.getBean("staticZoneManagerDnsRoutingNameInitializer");
		final File file = new File(getClass().getClassLoader().getResource("czmap.json").toURI());
		final JSONObject json = new JSONObject(new JSONTokener(new FileReader(file)));
		final JSONObject coverageZones = json.getJSONObject("coverageZones");

		for (String loc : JSONObject.getNames(coverageZones)) {
			final JSONObject locData = coverageZones.getJSONObject(loc);
			final JSONArray networks = locData.getJSONArray("network");
			String network = networks.getString(0).split("/")[0];
			InetAddress ip = InetAddresses.forString(network);
			ip = InetAddresses.increment(ip);

			netMap.put(loc, ip);
		}

		builder = new DNSAccessRecord.Builder(1, InetAddress.getByName("192.168.12.34"));
	}


	@Test
	public void testDynamicZoneCache() throws TextParseException, UnknownHostException {
		TrafficRouter trafficRouter = trafficRouterManager.getTrafficRouter();
		CacheRegister cacheRegister = trafficRouter.getCacheRegister();
		Map<String, Collection<CacheLocation>> edgeLocations = new HashMap<String, Collection<CacheLocation>>();
		Set<String> dnsLimited = new HashSet<String>();

		for (Cache c : cacheRegister.getCacheMap().values()) {
			for (DeliveryServiceReference dsr : c.getDeliveryServices()) {
				final DeliveryService ds = cacheRegister.getDeliveryService(dsr.getDeliveryServiceId());

				if (!ds.isDns()) continue;
				final String edgeName = dsr.getFqdn() + ".";

				if (ds.getMaxDnsIps() > 0) {
					dnsLimited.add(edgeName);
				}

				if (!edgeLocations.containsKey(edgeName)) {
					edgeLocations.put(edgeName, new HashSet<CacheLocation>());
				}

				for (CacheLocation location : cacheRegister.getCacheLocations()) {
					if (ds.isLocationAvailable(location)) {
						final Collection<CacheLocation> locations = edgeLocations.get(edgeName);
						locations.add(location);
						edgeLocations.put(edgeName, locations);
					}
				}
			}
		}

		for (String name : edgeLocations.keySet()) {
			for (CacheLocation location : edgeLocations.get(name)) {
				// need to iterate through the CZF and submit a bunch of these into a job queue to run repeatedly/fast
				final InetAddress source = netMap.get(location.getId());
				final Name n = new Name(name);
				final Map<Zone, Integer> zoneTracker = new HashMap<Zone, Integer>();
				final Set<Zone> zones = new HashSet<Zone>();
				final int seenThreshold = 5;

				while (true) {
					final Zone zone = trafficRouter.getZone(n, Type.A, source, true, builder); // this should load the zone into the dynamicZoneCache
					assertNotNull(zone);

					if (!zoneTracker.containsKey(zone)) {
						zoneTracker.put(zone, 1);
					} else {
						final int count = zoneTracker.get(zone);
						zoneTracker.put(zone, count + 1);
					}

					if (!dnsLimited.contains(n.toString())) {
						break;
					}

					/*
					 * If we have limits on the number of records, continue building the cache
					 */
					boolean allSeen = true;

					for (Integer count : zoneTracker.values()) {
						if (count < seenThreshold) {
							allSeen = false;
							break;
						}
					}

					if (allSeen) {
						break;
					}
				}

				zones.addAll(zoneTracker.keySet());
				//System.out.println("Generated " + zones.size() + " zone(s) for " + n.toString() + " for cache group " + location.getId());

				for (int i = 0; i < 500; i++) {
					final Zone cachedZone = trafficRouter.getZone(n, Type.A, source, true, builder); // this should be a cache hit
					assertTrue(zones.contains(cachedZone));
				}
			}
		}
	}
}
