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

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.core.IsCollectionContaining.hasItem;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertTrue;

import java.io.File;
import java.io.FileReader;
import java.math.BigInteger;
import java.net.InetAddress;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;

import com.comcast.cdn.traffic_control.traffic_router.core.util.IntegrationTest;
import org.json.JSONArray;
import org.json.JSONObject;
import org.json.JSONTokener;
import org.junit.Before;
import org.junit.BeforeClass;
import org.junit.Test;
import org.junit.experimental.categories.Category;
import org.springframework.context.ApplicationContext;
import org.xbill.DNS.Name;
import org.xbill.DNS.TextParseException;
import org.xbill.DNS.Type;
import org.xbill.DNS.Zone;

import com.comcast.cdn.traffic_control.traffic_router.core.TestBase;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.Cache;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheLocation;
import com.comcast.cdn.traffic_control.traffic_router.core.cache.CacheRegister;
import com.comcast.cdn.traffic_control.traffic_router.core.ds.DeliveryService;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouter;
import com.comcast.cdn.traffic_control.traffic_router.core.router.TrafficRouterManager;
import com.google.common.cache.CacheStats;
import com.google.common.net.InetAddresses;

@Category(IntegrationTest.class)
public class ZoneManagerTest {
	private static ApplicationContext context;
	private TrafficRouterManager trafficRouterManager;
	private Map<String, InetAddress> netMap = new HashMap<String, InetAddress>();

	@BeforeClass
	public static void setUpBeforeClass() throws Exception {
		assertThat("Copy core/src/main/conf/traffic_monitor.properties to core/src/test/conf and set 'traffic_monitor.bootstrap.hosts' to a real traffic monitor", Files.exists(Paths.get(TestBase.monitorPropertiesPath)), equalTo(true));
		context = TestBase.getContext();
	}

	@Before
	public void setUp() throws Exception {
		trafficRouterManager = (TrafficRouterManager) context.getBean("trafficRouterManager");
		trafficRouterManager.getTrafficRouter().setApplicationContext(context);
		final File file = new File("src/test/resources/czmap.json");
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

	}

	@Test
	public void testDynamicZoneCache() throws TextParseException {
		TrafficRouter trafficRouter = trafficRouterManager.getTrafficRouter();
		CacheRegister cacheRegister = trafficRouter.getCacheRegister();
		ZoneManager zoneManager = trafficRouter.getZoneManager();

		for (final DeliveryService ds : cacheRegister.getDeliveryServices().values()) {
			if (!ds.isDns()) {
				continue;
			}

			final JSONArray domains = ds.getDomains();

			for (int i = 0; i < domains.length(); i++) {
				final String domain = domains.optString(i);
				final Name edgeName = new Name(ZoneManager.getDnsRoutingName() + "." + domain + ".");

				for (CacheLocation location : cacheRegister.getCacheLocations()) {
					final List<Cache> caches = trafficRouter.selectCachesByCZ(ds, location);

					if (caches == null) {
						continue;
					}

					final InetAddress source = netMap.get(location.getId());
					final DNSAccessRecord.Builder builder = new DNSAccessRecord.Builder(1, source);
					final Set<Zone> zones = new HashSet<Zone>();
					final int maxDnsIps = ds.getMaxDnsIps();
					long combinations = 1;

					if (maxDnsIps > 0 && !trafficRouter.isConsistentDNSRouting() && caches.size() > maxDnsIps) {
						final BigInteger top = fact(caches.size());
						final BigInteger f = fact(caches.size() - maxDnsIps);
						final BigInteger s = fact(maxDnsIps);

						combinations = top.divide(f.multiply(s)).longValue();
						int c = 0;

						while (c < (combinations * 100)) {
							final Zone zone = trafficRouter.getZone(edgeName, Type.A, source, true, builder); // this should load the zone into the dynamicZoneCache if not already there
							assertNotNull(zone);
							zones.add(zone);
							c++;
						}
					}

					final CacheStats cacheStats = zoneManager.getDynamicCacheStats();

					for (int j = 0; j <= (combinations * 100); j++) {
						final long missCount = new Long(cacheStats.missCount());
						final Zone zone = trafficRouter.getZone(edgeName, Type.A, source, true, builder);
						assertNotNull(zone);
						assertEquals(missCount, cacheStats.missCount()); // should always be a cache hit so these should remain the same

						if (!zones.isEmpty()) {
							assertThat(zones, hasItem(zone));
							assertTrue(zones.contains(zone));
						}
					}
				}
			}
		}
	}

	private BigInteger fact(final int n) {
		BigInteger p = new BigInteger("1");

		for (long c = n; c > 0; c--) {
			p = p.multiply(BigInteger.valueOf(c));
		}

		return p;
	}
}