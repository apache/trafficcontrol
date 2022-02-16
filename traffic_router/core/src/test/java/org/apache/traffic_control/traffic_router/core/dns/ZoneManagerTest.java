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

package org.apache.traffic_control.traffic_router.core.dns;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.core.IsCollectionContaining.hasItem;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertTrue;

import java.io.File;
import java.math.BigInteger;
import java.net.InetAddress;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.Iterator;

import org.apache.traffic_control.traffic_router.core.util.IntegrationTest;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.AfterClass;
import org.junit.Before;
import org.junit.BeforeClass;
import org.junit.Test;
import org.junit.experimental.categories.Category;
import org.springframework.context.ApplicationContext;
import org.xbill.DNS.Name;
import org.xbill.DNS.TextParseException;
import org.xbill.DNS.Type;
import org.xbill.DNS.Zone;

import org.apache.traffic_control.traffic_router.core.TestBase;
import org.apache.traffic_control.traffic_router.core.edge.Cache;
import org.apache.traffic_control.traffic_router.core.edge.CacheLocation;
import org.apache.traffic_control.traffic_router.core.edge.CacheRegister;
import org.apache.traffic_control.traffic_router.core.edge.Node.IPVersions;
import org.apache.traffic_control.traffic_router.core.ds.DeliveryService;
import org.apache.traffic_control.traffic_router.core.router.TrafficRouter;
import org.apache.traffic_control.traffic_router.core.router.TrafficRouterManager;
import com.google.common.cache.CacheStats;
import com.google.common.net.InetAddresses;

@Category(IntegrationTest.class)
public class ZoneManagerTest {
	private static ApplicationContext context;
	private TrafficRouterManager trafficRouterManager;
	private Map<String, InetAddress> netMap = new HashMap<String, InetAddress>();

	@BeforeClass
	public static void setUpBeforeClass() throws Exception {
		TestBase.setupFakeServers();
		context = TestBase.getContext();
	}

	@Before
	public void setUp() throws Exception {
		trafficRouterManager = (TrafficRouterManager) context.getBean("trafficRouterManager");
		trafficRouterManager.getTrafficRouter().setApplicationContext(context);
		final File file = new File("src/test/resources/czmap.json");
		final ObjectMapper mapper = new ObjectMapper();
		final JsonNode jsonNode = mapper.readTree(file);
		final JsonNode coverageZones = jsonNode.get("coverageZones");

		final Iterator<String> czIter = coverageZones.fieldNames();
		while (czIter.hasNext()) {
			final String loc = czIter.next();
			final JsonNode locData = coverageZones.get(loc);
			final JsonNode networks = locData.get("network");
			String network = networks.get(0).asText().split("/")[0];
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

			final String domain = ds.getDomain();

			final Name edgeName = new Name(ds.getRoutingName() + "." + domain + ".");

			for (InetAddress source : netMap.values()) {
				final CacheLocation location = trafficRouter.getCoverageZoneCacheLocation(source.getHostAddress(), ds, IPVersions.IPV4ONLY);
				final List<Cache> caches = trafficRouter.selectCachesByCZ(ds, location, IPVersions.IPV4ONLY);

				if (caches == null) {
					continue;
				}

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

	@AfterClass
	public static void tearDown() throws Exception {
		TestBase.tearDownFakeServers();
	}

	private BigInteger fact(final int n) {
		BigInteger p = new BigInteger("1");

		for (long c = n; c > 0; c--) {
			p = p.multiply(BigInteger.valueOf(c));
		}

		return p;
	}
}
