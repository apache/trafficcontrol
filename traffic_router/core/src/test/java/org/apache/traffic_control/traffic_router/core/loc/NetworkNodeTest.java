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

package org.apache.traffic_control.traffic_router.core.loc;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

import static org.hamcrest.Matchers.greaterThanOrEqualTo;
import static org.hamcrest.MatcherAssert.assertThat;

import java.io.File;
import java.net.InetAddress;

import java.util.HashMap;
import java.util.Iterator;
import java.util.List;
import java.util.Map;
import java.util.ArrayList;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.junit.Before;
import org.junit.Test;

import com.google.common.net.InetAddresses;

public class NetworkNodeTest {
	private static final Logger LOGGER = LogManager.getLogger(NetworkNodeTest.class);
	private Map<String, List<String>> netMap = new HashMap<String, List<String>>();
	private Map<String, List<String>> deepNetMap = new HashMap<String, List<String>>();
	private NetworkNode root;
	private NetworkNode deepRoot;

	@Before
	public void setUp() throws Exception {
		root = setUp("czmap.json", false);
		deepRoot = setUp("dczmap.json", true);
	}

	private NetworkNode setUp(final String filename, final boolean useDeep) throws Exception {
		final Map<String, List<String>> testNetMap = useDeep ? deepNetMap : netMap;
		final File file = new File(getClass().getClassLoader().getResource(filename).toURI());
		final NetworkNode nn = NetworkNode.generateTree(file, false, useDeep);
		final ObjectMapper mapper = new ObjectMapper();
		final JsonNode jsonNode = mapper.readTree(file);
		final String czKey = useDeep ? "deepCoverageZones" : "coverageZones";
		final JsonNode coverageZones = jsonNode.get(czKey);

		final Iterator<String> networkIter = coverageZones.fieldNames();
		while (networkIter.hasNext()) {
		    final String loc = networkIter.next();
			final JsonNode locData = coverageZones.get(loc);
			for (String networkType : new String[]{"network", "network6"}) {
				final JsonNode networks = locData.get(networkType);
				final String network = networks.get(0).asText().split("/")[0];
				InetAddress ip = InetAddresses.forString(network);
				ip = InetAddresses.increment(ip);

				if (!testNetMap.containsKey(loc)) {
					testNetMap.put(loc, new ArrayList<String>());
				}

				final List<String> addressList = testNetMap.get(loc);
				addressList.add(InetAddresses.toAddrString(ip));

				testNetMap.put(loc, addressList);
			}
		}
		return nn;
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
		testNetworkNodePerformance(root, netMap);
	}

	@Test
	public void testDeepNetworkNodePerformance() {
		testNetworkNodePerformance(deepRoot, deepNetMap);
	}

	private void testNetworkNodePerformance(NetworkNode testRoot, Map<String, List<String>> testNetMap) {
		final int iterations = 100000;
		final long startTime = System.currentTimeMillis();
		final long nnTPS = Long.parseLong(System.getProperty("nnTPS", "12000"));

		for (int i = 0; i < iterations; i++) {
			for (final String location : testNetMap.keySet()) {
				try {
					for (final String address : testNetMap.get(location)) {
						final NetworkNode nn = testRoot.getNetwork(address);
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