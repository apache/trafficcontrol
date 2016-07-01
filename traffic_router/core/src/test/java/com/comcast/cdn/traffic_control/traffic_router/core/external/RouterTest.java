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

package com.comcast.cdn.traffic_control.traffic_router.core.external;

import com.comcast.cdn.traffic_control.traffic_router.core.util.ExternalTest;
import com.fasterxml.jackson.core.JsonFactory;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.http.Header;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.conn.ssl.SSLConnectionSocketFactory;
import org.apache.http.conn.ssl.TrustSelfSignedStrategy;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClientBuilder;
import org.apache.http.ssl.SSLContextBuilder;
import org.apache.http.util.EntityUtils;
import org.hamcrest.Matchers;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;
import org.junit.experimental.categories.Category;

import javax.net.ssl.HostnameVerifier;
import javax.net.ssl.SNIHostName;
import javax.net.ssl.SNIServerName;
import javax.net.ssl.SSLParameters;
import javax.net.ssl.SSLSession;
import javax.net.ssl.SSLSocket;
import java.io.IOException;
import java.io.InputStream;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;

import static org.hamcrest.CoreMatchers.containsString;
import static org.hamcrest.CoreMatchers.nullValue;
import static org.hamcrest.CoreMatchers.startsWith;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.isIn;
import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.core.IsNot.not;
import static org.junit.Assert.fail;

@Category(ExternalTest.class)
public class RouterTest {
	private CloseableHttpClient httpClient;
	private String deliveryServiceId;
	private List<String> validLocations = new ArrayList<>();
	private String deliveryServiceDomain;
	private String secureDeliveryServiceId;
	private List<String> secureValidLocations = new ArrayList<>();
	private String secureDeliveryServiceDomain;

	@Before
	public void before() throws Exception {
		ObjectMapper objectMapper = new ObjectMapper(new JsonFactory());

		String resourcePath = "internal/api/1.2/steering.json";
		InputStream inputStream = getClass().getClassLoader().getResourceAsStream(resourcePath);

		if (inputStream == null) {
			fail("Could not find file '" + resourcePath + "' needed for test from the current classpath as a resource!");
		}

		JsonNode steeringNode = objectMapper.readTree(inputStream).get("response").get(0);

		String steeringDeliveryServiceId = steeringNode.get("deliveryService").asText();

		resourcePath = "publish/CrConfig.json";
		inputStream = getClass().getClassLoader().getResourceAsStream(resourcePath);
		if (inputStream == null) {
			fail("Could not find file '" + resourcePath + "' needed for test from the current classpath as a resource!");
		}

		JsonNode jsonNode = objectMapper.readTree(inputStream);

		deliveryServiceId = null;

		Iterator<String> deliveryServices = jsonNode.get("deliveryServices").fieldNames();
		while (deliveryServices.hasNext()) {
			String dsId = deliveryServices.next();

			if (dsId.equals(steeringDeliveryServiceId)) {
				continue;
			}

			JsonNode deliveryServiceNode = jsonNode.get("deliveryServices").get(dsId);
			Iterator<JsonNode> matchsets = deliveryServiceNode.get("matchsets").iterator();
			while (matchsets.hasNext() && (deliveryServiceId == null || secureDeliveryServiceId == null)) {
				if ("HTTP".equals(matchsets.next().get("protocol").asText())) {
					if (deliveryServiceNode.get("sslEnabled").asBoolean(false)) {
						secureDeliveryServiceId = dsId;
						secureDeliveryServiceDomain = deliveryServiceNode.get("domains").get(0).asText();
					} else {
						deliveryServiceId = dsId;
						deliveryServiceDomain = deliveryServiceNode.get("domains").get(0).asText();
					}
				}
			}
		}

		assertThat(deliveryServiceId, not(nullValue()));
		assertThat(deliveryServiceDomain, not(nullValue()));
		assertThat(secureDeliveryServiceId, not(nullValue()));
		assertThat(secureDeliveryServiceDomain, not(nullValue()));

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
				validLocations.add("http://" + cacheId + "." + deliveryServiceDomain + portText + "/stuff?fakeClientIpAddress=12.34.56.78");
			}

			if (cacheNode.get("deliveryServices").has(secureDeliveryServiceId)) {
				int port = cacheNode.has("httpsPort") ? cacheNode.get("httpsPort").asInt(443) : 443;

				String portText = (port == 443) ? "" : ":" + port;
				secureValidLocations.add("https://" + cacheId + "." + secureDeliveryServiceDomain + portText + "/stuff?fakeClientIpAddress=12.34.56.78");
			}
		}

		assertThat(validLocations.isEmpty(), equalTo(false));
		assertThat(secureValidLocations.isEmpty(), equalTo(false));

		httpClient = HttpClientBuilder.create()
			.setSSLSocketFactory(new ClientSslSocketFactory("tr.https-test.thecdn.example.com"))
			.setSSLHostnameVerifier(new TestHostnameVerifier())
			.disableRedirectHandling()
			.build();
	}

	@After
	public void after() throws IOException {
	 	httpClient.close();
	}

	@Test
	public void itHasAHomePage() throws IOException {
		HttpGet httpGet = new HttpGet("http://localhost:8888/index.html");

		CloseableHttpResponse response = null;
		try {
			response = httpClient.execute(httpGet);
			assertThat(EntityUtils.toString(response.getEntity()), containsString("This is a test!"));
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itRedirectsValidHttpRequests() throws IOException, InterruptedException {
		HttpGet httpGet = new HttpGet("http://localhost:8888/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "tr." + deliveryServiceId + ".bar");
		CloseableHttpResponse response = null;

		try {
			response = httpClient.execute(httpGet);
			assertThat(response.getStatusLine().getStatusCode(), equalTo(302));
			Header header = response.getFirstHeader("Location");
			assertThat(header.getValue(), isIn(validLocations));
			assertThat(header.getValue(), Matchers.startsWith("http://"));
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itDoesRoutingThroughPathsStartingWithCrs() throws Exception {
		HttpGet httpGet = new HttpGet("http://localhost:8888/crs/stats?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "foo." + deliveryServiceId + ".bar");
		CloseableHttpResponse response = null;

		try {
			response = httpClient.execute(httpGet);
			assertThat(response.getStatusLine().getStatusCode(), equalTo(302));
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itConsistentlyRedirectsValidRequests() throws IOException, InterruptedException {
		HttpGet httpGet = new HttpGet("http://localhost:8888/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "tr." + deliveryServiceId + ".bar");
		CloseableHttpResponse response = null;

		try {
			response = httpClient.execute(httpGet);
			assertThat(response.getStatusLine().getStatusCode(), equalTo(302));
			String location = response.getFirstHeader("Location").getValue();

			response.close();

			for (int i = 0; i < 100; i++) {
				response = httpClient.execute(httpGet);
				assertThat(response.getFirstHeader("Location").getValue(), equalTo(location));
				response.close();
			}
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itRejectsInvalidRequests() throws IOException {
		HttpGet httpGet = new HttpGet("http://localhost:8888/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "foo.invalid-delivery-service-id.bar");
		CloseableHttpResponse response = null;

		try {
			response = httpClient.execute(httpGet);
			assertThat(response.getStatusLine().getStatusCode(), equalTo(503));
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itRedirectsHttpsRequests() throws Exception {
		HttpGet httpGet = new HttpGet("https://localhost:8443/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "tr." + secureDeliveryServiceId + ".thecdn.example.com");
		CloseableHttpResponse response = null;

		try {
			response = httpClient.execute(httpGet);
			assertThat(response.getStatusLine().getStatusCode(), equalTo(302));
			Header header = response.getFirstHeader("Location");
			assertThat(header.getValue(), isIn(secureValidLocations));
			assertThat(header.getValue(), startsWith("https://"));
			assertThat(header.getValue(), containsString(secureDeliveryServiceId + ".thecdn.example.com/stuff"));
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itRedirectsFromHttpToHttps() throws Exception {
		HttpGet httpGet = new HttpGet("http://localhost:8888/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "tr." + secureDeliveryServiceId + ".bar");
		CloseableHttpResponse response = null;

		try {
			response = httpClient.execute(httpGet);
			assertThat(response.getStatusLine().getStatusCode(), equalTo(302));
			Header header = response.getFirstHeader("Location");
			assertThat(header.getValue(), isIn(secureValidLocations));
			assertThat(header.getValue(), startsWith("https://"));
			assertThat(header.getValue(), containsString(secureDeliveryServiceId + ".thecdn.example.com/stuff"));
		} finally {
			if (response != null) response.close();
		}
	}

	@Test
	public void itRejectsHttpsRequestsForHttpDeliveryService() throws Exception {
		HttpGet httpGet = new HttpGet("https://localhost:8443/stuff?fakeClientIpAddress=12.34.56.78");
		httpGet.addHeader("Host", "tr." + deliveryServiceId + ".bar");
		CloseableHttpResponse response = null;

		try {
			response = httpClient.execute(httpGet);
			assertThat(response.getStatusLine().getStatusCode(), equalTo(503));
		} finally {
			if (response != null) response.close();
		}
	}

	// This is a workaround to get HttpClient to do the equivalent of
	// curl -v --resolve 'tr.https-test.thecdn.cdnlab.example.com:8443:127.0.0.1' https://tr.https-test.thecdn.example.com:8443/foo.json
	class ClientSslSocketFactory extends SSLConnectionSocketFactory {
		private final String host;

		public ClientSslSocketFactory(String host) throws Exception {
			super(SSLContextBuilder.create().loadTrustMaterial(null, new TrustSelfSignedStrategy()).build(),
				new TestHostnameVerifier());
			this.host = host;
		}

		protected void prepareSocket(final SSLSocket sslSocket) throws IOException {
			SNIHostName serverName = new SNIHostName(host);
			List<SNIServerName> serverNames = new ArrayList<>(1);
			serverNames.add(serverName);

			SSLParameters params = sslSocket.getSSLParameters();
			params.setServerNames(serverNames);
			sslSocket.setSSLParameters(params);
		}
	}

	// This is a workaround for the same reason as above
	// org.apache.http.conn.ssl.SSLConnectionSocketFactory.verifyHostname(<socket>, 'localhost') normally fails
	class TestHostnameVerifier implements HostnameVerifier {
		@Override
		public boolean verify(String s, SSLSession sslSession) {
			return true;
		}
	}
}
