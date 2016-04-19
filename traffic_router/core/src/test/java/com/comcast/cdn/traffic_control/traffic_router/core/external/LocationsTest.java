package com.comcast.cdn.traffic_control.traffic_router.core.external;

import com.comcast.cdn.traffic_control.traffic_router.core.util.ExternalTest;
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

import static org.hamcrest.CoreMatchers.containsString;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.core.AnyOf.anyOf;
import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.core.IsNot.not;

@Category(ExternalTest.class)
public class LocationsTest {
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
	public void itGetsAListOfLocations() throws Exception {
		HttpGet httpGet = new HttpGet("http://localhost:3333/crs/locations");

		CloseableHttpResponse response = null;
		try {
			response = closeableHttpClient.execute(httpGet);
			assertThat(response.getStatusLine().getStatusCode(), equalTo(200));

			ObjectMapper objectMapper = new ObjectMapper(new JsonFactory());
			JsonNode jsonNode = objectMapper.readTree(EntityUtils.toString(response.getEntity()));
			assertThat(jsonNode.get("locations").isArray(), equalTo(true));

		} finally {
			if (response != null) response.close();
		}
	}
	
	@Test
	public void itGetsAListOfCaches() throws Exception {
		HttpGet httpGet = new HttpGet("http://localhost:3333/crs/locations/caches");
		
		CloseableHttpResponse response = null;
		try {
			response = closeableHttpClient.execute(httpGet);
			
			ObjectMapper objectMapper = new ObjectMapper(new JsonFactory());
			JsonNode jsonNode = objectMapper.readTree(EntityUtils.toString(response.getEntity()));
			String locationName = jsonNode.get("locations").fieldNames().next();
			JsonNode cacheNode = jsonNode.get("locations").get(locationName).get(0);

			assertThat(cacheNode.get("cacheId").asText(), not(equalTo("")));
			assertThat(cacheNode.get("fqdn").asText(), not(equalTo("")));

			assertThat(cacheNode.get("ipAddresses").isArray(), equalTo(true));
			assertThat(cacheNode.has("adminStatus"), equalTo(true));

			assertThat(cacheNode.get("port").asInt(-123456), not(equalTo(-123456)));
			assertThat(cacheNode.get("lastUpdateTime").asInt(-123456), not(equalTo(-123456)));
			assertThat(cacheNode.get("connections").asInt(-123456), not(equalTo(-123456)));
			assertThat(cacheNode.get("currentBW").asInt(-123456), not(equalTo(-123456)));
			assertThat(cacheNode.get("availBW").asInt(-123456), not(equalTo(-123456)));

			assertThat(cacheNode.get("cacheOnline").asText(), anyOf(equalTo("true"), equalTo("false")) );
			assertThat(cacheNode.get("lastUpdateHealthy").asText(), anyOf(equalTo("true"), equalTo("false")) );

		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itGetsCachesForALocation() throws Exception {
		HttpGet httpGet = new HttpGet("http://localhost:3333/crs/locations");

		CloseableHttpResponse response = null;
		try {
			response = closeableHttpClient.execute(httpGet);

			ObjectMapper objectMapper = new ObjectMapper(new JsonFactory());
			JsonNode jsonNode = objectMapper.readTree(EntityUtils.toString(response.getEntity()));

			String location = jsonNode.get("locations").get(0).asText();

			httpGet = new HttpGet("http://localhost:3333/crs/locations/" + location + "/caches");

			response = closeableHttpClient.execute(httpGet);

			jsonNode = objectMapper.readTree(EntityUtils.toString(response.getEntity()));

			assertThat(jsonNode.get("caches").isArray(), equalTo(true));
			JsonNode cacheNode = jsonNode.get("caches").get(0);

			assertThat(cacheNode.get("cacheId").asText(), not(equalTo("")));
			assertThat(cacheNode.get("fqdn").asText(), not(equalTo("")));

			assertThat(cacheNode.get("ipAddresses").isArray(), equalTo(true));
			assertThat(cacheNode.has("adminStatus"), equalTo(true));

			assertThat(cacheNode.get("port").asInt(-123456), not(equalTo(-123456)));
			assertThat(cacheNode.get("lastUpdateTime").asInt(-123456), not(equalTo(-123456)));
			assertThat(cacheNode.get("connections").asInt(-123456), not(equalTo(-123456)));
			assertThat(cacheNode.get("currentBW").asInt(-123456), not(equalTo(-123456)));
			assertThat(cacheNode.get("availBW").asInt(-123456), not(equalTo(-123456)));

			assertThat(cacheNode.get("cacheOnline").asText(), anyOf(equalTo("true"), equalTo("false")) );
			assertThat(cacheNode.get("lastUpdateHealthy").asText(), anyOf(equalTo("true"), equalTo("false")) );

		} finally {
			if (response != null) response.close();
		}
	}
}
