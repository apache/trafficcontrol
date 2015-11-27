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
import java.util.List;
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
import com.comcast.cdn.traffic_control.traffic_router.core.cache.InetRecord;
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

		for (final DeliveryService ds : cacheRegister.getDeliveryServices().values()) {
			if (!ds.isDns()) continue;

			final JSONArray domains = ds.getDomains();

			for (int i = 0; i < domains.length(); i++) {
				final String domain = domains.optString(i);
				final Name edgeName = new Name(ZoneManager.getDnsRoutingName() + "." + domain + ".");

				for (CacheLocation location : cacheRegister.getCacheLocations()) {
					final List<Cache> caches = trafficRouter.selectCachesByCZ(ds, location);

					int p = 1;

					if (ds.getMaxDnsIps() > 0 && !trafficRouter.isConsistentDNSRouting() && caches.size() > ds.getMaxDnsIps()) {
						for (int c = caches.size(); c > (caches.size() - ds.getMaxDnsIps()); c--) {
							p *= c;
						}
					}

					final Set<Zone> zones = new HashSet<Zone>();
					final InetAddress source = netMap.get(location.getId());

					while (zones.size() != p) {
						final Zone zone = trafficRouter.getZone(edgeName, Type.A, source, true, builder); // this should load the zone into the dynamicZoneCache
						assertNotNull(zone);
						zones.add(zone);
					}

					for (int j = 0; j <= (p * 100); j++) {
						final Zone zone = trafficRouter.getZone(edgeName, Type.A, source, true, builder);
						assertTrue(zones.contains(zone));
					}
				}
			}
		}
	}
}