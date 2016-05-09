/*
 * Copyright 2016 Comcast Cable Communications Management, LLC
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
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClientBuilder;
import org.junit.Before;
import org.junit.Test;
import org.junit.experimental.categories.Category;

import java.io.InputStream;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.Iterator;
import java.util.List;
import java.util.Map;
import java.util.Random;

import static org.hamcrest.Matchers.endsWith;
import static org.hamcrest.Matchers.isIn;
import static org.hamcrest.Matchers.nullValue;
import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.core.IsNot.not;
import static org.hamcrest.number.IsCloseTo.closeTo;
import static org.junit.Assert.assertThat;
import static org.junit.Assert.fail;

@Category(ExternalTest.class)
public class SteeringTest {
	String steeringDeliveryServiceId;
	Map<String, String> targetDomains = new HashMap<String, String>();
	Map<String, Integer> targetWeights = new HashMap<String, Integer>();
	CloseableHttpClient httpClient;
	List<String> validLocations = new ArrayList<String>();

	@Before
	public void before() throws Exception {
		ObjectMapper objectMapper = new ObjectMapper(new JsonFactory());

		String resourcePath = "internal/api/1.2/steering.json";
		InputStream inputStream = getClass().getClassLoader().getResourceAsStream(resourcePath);

		if (inputStream == null) {
			fail("Could not find file '" + resourcePath + "' needed for test from the current classpath as a resource!");
		}

		JsonNode steeringNode = objectMapper.readTree(inputStream).get("response").get(0);

		steeringDeliveryServiceId = steeringNode.get("deliveryService").asText();
		Iterator<JsonNode> steeredDeliveryServices = steeringNode.get("targets").iterator();
		while (steeredDeliveryServices.hasNext()) {
			JsonNode steeredDeliveryService = steeredDeliveryServices.next();
			String targetId = steeredDeliveryService.get("deliveryService").asText();
			Integer targetWeight = steeredDeliveryService.get("weight").asInt();
			targetWeights.put(targetId, targetWeight);
			targetDomains.put(targetId, "");
		}

		resourcePath = "publish/CrConfig.json";
		inputStream = getClass().getClassLoader().getResourceAsStream(resourcePath);
		if (inputStream == null) {
			fail("Could not find file '" + resourcePath + "' needed for test from the current classpath as a resource!");
		}

		JsonNode jsonNode = objectMapper.readTree(inputStream);

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

		httpClient = HttpClientBuilder.create().disableRedirectHandling().build();
	}

	@Test
	public void itUsesSteeredDeliveryServiceIdInRedirect() throws Exception {
		HttpGet httpGet = new HttpGet("http://localhost:8888/stuff?fakeClientIpAddress=12.34.56.78");
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
		HttpGet httpGet = new HttpGet("http://localhost:8888/qwerytuiop/force-to-eight/asdfghjkl?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "foo.mm-test.thecdn.cdn.example.com");
		CloseableHttpResponse response = null;

		try {
			response = httpClient.execute(httpGet);
			assertThat("Failed getting 302 for request " + httpGet.getFirstHeader("Host").getValue(), response.getStatusLine().getStatusCode(), equalTo(302));
			assertThat(response.getFirstHeader("Location").getValue(), endsWith(".ds-08.thecdn.cdn.example.com:8090/qwerytuiop/force-to-eight/asdfghjkl?fakeClientIpAddress=12.34.56.78"));
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
			HttpGet httpGet = new HttpGet("http://localhost:8888" + path + "?fakeClientIpAddress=12.34.56.78");
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
