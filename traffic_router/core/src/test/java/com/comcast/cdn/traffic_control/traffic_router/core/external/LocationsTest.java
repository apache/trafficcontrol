package com.comcast.cdn.traffic_control.traffic_router.core.external;

import com.comcast.cdn.traffic_control.traffic_router.core.util.ExternalTest;
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
import static org.hamcrest.core.IsEqual.equalTo;

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

			String expected = "{" +
				"\"locations\":[" +
					"\"us-bb-central\"," +
					"\"us-bb-east\"," +
					"\"us-bb-west\"," +
					"\"us-ca-sanjose\"," +
					"\"us-co-denver\"," +
					"\"us-de-newcastle\"," +
					"\"us-fl-sarasota\"," +
					"\"us-ga-atlanta\"," +
					"\"us-il-chicago\"," +
					"\"us-ma-woburn\"," +
					"\"us-md-baltimore\"," +
					"\"us-mi-grand_rapids\"," +
					"\"us-mn-roseville\"," +
					"\"us-nj-plainfield\"," +
					"\"us-pa-pittsburgh\"," +
					"\"us-tx-houston\"," +
					"\"us-va-richmond\"," +
					"\"us-wa-seattle\"" +
				"]" +
			"}";

			assertThat(EntityUtils.toString(response.getEntity()), equalTo(expected));
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
			
			String expected = "{" +
				"\"cacheId\":\"odol-atsec-sim-030\"," +
				"\"fqdn\":\"odol-atsec-sim-030.jenkins-sim.cdnlab.comcast.net\"," +
				"\"ipAddresses\":[\"192.168.8.32\",\"2001:558:fee8:168:c:0:0:1e\"]," +
				"\"port\":0," +
				"\"adminStatus\":null," +
				"\"lastUpdateHealthy\":false," +
				"\"lastUpdateTime\":0," +
				"\"connections\":0," +
				"\"currentBW\":0," +
				"\"availBW\":0," +
				"\"cacheOnline\":true" +
				"}";
			
			assertThat(EntityUtils.toString(response.getEntity()), containsString(expected));
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itGetsCachesForALocation() throws Exception {
		HttpGet httpGet = new HttpGet("http://localhost:3333/crs/locations/us-co-denver/caches");

		CloseableHttpResponse response = null;
		try {
			response = closeableHttpClient.execute(httpGet);

			String expected = "{" +
				"\"cacheId\":\"odol-atsec-sim-119\"," +
				"\"fqdn\":\"odol-atsec-sim-119.jenkins-sim.cdnlab.comcast.net\"," +
				"\"ipAddresses\":[\"192.168.8.121\",\"2001:558:fee8:168:c:0:0:77\"]," +
				"\"port\":0," +
				"\"adminStatus\":null," +
				"\"lastUpdateHealthy\":false," +
				"\"lastUpdateTime\":0," +
				"\"connections\":0," +
				"\"currentBW\":0," +
				"\"availBW\":0," +
				"\"cacheOnline\":true" +
			"}";

			assertThat(EntityUtils.toString(response.getEntity()), containsString(expected));

		} finally {
			if (response != null) response.close();
		}
	}
}
