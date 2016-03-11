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

import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.Matchers.greaterThan;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

import static org.hamcrest.Matchers.greaterThanOrEqualTo;
import static org.hamcrest.MatcherAssert.assertThat;

import java.io.File;
import java.io.FileInputStream;
import java.io.FileReader;
import java.net.InetAddress;
import java.nio.file.Files;
import java.nio.file.StandardCopyOption;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Properties;

import com.comcast.cdn.traffic_control.traffic_router.core.util.IntegrationTest;
import org.apache.http.HttpResponse;
import org.apache.http.client.HttpClient;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.impl.client.DefaultHttpClient;
import org.apache.http.params.BasicHttpParams;
import org.apache.http.params.HttpConnectionParams;
import org.apache.http.params.HttpParams;
import org.apache.log4j.Logger;
import org.json.JSONArray;
import org.json.JSONObject;
import org.json.JSONTokener;
import org.junit.Before;
import org.junit.Test;

import com.google.common.net.InetAddresses;
import org.junit.experimental.categories.Category;

@Category(IntegrationTest.class)
public class NetworkNodeTest {
	private static final Logger LOGGER = Logger.getLogger(NetworkNodeTest.class);
	private Map<String, List<String>> netMap = new HashMap<String, List<String>>();
	private NetworkNode root;

	@Before
	public void setUp() throws Exception {
		File databasesDirectory = new File("src/test/db");
		if (!databasesDirectory.exists()) {
			assertThat(databasesDirectory.mkdirs(), equalTo(true));
		}

		File czmapFile = new File(databasesDirectory, "czmap.json");

		if (!czmapFile.exists()) {
			File crConfig = new File(databasesDirectory, "cr-config.json");

			if (!crConfig.exists()) {
				Properties trafficMonitorProperties = new Properties();
				trafficMonitorProperties.load(new FileInputStream("src/test/conf/traffic_monitor.properties"));
				String trafficMonitorHost = trafficMonitorProperties.getProperty("traffic_monitor.bootstrap.hosts").split(";")[0];

				HttpParams httpParams = new BasicHttpParams();
				HttpConnectionParams.setConnectionTimeout(httpParams, 5000);
				HttpConnectionParams.setSoTimeout(httpParams, 5000);
				HttpClient httpClient = new DefaultHttpClient(httpParams);

				HttpGet httpGet = new HttpGet("http://" + trafficMonitorHost + "/publish/CrConfig?json");
				HttpResponse httpResponse = httpClient.execute(httpGet);
				assertThat(httpResponse.getStatusLine().getStatusCode(), equalTo(200));
				long bytesCopied = Files.copy(httpResponse.getEntity().getContent(), crConfig.toPath(), StandardCopyOption.REPLACE_EXISTING);

				assertThat(bytesCopied, greaterThan(0L));
				httpResponse.getEntity().consumeContent();
			}

			JSONTokener jsonTokener = new JSONTokener(new FileReader(crConfig));
			String coverageZoneUrl = new JSONObject(jsonTokener).getJSONObject("config").getString("coveragezone.polling.url");

			HttpGet httpGet = new HttpGet(coverageZoneUrl);
			HttpClient httpClient = new DefaultHttpClient();
			HttpResponse httpResponse = httpClient.execute(httpGet);

			assertThat(httpResponse.getStatusLine().getStatusCode(), equalTo(200));
			long bytesCopied = Files.copy(httpResponse.getEntity().getContent(), czmapFile.toPath(), StandardCopyOption.REPLACE_EXISTING);
			assertThat(bytesCopied, greaterThan(0L));

			httpResponse.getEntity().consumeContent();
		}

		root = NetworkNode.generateTree(czmapFile);

		final JSONObject json = new JSONObject(new JSONTokener(new FileReader(czmapFile)));
		final JSONObject coverageZones = json.getJSONObject("coverageZones");

		for (String loc : JSONObject.getNames(coverageZones)) {
			final JSONObject locData = coverageZones.getJSONObject(loc);

			for (String networkType : JSONObject.getNames(locData)) {
				final JSONArray networks = locData.optJSONArray(networkType);
				if (networks == null) {
					continue;
				}
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

	@Test
	public void testNetworkNodePerformance() {
		final int iterations = 100000;
		final long startTime = System.currentTimeMillis();
		final long nnTPS = Long.parseLong(System.getProperty("nnTPS"));

		for (int i = 0; i < iterations; i++) {
			for (final String location : netMap.keySet()) {
				try {
					for (final String address : netMap.get(location)) {
						root.getNetwork(address);
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