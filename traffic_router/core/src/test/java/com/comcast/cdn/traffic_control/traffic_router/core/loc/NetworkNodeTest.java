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

package com.comcast.cdn.traffic_control.traffic_router.core.loc;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

import java.io.File;
import java.io.FileReader;
import java.net.InetAddress;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import org.apache.log4j.Logger;
import org.json.JSONArray;
import org.json.JSONObject;
import org.json.JSONTokener;
import org.junit.Before;
import org.junit.Test;

import com.google.common.net.InetAddresses;

public class NetworkNodeTest {
	private static final Logger LOGGER = Logger.getLogger(NetworkNodeTest.class);
	private Map<String, List<String>> netMap = new HashMap<String, List<String>>();
	private NetworkNode root;

	@Before
	public void setUp() throws Exception {
		final File file = new File(getClass().getClassLoader().getResource("czmap.json").toURI());
		root = NetworkNode.generateTree(file);

		final JSONObject json = new JSONObject(new JSONTokener(new FileReader(file)));
		final JSONObject coverageZones = json.getJSONObject("coverageZones");

		for (String loc : JSONObject.getNames(coverageZones)) {
			final JSONObject locData = coverageZones.getJSONObject(loc);

			for (String networkType : JSONObject.getNames(locData)) {
				final JSONArray networks = locData.getJSONArray(networkType);
				String network = networks.getString(0).split("/")[0];
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
}
