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

package com.comcast.cdn.traffic_control.traffic_router.core.loc;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

import static org.hamcrest.Matchers.greaterThanOrEqualTo;
import static org.hamcrest.MatcherAssert.assertThat;

import java.io.File;
import java.net.InetAddress;
import java.util.*;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.log4j.Logger;
import org.junit.Before;
import org.junit.Test;

import com.google.common.net.InetAddresses;
import org.junit.experimental.categories.Category;

public class NetworkNodeTest {
	private static final Logger LOGGER = Logger.getLogger(NetworkNodeTest.class);
	private Map<String, List<String>> netMap = new HashMap<String, List<String>>();
	private NetworkNode root;

	@Before
	public void setUp() throws Exception {
		final File file = new File(getClass().getClassLoader().getResource("czmap.json").toURI());
		root = NetworkNode.generateTree(file, false);
		final ObjectMapper mapper = new ObjectMapper();
		final JsonNode jsonNode = mapper.readTree(file);
		final JsonNode coverageZones = jsonNode.get("coverageZones");

		final Iterator<String> networkIter = coverageZones.fieldNames();
		while (networkIter.hasNext()) {
			final String loc = networkIter.next();
			final JsonNode locData = coverageZones.get(loc);
			for (final JsonNode networkType : locData) {
				final String network = networkType.get(0).asText().split("/")[0];
				InetAddress ip = InetAddresses.forString(network);
				ip = InetAddresses.increment(ip);

				if (!netMap.containsKey(loc)) {
					netMap.put(loc, new ArrayList<String>());
				}

				final List<String> addressList = netMap.get(loc);
				addressList.add(InetAddresses.toAddrString(ip));

				netMap.put(loc, addressList);

			}
		}
	}

	@Test
	public void testIps() {
		try {
			for (String location : netMap.keySet()) {
				for (String address : netMap.get(location)) {
					final NetworkNode nn = root.getNetwork(address);
					assertNotNull(nn);
					final String loc = nn.getLoc();
					assertEquals(loc, location);
					LOGGER.info(String.format("result for ip=%s: %s", address, loc));
				}
			}
		} catch (Exception e) {
			e.printStackTrace();
		}
	}

	@Test
	public void testNetworkNodePerformance() {
		final int iterations = 100000;
		final long startTime = System.currentTimeMillis();
		final long nnTPS = Long.parseLong(System.getProperty("nnTPS", "12000"));

		for (int i = 0; i < iterations; i++) {
			for (final String location : netMap.keySet()) {
				try {
					for (final String address : netMap.get(location)) {
						final NetworkNode nn = root.getNetwork(address);
					}
				} catch (final NetworkNodeException e) {
					e.printStackTrace();
				}
			}
		}

		final long runTime = System.currentTimeMillis() - startTime;
		final long tps = (iterations / runTime) * 1000;

		assertThat(tps, greaterThanOrEqualTo(nnTPS));
	}
}