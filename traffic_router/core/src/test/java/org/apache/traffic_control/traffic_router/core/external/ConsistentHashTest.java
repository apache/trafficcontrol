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

package org.apache.traffic_control.traffic_router.core.external;

import org.apache.traffic_control.traffic_router.core.util.CidrAddress;
import org.apache.traffic_control.traffic_router.core.util.ExternalTest;
import org.apache.traffic_control.traffic_router.core.util.JsonUtils;
import org.apache.traffic_control.traffic_router.core.util.TrafficOpsUtils;
import com.fasterxml.jackson.core.JsonFactory;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClientBuilder;
import org.apache.http.util.EntityUtils;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;
import org.junit.experimental.categories.Category;

import java.io.InputStream;
import java.net.URLEncoder;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;

import static org.hamcrest.CoreMatchers.nullValue;
import static org.hamcrest.Matchers.isIn;
import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.core.IsNot.not;
import static org.hamcrest.number.OrderingComparison.greaterThan;
import static org.junit.Assert.assertThat;
import static org.junit.Assert.fail;

@Category(ExternalTest.class)
public class ConsistentHashTest {
	private CloseableHttpClient closeableHttpClient;
	String deliveryServiceId;
	String ipAddressInCoverageZone;
	String steeringDeliveryServiceId;
	String consistentHashRegex;
	List<String> steeredDeliveryServices = new ArrayList<String>();

	@Before
	public void before() throws Exception {
		closeableHttpClient = HttpClientBuilder.create().build();

		String resourcePath = "api/"+TrafficOpsUtils.TO_API_VERSION+"/steering";
		InputStream inputStream = getClass().getClassLoader().getResourceAsStream(resourcePath);

		if (inputStream == null) {
			fail("Could not find file '" + resourcePath + "' needed for test from the current classpath as a resource!");
		}

		ObjectMapper objectMapper = new ObjectMapper(new JsonFactory());
		JsonNode steeringNode = objectMapper.readTree(inputStream).get("response").get(0);

		steeringDeliveryServiceId = steeringNode.get("deliveryService").asText();
		Iterator<JsonNode> iterator = steeringNode.get("targets").iterator();
		while (iterator.hasNext()) {
			JsonNode target = iterator.next();
			steeredDeliveryServices.add(target.get("deliveryService").asText());
		}

		resourcePath = "publish/CrConfig.json";
		inputStream = getClass().getClassLoader().getResourceAsStream(resourcePath);
		if (inputStream == null) {
			fail("Could not find file '" + resourcePath + "' needed for test from the current classpath as a resource!");
		}

		JsonNode jsonNode = objectMapper.readTree(inputStream);

		deliveryServiceId = null;

		Iterator<String> deliveryServices = jsonNode.get("deliveryServices").fieldNames();
		while (deliveryServices.hasNext() && deliveryServiceId == null) {
			String dsId = deliveryServices.next();

			JsonNode deliveryServiceNode = jsonNode.get("deliveryServices").get(dsId);

			if (deliveryServiceNode.has("steeredDeliveryServices")) {
				continue;
			}

			JsonNode dispersionNode = deliveryServiceNode.get("dispersion");

			if (dispersionNode == null || dispersionNode.get("limit").asInt() != 1 && dispersionNode.get("shuffled").asText().equals("true")) {
				continue;
			}

			Iterator<JsonNode> matchsets = deliveryServiceNode.get("matchsets").iterator();
			while (matchsets.hasNext() && deliveryServiceId == null) {
				if ("HTTP".equals(matchsets.next().get("protocol").asText())) {
					if (deliveryServiceNode.has("consistentHashRegex")) {
						deliveryServiceId = dsId;
						consistentHashRegex = deliveryServiceNode.get("consistentHashRegex").asText();
					}
				}
			}

			if (deliveryServiceId == null) {
				System.out.println("Skipping " + deliveryServiceId + " no http protocol matchset");
			}
		}

		assertThat(deliveryServiceId, not(nullValue()));
		assertThat(steeringDeliveryServiceId, not(nullValue()));
		assertThat(steeredDeliveryServices.isEmpty(), equalTo(false));

		resourcePath = "czf.json";
		inputStream = getClass().getClassLoader().getResourceAsStream(resourcePath);
		if (inputStream == null) {
			fail("Could not find file '" + resourcePath + "' needed for test from the current classpath as a resource!");
		}

		jsonNode = objectMapper.readTree(inputStream);

		JsonNode network = jsonNode.get("coverageZones").get(jsonNode.get("coverageZones").fieldNames().next()).get("network");

		for (int i = 0; i < network.size(); i++) {
			String cidrString = network.get(i).asText();
			CidrAddress cidrAddress = CidrAddress.fromString(cidrString);
			if (cidrAddress.getNetmaskLength() == 24) {
				byte[] hostBytes = cidrAddress.getHostBytes();
				ipAddressInCoverageZone = String.format("%d.%d.%d.123", hostBytes[0], hostBytes[1], hostBytes[2]);
				break;
			}
		}

		assertThat(ipAddressInCoverageZone.length(), greaterThan(0));
	}

	@After
	public void after() throws Exception {
		if (closeableHttpClient != null) closeableHttpClient.close();
	}

	@Test
	public void itAppliesConsistentHashingToRequestsForCoverageZone() throws Exception {
		CloseableHttpResponse response = null;

		try {
			String requestPath = URLEncoder.encode("/some/path/thing", "UTF-8");
			HttpGet httpGet = new HttpGet("http://localhost:3333/crs/consistenthash/cache/coveragezone?ip=" + ipAddressInCoverageZone + "&deliveryServiceId=" + deliveryServiceId + "&requestPath=" + requestPath);

			response = closeableHttpClient.execute(httpGet);

			assertThat("Expected to find " + ipAddressInCoverageZone + " in coverage zone using delivery service id " + deliveryServiceId, response.getStatusLine().getStatusCode(), equalTo(200));

			ObjectMapper objectMapper = new ObjectMapper(new JsonFactory());
			JsonNode cacheNode = objectMapper.readTree(EntityUtils.toString(response.getEntity()));

			String cacheId = cacheNode.get("id").asText();
			assertThat(cacheId, not(equalTo("")));

			response.close();

			response = closeableHttpClient.execute(httpGet);
			cacheNode = objectMapper.readTree(EntityUtils.toString(response.getEntity()));
			assertThat(cacheNode.get("id").asText(), equalTo(cacheId));

			response.close();

			requestPath = URLEncoder.encode("/another/different/path", "UTF-8");
			httpGet = new HttpGet("http://localhost:3333/crs/consistenthash/cache/coveragezone?ip=" + ipAddressInCoverageZone + "&deliveryServiceId=" + deliveryServiceId + "&requestPath=" + requestPath);

			response = closeableHttpClient.execute(httpGet);
			cacheNode = objectMapper.readTree(EntityUtils.toString(response.getEntity()));
			assertThat(cacheNode.get("id").asText(), not(equalTo(cacheId)));
			assertThat(cacheNode.get("id").asText(), not(equalTo("")));
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itAppliesConsistentHashingForRequestsOutsideCoverageZone() throws Exception {
		CloseableHttpResponse response = null;

		try {
			String requestPath = URLEncoder.encode("/some/path/thing", "UTF-8");
			HttpGet httpGet = new HttpGet("http://localhost:3333/crs/consistenthash/cache/geolocation?ip=8.8.8.8&deliveryServiceId=" + deliveryServiceId + "&requestPath=" + requestPath);

			response = closeableHttpClient.execute(httpGet);

			assertThat(response.getStatusLine().getStatusCode(), equalTo(200));

			ObjectMapper objectMapper = new ObjectMapper(new JsonFactory());
			JsonNode cacheNode = objectMapper.readTree(EntityUtils.toString(response.getEntity()));

			String cacheId = cacheNode.get("id").asText();
			assertThat(cacheId, not(equalTo("")));

			response.close();

			response = closeableHttpClient.execute(httpGet);
			cacheNode = objectMapper.readTree(EntityUtils.toString(response.getEntity()));
			assertThat(cacheNode.get("id").asText(), equalTo(cacheId));

			response.close();

			requestPath = URLEncoder.encode("/another/different/path", "UTF-8");
			httpGet = new HttpGet("http://localhost:3333/crs/consistenthash/cache/geolocation?ip=8.8.8.8&deliveryServiceId=" + deliveryServiceId + "&requestPath=" + requestPath);

			response = closeableHttpClient.execute(httpGet);
			cacheNode = objectMapper.readTree(EntityUtils.toString(response.getEntity()));
			assertThat(cacheNode.get("id").asText(), not(equalTo(cacheId)));
			assertThat(cacheNode.get("id").asText(), not(equalTo("")));
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itAppliesConsistentHashingToSteeringDeliveryService() throws Exception {
		CloseableHttpResponse response = null;
		try {
			String requestPath = URLEncoder.encode("/some/path/thing", "UTF-8");
			HttpGet httpGet = new HttpGet("http://localhost:3333/crs/consistenthash/deliveryservice?ip=98.76.54.123&deliveryServiceId=" + steeringDeliveryServiceId + "&requestPath=" + requestPath);
			response = closeableHttpClient.execute(httpGet);

			ObjectMapper objectMapper = new ObjectMapper(new JsonFactory());
			JsonNode deliveryServiceNode = objectMapper.readTree(EntityUtils.toString(response.getEntity()));
			assertThat(deliveryServiceNode.get("id").asText(), isIn(steeredDeliveryServices));

		} finally {
			if (response != null) response.close();
		}
	}
	
	@Test
	public void itUsesBypassFiltersWithDeliveryServiceSteering() throws Exception {
		CloseableHttpResponse response = null;
		try {
			String requestPath = URLEncoder.encode("/some/path/force-to-target-2/more/asdfasdf", "UTF-8");
			HttpGet httpGet = new HttpGet("http://localhost:3333/crs/consistenthash/deliveryservice?ip=98.76.54.123&deliveryServiceId=" + steeringDeliveryServiceId + "&requestPath=" + requestPath);
			response = closeableHttpClient.execute(httpGet);

			ObjectMapper objectMapper = new ObjectMapper(new JsonFactory());
			JsonNode deliveryServiceNode = objectMapper.readTree(EntityUtils.toString(response.getEntity()));
			assertThat(deliveryServiceNode.get("id").asText(), equalTo("steering-target-2"));

			requestPath = URLEncoder.encode("/some/path/force-to-target-1/more/asdfasdf", "UTF-8");
			httpGet = new HttpGet("http://localhost:3333/crs/consistenthash/deliveryservice?ip=98.76.54.123&deliveryServiceId=" + steeringDeliveryServiceId + "&requestPath=" + requestPath);
			response = closeableHttpClient.execute(httpGet);

			deliveryServiceNode = objectMapper.readTree(EntityUtils.toString(response.getEntity()));
			assertThat(deliveryServiceNode.get("id").asText(), equalTo("steering-target-1"));
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itUsesRegexToStandardizeRequestPath() throws Exception {
		CloseableHttpResponse response = null;

		try {
			String requestPath = URLEncoder.encode("/some/path/thing.m3u8", "UTF-8");
			String encodedConsistentHashRegex = URLEncoder.encode(consistentHashRegex, "UTF-8");
			HttpGet httpGet = new HttpGet("http://localhost:3333/crs/consistenthash/patternbased/regex?regex=" + encodedConsistentHashRegex + "&requestPath=" + requestPath);

			response = closeableHttpClient.execute(httpGet);

			assertThat("Expected to get 200 response from /consistenthash/patternbased/regex endpoint", response.getStatusLine().getStatusCode(), equalTo(200));

			ObjectMapper objectMapper = new ObjectMapper(new JsonFactory());
			JsonNode resp = objectMapper.readTree(EntityUtils.toString(response.getEntity()));
			String resultingPathToConsistentHash = resp.get("resultingPathToConsistentHash").asText();

			requestPath = URLEncoder.encode("/other/path/other_thing.m3u8", "UTF-8");
			httpGet = new HttpGet("http://localhost:3333/crs/consistenthash/patternbased/regex?regex=" + encodedConsistentHashRegex + "&requestPath=" + requestPath);

			response = closeableHttpClient.execute(httpGet);

			resp = objectMapper.readTree(EntityUtils.toString(response.getEntity()));

			assertThat(JsonUtils.optString(resp, "resultingPathToConsistentHash"),equalTo(resultingPathToConsistentHash));
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itAppliesPatternBasedConsistentHashingToRequestsForCoverageZone() throws Exception {
		CloseableHttpResponse response = null;

		try {
			String requestPath = URLEncoder.encode("/some/path/thing.m3u8", "UTF-8");
			HttpGet httpGet = new HttpGet("http://localhost:3333/crs/consistenthash/cache/coveragezone?ip=" + ipAddressInCoverageZone + "&deliveryServiceId=" + deliveryServiceId + "&requestPath=" + requestPath);

			response = closeableHttpClient.execute(httpGet);

			assertThat("Expected to find " + ipAddressInCoverageZone + " in coverage zone using delivery service id " + deliveryServiceId, response.getStatusLine().getStatusCode(), equalTo(200));

			ObjectMapper objectMapper = new ObjectMapper(new JsonFactory());
			JsonNode cacheNode = objectMapper.readTree(EntityUtils.toString(response.getEntity()));

			String cacheId = cacheNode.get("id").asText();
			assertThat(cacheId, not(equalTo("")));

			response.close();

			response = closeableHttpClient.execute(httpGet);
			cacheNode = objectMapper.readTree(EntityUtils.toString(response.getEntity()));
			assertThat(cacheNode.get("id").asText(), equalTo(cacheId));

			response.close();

			requestPath = URLEncoder.encode("/other/path/other_thing.m3u8", "UTF-8");
			httpGet = new HttpGet("http://localhost:3333/crs/consistenthash/cache/coveragezone?ip=" + ipAddressInCoverageZone + "&deliveryServiceId=" + deliveryServiceId + "&requestPath=" + requestPath);

			response = closeableHttpClient.execute(httpGet);
			cacheNode = objectMapper.readTree(EntityUtils.toString(response.getEntity()));
			assertThat(cacheNode.get("id").asText(), equalTo(cacheId));
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itAppliesPatternBasedConsistentHashingForRequestsOutsideCoverageZone() throws Exception {
		CloseableHttpResponse response = null;

		try {
			String requestPath = URLEncoder.encode("/some/path/thing.m3u8", "UTF-8");
			HttpGet httpGet = new HttpGet("http://localhost:3333/crs/consistenthash/cache/geolocation?ip=8.8.8.8&deliveryServiceId=" + deliveryServiceId + "&requestPath=" + requestPath);

			response = closeableHttpClient.execute(httpGet);

			assertThat(response.getStatusLine().getStatusCode(), equalTo(200));

			ObjectMapper objectMapper = new ObjectMapper(new JsonFactory());
			JsonNode cacheNode = objectMapper.readTree(EntityUtils.toString(response.getEntity()));

			String cacheId = cacheNode.get("id").asText();
			assertThat(cacheId, not(equalTo("")));

			response.close();

			response = closeableHttpClient.execute(httpGet);
			cacheNode = objectMapper.readTree(EntityUtils.toString(response.getEntity()));
			assertThat(cacheNode.get("id").asText(), equalTo(cacheId));

			response.close();

			requestPath = URLEncoder.encode("/other/path/other_thing.m3u8", "UTF-8");
			httpGet = new HttpGet("http://localhost:3333/crs/consistenthash/cache/geolocation?ip=8.8.8.8&deliveryServiceId=" + deliveryServiceId + "&requestPath=" + requestPath);

			response = closeableHttpClient.execute(httpGet);
			cacheNode = objectMapper.readTree(EntityUtils.toString(response.getEntity()));
			assertThat(cacheNode.get("id").asText(), equalTo(cacheId));
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itAppliesPatternBasedConsistentHashingToSteeringDeliveryService() throws Exception {
		CloseableHttpResponse response = null;
		try {
			String requestPath = URLEncoder.encode("/some/path/thing.m3u8", "UTF-8");
			HttpGet httpGet = new HttpGet("http://localhost:3333/crs/consistenthash/deliveryservice?ip=98.76.54.123&deliveryServiceId=" + steeringDeliveryServiceId + "&requestPath=" + requestPath);
			response = closeableHttpClient.execute(httpGet);

			ObjectMapper objectMapper = new ObjectMapper(new JsonFactory());
			JsonNode deliveryServiceNode = objectMapper.readTree(EntityUtils.toString(response.getEntity()));
			String deliveryServiceId = deliveryServiceNode.get("id").asText();
			assertThat(deliveryServiceId, isIn(steeredDeliveryServices));

			response.close();

			requestPath = URLEncoder.encode("/other_different_path_12344321/path/other_thing_to_hash_differently.m3u8", "UTF-8");
			httpGet = new HttpGet("http://localhost:3333/crs/consistenthash/deliveryservice?ip=98.76.54.123&deliveryServiceId=" + steeringDeliveryServiceId + "&requestPath=" + requestPath);
			response = closeableHttpClient.execute(httpGet);

			deliveryServiceNode = objectMapper.readTree(EntityUtils.toString(response.getEntity()));
			assertThat(deliveryServiceNode.get("id").asText(), equalTo(deliveryServiceId));

		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itAppliesPatternBasedConsistentHashingToSteeringRequestsForCoverageZone() throws Exception {
		CloseableHttpResponse response = null;

		try {
			String requestPath = URLEncoder.encode("/some/path/thing.m3u8", "UTF-8");
			HttpGet httpGet = new HttpGet("http://localhost:3333/crs/consistenthash/cache/coveragezone/steering?ip=" + ipAddressInCoverageZone + "&deliveryServiceId=" + steeringDeliveryServiceId + "&requestPath=" + requestPath);

			response = closeableHttpClient.execute(httpGet);

			assertThat("Expected to find " + ipAddressInCoverageZone + " in coverage zone using delivery service id " + deliveryServiceId, response.getStatusLine().getStatusCode(), equalTo(200));

			ObjectMapper objectMapper = new ObjectMapper(new JsonFactory());
			JsonNode cacheNode = objectMapper.readTree(EntityUtils.toString(response.getEntity()));

			String cacheId = cacheNode.get("id").asText();
			assertThat(cacheId, not(equalTo("")));

			response.close();

			response = closeableHttpClient.execute(httpGet);
			cacheNode = objectMapper.readTree(EntityUtils.toString(response.getEntity()));
			assertThat(cacheNode.get("id").asText(), equalTo(cacheId));

			response.close();

			requestPath = URLEncoder.encode("/other_different_path_12344321/path/other_thing_to_hash_differently.m3u8", "UTF-8");
			httpGet = new HttpGet("http://localhost:3333/crs/consistenthash/cache/coveragezone/steering?ip=" + ipAddressInCoverageZone + "&deliveryServiceId=" + steeringDeliveryServiceId + "&requestPath=" + requestPath);

			response = closeableHttpClient.execute(httpGet);
			cacheNode = objectMapper.readTree(EntityUtils.toString(response.getEntity()));
			assertThat(cacheNode.get("id").asText(), equalTo(cacheId));
		} finally {
			if (response != null) response.close();
		}
	}
}
