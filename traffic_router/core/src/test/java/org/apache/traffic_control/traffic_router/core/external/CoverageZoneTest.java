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

import org.apache.traffic_control.traffic_router.core.util.ExternalTest;
import com.fasterxml.jackson.core.JsonFactory;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.catalina.LifecycleException;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClientBuilder;
import org.apache.http.util.EntityUtils;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;
import org.junit.experimental.categories.Category;
import org.powermock.core.classloader.annotations.PowerMockIgnore;

import static org.hamcrest.Matchers.endsWith;
import static org.hamcrest.Matchers.greaterThan;
import static org.hamcrest.Matchers.nullValue;
import static org.hamcrest.Matchers.startsWith;
import static org.hamcrest.core.AnyOf.anyOf;
import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.core.IsNot.not;
import static org.junit.Assert.assertThat;

@PowerMockIgnore("javax.management.*")
@Category(ExternalTest.class)
public class CoverageZoneTest {
	private CloseableHttpClient closeableHttpClient;

	@Before
	public void before() throws LifecycleException {
		closeableHttpClient = HttpClientBuilder.create().build();
	}

	@After
	public void after() throws Exception{
		if (closeableHttpClient != null) closeableHttpClient.close();
	}

	@Test
	public void itGetsCacheLocation() throws Exception {
		HttpGet httpGet = new HttpGet("http://localhost:3333/crs/coveragezone/cachelocation?ip=100.3.3.123&deliveryServiceId=steering-target-1");

		CloseableHttpResponse response = null;
		try {
			response = closeableHttpClient.execute(httpGet);
			String jsonString = EntityUtils.toString(response.getEntity());
			ObjectMapper objectMapper = new ObjectMapper(new JsonFactory());
			JsonNode jsonNode = objectMapper.readTree(jsonString);

			assertThat(jsonNode.get("id").asText(), equalTo("location-3"));
			assertThat(jsonNode.get("geolocation"), not(nullValue()));
			assertThat(jsonNode.get("caches").get(0).get("id").asText(), startsWith("edge-cache-03"));
			assertThat(jsonNode.get("caches").get(0).get("fqdn").asText(), startsWith("edge-cache-03"));
			assertThat(jsonNode.get("caches").get(0).get("fqdn").asText(), endsWith("thecdn.example.com"));
			assertThat(jsonNode.get("caches").get(0).get("port").asInt(), greaterThan(1024));
			assertThat(jsonNode.get("caches").get(0).get("hashValues").get(0).asDouble(), greaterThan(1.0));
			assertThat(isValidIpV4String(jsonNode.get("caches").get(0).get("ip4").asText()), equalTo(true));
			assertThat(jsonNode.get("caches").get(0).get("ip6").asText(), not(equalTo("")));
			assertThat(jsonNode.get("caches").get(0).has("available"), equalTo(true));
			assertThat(jsonNode.has("properties"), equalTo(true));
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itGetsCaches() throws Exception {
		HttpGet httpGet = new HttpGet("http://localhost:3333/crs/coveragezone/caches?deliveryServiceId=steering-target-4&cacheLocationId=location-3");

		CloseableHttpResponse response = null;
		try {
			response = closeableHttpClient.execute(httpGet);
			assertThat(response.getStatusLine().getStatusCode(), equalTo(200));

			ObjectMapper objectMapper = new ObjectMapper(new JsonFactory());
			JsonNode jsonNode = objectMapper.readTree(EntityUtils.toString(response.getEntity()));

			assertThat(jsonNode.isArray(), equalTo(true));
			JsonNode cacheNode = jsonNode.get(0);
			assertThat(cacheNode.get("id").asText(), not(nullValue()));
			assertThat(cacheNode.get("fqdn").asText(), not(nullValue()));
			assertThat(cacheNode.get("ip4").asText(), not(nullValue()));
			assertThat(cacheNode.get("ip6").asText(), not(nullValue()));
			// If the value is null or otherwise not an int we'll get back -123456, so any other value returned means success
			assertThat(cacheNode.get("port").asInt(-123456), not(equalTo(-123456)));
			assertThat(cacheNode.get("deliveryServices").isArray(), equalTo(true));
			assertThat(cacheNode.get("hashValues").get(0).asDouble(-1024.1024), not(equalTo(-1024.1024)));
			assertThat(cacheNode.get("available").asText(), anyOf(equalTo("true"), equalTo("false")));
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itReturns404ForMissingDeliveryService() throws Exception {
		HttpGet httpGet = new HttpGet("http://localhost:3333/crs/coveragezone/caches?deliveryServiceId=ds-07&cacheLocationId=location-5");

		try (CloseableHttpResponse response = closeableHttpClient.execute(httpGet)) {
			assertThat(response.getStatusLine().getStatusCode(), equalTo(404));
		}
	}

	boolean isValidIpV4String(String ip) {
		String[] octets = ip.split("\\.");

		if (octets.length != 4) {
			return false;
		}

		for (String octet : octets) {
			int b = Integer.parseInt(octet);
			if (b < 0 || 255 < b) {
				return false;
			}
		}

		return true;
	}
}

