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

package com.comcast.cdn.traffic_control.traffic_router.core.external;

import com.comcast.cdn.traffic_control.traffic_router.core.util.ExternalTest;
import com.fasterxml.jackson.core.JsonFactory;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClientBuilder;
import org.junit.Before;
import org.junit.FixMethodOrder;
import org.junit.Test;
import org.junit.experimental.categories.Category;
import org.junit.runners.MethodSorters;

import java.io.IOException;
import java.io.InputStream;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.Iterator;
import java.util.List;
import java.util.Map;
import java.util.Random;

import static org.hamcrest.Matchers.endsWith;
import static org.hamcrest.Matchers.greaterThan;
import static org.hamcrest.Matchers.isIn;
import static org.hamcrest.Matchers.lessThan;
import static org.hamcrest.Matchers.nullValue;
import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.core.IsNot.not;
import static org.hamcrest.number.IsCloseTo.closeTo;
import static org.junit.Assert.assertThat;
import static org.junit.Assert.fail;

@Category(ExternalTest.class)
@FixMethodOrder(MethodSorters.NAME_ASCENDING)
public class SteeringTest {
	String steeringDeliveryServiceId;
	Map<String, String> targetDomains = new HashMap<String, String>();
	Map<String, Integer> targetWeights = new HashMap<String, Integer>();
	CloseableHttpClient httpClient;
	List<String> validLocations = new ArrayList<String>();
	String routerHttpPort = System.getProperty("routerHttpPort", "8888");
	String testHttpPort = System.getProperty("testHttpServerPort", "8889");

	JsonNode getJsonForResourcePath(String resourcePath) throws IOException {
		ObjectMapper objectMapper = new ObjectMapper(new JsonFactory());
		InputStream inputStream = getClass().getClassLoader().getResourceAsStream(resourcePath);

		if (inputStream == null) {
			fail("Could not find file '" + resourcePath + "' needed for test from the current classpath as a resource!");
		}

		return objectMapper.readTree(inputStream).get("response").get(0);
	}

	public String setupSteering(Map<String, String> domains, Map<String, Integer> weights, String resourcePath) throws IOException {
		domains.clear();
		weights.clear();

		JsonNode steeringNode = getJsonForResourcePath(resourcePath);

		Iterator<JsonNode> steeredDeliveryServices = steeringNode.get("targets").iterator();
		while (steeredDeliveryServices.hasNext()) {
			JsonNode steeredDeliveryService = steeredDeliveryServices.next();
			String targetId = steeredDeliveryService.get("deliveryService").asText();
			Integer targetWeight = steeredDeliveryService.get("weight").asInt();
			weights.put(targetId, targetWeight);
			domains.put(targetId, "");
		}

		return steeringNode.get("deliveryService").asText();
	}

	public void setupCrConfig() throws IOException {
		String resourcePath = "publish/CrConfig.json";
		InputStream inputStream = getClass().getClassLoader().getResourceAsStream(resourcePath);
		if (inputStream == null) {
			fail("Could not find file '" + resourcePath + "' needed for test from the current classpath as a resource!");
		}

		JsonNode jsonNode = new ObjectMapper(new JsonFactory()).readTree(inputStream);

		Iterator<String> deliveryServices = jsonNode.get("deliveryServices").fieldNames();
		while (deliveryServices.hasNext()) {
			String dsId = deliveryServices.next();
			if (targetDomains.containsKey(dsId)) {
				targetDomains.put(dsId, jsonNode.get("deliveryServices").get(dsId).get("domains").get(0).asText());
			}
		}

		assertThat(steeringDeliveryServiceId, not(nullValue()));
		assertThat(targetDomains.isEmpty(), equalTo(false));

		for (String deliveryServiceId : targetDomains.keySet()) {
			Iterator<String> cacheIds = jsonNode.get("contentServers").fieldNames();
			while (cacheIds.hasNext()) {
				String cacheId = cacheIds.next();
				JsonNode cacheNode = jsonNode.get("contentServers").get(cacheId);

				if (!cacheNode.has("deliveryServices")) {
					continue;
				}

				if (cacheNode.get("deliveryServices").has(deliveryServiceId)) {
					int port = cacheNode.get("port").asInt();
					String portText = (port == 80) ? "" : ":" + port;
					validLocations.add("http://" + cacheId + "." + targetDomains.get(deliveryServiceId) + portText + "/stuff?fakeClientIpAddress=12.34.56.78");
				}
			}
		}

		assertThat(validLocations.isEmpty(), equalTo(false));
	}

	@Before
	public void before() throws Exception {
		steeringDeliveryServiceId = setupSteering(targetDomains, targetWeights, "internal/api/1.2/steering.json");
		setupCrConfig();

		httpClient = HttpClientBuilder.create().disableRedirectHandling().build();
	}

	@Test
	public void itUsesSteeredDeliveryServiceIdInRedirect() throws Exception {
		HttpGet httpGet = new HttpGet("http://localhost:" + routerHttpPort + "/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "foo." + steeringDeliveryServiceId + ".bar");
		CloseableHttpResponse response = null;

		try {
			response = httpClient.execute(httpGet);
			assertThat("Failed getting 302 for request " + httpGet.getFirstHeader("Host").getValue(), response.getStatusLine().getStatusCode(), equalTo(302));
			assertThat(response.getFirstHeader("Location").getValue(), isIn(validLocations));
		} finally {
			if (response != null) { response.close(); }
		}
	}

	@Test
	public void itUsesTargetFiltersForSteering() throws Exception {
		HttpGet httpGet = new HttpGet("http://localhost:" + routerHttpPort + "/qwerytuiop/force-to-target-2/asdfghjkl?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "foo.steering-test-1.thecdn.example.com");
		CloseableHttpResponse response = null;

		try {
			response = httpClient.execute(httpGet);
			assertThat("Failed getting 302 for request " + httpGet.getFirstHeader("Host").getValue(), response.getStatusLine().getStatusCode(), equalTo(302));
			assertThat(response.getFirstHeader("Location").getValue(), endsWith(".steering-target-2.thecdn.example.com:8090/qwerytuiop/force-to-target-2/asdfghjkl?fakeClientIpAddress=12.34.56.78"));
		} finally {
			if (response != null) { response.close(); }
		}
	}

	@Test
	public void itUsesXtcSteeringOptionForOverride() throws Exception {
		HttpGet httpGet = new HttpGet("http://localhost:" + routerHttpPort + "/qwerytuiop/force-to-target-2/asdfghjkl?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "foo.steering-test-1.thecdn.example.com");
		httpGet.addHeader("X-TC-Steering-Option", "steering-target-1");

		CloseableHttpResponse response = null;

		try {
			response = httpClient.execute(httpGet);
			assertThat("Failed getting 302 for request " + httpGet.getFirstHeader("Host").getValue(), response.getStatusLine().getStatusCode(), equalTo(302));
			assertThat(response.getFirstHeader("Location").getValue(), endsWith(".steering-target-1.thecdn.example.com:8090/qwerytuiop/force-to-target-2/asdfghjkl?fakeClientIpAddress=12.34.56.78"));
		} finally {
			if (response != null) { response.close(); }
		}
	}

	@Test
	public void itReturns503ForBadDeliveryServiceInXtcSteeringOption() throws Exception {
		HttpGet httpGet = new HttpGet("http://localhost:" + routerHttpPort + "/qwerytuiop/asdfghjkl?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "foo.steering-test-1.thecdn.example.com");
		httpGet.addHeader("X-TC-Steering-Option", "ds-02");
		CloseableHttpResponse response = null;

		try {
			response = httpClient.execute(httpGet);
			assertThat(response.getStatusLine().getStatusCode(), equalTo(503));
		} finally {
			if (response != null) { response.close(); }
		}
	}

	@Test
	public void itUsesWeightedDistributionForRequestPath() throws Exception {
		int count = 0;
		for (int weight : targetWeights.values()) {
			count += weight;
		}

		count *= 1000;

		if (count > 100000) {
			count = 100000;
		}

		Map<String, Integer> results = new HashMap<String, Integer>();
		for (String steeredId : targetWeights.keySet()) {
			results.put(steeredId, 0);
		}

		System.out.println("Going to execute " + count + " requests through steering delivery service '" + steeringDeliveryServiceId + "'");

		for (int i = 0; i < count; i++) {
			String path = generateRandomPath();
			HttpGet httpGet = new HttpGet("http://localhost:" + routerHttpPort + path + "?fakeClientIpAddress=12.34.56.78");
			httpGet.addHeader("Host", "foo." + steeringDeliveryServiceId + ".bar");
			CloseableHttpResponse response = null;

			try {
				response = httpClient.execute(httpGet);
				assertThat("Did not get 302 for request '" + httpGet.getURI() + "'", response.getStatusLine().getStatusCode(), equalTo(302));
				String location = response.getFirstHeader("Location").getValue();

				for (String id : results.keySet()) {
					if (location.contains(id)) {
						results.put(id, results.get(id) + 1);
					}
				}
			} finally {
				if (response != null) { response.close(); }
			}
		}

		double totalWeight = 0;
		for (int weight : targetWeights.values()) {
			totalWeight += weight;
		}

		Map<String, Double> expectedHitRates = new HashMap<String, Double>();
		for (String id : targetWeights.keySet()) {
			expectedHitRates.put(id, targetWeights.get(id) / totalWeight);
		}

		for (String id : results.keySet()) {
			int hits = results.get(id);
			double hitRate = (double) hits / count;
			assertThat(hitRate, closeTo(expectedHitRates.get(id), 0.009));
		}
	}

	@Test
	public void z_itemsMigrateFromSmallerToLargerBucket() throws Exception {
		Map<String, String> domains = new HashMap<>();
		Map<String, Integer> weights = new HashMap<>();

		setupSteering(domains, weights, "internal/api/1.2/steering2.json");

		List<String> randomPaths = new ArrayList<>();

		for (int i = 0; i < 10000; i++) {
			randomPaths.add(generateRandomPath());
		}


		String smallerTarget = null;
		String largerTarget = null;
		for (String target : weights.keySet()) {
			if (smallerTarget == null && largerTarget == null) {
				smallerTarget = target;
				largerTarget = target;
			}

			if (weights.get(smallerTarget) > weights.get(target)) {
				smallerTarget = target;
			}

			if (weights.get(largerTarget) < weights.get(target)) {
				largerTarget = target;
			}
		}

		Map<String, List<String>> hashedPaths = new HashMap<>();
		hashedPaths.put(smallerTarget, new ArrayList<String>());
		hashedPaths.put(largerTarget, new ArrayList<String>());

		for (String path : randomPaths) {
			HttpGet httpGet = new HttpGet("http://localhost:" + routerHttpPort + path + "?fakeClientIpAddress=12.34.56.78");
			httpGet.addHeader("Host", "foo." + steeringDeliveryServiceId + ".bar");
			CloseableHttpResponse response = null;

			try {
				response = httpClient.execute(httpGet);
				assertThat("Did not get 302 for request '" + httpGet.getURI() + "'", response.getStatusLine().getStatusCode(), equalTo(302));
				String location = response.getFirstHeader("Location").getValue();

				for (String targetXmlId : hashedPaths.keySet()) {
					if (location.contains(targetXmlId)) {
						hashedPaths.get(targetXmlId).add(path);
					}
				}
			} finally {
				if (response != null) { response.close(); }
			}
		}

		// Change the steering attributes
		HttpPost httpPost = new HttpPost("http://localhost:" + testHttpPort + "/steering");
		httpClient.execute(httpPost).close();

		// steering is checked every 15 seconds by default.
		Thread.sleep(30 * 1000);

		Map<String, List<String>> rehashedPaths = new HashMap<>();
		rehashedPaths.put(smallerTarget, new ArrayList<String>());
		rehashedPaths.put(largerTarget, new ArrayList<String>());

		for (String path : randomPaths) {
			HttpGet httpGet = new HttpGet("http://localhost:" + routerHttpPort + path + "?fakeClientIpAddress=12.34.56.78");
			httpGet.addHeader("Host", "foo." + steeringDeliveryServiceId + ".bar");
			CloseableHttpResponse response = null;

			try {
				response = httpClient.execute(httpGet);
				assertThat("Did not get 302 for request '" + httpGet.getURI() + "'", response.getStatusLine().getStatusCode(), equalTo(302));
				String location = response.getFirstHeader("Location").getValue();

				for (String targetXmlId : rehashedPaths.keySet()) {
					if (location.contains(targetXmlId)) {
						rehashedPaths.get(targetXmlId).add(path);
					}
				}
			} finally {
				if (response != null) { response.close(); }
			}
		}

		assertThat(rehashedPaths.get(smallerTarget).size(), greaterThan(hashedPaths.get(smallerTarget).size()));
		assertThat(rehashedPaths.get(largerTarget).size(), lessThan(hashedPaths.get(largerTarget).size()));

		for (String path : hashedPaths.get(smallerTarget)) {
			assertThat(rehashedPaths.get(smallerTarget).contains(path), equalTo(true));
			assertThat(rehashedPaths.get(largerTarget).contains(path), equalTo(false));
		}
	}

	String alphanumericCharacters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWZYZ";
	String exampleValidPathCharacters = alphanumericCharacters + "/=;()-.";

	Random random = new Random(1462307930227L);
	String generateRandomPath() {
		int pathLength = 60 + random.nextInt(61);

		StringBuilder stringBuilder = new StringBuilder("/");
		for (int i = 0; i < 4; i++) {
			int index = random.nextInt(alphanumericCharacters.length());
			stringBuilder.append(alphanumericCharacters.charAt(index));
		}

		stringBuilder.append("/");

		for (int i = 0; i < pathLength; i++) {
			int index = random.nextInt(exampleValidPathCharacters.length());
			stringBuilder.append(exampleValidPathCharacters.charAt(index));
		}

		return stringBuilder.toString();
	}
}
